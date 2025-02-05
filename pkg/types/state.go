package types

const (
	BUTTON_PRESS = iota
	BUTTON_RELEASE
	BUTTON_LONG_PRESS
	BUTTON_LONG_PRESS_RELEASE
)

type CmdPos struct {
	row   int
	col   int
	state int
}
type ButtonPos struct {
	row int
	col int
}

type State struct {
	cmds          map[CmdPos]string
	buttons       map[ButtonPos]*Button
	buttonConfigs map[ButtonPos]*ButtonConfig
}

func (s *State) Init(rows, cols int) {
	actionCount := rows * cols * 4
	buttonsCount := rows * cols
	s.cmds = make(map[CmdPos]string, actionCount)
	s.buttons = make(map[ButtonPos]*Button, buttonsCount)
	s.buttonConfigs = make(map[ButtonPos]*ButtonConfig, buttonsCount)
}

func (s *State) AddButton(btn *Button, cfg *ButtonConfig) {
	s.buttons[ButtonPos{btn.PositionY, btn.PositionX}] = btn
	s.buttonConfigs[ButtonPos{btn.PositionY, btn.PositionX}] = cfg
	s.AddCmd(btn.PositionY, btn.PositionX, BUTTON_PRESS, cfg.ButtonPress)
	s.AddCmd(btn.PositionY, btn.PositionX, BUTTON_RELEASE, cfg.ButtonRelease)
	s.AddCmd(btn.PositionY, btn.PositionX, BUTTON_LONG_PRESS, cfg.ButtonLongPress)
	s.AddCmd(btn.PositionY, btn.PositionX, BUTTON_LONG_PRESS_RELEASE, cfg.ButtonLongPressRelease)
}

func (s *State) GetButton(row, col int) *Button {
	return s.buttons[ButtonPos{row, col}]
}

func (s *State) GetButtonConfig(row, col int) *ButtonConfig {
	return s.buttonConfigs[ButtonPos{row, col}]
}

func (s *State) AddCmd(row, col, state int, cmd string) {
	s.cmds[CmdPos{row, col, state}] = cmd
}

func (s *State) GetCmd(row, col, state int) string {
	return s.cmds[CmdPos{row, col, state}]
}
