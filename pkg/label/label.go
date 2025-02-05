package label

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"math"
	"os"
	"strings"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

func GenerateLabel(text string) string {
	margin := 20
	width := 250
	height := 250
	size := 50.0
	dpi := 72.0
	lineheight := 1.2
	// fg := image.Black
	fg := image.White

	fontfile := "./luxisr.ttf"
	hinting := "none"

	// Read the font data
	fontBytes, err := os.ReadFile(fontfile)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	f, err := truetype.Parse(fontBytes)
	if err != nil {
		fmt.Println(err)
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
	d := &font.Drawer{
		Dst: rgba,
		Src: fg,
		Face: truetype.NewFace(f, &truetype.Options{
			Size:    size,
			DPI:     dpi,
			Hinting: h,
		}),
	}
	x := margin
	y := margin
	dy := int(math.Ceil(size * lineheight))
	y += dy
	for _, line := range strings.Split(text, "\n") {
		d.Dot = fixed.P(x, y)
		d.DrawString(line)
		y += dy
	}

	buffer := new(bytes.Buffer)
	if err := png.Encode(buffer, rgba); err != nil {
		fmt.Println("unable to encode image.")
	}

	result := base64.StdEncoding.EncodeToString(buffer.Bytes())
	return result
}
