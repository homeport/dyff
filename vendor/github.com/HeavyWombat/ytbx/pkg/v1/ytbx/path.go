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
	"strconv"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

// PathStyle is a custom type for supported path styles
type PathStyle int

// Supported styles are the Dot-Style (used by Spruce for example) and GoPatch
// Style which is used by BOSH
const (
	DotStyle PathStyle = iota
	GoPatchStyle
)

// Path points to a section in a data struture by using names to identify the
// location.
// Example:
//   ---
//   sizing:
//     api:
//       count: 2
// For example, `sizing.api.count` points to the key `sizing` of the root
// element and in there to the key `api` and so on and so forth.
type Path struct {
	DocumentIdx  int
	PathElements []PathElement
}

// PathElement represents one part of a path, which can either address an entry
// in a map (by name), a named-entry list entry (key and name), or an entry in a
// list (by index).
type PathElement struct {
	Idx  int
	Key  string
	Name string
}

func (path Path) String() string {
	return path.ToGoPatchStyle()
}

// ToGoPatchStyle returns the path as a GoPatch style string.
func (path *Path) ToGoPatchStyle() string {
	sections := []string{""}

	for _, element := range path.PathElements {
		switch {
		case element.Name != "" && element.Key == "":
			sections = append(sections, element.Name)

		case element.Name != "" && element.Key != "":
			sections = append(sections, fmt.Sprintf("%s=%s", element.Key, element.Name))

		default:
			sections = append(sections, strconv.Itoa(element.Idx))
		}
	}

	return strings.Join(sections, "/")
}

// ToDotStyle returns the path as a Dot-Style string.
func (path *Path) ToDotStyle() string {
	sections := []string{}

	for _, element := range path.PathElements {
		switch {
		case element.Name != "":
			sections = append(sections, element.Name)

		case element.Idx >= 0:
			sections = append(sections, strconv.Itoa(element.Idx))
		}
	}

	return strings.Join(sections, ".")
}

// NewPathWithPathElement returns a new path based on a given path adding a new
// path element.
func NewPathWithPathElement(path Path, pathElement PathElement) Path {
	result := make([]PathElement, len(path.PathElements))
	copy(result, path.PathElements)

	return Path{
		DocumentIdx:  path.DocumentIdx,
		PathElements: append(result, pathElement)}
}

// NewPathWithNamedElement returns a new path based on a given path adding a new
// of type entry in map using the name.
func NewPathWithNamedElement(path Path, name interface{}) Path {
	return NewPathWithPathElement(path, PathElement{
		Idx:  -1,
		Name: fmt.Sprintf("%v", name)})
}

// NewPathWithNamedListElement returns a new path based on a given path adding a
// new of type entry in a named-entry list by using key and name.
func NewPathWithNamedListElement(path Path, identifier interface{}, name interface{}) Path {
	return NewPathWithPathElement(path, PathElement{
		Idx:  -1,
		Key:  fmt.Sprintf("%v", identifier),
		Name: fmt.Sprintf("%v", name)})
}

// NewPathWithIndexedListElement returns a new path based on a given path adding
// a new of type list entry using the index.
func NewPathWithIndexedListElement(path Path, idx int) Path {
	return NewPathWithPathElement(path, PathElement{
		Idx: idx,
	})
}

// ListPaths returns all paths in the documents using the provided choice of
// path style.
func ListPaths(location string, style PathStyle) ([]Path, error) {
	inputfile, err := LoadFile(location)
	if err != nil {
		return nil, err
	}

	paths := []Path{}
	for idx, document := range inputfile.Documents {
		root := Path{DocumentIdx: idx}

		traverseTree(root, document, func(path Path, _ interface{}) {
			paths = append(paths, path)
		})
	}

	return paths, nil
}

func traverseTree(path Path, obj interface{}, leafFunc func(path Path, value interface{})) {
	switch obj.(type) {
	case []interface{}:
		if identifier := GetIdentifierFromNamedList(obj.([]interface{})); identifier != "" {
			for _, entry := range obj.([]interface{}) {
				name, data := splitEntryIntoNameAndData(entry.(yaml.MapSlice), identifier)
				traverseTree(NewPathWithNamedListElement(path, identifier, name), data, leafFunc)
			}

		} else {
			for idx, entry := range obj.([]interface{}) {
				traverseTree(NewPathWithIndexedListElement(path, idx), entry, leafFunc)
			}
		}

	case yaml.MapSlice:
		for _, mapitem := range obj.(yaml.MapSlice) {
			traverseTree(NewPathWithNamedElement(path, mapitem.Key), mapitem.Value, leafFunc)
		}

	default:
		leafFunc(path, obj)
	}
}

// ParseGoPatchStylePathString returns a path by parsing a string representation
// which is assumed to be a GoPatch style path.
func ParseGoPatchStylePathString(path string) (Path, error) {
	// Special case for root path
	if path == "/" {
		return Path{DocumentIdx: 0, PathElements: nil}, nil
	}

	// Poor mans solution to deal with escaped slashes, replace them with a "safe"
	// replacement string that is later resolved into a simple slash
	path = strings.Replace(path, `\/`, `%2F`, -1)

	elements := make([]PathElement, 0)
	for i, section := range strings.Split(path, "/") {
		if i == 0 {
			continue
		}

		keyNameSplit := strings.Split(section, "=")
		switch len(keyNameSplit) {
		case 1:
			if idx, err := strconv.Atoi(keyNameSplit[0]); err == nil {
				elements = append(elements, PathElement{
					Idx: idx,
				})

			} else {
				elements = append(elements, PathElement{
					Idx:  -1,
					Name: strings.Replace(keyNameSplit[0], `%2F`, "/", -1),
				})
			}

		case 2:
			elements = append(elements, PathElement{Idx: -1,
				Key:  strings.Replace(keyNameSplit[0], `%2F`, "/", -1),
				Name: strings.Replace(keyNameSplit[1], `%2F`, "/", -1),
			})

		default:
			return Path{}, &InvalidPathString{
				Style:       GoPatchStyle,
				PathString:  path,
				Explanation: fmt.Sprintf("element '%s' cannot contain more than one equal sign", section),
			}
		}
	}

	return Path{DocumentIdx: 0, PathElements: elements}, nil
}

// ParseDotStylePathString returns a path by parsing a string representation
// which is assumed to be a Dot-Style path.
func ParseDotStylePathString(path string, obj interface{}) (Path, error) {
	elements := make([]PathElement, 0)

	pointer := obj
	for _, section := range strings.Split(path, ".") {
		switch {
		case isMapSlice(pointer):
			mapslice := pointer.(yaml.MapSlice)
			if value, err := getValueByKey(mapslice, section); err == nil {
				pointer = value
				elements = append(elements, PathElement{Idx: -1, Name: section})

			} else {
				pointer = nil
				elements = append(elements, PathElement{Idx: -1, Name: section})
			}

		case isList(pointer):
			list := pointer.([]interface{})
			if id, err := strconv.Atoi(section); err == nil {
				if id < 0 || id >= len(list) {
					return Path{}, &InvalidPathString{
						Style:       DotStyle,
						PathString:  path,
						Explanation: fmt.Sprintf("provided list index %d is not in range: 0..%d", id, len(list)-1),
					}
				}

				pointer = list[id]
				elements = append(elements, PathElement{Idx: id})

			} else {
				identifier := GetIdentifierFromNamedList(list)
				value, ok := getEntryFromNamedList(list, identifier, section)
				if !ok {
					names, err := listNamesOfNamedList(list, identifier)
					if err != nil {
						return Path{}, &InvalidPathString{
							Style:       DotStyle,
							PathString:  path,
							Explanation: fmt.Sprintf("provided named list entry '%s' cannot be found in list", section),
						}
					}

					return Path{}, &InvalidPathString{
						Style:       DotStyle,
						PathString:  path,
						Explanation: fmt.Sprintf("provided named list entry '%s' cannot be found in list, available names are: %s", section, strings.Join(names, ", ")),
					}
				}

				pointer = value
				elements = append(elements, PathElement{Idx: -1, Key: identifier, Name: section})
			}

		case pointer == nil:
			// If the pointer is nil, it means that the previous section of the path
			// string could not be found in the data structure and that all remaining
			// sections are assumed to be of type map.
			elements = append(elements, PathElement{Idx: -1, Name: section})
		}
	}

	return Path{DocumentIdx: 0, PathElements: elements}, nil
}

// ParsePathString returns a path by parsing a string representation
// of a path, which can be one of the supported types.
func ParsePathString(pathString string, obj interface{}) (Path, error) {
	if strings.HasPrefix(pathString, "/") {
		return ParseGoPatchStylePathString(pathString)
	}

	return ParseDotStylePathString(pathString, obj)
}
