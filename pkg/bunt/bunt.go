// Copyright © 2018 Matthias Diester
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

package bunt

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	isatty "github.com/mattn/go-isatty"
)

// The named colors are based upon https://en.wikipedia.org/wiki/Web_colors

// NoColor specifies whether color sequences are used or not. It is automatically initialized based on the detected terminal is use
var NoColor = os.Getenv("TERM") == "dumb" ||
	!isatty.IsTerminal(os.Stdout.Fd()) &&
		!isatty.IsCygwinTerminal(os.Stdout.Fd())

// seq is the ANSI escape sequence used for the coloring
const seq = "\x1b"

// Pink colors
const (
	Pink            = 0x00FFC0CB
	LightPink       = 0x00FFB6C1
	HotPink         = 0x00FF69B4
	DeepPink        = 0x00FF1493
	PaleVioletRed   = 0x00DB7093
	MediumVioletRed = 0x00C71585
)

// Red colors
const (
	LightSalmon = 0x00FFA07A
	Salmon      = 0x00FA8072
	DarkSalmon  = 0x00E9967A
	LightCoral  = 0x00F08080
	IndianRed   = 0x00CD5C5C
	Crimson     = 0x00DC143C
	FireBrick   = 0x00B22222
	DarkRed     = 0x008B0000
	Red         = 0x00FF0000
)

// Orange colors
const (
	OrangeRed  = 0x00FF4500
	Tomato     = 0x00FF6347
	Coral      = 0x00FF7F50
	DarkOrange = 0x00FF8C00
	Orange     = 0x00FFA500
)

// Yellow colors
const (
	Yellow               = 0x00FFFF00
	LightYellow          = 0x00FFFFE0
	LemonChiffon         = 0x00FFFACD
	LightGoldenrodYellow = 0x00FAFAD2
	PapayaWhip           = 0x00FFEFD5
	Moccasin             = 0x00FFE4B5
	PeachPuff            = 0x00FFDAB9
	PaleGoldenrod        = 0x00EEE8AA
	Khaki                = 0x00F0E68C
	DarkKhaki            = 0x00BDB76B
	Gold                 = 0x00FFD700
)

// Brown colors
const (
	Cornsilk       = 0x00FFF8DC
	BlanchedAlmond = 0x00FFEBCD
	Bisque         = 0x00FFE4C4
	NavajoWhite    = 0x00FFDEAD
	Wheat          = 0x00F5DEB3
	BurlyWood      = 0x00DEB887
	Tan            = 0x00D2B48C
	RosyBrown      = 0x00BC8F8F
	SandyBrown     = 0x00F4A460
	Goldenrod      = 0x00DAA520
	DarkGoldenrod  = 0x00B8860B
	Peru           = 0x00CD853F
	Chocolate      = 0x00D2691E
	SaddleBrown    = 0x008B4513
	Sienna         = 0x00A0522D
	Brown          = 0x00A52A2A
	Maroon         = 0x00800000
)

// Green colors
const (
	DarkOliveGreen    = 0x00556B2F
	Olive             = 0x00808000
	OliveDrab         = 0x006B8E23
	YellowGreen       = 0x009ACD32
	LimeGreen         = 0x0032CD32
	Lime              = 0x0000FF00
	LawnGreen         = 0x007CFC00
	Chartreuse        = 0x007FFF00
	GreenYellow       = 0x00ADFF2F
	SpringGreen       = 0x0000FF7F
	MediumSpringGreen = 0x0000FA9A
	LightGreen        = 0x0090EE90
	PaleGreen         = 0x0098FB98
	DarkSeaGreen      = 0x008FBC8F
	MediumAquamarine  = 0x0066CDAA
	MediumSeaGreen    = 0x003CB371
	SeaGreen          = 0x002E8B57
	ForestGreen       = 0x00228B22
	Green             = 0x00008000
	DarkGreen         = 0x00006400
)

// Cyan colors
const (
	Aqua            = 0x0000FFFF
	Cyan            = 0x0000FFFF
	LightCyan       = 0x00E0FFFF
	PaleTurquoise   = 0x00AFEEEE
	Aquamarine      = 0x007FFFD4
	Turquoise       = 0x0040E0D0
	MediumTurquoise = 0x0048D1CC
	DarkTurquoise   = 0x0000CED1
	LightSeaGreen   = 0x0020B2AA
	CadetBlue       = 0x005F9EA0
	DarkCyan        = 0x00008B8B
	Teal            = 0x00008080
)

// Blue colors
const (
	LightSteelBlue = 0x00B0C4DE
	PowderBlue     = 0x00B0E0E6
	LightBlue      = 0x00ADD8E6
	SkyBlue        = 0x0087CEEB
	LightSkyBlue   = 0x0087CEFA
	DeepSkyBlue    = 0x0000BFFF
	DodgerBlue     = 0x001E90FF
	CornflowerBlue = 0x006495ED
	SteelBlue      = 0x004682B4
	RoyalBlue      = 0x004169E1
	Blue           = 0x000000FF
	MediumBlue     = 0x000000CD
	DarkBlue       = 0x0000008B
	Navy           = 0x00000080
	MidnightBlue   = 0x00191970
)

// Purple, violet, and magenta colors
const (
	Lavender        = 0x00E6E6FA
	Thistle         = 0x00D8BFD8
	Plum            = 0x00DDA0DD
	Violet          = 0x00EE82EE
	Orchid          = 0x00DA70D6
	Fuchsia         = 0x00FF00FF
	Magenta         = 0x00FF00FF
	MediumOrchid    = 0x00BA55D3
	MediumPurple    = 0x009370DB
	BlueViolet      = 0x008A2BE2
	DarkViolet      = 0x009400D3
	DarkOrchid      = 0x009932CC
	DarkMagenta     = 0x008B008B
	Purple          = 0x00800080
	Indigo          = 0x004B0082
	DarkSlateBlue   = 0x00483D8B
	SlateBlue       = 0x006A5ACD
	MediumSlateBlue = 0x007B68EE
)

// White colors
const (
	White         = 0x00FFFFFF
	Snow          = 0x00FFFAFA
	Honeydew      = 0x00F0FFF0
	MintCream     = 0x00F5FFFA
	Azure         = 0x00F0FFFF
	AliceBlue     = 0x00F0F8FF
	GhostWhite    = 0x00F8F8FF
	WhiteSmoke    = 0x00F5F5F5
	Seashell      = 0x00FFF5EE
	Beige         = 0x00F5F5DC
	OldLace       = 0x00FDF5E6
	FloralWhite   = 0x00FFFAF0
	Ivory         = 0x00FFFFF0
	AntiqueWhite  = 0x00FAEBD7
	Linen         = 0x00FAF0E6
	LavenderBlush = 0x00FFF0F5
	MistyRose     = 0x00FFE4E1
)

// Gray and black colors
const (
	Gainsboro      = 0x00DCDCDC
	LightGray      = 0x00D3D3D3
	Silver         = 0x00C0C0C0
	DarkGray       = 0x00A9A9A9
	Gray           = 0x00808080
	DimGray        = 0x00696969
	LightSlateGray = 0x00778899
	SlateGray      = 0x00708090
	DarkSlateGray  = 0x002F4F4F
	Black          = 0x00000000
)

// Special additional colors
const (
	AdditionGreen      = 0x0058BF38
	RemovalRed         = 0x00B9311B
	ModificationYellow = 0x00C7C43F
)

// Modifiers
const (
	Bold      = 1
	Italic    = 3
	Underline = 4
)

// String defines a string that consists of differently colored text segments
type String []Segment

// Segment defines a piece of text and its designed style attributes (which can be empty)
type Segment struct {
	Attributes []uint32
	Data       string
}

// https://regex101.com/segmentRegexp/ulipXZ/3
var segmentRegexp = regexp.MustCompile(fmt.Sprintf(`(?m)(.*?)((%s\[(\d+(;\d+)*)m)(.+?)(%s\[0m))`, seq, seq))

// Colorize applies an ANSI truecolor sequence for the provided color to the given text.
func Colorize(text string, color uint32, modifiers ...uint32) string {
	modifiers = keepAttributes(modifiers, []uint32{Bold, Italic, Underline})
	sort.Slice(modifiers, func(i, j int) bool {
		return modifiers[i] < modifiers[j]
	})

	r, g, b := BreakUpColorIntoChannels(color)
	colorCoding := []uint32{38, 2, uint32(r), uint32(g), uint32(b)}

	return wrapTextInSeq(text, append(modifiers, colorCoding...)...)
}

// Style applies only text modifications like Bold, Italic, or Underline to the text
func Style(text string, modifiers ...uint32) string {
	return wrapTextInSeq(text, keepAttributes(modifiers, []uint32{Bold, Italic, Underline})...)
}

// BoldText is a convenience function to make the string bold
func BoldText(text string) string {
	return Style(text, Bold)
}

// RemoveAllEscapeSequences return the input string with all escape sequences removed
func RemoveAllEscapeSequences(input string) string {
	regexp := regexp.MustCompile(seq + `\[\d+(;\d+)*m`)

	for loc := regexp.FindStringIndex(input); loc != nil; loc = regexp.FindStringIndex(input) {
		start := loc[0]
		end := loc[1]
		input = strings.Replace(input, input[start:end], "", -1)
	}

	return input
}

// BreakUpColorIntoChannels takes the provided color and breaks into the its red, green, and blue color channels
func BreakUpColorIntoChannels(color uint32) (uint8, uint8, uint8) {
	b := uint8(color & 0x000000FF)
	color >>= 8
	g := uint8(color & 0x000000FF)
	color >>= 8
	r := uint8(color & 0x000000FF)

	return r, g, b
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
		attributes := make([]uint32, 0, len(splitted))
		for _, attribute := range splitted {
			var value int
			var err error
			if value, err = strconv.Atoi(attribute); err != nil {
				return nil, fmt.Errorf("Failed to split input string in its colored segments: %s", err.Error())
			}

			attributes = append(attributes, uint32(value))
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

// func colorEachLine(color *color.Color, text string) string {
// 	var buf bytes.Buffer
//
// 	splitted := strings.Split(text, "\n")
// 	length := len(splitted)
// 	for idx, line := range splitted {
// 		buf.WriteString(color.Sprint(line))
//
// 		if idx < length-1 {
// 			buf.WriteString("\n")
// 		}
// 	}
//
// 	return buf.String()
// }

func wrapTextInSeq(text string, attributes ...uint32) string {
	if NoColor {
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

func wrapLineInSeq(line string, attributes ...uint32) string {
	cstring, err := BreakUpStringIntoColorSegments(line)
	if err != nil { // TODO Debug, or trace the error
		return line
	}

	var buf bytes.Buffer
	for _, segment := range cstring {
		segmentAttributes := keepAttributes(segment.Attributes, []uint32{Bold, Italic, Underline})

		buf.WriteString(fmt.Sprintf("%s[%sm", seq, attributesAsList(append(segmentAttributes, attributes...))))
		buf.WriteString(segment.Data)
		buf.WriteString(fmt.Sprintf("%s[%dm", seq, 0))
	}

	return buf.String()
}

func attributesAsList(attributes []uint32) string {
	fields := make([]string, len(attributes))
	for i, attribute := range attributes {
		fields[i] = strconv.FormatUint(uint64(attribute), 10)
	}

	return strings.Join(fields, ";")
}

func keepAttributes(attributes []uint32, keepers []uint32) []uint32 {
	result := make([]uint32, 0)
	for _, toBeKept := range keepers {
		if i := getIndexOf(attributes, toBeKept); i >= 0 {
			result = append(result, toBeKept)
		}
	}

	return result
}

func getIndexOf(list []uint32, entry uint32) int {
	for i, element := range list {
		if element == entry {
			return i
		}
	}

	return -1
}
