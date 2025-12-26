package dyff

import (
	"testing"

	"github.com/gonvenience/ytbx"
	yamlv3 "gopkg.in/yaml.v3"
)

// TestBuildParentMapAscendAndInsertMapping ensures we can reconstruct a simple mapping path.
func TestBuildParentMapAscendAndInsertMapping(t *testing.T) {
	// Build YAML structure: a.b.c: 1
	val := &yamlv3.Node{Kind: yamlv3.ScalarNode, Tag: "!!int", Value: "1"}
	cKey := &yamlv3.Node{Kind: yamlv3.ScalarNode, Tag: "!!str", Value: "c"}
	cMap := &yamlv3.Node{Kind: yamlv3.MappingNode, Tag: "!!map", Content: []*yamlv3.Node{cKey, val}}

	bKey := &yamlv3.Node{Kind: yamlv3.ScalarNode, Tag: "!!str", Value: "b"}
	bMap := &yamlv3.Node{Kind: yamlv3.MappingNode, Tag: "!!map", Content: []*yamlv3.Node{bKey, cMap}}

	aKey := &yamlv3.Node{Kind: yamlv3.ScalarNode, Tag: "!!str", Value: "a"}
	rootMap := &yamlv3.Node{Kind: yamlv3.MappingNode, Tag: "!!map", Content: []*yamlv3.Node{aKey, bMap}}
	doc := &yamlv3.Node{Kind: yamlv3.DocumentNode, Content: []*yamlv3.Node{rootMap}}

	parentMap := buildParentMap(rootMap)
	steps := ascendPath(val, parentMap, rootMap)
	if len(steps) != 3 {
		t.Fatalf("expected 3 steps, got %d", len(steps))
	}
	if steps[0].key != "a" || steps[1].key != "b" || steps[2].key != "c" {
		t.Fatalf("unexpected keys in path order: %q, %q, %q", steps[0].key, steps[1].key, steps[2].key)
	}

	outRoot := &yamlv3.Node{Kind: yamlv3.MappingNode, Tag: "!!map"}
	insertPath(outRoot, steps)

	// outRoot should now contain a.b.c with value 1
	if len(outRoot.Content) != 2 || outRoot.Content[0].Value != "a" {
		t.Fatalf("expected top-level key 'a', got %#v", outRoot.Content)
	}
	bOut := outRoot.Content[1]
	if len(bOut.Content) != 2 || bOut.Content[0].Value != "b" {
		t.Fatalf("expected nested key 'b', got %#v", bOut.Content)
	}
	cOut := bOut.Content[1]
	if len(cOut.Content) != 2 || cOut.Content[0].Value != "c" {
		t.Fatalf("expected nested key 'c', got %#v", cOut.Content)
	}
	if got := cOut.Content[1].Value; got != "1" {
		t.Fatalf("expected final scalar '1', got %q", got)
	}

	// Sanity-check that buildChangedDocuments can use this machinery end-to-end.
	report := ChangedEntriesReport{
		Report: Report{
			To: ytbx.InputFile{Documents: []*yamlv3.Node{doc}},
			Diffs: []Diff{{
				Path:    &ytbx.Path{PathElements: []ytbx.PathElement{{Key: "a"}, {Key: "b"}, {Key: "c"}}},
				Details: []Detail{{Kind: MODIFICATION, To: val}},
			}},
		},
	}

	docs := report.buildChangedDocuments()
	if len(docs) != 1 {
		t.Fatalf("expected one changed document, got %d", len(docs))
	}
}

// TestInsertPathSequenceAndAppendIfNotPresent covers the sequence branch and appendIfNotPresent.
func TestInsertPathSequenceAndAppendIfNotPresent(t *testing.T) {
	seqParent := &yamlv3.Node{Kind: yamlv3.SequenceNode, Tag: "!!seq"}
	item := &yamlv3.Node{Kind: yamlv3.ScalarNode, Tag: "!!str", Value: "x"}
	steps := []pathStep{{parent: seqParent, index: 0, node: item}}

	outRoot := &yamlv3.Node{Kind: yamlv3.MappingNode, Tag: "!!map"}
	insertPath(outRoot, steps)
	if outRoot.Kind != yamlv3.SequenceNode {
		t.Fatalf("expected outRoot to become sequence, got kind %d", outRoot.Kind)
	}
	if len(outRoot.Content) != 1 {
		t.Fatalf("expected one item in sequence, got %d", len(outRoot.Content))
	}

	// appendIfNotPresent: same pointer should not be added again
	list := []*yamlv3.Node{item}
	list2 := appendIfNotPresent(list, item)
	if len(list2) != 1 {
		t.Fatalf("expected appendIfNotPresent to skip existing node, len=%d", len(list2))
	}

	// different node should be appended
	other := &yamlv3.Node{Kind: yamlv3.ScalarNode, Tag: "!!str", Value: "y"}
	list3 := appendIfNotPresent(list2, other)
	if len(list3) != 2 {
		t.Fatalf("expected different node to be appended, len=%d", len(list3))
	}
}

// TestComparePathSteps exercises mapping and sequence comparisons and length differences.
func TestComparePathSteps(t *testing.T) {
	mParent := &yamlv3.Node{Kind: yamlv3.MappingNode}
	sParent := &yamlv3.Node{Kind: yamlv3.SequenceNode}

	// mapping vs mapping, key order
	a := []pathStep{{parent: mParent, key: "a"}}
	b := []pathStep{{parent: mParent, key: "b"}}
	if comparePathSteps(a, b) >= 0 {
		t.Fatalf("expected 'a' < 'b'")
	}
	if comparePathSteps(b, a) <= 0 {
		t.Fatalf("expected 'b' > 'a'")
	}

	// sequence vs sequence, index order
	s1 := []pathStep{{parent: sParent, index: 0}}
	s2 := []pathStep{{parent: sParent, index: 1}}
	if comparePathSteps(s1, s2) >= 0 || comparePathSteps(s2, s1) <= 0 {
		t.Fatalf("expected index-based ordering")
	}

	// mapping < sequence
	mixed1 := []pathStep{{parent: mParent, key: "k"}}
	mixed2 := []pathStep{{parent: sParent, index: 0}}
	if comparePathSteps(mixed1, mixed2) >= 0 {
		t.Fatalf("expected mapping parent < sequence parent")
	}

	// length difference
	short := []pathStep{{parent: mParent, key: "k"}}
	long := []pathStep{{parent: mParent, key: "k"}, {parent: mParent, key: "x"}}
	if comparePathSteps(short, long) >= 0 || comparePathSteps(long, short) <= 0 {
		t.Fatalf("expected shorter path < longer path")
	}
}

// TestChangedEntriesReport_DefaultAnchorBranch ensures that the fallback anchor
// handling path (for additions that are neither sequences nor mappings) is
// executed without panicking, even though such changes currently do not
// contribute any entries to the output documents.
func TestChangedEntriesReport_DefaultAnchorBranch(t *testing.T) {
	// Single scalar document node used directly as the target of an addition.
	val := &yamlv3.Node{Kind: yamlv3.ScalarNode, Tag: "!!str", Value: "x"}
	doc := &yamlv3.Node{Kind: yamlv3.DocumentNode, Content: []*yamlv3.Node{val}}

	report := ChangedEntriesReport{
		Report: Report{
			To: ytbx.InputFile{Documents: []*yamlv3.Node{doc}},
			Diffs: []Diff{{
				Details: []Detail{{Kind: ADDITION, To: val}},
			}},
		},
	}

	docs := report.buildChangedDocuments()
	if len(docs) != 0 {
		t.Fatalf("expected no changed documents for root-level scalar addition, got %d", len(docs))
	}
}

// TestChangedEntriesReport_AdditionInMapping verifies that added mapping entries
// are turned into a minimal tree containing only the new key.
func TestChangedEntriesReport_AdditionInMapping(t *testing.T) {
	fromYAML := "---\n" +
		"root:\n" +
		"  a: 1\n"
	toYAML := "---\n" +
		"root:\n" +
		"  a: 1\n" +
		"  b: 2\n"

	fromDocs, err := ytbx.LoadYAMLDocuments([]byte(fromYAML))
	if err != nil {
		t.Fatalf("failed to load from YAML: %v", err)
	}
	toDocs, err := ytbx.LoadYAMLDocuments([]byte(toYAML))
	if err != nil {
		t.Fatalf("failed to load to YAML: %v", err)
	}

	report, err := CompareInputFiles(
		ytbx.InputFile{Documents: fromDocs},
		ytbx.InputFile{Documents: toDocs},
	)
	if err != nil {
		t.Fatalf("CompareInputFiles failed: %v", err)
	}

	changed := ChangedEntriesReport{Report: report}
	docs := changed.buildChangedDocuments()
	if len(docs) != 1 {
		t.Fatalf("expected one changed document, got %d", len(docs))
	}

	rootVal, ok := findValueByKey(docs[0], "root")
	if !ok {
		t.Fatalf("expected root mapping in changed document")
	}
	if rootVal.Kind != yamlv3.MappingNode {
		t.Fatalf("expected root value to be mapping, got kind %d", rootVal.Kind)
	}
	if len(rootVal.Content) != 2 {
		t.Fatalf("expected only new key 'b' in root mapping, got %d nodes", len(rootVal.Content))
	}
	if rootVal.Content[0].Value != "b" || rootVal.Content[1].Value != "2" {
		t.Fatalf("unexpected root mapping content: key=%q value=%q", rootVal.Content[0].Value, rootVal.Content[1].Value)
	}
}

// TestChangedEntriesReport_AdditionInSimpleList verifies that added list items
// are included as sequence entries in the result.
func TestChangedEntriesReport_AdditionInSimpleList(t *testing.T) {
	fromYAML := "---\n" +
		"list: [ A, B ]\n"
	toYAML := "---\n" +
		"list: [ A, B, C ]\n"

	fromDocs, err := ytbx.LoadYAMLDocuments([]byte(fromYAML))
	if err != nil {
		t.Fatalf("failed to load from YAML: %v", err)
	}
	toDocs, err := ytbx.LoadYAMLDocuments([]byte(toYAML))
	if err != nil {
		t.Fatalf("failed to load to YAML: %v", err)
	}

	report, err := CompareInputFiles(
		ytbx.InputFile{Documents: fromDocs},
		ytbx.InputFile{Documents: toDocs},
	)
	if err != nil {
		t.Fatalf("CompareInputFiles failed: %v", err)
	}

	changed := ChangedEntriesReport{Report: report}
	docs := changed.buildChangedDocuments()
	if len(docs) != 1 {
		t.Fatalf("expected one changed document, got %d", len(docs))
	}

	listVal, ok := findValueByKey(docs[0], "list")
	if !ok {
		t.Fatalf("expected list key in changed document")
	}
	if listVal.Kind != yamlv3.SequenceNode {
		t.Fatalf("expected list value to be sequence, got kind %d", listVal.Kind)
	}
	if len(listVal.Content) != 1 {
		t.Fatalf("expected only newly added element in list, got %d entries", len(listVal.Content))
	}
	if listVal.Content[0].Value != "C" {
		t.Fatalf("expected added list element 'C', got %q", listVal.Content[0].Value)
	}
}

// TestChangedEntriesReport_OrderChangeInSimpleList verifies that order changes
// in simple lists are reflected in the changed-entries document.
func TestChangedEntriesReport_OrderChangeInSimpleList(t *testing.T) {
	fromYAML := "---\n" +
		"list: [ A, C, B, D ]\n"
	toYAML := "---\n" +
		"list: [ A, B, C, D ]\n"

	fromDocs, err := ytbx.LoadYAMLDocuments([]byte(fromYAML))
	if err != nil {
		t.Fatalf("failed to load from YAML: %v", err)
	}
	toDocs, err := ytbx.LoadYAMLDocuments([]byte(toYAML))
	if err != nil {
		t.Fatalf("failed to load to YAML: %v", err)
	}

	report, err := CompareInputFiles(
		ytbx.InputFile{Documents: fromDocs},
		ytbx.InputFile{Documents: toDocs},
	)
	if err != nil {
		t.Fatalf("CompareInputFiles failed: %v", err)
	}

	changed := ChangedEntriesReport{Report: report}
	docs := changed.buildChangedDocuments()
	if len(docs) != 1 {
		t.Fatalf("expected one changed document, got %d", len(docs))
	}

	listVal, ok := findValueByKey(docs[0], "list")
	if !ok {
		t.Fatalf("expected list key in changed document")
	}
	if listVal.Kind != yamlv3.SequenceNode {
		t.Fatalf("expected list value to be sequence, got kind %d", listVal.Kind)
	}
	if len(listVal.Content) != 4 {
		t.Fatalf("expected four list entries involved in order change, got %d", len(listVal.Content))
	}
}
