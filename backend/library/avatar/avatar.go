//
// Date: 2018-03-21
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-12-29
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//
// Stolen From: https://github.com/ae0000/avatar
//

package avatar

import (
	"bufio"
	"flag"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

// init sets default environment variables for avatar functionality during tests
func init() {
	// Only set defaults during tests
	if flag.Lookup("test.v") != nil {
		// Set font path relative to this source file
		_, b, _, _ := runtime.Caller(0)
		basepath := filepath.Dir(b)
		fontPath := filepath.Join(basepath, "..", "..", "fonts")
		setDefaultIfEmpty("FONT_PATH", fontPath)
	}
}

// setDefaultIfEmpty sets an environment variable to a default value if it's not already set
func setDefaultIfEmpty(key, defaultValue string) {
	if os.Getenv(key) == "" {
		os.Setenv(key, defaultValue)
	}
}

const (
	defaultfontFace = "Roboto-Bold.ttf" //SourceSansVariable-Roman.ttf"
	fontSize        = 210.0
	imageWidth      = 600.0
	imageHeight     = 600.0
	dpi             = 72.0
	spacer          = 20
	textY           = 370
)

var fontFacePath = ""

// SetFontFacePath sets the font to do the business with
func SetFontFacePath(f string) {
	fontFacePath = f
}

// ToDisk saves the image to disk
func ToDisk(initials, path string) error {
	rgba, err := createAvatar(initials)
	if err != nil {
		return err
	}

	// Save image to disk
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	b := bufio.NewWriter(out)

	err = png.Encode(b, rgba)
	if err != nil {
		return err
	}

	err = b.Flush()
	if err != nil {
		return err
	}

	return nil
}

func cleanString(incoming string) string {
	incoming = strings.TrimSpace(incoming)

	// If its something like "firstname surname" get the initials out
	split := strings.Split(incoming, " ")
	if len(split) == 2 {
		incoming = split[0][0:1] + split[1][0:1]
	}

	// Max length of 2
	if len(incoming) > 2 {
		incoming = incoming[0:2]
	}

	// To upper and trimmed
	return strings.ToUpper(strings.TrimSpace(incoming))
}

func getFont(fontPath string) (*truetype.Font, error) {
	if fontPath == "" {
		fontPath = os.Getenv("FONT_PATH") + "/" + defaultfontFace
	}
	// Read the font data.
	fontBytes, err := ioutil.ReadFile(fontPath) //fmt.Sprintf("%s/%s", sourceDir, fontFaceName))
	if err != nil {
		return nil, err
	}

	return freetype.ParseFont(fontBytes)
}

func createAvatar(initials string) (*image.RGBA, error) {
	// Make a Regex to say we only want letters and numbers
	reg, err := regexp.Compile(`[^a-zA-Z0-9\s]+`)

	if err != nil {
		return nil, err
	}

	initials = strings.TrimSpace(reg.ReplaceAllString(initials, ""))

	// Make sure initials is not empty
	if len(initials) <= 1 {
		initials = "**"
	}

	// Make sure the string is OK
	text := cleanString(initials)

	// Load and get the font
	f, err := getFont(fontFacePath)
	if err != nil {
		return nil, err
	}

	// Setup the colors, text white, background based on first initial
	textColor := image.White
	background := defaultColor(text[0:1])
	rgba := image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))
	draw.Draw(rgba, rgba.Bounds(), &background, image.ZP, draw.Src)
	c := freetype.NewContext()
	c.SetDPI(dpi)
	c.SetFont(f)
	c.SetFontSize(fontSize)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(textColor)
	c.SetHinting(font.HintingFull)

	// We need to convert the font into a "font.Face" so we can read the glyph
	// info
	to := truetype.Options{}
	to.Size = fontSize
	face := truetype.NewFace(f, &to)

	// Calculate the widths and print to image
	xPoints := []int{0, 0}
	textWidths := []int{0, 0}

	// Get the widths of the text characters
	for i, char := range text {
		width, ok := face.GlyphAdvance(rune(char))
		if !ok {
			return nil, err
		}

		textWidths[i] = int(float64(width) / 64)
	}

	// TODO need some tests for this
	if len(textWidths) == 1 {
		textWidths[1] = 0
	}

	// Get the combined width of the characters
	combinedWidth := textWidths[0] + spacer + textWidths[1]

	// Draw first character
	xPoints[0] = int((imageWidth - combinedWidth) / 2)
	xPoints[1] = int(xPoints[0] + textWidths[0] + spacer)

	for i, char := range text {
		pt := freetype.Pt(xPoints[i], textY)
		_, err := c.DrawString(string(char), pt)
		if err != nil {
			return nil, err
		}
	}

	return rgba, nil
}
