package label

import (
	"bytes"
	_ "embed"
	"encoding/base64"
	"errors"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"strings"

	"github.com/golang/freetype/truetype"
	"github.com/pbogut/hackdeck/pkg/logger"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

//go:embed JetBrainsMonoNerdFontMono-Regular.ttf
var fontBytes []byte

func GenerateIcon(icon string, hexColor string) string {
	img := image.White
	color, err := parseHexColor(hexColor)
	if err == nil {
		img = &image.Uniform{C: color}
	}
	return generateImage(icon, 250.0, 30, img)
}

func GenerateLabel(text string, size float64, hexColor string) string {
	img := image.White
	color, err := parseHexColor(hexColor)
	if err == nil {
		img = &image.Uniform{C: color}
	}
	return generateImage(text, size, 10, img)
}

func generateImage(text string, size float64, margin int, fg *image.Uniform) string {
	width := 250
	height := 250
	dpi := 72.0
	lineheight := 1.2
	hinting := "none"

	f, err := truetype.Parse(fontBytes)
	if err != nil {
		logger.Errorf("Fant parsing error: %s", err)
		return ""
	}
	// Draw the background
	rgba := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(rgba, rgba.Bounds(), image.Transparent, image.Point{0, 0}, draw.Src)

	h := font.HintingFull
	switch hinting {
	case "full":
		h = font.HintingFull
	}
	drawer := &font.Drawer{
		Dst: rgba,
		Src: fg,
		Face: truetype.NewFace(f, &truetype.Options{
			Size:    size,
			DPI:     dpi,
			Hinting: h,
		}),
	}
	outline := &font.Drawer{
		Dst: rgba,
		Src: image.Black,
		Face: truetype.NewFace(f, &truetype.Options{
			Size:    size,
			DPI:     dpi,
			Hinting: h,
		}),
	}
	lines := strings.Split(strings.Trim(text, "\n"), "\n")
	dy := int(math.Ceil(size * lineheight))
	x := margin
	y := height - margin - ((len(lines) - 1) * dy)

	for _, line := range lines {
		// @todo figure better way for outline
		textWidth := drawer.MeasureString(line).Round()
		centerX := ((width - textWidth) / 2) - margin

		for _, shiftX := range []int{-3, 0, 3} {
			for _, shiftY := range []int{-3, 0, 3} {
				outline.Dot = fixed.P(x+shiftX+centerX, y+shiftY)
				outline.DrawString(line)
			}
		}
		drawer.Dot = fixed.P(x+centerX, y)
		drawer.DrawString(line)
		y += dy
	}

	buffer := new(bytes.Buffer)
	if err := png.Encode(buffer, rgba); err != nil {
		logger.Error("Unable to encode image.")
	}

	result := base64.StdEncoding.EncodeToString(buffer.Bytes())
	return result
}

var errInvalidFormat = errors.New("invalid format")

func parseHexColor(s string) (c color.RGBA, err error) {
	c.A = 0xff

	if s == "" {
		return c, errInvalidFormat
	}

	if s[0] != '#' {
		logger.Errorf("Color must start with #: %s", s)
		return c, errInvalidFormat
	}

	hexToByte := func(b byte) byte {
		switch {
		case b >= '0' && b <= '9':
			return b - '0'
		case b >= 'a' && b <= 'f':
			return b - 'a' + 10
		case b >= 'A' && b <= 'F':
			return b - 'A' + 10
		}
		err = errInvalidFormat
		return 0
	}

	switch len(s) {
	case 7:
		c.R = hexToByte(s[1])<<4 + hexToByte(s[2])
		c.G = hexToByte(s[3])<<4 + hexToByte(s[4])
		c.B = hexToByte(s[5])<<4 + hexToByte(s[6])
	case 4:
		c.R = hexToByte(s[1]) * 17
		c.G = hexToByte(s[2]) * 17
		c.B = hexToByte(s[3]) * 17
	default:
		err = errInvalidFormat
		logger.Errorf("Color format is invalid: %s", s)
		return c, err
	}
	return c, nil
}
