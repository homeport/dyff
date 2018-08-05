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

/*
Package neat provides a YAML Marshaller that supports colors.

The `ToYAML` function returns neat looking YAML string output using text
highlighting with emphasis, colors, and indent helper guide lines to create
pleasing and easy to read YAML.
*/
package neat

import (
	"bufio"
	"bytes"
	"strings"

	"github.com/HeavyWombat/dyff/pkg/v1/bunt"
	colorful "github.com/lucasb-eyer/go-colorful"
	yaml "gopkg.in/yaml.v2"
)

// DefaultColorSchema is a prepared usable color schema for the neat output
// processor which is loosly based upon the colors used by Atom
var DefaultColorSchema = map[string]colorful.Color{
	"keyColor":           bunt.IndianRed,
	"indentLineColor":    {R: 0.14, G: 0.14, B: 0.14},
	"scalarDefaultColor": bunt.PaleGreen,
	"boolColor":          bunt.Moccasin,
	"floatColor":         bunt.Orange,
	"intColor":           bunt.MediumPurple,
	"multiLineTextColor": bunt.Aquamarine,
	"nullColor":          bunt.DarkOrange,
	"emptyStructures":    bunt.PaleGoldenrod,
}

// OutputProcessor provides the functionality to output neat YAML strings using
// colors and text emphasis
type OutputProcessor struct {
	data           *bytes.Buffer
	out            *bufio.Writer
	colorSchema    *map[string]colorful.Color
	useIndentLines bool
	boldKeys       bool
}

// NewOutputProcessor creates a new output processor including the required
// internals using the provided preferences
func NewOutputProcessor(useIndentLines bool, boldKeys bool, colorSchema *map[string]colorful.Color) *OutputProcessor {
	bytesBuffer := &bytes.Buffer{}
	writer := bufio.NewWriter(bytesBuffer)

	// Only use indent lines in color mode
	if !bunt.UseColors() {
		useIndentLines = false
	}

	return &OutputProcessor{
		data:           bytesBuffer,
		out:            writer,
		useIndentLines: useIndentLines,
		boldKeys:       boldKeys,
		colorSchema:    colorSchema,
	}
}

func (p *OutputProcessor) colorize(text string, colorName string) string {
	if p.colorSchema != nil {
		if value, ok := (*p.colorSchema)[colorName]; ok {
			return bunt.Colorize(text, value)
		}
	}

	return text
}

func (p *OutputProcessor) determineColorByType(obj interface{}) string {
	color := "scalarDefaultColor"

	switch obj.(type) {
	case bool:
		color = "boolColor"

	case float32, float64:
		color = "floatColor"

	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr:
		color = "intColor"

	case string:
		if len(strings.Split(strings.TrimSpace(obj.(string)), "\n")) > 1 {
			color = "multiLineTextColor"
		}
	}

	return color
}

func (p *OutputProcessor) isScalar(obj interface{}) bool {
	switch obj.(type) {
	case yaml.MapSlice, []interface{}, []yaml.MapSlice:
		return false

	default:
		return true
	}
}

func (p *OutputProcessor) simplify(list []yaml.MapSlice) []interface{} {
	result := make([]interface{}, len(list))
	for idx, value := range list {
		result[idx] = value
	}

	return result
}

func (p *OutputProcessor) prefixAdd() string {
	if p.useIndentLines {
		return p.colorize("│ ", "indentLineColor")
	}

	return p.colorize("  ", "indentLineColor")
}
