package handlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pbogut/hackdeck/pkg/types"
)

var state types.State

var config types.Config

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections (for development purposes)
	},
}

func handleConnected(msg []byte) types.GetConfig {
	var connected types.Connected
	json.Unmarshal(msg, &connected)

	get_config := types.GetConfig{
		Method:                        "GET_CONFIG",
		Rows:                          config.Rows,
		Columns:                       config.Columns,
		ButtonSpacing:                 config.ButtonSpacing,
		ButtonRadius:                  config.ButtonRadius,
		ButtonBackground:              config.ButtonBackground,
		Brightness:                    config.Brightness,
		AutoConnect:                   false,
		WakeLock:                      "Connected",
		SupportButtonReleaseLongPress: config.SupportButtonReleaseLongPress,
	}

	return get_config
}

func handleGetButtons() types.Buttons {
	buttons := types.NewGetButtons()

	for _, btn := range state.GetButtons() {
		buttons.AddButton(*btn)
	}

	return buttons
}

func execCommand(conn *websocket.Conn, row, col int, command string) {
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
			if strings.HasPrefix(m, "!ICON!") {
				btn.SetIconFromPath(strings.TrimPrefix(m, "!ICON!"))
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
				conn.WriteMessage(websocket.TextMessage, response)
			}
		}
		cmd.Wait()
	}
}
func execAction(conn *websocket.Conn, row, col, status int) {
	command := state.GetCmd(row, col, status)
	execCommand(conn, row, col, command)
}

func handleButtonPress(conn *websocket.Conn, row, col int) {
	fmt.Printf("Button pressed at row: %d, col: %d\n", row, col)
	go execAction(conn, row, col, types.BUTTON_PRESS)
}

func handleButtonLongPress(conn *websocket.Conn, row, col int) {
	fmt.Printf("Button long pressed at row: %d, col: %d\n", row, col)
	go execAction(conn, row, col, types.BUTTON_LONG_PRESS)
}

func handleButtonRelease(conn *websocket.Conn, row, col int) {
	fmt.Printf("Button released at row: %d, col: %d\n", row, col)
	go execAction(conn, row, col, types.BUTTON_RELEASE)
}

func handleButtonLongPressRelease(conn *websocket.Conn, row, col int) {
	fmt.Printf("Button long press released at row: %d, col: %d\n", row, col)
	go execAction(conn, row, col, types.BUTTON_LONG_PRESS_RELEASE)
}

func msgToRowCol(msg []byte) (int, int) {
	var action types.ClickAction
	json.Unmarshal(msg, &action)
	return action.GetXY()
}

func startExecute(conn *websocket.Conn) {
	fmt.Println("Start Execute", state.GetButtonConfigs())
	for _, btnCfg := range state.GetButtonConfigs() {
		fmt.Println("Execute:", btnCfg.Execute)
		if btnCfg.Execute != "" {
			if btnCfg.Interval > 0 {
				go handleCommandInterval(conn, btnCfg.Row, btnCfg.Column, btnCfg.Execute)
			} else {
				go execCommand(conn, btnCfg.Row, btnCfg.Column, btnCfg.Execute)
			}
		}
	}
}

func handleCommandInterval(conn *websocket.Conn, row, col int, command string) {
	for range time.Tick(time.Second * 1) {
		execCommand(conn, row, col, command)
	}
}

// WebSocket handler
func WsHandler(w http.ResponseWriter, r *http.Request) {
	config = types.ReadConfig()
	state.Init(config)

	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error while upgrading connection:", err)
		return
	}
	defer conn.Close()
	fmt.Println("Client connected")
	// Echo messages back to the client

	startExecute(conn)

	for {
		// Read message from the client
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error while reading message:", err)
			break
		}

		var result interface{}
		var method types.Method
		json.Unmarshal(msg, &method)

		switch method.Method {
		case "BUTTON_LONG_PRESS":
			row, col := msgToRowCol(msg)
			handleButtonLongPress(conn, row, col)
		case "BUTTON_LONG_PRESS_RELEASE":
			row, col := msgToRowCol(msg)
			handleButtonLongPressRelease(conn, row, col)
		case "BUTTON_PRESS":
			row, col := msgToRowCol(msg)
			handleButtonPress(conn, row, col)
		case "BUTTON_RELEASE":
			orw, col := msgToRowCol(msg)
			handleButtonRelease(conn, orw, col)
		case "GET_BUTTONS":
			result = handleGetButtons()
		case "CONNECTED":
			result = handleConnected(msg)
		}

		if result != nil {
			response, err := json.Marshal(result)
			if err != nil {
				break
			}
			conn.WriteMessage(messageType, response)
		}

		// Print the received message
		fmt.Printf("Received: %s\n", msg)
	}
}
