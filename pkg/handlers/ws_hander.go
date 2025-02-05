package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

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

func handleButtonLongPress(row, col int) {
	fmt.Printf("Button long pressed at row: %d, col: %d\n", row, col)
	go execAction(row, col, types.BUTTON_LONG_PRESS)
}

func handleButtonRelease(row, col int) {
	fmt.Printf("Button released at row: %d, col: %d\n", row, col)
	go execAction(row, col, types.BUTTON_RELEASE)
}

func handleButtonLongPressRelease(row, col int) {
	fmt.Printf("Button long press released at row: %d, col: %d\n", row, col)
	go execAction(row, col, types.BUTTON_LONG_PRESS_RELEASE)
}

func msgToRowCol(msg []byte) (int, int) {
	var action types.ClickAction
	json.Unmarshal(msg, &action)
	return action.GetXY()
}

// WebSocket handler
func WsHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error while upgrading connection:", err)
		return
	}
	defer conn.Close()
	fmt.Println("Client connected")
	// Echo messages back to the client

	RegisterClient(conn)

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
			handleButtonLongPress(row, col)
		case "BUTTON_LONG_PRESS_RELEASE":
			row, col := msgToRowCol(msg)
			handleButtonLongPressRelease(row, col)
		case "BUTTON_PRESS":
			row, col := msgToRowCol(msg)
			handleButtonPress(row, col)
		case "BUTTON_RELEASE":
			row, col := msgToRowCol(msg)
			handleButtonRelease(row, col)
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
