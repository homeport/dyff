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
	"fmt"
	"strconv"
	"strings"

	"github.com/HeavyWombat/dyff/pkg/v1/bunt"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

// UseGoPatchPaths style paths instead of Spruce Dot-Style
var UseGoPatchPaths = false

// NewPath creates a new Path using the provided serialized path string. In case of Spruce paths, we need the actual tree as a reference to create the correct path.
func NewPath(path string, obj interface{}) (Path, error) {
	// Go-path path in case it starts with a slash
	if strings.HasPrefix(path, "/") {
		return parseGoPatchString(path, obj)
	}

	// In any other case, try to parse as Spruce path
	return parseSpruceString(path, obj)
}

// ToDotStyle returns a path as a string in dot style separating each path element by a dot.
// Please note that path elements that are named "." will look ugly.
func (path *Path) ToDotStyle(showDocumentIdx bool) string {
	pathLength := len(path.PathElements)

	// The Dot style does not really support the root level. An empty path
	// will just return a text indicating the root level is meant
	if pathLength == 0 {
		restultString := bunt.Style("(root level)", bunt.Italic, bunt.Bold)

		if showDocumentIdx {
			restultString += bunt.Colorize(fmt.Sprintf("  (document #%d)", path.DocumentIdx+1), bunt.Aquamarine)
		}

		return restultString
	}

	result := make([]string, 0, pathLength)
	for _, element := range path.PathElements {
		if element.Key != "" {
			result = append(result, bunt.Style(element.Name, bunt.Italic, bunt.Bold))
		} else {
			result = append(result, bunt.Style(element.Name, bunt.Bold))
		}
	}

	if showDocumentIdx {
		return strings.Join(result, ".") + bunt.Colorize(fmt.Sprintf("  (document #%d)", path.DocumentIdx+1), bunt.Aquamarine)
	}

	return strings.Join(result, ".")
}

// ToGoPatchStyle returns a path as a string in Go-Patch (https://github.com/cppforlife/go-patch) style separating each path element by a slash. Named list entries will be shown with their respecitive identifier name such as "name", "key", or "id".
func (path *Path) ToGoPatchStyle(showDocumentIdx bool) string {
	result := make([]string, 0, len(path.PathElements))
	for _, element := range path.PathElements {
		if element.Key != "" {
			result = append(result, fmt.Sprintf("%s=%s", bunt.Style(element.Key, bunt.Italic), bunt.Style(element.Name, bunt.Bold, bunt.Italic)))
		} else {
			result = append(result, bunt.Style(element.Name, bunt.Bold))
		}
	}

	if showDocumentIdx {
		return "/" + strings.Join(result, "/") + bunt.Colorize(fmt.Sprintf("  (document #%d)", path.DocumentIdx+1), bunt.Aquamarine)
	}

	return "/" + strings.Join(result, "/")
}

// ToString returns a nicely formatted version of the system default style (Go-patch)
func (path *Path) String() string {
	return path.ToGoPatchStyle(true)
}

// ToString returns a nicely formatted version of the provided path based on the user-preference for the style
func (path *Path) ToString(showDocumentIdx bool) string {
	if UseGoPatchPaths {
		return path.ToGoPatchStyle(showDocumentIdx)
	}

	return path.ToDotStyle(showDocumentIdx)
}

func parseGoPatchString(path string, obj interface{}) (Path, error) {
	elements := make([]PathElement, 0)

	for i, section := range strings.Split(path, "/") {
		if i == 0 {
			continue
		}

		keyNameSplit := strings.Split(section, "=")
		switch len(keyNameSplit) {
		case 1:
			elements = append(elements, PathElement{Name: keyNameSplit[0]})

		case 2:
			elements = append(elements, PathElement{Key: keyNameSplit[0], Name: keyNameSplit[1]})

		default:
			return Path{}, fmt.Errorf("invalid Go-patch style path, element '%s' cannot contain more than one equal sign", section)
		}
	}

	return Path{DocumentIdx: 0, PathElements: elements}, nil
}

func parseSpruceString(path string, obj interface{}) (Path, error) {
	elements := make([]PathElement, 0)

	pointer := obj
	for _, section := range strings.Split(path, ".") {
		if isMapSlice(pointer) {
			mapslice := pointer.(yaml.MapSlice)
			value, err := getValueByKey(mapslice, section)
			if err != nil {
				return Path{}, errors.Wrap(err, fmt.Sprintf("failed to parse path %s", path))
			}

			pointer = value
			elements = append(elements, PathElement{Name: section})

		} else if isList(pointer) {
			list := pointer.([]interface{})
			if id, err := strconv.Atoi(section); err == nil {
				if id < 0 || id >= len(list) {
					return Path{}, fmt.Errorf("failed to parse path %s, provided list index %d is not in range: 0..%d", path, id, len(list)-1)
				}

				pointer = list[id]
				elements = append(elements, PathElement{Name: section})

			} else {
				identifier := GetIdentifierFromNamedList(list)
				value, ok := getEntryFromNamedList(list, identifier, section)
				if !ok {
					names, err := listNamesOfNamedList(list, identifier)
					if err != nil {
						return Path{}, fmt.Errorf("failed to parse path %s, provided named list entry '%s' cannot be found in list", path, section)
					}

					return Path{}, fmt.Errorf("failed to parse path %s, provided named list entry '%s' cannot be found in list, available names are: %s", path, section, strings.Join(names, ", "))
				}

				pointer = value
				elements = append(elements, PathElement{Key: identifier, Name: section})
			}
		}
	}

	return Path{DocumentIdx: 0, PathElements: elements}, nil
}
