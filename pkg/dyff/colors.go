// Copyright © 2019 The Homeport Team
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

package dyff

import (
	"fmt"

	"github.com/gonvenience/bunt"
	"github.com/gonvenience/neat"
	"github.com/lucasb-eyer/go-colorful"
)

// Colorizer is the interface for applying colors to diff output. Implement this
// to replace the default bunt-based coloring with an alternative (e.g. ANSI).
type Colorizer interface {
	Green(text string) string
	Greenf(format string, a ...interface{}) string
	Red(text string) string
	Redf(format string, a ...interface{}) string
	Yellowf(format string, a ...interface{}) string
	LightGreen(text string) string
	LightRed(text string) string
	DimGray(text string) string
	Bold(text string) string
	Italic(text string) string
	BoldGreen(text string) string
	BoldRed(text string) string
	StylizeHeader(header string) string
	YAMLInRedishColors(input interface{}, useIndentLines bool) (string, error)
	YAMLInGreenishColors(input interface{}, useIndentLines bool) (string, error)
}

// BuntColorizer is the default Colorizer using bunt and go-colorful.
type BuntColorizer struct{}

var (
	additionGreen      = hexColor("#58BF38")
	modificationYellow = hexColor("#C7C43F")
	removalRed         = hexColor("#B9311B")
)

func hexColor(hex string) colorful.Color {
	c, _ := colorful.Hex(hex)
	return c
}

func render(format string, a ...interface{}) string {
	if len(a) == 0 {
		return format
	}

	return fmt.Sprintf(format, a...)
}

func colored(color colorful.Color, text string) string {
	return bunt.Style(text,
		bunt.EachLine(),
		bunt.Foreground(color),
	)
}

func coloredf(color colorful.Color, format string, a ...interface{}) string {
	return bunt.Style(render(format, a...),
		bunt.EachLine(),
		bunt.Foreground(color),
	)
}

func defaultColorizer() Colorizer {
	return &BuntColorizer{}
}

func (c *BuntColorizer) Green(text string) string {
	return colored(additionGreen, text)
}

func (c *BuntColorizer) Greenf(format string, a ...interface{}) string {
	return coloredf(additionGreen, format, a...)
}

func (c *BuntColorizer) Red(text string) string {
	return colored(removalRed, text)
}

func (c *BuntColorizer) Redf(format string, a ...interface{}) string {
	return coloredf(removalRed, format, a...)
}

func (c *BuntColorizer) Yellowf(format string, a ...interface{}) string {
	return coloredf(modificationYellow, format, a...)
}

func (c *BuntColorizer) LightGreen(text string) string {
	return colored(bunt.LightGreen, text)
}

func (c *BuntColorizer) LightRed(text string) string {
	return colored(bunt.LightSalmon, text)
}

func (c *BuntColorizer) DimGray(text string) string {
	return colored(bunt.DimGray, text)
}

func (c *BuntColorizer) Bold(text string) string {
	return bunt.Style(text,
		bunt.EachLine(),
		bunt.Bold(),
	)
}

func (c *BuntColorizer) Italic(text string) string {
	return bunt.Style(text,
		bunt.EachLine(),
		bunt.Italic(),
	)
}

func (c *BuntColorizer) BoldGreen(text string) string {
	return c.Bold(c.Green(text))
}

func (c *BuntColorizer) BoldRed(text string) string {
	return c.Bold(c.Red(text))
}

func (c *BuntColorizer) StylizeHeader(header string) string {
	return bunt.Style(
		header,
		bunt.ForegroundFunc(func(x int, _ int, _ rune) *colorful.Color {
			switch {
			case x < 7:
				return &colorful.Color{R: .45, G: .71, B: .30}
			case x < 13:
				return &colorful.Color{R: .79, G: .76, B: .38}
			case x < 21:
				return &colorful.Color{R: .65, G: .17, B: .17}
			}
			return nil
		}),
	)
}

func (c *BuntColorizer) YAMLInRedishColors(input interface{}, useIndentLines bool) (string, error) {
	return neat.NewOutputProcessor(useIndentLines, true, &map[string]colorful.Color{
		"keyColor":           bunt.FireBrick,
		"indentLineColor":    {R: 0.2, G: 0, B: 0},
		"scalarDefaultColor": bunt.LightCoral,
		"boolColor":          bunt.LightCoral,
		"floatColor":         bunt.LightCoral,
		"intColor":           bunt.LightCoral,
		"multiLineTextColor": bunt.DarkSalmon,
		"nullColor":          bunt.Salmon,
		"emptyStructures":    bunt.LightSalmon,
		"dashColor":          bunt.FireBrick,
	}).ToYAML(input)
}

func (c *BuntColorizer) YAMLInGreenishColors(input interface{}, useIndentLines bool) (string, error) {
	return neat.NewOutputProcessor(useIndentLines, true, &map[string]colorful.Color{
		"keyColor":           bunt.Green,
		"indentLineColor":    {R: 0, G: 0.2, B: 0},
		"scalarDefaultColor": bunt.LimeGreen,
		"boolColor":          bunt.LimeGreen,
		"floatColor":         bunt.LimeGreen,
		"intColor":           bunt.LimeGreen,
		"multiLineTextColor": bunt.OliveDrab,
		"nullColor":          bunt.Olive,
		"emptyStructures":    bunt.DarkOliveGreen,
		"dashColor":          bunt.Green,
	}).ToYAML(input)
}
