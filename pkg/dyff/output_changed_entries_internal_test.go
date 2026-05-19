package dyff

import (
	"bytes"
	"fmt"
	"strings"
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
	if comparePathSteps(mixed2, mixed1) <= 0 {
		t.Fatalf("expected sequence parent > mapping parent")
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

// TestWriteReportFlushError ensures the deferred flush error handling branch is
// executed when the underlying writer fails.
type failingWriter struct{}

func (w *failingWriter) Write(p []byte) (int, error) {
	return 0, fmt.Errorf("write failed")
}

func TestWriteReportFlushError(t *testing.T) {
	val := &yamlv3.Node{Kind: yamlv3.ScalarNode, Tag: "!!str", Value: "x"}
	key := &yamlv3.Node{Kind: yamlv3.ScalarNode, Tag: "!!str", Value: "k"}
	root := &yamlv3.Node{Kind: yamlv3.MappingNode, Tag: "!!map", Content: []*yamlv3.Node{key, val}}
	doc := &yamlv3.Node{Kind: yamlv3.DocumentNode, Content: []*yamlv3.Node{root}}

	report := ChangedEntriesReport{
		Report: Report{
			To:    ytbx.InputFile{Documents: []*yamlv3.Node{doc}},
			Diffs: []Diff{{Details: []Detail{{Kind: MODIFICATION, To: val}}}},
		},
	}

	// Ensure normal YAML rendering is used so the error originates from the
	// writer flush, not from marshalToYAML.
	oldMarshal := marshalToYAML
	marshalToYAML = func(doc *yamlv3.Node) (string, error) {
		return "ok", nil
	}
	defer func() { marshalToYAML = oldMarshal }()

	var w failingWriter
	if err := report.WriteReport(&w); err == nil {
		t.Fatalf("expected error from WriteReport when underlying writer fails")
	}
}

// TestWriteReportYAMLError ensures the YAML conversion error branch is
// exercised by injecting a failing marshalToYAML implementation.
func TestWriteReportYAMLError(t *testing.T) {
	val := &yamlv3.Node{Kind: yamlv3.ScalarNode, Tag: "!!str", Value: "x"}
	key := &yamlv3.Node{Kind: yamlv3.ScalarNode, Tag: "!!str", Value: "k"}
	root := &yamlv3.Node{Kind: yamlv3.MappingNode, Tag: "!!map", Content: []*yamlv3.Node{key, val}}
	doc := &yamlv3.Node{Kind: yamlv3.DocumentNode, Content: []*yamlv3.Node{root}}

	report := ChangedEntriesReport{
		Report: Report{
			To:    ytbx.InputFile{Documents: []*yamlv3.Node{doc}},
			Diffs: []Diff{{Details: []Detail{{Kind: MODIFICATION, To: val}}}},
		},
	}

	oldMarshal := marshalToYAML
	marshalToYAML = func(doc *yamlv3.Node) (string, error) {
		return "", fmt.Errorf("marshal error")
	}
	defer func() { marshalToYAML = oldMarshal }()

	var buf bytes.Buffer
	if err := report.WriteReport(&buf); err == nil || !strings.Contains(err.Error(), "failed to convert document to YAML") {
		t.Fatalf("expected YAML conversion error from WriteReport, got: %v", err)
	}
}

// TestWriteReportMultiDocumentSeparator verifies that multi-document output
// uses the '---' separator and therefore exercises the i>0 branch.
func TestWriteReportMultiDocumentSeparator(t *testing.T) {
	aVal := &yamlv3.Node{Kind: yamlv3.ScalarNode, Tag: "!!str", Value: "1"}
	aKey := &yamlv3.Node{Kind: yamlv3.ScalarNode, Tag: "!!str", Value: "a"}
	aMap := &yamlv3.Node{Kind: yamlv3.MappingNode, Tag: "!!map", Content: []*yamlv3.Node{aKey, aVal}}
	doc0 := &yamlv3.Node{Kind: yamlv3.DocumentNode, Content: []*yamlv3.Node{aMap}}

	bVal := &yamlv3.Node{Kind: yamlv3.ScalarNode, Tag: "!!str", Value: "2"}
	bKey := &yamlv3.Node{Kind: yamlv3.ScalarNode, Tag: "!!str", Value: "b"}
	bMap := &yamlv3.Node{Kind: yamlv3.MappingNode, Tag: "!!map", Content: []*yamlv3.Node{bKey, bVal}}
	doc1 := &yamlv3.Node{Kind: yamlv3.DocumentNode, Content: []*yamlv3.Node{bMap}}

	report := ChangedEntriesReport{
		Report: Report{
			To: ytbx.InputFile{Documents: []*yamlv3.Node{doc0, doc1}},
			Diffs: []Diff{
				{Path: &ytbx.Path{DocumentIdx: 0}, Details: []Detail{{Kind: MODIFICATION, To: aVal}}},
				{Path: &ytbx.Path{DocumentIdx: 1}, Details: []Detail{{Kind: MODIFICATION, To: bVal}}},
			},
		},
	}

	var buf bytes.Buffer
	if err := report.WriteReport(&buf); err != nil {
		t.Fatalf("unexpected error from WriteReport: %v", err)
	}
	if !strings.Contains(buf.String(), "---\n") {
		t.Fatalf("expected multi-document separator '---' in output, got: %q", buf.String())
	}
}

// TestBuildChangedDocumentsSkipsNilTo ensures details with To == nil are
// ignored when collecting targets.
func TestBuildChangedDocumentsSkipsNilTo(t *testing.T) {
	val := &yamlv3.Node{Kind: yamlv3.ScalarNode, Tag: "!!str", Value: "x"}
	key := &yamlv3.Node{Kind: yamlv3.ScalarNode, Tag: "!!str", Value: "k"}
	root := &yamlv3.Node{Kind: yamlv3.MappingNode, Tag: "!!map", Content: []*yamlv3.Node{key, val}}
	doc := &yamlv3.Node{Kind: yamlv3.DocumentNode, Content: []*yamlv3.Node{root}}

	report := ChangedEntriesReport{
		Report: Report{
			To: ytbx.InputFile{Documents: []*yamlv3.Node{doc}},
			Diffs: []Diff{{
				Details: []Detail{
					{Kind: MODIFICATION, To: nil},
					{Kind: MODIFICATION, To: val},
				},
			}},
		},
	}

	docs := report.buildChangedDocuments()
	if len(docs) != 1 {
		t.Fatalf("expected one changed document when only non-nil detail contributes, got %d", len(docs))
	}
}

// TestBuildChangedDocumentsSkipsInvalidDocIndex exercises the guard against
// out-of-range document indices.
func TestBuildChangedDocumentsSkipsInvalidDocIndex(t *testing.T) {
	report := ChangedEntriesReport{
		Report: Report{
			To: ytbx.InputFile{Documents: []*yamlv3.Node{}},
			Diffs: []Diff{{
				Path:    &ytbx.Path{DocumentIdx: 1},
				Details: []Detail{{Kind: MODIFICATION, To: &yamlv3.Node{Kind: yamlv3.ScalarNode, Tag: "!!str", Value: "x"}}},
			}},
		},
	}

	docs := report.buildChangedDocuments()
	if len(docs) != 0 {
		t.Fatalf("expected no changed documents for out-of-range document index, got %d", len(docs))
	}
}

// TestBuildChangedDocumentsSkipsNilAnchor ensures the nil-anchor guard is
// exercised when a sequence of anchors contains a nil element.
func TestBuildChangedDocumentsSkipsNilAnchor(t *testing.T) {
	// Build document: list: [ A ]
	aNode := &yamlv3.Node{Kind: yamlv3.ScalarNode, Tag: "!!str", Value: "A"}
	seq := &yamlv3.Node{Kind: yamlv3.SequenceNode, Tag: "!!seq", Content: []*yamlv3.Node{aNode}}
	key := &yamlv3.Node{Kind: yamlv3.ScalarNode, Tag: "!!str", Value: "list"}
	root := &yamlv3.Node{Kind: yamlv3.MappingNode, Tag: "!!map", Content: []*yamlv3.Node{key, seq}}
	doc := &yamlv3.Node{Kind: yamlv3.DocumentNode, Content: []*yamlv3.Node{root}}

	// Detail.To is a sequence that reuses the same A node plus a nil anchor.
	seqWithNil := &yamlv3.Node{Kind: yamlv3.SequenceNode, Tag: "!!seq", Content: []*yamlv3.Node{aNode, nil}}

	report := ChangedEntriesReport{
		Report: Report{
			To: ytbx.InputFile{Documents: []*yamlv3.Node{doc}},
			Diffs: []Diff{{
				Path:    &ytbx.Path{DocumentIdx: 0},
				Details: []Detail{{Kind: ADDITION, To: seqWithNil}},
			}},
		},
	}

	docs := report.buildChangedDocuments()
	if len(docs) != 1 {
		t.Fatalf("expected one changed document when skipping nil anchors, got %d", len(docs))
	}
}

// TestAscendPathMissingParent exercises the branch where ascendPath encounters
// a node without a parent entry in the parent map.
func TestAscendPathMissingParent(t *testing.T) {
	target := &yamlv3.Node{Kind: yamlv3.ScalarNode, Tag: "!!str", Value: "x"}
	fullRoot := &yamlv3.Node{Kind: yamlv3.MappingNode, Tag: "!!map"}
	parentMap := map[*yamlv3.Node]*yamlv3.Node{}

	steps := ascendPath(target, parentMap, fullRoot)
	if len(steps) != 0 {
		t.Fatalf("expected no steps when parent map does not contain target, got %d", len(steps))
	}
}

// TestComparePathStepsEqual ensures the final return-0 path in
// comparePathSteps is exercised.
func TestComparePathStepsEqual(t *testing.T) {
	parent := &yamlv3.Node{Kind: yamlv3.MappingNode}
	pathA := []pathStep{{parent: parent, key: "k"}}
	pathB := []pathStep{{parent: parent, key: "k"}}
	if got := comparePathSteps(pathA, pathB); got != 0 {
		t.Fatalf("expected equal paths to compare as 0, got %d", got)
	}
}

// TestBuildParentMapNilRoot covers the early return in buildParentMap when the
// root node is nil.
func TestBuildParentMapNilRoot(t *testing.T) {
	parentMap := buildParentMap(nil)
	if len(parentMap) != 0 {
		t.Fatalf("expected empty parent map for nil root, got %d entries", len(parentMap))
	}
}

// TestCloneNodeNil covers the early return in cloneNode when the input node is
// nil.
func TestCloneNodeNil(t *testing.T) {
	if cloneNode(nil) != nil {
		t.Fatalf("expected cloneNode(nil) to return nil")
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
