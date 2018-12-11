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

package bunt

import (
	colorful "github.com/lucasb-eyer/go-colorful"
)

// Pink colors
var (
	Pink            = hexColor("#FFC0CB")
	LightPink       = hexColor("#FFB6C1")
	HotPink         = hexColor("#FF69B4")
	DeepPink        = hexColor("#FF1493")
	PaleVioletRed   = hexColor("#DB7093")
	MediumVioletRed = hexColor("#C71585")
)

// Red colors
var (
	LightSalmon = hexColor("#FFA07A")
	Salmon      = hexColor("#FA8072")
	DarkSalmon  = hexColor("#E9967A")
	LightCoral  = hexColor("#F08080")
	IndianRed   = hexColor("#CD5C5C")
	Crimson     = hexColor("#DC143C")
	FireBrick   = hexColor("#B22222")
	DarkRed     = hexColor("#8B0000")
	Red         = hexColor("#FF0000")
)

// Orange colors
var (
	OrangeRed  = hexColor("#FF4500")
	Tomato     = hexColor("#FF6347")
	Coral      = hexColor("#FF7F50")
	DarkOrange = hexColor("#FF8C00")
	Orange     = hexColor("#FFA500")
)

// Yellow colors
var (
	Yellow               = hexColor("#FFFF00")
	LightYellow          = hexColor("#FFFFE0")
	LemonChiffon         = hexColor("#FFFACD")
	LightGoldenrodYellow = hexColor("#FAFAD2")
	PapayaWhip           = hexColor("#FFEFD5")
	Moccasin             = hexColor("#FFE4B5")
	PeachPuff            = hexColor("#FFDAB9")
	PaleGoldenrod        = hexColor("#EEE8AA")
	Khaki                = hexColor("#F0E68C")
	DarkKhaki            = hexColor("#BDB76B")
	Gold                 = hexColor("#FFD700")
)

// Brown colors
var (
	Cornsilk       = hexColor("#FFF8DC")
	BlanchedAlmond = hexColor("#FFEBCD")
	Bisque         = hexColor("#FFE4C4")
	NavajoWhite    = hexColor("#FFDEAD")
	Wheat          = hexColor("#F5DEB3")
	BurlyWood      = hexColor("#DEB887")
	Tan            = hexColor("#D2B48C")
	RosyBrown      = hexColor("#BC8F8F")
	SandyBrown     = hexColor("#F4A460")
	Goldenrod      = hexColor("#DAA520")
	DarkGoldenrod  = hexColor("#B8860B")
	Peru           = hexColor("#CD853F")
	Chocolate      = hexColor("#D2691E")
	SaddleBrown    = hexColor("#8B4513")
	Sienna         = hexColor("#A0522D")
	Brown          = hexColor("#A52A2A")
	Maroon         = hexColor("#800000")
)

// Green colors
var (
	DarkOliveGreen    = hexColor("#556B2F")
	Olive             = hexColor("#808000")
	OliveDrab         = hexColor("#6B8E23")
	YellowGreen       = hexColor("#9ACD32")
	LimeGreen         = hexColor("#32CD32")
	Lime              = hexColor("#00FF00")
	LawnGreen         = hexColor("#7CFC00")
	Chartreuse        = hexColor("#7FFF00")
	GreenYellow       = hexColor("#ADFF2F")
	SpringGreen       = hexColor("#00FF7F")
	MediumSpringGreen = hexColor("#00FA9A")
	LightGreen        = hexColor("#90EE90")
	PaleGreen         = hexColor("#98FB98")
	DarkSeaGreen      = hexColor("#8FBC8F")
	MediumAquamarine  = hexColor("#66CDAA")
	MediumSeaGreen    = hexColor("#3CB371")
	SeaGreen          = hexColor("#2E8B57")
	ForestGreen       = hexColor("#228B22")
	Green             = hexColor("#008000")
	DarkGreen         = hexColor("#006400")
)

// Cyan colors
var (
	Aqua            = hexColor("#00FFFF")
	Cyan            = hexColor("#00FFFF")
	LightCyan       = hexColor("#E0FFFF")
	PaleTurquoise   = hexColor("#AFEEEE")
	Aquamarine      = hexColor("#7FFFD4")
	Turquoise       = hexColor("#40E0D0")
	MediumTurquoise = hexColor("#48D1CC")
	DarkTurquoise   = hexColor("#00CED1")
	LightSeaGreen   = hexColor("#20B2AA")
	CadetBlue       = hexColor("#5F9EA0")
	DarkCyan        = hexColor("#008B8B")
	Teal            = hexColor("#008080")
)

// Blue colors
var (
	LightSteelBlue = hexColor("#B0C4DE")
	PowderBlue     = hexColor("#B0E0E6")
	LightBlue      = hexColor("#ADD8E6")
	SkyBlue        = hexColor("#87CEEB")
	LightSkyBlue   = hexColor("#87CEFA")
	DeepSkyBlue    = hexColor("#00BFFF")
	DodgerBlue     = hexColor("#1E90FF")
	CornflowerBlue = hexColor("#6495ED")
	SteelBlue      = hexColor("#4682B4")
	RoyalBlue      = hexColor("#4169E1")
	Blue           = hexColor("#0000FF")
	MediumBlue     = hexColor("#0000CD")
	DarkBlue       = hexColor("#00008B")
	Navy           = hexColor("#000080")
	MidnightBlue   = hexColor("#191970")
)

// Purple, violet, and magenta colors
var (
	Lavender        = hexColor("#E6E6FA")
	Thistle         = hexColor("#D8BFD8")
	Plum            = hexColor("#DDA0DD")
	Violet          = hexColor("#EE82EE")
	Orchid          = hexColor("#DA70D6")
	Fuchsia         = hexColor("#FF00FF")
	Magenta         = hexColor("#FF00FF")
	MediumOrchid    = hexColor("#BA55D3")
	MediumPurple    = hexColor("#9370DB")
	BlueViolet      = hexColor("#8A2BE2")
	DarkViolet      = hexColor("#9400D3")
	DarkOrchid      = hexColor("#9932CC")
	DarkMagenta     = hexColor("#8B008B")
	Purple          = hexColor("#800080")
	Indigo          = hexColor("#4B0082")
	DarkSlateBlue   = hexColor("#483D8B")
	SlateBlue       = hexColor("#6A5ACD")
	MediumSlateBlue = hexColor("#7B68EE")
)

// White colors
var (
	White         = hexColor("#FFFFFF")
	Snow          = hexColor("#FFFAFA")
	Honeydew      = hexColor("#F0FFF0")
	MintCream     = hexColor("#F5FFFA")
	Azure         = hexColor("#F0FFFF")
	AliceBlue     = hexColor("#F0F8FF")
	GhostWhite    = hexColor("#F8F8FF")
	WhiteSmoke    = hexColor("#F5F5F5")
	Seashell      = hexColor("#FFF5EE")
	Beige         = hexColor("#F5F5DC")
	OldLace       = hexColor("#FDF5E6")
	FloralWhite   = hexColor("#FFFAF0")
	Ivory         = hexColor("#FFFFF0")
	AntiqueWhite  = hexColor("#FAEBD7")
	Linen         = hexColor("#FAF0E6")
	LavenderBlush = hexColor("#FFF0F5")
	MistyRose     = hexColor("#FFE4E1")
)

// Gray and black colors
var (
	Gainsboro      = hexColor("#DCDCDC")
	LightGray      = hexColor("#D3D3D3")
	Silver         = hexColor("#C0C0C0")
	DarkGray       = hexColor("#A9A9A9")
	Gray           = hexColor("#808080")
	DimGray        = hexColor("#696969")
	LightSlateGray = hexColor("#778899")
	SlateGray      = hexColor("#708090")
	DarkSlateGray  = hexColor("#2F4F4F")
	Black          = hexColor("#000000")
)

// Special additional colors
var (
	AdditionGreen      = hexColor("#58BF38")
	RemovalRed         = hexColor("#B9311B")
	ModificationYellow = hexColor("#C7C43F")
)

// translateMap4bitColor is an internal helper map to translate specific well
// defined colors to matching 4-bit color attributes
var translateMap4bitColor = map[colorful.Color]Attribute{
	hexColor("#000000"): Attribute(30),
	hexColor("#AA0000"): Attribute(31),
	hexColor("#00AA00"): Attribute(32),
	hexColor("#FFFF00"): Attribute(33),
	hexColor("#0000AA"): Attribute(34),
	hexColor("#AA00AA"): Attribute(35),
	hexColor("#00AAAA"): Attribute(36),
	hexColor("#AAAAAA"): Attribute(37),
	hexColor("#555555"): Attribute(90),
	hexColor("#FF5555"): Attribute(91),
	hexColor("#55FF55"): Attribute(92),
	hexColor("#FFFF55"): Attribute(93),
	hexColor("#5555FF"): Attribute(94),
	hexColor("#FF55FF"): Attribute(95),
	hexColor("#55FFFF"): Attribute(96),
	hexColor("#FFFFFF"): Attribute(97),
}

var colorByNameMap = map[string]colorful.Color{
	"Pink":                 Pink,
	"LightPink":            LightPink,
	"HotPink":              HotPink,
	"DeepPink":             DeepPink,
	"PaleVioletRed":        PaleVioletRed,
	"MediumVioletRed":      MediumVioletRed,
	"LightSalmon":          LightSalmon,
	"Salmon":               Salmon,
	"DarkSalmon":           DarkSalmon,
	"LightCoral":           LightCoral,
	"IndianRed":            IndianRed,
	"Crimson":              Crimson,
	"FireBrick":            FireBrick,
	"DarkRed":              DarkRed,
	"Red":                  Red,
	"OrangeRed":            OrangeRed,
	"Tomato":               Tomato,
	"Coral":                Coral,
	"DarkOrange":           DarkOrange,
	"Orange":               Orange,
	"Yellow":               Yellow,
	"LightYellow":          LightYellow,
	"LemonChiffon":         LemonChiffon,
	"LightGoldenrodYellow": LightGoldenrodYellow,
	"PapayaWhip":           PapayaWhip,
	"Moccasin":             Moccasin,
	"PeachPuff":            PeachPuff,
	"PaleGoldenrod":        PaleGoldenrod,
	"Khaki":                Khaki,
	"DarkKhaki":            DarkKhaki,
	"Gold":                 Gold,
	"Cornsilk":             Cornsilk,
	"BlanchedAlmond":       BlanchedAlmond,
	"Bisque":               Bisque,
	"NavajoWhite":          NavajoWhite,
	"Wheat":                Wheat,
	"BurlyWood":            BurlyWood,
	"Tan":                  Tan,
	"RosyBrown":            RosyBrown,
	"SandyBrown":           SandyBrown,
	"Goldenrod":            Goldenrod,
	"DarkGoldenrod":        DarkGoldenrod,
	"Peru":                 Peru,
	"Chocolate":            Chocolate,
	"SaddleBrown":          SaddleBrown,
	"Sienna":               Sienna,
	"Brown":                Brown,
	"Maroon":               Maroon,
	"DarkOliveGreen":       DarkOliveGreen,
	"Olive":                Olive,
	"OliveDrab":            OliveDrab,
	"YellowGreen":          YellowGreen,
	"LimeGreen":            LimeGreen,
	"Lime":                 Lime,
	"LawnGreen":            LawnGreen,
	"Chartreuse":           Chartreuse,
	"GreenYellow":          GreenYellow,
	"SpringGreen":          SpringGreen,
	"MediumSpringGreen":    MediumSpringGreen,
	"LightGreen":           LightGreen,
	"PaleGreen":            PaleGreen,
	"DarkSeaGreen":         DarkSeaGreen,
	"MediumAquamarine":     MediumAquamarine,
	"MediumSeaGreen":       MediumSeaGreen,
	"SeaGreen":             SeaGreen,
	"ForestGreen":          ForestGreen,
	"Green":                Green,
	"DarkGreen":            DarkGreen,
	"Aqua":                 Aqua,
	"Cyan":                 Cyan,
	"LightCyan":            LightCyan,
	"PaleTurquoise":        PaleTurquoise,
	"Aquamarine":           Aquamarine,
	"Turquoise":            Turquoise,
	"MediumTurquoise":      MediumTurquoise,
	"DarkTurquoise":        DarkTurquoise,
	"LightSeaGreen":        LightSeaGreen,
	"CadetBlue":            CadetBlue,
	"DarkCyan":             DarkCyan,
	"Teal":                 Teal,
	"LightSteelBlue":       LightSteelBlue,
	"PowderBlue":           PowderBlue,
	"LightBlue":            LightBlue,
	"SkyBlue":              SkyBlue,
	"LightSkyBlue":         LightSkyBlue,
	"DeepSkyBlue":          DeepSkyBlue,
	"DodgerBlue":           DodgerBlue,
	"CornflowerBlue":       CornflowerBlue,
	"SteelBlue":            SteelBlue,
	"RoyalBlue":            RoyalBlue,
	"Blue":                 Blue,
	"MediumBlue":           MediumBlue,
	"DarkBlue":             DarkBlue,
	"Navy":                 Navy,
	"MidnightBlue":         MidnightBlue,
	"Lavender":             Lavender,
	"Thistle":              Thistle,
	"Plum":                 Plum,
	"Violet":               Violet,
	"Orchid":               Orchid,
	"Fuchsia":              Fuchsia,
	"Magenta":              Magenta,
	"MediumOrchid":         MediumOrchid,
	"MediumPurple":         MediumPurple,
	"BlueViolet":           BlueViolet,
	"DarkViolet":           DarkViolet,
	"DarkOrchid":           DarkOrchid,
	"DarkMagenta":          DarkMagenta,
	"Purple":               Purple,
	"Indigo":               Indigo,
	"DarkSlateBlue":        DarkSlateBlue,
	"SlateBlue":            SlateBlue,
	"MediumSlateBlue":      MediumSlateBlue,
	"White":                White,
	"Snow":                 Snow,
	"Honeydew":             Honeydew,
	"MintCream":            MintCream,
	"Azure":                Azure,
	"AliceBlue":            AliceBlue,
	"GhostWhite":           GhostWhite,
	"WhiteSmoke":           WhiteSmoke,
	"Seashell":             Seashell,
	"Beige":                Beige,
	"OldLace":              OldLace,
	"FloralWhite":          FloralWhite,
	"Ivory":                Ivory,
	"AntiqueWhite":         AntiqueWhite,
	"Linen":                Linen,
	"LavenderBlush":        LavenderBlush,
	"MistyRose":            MistyRose,
	"Gainsboro":            Gainsboro,
	"LightGray":            LightGray,
	"Silver":               Silver,
	"DarkGray":             DarkGray,
	"Gray":                 Gray,
	"DimGray":              DimGray,
	"LightSlateGray":       LightSlateGray,
	"SlateGray":            SlateGray,
	"DarkSlateGray":        DarkSlateGray,
	"Black":                Black,
	"AdditionGreen":        AdditionGreen,
	"RemovalRed":           RemovalRed,
	"ModificationYellow":   ModificationYellow,
}

func hexColor(scol string) colorful.Color {
	c, _ := colorful.Hex(scol)
	return c
}

func lookupColorByName(colorName string) *colorful.Color {
	if color, ok := colorByNameMap[colorName]; ok {
		return &color
	}

	return nil
}
