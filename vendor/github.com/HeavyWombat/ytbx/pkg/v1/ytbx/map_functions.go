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
	"fmt"

	yaml "gopkg.in/yaml.v2"
)

// listKeys returns a list of the keys of the YAML MapSlice (map).
func listKeys(mapslice yaml.MapSlice) []string {
	keys := make([]string, len(mapslice))
	for i, mapitem := range mapslice {
		keys[i] = fmt.Sprintf("%v", mapitem.Key)
	}

	return keys
}

// getValueByKey returns the value for a given key in a provided MapSlice, or nil with an error if there is no such entry. This is comparable to getting a value from a map with `foobar[key]`.
func getValueByKey(mapslice yaml.MapSlice, key string) (interface{}, error) {
	for _, element := range mapslice {
		if element.Key == key {
			return element.Value, nil
		}
	}

	return nil, &KeyNotFoundInMapError{MissingKey: key, AvailableKeys: listKeys(mapslice)}
}

func getEntryByIdentifierAndName(list []yaml.MapSlice, identifier string, name interface{}) (yaml.MapSlice, error) {
	for _, mapslice := range list {
		for _, element := range mapslice {
			if element.Key == identifier && element.Value == name {
				return mapslice, nil
			}
		}
	}

	return nil, fmt.Errorf("there is no entry %s=%v in the list", identifier, name)
}
