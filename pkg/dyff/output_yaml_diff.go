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
	"bufio"
	"io"

	"github.com/gonvenience/ytbx"
	yamlv3 "go.yaml.in/yaml/v3"
)

// YAMLDiffReport is a reporter that outputs only the differences as valid YAML
type YAMLDiffReport struct {
	Report
}

// WriteReport writes the differences as a valid YAML document to the provided writer
func (report *YAMLDiffReport) WriteReport(out io.Writer) error {
	writer := bufio.NewWriter(out)
	defer func() { _ = writer.Flush() }()

	// Create a map to store all the differences
	diffs := make(map[string]interface{})

	for _, diff := range report.Diffs {
		for _, detail := range diff.Details {
			path := diff.Path.String()

			switch detail.Kind {
			case ADDITION:
				// For additions, use the 'To' value
				if detail.To != nil {
					setValue(diffs, path, detail.To)
				}
			case REMOVAL:
				// For removals, we could either omit or mark as null
				// Here we omit them since they don't exist in the 'to' version
				continue
			case MODIFICATION:
				// For modifications, use the 'To' value (the new value)
				if detail.To != nil {
					setValue(diffs, path, detail.To)
				}
			case ORDERCHANGE:
				// For order changes, use the 'To' value
				if detail.To != nil {
					setValue(diffs, path, detail.To)
				}
			}
		}
	}

	// Write the differences as YAML
	encoder := yamlv3.NewEncoder(writer)
	encoder.SetIndent(2)

	if err := encoder.Encode(diffs); err != nil {
		return err
	}

	return encoder.Close()
}

// setValue sets a value in the nested map structure based on the path
func setValue(root map[string]interface{}, path string, node *yamlv3.Node) {
	// Parse the path and navigate through the structure
	parsedPath, err := ytbx.ParsePathStringUnsafe(path)
	if err != nil {
		return
	}

	current := root
	pathElements := parsedPath.PathElements

	for i, element := range pathElements {
		isLast := i == len(pathElements)-1
		key := element.Name

		if element.Key != "" {
			key = element.Key
		}

		if isLast {
			// Set the value at the final key
			var value interface{}
			_ = node.Decode(&value)
			current[key] = value
		} else {
			// Navigate or create intermediate maps
			if _, exists := current[key]; !exists {
				current[key] = make(map[string]interface{})
			}
			// Type assertion to continue navigation
			if next, ok := current[key].(map[string]interface{}); ok {
				current = next
			}
		}
	}
}
