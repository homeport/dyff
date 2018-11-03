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
	yaml "gopkg.in/yaml.v2"
)

// GetIdentifierFromNamedList returns the identifier key used in the provided list, or an empty string if there is none. The identifier key is either 'name', 'key', or 'id'.
func GetIdentifierFromNamedList(list []interface{}) string {
	counters := map[interface{}]int{}

	for _, sliceEntry := range list {
		switch sliceEntry.(type) {
		case yaml.MapSlice:
			for _, mapSliceEntry := range sliceEntry.(yaml.MapSlice) {
				if _, ok := counters[mapSliceEntry.Key]; !ok {
					counters[mapSliceEntry.Key] = 0
				}

				counters[mapSliceEntry.Key]++
			}
		}
	}

	sliceLength := len(list)
	for _, identifier := range []string{"name", "key", "id"} {
		if count, ok := counters[identifier]; ok && count == sliceLength {
			return identifier
		}
	}

	return ""
}

// getEntryFromNamedList returns the entry that is identified by the identifier key and a name, for example: `name: one` where name is the identifier key and one the name. Function will return nil with bool false if there is no such entry.
func getEntryFromNamedList(list []interface{}, identifier string, name interface{}) (interface{}, bool) {
	for _, listEntry := range list {
		mapslice := listEntry.(yaml.MapSlice)

		for _, element := range mapslice {
			if element.Key == identifier && element.Value == name {
				return mapslice, true
			}
		}
	}

	return nil, false
}

func listNamesOfNamedList(list []interface{}, identifier string) ([]string, error) {
	result := make([]string, len(list))
	for i, entry := range list {
		switch entry.(type) {
		case yaml.MapSlice:
			value, err := getValueByKey(entry.(yaml.MapSlice), identifier)
			if err != nil {
				return nil, err
			}

			result[i] = value.(string)

		default:
			return nil, &NoNamedEntryListError{}
		}
	}

	return result, nil
}

func splitEntryIntoNameAndData(mapslice yaml.MapSlice, identifier string) (name interface{}, result yaml.MapSlice) {
	for _, mapitem := range mapslice {
		if key, ok := mapitem.Key.(string); ok && key == identifier {
			name = mapitem.Value

		} else {
			result = append(result, mapitem)
		}
	}

	return name, result
}
