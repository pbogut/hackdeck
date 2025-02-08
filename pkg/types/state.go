package types

import (
	"io"
	"sync"
)

const (
	BUTTON_PRESS = iota
	BUTTON_RELEASE
	BUTTON_LONG_PRESS
	BUTTON_LONG_PRESS_RELEASE
)

type ButtonPos struct {
	row int
	col int
}

type State struct {
	mutex         sync.Mutex
	buttons       map[ButtonPos]*Button
	buttonConfigs map[ButtonPos]*ButtonConfig
	pipes         map[ButtonPos]*io.WriteCloser
}

func (s *State) AddPipe(row, col int, pipe *io.WriteCloser) {
	s.mutex.Lock()
	s.pipes[ButtonPos{row, col}] = pipe
	s.mutex.Unlock()
}

func (s *State) GetPipe(row, col int) *io.WriteCloser {
	return s.pipes[ButtonPos{row, col}]
}

func (s *State) Init(config Config) {
	buttonsCount := config.Rows * config.Columns
	s.buttons = make(map[ButtonPos]*Button, buttonsCount)
	s.buttonConfigs = make(map[ButtonPos]*ButtonConfig, buttonsCount)
	s.pipes = make(map[ButtonPos]*io.WriteCloser, buttonsCount)

	for _, btnCfg := range config.Buttons {
		button := NewButton(btnCfg.Row, btnCfg.Column)
		button.SetColor(btnCfg.Color)
		button.SetIconFromPath(btnCfg.IconPath)
		button.SetIconFromText(btnCfg.IconText)
		button.SetIconColor(btnCfg.IconColor)
		button.SetLabel(btnCfg.Label)
		s.AddButton(&button, &btnCfg)
	}
}

func (s *State) AddButton(btn *Button, cfg *ButtonConfig) {
	s.mutex.Lock()
	s.buttons[ButtonPos{btn.PositionY, btn.PositionX}] = btn
	s.buttonConfigs[ButtonPos{btn.PositionY, btn.PositionX}] = cfg
	s.mutex.Unlock()
}

func (s *State) GetButton(row, col int) *Button {
	return s.buttons[ButtonPos{row, col}]
}

func (s *State) GetButtons() map[ButtonPos]*Button {
	return s.buttons
}

func (s *State) GetButtonConfig(row, col int) *ButtonConfig {
	return s.buttonConfigs[ButtonPos{row, col}]
}

func (s *State) GetButtonConfigs() map[ButtonPos]*ButtonConfig {
	return s.buttonConfigs
}

func (s *State) GetCmd(row, col, state int) string {
	switch state {
	case BUTTON_PRESS:
		return s.GetButtonConfig(row, col).ButtonPress
	case BUTTON_RELEASE:
		return s.GetButtonConfig(row, col).ButtonRelease
	case BUTTON_LONG_PRESS:
		return s.GetButtonConfig(row, col).ButtonLongPress
	case BUTTON_LONG_PRESS_RELEASE:
		return s.GetButtonConfig(row, col).ButtonLongPressRelease
	}
	return ""
}
