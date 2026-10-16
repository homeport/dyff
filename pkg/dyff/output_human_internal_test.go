package dyff

import (
	"testing"

	"github.com/gonvenience/ytbx"
)

// TestGetPlainPathString covers all branches of getPlainPathString.
func TestGetPlainPathString(t *testing.T) {
	cases := []struct {
		name string
		path *ytbx.Path
		want string
	}{
		{"nilPath", nil, ""},
		{"noElements", &ytbx.Path{}, ""},
		{"nameOnly", &ytbx.Path{PathElements: []ytbx.PathElement{{Name: "obj", Idx: -1}}}, "obj"},
		{"keyAndName", &ytbx.Path{PathElements: []ytbx.PathElement{{Key: "k", Name: "named", Idx: -1}}}, "named"},
		{"idxOnly", &ytbx.Path{PathElements: []ytbx.PathElement{{Idx: 3}}}, "3"},
		{
			"mixed",
			&ytbx.Path{PathElements: []ytbx.PathElement{
				{Key: "root", Idx: -1},
				{Name: "child", Idx: -1},
				{Key: "k", Name: "leaf", Idx: -1},
				{Idx: 7},
			}},
			"root.child.leaf.7",
		},
	}

	for _, tc := range cases {
		got := getPlainPathString(tc.path)
		if got != tc.want {
			t.Errorf("%s: got %q, want %q", tc.name, got, tc.want)
		}
	}
}
