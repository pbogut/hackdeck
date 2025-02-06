package types

import (
	"encoding/base64"
	"fmt"
	"github.com/pbogut/hackdeck/pkg/label"
	"os"
	"strconv"
	"strings"
)

type Method struct {
	Method string `json:"Method"`
}

type ClickAction struct {
	Method  string `json:"Method"`
	Message string `json:"Message"`
}

func (c *ClickAction) GetXY() (int, int) {
	pos := strings.Split(c.Message, "_")
	x, err := strconv.Atoi(pos[0])
	if err != nil {
		return -1, -1
	}
	y, err := strconv.Atoi(pos[1])
	if err != nil {
		return -1, -1
	}

	return x, y
}

type Connected struct {
	// "Method": "CONNECTED",
	Method string `json:"Method"`
	// "Client-Id": "q9pttrzqq",
	ClientId string `json:"Client-Id"`
	// "API": 20,
	API int `json:"API"`
	// "Device-Type": "Web"
	DeviceType string `json:"Device-Type"`
}

type GetConfig struct {
	// "Method": "GET_CONFIG",
	Method string `json:"Method"`
	// "Rows": 3,
	Rows int `json:"Rows"`
	// "Columns": 5,
	Columns int `json:"Columns"`
	// "ButtonSpacing": 10,
	ButtonSpacing int `json:"ButtonSpacing"`
	// "ButtonRadius": 40,
	ButtonRadius int `json:"ButtonRadius"`
	// "ButtonBackground": true,
	ButtonBackground bool `json:"ButtonBackground"`
	// "Brightness": 0.3,
	Brightness float32 `json:"Brightness"`
	// "AutoConnect": false,
	AutoConnect bool `json:"AutoConnect"`
	// "WakeLock": "Connected",
	WakeLock string `json:"WakeLock"`
	// "SupportButtonReleaseLongPress": true
	SupportButtonReleaseLongPress bool `json:"SupportButtonReleaseLongPress"`
}

type Button struct {
	// "IconBase64": "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",
	IconBase64 string `json:"IconBase64"`
	// "Position_X": 0,
	PositionX int `json:"Position_X"`
	// "Position_Y": 0,
	PositionY int `json:"Position_Y"`
	// "LabelBase64": "",
	LabelBase64 string `json:"LabelBase64"`
	// "BackgroundColorHex": "#232323"
	BackgroundColorHex string `json:"BackgroundColorHex"`

	changed   bool   `json:"-"`
	iconPath  string `json:"-"`
	iconText  string `json:"-"`
	iconColor string `json:"-"`
	label     string `json:"-"`
}

type Buttons struct {
	// "Method": "GET_BUTTONS",
	Method string `json:"Method"`
	// "Buttons": []
	Buttons []Button `json:"Buttons"`
}

func (b *Buttons) AddButton(button Button) *Buttons {
	b.Buttons = append(b.Buttons, button)
	return b
}

func NewButton(row, col int) Button {
	return Button{
		PositionX:          col, // !sic
		PositionY:          row, // !sic
		IconBase64:         "",
		LabelBase64:        "",
		BackgroundColorHex: "#232323",

		iconPath: "",
		label:    "",
		changed:  false,
	}
}

func (b *Button) SetColor(color string) {
	if b.BackgroundColorHex != color {
		b.BackgroundColorHex = color
		b.changed = true
	}
}

func (b *Button) SetIconFromPath(path string) {
	if b.iconPath == path {
		return
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error while reading icon file:", err, path)
		return
	}

	b.changed = true
	b.iconPath = path
	b.iconText = ""
	b.IconBase64 = base64.StdEncoding.EncodeToString(bytes)
}

func (b *Button) SetIconFromText(text string) {
	if b.iconText == text {
		return
	}

	b.changed = true
	b.iconText = text
	b.iconPath = ""
	b.IconBase64 = label.GenerateIcon(text, b.iconColor)
}

func (b *Button) SetIconColor(color string) {
	if b.iconColor == color {
		return
	}
	b.iconColor = color

	if b.iconText == "" {
		return
	}

	b.changed = true
	b.IconBase64 = label.GenerateIcon(b.iconText, color)
}

func (b *Button) SetLabel(text string) {
	if b.label == text {
		return
	}

	b.changed = true
	b.label = text
	b.LabelBase64 = label.GenerateLabel(text)
}

func (b *Button) ResetChanged() {
	b.changed = false
}

func (b *Button) IsChanged() bool {
	return b.changed
}

func NewGetButtons() Buttons {
	return Buttons{
		Method:  "GET_BUTTONS",
		Buttons: []Button{},
	}
}

func NewUpdateButton() Buttons {
	return Buttons{
		Method:  "UPDATE_BUTTON",
		Buttons: []Button{},
	}
}
