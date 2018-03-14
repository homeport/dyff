package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"
	"unicode/utf8"

	"github.com/fatih/color"
	"github.com/texttheater/golang-levenshtein/levenshtein"
	yaml "gopkg.in/yaml.v2"
)

// Debug log output
var Debug = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)

// Info log output
var Info = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

// Warning log output
var Warning = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)

// Error log output
var Error = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

// Constants to differenciate between the different kinds of differences
const (
	ADDITION     = '+'
	REMOVAL      = '-'
	MODIFICATION = '±'
	ILLEGAL      = '✕'
	ATTENTION    = '⚠'
)

// PathElement describes a part of a path, meaning its name. In this case the "Key" string is empty. Named list entries such as "name: one" use both "Key" and "Name" to properly specify the path element.
type PathElement struct {
	Key  string
	Name string
}

// Path describes a position inside a YAML (or JSON) structure by providing a name to each hierarchy level (tree structure).
type Path []PathElement

// Diff encapsulates everything noteworthy about a difference
type Diff struct {
	Kind rune
	Path string
	From interface{}
	To   interface{}
}

// ANSI coloring convenience helpers
var bold = color.New(color.Bold)
var italic = color.New(color.Italic)

// Bold returns the provided string in 'bold' format
func Bold(text string) string {
	return bold.Sprint(text)
}

// Italic returns the provided string in 'italic' format
func Italic(text string) string {
	return italic.Sprint(text)
}

// ToDotStyle returns a path as a string in dot style separating each path element by a dot.
// Please note that path elements that are named "." will look ugly.
func ToDotStyle(path Path) string {
	result := make([]string, 0, len(path))
	for _, element := range path {
		if element.Key != "" {
			result = append(result, Italic(element.Name))
		} else {
			result = append(result, element.Name)
		}
	}

	return strings.Join(result, ".")
}

// ToGoPatchStyle returns a path as a string in Go-Patch (https://github.com/cppforlife/go-patch) style separating each path element by a slash. Named list entries will be shown with their respecitive identifier name such as "name", "key", or "id".
func ToGoPatchStyle(path Path) string {
	result := make([]string, 0, len(path))
	for _, element := range path {
		if element.Key != "" {
			result = append(result, fmt.Sprintf("%s=%s", element.Key, element.Name))
		} else {
			result = append(result, element.Name)
		}
	}

	return "/" + strings.Join(result, "/")
}

func (path Path) String() string {
	return ToGoPatchStyle(path)
}

// CompareObjects returns a list of differences between `from` and `to`
func CompareObjects(from interface{}, to interface{}) []Diff {
	result := make([]Diff, 0)

	Debug.Printf("Entering Compare(from %s, to %s)", reflect.TypeOf(from), reflect.TypeOf(to))
	switch from.(type) {

	case yaml.MapSlice:
		switch to.(type) {
		case yaml.MapSlice:
			result = append(result, compareMapSlices(from.(yaml.MapSlice), to.(yaml.MapSlice))...)

		}

	case []interface{}:
		switch to.(type) {
		case []interface{}:
			result = append(result, compareLists(from.([]interface{}), to.([]interface{}))...)
		}

	case string:
		switch to.(type) {
		case string:
			result = append(result, compareStrings(from.(string), to.(string))...)

		}

	case bool, float32, float64, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr:
		switch to.(type) {
		case bool, float32, float64, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr:
			if from != to {
				result = append(result, Diff{Kind: MODIFICATION, From: from, To: to})
			}
		}

	default:
		panic(fmt.Sprintf("Unsupported type %s", reflect.TypeOf(from)))
	}

	return result
}

func compareMapSlices(from yaml.MapSlice, to yaml.MapSlice) []Diff {
	return compareMaps(convertMapSliceToMap(from), convertMapSliceToMap(to))
}

func compareMaps(from map[interface{}]interface{}, to map[interface{}]interface{}) []Diff {
	result := make([]Diff, 0)

	for fromKey, fromValue := range from {
		if toValue, ok := to[fromKey]; ok {
			// `from` and `to` contain the same `key` -> require comparison
			result = append(result, CompareObjects(fromValue, toValue)...)

		} else {
			// `from` contain the `key`, but `to` does not -> removal
			result = append(result, Diff{Kind: REMOVAL, From: fromValue, To: nil})
		}
	}

	for toKey, toValue := range to {
		if _, ok := from[toKey]; !ok {
			// `to` contains a `key` that `from` does not have -> addition
			result = append(result, Diff{Kind: ADDITION, From: nil, To: toValue})
		}
	}

	return result
}

func compareLists(from []interface{}, to []interface{}) []Diff {
	result := make([]Diff, 0)

	fromLookup := createLookUpMap(from)
	toLookup := createLookUpMap(to)

	for fromValue := range fromLookup {
		if _, ok := toLookup[fromValue]; !ok {
			// `from` entry does not exist in `to` list
			result = append(result, Diff{Kind: REMOVAL, From: fromValue, To: nil})
		}
	}

	for toValue := range toLookup {
		if _, ok := fromLookup[toValue]; !ok {
			// `to` entry does not exist in `from` list
			result = append(result, Diff{Kind: ADDITION, From: nil, To: toValue})
		}
	}

	return result
}

func compareStrings(from string, to string) []Diff {
	distance := levenshtein.DistanceForStrings([]rune(from), []rune(to), levenshtein.DefaultOptions)
	relative := float64(distance) / float64(utf8.RuneCountInString(to))
	Debug.Printf("levenshtein distance between %s and %s is %d (relative: %f)", from, to, distance, relative)

	result := make([]Diff, 0)
	if strings.Compare(from, to) != 0 {
		result = append(result, Diff{Kind: MODIFICATION, From: from, To: to})
	}

	return result
}

func createLookUpMap(list []interface{}) map[interface{}]struct{} {
	result := make(map[interface{}]struct{}, len(list))
	for _, entry := range list {
		result[entry] = struct{}{}
	}

	return result
}

func convertMapSliceToMap(mapslice yaml.MapSlice) map[interface{}]interface{} {
	result := make(map[interface{}]interface{})
	for _, entry := range mapslice {
		result[entry.Key] = entry.Value
	}

	return result
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
