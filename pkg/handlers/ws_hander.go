package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/pbogut/hackdeck/pkg/types"
)

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

	for _, b := range config.Buttons {
		buttons.AddButton(types.NewButton(b.Row, b.Column).SetColor(b.Color).SetIconFromPath(b.Icon))
	}

	return buttons
}

func handleButtonPress(row, col int) types.Buttons {
	fmt.Printf("Button pressed at row: %d, col: %d\n", row, col)

	return types.NewUpdateButton()
}

func handleButtonLongPress(row, col int) types.Buttons {
	fmt.Printf("Button long pressed at row: %d, col: %d\n", row, col)

	return types.NewUpdateButton()
}

func handleButtonRelease(row, col int) types.Buttons {
	fmt.Printf("Button released at row: %d, col: %d\n", row, col)

	return types.NewUpdateButton()
}

func handleButtonLongPressRelease(row, col int) types.Buttons {
	fmt.Printf("Button long press released at row: %d, col: %d\n", row, col)

	return types.NewUpdateButton()
}

func msgToRowCol(msg []byte) (int, int) {
	var action types.ClickAction
	json.Unmarshal(msg, &action)
	return action.GetXY()
}

// WebSocket handler
func WsHandler(w http.ResponseWriter, r *http.Request) {
	config = types.ReadConfig()

	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error while upgrading connection:", err)
		return
	}
	defer conn.Close()
	fmt.Println("Client connected")
	// Echo messages back to the client
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
			result = handleButtonLongPress(row, col)
		case "BUTTON_LONG_PRESS_RELEASE":
			row, col := msgToRowCol(msg)
			result = handleButtonLongPressRelease(row, col)
		case "BUTTON_PRESS":
			row, col := msgToRowCol(msg)
			result = handleButtonPress(row, col)
		case "BUTTON_RELEASE":
			orw, col := msgToRowCol(msg)
			result = handleButtonRelease(orw, col)
		case "GET_BUTTONS":
			result = handleGetButtons()
		case "CONNECTED":
			result = handleConnected(msg)
		}

		response, err := json.Marshal(result)
		if err != nil {
			break
		}
		conn.WriteMessage(messageType, response)

		// Print the received message
		fmt.Printf("Received: %s\n", msg)
	}
}
