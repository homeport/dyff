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
	"encoding/json"
	"fmt"
	"strings"

	"github.com/HeavyWombat/dyff/pkg/v1/bunt"
	yaml "gopkg.in/yaml.v2"
)

// ToJSONString marshals the provided object into JSON with text decorations
// and is basically just a convenience function to create the output processor
// and call its `ToJSON` function.
func ToJSONString(obj interface{}) (string, error) {
	return NewOutputProcessor(true, true, &DefaultColorSchema).ToCompactJSON(obj)
}

// ToJSON processes the provided input object and tries to neatly output it as
// human readable JSON honoring the preferences provided to the output processor
func (p *OutputProcessor) ToJSON(obj interface{}) (string, error) {
	var out string
	var err error

	if out, err = p.neatJSON("", obj); err != nil {
		return "", err
	}

	return out, nil
}

// ToCompactJSON processed the provided input object and tries to create a as compact
// as possible output.
func (p *OutputProcessor) ToCompactJSON(obj interface{}) (string, error) {
	switch v := obj.(type) {

	case []interface{}:
		result := make([]string, 0)
		for _, i := range v {
			value, err := p.ToCompactJSON(i)
			if err != nil {
				return "", err
			}
			result = append(result, value)
		}

		return fmt.Sprintf("[%s]", strings.Join(result, ", ")), nil

	case yaml.MapSlice:
		result := make([]string, 0)
		for _, i := range v {
			value, err := p.ToCompactJSON(i)
			if err != nil {
				return "", err
			}
			result = append(result, value)
		}

		return fmt.Sprintf("{%s}", strings.Join(result, ", ")), nil

	case yaml.MapItem:
		key, keyError := p.ToCompactJSON(v.Key)
		if keyError != nil {
			return "", keyError
		}

		value, valueError := p.ToCompactJSON(v.Value)
		if valueError != nil {
			return "", valueError
		}

		return fmt.Sprintf("%s: %s", key, value), nil

	default:
		bytes, err := json.Marshal(v)
		if err != nil {
			return "", err
		}

		return string(bytes), nil
	}
}

func (p *OutputProcessor) neatJSON(prefix string, obj interface{}) (string, error) {
	switch obj.(type) {
	case yaml.MapSlice:
		if err := p.neatJSONofYAMLMapSlice(prefix, obj.(yaml.MapSlice)); err != nil {
			return "", err
		}

	case []interface{}:
		if err := p.neatJSONofSlice(prefix, obj.([]interface{})); err != nil {
			return "", err
		}

	case []yaml.MapSlice:
		if err := p.neatJSONofSlice(prefix, p.simplify(obj.([]yaml.MapSlice))); err != nil {
			return "", err
		}

	default:
		if err := p.neatJSONofScalar(prefix, obj); err != nil {
			return "", nil
		}
	}

	p.out.Flush()
	return p.data.String(), nil
}

func (p *OutputProcessor) neatJSONofYAMLMapSlice(prefix string, mapslice yaml.MapSlice) error {
	if len(mapslice) == 0 {
		p.out.WriteString(p.colorize("{}", "emptyStructures"))
		return nil
	}

	p.out.WriteString(bunt.BoldText("{"))
	p.out.WriteString("\n")

	for idx, mapitem := range mapslice {
		keyString := fmt.Sprintf("\"%v\": ", mapitem.Key)

		p.out.WriteString(prefix + p.prefixAdd())
		p.out.WriteString(p.colorize(keyString, "keyColor"))

		if p.isScalar(mapitem.Value) {
			p.neatJSONofScalar("", mapitem.Value)

		} else {
			p.neatJSON(prefix+p.prefixAdd(), mapitem.Value)
		}

		if idx < len(mapslice)-1 {
			p.out.WriteString(",")
		}

		p.out.WriteString("\n")
	}

	p.out.WriteString(prefix)
	p.out.WriteString(bunt.BoldText("}"))

	return nil
}

func (p *OutputProcessor) neatJSONofSlice(prefix string, list []interface{}) error {
	if len(list) == 0 {
		p.out.WriteString(p.colorize("[]", "emptyStructures"))
		return nil
	}

	p.out.WriteString(bunt.BoldText("["))
	p.out.WriteString("\n")

	for idx, value := range list {
		if p.isScalar(value) {
			p.neatJSONofScalar(prefix+p.prefixAdd(), value)

		} else {
			p.out.WriteString(prefix + p.prefixAdd())
			p.neatJSON(prefix+p.prefixAdd(), value)
		}

		if idx < len(list)-1 {
			p.out.WriteString(",")
		}

		p.out.WriteString("\n")
	}

	p.out.WriteString(prefix)
	p.out.WriteString(bunt.BoldText("]"))

	return nil
}

func (p *OutputProcessor) neatJSONofScalar(prefix string, obj interface{}) error {
	if obj == nil {
		p.out.WriteString(p.colorize("null", "nullColor"))
		return nil
	}

	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	color := p.determineColorByType(obj)

	p.out.WriteString(prefix)
	parts := strings.Split(string(data), "\\n")
	for idx, part := range parts {
		p.out.WriteString(p.colorize(part, color))

		if idx < len(parts)-1 {
			p.out.WriteString(p.colorize("\\n", "emptyStructures"))
		}
	}

	return nil
}
