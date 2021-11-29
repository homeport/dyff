// Copyright © 2019 The Homeport Team
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
	"sort"
	"strings"

	"github.com/gonvenience/bunt"
	"github.com/gonvenience/text"
	"github.com/gonvenience/wrap"
	"github.com/gonvenience/ytbx"
	"github.com/mitchellh/hashstructure"
	yamlv3 "gopkg.in/yaml.v3"
)

// CompareOption sets a specific compare setting for the object comparison
type CompareOption func(*compareSettings)

type compareSettings struct {
	NonStandardIdentifierGuessCountThreshold int
	IgnoreOrderChanges                       bool
	KubernetesEntityDetection                bool
}

type compare struct {
	settings compareSettings
}

// ListItemIdentifierField names the field that identifies a list.
type ListItemIdentifierField string

// NonStandardIdentifierGuessCountThreshold specifies how many list entries are
// needed for the guess-the-identifier function to actually consider the key
// name. Or in short, if the lists only contain two entries each, there are more
// possibilities to find unique enough key, which might no qualify as such.
func NonStandardIdentifierGuessCountThreshold(nonStandardIdentifierGuessCountThreshold int) CompareOption {
	return func(settings *compareSettings) {
		settings.NonStandardIdentifierGuessCountThreshold = nonStandardIdentifierGuessCountThreshold
	}
}

// IgnoreOrderChanges disables the detection for changes of the order in lists
func IgnoreOrderChanges(value bool) CompareOption {
	return func(settings *compareSettings) {
		settings.IgnoreOrderChanges = value

	}
}

// KubernetesEntityDetection enabled detecting entity identifiers from Kubernetes "kind:" and "metadata:" fields.
func KubernetesEntityDetection(value bool) CompareOption {
	return func(settings *compareSettings) {
		settings.KubernetesEntityDetection = value

	}
}

// CompareInputFiles is one of the convenience main entry points for comparing
// objects. In this case the representation of an input file, which might
// contain multiple documents. It returns a report with the list of differences.
func CompareInputFiles(from ytbx.InputFile, to ytbx.InputFile, compareOptions ...CompareOption) (Report, error) {
	if len(from.Documents) != len(to.Documents) {
		return Report{}, fmt.Errorf("comparing YAMLs with a different number of documents is currently not supported")
	}

	// initialize the comparator with the tool defaults
	compare := compare{
		settings: compareSettings{
			NonStandardIdentifierGuessCountThreshold: 3,
			IgnoreOrderChanges:                       false,
			KubernetesEntityDetection:                true,
		},
	}

	// apply the optional compare options provided to this function call
	for _, compareOption := range compareOptions {
		compareOption(&compare.settings)
	}

	result := make([]Diff, 0)
	for idx := range from.Documents {
		diffs, err := compare.objects(
			ytbx.Path{
				Root:        &from,
				DocumentIdx: idx,
			},
			from.Documents[idx],
			to.Documents[idx],
		)

		if err != nil {
			return Report{}, err
		}

		result = append(result, diffs...)
	}

	return Report{from, to, result}, nil
}

func (compare *compare) objects(path ytbx.Path, from *yamlv3.Node, to *yamlv3.Node) ([]Diff, error) {
	switch {
	case from == nil && to == nil:
		return []Diff{}, nil

	case (from == nil && to != nil) || (from != nil && to == nil):
		return []Diff{{
			&path,
			[]Detail{{
				Kind: MODIFICATION,
				From: from,
				To:   to,
			}},
		}}, nil

	case (from.Kind != to.Kind) || (from.Tag != to.Tag):
		return []Diff{{
			&path,
			[]Detail{{
				Kind: MODIFICATION,
				From: from,
				To:   to,
			}},
		}}, nil
	}

	return compare.nonNilSameKindNodes(path, from, to)
}

func (compare *compare) nonNilSameKindNodes(path ytbx.Path, from *yamlv3.Node, to *yamlv3.Node) ([]Diff, error) {
	var diffs []Diff
	var err error

	switch from.Kind {
	case yamlv3.DocumentNode:
		diffs, err = compare.objects(path, from.Content[0], to.Content[0])

	case yamlv3.MappingNode:
		diffs, err = compare.mappingNodes(path, from, to)

	case yamlv3.SequenceNode:
		diffs, err = compare.sequenceNodes(path, from, to)

	case yamlv3.ScalarNode:
		switch from.Tag {
		case "!!str":
			diffs, err = compare.nodeValues(path, from, to)

		case "!!null":
			// Ignore different ways to define a null value

		default:
			if from.Value != to.Value {
				diffs, err = []Diff{{
					&path,
					[]Detail{{
						Kind: MODIFICATION,
						From: from,
						To:   to,
					}},
				}}, nil
			}
		}

	case yamlv3.AliasNode:
		diffs, err = compare.objects(path, from.Alias, to.Alias)

	default:
		err = fmt.Errorf("failed to compare objects due to unsupported kind %v", from.Kind)
	}

	return diffs, err
}

func (compare *compare) mappingNodes(path ytbx.Path, from *yamlv3.Node, to *yamlv3.Node) ([]Diff, error) {
	result := make([]Diff, 0)
	removals := []*yamlv3.Node{}
	additions := []*yamlv3.Node{}

	for i := 0; i < len(from.Content); i += 2 {
		key, fromItem := from.Content[i], from.Content[i+1]
		if toItem, ok := findValueByKey(to, key.Value); ok {
			// `from` and `to` contain the same `key` -> require comparison
			diffs, err := compare.objects(
				ytbx.NewPathWithNamedElement(path, key.Value),
				followAlias(fromItem),
				followAlias(toItem),
			)

			if err != nil {
				return nil, err
			}

			result = append(result, diffs...)

		} else {
			// `from` contain the `key`, but `to` does not -> removal
			removals = append(removals, key, fromItem)
		}
	}

	for i := 0; i < len(to.Content); i += 2 {
		key, toItem := to.Content[i], to.Content[i+1]
		if _, ok := findValueByKey(from, key.Value); !ok {
			// `to` contains a `key` that `from` does not have -> addition
			additions = append(additions, key, toItem)
		}
	}

	diff := Diff{Path: &path, Details: []Detail{}}

	if len(removals) > 0 {
		diff.Details = append(diff.Details,
			Detail{
				Kind: REMOVAL,
				From: &yamlv3.Node{
					Kind:    from.Kind,
					Tag:     from.Tag,
					Content: removals,
				},
				To: nil,
			},
		)
	}

	if len(additions) > 0 {
		diff.Details = append(diff.Details,
			Detail{
				Kind: ADDITION,
				From: nil,
				To: &yamlv3.Node{
					Kind:    to.Kind,
					Tag:     to.Tag,
					Content: additions,
				},
			},
		)
	}

	if len(diff.Details) > 0 {
		result = append([]Diff{diff}, result...)
	}

	return result, nil
}

func (compare *compare) sequenceNodes(path ytbx.Path, from *yamlv3.Node, to *yamlv3.Node) ([]Diff, error) {
	// Bail out quickly if there is nothing to check
	if len(from.Content) == 0 && len(to.Content) == 0 {
		return []Diff{}, nil
	}

	if identifier, err := compare.getIdentifierFromNamedLists(from, to); err == nil {
		return compare.namedEntryLists(path, identifier, from, to)
	}

	if identifier := getNonStandardIdentifierFromNamedLists(from, to, compare.settings.NonStandardIdentifierGuessCountThreshold); identifier != "" {
		d, err := compare.namedEntryLists(path, identifier, from, to)
		if err != nil {
			return nil, fmt.Errorf("sequenceNodes(nonstd): %w", err)
		}

		return d, nil
	}

	if compare.settings.KubernetesEntityDetection {
		if identifier, err := getIdentifierFromKubernetesEntityList(from, to); err == nil {
			return compare.namedEntryLists(path, identifier, from, to)
		}
	}

	return compare.simpleLists(path, from, to)
}

func (compare *compare) simpleLists(path ytbx.Path, from *yamlv3.Node, to *yamlv3.Node) ([]Diff, error) {
	removals := make([]*yamlv3.Node, 0)
	additions := make([]*yamlv3.Node, 0)

	fromLength := len(from.Content)
	toLength := len(to.Content)

	// Special case if both lists only contain one entry, then directly compare
	// the two entries with each other
	if fromLength == 1 && fromLength == toLength {
		return compare.objects(
			ytbx.NewPathWithIndexedListElement(path, 0),
			followAlias(from.Content[0]),
			followAlias(to.Content[0]),
		)
	}

	fromLookup := compare.createLookUpMap(from)
	toLookup := compare.createLookUpMap(to)

	// Fill two lists with the hashes of the entries of each list
	fromCommon := make([]*yamlv3.Node, 0, fromLength)
	toCommon := make([]*yamlv3.Node, 0, toLength)

	for idxPos, fromValue := range from.Content {
		hash := compare.calcNodeHash(fromValue)
		_, ok := toLookup[hash]
		if ok {
			fromCommon = append(fromCommon, fromValue)
		}

		switch {
		case !ok:
			// `from` entry does not exist in `to` list
			removals = append(removals, from.Content[idxPos])

		case len(fromLookup[hash]) > len(toLookup[hash]):
			// `from` entry exists in `to` list, but there are duplicates and
			// the number of duplicates is smaller
			if !compare.hasEntry(removals, from.Content[idxPos]) {
				for i := 0; i < len(fromLookup[hash])-len(toLookup[hash]); i++ {
					removals = append(removals, from.Content[idxPos])
				}
			}
		}
	}

	for idxPos, toValue := range to.Content {
		hash := compare.calcNodeHash(toValue)
		_, ok := fromLookup[hash]
		if ok {
			toCommon = append(toCommon, toValue)
		}

		switch {
		case !ok:
			// `to` entry does not exist in `from` list
			additions = append(additions, to.Content[idxPos])

		case len(fromLookup[hash]) < len(toLookup[hash]):
			// `to` entry exists in `from` list, but there are duplicates and
			// the number of duplicates is increased
			if !compare.hasEntry(additions, to.Content[idxPos]) {
				for i := 0; i < len(toLookup[hash])-len(fromLookup[hash]); i++ {
					additions = append(additions, to.Content[idxPos])
				}
			}
		}
	}

	var orderChanges []Detail
	if !compare.settings.IgnoreOrderChanges {
		orderChanges = compare.findOrderChangesInSimpleList(fromCommon, toCommon)
	}

	return packChangesAndAddToResult([]Diff{}, path, orderChanges, additions, removals)
}

func nameFromPath(node *yamlv3.Node, field ListItemIdentifierField) (string, error) {
	parts := strings.SplitN(string(field), ".", 2)
	key := parts[0]
	val, err := getValueByKey(node, key)
	if err != nil {
		return "", fmt.Errorf("nameFromPath issue: %w", err)
	}
	if len(parts) == 1 {
		return val.Value, nil
	}
	return nameFromPath(val, ListItemIdentifierField(parts[1]))
}

func (compare *compare) namedEntryLists(path ytbx.Path, identifier ListItemIdentifierField, from *yamlv3.Node, to *yamlv3.Node) ([]Diff, error) {
	removals := make([]*yamlv3.Node, 0)
	additions := make([]*yamlv3.Node, 0)

	result := make([]Diff, 0)

	// Fill two lists with the names of the entries that are common in both lists
	fromLength := len(from.Content)
	fromNames := make([]string, 0, fromLength)
	toNames := make([]string, 0, fromLength)

	// Find entries that are common to both lists to compare them separately, and
	// find entries that are only in from, but not to and are therefore removed
	for _, fromEntry := range from.Content {
		name, err := nameFromPath(fromEntry, identifier)
		if err != nil {
			return nil, fmt.Errorf("nameEntryList from issue: %w", err)
		}

		if toEntry, ok := getEntryFromNamedList(to, identifier, name); ok {
			// `from` and `to` have the same entry identified by identifier and name -> require comparison
			diffs, err := compare.objects(
				ytbx.NewPathWithNamedListElement(path, identifier, name),
				followAlias(fromEntry),
				followAlias(toEntry),
			)
			if err != nil {
				return nil, err
			}
			result = append(result, diffs...)
			fromNames = append(fromNames, name)

		} else {
			// `from` has an entry (identified by identifier and name), but `to` does not -> removal
			removals = append(removals, fromEntry)
		}
	}

	// Find entries that are only in to, but not from and are therefore added
	for _, toEntry := range to.Content {
		name, err := nameFromPath(toEntry, identifier)
		if err != nil {
			return nil, fmt.Errorf("nameEntryList to issue: %w", err)
		}

		if _, ok := getEntryFromNamedList(from, identifier, name); ok {
			// `to` and `from` have the same entry identified by identifier and name (comparison already covered by previous range)
			toNames = append(toNames, name)

		} else {
			// `to` has an entry (identified by identifier and name), but `from` does not -> addition
			additions = append(additions, toEntry)
		}
	}

	var orderChanges []Detail
	if !compare.settings.IgnoreOrderChanges {
		orderChanges = findOrderChangesInNamedEntryLists(fromNames, toNames)
	}

	return packChangesAndAddToResult(result, path, orderChanges, additions, removals)
}

func (compare *compare) nodeValues(path ytbx.Path, from *yamlv3.Node, to *yamlv3.Node) ([]Diff, error) {
	result := make([]Diff, 0)
	if strings.Compare(from.Value, to.Value) != 0 {
		result = append(result, Diff{
			&path,
			[]Detail{{
				Kind: MODIFICATION,
				From: from,
				To:   to,
			}},
		})
	}

	return result, nil
}

func (compare *compare) findOrderChangesInSimpleList(fromCommon, toCommon []*yamlv3.Node) []Detail {
	// Try to find order changes ...
	if len(fromCommon) == len(toCommon) {
		for idx := range fromCommon {
			if compare.calcNodeHash(fromCommon[idx]) != compare.calcNodeHash(toCommon[idx]) {
				return []Detail{{
					Kind: ORDERCHANGE,
					From: &yamlv3.Node{Kind: yamlv3.SequenceNode, Content: fromCommon},
					To:   &yamlv3.Node{Kind: yamlv3.SequenceNode, Content: toCommon},
				}}
			}
		}
	}

	return []Detail{}
}

// hasEntry returns whether the given node is in the provided list. Not exactly
// a fast or efficient way to verify that a node is already in a list, but
// given that this should rarely be used it is ok for now.
func (compare *compare) hasEntry(list []*yamlv3.Node, searchEntry *yamlv3.Node) bool {
	var searchEntryHash = compare.calcNodeHash(searchEntry)
	for _, listEntry := range list {
		if searchEntryHash == compare.calcNodeHash(listEntry) {
			return true
		}
	}

	return false
}

// AsSequenceNode translates a string list into a SequenceNode
func AsSequenceNode(list []string) *yamlv3.Node {
	result := make([]*yamlv3.Node, len(list))
	for i, entry := range list {
		result[i] = &yamlv3.Node{
			Kind:  yamlv3.ScalarNode,
			Tag:   "!!str",
			Value: entry,
		}
	}

	return &yamlv3.Node{
		Kind:    yamlv3.SequenceNode,
		Content: result,
	}
}

func findOrderChangesInNamedEntryLists(fromNames, toNames []string) []Detail {
	orderchanges := make([]Detail, 0)

	idxLookupMap := make(map[string]int, len(toNames))
	for idx, name := range toNames {
		idxLookupMap[name] = idx
	}

	// Try to find order changes ...
	for idx, name := range fromNames {
		if idxLookupMap[name] != idx {
			orderchanges = append(orderchanges, Detail{
				Kind: ORDERCHANGE,
				From: AsSequenceNode(fromNames),
				To:   AsSequenceNode(toNames),
			})
			break
		}
	}

	return orderchanges
}

func packChangesAndAddToResult(list []Diff, path ytbx.Path, orderchanges []Detail, additions, removals []*yamlv3.Node) ([]Diff, error) {
	// Prepare a diff for this path to added to the result set (if there are changes)
	diff := Diff{Path: &path, Details: []Detail{}}

	if len(orderchanges) > 0 {
		diff.Details = append(diff.Details, orderchanges...)
	}

	if len(removals) > 0 {
		diff.Details = append(diff.Details, Detail{
			Kind: REMOVAL,
			From: &yamlv3.Node{
				Kind:    yamlv3.SequenceNode,
				Tag:     "!!seq",
				Content: removals,
			},
			To: nil,
		})
	}

	if len(additions) > 0 {
		diff.Details = append(diff.Details, Detail{
			Kind: ADDITION,
			From: nil,
			To: &yamlv3.Node{
				Kind:    yamlv3.SequenceNode,
				Tag:     "!!seq",
				Content: additions,
			},
		})
	}

	// If there were changes added to the details list, we can safely add it to
	// the result set. Otherwise it the result set will be returned as-is.
	if len(diff.Details) > 0 {
		list = append([]Diff{diff}, list...)
	}

	return list, nil
}

func followAlias(node *yamlv3.Node) *yamlv3.Node {
	if node != nil && node.Alias != nil {
		return followAlias(node.Alias)
	}

	return node
}

func findValueByKey(mappingNode *yamlv3.Node, key string) (*yamlv3.Node, bool) {
	for i := 0; i < len(mappingNode.Content); i += 2 {
		k, v := followAlias(mappingNode.Content[i]), followAlias(mappingNode.Content[i+1])
		if k.Value == key {
			return v, true
		}
	}

	return nil, false
}

// getValueByKey returns the value for a given key in a provided mapping node,
// or nil with an error if there is no such entry. This is comparable to getting
// a value from a map with `foobar[key]`.
func getValueByKey(mappingNode *yamlv3.Node, key string) (*yamlv3.Node, error) {
	for i := 0; i < len(mappingNode.Content); i += 2 {
		k, v := followAlias(mappingNode.Content[i]), followAlias(mappingNode.Content[i+1])
		if k.Value == key {
			return v, nil
		}
	}

	if names, err := ytbx.ListStringKeys(mappingNode); err == nil {
		return nil, fmt.Errorf("no key '%s' found in map, available keys are: %s", key, strings.Join(names, ", "))
	}

	return nil, fmt.Errorf("no key '%s' found in map and also failed to get a list of key for this map", key)
}

// getEntryFromNamedList returns the entry that is identified by the identifier
// key and a name, for example: `name: one` where name is the identifier key and
// one the name. Function will return nil with bool false if there is no entry.
func getEntryFromNamedList(sequenceNode *yamlv3.Node, identifier ListItemIdentifierField, name string) (*yamlv3.Node, bool) {
	for _, mappingNode := range sequenceNode.Content {
		nodeName, _ := nameFromPath(mappingNode, identifier)
		if nodeName == name {
			return mappingNode, true
		}
	}
	return nil, false
}

func (compare *compare) listItemIdentifierCandidates() []ListItemIdentifierField {
	var candidates = []ListItemIdentifierField{"name", "key", "id"}

	// Add Kubernetes specific extra candidate
	if compare.settings.KubernetesEntityDetection {
		candidates = append(candidates, "manager")
	}

	return candidates
}

func (compare *compare) getIdentifierFromNamedLists(listA, listB *yamlv3.Node) (ListItemIdentifierField, error) {
	isCandidate := func(node *yamlv3.Node) bool {
		if node.Kind == yamlv3.ScalarNode {
			for _, entry := range compare.listItemIdentifierCandidates() {
				if node.Value == string(entry) {
					return true
				}
			}
		}

		return false
	}

	createKeyCountMap := func(sequenceNode *yamlv3.Node) map[ListItemIdentifierField]map[string]struct{} {
		result := map[ListItemIdentifierField]map[string]struct{}{}
		for _, entry := range sequenceNode.Content {
			switch entry.Kind {
			case yamlv3.MappingNode:
				for i := 0; i < len(entry.Content); i += 2 {
					k, v := followAlias(entry.Content[i]), followAlias(entry.Content[i+1])
					if isCandidate(k) {
						if _, found := result[ListItemIdentifierField(k.Value)]; !found {
							result[ListItemIdentifierField(k.Value)] = map[string]struct{}{}
						}

						result[ListItemIdentifierField(k.Value)][v.Value] = struct{}{}
					}
				}
			}
		}

		return result
	}

	counterA := createKeyCountMap(listA)
	counterB := createKeyCountMap(listB)

	// Check for the usual suspects: name, key, and id
	for _, identifier := range compare.listItemIdentifierCandidates() {
		if countA, okA := counterA[identifier]; okA && len(countA) == len(listA.Content) {
			if countB, okB := counterB[identifier]; okB && len(countB) == len(listB.Content) {
				return identifier, nil
			}
		}
	}

	return "", fmt.Errorf("unable to find a key that can serve as an unique identifier")
}

// getIdentifierFromKubernetesEntityList returns 'metadata.name' as a field identifier if the provided objects all have the key.
func getIdentifierFromKubernetesEntityList(listA, listB *yamlv3.Node) (ListItemIdentifierField, error) {
	key := ListItemIdentifierField("metadata.name")
	allHaveMetadataName := func(sequenceNode *yamlv3.Node) bool {
		numWithMetadata := 0
		for _, entry := range sequenceNode.Content {
			switch entry.Kind {
			case yamlv3.MappingNode:
				_, err := nameFromPath(entry, key)
				if err == nil {
					numWithMetadata++
				}
			}
		}
		return numWithMetadata == len(sequenceNode.Content)
	}

	listAHasKey := allHaveMetadataName(listA)
	listBHasKey := allHaveMetadataName(listB)
	if listAHasKey && listBHasKey {
		return key, nil
	}

	return "", fmt.Errorf("not all entities appear to have %q fields", key)
}

func getNonStandardIdentifierFromNamedLists(listA, listB *yamlv3.Node, nonStandardIdentifierGuessCountThreshold int) ListItemIdentifierField {
	createKeyCountMap := func(list *yamlv3.Node) map[string]int {
		tmp := map[string]map[string]struct{}{}
		for _, entry := range list.Content {
			if entry.Kind != yamlv3.MappingNode {
				return map[string]int{}
			}

			for i := 0; i < len(entry.Content); i += 2 {
				k, v := followAlias(entry.Content[i]), followAlias(entry.Content[i+1])
				if k.Kind == yamlv3.ScalarNode && k.Tag == "!!str" &&
					v.Kind == yamlv3.ScalarNode && v.Tag == "!!str" {
					if _, ok := tmp[k.Value]; !ok {
						tmp[k.Value] = map[string]struct{}{}
					}

					tmp[k.Value][v.Value] = struct{}{}
				}
			}
		}

		result := map[string]int{}
		for key, value := range tmp {
			result[key] = len(value)
		}

		return result
	}

	listALength := len(listA.Content)
	listBLength := len(listB.Content)
	counterA := createKeyCountMap(listA)
	counterB := createKeyCountMap(listB)

	for keyA, countA := range counterA {
		if countB, ok := counterB[keyA]; ok {
			if countA == listALength && countB == listBLength && countA > nonStandardIdentifierGuessCountThreshold {
				return ListItemIdentifierField(keyA)
			}
		}
	}

	return ""
}

func (compare *compare) createLookUpMap(sequenceNode *yamlv3.Node) map[uint64][]int {
	result := make(map[uint64][]int, len(sequenceNode.Content))
	for idx, entry := range sequenceNode.Content {
		hash := compare.calcNodeHash(entry)
		if _, ok := result[hash]; !ok {
			result[hash] = []int{}
		}

		result[hash] = append(result[hash], idx)
	}

	return result
}

func (compare *compare) basicType(node *yamlv3.Node) interface{} {
	switch node.Kind {
	case yamlv3.DocumentNode:
		panic("document nodes are not supported to be translated into a basic type")

	case yamlv3.MappingNode:
		result := map[interface{}]interface{}{}
		for i := 0; i < len(node.Content); i += 2 {
			k, v := followAlias(node.Content[i]), followAlias(node.Content[i+1])
			result[compare.basicType(k)] = compare.basicType(v)
		}

		return result

	case yamlv3.SequenceNode:
		result := []interface{}{}

		if compare.settings.IgnoreOrderChanges {
			sortNode(node)
		}

		for _, entry := range node.Content {
			result = append(result, compare.basicType(followAlias(entry)))
		}

		return result

	case yamlv3.ScalarNode:
		return node.Value

	case yamlv3.AliasNode:
		return compare.basicType(node.Alias)

	default:
		panic("should be unreachable")
	}
}

func (compare *compare) calcNodeHash(node *yamlv3.Node) (hash uint64) {
	var err error

	switch node.Kind {
	case yamlv3.MappingNode, yamlv3.SequenceNode:
		hash, err = hashstructure.Hash(compare.basicType(node), nil)

	case yamlv3.ScalarNode:
		hash, err = hashstructure.Hash(node.Value, nil)

	case yamlv3.AliasNode:
		hash = compare.calcNodeHash(followAlias(node))

	default:
		err = fmt.Errorf("kind %v is not supported", node.Kind)
	}

	if err != nil {
		panic(wrap.Errorf(err, "failed to calculate hash of %#v", node.Value))
	}

	return hash
}

func sortNode(node *yamlv3.Node) {
	sort.Slice(node.Content, func(i, j int) bool {
		a, b := node.Content[i], node.Content[j]

		if a.Kind != b.Kind {
			return a.Kind < b.Kind
		}

		if a.Tag != b.Tag {
			return strings.Compare(a.Tag, b.Tag) < 0
		}

		switch a.Kind {
		case yamlv3.ScalarNode:
			return strings.Compare(a.Value, b.Value) < 0
		}

		return len(a.Content) < len(b.Content)
	})
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

func isList(node *yamlv3.Node) bool {
	switch node.Kind {
	case yamlv3.SequenceNode:
		return true
	}

	return false
}

// ChangeRoot changes the root of an input file to a position inside its
// document based on the given path. Input files with more than one document are
// not supported, since they could have multiple elements with that path.
func ChangeRoot(inputFile *ytbx.InputFile, path string, useGoPatchPaths bool, translateListToDocuments bool) error {
	multipleDocuments := len(inputFile.Documents) != 1

	if multipleDocuments {
		return fmt.Errorf("change root for an input file is only possible if there is only one document, but %s contains %s",
			inputFile.Location,
			text.Plural(len(inputFile.Documents), "document"))
	}

	// For reference reasons, keep the original root level
	originalRoot := inputFile.Documents[0]

	// Find the object at the given path
	obj, err := ytbx.Grab(inputFile.Documents[0], path)
	if err != nil {
		return err
	}

	wrapInDocumentNodes := func(list []*yamlv3.Node) []*yamlv3.Node {
		result := make([]*yamlv3.Node, len(list))
		for i := range list {
			result[i] = &yamlv3.Node{
				Kind:    yamlv3.DocumentNode,
				Content: []*yamlv3.Node{list[i]},
			}
		}

		return result
	}

	if translateListToDocuments && isList(obj) {
		// Change root of input file main document to a new list of documents based on the the list that was found
		inputFile.Documents = wrapInDocumentNodes(obj.Content)

	} else {
		// Change root of input file main document to the object that was found
		inputFile.Documents = wrapInDocumentNodes([]*yamlv3.Node{obj})
	}

	// Parse path string and create nicely formatted output path
	if resolvedPath, err := ytbx.ParsePathString(path, originalRoot); err == nil {
		path = pathToString(&resolvedPath, useGoPatchPaths, multipleDocuments)
	}

	inputFile.Note = fmt.Sprintf("YAML root was changed to %s", path)

	return nil
}

func pathToString(path *ytbx.Path, useGoPatchPaths bool, showPathRoot bool) string {
	var result string

	if useGoPatchPaths {
		result = styledGoPatchPath(path)

	} else {
		result = styledDotStylePath(path)
	}

	if path != nil && showPathRoot {
		result += bunt.Sprintf("  LightSteelBlue{(%s)}", path.RootDescription())
	}

	return result
}
