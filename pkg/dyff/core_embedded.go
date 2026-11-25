// Copyright © 2026 The Homeport Team
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
	"bytes"
	"errors"
	"io"
	"path/filepath"
	"strings"

	"github.com/gonvenience/ytbx"
	yamlv3 "go.yaml.in/yaml/v3"
)

// EmbeddedFiles enables comparing the YAML and JSON files embedded in
// Kubernetes ConfigMaps and Secrets (stringData) as structured data rather
// than as text, recognized by the file name extension of their key.
func EmbeddedFiles(value bool) CompareOption {
	return func(settings *compareSettings) {
		settings.EmbeddedFiles = value
	}
}

// embeddedFiles compares two string values as structured documents in case
// they are the content of a YAML or JSON file embedded in a Kubernetes
// ConfigMap or Secret. It returns false if the values are not such files, or
// if one of them cannot be parsed, in which case they are compared as text.
func (compare *compare) embeddedFiles(path ytbx.Path, from *yamlv3.Node, to *yamlv3.Node) ([]Diff, bool, error) {
	if !compare.settings.EmbeddedFiles || !isEmbeddedFile(path) {
		return nil, false, nil
	}

	fromNode, fromOK := parseEmbeddedFile(from.Value)
	toNode, toOK := parseEmbeddedFile(to.Value)
	if !fromOK || !toOK || fromNode.Kind != toNode.Kind {
		return nil, false, nil
	}

	diffs, err := compare.objects(path, fromNode, toNode)
	return diffs, true, err
}

// isEmbeddedFile checks whether the path points to an entry of a ConfigMap's
// data, or a Secret's stringData, whose key has a YAML or JSON file extension.
func isEmbeddedFile(path ytbx.Path) bool {
	if path.Root == nil || path.DocumentIdx < 0 || path.DocumentIdx >= len(path.Root.Documents) {
		return false
	}

	if len(path.PathElements) != 2 {
		return false
	}

	var field, key = path.PathElements[0].Name, path.PathElements[1].Name
	switch strings.ToLower(filepath.Ext(key)) {
	case ".yaml", ".yml", ".json":
	default:
		return false
	}

	document := path.Root.Documents[path.DocumentIdx]
	if document.Kind == yamlv3.DocumentNode && len(document.Content) == 1 {
		document = document.Content[0]
	}

	kind, ok := findValueByKey(document, "kind")
	if !ok {
		return false
	}

	switch {
	case kind.Value == "ConfigMap" && field == "data":
		return true

	case kind.Value == "Secret" && field == "stringData":
		return true

	default:
		return false
	}
}

// parseEmbeddedFile parses a file of exactly one document whose content is a
// map or a list: a plain string is valid YAML too, but nothing to compare
// structurally.
func parseEmbeddedFile(content string) (*yamlv3.Node, bool) {
	decoder := yamlv3.NewDecoder(bytes.NewReader([]byte(content)))

	var document yamlv3.Node
	if err := decoder.Decode(&document); err != nil {
		return nil, false
	}

	var next yamlv3.Node
	if err := decoder.Decode(&next); !errors.Is(err, io.EOF) {
		return nil, false
	}

	if len(document.Content) != 1 {
		return nil, false
	}

	switch node := document.Content[0]; node.Kind {
	case yamlv3.MappingNode, yamlv3.SequenceNode:
		return node, true

	default:
		return nil, false
	}
}
