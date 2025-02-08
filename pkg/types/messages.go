package types

import (
	"encoding/base64"
	"os"
	"strconv"
	"strings"

	"github.com/pbogut/hackdeck/pkg/label"
	"github.com/pbogut/hackdeck/pkg/logger"
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
	Method     string `json:"Method"`
	ClientId   string `json:"Client-Id"`
	API        int    `json:"API"`
	DeviceType string `json:"Device-Type"`
}

type GetConfig struct {
	Method                        string  `json:"Method"`
	Rows                          int     `json:"Rows"`
	Columns                       int     `json:"Columns"`
	ButtonSpacing                 int     `json:"ButtonSpacing"`
	ButtonRadius                  int     `json:"ButtonRadius"`
	ButtonBackground              bool    `json:"ButtonBackground"`
	Brightness                    float32 `json:"Brightness"`
	AutoConnect                   bool    `json:"AutoConnect"`
	WakeLock                      string  `json:"WakeLock"`
	SupportButtonReleaseLongPress bool    `json:"SupportButtonReleaseLongPress"`
}

type Button struct {
	IconBase64         string `json:"IconBase64"`
	PositionX          int    `json:"Position_X"`
	PositionY          int    `json:"Position_Y"`
	LabelBase64        string `json:"LabelBase64"`
	BackgroundColorHex string `json:"BackgroundColorHex"`

	pid        int     `json:"-"`
	changed    bool    `json:"-"`
	iconPath   string  `json:"-"`
	iconText   string  `json:"-"`
	iconColor  string  `json:"-"`
	label      string  `json:"-"`
	labelSize  float64 `json:"-"`
	labelColor string  `json:"-"`
}

type Buttons struct {
	Method  string   `json:"Method"`
	Buttons []Button `json:"Buttons"`
}

func (b *Button) SetPid(pid int) {
	b.pid = pid
}

func (b *Button) GetPid() int {
	return b.pid
}

func (b *Buttons) AddButton(button Button) *Buttons {
	b.Buttons = append(b.Buttons, button)
	return b
}

func NewButton(row, col int) Button {
	return Button{
		PositionX:          col,
		PositionY:          row,
		IconBase64:         "",
		LabelBase64:        "",
		BackgroundColorHex: "#232323",

		iconPath:   "",
		iconText:   "",
		iconColor:  "#fff",
		label:      "",
		labelSize:  35,
		labelColor: "#fff",
		changed:    false,
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
		logger.Error("Error while reading icon file:", err, path)
		return
	}

	b.changed = true
	b.iconPath = path
	b.iconText = ""
	b.IconBase64 = base64.StdEncoding.EncodeToString(bytes)
}

func (b *Button) SetIconFromText(text, color string) {
	if b.iconText == text && b.iconColor == color {
		return
	}

	b.changed = true
	b.iconText = text
	b.iconColor = color
	b.iconPath = ""
	b.IconBase64 = label.GenerateIcon(text, color)
}

func (b *Button) GetIconText() (text string, color string) {
	return b.iconText, b.iconColor
}

func (b *Button) GetLabel() (text string, size float64, color string) {
	return b.label, b.labelSize, b.labelColor
}

func (b *Button) SetLabel(text string, size float64, color string) {
	if b.label == text && b.labelSize == size && b.labelColor == color {
		return
	}

	b.changed = true
	b.label = text
	b.labelSize = size
	b.labelColor = color
	b.LabelBase64 = label.GenerateLabel(text, size, color)
}

func (b *Button) UpdateFromAnyMap(m map[string]any) {
	if v, ok := m["color"].(string); ok {
		b.SetColor(v)
	}
	if v, ok := m["icon_path"].(string); ok {
		b.SetIconFromPath(v)
	}
	if m["icon_text"] != nil || m["icon_color"] != nil {
		text, color := b.GetIconText()
		if v, ok := m["icon_text"].(string); ok {
			text = v
		}
		if v, ok := m["icon_color"].(string); ok {
			color = v
		}
		b.SetIconFromText(text, color)
	}
	if m["label"] != nil || m["label_size"] != nil || m["label_color"] != nil {
		text, size, color := b.GetLabel()
		if v, ok := m["label"].(string); ok {
			text = v
		}
		if v, ok := m["label_size"].(float64); ok {
			size = v
		}
		if v, ok := m["label_color"].(string); ok {
			color = v
		}
		b.SetLabel(text, size, color)
	}
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
