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

package neat

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"

	"github.com/HeavyWombat/dyff/pkg/bunt"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

var (
	keyColor           = bunt.Coral
	indentLineColor    = bunt.Color(0x00242424)
	scalarDefaultColor = bunt.PaleGreen
	boolColor          = bunt.Moccasin
	floatColor         = bunt.Orange
	intColor           = bunt.MediumPurple
	multiLineTextColor = bunt.Aquamarine
	documentStartColor = bunt.Goldenrod
)

// ToYAMLString marshals the provided object into YAML with text decorations
func ToYAMLString(obj interface{}, plainYAML bool, note string) (string, error) {
	if plainYAML {
		// Use default YAML marshaling
		output, err := yaml.Marshal(obj)
		if err != nil {
			return "", errors.Wrap(err, fmt.Sprintf("Failed to marshal input object %#v", obj))
		}

		return fmt.Sprintf("---\n%s\n", string(output)), nil
	}

	// Use internal custom YAML marshaling with colors
	buf := &bytes.Buffer{}
	writer := bufio.NewWriter(buf)

	writer.WriteString(bunt.Colorize("---", documentStartColor, bunt.Bold))
	if len(note) > 0 {
		writer.WriteString(" # ")
		writer.WriteString(note)
	}
	writer.WriteString("\n")

	if err := neat(writer, "", false, obj); err != nil {
		return "", err
	}

	writer.WriteString("\n")
	writer.Flush()

	return buf.String(), nil
}

func neat(out *bufio.Writer, prefix string, skipIndentOnFirstLine bool, obj interface{}) error {
	switch obj.(type) {
	case yaml.MapSlice:
		if err := neatMapSlice(out, prefix, skipIndentOnFirstLine, obj.(yaml.MapSlice)); err != nil {
			return err
		}

	case []interface{}:
		if err := neatSlice(out, prefix, skipIndentOnFirstLine, obj.([]interface{})); err != nil {
			return err
		}

	case []yaml.MapSlice:
		if err := neatMapSliceSlice(out, prefix, skipIndentOnFirstLine, obj.([]yaml.MapSlice)); err != nil {
			return err
		}

	default:
		if err := neatScalar(out, prefix, skipIndentOnFirstLine, obj); err != nil {
			return err
		}
	}

	return nil
}

func neatMapSlice(out *bufio.Writer, prefix string, skipIndentOnFirstLine bool, mapslice yaml.MapSlice) error {
	for i, mapitem := range mapslice {
		if !skipIndentOnFirstLine || i > 0 {
			out.WriteString(prefix)
		}

		out.WriteString(bunt.Colorize(fmt.Sprintf("%v:", mapitem.Key), keyColor, bunt.Bold))

		switch mapitem.Value.(type) {
		case yaml.MapSlice:
			out.WriteString("\n")
			if err := neatMapSlice(out, prefix+prefixAdd(), false, mapitem.Value.(yaml.MapSlice)); err != nil {
				return err
			}

		case []interface{}:
			out.WriteString("\n")
			if err := neatSlice(out, prefix, false, mapitem.Value.([]interface{})); err != nil {
				return err
			}

		default:
			out.WriteString(" ")
			if err := neatScalar(out, prefix, false, mapitem.Value); err != nil {
				return err
			}
		}
	}

	return nil
}

func neatSlice(out *bufio.Writer, prefix string, skipIndentOnFirstLine bool, list []interface{}) error {
	for _, entry := range list {
		out.WriteString(prefix)
		out.WriteString(bunt.Style("- ", bunt.Bold))
		if err := neat(out, prefix+prefixAdd(), true, entry); err != nil {
			return err
		}
	}

	return nil
}

func neatMapSliceSlice(out *bufio.Writer, prefix string, skipIndentOnFirstLine bool, list []yaml.MapSlice) error {
	for _, entry := range list {
		out.WriteString(prefix)
		out.WriteString(bunt.Style("- ", bunt.Bold))
		if err := neat(out, prefix+prefixAdd(), true, entry); err != nil {
			return err
		}
	}

	return nil
}

func neatScalar(out *bufio.Writer, prefix string, skipIndentOnFirstLine bool, obj interface{}) error {
	data, err := yaml.Marshal(obj)
	if err != nil {
		return err
	}

	color := scalarDefaultColor
	switch obj.(type) {
	case bool:
		color = boolColor

	case float32, float64:
		color = floatColor

	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr:
		color = intColor
	}

	// Cast byte slice to string, remove trailing newlines, split into lines
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")

	if len(lines) > 1 {
		color = multiLineTextColor
	}

	for i, line := range lines {
		if i > 0 {
			out.WriteString(prefix)
		}

		out.WriteString(bunt.Colorize(line, color))
		out.WriteString("\n")
	}

	return nil
}

func prefixAdd() string {
	return bunt.Colorize("│ ", indentLineColor)
}
