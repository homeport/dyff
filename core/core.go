package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/fatih/color"
	yaml "gopkg.in/yaml.v2"
)

// ANSI coloring convenience helpers
var bold = color.New(color.Bold)

// Bold returns the provided string in 'bold' format
func Bold(text string) string {
	return bold.Sprint(text)
}

// LoadFile Processes the provided input location to load a YAML (or JSON) into a yaml.MapSlice
func LoadFile(location string) (yaml.MapSlice, error) {
	// TODO Support URIs as loaction
	// TODO Support STDIN as location
	// TODO Generate error if file contains more than one document

	data, ioerr := ioutil.ReadFile(location)
	if ioerr != nil {
		return nil, ioerr
	}

	content := yaml.MapSlice{}
	if err := yaml.UnmarshalStrict([]byte(data), &content); err != nil {
		return nil, err
	}

	return content, nil
}

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

		return fmt.Sprintf("%s", string(bytes)), nil
	}
}

// ToYAMLString converts the provided YAML MapSlice into a human readable YAML string.
func ToYAMLString(content yaml.MapSlice) (string, error) {
	out, err := yaml.Marshal(content)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("---\n%s\n", string(out)), nil
}
