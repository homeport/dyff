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

The `ToString` function returns neat looking YAML string output using text
highlighting with emphasis, colors, and indent helper guide lines to create
pleasing and easy to read YAML.
*/
package neat

import (
	"fmt"
	"strings"

	"github.com/HeavyWombat/dyff/pkg/v1/bunt"
	yaml "gopkg.in/yaml.v2"
)

// ToYAMLString marshals the provided object into YAML with text decorations
// and is basically just a convenience function to create the output processor
// and call its `ToYAML` function.
func ToYAMLString(obj interface{}) (string, error) {
	return NewOutputProcessor(true, true, &DefaultColorSchema).ToYAML(obj)
}

// ToYAML processes the provided input object and tries to neatly output it as
// human readable YAML honoring the preferences provided to the output processor
func (p *OutputProcessor) ToYAML(obj interface{}) (string, error) {
	if err := p.neatYAML("", false, obj); err != nil {
		return "", err
	}

	p.out.Flush()
	return p.data.String(), nil
}

func (p *OutputProcessor) neatYAML(prefix string, skipIndentOnFirstLine bool, obj interface{}) error {
	switch obj.(type) {
	case yaml.MapSlice:
		if err := p.neatYAMLofMapSlice(prefix, skipIndentOnFirstLine, obj.(yaml.MapSlice)); err != nil {
			return err
		}

	case []interface{}:
		if err := p.neatYAMLofSlice(prefix, skipIndentOnFirstLine, obj.([]interface{})); err != nil {
			return err
		}

	case []yaml.MapSlice:
		if err := p.neatYAMLofSlice(prefix, skipIndentOnFirstLine, p.simplify(obj.([]yaml.MapSlice))); err != nil {
			return err
		}

	default:
		if err := p.neatYAMLofScalar(prefix, skipIndentOnFirstLine, obj); err != nil {
			return err
		}
	}

	return nil
}

func (p *OutputProcessor) neatYAMLofMapSlice(prefix string, skipIndentOnFirstLine bool, mapslice yaml.MapSlice) error {
	for i, mapitem := range mapslice {
		if !skipIndentOnFirstLine || i > 0 {
			p.out.WriteString(prefix)
		}

		keyString := fmt.Sprintf("%v:", mapitem.Key)
		if p.boldKeys {
			keyString = bunt.Style(keyString, bunt.Bold)
		}

		p.out.WriteString(p.colorize(keyString, "keyColor"))

		switch mapitem.Value.(type) {
		case yaml.MapSlice:
			if len(mapitem.Value.(yaml.MapSlice)) == 0 {
				p.out.WriteString(" ")
				p.out.WriteString(p.colorize("{}", "emptyStructures"))
				p.out.WriteString("\n")

			} else {
				p.out.WriteString("\n")
				if err := p.neatYAMLofMapSlice(prefix+p.prefixAdd(), false, mapitem.Value.(yaml.MapSlice)); err != nil {
					return err
				}
			}

		case []interface{}:
			if len(mapitem.Value.([]interface{})) == 0 {
				p.out.WriteString(" ")
				p.out.WriteString(p.colorize("[]", "emptyStructures"))
				p.out.WriteString("\n")
			} else {
				p.out.WriteString("\n")
				if err := p.neatYAMLofSlice(prefix, false, mapitem.Value.([]interface{})); err != nil {
					return err
				}
			}

		default:
			p.out.WriteString(" ")
			if err := p.neatYAMLofScalar(prefix, false, mapitem.Value); err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *OutputProcessor) neatYAMLofSlice(prefix string, skipIndentOnFirstLine bool, list []interface{}) error {
	for _, entry := range list {
		p.out.WriteString(prefix)
		p.out.WriteString(p.colorize("-", "dashColor"))
		p.out.WriteString(" ")
		if err := p.neatYAML(prefix+p.prefixAdd(), true, entry); err != nil {
			return err
		}
	}

	return nil
}

func (p *OutputProcessor) neatYAMLofScalar(prefix string, skipIndentOnFirstLine bool, obj interface{}) error {
	// Process nil values immediately and return afterwards
	if obj == nil {
		p.out.WriteString(p.colorize("null", "nullColor"))
		p.out.WriteString("\n")
		return nil
	}

	// Any other value: Run through Go YAML marshaller and colorize afterwards
	data, err := yaml.Marshal(obj)
	if err != nil {
		return err
	}

	// Decide on one color to be used
	color := p.determineColorByType(obj)

	// Cast byte slice to string, remove trailing newlines, split into lines
	for i, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		if i > 0 {
			p.out.WriteString(prefix)
		}

		p.out.WriteString(p.colorize(line, color))
		p.out.WriteString("\n")
	}

	return nil
}
