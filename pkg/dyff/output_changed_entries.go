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
	"strings"

	"github.com/gonvenience/neat"
	"github.com/gonvenience/ytbx"
	yamlv3 "gopkg.in/yaml.v3"
)

// ChangedEntriesReport is a reporter that outputs complete final state of entries involved in changes
type ChangedEntriesReport struct {
	Report
}

// WriteReport writes the changed entries to the provided writer
func (report *ChangedEntriesReport) WriteReport(out io.Writer) error {
	writer := bufio.NewWriter(out)
	defer writer.Flush()

	documents := report.buildChangedDocuments()

	if len(documents) == 0 {
		_, _ = writer.WriteString("No changed entries found.\n")
		return nil
	}

	for i, doc := range documents {
		if i > 0 {
			_, _ = writer.WriteString("---\n")
		}

		// Convert the document to YAML
		ytbx.RestructureObject(doc)
		yamlOutput, err := neat.NewOutputProcessor(false, true, nil).ToYAML(doc)
		if err != nil {
			return fmt.Errorf("failed to convert document to YAML: %w", err)
		}

		_, _ = writer.WriteString(yamlOutput)
	}

	return nil
}

// buildChangedDocuments creates new documents containing only the changed fields with their final values
func (report *ChangedEntriesReport) buildChangedDocuments() []*yamlv3.Node {
	var documents []*yamlv3.Node

	// Group changed paths by document index
	docChanges := make(map[int]map[string]*yamlv3.Node)

	for _, diff := range report.Diffs {
		if diff.Path == nil {
			continue
		}

		pathStr := diff.Path.String()
		docIndex := 0
		if diff.Path.RootDescription() != "" && strings.Contains(diff.Path.RootDescription(), "#2") {
			docIndex = 1
		}

		// Initialize document changes if not exists
		if docChanges[docIndex] == nil {
			docChanges[docIndex] = make(map[string]*yamlv3.Node)
		}

		for _, detail := range diff.Details {
			if detail.Kind == MODIFICATION || detail.Kind == ADDITION || detail.Kind == ORDERCHANGE {
				// Get the final value from the "To" document
				finalValue := report.getFinalValueAtPath(pathStr, docIndex)
				if finalValue != nil {
					docChanges[docIndex][pathStr] = finalValue
					
					// Also capture parent objects to include all sibling fields
					report.captureParentPath(pathStr, docIndex, docChanges[docIndex])
				}
			} else if detail.Kind == REMOVAL && detail.To != nil {
				// For root level removals that result in additions (like list changes)
				finalValue := report.getFinalValueAtPath(pathStr, docIndex)
				if finalValue != nil {
					docChanges[docIndex][pathStr] = finalValue
				}
			}
		}
	}

	// Build output documents
	for docIndex := 0; docIndex < len(report.To.Documents); docIndex++ {
		if changes, hasChanges := docChanges[docIndex]; hasChanges && len(changes) > 0 {
			doc := report.buildDocumentFromChanges(changes, docIndex)
			if doc != nil {
				documents = append(documents, doc)
			}
		}
	}

	return documents
}

// captureParentPath captures the parent object when a child field changes
func (report *ChangedEntriesReport) captureParentPath(pathStr string, docIndex int, changes map[string]*yamlv3.Node) {
	// For paths like "/nil-tests/something", capture "/nil-tests" as well
	parts := strings.Split(strings.TrimPrefix(pathStr, "/"), "/")
	if len(parts) > 1 {
		parentPath := "/" + strings.Join(parts[:len(parts)-1], "/")
		if _, exists := changes[parentPath]; !exists {
			parentValue := report.getFinalValueAtPath(parentPath, docIndex)
			if parentValue != nil {
				changes[parentPath] = parentValue
			}
		}
	}
}

// getFinalValueAtPath extracts the final value at the given path from the "To" document
func (report *ChangedEntriesReport) getFinalValueAtPath(pathStr string, docIndex int) *yamlv3.Node {
	if docIndex >= len(report.To.Documents) {
		return nil
	}

	doc := report.To.Documents[docIndex]
	if doc.Kind != yamlv3.DocumentNode || len(doc.Content) == 0 {
		return nil
	}

	// Remove leading slash for ytbx.Grab
	path := strings.TrimPrefix(pathStr, "/")
	if path == "" {
		// Root level change
		return doc.Content[0]
	}

	value, err := ytbx.Grab(doc, "/"+path)
	if err != nil {
		return nil
	}

	return value
}

// buildDocumentFromChanges constructs a new document containing only the changed paths
func (report *ChangedEntriesReport) buildDocumentFromChanges(changes map[string]*yamlv3.Node, docIndex int) *yamlv3.Node {
	root := &yamlv3.Node{
		Kind: yamlv3.MappingNode,
		Tag:  "!!map",
	}

	// Process each changed path in sorted order to ensure deterministic output
	var sortedPaths []string
	for pathStr := range changes {
		sortedPaths = append(sortedPaths, pathStr)
	}
	sort.Strings(sortedPaths)
	
	for _, pathStr := range sortedPaths {
		value := changes[pathStr]
		report.setValueAtPath(root, pathStr, value)
	}

	// Return nil if no content was added
	if len(root.Content) == 0 {
		return nil
	}

	return root
}

// setValueAtPath sets a value at the specified path in the target node
func (report *ChangedEntriesReport) setValueAtPath(target *yamlv3.Node, pathStr string, value *yamlv3.Node) {
	path := strings.TrimPrefix(pathStr, "/")
	if path == "" {
		// Root level - copy content directly
		if value.Kind == yamlv3.SequenceNode {
			// Copy the sequence content
			*target = *value
		}
		return
	}

	parts := strings.Split(path, "/")
	current := target

	// Navigate/create the path
	for i, part := range parts {
		isLast := i == len(parts)-1

		if current.Kind != yamlv3.MappingNode {
			current.Kind = yamlv3.MappingNode
			current.Tag = "!!map"
		}

		// Find existing key or create new one
		var valueNode *yamlv3.Node
		found := false

		for j := 0; j < len(current.Content); j += 2 {
			if current.Content[j].Value == part {
				valueNode = current.Content[j+1]
				found = true
				break
			}
		}

		if !found {
			keyNode := &yamlv3.Node{
				Kind:  yamlv3.ScalarNode,
				Tag:   "!!str",
				Value: part,
			}
			valueNode = &yamlv3.Node{
				Kind: yamlv3.MappingNode,
				Tag:  "!!map",
			}
			current.Content = append(current.Content, keyNode, valueNode)
		}

		if isLast {
			// Set the final value
			*valueNode = *report.cloneNode(value)
		} else {
			current = valueNode
		}
	}
}

// cloneNode creates a deep copy of a YAML node
func (report *ChangedEntriesReport) cloneNode(node *yamlv3.Node) *yamlv3.Node {
	if node == nil {
		return nil
	}

	clone := &yamlv3.Node{
		Kind:        node.Kind,
		Style:       node.Style,
		Tag:         node.Tag,
		Value:       node.Value,
		Anchor:      node.Anchor,
		Alias:       node.Alias,
		HeadComment: node.HeadComment,
		LineComment: node.LineComment,
		FootComment: node.FootComment,
		Line:        node.Line,
		Column:      node.Column,
	}

	if node.Content != nil {
		clone.Content = make([]*yamlv3.Node, len(node.Content))
		for i, child := range node.Content {
			clone.Content[i] = report.cloneNode(child)
		}
	}

	return clone
}
