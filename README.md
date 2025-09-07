# HackDeck

This is alternative server implementation for Macro Deck client apps.
Tested with Web and Android clients, written in Go.

## Demo

https://github.com/user-attachments/assets/049a1e5a-1dc0-47d4-a483-c8d17e7412d0

## Installation

```bash
go install github.com/pbogut/hackdeck@latest
```

## Configuration

Configuration file is located at:

Linux: `$HOME/.config/hackdeck/hackdeck.toml`
Windows: `C:\Users\<username>\AppData\Roaming\hackdeck\hackdeck.toml`)
macOS: `$HOME/Library/Application Support/hackdeck/hackdeck.toml`

### Properties

- `column` - number of columns in the grid (default 5)
- `rows` - number of rows in the grid (default 3)
- `port` - port to listen on (default 8191)
- `button_spacing` - space between buttons (default 10)
- `button_radius` - radius of buttons (default 40)
- `button_background` - (default true)
- `brightness` - (default 0.3)
- `support_button_release_long_press` - support for long pressing buttons (default true)
- `shell_command` - shell command to run commands when button is pressed (default bash)
- `shell_arguments` - arguments to pass to shell command (default -c)
- `[[buttons]]` - list of button configurations

### Button configuration

- `row` - row number of the button
- `column` - column number of the button
- `color` - hex color of the button (for example #000000)
- `icon_path` - path to icon file (for example /home/user/icon.png)
- `icon_text` - create icon from text, supports nerdfont glyphs
- `icon_color` - hex color of the generated icon (default #FFFFFF)
- `button_press` - command to run when button is pressed
- `button_release` - command to run when button is released
- `button_long_press` - command to run when button is long press
- `button_long_press_release` - command to run when button is released after long press
- `execute` - execute command when button is created, output can update icon parameters
- `interval` - how often to execute command (if not set, command is assumed to run in "daemon" mode)
- `label` - text to display on the button
- `label_size` - size of the label text (default 35)
- `label_color` - hex color of the label text (default #FFFFFF)

## License

MIT License;
The software is provided "as is", without warranty of any kind.
