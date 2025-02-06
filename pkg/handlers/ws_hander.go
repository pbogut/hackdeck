package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/pbogut/hackdeck/pkg/logger"
	"github.com/pbogut/hackdeck/pkg/types"
)

var state types.State

var config types.Config

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections (for development purposes)
	},
}

func getConfigResponse() []byte {
	getConfig := types.GetConfig{
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

	response, err := json.Marshal(getConfig)
	if err != nil {
		logger.Fatalf("Error while marshaling buttons: %s", err)
	}

	return response
}

func getButtonsResponse() []byte {
	buttons := types.NewGetButtons()

	for _, btn := range state.GetButtons() {
		buttons.AddButton(*btn)
	}

	response, err := json.Marshal(buttons)
	if err != nil {
		logger.Fatalf("Error while marshaling buttons: %s", err)
	}

	return response
}

func handleButtonPress(row, col int) {
	logger.Debugf("Button pressed (row: %d, col: %d)", row, col)
	go execAction(row, col, types.BUTTON_PRESS)
}

func handleButtonLongPress(row, col int) {
	logger.Debugf("Button long pressed (row: %d, col: %d)", row, col)
	go execAction(row, col, types.BUTTON_LONG_PRESS)
}

func handleButtonRelease(row, col int) {
	logger.Debugf("Button released (row: %d, col: %d)", row, col)
	go execAction(row, col, types.BUTTON_RELEASE)
}

func handleButtonLongPressRelease(row, col int) {
	logger.Debugf("Button long press released (row: %d, col: %d)", row, col)
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
		logger.Error("Error while upgrading connection:", err)
		return
	}
	defer conn.Close()
	logger.Info("New client connected")
	// Echo messages back to the client

	RegisterClient(conn)

	for {
		// Read message from the client
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			logger.Error("Error while reading message:", err)
			break
		}

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
			conn.WriteMessage(messageType, getButtonsResponse())
		case "CONNECTED":
			conn.WriteMessage(messageType, getConfigResponse())
		}

		// Print the received message
		logger.Debugf("Message Received: %s", msg)
	}
}
