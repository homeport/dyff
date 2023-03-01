package dyff

import (
	"fmt"
	"github.com/gonvenience/ytbx"
	yamlv3 "gopkg.in/yaml.v3"
)

var opMap = map[rune]string{
	ADDITION:     "add",
	REMOVAL:      "remove",
	MODIFICATION: "replace",
	ORDERCHANGE:  "reorder",
}

var kindMap = map[yamlv3.Kind]string{
	yamlv3.DocumentNode: "document",
	yamlv3.SequenceNode: "sequence",
	yamlv3.MappingNode:  "mapping",
	yamlv3.ScalarNode:   "scalar",
	yamlv3.AliasNode:    "alias",
}

type PatchOp struct {
	Op        string
	FromKind  string
	ToKind    string
	Path      string
	ToValue   yamlv3.Node
	FromValue yamlv3.Node
}

func GeneratePatch(r *Report) ([]PatchOp, error) {
	var out []PatchOp

	for _, d := range r.Diffs {
		for _, dd := range d.Details {
			po := PatchOp{
				Op:   opMap[dd.Kind],
				Path: d.Path.String(),
			}
			if dd.From != nil {
				po.FromKind = kindMap[dd.From.Kind]
				po.FromValue = *dd.From
			}
			if dd.To != nil {
				po.ToKind = kindMap[dd.To.Kind]
				po.ToValue = *dd.To
			}
			out = append(out, po)
		}
	}

	return out, nil
}

func ApplyPatch(input *yamlv3.Node, patch []PatchOp) error {
	for i, op := range patch {
		n, err := ytbx.Grab(input, op.Path)
		if err != nil {
			return fmt.Errorf("offset %d: error getting node at path: %v: %w", i, op.Path, err)
		}
		switch op.Op {
		case "add":
			n.Content = append(n.Content, op.ToValue.Content...)
		case "remove":
			err := rmNode(n, &op.FromValue)
			if err != nil {
				return fmt.Errorf("offset %d: error removing node: %v: %w", i, op.Path, err)
			}
		case "replace":
			if !matchNodes(n, &op.FromValue) {
				return fmt.Errorf("offset %d: error replacing value: from value doesn't match: %v (wanted %v)", i, n.Value, op.FromValue.Value)
			}
			n.Value = op.ToValue.Value
		case "reorder":
			if n.Kind != yamlv3.SequenceNode {
				return fmt.Errorf("offset %d: incorrect kind of node for reorder (wanted Sequence): %v", i, n.Kind)
			}
			found, idx := matchSubslice(n.Content, op.FromValue.Content)
			if !found {
				return fmt.Errorf("offset %d: error reordering: from value not found at path: %v", i, op.Path)
			}
			for i := range op.ToValue.Content {
				n.Content[idx+i] = op.ToValue.Content[i]
			}
		default:
			return fmt.Errorf("offset %d: unknown op: %v", i, op.Op)
		}
	}
	return nil
}

// rmNode finds from in n and removes it
func rmNode(n *yamlv3.Node, from *yamlv3.Node) error {
	found, idx := matchSubslice(n.Content, from.Content)
	if !found {
		return fmt.Errorf("from node not found: %+v", from)
	}
	n.Content = append(n.Content[:idx], n.Content[idx+len(from.Content):]...)
	return nil
}

// matchSubslice checks if a contains subslice b, returning index of b if found
func matchSubslice(a, b []*yamlv3.Node) (bool, int) {
	if len(b) > len(a) {
		return false, 0
	}
Loop:
	for i := range a {
		if matchNodes(a[i], b[0]) {
			for j := range b[1:] {
				if !matchNodes(a[i+j+1], b[j+1]) {
					continue Loop
				}
			}
			return true, i
		}
	}
	return false, 0
}

// matchNodes determines whether a and b are equal
func matchNodes(a, b *yamlv3.Node) bool {
	if a.Kind != b.Kind {
		return false
	}
	if len(a.Content) != len(b.Content) {
		return false
	}
	switch a.Kind {
	case yamlv3.ScalarNode:
		return a.Value == b.Value
	case yamlv3.DocumentNode:
		fallthrough
	case yamlv3.MappingNode:
		fallthrough
	case yamlv3.SequenceNode:
		for i := range a.Content {
			if !matchNodes(a.Content[i], b.Content[i]) {
				return false
			}
		}
		return true
	case yamlv3.AliasNode:
		return matchNodes(a.Alias, b.Alias)
	default:
		panic(fmt.Sprintf("unknown node kind: %v", a.Kind))
	}
}
