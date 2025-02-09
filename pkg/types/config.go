package types

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/kirsle/configdir"
	"github.com/pbogut/hackdeck/pkg/logger"
)

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
	IconPath               string  `toml:"icon_path"`
	IconText               string  `toml:"icon_text"`
	IconColor              string  `toml:"icon_color"`
	ButtonPress            string  `toml:"button_press"`
	ButtonRelease          string  `toml:"button_release"`
	ButtonLongPress        string  `toml:"button_long_press"`
	ButtonLongPressRelease string  `toml:"button_long_press_release"`
	Interval               int     `toml:"interval"`
	Execute                string  `toml:"execute"`
	Label                  string  `toml:"label"`
	LabelSize              float64 `toml:"label_size"`
	LabelColor             string  `toml:"label_color"`
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

	configFile := configdir.LocalConfig("hackdeck", "hackdeck.toml")
	logger.Debug("Looking for config file:", configFile)
	_, err := os.Stat(configFile)
	if err != nil {
		logger.Debug("Config not found:", configFile)
		configFile = "hackdeck.toml"
		logger.Debug("Looking for config file:", configFile)
	}

	_, err = os.Stat(configFile)
	if err != nil {
		logger.Debug("Config not found:", configFile)
		return config
	}

	logger.Debug("Config file found:", configFile)

	if _, err := toml.DecodeFile(configFile, &config); err != nil {
		logger.Error("Error while decoding config file:", err)
		return config
	}

	return config
}
