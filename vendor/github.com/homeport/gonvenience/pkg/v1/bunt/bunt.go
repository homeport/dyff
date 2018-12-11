// Copyright © 2018 The Homeport Team
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

/*
Package bunt provides terminal true color support.

Terminals that support 24 bit (true color) colors with ANSI sequences, the
`bunt` package provides functions to wrap text into the respective sequence.

In case the terminal does not support true color, a fallback option is used to
calcucate amd pick the best matching color out of the 16 color basic color
schema instead.

If the input string already contains ANSI sequences, the wrapping function will
break up the string first and apply the target color to each respective segment
of the input string to preserve any existing text emphasis like `bold`, or
`italic` while additionally applying the desired color.
*/
package bunt

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/homeport/gonvenience/pkg/v1/term"
	colorful "github.com/lucasb-eyer/go-colorful"
	ciede2000 "github.com/mattn/go-ciede2000"
)

// The named colors are based upon https://en.wikipedia.org/wiki/Web_colors
// Nice page for color code conversions: https://convertingcolors.com/

// ColorSetting defines the coloring setting to be used
var ColorSetting = AUTO

// TrueColorSetting defines the true color usage setting to be used
var TrueColorSetting = AUTO

// seq is the ANSI escape sequence used for the coloring
const seq = "\x1b"

// Attribute is used to specify the attributes (codes) used in ANSI sequences
type Attribute uint8

// SwitchState is the type to cover different preferences/settings like:
// on, off, or auto
type SwitchState int

// Supported setting states
const (
	OFF  = SwitchState(-1)
	AUTO = SwitchState(0)
	ON   = SwitchState(+1)
)

// Text modifiers (emphasis)
const (
	Bold      = 1
	Italic    = 3
	Underline = 4
)

// String defines a string that consists of differently colored text segments
type String []Segment

// Segment defines a piece of text and its designed style attributes (which can be empty)
type Segment struct {
	Attributes []Attribute
	Data       string
}

// https://regex101.com/segmentRegexp/ulipXZ/3
var segmentRegexp = regexp.MustCompile(fmt.Sprintf(`(?m)(.*?)((%s\[(\d+(;\d+)*)m)(.+?)(%s\[0m))`, seq, seq))

// UseColors return whether colors are used or not based on the configured color setting
func UseColors() bool {
	return (ColorSetting == ON) ||
		(ColorSetting == AUTO && term.IsTerminal() && !term.IsDumbTerminal())
}

// UseTrueColor returns whether true color colors should be used or not base on
// the configured true color usage setting
func UseTrueColor() bool {
	return (TrueColorSetting == ON) ||
		(TrueColorSetting == AUTO && term.IsTrueColor())
}

// Colorize applies an ANSI truecolor sequence for the provided color to the given text.
func Colorize(text string, color colorful.Color, modifiers ...Attribute) string {
	return applyColorAndAttributes(text, false, color, modifiers...)
}

// Colorize4bit applies the best suitable 4-bit color attribute to the given text.
func Colorize4bit(text string, color colorful.Color, modifiers ...Attribute) string {
	return wrapTextInSeq(text, append(modifiers, Get4bitEquivalentColorAttribute(color))...)
}

// ColorizeFgBg applies an ANSI truecolor sequence for the provided foreground and background colors to the given text.
func ColorizeFgBg(text string, foreground colorful.Color, background colorful.Color, modifiers ...Attribute) string {
	modifiers = keepAttributes(modifiers, []Attribute{Bold, Italic, Underline})
	sort.Slice(modifiers, func(i, j int) bool {
		return modifiers[i] < modifiers[j]
	})

	fgr, fgg, fgb := foreground.RGB255()
	fgColorCoding := []Attribute{38, 2, Attribute(fgr), Attribute(fgg), Attribute(fgb)}

	bgr, bgg, bgb := background.RGB255()
	bgColorCoding := []Attribute{48, 2, Attribute(bgr), Attribute(bgg), Attribute(bgb)}

	return wrapTextInSeq(text, append(append(modifiers, bgColorCoding...), fgColorCoding...)...)
}

// Style applies only text modifications like Bold, Italic, or Underline to the text
func Style(text string, modifiers ...Attribute) string {
	return wrapTextInSeq(text, keepAttributes(modifiers, []Attribute{Bold, Italic, Underline})...)
}

// BoldText is a convenience function to make the string bold
func BoldText(text string) string {
	return Style(text, Bold)
}

// RemoveAllEscapeSequences return the input string with all escape sequences removed
func RemoveAllEscapeSequences(input string) string {
	escapeSeqFinderRegExp := regexp.MustCompile(seq + `\[\d+(;\d+)*m`)

	for loc := escapeSeqFinderRegExp.FindStringIndex(input); loc != nil; loc = escapeSeqFinderRegExp.FindStringIndex(input) {
		start := loc[0]
		end := loc[1]
		input = strings.Replace(input, input[start:end], "", -1)
	}

	return input
}

// BreakUpStringIntoColorSegments takes the input string and breaks it up into individual substrings that each have their own coloring already applied
func BreakUpStringIntoColorSegments(input string) (String, error) {
	var result String

	lastPos := 0
	for _, submatch := range segmentRegexp.FindAllStringSubmatchIndex(input, -1) {
		// 0:1 full capture group
		// 2:3 text before first espace sequence
		// 4:5 text in espace sequence including the start and end escape sequence
		// 6:7 start escape sequence
		// 8:9 attributes of start escape sequence
		// 10:11 optinal attributes of start escape sequence
		// 12:13 the actual text
		// 14:15 end escape sequence
		headTextSegment := input[submatch[2]:submatch[3]]
		attributesPart := input[submatch[8]:submatch[9]]
		textSegment := input[submatch[12]:submatch[13]]

		// Break up the style attributes to convert it from string into proper int attributes
		splitted := strings.Split(attributesPart, ";")
		attributes := make([]Attribute, 0, len(splitted))
		for _, attribute := range splitted {
			var value int
			var err error
			if value, err = strconv.Atoi(attribute); err != nil {
				return nil, fmt.Errorf("Failed to split input string in its colored segments: %s", err.Error())
			}

			attributes = append(attributes, Attribute(value))
		}

		// Append any piece of text that is between two identified segments (mostely line breaks and stuff)
		if lastPos != submatch[0] {
			result = append(result, Segment{Data: input[lastPos:submatch[0]]})
		}

		// Append text segment that was identified as being before the styled segment
		result = append(result, Segment{Data: headTextSegment})

		// Append the segment including its identified style attributes
		result = append(result, Segment{Data: textSegment, Attributes: attributes})

		lastPos = submatch[1]
	}

	// Append any trailing text that is still left
	if lastPos < len(input) {
		result = append(result, Segment{Data: input[lastPos:]})
	}

	return result, nil
}

// Get4bitEquivalentColorAttribute returns the color attribute which matches the
// best with the provided color.
func Get4bitEquivalentColorAttribute(targetColor colorful.Color) Attribute {
	red, green, blue := targetColor.RGB255()
	target := &color.RGBA{red, green, blue, 0xFF}

	min := math.MaxFloat64
	result := Attribute(0)

	// Calculate the distance between the target color and the available 4-bit
	// colors using the `deltaE` algorithm to find the best match.
	for candidate, attribute := range translateMap4bitColor {
		r, g, b := candidate.RGB255()
		test := &color.RGBA{r, g, b, 0xFF}

		if distance := ciede2000.Diff(target, test); distance < min {
			min = distance
			result = attribute
		}
	}

	return result
}

// ParseSetting converts a string with a setting into its respective state
func ParseSetting(setting string) (SwitchState, error) {
	switch strings.ToLower(setting) {
	case "auto":
		return AUTO, nil

	case "off", "no", "false":
		return OFF, nil

	case "on", "yes", "true":
		return ON, nil

	default:
		return OFF, fmt.Errorf("invalid state '%s' used, supported modes are: auto, on, or off", setting)
	}
}

func applyColorAndAttributes(text string, blendColors bool, color colorful.Color, modifiers ...Attribute) string {
	modifiers = keepAttributes(modifiers, []Attribute{Bold, Italic, Underline})
	sort.Slice(modifiers, func(i, j int) bool {
		return modifiers[i] < modifiers[j]
	})

	if UseTrueColor() {
		r, g, b := color.RGB255()
		colorCoding := []Attribute{38, 2, Attribute(r), Attribute(g), Attribute(b)}
		return wrapTextInSeq(text, append(modifiers, colorCoding...)...)
	}

	colorAttribute := Get4bitEquivalentColorAttribute(color)
	return wrapTextInSeq(text, append(modifiers, colorAttribute)...)
}

func wrapTextInSeq(text string, attributes ...Attribute) string {
	if !UseColors() {
		return text
	}

	var buf bytes.Buffer

	splitted := strings.Split(text, "\n")
	length := len(splitted)

	for idx, line := range splitted {
		buf.WriteString(wrapLineInSeq(line, attributes...))

		if idx < length-1 {
			buf.WriteString("\n")
		}
	}

	return buf.String()
}

func wrapLineInSeq(line string, attributes ...Attribute) string {
	cstring, err := BreakUpStringIntoColorSegments(line)
	if err != nil {
		log.New(os.Stderr, "Error: ", log.LstdFlags).Printf("Text does have incorrect escape sequence: %s", err.Error())
		return line
	}

	var buf bytes.Buffer
	for _, segment := range cstring {
		segmentAttributes := keepAttributes(segment.Attributes, []Attribute{Bold, Italic, Underline})

		buf.WriteString(fmt.Sprintf("%s[%sm", seq, attributesAsList(append(segmentAttributes, attributes...))))
		buf.WriteString(segment.Data)
		buf.WriteString(fmt.Sprintf("%s[%dm", seq, 0))
	}

	return buf.String()
}

func attributesAsList(attributes []Attribute) string {
	fields := make([]string, len(attributes))
	for i, attribute := range attributes {
		fields[i] = strconv.FormatUint(uint64(attribute), 10)
	}

	return strings.Join(fields, ";")
}

func keepAttributes(attributes []Attribute, keepers []Attribute) []Attribute {
	result := make([]Attribute, 0)
	for _, toBeKept := range keepers {
		if i := getIndexOf(attributes, toBeKept); i >= 0 {
			result = append(result, toBeKept)
		}
	}

	return result
}

func getIndexOf(list []Attribute, entry Attribute) int {
	for i, element := range list {
		if element == entry {
			return i
		}
	}

	return -1
}
