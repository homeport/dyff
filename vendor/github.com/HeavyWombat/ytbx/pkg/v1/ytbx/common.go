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

package ytbx

import (
	"reflect"

	yaml "gopkg.in/yaml.v2"
)

// Internal string constants for type names and type decisions
const (
	typeMap         = "map"
	typeSimpleList  = "list"
	typeComplexList = "complex-list"
	typeString      = "string"
)

// GetType returns the type of the input value with a YAML specific view
func GetType(value interface{}) string {
	switch value.(type) {
	case yaml.MapSlice:
		return typeMap

	case []interface{}:
		if IsComplexSlice(value.([]interface{})) {
			return typeComplexList
		}

		return typeSimpleList

	case []yaml.MapSlice:
		return typeComplexList

	case string:
		return typeString

	default:
		return reflect.TypeOf(value).Kind().String()
	}
}

// IsComplexSlice returns whether the slice contains (hash)map entries, otherwise the slice is called a simple list.
func IsComplexSlice(slice []interface{}) bool {
	// This is kind of a weird case, but by definition an empty list is a simple slice
	if len(slice) == 0 {
		return false
	}

	// Count the number of entries which are maps or YAML MapSlices
	counter := 0
	for _, entry := range slice {
		switch entry.(type) {
		case map[string]interface{}, map[interface{}]interface{}, yaml.MapSlice:
			counter++
		}
	}

	return counter == len(slice)
}

// SimplifyList will cast a slice of YAML MapSlices into a slice of interfaces.
func SimplifyList(input []yaml.MapSlice) []interface{} {
	result := make([]interface{}, len(input))
	for i := range input {
		result[i] = input[i]
	}

	return result
}

func isList(obj interface{}) bool {
	switch obj.(type) {
	case []interface{}:
		return true

	default:
		return false
	}
}

func isMapSlice(obj interface{}) bool {
	switch obj.(type) {
	case yaml.MapSlice:
		return true

	default:
		return false
	}
}
