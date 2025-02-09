package handlers

import (
	"bufio"
	"encoding/json"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pbogut/hackdeck/pkg/logger"
	"github.com/pbogut/hackdeck/pkg/types"
)

type CommandType int

const (
	MAIN_COMMAND = iota
	ACTION_COMMAND
)

var clients []*websocket.Conn
var commands []*exec.Cmd

var clientsMutex = &sync.Mutex{}
var commandsMutex = &sync.Mutex{}

func sendStdin(row, col int, command string) {
	btn := state.GetButton(row, col)
	pid := btn.GetPid()
	for _, cmd := range commands {
		if cmd != nil && cmd.Process.Pid == pid {
			logger.Debugf("Send stdin to command (row: %d, col: %d, pid: %d) %s", row, col, pid, command)
			pipe := state.GetPipe(row, col)
			if pipe == nil {
				logger.Errorf("Process not found (row: %d, col: %d, pid: %d)", row, col, pid)
			} else {
				(*pipe).Write([]byte(command))
			}
			break
		}
	}
}

func execCommand(row, col int, command string, cmdType CommandType) {
	if command != "" {
		btn := state.GetButton(row, col)
		pid := btn.GetPid()
		if cmdType == MAIN_COMMAND && pid > 0 {
			command = strings.ReplaceAll(command, "%pid%", strconv.Itoa(pid))
		}

		args := config.ShellArguments
		args = append(args, command)

		cmd := exec.Command(config.ShellCommand, args...)

		logger.Debugf("Execute command (row: %d, col: %d, cmd: %s)", row, col, cmd)

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			logger.Error("Error while getting stdout pipe:", err)
		}
		if cmdType == MAIN_COMMAND {
			stdin, err := cmd.StdinPipe()
			if err != nil {
				logger.Error("Error while getting stdin pipe:", err)
			}
			state.AddPipe(row, col, &stdin)
		}

		monitorCommand(cmd)
		cmd.Start()
		pid = cmd.Process.Pid
		logger.Debugf("Command started (row: %d, col: %d, pid: %d, cmd: %s)", row, col, pid, cmd)

		if cmdType == MAIN_COMMAND {
			btn.SetPid(pid)
		}

		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			m := scanner.Text()

			logger.Debugf("Recievied response (row: %d, col: %d, pid: %d) %s", row, col, pid, m)

			update := types.NewUpdateButton()
			btn := state.GetButton(row, col)

			var result map[string]any
			err := json.Unmarshal([]byte(m), &result)
			if err == nil {
				btn.UpdateFromAnyMap(result)

				if btn.IsChanged() {
					btn.ResetChanged()
					update.AddButton(*btn)
					response, err := json.Marshal(update)
					if err != nil {
						break
					}
					Broadcast(response)
				}
			}
		}
		cmd.Wait()
		if cmdType == MAIN_COMMAND {
			btn.SetPid(0)
		}
		releaseCommand(cmd)
	}
}

func handleCommandInterval(row, col int, command string, interval int) {
	ticker := time.NewTicker(time.Second * time.Duration(interval))
	for ; true; <-ticker.C {
		execCommand(row, col, command, MAIN_COMMAND)
	}
}

func startExecute() {
	logger.Debug("Start Execute commands")
	for _, btnCfg := range state.GetButtonConfigs() {
		if btnCfg.Execute != "" {
			if btnCfg.Interval > 0 {
				logger.Debugf("Interval (%d): %s", btnCfg.Interval, btnCfg.Execute)
				go handleCommandInterval(btnCfg.Row, btnCfg.Column, btnCfg.Execute, btnCfg.Interval)
			} else {
				logger.Debug("Execute:", btnCfg.Execute)
				go execCommand(btnCfg.Row, btnCfg.Column, btnCfg.Execute, MAIN_COMMAND)
			}
		}
	}
}

func execAction(row, col, status int) {
	command := state.GetCmd(row, col, status)
	command, pipe := strings.CutPrefix(command, "<|")
	if pipe {
		sendStdin(row, col, command)
	} else {
		execCommand(row, col, command, ACTION_COMMAND)
	}
}

func monitorCommand(cmd *exec.Cmd) {
	logger.Debug("Add command:", cmd)
	commandsMutex.Lock()
	commands = append(commands, cmd)
	commandsMutex.Unlock()
}

func releaseCommand(cmd *exec.Cmd) {
	commandsMutex.Lock()
	for i, c := range commands {
		if c == cmd {
			logger.Debug("Remove command:", cmd)
			commands[i] = commands[len(commands)-1]
			commands = commands[:len(commands)-1]
		}
	}
	commandsMutex.Unlock()
}

func killMonitoredCommands() {
	logger.Debugf("Kill all commands #%d", len(commands))
	for _, cmd := range commands {
		logger.Debug("Kill command:", cmd)
		cmd.Process.Kill()
	}
}

func Init() {
	config = types.ReadConfig()
	state.Init(config)
	commandsSize := config.Rows * config.Columns
	clients = make([]*websocket.Conn, 0)
	commands = make([]*exec.Cmd, commandsSize*2)

	startExecute()
}

func ReloadConfig() {
	killMonitoredCommands()
	config = types.ReadConfig()
	state.Init(config)

	Broadcast(getConfigResponse())
	Broadcast(getButtonsResponse())
}

func RegisterClient(client *websocket.Conn) {
	client.SetCloseHandler(func(code int, text string) error {
		clientsMutex.Lock()
		for i, c := range clients {
			if c == client {
				clients[i] = clients[len(clients)-1]
				clients = clients[:len(clients)-1]
			}
		}
		clientsMutex.Unlock()
		return nil
	})

	clientsMutex.Lock()
	clients = append(clients, client)
	clientsMutex.Unlock()
}

func Broadcast(msg []byte) {
	for i, client := range clients {
		logger.Debug("Broadcasting message to client:", i)
		client.WriteMessage(websocket.TextMessage, msg)
	}
}
