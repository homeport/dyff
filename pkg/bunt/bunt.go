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
	"image/color"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	ciede2000 "github.com/mattn/go-ciede2000"
	isatty "github.com/mattn/go-isatty"
)

// The named colors are based upon https://en.wikipedia.org/wiki/Web_colors
// Nice page for color code conversions: https://convertingcolors.com/

// isDumbTerminal is initialised with true if the current terminal has limited support for escape sequences, or false otherwise
var isDumbTerminal = os.Getenv("TERM") == "dumb"

// isTerminal is initialised with true if the current STDOUT stream writes to a terminal (not a redirect), or false otherwise
var isTerminal = isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())

// isTrueColor is initialised with true if the current terminal reports to support 24-bit colors, or false otherwise
var isTrueColor = func() bool {
	switch os.Getenv("COLORTERM") {
	case "truecolor", "24bit":
		return true
	}

	return false
}()

// ColorSetting defines the coloring setting to be used
var ColorSetting = AUTO

// TrueColorSetting defines the true color usage setting to be used
var TrueColorSetting = AUTO

// seq is the ANSI escape sequence used for the coloring
const seq = "\x1b"

// Color is just use to have a custom type to handle 32-bit RGB definitions
type Color uint32

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

// Pink colors
const (
	Pink            = Color(0x00FFC0CB)
	LightPink       = Color(0x00FFB6C1)
	HotPink         = Color(0x00FF69B4)
	DeepPink        = Color(0x00FF1493)
	PaleVioletRed   = Color(0x00DB7093)
	MediumVioletRed = Color(0x00C71585)
)

// Red colors
const (
	LightSalmon = Color(0x00FFA07A)
	Salmon      = Color(0x00FA8072)
	DarkSalmon  = Color(0x00E9967A)
	LightCoral  = Color(0x00F08080)
	IndianRed   = Color(0x00CD5C5C)
	Crimson     = Color(0x00DC143C)
	FireBrick   = Color(0x00B22222)
	DarkRed     = Color(0x008B0000)
	Red         = Color(0x00FF0000)
)

// Orange colors
const (
	OrangeRed  = Color(0x00FF4500)
	Tomato     = Color(0x00FF6347)
	Coral      = Color(0x00FF7F50)
	DarkOrange = Color(0x00FF8C00)
	Orange     = Color(0x00FFA500)
)

// Yellow colors
const (
	Yellow               = Color(0x00FFFF00)
	LightYellow          = Color(0x00FFFFE0)
	LemonChiffon         = Color(0x00FFFACD)
	LightGoldenrodYellow = Color(0x00FAFAD2)
	PapayaWhip           = Color(0x00FFEFD5)
	Moccasin             = Color(0x00FFE4B5)
	PeachPuff            = Color(0x00FFDAB9)
	PaleGoldenrod        = Color(0x00EEE8AA)
	Khaki                = Color(0x00F0E68C)
	DarkKhaki            = Color(0x00BDB76B)
	Gold                 = Color(0x00FFD700)
)

// Brown colors
const (
	Cornsilk       = Color(0x00FFF8DC)
	BlanchedAlmond = Color(0x00FFEBCD)
	Bisque         = Color(0x00FFE4C4)
	NavajoWhite    = Color(0x00FFDEAD)
	Wheat          = Color(0x00F5DEB3)
	BurlyWood      = Color(0x00DEB887)
	Tan            = Color(0x00D2B48C)
	RosyBrown      = Color(0x00BC8F8F)
	SandyBrown     = Color(0x00F4A460)
	Goldenrod      = Color(0x00DAA520)
	DarkGoldenrod  = Color(0x00B8860B)
	Peru           = Color(0x00CD853F)
	Chocolate      = Color(0x00D2691E)
	SaddleBrown    = Color(0x008B4513)
	Sienna         = Color(0x00A0522D)
	Brown          = Color(0x00A52A2A)
	Maroon         = Color(0x00800000)
)

// Green colors
const (
	DarkOliveGreen    = Color(0x00556B2F)
	Olive             = Color(0x00808000)
	OliveDrab         = Color(0x006B8E23)
	YellowGreen       = Color(0x009ACD32)
	LimeGreen         = Color(0x0032CD32)
	Lime              = Color(0x0000FF00)
	LawnGreen         = Color(0x007CFC00)
	Chartreuse        = Color(0x007FFF00)
	GreenYellow       = Color(0x00ADFF2F)
	SpringGreen       = Color(0x0000FF7F)
	MediumSpringGreen = Color(0x0000FA9A)
	LightGreen        = Color(0x0090EE90)
	PaleGreen         = Color(0x0098FB98)
	DarkSeaGreen      = Color(0x008FBC8F)
	MediumAquamarine  = Color(0x0066CDAA)
	MediumSeaGreen    = Color(0x003CB371)
	SeaGreen          = Color(0x002E8B57)
	ForestGreen       = Color(0x00228B22)
	Green             = Color(0x00008000)
	DarkGreen         = Color(0x00006400)
)

// Cyan colors
const (
	Aqua            = Color(0x0000FFFF)
	Cyan            = Color(0x0000FFFF)
	LightCyan       = Color(0x00E0FFFF)
	PaleTurquoise   = Color(0x00AFEEEE)
	Aquamarine      = Color(0x007FFFD4)
	Turquoise       = Color(0x0040E0D0)
	MediumTurquoise = Color(0x0048D1CC)
	DarkTurquoise   = Color(0x0000CED1)
	LightSeaGreen   = Color(0x0020B2AA)
	CadetBlue       = Color(0x005F9EA0)
	DarkCyan        = Color(0x00008B8B)
	Teal            = Color(0x00008080)
)

// Blue colors
const (
	LightSteelBlue = Color(0x00B0C4DE)
	PowderBlue     = Color(0x00B0E0E6)
	LightBlue      = Color(0x00ADD8E6)
	SkyBlue        = Color(0x0087CEEB)
	LightSkyBlue   = Color(0x0087CEFA)
	DeepSkyBlue    = Color(0x0000BFFF)
	DodgerBlue     = Color(0x001E90FF)
	CornflowerBlue = Color(0x006495ED)
	SteelBlue      = Color(0x004682B4)
	RoyalBlue      = Color(0x004169E1)
	Blue           = Color(0x000000FF)
	MediumBlue     = Color(0x000000CD)
	DarkBlue       = Color(0x0000008B)
	Navy           = Color(0x00000080)
	MidnightBlue   = Color(0x00191970)
)

// Purple, violet, and magenta colors
const (
	Lavender        = Color(0x00E6E6FA)
	Thistle         = Color(0x00D8BFD8)
	Plum            = Color(0x00DDA0DD)
	Violet          = Color(0x00EE82EE)
	Orchid          = Color(0x00DA70D6)
	Fuchsia         = Color(0x00FF00FF)
	Magenta         = Color(0x00FF00FF)
	MediumOrchid    = Color(0x00BA55D3)
	MediumPurple    = Color(0x009370DB)
	BlueViolet      = Color(0x008A2BE2)
	DarkViolet      = Color(0x009400D3)
	DarkOrchid      = Color(0x009932CC)
	DarkMagenta     = Color(0x008B008B)
	Purple          = Color(0x00800080)
	Indigo          = Color(0x004B0082)
	DarkSlateBlue   = Color(0x00483D8B)
	SlateBlue       = Color(0x006A5ACD)
	MediumSlateBlue = Color(0x007B68EE)
)

// White colors
const (
	White         = Color(0x00FFFFFF)
	Snow          = Color(0x00FFFAFA)
	Honeydew      = Color(0x00F0FFF0)
	MintCream     = Color(0x00F5FFFA)
	Azure         = Color(0x00F0FFFF)
	AliceBlue     = Color(0x00F0F8FF)
	GhostWhite    = Color(0x00F8F8FF)
	WhiteSmoke    = Color(0x00F5F5F5)
	Seashell      = Color(0x00FFF5EE)
	Beige         = Color(0x00F5F5DC)
	OldLace       = Color(0x00FDF5E6)
	FloralWhite   = Color(0x00FFFAF0)
	Ivory         = Color(0x00FFFFF0)
	AntiqueWhite  = Color(0x00FAEBD7)
	Linen         = Color(0x00FAF0E6)
	LavenderBlush = Color(0x00FFF0F5)
	MistyRose     = Color(0x00FFE4E1)
)

// Gray and black colors
const (
	Gainsboro      = Color(0x00DCDCDC)
	LightGray      = Color(0x00D3D3D3)
	Silver         = Color(0x00C0C0C0)
	DarkGray       = Color(0x00A9A9A9)
	Gray           = Color(0x00808080)
	DimGray        = Color(0x00696969)
	LightSlateGray = Color(0x00778899)
	SlateGray      = Color(0x00708090)
	DarkSlateGray  = Color(0x002F4F4F)
	Black          = Color(0x00000000)
)

// Special additional colors
const (
	AdditionGreen      = Color(0x0058BF38)
	RemovalRed         = Color(0x00B9311B)
	ModificationYellow = Color(0x00C7C43F)
)

// Modifiers
const (
	Bold      = 1
	Italic    = 3
	Underline = 4
)

// translateMap4bitColor is an internal helper map to translate specific well
// defined colors to matching 4-bit color attributes
var translateMap4bitColor = map[Color]Attribute{
	Color(0x00000000): Attribute(30),
	Color(0x00AA0000): Attribute(31),
	Color(0x0000AA00): Attribute(32),
	Color(0x00FFFF00): Attribute(33),
	Color(0x000000AA): Attribute(34),
	Color(0x00AA00AA): Attribute(35),
	Color(0x0000AAAA): Attribute(36),
	Color(0x00AAAAAA): Attribute(37),
	Color(0x00555555): Attribute(90),
	Color(0x00FF5555): Attribute(91),
	Color(0x0055FF55): Attribute(92),
	Color(0x00FFFF55): Attribute(93),
	Color(0x005555FF): Attribute(94),
	Color(0x00FF55FF): Attribute(95),
	Color(0x0055FFFF): Attribute(96),
	Color(0x00FFFFFF): Attribute(97),
}

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
		(ColorSetting == AUTO && !isDumbTerminal && isTerminal)
}

// UseTrueColor returns whether true color colors should be used or not base on
// the configured true color usage setting
func UseTrueColor() bool {
	return (TrueColorSetting == ON) ||
		(TrueColorSetting == AUTO && isTrueColor)
}

// Colorize applies an ANSI truecolor sequence for the provided color to the given text.
func Colorize(text string, color Color, modifiers ...Attribute) string {
	modifiers = keepAttributes(modifiers, []Attribute{Bold, Italic, Underline})
	sort.Slice(modifiers, func(i, j int) bool {
		return modifiers[i] < modifiers[j]
	})

	if UseTrueColor() {
		r, g, b := BreakUpColorIntoChannels(color)
		colorCoding := []Attribute{38, 2, Attribute(r), Attribute(g), Attribute(b)}
		return wrapTextInSeq(text, append(modifiers, colorCoding...)...)
	}

	colorAttribute := Get4bitEquivalentColorAttribute(color)
	return wrapTextInSeq(text, append(modifiers, colorAttribute)...)
}

// ColorizeFgBg applies an ANSI truecolor sequence for the provided foreground and background colors to the given text.
func ColorizeFgBg(text string, foreground Color, background Color, modifiers ...Attribute) string {
	modifiers = keepAttributes(modifiers, []Attribute{Bold, Italic, Underline})
	sort.Slice(modifiers, func(i, j int) bool {
		return modifiers[i] < modifiers[j]
	})

	fgr, fgg, fgb := BreakUpColorIntoChannels(foreground)
	fgColorCoding := []Attribute{38, 2, Attribute(fgr), Attribute(fgg), Attribute(fgb)}

	bgr, bgg, bgb := BreakUpColorIntoChannels(background)
	bgColorCoding := []Attribute{48, 2, Attribute(bgr), Attribute(bgg), Attribute(bgb)}

	return wrapTextInSeq(text, append(append(modifiers, bgColorCoding...), fgColorCoding...)...)
}

// Style applies only text modifications like Bold, Italic, or Underline to the text
func Style(text string, modifiers ...Attribute) string {
	return wrapTextInSeq(text, keepAttributes(modifiers, []Attribute{Bold, Italic, Underline})...)
}

// ColorizeWithAttributes applies the provided attributes with any filtering or checks
func ColorizeWithAttributes(text string, attributes ...Attribute) string {
	return wrapTextInSeq(text, attributes...)
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
func BreakUpColorIntoChannels(color Color) (uint8, uint8, uint8) {
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
func Get4bitEquivalentColorAttribute(targetColor Color) Attribute {
	red, green, blue := BreakUpColorIntoChannels(targetColor)
	target := &color.RGBA{red, green, blue, 0xFF}

	min := math.MaxFloat64
	result := Attribute(0)

	// Calculate the distance between the target color and the available 4-bit
	// colors using the `deltaE` algorithm to find the best match.
	for candidate, attribute := range translateMap4bitColor {
		r, g, b := BreakUpColorIntoChannels(candidate)
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
	if err != nil { // TODO Debug, or trace the error
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
