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

	changedEntries := report.extractChangedEntries()

	if len(changedEntries) == 0 {
		_, _ = writer.WriteString("No changed entries found.\n")
		return nil
	}

	for listPath, entries := range changedEntries {
		// Clean up the list path for display (remove leading slash)
		displayPath := strings.TrimPrefix(listPath, "/")
		_, _ = writer.WriteString(fmt.Sprintf("# Changed entries from '%s':\n", displayPath))

		for _, entry := range entries {
			// Convert the node to YAML using RestructureObject and neat
			ytbx.RestructureObject(entry)
			yamlOutput, err := neat.NewOutputProcessor(false, true, nil).ToYAML(entry)
			if err != nil {
				return fmt.Errorf("failed to convert entry to YAML: %w", err)
			}

			// Add leading dash to make it a proper YAML list entry
			lines := strings.Split(strings.TrimSuffix(yamlOutput, "\n"), "\n")
			for i, line := range lines {
				if i == 0 {
					_, _ = writer.WriteString(fmt.Sprintf("- %s\n", line))
				} else {
					_, _ = writer.WriteString(fmt.Sprintf("  %s\n", line))
				}
			}
		}
		_, _ = writer.WriteString("\n")
	}

	return nil
}

// extractChangedEntries analyzes the diff report to find complete entries that were changed
func (report *ChangedEntriesReport) extractChangedEntries() map[string][]*yamlv3.Node {
	modifiedEntries := make(map[string][]*yamlv3.Node)
	entryPaths := make(map[string]bool) // Track unique entry paths to avoid duplicates

	for _, diff := range report.Diffs {
		if diff.Path == nil {
			continue
		}

		pathStr := diff.Path.String()

		for _, detail := range diff.Details {
			if detail.Kind == ADDITION && detail.To != nil {
				// Check if this is a list entry addition
				if detail.To.Kind == yamlv3.SequenceNode {
					// This is a sequence of entries being added
					listPath := pathStr

					// Extract all entries from the added sequence
					for _, entry := range detail.To.Content {
						if entry.Kind == yamlv3.MappingNode {
							entryKey := report.getEntryKey(listPath, entry)
							if !entryPaths[entryKey] {
								modifiedEntries[listPath] = append(modifiedEntries[listPath], entry)
								entryPaths[entryKey] = true
							}
						}
					}
				}
			} else if detail.Kind == MODIFICATION {
				// For field modifications, extract the complete entry from the "To" document
				entryPath := report.extractEntryPathFromFieldPath(pathStr)
				if entryPath != "" {
					entry := report.findEntryByPath(entryPath)
					if entry != nil {
						listPath := report.extractListPath(entryPath)
						entryKey := report.getEntryKey(listPath, entry)
						if !entryPaths[entryKey] {
							modifiedEntries[listPath] = append(modifiedEntries[listPath], entry)
							entryPaths[entryKey] = true
						}
					}
				}
			}
		}
	}

	return modifiedEntries
}

// extractEntryPathFromFieldPath extracts entry path from a field modification path
// e.g., "/allowed/image=name/container/tag" -> "/allowed/image=name/container"
func (report *ChangedEntriesReport) extractEntryPathFromFieldPath(fieldPath string) string {
	lastSlash := strings.LastIndex(fieldPath, "/")
	if lastSlash == -1 {
		return ""
	}
	return fieldPath[:lastSlash]
}

// extractListPath extracts the list name from an entry path
// e.g., "/allowed/image=name/container" -> "/allowed"
func (report *ChangedEntriesReport) extractListPath(entryPath string) string {
	parts := strings.Split(entryPath, "/")
	if len(parts) < 3 {
		return entryPath
	}
	return "/" + parts[1]
}

// getEntryKey creates a unique key for an entry to avoid duplicates
func (report *ChangedEntriesReport) getEntryKey(listPath string, entry *yamlv3.Node) string {
	identifier := report.getEntryIdentifier(entry)
	return fmt.Sprintf("%s/%s", listPath, identifier)
}

// getEntryIdentifier extracts the identifier for a list entry
func (report *ChangedEntriesReport) getEntryIdentifier(entry *yamlv3.Node) string {
	if entry.Kind != yamlv3.MappingNode {
		return ""
	}

	// Common identifier fields to check
	identifierFields := []string{"image", "name", "id", "key", "digest"}

	for i := 0; i < len(entry.Content); i += 2 {
		if i+1 < len(entry.Content) {
			key := entry.Content[i].Value
			value := entry.Content[i+1].Value

			for _, field := range identifierFields {
				if key == field {
					return fmt.Sprintf("%s=%s", key, value)
				}
			}
		}
	}

	return "unknown"
}

// findEntryByPath finds the complete entry node at the specified path in the "To" document
func (report *ChangedEntriesReport) findEntryByPath(entryPath string) *yamlv3.Node {
	// Parse paths like "/allowed/image=name/container"
	if !strings.HasPrefix(entryPath, "/") {
		return nil
	}

	// Remove leading slash
	pathWithoutSlash := entryPath[1:]

	// Find the first slash - everything before is the list name
	firstSlash := strings.Index(pathWithoutSlash, "/")
	if firstSlash == -1 {
		return nil
	}

	listName := pathWithoutSlash[:firstSlash]
	remainder := pathWithoutSlash[firstSlash+1:]

	// Now find the identifier key=value
	equalIndex := strings.Index(remainder, "=")
	if equalIndex == -1 {
		return nil
	}

	identifierKey := remainder[:equalIndex]
	identifierValue := remainder[equalIndex+1:]

	// Start from the root of the "To" document
	if len(report.To.Documents) == 0 {
		return nil
	}

	current := report.To.Documents[0]
	if current.Kind != yamlv3.DocumentNode || len(current.Content) == 0 {
		return nil
	}

	current = current.Content[0] // Get the actual document content

	// Find the list in the document
	if current.Kind == yamlv3.MappingNode {
		for i := 0; i < len(current.Content); i += 2 {
			if current.Content[i].Value == listName && i+1 < len(current.Content) {
				listNode := current.Content[i+1]
				if listNode.Kind == yamlv3.SequenceNode {
					// Look for the entry with the matching identifier
					return report.findEntryInSequenceByIdentifier(listNode, identifierKey, identifierValue)
				}
			}
		}
	}

	return nil
}

// findEntryInSequenceByIdentifier finds an entry in a sequence by identifier key-value pair
func (report *ChangedEntriesReport) findEntryInSequenceByIdentifier(sequence *yamlv3.Node, identifierKey, identifierValue string) *yamlv3.Node {
	if sequence.Kind != yamlv3.SequenceNode {
		return nil
	}

	for _, item := range sequence.Content {
		if item.Kind == yamlv3.MappingNode {
			// Look for the identifier key-value pair in this mapping
			for i := 0; i < len(item.Content); i += 2 {
				if i+1 < len(item.Content) &&
					item.Content[i].Value == identifierKey &&
					item.Content[i+1].Value == identifierValue {
					return item
				}
			}
		}
	}

	return nil
}
