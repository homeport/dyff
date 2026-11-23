package dyff

import (
	"testing"

	yamlv3 "gopkg.in/yaml.v3"
)

func TestNodesEqual(t *testing.T) {
	// both nil should be equal
	if !nodesEqual(nil, nil) {
		t.Fatalf("expected two nil nodes to be equal")
	}

	// one nil and one non-nil should not be equal
	nonNil := &yamlv3.Node{Kind: yamlv3.ScalarNode, Tag: "!!str", Value: "a"}
	if nodesEqual(nonNil, nil) || nodesEqual(nil, nonNil) {
		t.Fatalf("expected nil and non-nil nodes to be different")
	}

	// different scalar values should not be equal
	n1 := &yamlv3.Node{Kind: yamlv3.ScalarNode, Tag: "!!str", Value: "a"}
	n2 := &yamlv3.Node{Kind: yamlv3.ScalarNode, Tag: "!!str", Value: "b"}
	if nodesEqual(n1, n2) {
		t.Fatalf("expected scalar nodes with different values to be different")
	}

	// equal nested trees should be equal
	child1 := &yamlv3.Node{Kind: yamlv3.ScalarNode, Tag: "!!str", Value: "child"}
	child2 := &yamlv3.Node{Kind: yamlv3.ScalarNode, Tag: "!!str", Value: "child"}
	parent1 := &yamlv3.Node{Kind: yamlv3.SequenceNode, Content: []*yamlv3.Node{child1}}
	parent2 := &yamlv3.Node{Kind: yamlv3.SequenceNode, Content: []*yamlv3.Node{child2}}
	if !nodesEqual(parent1, parent2) {
		t.Fatalf("expected parents with identical children to be equal")
	}

	// nested trees that differ in a child should not be equal
	diffChild := &yamlv3.Node{Kind: yamlv3.ScalarNode, Tag: "!!str", Value: "other"}
	parent3 := &yamlv3.Node{Kind: yamlv3.SequenceNode, Content: []*yamlv3.Node{diffChild}}
	if nodesEqual(parent1, parent3) {
		t.Fatalf("expected parents with different children to be different")
	}
}
