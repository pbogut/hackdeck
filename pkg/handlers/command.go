package handlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pbogut/hackdeck/pkg/types"
)

var clients []*websocket.Conn

func execCommand(row, col int, command string) {
	if command != "" {
		args := config.ShellArguments
		args = append(args, command)

		fmt.Println("Execute command:", command)
		fmt.Println("Execute:", config.ShellCommand, args)

		cmd := exec.Command(config.ShellCommand, args...)

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			fmt.Println("Error while getting stdout pipe:", err)
		}

		cmd.Start()

		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			m := scanner.Text()

			fmt.Println("Recievied response from command:", m)

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
	fmt.Println("Start Execute", state.GetButtonConfigs())
	for _, btnCfg := range state.GetButtonConfigs() {
		fmt.Println("Execute:", btnCfg.Execute)
		if btnCfg.Execute != "" {
			if btnCfg.Interval > 0 {
				go handleCommandInterval(btnCfg.Row, btnCfg.Column, btnCfg.Execute)
			} else {
				go execCommand(btnCfg.Row, btnCfg.Column, btnCfg.Execute)
			}
		}
	}
}

func execAction(row, col, status int) {
	command := state.GetCmd(row, col, status)
	execCommand(row, col, command)
}

func handleButtonPress(row, col int) {
	fmt.Printf("Button pressed at row: %d, col: %d\n", row, col)
	go execAction(row, col, types.BUTTON_PRESS)
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
		fmt.Println("Broadcasting message to client:", i)

		client.WriteMessage(websocket.TextMessage, msg)
	}
}
