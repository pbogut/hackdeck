package types

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
	buttons       map[ButtonPos]*Button
	buttonConfigs map[ButtonPos]*ButtonConfig
}

func (s *State) Init(rows, cols int) {
	buttonsCount := rows * cols
	s.buttons = make(map[ButtonPos]*Button, buttonsCount)
	s.buttonConfigs = make(map[ButtonPos]*ButtonConfig, buttonsCount)
}

func (s *State) AddButton(btn *Button, cfg *ButtonConfig) {
	s.buttons[ButtonPos{btn.PositionY, btn.PositionX}] = btn
	s.buttonConfigs[ButtonPos{btn.PositionY, btn.PositionX}] = cfg
}

func (s *State) GetButton(row, col int) *Button {
	return s.buttons[ButtonPos{row, col}]
}

func (s *State) GetButtonConfig(row, col int) *ButtonConfig {
	return s.buttonConfigs[ButtonPos{row, col}]
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
