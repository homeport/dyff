// Copyright © 2020 The Homeport Team
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
	"bufio"
	"fmt"
	"io"
	"sort"

	"github.com/gonvenience/neat"
	"github.com/gonvenience/ytbx"
	yamlv3 "gopkg.in/yaml.v3"
)

// ChangedEntriesReport is a reporter that outputs complete final state of entries involved in changes
type ChangedEntriesReport struct {
	Report
}

// marshalToYAML is a small indirection around neat's YAML rendering so tests
// can inject failures and exercise error-handling branches.
var marshalToYAML = func(doc *yamlv3.Node) (string, error) {
	return neat.NewOutputProcessor(false, true, nil).ToYAML(doc)
}

// WriteReport writes the changed entries to the provided writer
func (report *ChangedEntriesReport) WriteReport(out io.Writer) (err error) {
	writer := bufio.NewWriter(out)
	defer func() {
		if flushErr := writer.Flush(); err == nil && flushErr != nil {
			err = flushErr
		}
	}()

	documents := report.buildChangedDocuments()

	if len(documents) == 0 {
		_, _ = writer.WriteString("No changed entries found.\n")
		return nil
	}

	for i, doc := range documents {
		if i > 0 {
			_, _ = writer.WriteString("---\n")
		}

		// Restructure & render
		ytbx.RestructureObject(doc)
		yamlOutput, err := marshalToYAML(doc)
		if err != nil {
			return fmt.Errorf("failed to convert document to YAML: %w", err)
		}
		_, _ = writer.WriteString(yamlOutput)
	}

	return nil
}

// buildChangedDocuments builds one output document per input document containing only changed list items / map entries
func (report *ChangedEntriesReport) buildChangedDocuments() []*yamlv3.Node {
	// parent maps per document (node ptr -> parent ptr)
	parentMaps := make([]map[*yamlv3.Node]*yamlv3.Node, len(report.To.Documents))
	for i := range report.To.Documents {
		if report.To.Documents[i] != nil && len(report.To.Documents[i].Content) > 0 {
			parentMaps[i] = buildParentMap(report.To.Documents[i].Content[0])
		}
	}

	// changed roots per document (set of nodes we want to include)
	targetsPerDoc := make([]map[*yamlv3.Node]struct{}, len(report.To.Documents))
	for i := range targetsPerDoc {
		targetsPerDoc[i] = make(map[*yamlv3.Node]struct{})
	}

	// collect nodes
	for _, diff := range report.Diffs {
		for _, detail := range diff.Details {
			if detail.Kind != MODIFICATION && detail.Kind != ADDITION && detail.Kind != ORDERCHANGE {
				continue
			}
			if detail.To == nil { // nothing in final state
				continue
			}

			// Determine document index from diff.Path if possible
			idx := 0
			if diff.Path != nil {
				idx = diff.Path.DocumentIdx
			}
			if idx >= len(report.To.Documents) || parentMaps[idx] == nil {
				continue
			}

			parentMap := parentMaps[idx]

			// Build list of candidate anchor nodes in the real document tree
			var anchors []*yamlv3.Node
			switch {
			case detail.Kind == ADDITION && detail.To.Kind == yamlv3.SequenceNode:
				// Added list entries: take each child (they are pointers into the target doc sequence)
				anchors = append(anchors, detail.To.Content...)
			case detail.Kind == ADDITION && detail.To.Kind == yamlv3.MappingNode:
				// Added mapping entries: take each value node so key+value path is reconstructed
				for i := 0; i < len(detail.To.Content); i += 2 {
					if i+1 < len(detail.To.Content) {
						anchors = append(anchors, detail.To.Content[i+1])
					}
				}
			case detail.Kind == MODIFICATION:
				anchors = append(anchors, detail.To)
			case detail.Kind == ORDERCHANGE && detail.To.Kind == yamlv3.SequenceNode:
				// Order changes: include each involved entry if they are mapping nodes referencing real list items
				for _, n := range detail.To.Content {
					if n.Kind == yamlv3.MappingNode || n.Kind == yamlv3.ScalarNode || n.Kind == yamlv3.SequenceNode {
						anchors = append(anchors, n)
					}
				}
			default:
				anchors = append(anchors, detail.To)
			}

			for _, anchor := range anchors {
				if anchor == nil { continue }
				// If anchor not part of document (no parent), skip
				if _, ok := parentMap[anchor]; !ok {
					// attempt to see if anchor is itself root (rare) – skip otherwise
					continue
				}
				// find enclosing mapping that represents list item if parent is a sequence
				candidate := anchor
				for candidate != nil {
					p := parentMap[candidate]
					if p == nil || p.Kind == yamlv3.DocumentNode {
						break
					}
					if p.Kind == yamlv3.SequenceNode { // candidate is list item
						break
					}
					candidate = p
				}
				if parent := parentMap[candidate]; parent != nil && parent.Kind == yamlv3.SequenceNode {
					// include entire list item mapping
					targetsPerDoc[idx][candidate] = struct{}{}
				} else {
					targetsPerDoc[idx][anchor] = struct{}{}
				}
			}
		}
	}

	// build output docs
	var result []*yamlv3.Node
	for docIdx, targets := range targetsPerDoc {
		if len(targets) == 0 {
			continue
		}
		rootDoc := report.To.Documents[docIdx]
		fullRoot := rootDoc.Content[0]
		parentMap := parentMaps[docIdx]

		// reconstruct minimal tree
		outRoot := &yamlv3.Node{Kind: yamlv3.MappingNode, Tag: "!!map"}
		paths := make([][]pathStep, 0, len(targets))
		for target := range targets {
			paths = append(paths, ascendPath(target, parentMap, fullRoot))
		}
		// sort paths for deterministic output
		sort.Slice(paths, func(i, j int) bool { return comparePathSteps(paths[i], paths[j]) < 0 })
		for _, p := range paths {
			insertPath(outRoot, p)
		}
		result = append(result, outRoot)
	}
	return result
}

// --- helpers ---

type pathStep struct {
	parent *yamlv3.Node
	// for mapping parent
	key string
	// for sequence parent
	index int
	// node itself
	node *yamlv3.Node
}

// ascendPath collects steps from target up to the fullRoot (excluded) then returns them top-down
func ascendPath(target *yamlv3.Node, parentMap map[*yamlv3.Node]*yamlv3.Node, fullRoot *yamlv3.Node) []pathStep {
	var rev []pathStep
	cur := target
	for cur != nil && cur != fullRoot {
		p := parentMap[cur]
		if p == nil { // reached doc root
			break
		}
		step := pathStep{parent: p, node: cur, index: -1}
		if p.Kind == yamlv3.MappingNode {
			// find key
			for i := 0; i < len(p.Content); i += 2 {
				if p.Content[i+1] == cur { step.key = p.Content[i].Value; break }
			}
		} else if p.Kind == yamlv3.SequenceNode {
			for i := 0; i < len(p.Content); i++ { if p.Content[i] == cur { step.index = i; break } }
		}
		rev = append(rev, step)
		cur = p
	}
	// now cur should be fullRoot or nil; we do not include fullRoot itself unless target IS fullRoot
	for i,j:=0,len(rev)-1; i<j; i,j = i+1,j-1 { rev[i], rev[j] = rev[j], rev[i] }
	return rev
}

func comparePathSteps(a, b []pathStep) int {
	// compare lexicographically by keys/indices sequence
	la, lb := len(a), len(b)
	for i:=0; i<la && i<lb; i++ {
		if a[i].parent.Kind == yamlv3.SequenceNode && b[i].parent.Kind == yamlv3.SequenceNode {
			if a[i].index != b[i].index { if a[i].index < b[i].index { return -1 } ; return 1 }
		} else if a[i].parent.Kind == yamlv3.MappingNode && b[i].parent.Kind == yamlv3.MappingNode {
			if a[i].key != b[i].key { if a[i].key < b[i].key { return -1 } ; return 1 }
		} else {
			// different kinds, define mapping < sequence
			if a[i].parent.Kind != b[i].parent.Kind { if a[i].parent.Kind == yamlv3.MappingNode { return -1 } ; return 1 }
		}
	}
	if la != lb { if la < lb { return -1 } ; return 1 }
	return 0
}

// insertPath inserts a path (list of steps) into outRoot cloning nodes
func insertPath(outRoot *yamlv3.Node, steps []pathStep) {
	cur := outRoot
	for i, step := range steps {
		isLast := i == len(steps)-1
		if step.parent.Kind == yamlv3.MappingNode {
			// find/create key
			var val *yamlv3.Node
			for j:=0; j < len(cur.Content); j+=2 { if cur.Content[j].Value == step.key { val = cur.Content[j+1]; break } }
			if val == nil {
				keyNode := &yamlv3.Node{Kind: yamlv3.ScalarNode, Tag: "!!str", Value: step.key}
				val = &yamlv3.Node{Kind: yamlv3.MappingNode, Tag: "!!map"}
				cur.Content = append(cur.Content, keyNode, val)
			}
			if isLast {
				*val = *cloneNode(step.node)
			} else { cur = val }
		} else if step.parent.Kind == yamlv3.SequenceNode {
			// ensure sequence node exists in cur (last mapping value added previously)
			// we model sequence parent as mapping with synthetic key? Instead compact: represent sequence as list under existing key already set in prior mapping step.
			// If current node is mapping (from previous key) and we need a sequence root, transform it.
			if cur.Kind != yamlv3.SequenceNode { cur.Kind = yamlv3.SequenceNode; cur.Tag = "!!seq" }
			// append clone of item if last (avoid duplicates: check digest of original pointer clone uniqueness by pointer) – simple append
			if isLast { cur.Content = appendIfNotPresent(cur.Content, step.node) }
		}
	}
}

func appendIfNotPresent(list []*yamlv3.Node, node *yamlv3.Node) []*yamlv3.Node {
	for _, n := range list { if n == node { return list } }
	return append(list, cloneNode(node))
}

// buildParentMap builds a child->parent map
func buildParentMap(root *yamlv3.Node) map[*yamlv3.Node]*yamlv3.Node {
	result := make(map[*yamlv3.Node]*yamlv3.Node)
	var walk func(parent, n *yamlv3.Node)
	walk = func(parent, n *yamlv3.Node) {
		if n == nil { return }
		if parent != nil { result[n] = parent }
		for _, c := range n.Content { walk(n, c) }
	}
	walk(nil, root)
	return result
}

// cloneNode deep-copies a node
func cloneNode(node *yamlv3.Node) *yamlv3.Node {
	if node == nil { return nil }
	c := *node
	if node.Content != nil { c.Content = make([]*yamlv3.Node, len(node.Content)); for i, ch := range node.Content { c.Content[i] = cloneNode(ch) } }
	return &c
}
