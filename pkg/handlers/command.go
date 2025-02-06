package handlers

import (
	"bufio"
	"encoding/json"
	"os/exec"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pbogut/hackdeck/pkg/logger"
	"github.com/pbogut/hackdeck/pkg/types"
)

var clients []*websocket.Conn

func execCommand(row, col int, command string) {
	if command != "" {
		args := config.ShellArguments
		args = append(args, command)

		cmd := exec.Command(config.ShellCommand, args...)

		logger.Debugf("Execute command (row: %d, col: %d, cmd: %s)", row, col, cmd)

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			logger.Error("Error while getting stdout pipe:", err)
		}

		cmd.Start()

		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			m := scanner.Text()

			logger.Debugf("Recievied response (row: %d, col: %d, cmd: %s)", row, col, cmd)
			logger.Debugf("Response: %s", m)

			update := types.NewUpdateButton()
			btn := state.GetButton(row, col)

			if strings.HasPrefix(m, "!COLOR!") {
				btn.SetColor(strings.TrimPrefix(m, "!COLOR!"))
			}
			if strings.HasPrefix(m, "!ICON_PATH!") {
				btn.SetIconFromPath(strings.TrimPrefix(m, "!ICON_PATH!"))
			}
			if strings.HasPrefix(m, "!ICON_TEXT!") {
				btn.SetIconFromText(strings.TrimPrefix(m, "!ICON_TEXT!"))
			}
			if strings.HasPrefix(m, "!ICON_COLOR!") {
				btn.SetIconColor(strings.TrimPrefix(m, "!ICON_COLOR!"))
			}
			if strings.HasPrefix(m, "!LABEL!") {
				btn.SetLabel(strings.TrimPrefix(m, "!LABEL!"))
			}

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
		cmd.Wait()
	}
}

func handleCommandInterval(row, col int, command string) {
	for range time.Tick(time.Second * 1) {
		execCommand(row, col, command)
	}
}

func startExecute() {
	logger.Debug("Start Execute commands")
	for _, btnCfg := range state.GetButtonConfigs() {
		if btnCfg.Execute != "" {
			if btnCfg.Interval > 0 {
				logger.Debugf("Interval (%d): %s", btnCfg.Interval, btnCfg.Execute)
				go handleCommandInterval(btnCfg.Row, btnCfg.Column, btnCfg.Execute)
			} else {
				logger.Debugf("Execute: %s", btnCfg.Execute)
				go execCommand(btnCfg.Row, btnCfg.Column, btnCfg.Execute)
			}
		}
	}
}

func execAction(row, col, status int) {
	command := state.GetCmd(row, col, status)
	execCommand(row, col, command)
}

func Init() {
	config = types.ReadConfig()
	state.Init(config)
	clients = make([]*websocket.Conn, 0)

	startExecute()
}

func RegisterClient(client *websocket.Conn) {
	client.SetCloseHandler(func(code int, text string) error {
		for i, c := range clients {
			if c == client {
				clients[i] = clients[len(clients)-1]
				clients = clients[:len(clients)-1]
			}
		}
		return nil
	})

	clients = append(clients, client)
}

func Broadcast(msg []byte) {
	for i, client := range clients {
		logger.Debugf("Broadcasting message to client: %d", i)
		client.WriteMessage(websocket.TextMessage, msg)
	}
}
