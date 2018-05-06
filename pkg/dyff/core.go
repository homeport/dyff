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
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/HeavyWombat/dyff/pkg/bunt"
	"github.com/mitchellh/hashstructure"
	"github.com/pkg/errors"
	"github.com/texttheater/golang-levenshtein/levenshtein"
	"golang.org/x/crypto/ssh/terminal"
	yaml "gopkg.in/yaml.v2"
)

// DebugMode is the global switch to enable debug output
var DebugMode = false

// NoColor is the gobal switch to decide whether strings should be colored in the output
var NoColor = false

// FixedTerminalWidth disables terminal width detection and reset it with a fixed given value
var FixedTerminalWidth = -1

// Debug log output
var Debug = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)

// Info log output
var Info = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

// Warning log output
var Warning = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)

// Error log output
var Error = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

// Constants to distinguish between the different kinds of differences
const (
	ADDITION     = '+'
	REMOVAL      = '-'
	MODIFICATION = '±'
	ORDERCHANGE  = '⇆'
	ILLEGAL      = '✕'
	ATTENTION    = '⚠'
)

// PathElement describes a part of a path, meaning its name. In this case the "Key" string is empty. Named list entries such as "name: one" use both "Key" and "Name" to properly specify the path element.
type PathElement struct {
	Key  string
	Name string
}

// Path describes a position inside a YAML (or JSON) structure by providing a name to each hierarchy level (tree structure).
type Path struct {
	DocumentIdx  int
	PathElements []PathElement
}

// Detail encapsulate the actual details of a change, mainly the kind of difference and the values.
type Detail struct {
	Kind rune
	From interface{}
	To   interface{}
}

// Diff encapsulates everything noteworthy about a difference
type Diff struct {
	Path    Path
	Details []Detail
}

// Report encapsulates the actual end-result of the comparison: The input data and the list of differences.
type Report struct {
	From  InputFile
	To    InputFile
	Diffs []Diff
}

func getTerminalWidth() int {
	if FixedTerminalWidth > 0 {
		return FixedTerminalWidth
	}

	termWidth, _, err := terminal.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 80
	}

	return termWidth
}

// bold returns the provided string in 'bold' format
func bold(text string) string {
	return bunt.Style(text, bunt.Bold)
}

// italic returns the provided string in 'italic' format
func italic(text string) string {
	return bunt.Style(text, bunt.Italic)
}

func green(text string) string {
	return bunt.Colorize(text, bunt.AdditionGreen)
}

func red(text string) string {
	return bunt.Colorize(text, bunt.RemovalRed)
}

func yellow(text string) string {
	return bunt.Colorize(text, bunt.ModificationYellow)
}

// Plural returns a string with the number and noun in either singular or plural form.
// If one text argument is given, the plural will be done with the plural s. If two
// arguments are provided, the second text is the irregular plural. If more than two
// are provided, then the additional ones are simply ignored.
func Plural(amount int, text ...string) string {
	words := [...]string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten", "eleven", "twelve"}

	var number string
	if amount < len(words) {
		number = words[amount]
	} else {
		number = strconv.Itoa(amount)
	}

	switch len(text) {
	case 1:
		if amount == 1 {
			return fmt.Sprintf("%s %s", number, text[0])
		}

		return fmt.Sprintf("%s %ss", number, text[0])

	default:
		if amount == 1 {
			return fmt.Sprintf("%s %s", number, text[0])
		}

		return fmt.Sprintf("%s %s", number, text[1])
	}
}

// ToDotStyle returns a path as a string in dot style separating each path element by a dot.
// Please note that path elements that are named "." will look ugly.
func ToDotStyle(path Path, showDocumentIdx bool) string {
	pathLength := len(path.PathElements)

	// The Dot style does not really support the root level. An empty path
	// will just return a text indicating the root level is meant
	if pathLength == 0 {
		return bunt.Style("(root level)", bunt.Italic, bunt.Bold)
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
func ToGoPatchStyle(path Path, showDocumentIdx bool) string {
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

func (path Path) String() string {
	return ToGoPatchStyle(path, true)
}

// CompareInputFiles is one of the convenience main entry points for comparing objects. In this case the representation of an input file, which might contain multiple documents. It returns a report with the list of differences. Each difference describes a change to comes from "from" to "to", hence the names.
func CompareInputFiles(from InputFile, to InputFile) (Report, error) {
	if len(from.Documents) != len(to.Documents) {
		return Report{}, fmt.Errorf("Comparing YAMLs with a different number of documents is currently not supported")
	}

	result := make([]Diff, 0)
	for idx := range from.Documents {
		diffs, err := compareObjects(Path{DocumentIdx: idx}, from.Documents[idx], to.Documents[idx])
		if err != nil {
			return Report{}, err
		}

		result = append(result, diffs...)
	}

	return Report{from, to, result}, nil
}

func compareObjects(path Path, from interface{}, to interface{}) ([]Diff, error) {
	// Save some time and process some simple nil and type-change use cases immediately
	if from == nil && to != nil {
		return []Diff{{path, []Detail{{Kind: ADDITION, From: from, To: to}}}}, nil

	} else if from != nil && to == nil {
		return []Diff{{path, []Detail{{Kind: REMOVAL, From: from, To: to}}}}, nil

	} else if from == nil && to == nil {
		return []Diff{}, nil

	} else if reflect.TypeOf(from) != reflect.TypeOf(to) {
		return []Diff{{path, []Detail{{Kind: MODIFICATION, From: from, To: to}}}}, nil
	}

	var diffs []Diff
	var err error

	switch from.(type) {
	case yaml.MapSlice:
		diffs, err = compareMapSlices(path, from.(yaml.MapSlice), to.(yaml.MapSlice))

	case []interface{}:
		diffs, err = compareLists(path, from.([]interface{}), to.([]interface{}))

	case []yaml.MapSlice:
		diffs, err = compareListOfMapSlices(path, from.([]yaml.MapSlice), to.([]yaml.MapSlice))

	case string:
		diffs, err = compareStrings(path, from.(string), to.(string))

	case bool, float32, float64, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr:
		if from != to {
			diffs = []Diff{{path, []Detail{{Kind: MODIFICATION, From: from, To: to}}}}
			err = nil
		}

	default:
		err = fmt.Errorf("Failed to compare objects due to unsupported type %s", reflect.TypeOf(from))
	}

	return diffs, err
}

func compareMapSlices(path Path, from yaml.MapSlice, to yaml.MapSlice) ([]Diff, error) {
	removals := yaml.MapSlice{}
	additions := yaml.MapSlice{}

	result := make([]Diff, 0)

	for _, fromItem := range from {
		key := fromItem.Key
		if toItem, ok := getMapItemByKeyFromMapSlice(key, to); ok {
			// `from` and `to` contain the same `key` -> require comparison
			diffs, err := compareObjects(newPath(path, "", key), fromItem.Value, toItem.Value)
			if err != nil {
				return nil, err
			}
			result = append(result, diffs...)

		} else {
			// `from` contain the `key`, but `to` does not -> removal
			removals = append(removals, fromItem)
		}
	}

	for _, toItem := range to {
		key := toItem.Key
		if _, ok := getMapItemByKeyFromMapSlice(key, from); !ok {
			// `to` contains a `key` that `from` does not have -> addition
			additions = append(additions, toItem)
		}
	}

	diff := Diff{Path: path, Details: []Detail{}}

	if len(removals) > 0 {
		diff.Details = append(diff.Details, Detail{Kind: REMOVAL, From: removals, To: nil})
	}

	if len(additions) > 0 {
		diff.Details = append(diff.Details, Detail{Kind: ADDITION, From: nil, To: additions})
	}

	if len(diff.Details) > 0 {
		result = append([]Diff{diff}, result...)
	}

	return result, nil
}

func compareLists(path Path, from []interface{}, to []interface{}) ([]Diff, error) {
	if fromIdentifier := GetIdentifierFromNamedList(from); fromIdentifier != "" {
		if toIdentifier := GetIdentifierFromNamedList(to); fromIdentifier == toIdentifier {
			return compareNamedEntryLists(path, fromIdentifier, from, to)
		}
	}

	return compareSimpleLists(path, from, to)
}

func compareListOfMapSlices(path Path, from []yaml.MapSlice, to []yaml.MapSlice) ([]Diff, error) {
	// TODO Check if there is another way to do this, or if we can save time by doing something else
	return compareLists(path, SimplifyList(from), SimplifyList(to))
}

func compareSimpleLists(path Path, from []interface{}, to []interface{}) ([]Diff, error) {
	removals := make([]interface{}, 0)
	additions := make([]interface{}, 0)

	result := make([]Diff, 0)

	fromLength := len(from)
	toLength := len(to)

	// Back out immediately if both lists are empty
	if fromLength == 0 && fromLength == toLength {
		return result, nil
	}

	// Special case if both lists only contain one entry: directly compare the two entries with each other
	if fromLength == 1 && fromLength == toLength {
		return compareObjects(newPath(path, "", 0), from[0], to[0])
	}

	fromLookup, err := createLookUpMap(from)
	if err != nil {
		return nil, err
	}

	toLookup, err := createLookUpMap(to)
	if err != nil {
		return nil, err
	}

	// Fill two lists with the names of the entries that are common to both provided lists
	fromNames := make([]uint64, 0, fromLength)
	toNames := make([]uint64, 0, fromLength)

	for idxPos, fromValue := range from {
		hash, err := calcHash(fromValue)
		if err != nil {
			return nil, err
		}

		if _, ok := toLookup[hash]; !ok {
			// `from` entry does not exist in `to` list
			removals = append(removals, from[idxPos])

		} else {
			fromNames = append(fromNames, hash)
		}
	}

	for idxPos, toValue := range to {
		hash, err := calcHash(toValue)
		if err != nil {
			return nil, err
		}

		if _, ok := fromLookup[hash]; !ok {
			// `to` entry does not exist in `from` list
			additions = append(additions, to[idxPos])

		} else {
			toNames = append(toNames, hash)
		}
	}

	orderchanges := findOrderChangesInSimpleList(from, to, fromNames, toNames, fromLookup, toLookup)

	return packChangesAndAddToResult(result, true, path, orderchanges, additions, removals)
}

func compareNamedEntryLists(path Path, identifier string, from []interface{}, to []interface{}) ([]Diff, error) {
	removals := make([]interface{}, 0)
	additions := make([]interface{}, 0)

	fromLength := len(from)
	toLength := len(to)

	result := make([]Diff, 0)

	// Bail out quickly if there is nothing to check
	if fromLength == 0 && toLength == 0 {
		return result, nil
	}

	// Fill two lists with the names of the entries that are common to both provided lists
	fromNames := make([]string, 0, fromLength)
	toNames := make([]string, 0, fromLength)

	// Find entries that are common to both lists to compare them separately, and find entries that are only in from, but not to and are therefore removed
	for _, fromEntry := range from {
		name, err := getValueByKey(fromEntry.(yaml.MapSlice), identifier)
		if err != nil {
			return nil, err
		}

		if toEntry, ok := getEntryFromNamedList(to, identifier, name); ok {
			// `from` and `to` have the same entry idenfified by identifier and name -> require comparison
			diffs, err := compareObjects(newPath(path, identifier, name), fromEntry, toEntry)
			if err != nil {
				return nil, err
			}
			result = append(result, diffs...)
			fromNames = append(fromNames, name.(string))

		} else {
			// `from` has an entry (identified by identifier and name), but `to` does not -> removal
			removals = append(removals, fromEntry)
		}
	}

	// Find entries that are only in to, but not from and are therefore added
	for _, toEntry := range to {
		name, err := getValueByKey(toEntry.(yaml.MapSlice), identifier)
		if err != nil {
			return nil, err
		}

		if _, ok := getEntryFromNamedList(from, identifier, name); ok {
			// `to` and `from` have the same entry idenfified by identifier and name (comparison already covered by previous range)
			toNames = append(toNames, name.(string))

		} else {
			// `to` has an entry (identified by identifier and name), but `from` does not -> addition
			additions = append(additions, toEntry)
		}
	}

	orderchanges := findOrderChangesInNamedEntryLists(fromNames, toNames)

	return packChangesAndAddToResult(result, true, path, orderchanges, additions, removals)
}

func compareStrings(path Path, from string, to string) ([]Diff, error) {
	result := make([]Diff, 0)
	if strings.Compare(from, to) != 0 {
		result = append(result, Diff{path, []Detail{{Kind: MODIFICATION, From: from, To: to}}})
	}

	return result, nil
}

func findOrderChangesInSimpleList(from, to []interface{}, fromNames, toNames []uint64, fromLookup, toLookup map[uint64]int) []Detail {
	orderchanges := make([]Detail, 0)

	// Try to find order changes ...
	if len(fromNames) == len(toNames) {
		for idx, hash := range fromNames {
			if toNames[idx] != hash {
				cnv := func(list []uint64, lookup map[uint64]int, content []interface{}) []interface{} {
					result := make([]interface{}, 0, len(list))
					for _, hash := range list {
						result = append(result, content[lookup[hash]])
					}

					return result
				}

				orderchanges = append(orderchanges,
					Detail{
						Kind: ORDERCHANGE,
						From: cnv(fromNames, fromLookup, from),
						To:   cnv(toNames, toLookup, to),
					})
				break
			}
		}
	}

	return orderchanges
}

func findOrderChangesInNamedEntryLists(fromNames, toNames []string) []Detail {
	orderchanges := make([]Detail, 0)

	// Try to find order changes ...
	idxLookupMap := make(map[string]int, len(toNames))
	for idx, name := range toNames {
		idxLookupMap[name] = idx
	}

	for idx, name := range fromNames {
		if idxLookupMap[name] != idx {
			orderchanges = append(orderchanges, Detail{Kind: ORDERCHANGE, From: fromNames, To: toNames})
			break
		}
	}

	return orderchanges
}

func packChangesAndAddToResult(list []Diff, prepend bool, path Path, orderchanges []Detail, additions, removals []interface{}) ([]Diff, error) {
	// Prepare a diff for this path to added to the result set (if there are changes)
	diff := Diff{Path: path, Details: []Detail{}}

	if len(orderchanges) > 0 {
		diff.Details = append(diff.Details, orderchanges...)
	}

	if len(removals) > 0 {
		diff.Details = append(diff.Details, Detail{Kind: REMOVAL, From: removals, To: nil})
	}

	if len(additions) > 0 {
		diff.Details = append(diff.Details, Detail{Kind: ADDITION, From: nil, To: additions})
	}

	// If there were changes added to the details list,
	// we can safely add it to the result set.
	// Otherwise it the result set will be returned as-is.
	if len(diff.Details) > 0 {
		switch prepend {
		case true:
			list = append([]Diff{diff}, list...)

		case false:
			list = append(list, diff)
		}
	}

	return list, nil
}

func newPath(path Path, key interface{}, name interface{}) Path {
	result := make([]PathElement, len(path.PathElements))
	copy(result, path.PathElements)

	result = append(result, PathElement{
		Key:  fmt.Sprintf("%v", key),
		Name: fmt.Sprintf("%v", name)})

	return Path{
		DocumentIdx:  path.DocumentIdx,
		PathElements: result}
}

// getMapItemByKeyFromMapSlice returns the MapItem (tuple of key/value) where the MapItem key matches the provided key. It will return an empty MapItem and bool false if the given MapSlice does not contain a suitable MapItem.
func getMapItemByKeyFromMapSlice(key interface{}, mapslice yaml.MapSlice) (yaml.MapItem, bool) {
	for _, mapitem := range mapslice {
		if mapitem.Key == key {
			return mapitem, true
		}
	}

	return yaml.MapItem{}, false
}

// getValueByKey returns the value for a given key in a provided MapSlice, or nil with an error if there is no such entry. This is comparable to getting a value from a map with `foobar[key]`.
func getValueByKey(mapslice yaml.MapSlice, key string) (interface{}, error) {
	for _, element := range mapslice {
		if element.Key == key {
			return element.Value, nil
		}
	}

	if names, err := ListStringKeys(mapslice); err == nil {
		return nil, fmt.Errorf("no key '%s' found in map, available keys are: %s", key, strings.Join(names, ", "))
	}

	return nil, fmt.Errorf("no key '%s' found in map and also failed to get a list of key for this map", key)
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

// GetIdentifierFromNamedList returns the identifier key used in the provided list, or an empty string if there is none. The identifier key is either 'name', 'key', or 'id'.
func GetIdentifierFromNamedList(list []interface{}) string {
	// TODO Write additional logic to detect an identifier that is not a known one but something completely different
	// TODO Check whether there is a way to support Concourse YAMLs which do not come with one unique identifier per list

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

func listNamesOfNamedList(list []interface{}, identifier string) ([]string, error) {
	result := make([]string, len(list))
	for i, entry := range list {
		switch entry.(type) {
		case yaml.MapSlice:
			value, err := getValueByKey(entry.(yaml.MapSlice), identifier)
			if err != nil {
				return nil, errors.Wrap(err, "unable to list names of a names list")
			}

			result[i] = value.(string)

		default:
			return nil, fmt.Errorf("unable to list names of a names list, because list entry #%d is not a YAML map but %s", i, typeToName(entry))
		}
	}

	return result, nil
}

func createLookUpMap(list []interface{}) (map[uint64]int, error) {
	result := make(map[uint64]int, len(list))
	for idx, entry := range list {
		hash, err := calcHash(entry)
		if err != nil {
			return nil, err
		}
		result[hash] = idx
	}

	return result, nil
}

func calcHash(obj interface{}) (uint64, error) {
	var hash uint64
	var err error

	// Convert YAML MapSlices to maps first so that the order of keys does not matter for the hash value of this object
	switch obj.(type) {
	case yaml.MapSlice:
		tmp := make(map[interface{}]interface{}, len(obj.(yaml.MapSlice)))
		for _, entry := range obj.(yaml.MapSlice) {
			tmp[entry.Key] = entry.Value
		}
		obj = tmp
	}

	if hash, err = hashstructure.Hash(obj, nil); err != nil {
		return 0, errors.Wrap(err, "Failed to calculate hash")
	}

	return hash, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func isMinorChange(from string, to string) bool {
	levenshteinDistance := levenshtein.DistanceForStrings([]rune(from), []rune(to), levenshtein.DefaultOptions)
	referenceLength := min(len(from), len(to))

	distanceVsLengthFactor := float64(levenshteinDistance) / float64(referenceLength)
	threshold := 0.1

	return distanceVsLengthFactor < threshold
}

func isMultiLine(from string, to string) bool {
	return strings.Contains(from, "\n") || strings.Contains(to, "\n")
}

// SimplifyList will cast a slice of YAML MapSlices into a slice of interfaces.
func SimplifyList(input []yaml.MapSlice) []interface{} {
	result := make([]interface{}, len(input))
	for i := range input {
		result[i] = input[i]
	}

	return result
}

func goPatchStringToPath(path string, obj interface{}) (Path, error) {
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

func spruceStringToPath(path string, obj interface{}) (Path, error) {
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

// StringToPath creates a new Path using the provided serialized path string. In case of Spruce paths, we need the actual tree as a reference to create the correct path.
func StringToPath(path string, obj interface{}) (Path, error) {
	if strings.HasPrefix(path, "/") { // Go-path path in case it starts with a slash
		return goPatchStringToPath(path, obj)
	}

	// In any other case, try to parse as Spruce path
	return spruceStringToPath(path, obj)
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

// Grab get the value from the provided YAML tree using a path to traverse through the tree structure
func Grab(obj interface{}, pathString string) (interface{}, error) {
	path, err := StringToPath(pathString, obj)
	if err != nil {
		return nil, err
	}

	pointer := obj
	pointerPath := Path{DocumentIdx: path.DocumentIdx}

	for _, element := range path.PathElements {
		if element.Key != "" { // List
			if !isList(pointer) {
				return nil, fmt.Errorf("failed to traverse tree, expected a list but found type %s at %s", typeToName(pointer), ToGoPatchStyle(pointerPath, false))
			}

			entry, ok := getEntryFromNamedList(pointer.([]interface{}), element.Key, element.Name)
			if !ok {
				return nil, fmt.Errorf("there is no entry %s: %s in the list", element.Key, element.Name)
			}

			pointer = entry

		} else if id, err := strconv.Atoi(element.Name); err == nil { // List (entry referenced by its index)
			if !isList(pointer) {
				return nil, fmt.Errorf("failed to traverse tree, expected a list but found type %s at %s", typeToName(pointer), ToGoPatchStyle(pointerPath, false))
			}

			list := pointer.([]interface{})
			if id < 0 || id >= len(list) {
				return nil, fmt.Errorf("failed to traverse tree, provided list index %d is not in range: 0..%d", id, len(list)-1)
			}

			pointer = list[id]

		} else { // Map
			if !isMapSlice(pointer) {
				return nil, fmt.Errorf("failed to traverse tree, expected a YAML map but found type %s at %s", typeToName(pointer), ToGoPatchStyle(pointerPath, false))
			}

			entry, err := getValueByKey(pointer.(yaml.MapSlice), element.Name)
			if err != nil {
				return nil, err
			}

			pointer = entry
		}

		// Update the path that the current pointer has (only used in error case to point to the right position)
		pointerPath.PathElements = append(pointerPath.PathElements, element)
	}

	return pointer, nil
}

// ChangeRoot changes the root of an input file to a position inside its document based on the given path. Input files with more than one document are not supported, since they could have multiple elements with that path.
func ChangeRoot(inputFile *InputFile, path string, translateListToDocuments bool) error {
	if len(inputFile.Documents) != 1 {
		return fmt.Errorf("change root for an input file is only possible if there is only one document, but %s contains %s",
			inputFile.Location,
			Plural(len(inputFile.Documents), "document"))
	}

	// Find the object at the given path
	obj, err := Grab(inputFile.Documents[0], path)
	if err != nil {
		return err
	}

	if translateListToDocuments && isList(obj) {
		// Change root of input file main document to a new list of documents based on the the list that was found
		inputFile.Documents = obj.([]interface{})

	} else {
		// Change root of input file main document to the object that was found
		inputFile.Documents = []interface{}{obj}
	}

	// Parse path string and create nicely formatted output path
	if resolvedPath, err := StringToPath(path, obj); err == nil {
		path = PathToString(resolvedPath, false)
	}

	inputFile.Note = fmt.Sprintf("YAML root was changed to %s", path)

	return nil
}

func typeToName(obj interface{}) string {
	switch obj.(type) {
	case yaml.MapSlice:
		return "YAML map"

	case []yaml.MapSlice, []interface{}:
		return "YAML list"

	default:
		return reflect.TypeOf(obj).Kind().String()
	}
}
