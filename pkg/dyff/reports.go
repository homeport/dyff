package dyff

import "github.com/gonvenience/ytbx"

// Filter accepts YAML paths as input and returns a new report with differences for those paths only
func (r Report) Filter(paths ...*ytbx.Path) (result Report) {
	if len(paths) == 0 {
		return r
	}

	result = Report{
		From: r.From,
		To:   r.To,
	}

	pathsMap := make(map[string]struct{})

	for _, path := range paths {
		pathsMap[path.String()] = struct{}{}
	}

	for _, diff := range r.Diffs {
		diffPathString := diff.Path.String()
		if _, ok := pathsMap[diffPathString]; ok {
			result.Diffs = append(result.Diffs, diff)
		}
	}

	return result
}
