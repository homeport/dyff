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

package dyff

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/HeavyWombat/dyff/pkg/bunt"
	"github.com/HeavyWombat/dyff/pkg/neat"
	colorful "github.com/lucasb-eyer/go-colorful"
	yaml "gopkg.in/yaml.v2"
)

// ToJSONString converts the provided object into a human readable JSON string.
func ToJSONString(obj interface{}) (string, error) {
	switch v := obj.(type) {

	case []interface{}:
		result := make([]string, 0)
		for _, i := range v {
			value, err := ToJSONString(i)
			if err != nil {
				return "", err
			}
			result = append(result, value)
		}

		return fmt.Sprintf("[%s]", strings.Join(result, ", ")), nil

	case yaml.MapSlice:
		result := make([]string, 0)
		for _, i := range v {
			value, err := ToJSONString(i)
			if err != nil {
				return "", err
			}
			result = append(result, value)
		}

		return fmt.Sprintf("{%s}", strings.Join(result, ", ")), nil

	case yaml.MapItem:
		key, keyError := ToJSONString(v.Key)
		if keyError != nil {
			return "", keyError
		}

		value, valueError := ToJSONString(v.Value)
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

// ToYAMLString converts the provided data into a human readable YAML string.
func ToYAMLString(content interface{}) (string, error) {
	output, err := yamlString(content)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("---\n%s\n", output), nil
}

func yamlString(input interface{}) (string, error) {
	return neat.NewOutputProcessor(false, true, nil).ToString(input)
}

func yamlStringInRedishColors(input interface{}) (string, error) {
	return neat.NewOutputProcessor(true, true, &map[string]colorful.Color{
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
	}).ToString(input)
}

func yamlStringInGreenishColors(input interface{}) (string, error) {
	return neat.NewOutputProcessor(true, true, &map[string]colorful.Color{
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
	}).ToString(input)
}
