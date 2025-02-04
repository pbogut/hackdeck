package types

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

var configFile = "hackdeck.toml"

type Config struct {
	Rows                          int
	Columns                       int
	ButtonSpacing                 int  `toml:"button_spacing"`
	ButtonRadius                  int  `toml:"button_radius"`
	ButtonBackground              bool `toml:"button_background"`
	Brightness                    float32
	AutoConnect                   bool     `toml:"auto_connect"`
	SupportButtonReleaseLongPress bool     `toml:"support_button_release_long_press"`
	ShellCommand                  string   `toml:"shell_command"`
	ShellArguments                []string `toml:"shell_arguments"`

	Buttons []ButtonConfig
}

type ButtonConfig struct {
	Row                    int
	Column                 int
	Color                  string
	Icon                   string
	ButtonPress            string `toml:"button_press"`
	ButtonRelease          string `toml:"button_release"`
	ButtonLongPress        string `toml:"button_long_press"`
	ButtonLongPressRelease string `toml:"button_long_press_release"`
}

func ReadConfig() Config {
	// default config
	config := Config{
		Rows:                          3,
		Columns:                       5,
		ButtonSpacing:                 10,
		ButtonRadius:                  40,
		ButtonBackground:              true,
		Brightness:                    0.3,
		AutoConnect:                   false,
		SupportButtonReleaseLongPress: true,

		ShellCommand:   "bash",
		ShellArguments: []string{"-c"},

		Buttons: []ButtonConfig{},
	}

	_, err := os.Stat(configFile)
	if err != nil {
		fmt.Println("Config file is missing: ", configFile)
		return config
	}

	if _, err := toml.DecodeFile(configFile, &config); err != nil {
		fmt.Println("Error while decoding config file:", err)
		return config
	}

	return config
}
