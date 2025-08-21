package dyff

import (
	"regexp"

	"github.com/gonvenience/ytbx"
	yamlv3 "gopkg.in/yaml.v3"
)

func (r Report) filter(hasPath func(*ytbx.Path) bool) (result Report) {
	result = Report{
		From: r.From,
		To:   r.To,
	}

	for _, diff := range r.Diffs {
		if hasPath(diff.Path) {
			result.Diffs = append(result.Diffs, diff)
		}
	}

	return result
}

// Filter accepts YAML paths as input and returns a new report with differences for those paths only
func (r Report) Filter(paths ...string) (result Report) {
	if len(paths) == 0 {
		return r
	}

	return r.filter(func(filterPath *ytbx.Path) bool {
		for _, pathString := range paths {
			path, err := ytbx.ParsePathStringUnsafe(pathString)
			if err == nil && filterPath != nil && path.String() == filterPath.String() {
				return true
			}
		}

		return false
	})
}

// Exclude accepts YAML paths as input and returns a new report with differences without those paths
func (r Report) Exclude(paths ...string) (result Report) {
	if len(paths) == 0 {
		return r
	}

	return r.filter(func(filterPath *ytbx.Path) bool {
		for _, pathString := range paths {
			path, err := ytbx.ParsePathStringUnsafe(pathString)
			if err == nil && filterPath != nil && path.String() == filterPath.String() {
				return false
			}
		}

		return true
	})
}

// FilterRegexp accepts regular expressions as input and returns a new report with differences for matching those patterns
func (r Report) FilterRegexp(pattern ...string) (result Report) {
	if len(pattern) == 0 {
		return r
	}

	regexps := make([]*regexp.Regexp, len(pattern))
	for i := range pattern {
		regexps[i] = regexp.MustCompile(pattern[i])
	}

	return r.filter(func(filterPath *ytbx.Path) bool {
		for _, regexp := range regexps {
			if filterPath != nil && regexp.MatchString(filterPath.String()) {
				return true
			}
		}
		return false
	})
}

// ExcludeRegexp accepts regular expressions as input and returns a new report with differences for not matching those patterns
func (r Report) ExcludeRegexp(pattern ...string) (result Report) {
	if len(pattern) == 0 {
		return r
	}

	regexps := make([]*regexp.Regexp, len(pattern))
	for i := range pattern {
		regexps[i] = regexp.MustCompile(pattern[i])
	}

	result = Report{
		From: r.From,
		To:   r.To,
	}

	for _, diff := range r.Diffs {
		var shouldExclude = false

		// Check if the path itself matches any pattern
		if diff.Path != nil {
			for _, regexp := range regexps {
				if regexp.MatchString(diff.Path.String()) {
					shouldExclude = true
					break
				}
			}
		}

		// For additions and removals, also check the specific keys
		if !shouldExclude {
			for _, detail := range diff.Details {
				var node *yamlv3.Node
				if detail.Kind == ADDITION && detail.To != nil {
					node = detail.To
				} else if detail.Kind == REMOVAL && detail.From != nil {
					node = detail.From
				}

				if node != nil {
					// Construct the full path including the key
					var fullPath string
					if diff.Path != nil {
						fullPath = diff.Path.String()
						// If it's a map entry, append the key name
						if node.Kind == yamlv3.MappingNode && len(node.Content) >= 2 {
							for i := 0; i < len(node.Content); i += 2 {
								keyNode := node.Content[i]
								if keyNode.Value != "" {
									keyPath := fullPath + "." + keyNode.Value
									for _, regexp := range regexps {
										if regexp.MatchString(keyPath) {
											shouldExclude = true
											break
										}
									}
									if shouldExclude {
										break
									}
								}
							}
						} else {
							// For non-map entries, check the path directly
							for _, regexp := range regexps {
								if regexp.MatchString(fullPath) {
									shouldExclude = true
									break
								}
							}
						}
					}
					if shouldExclude {
						break
					}
				}
			}
		}

		if !shouldExclude {
			result.Diffs = append(result.Diffs, diff)
		}
	}

	return result
}



func (r Report) IgnoreValueChanges() (result Report) {
	result = Report{
		From: r.From,
		To:   r.To,
	}

	for _, diff := range r.Diffs {
		var hasValChange = false
		for _, detail := range diff.Details {
			if detail.Kind == MODIFICATION {
				hasValChange = true
				break
			}
  		}

		if !hasValChange {
			result.Diffs = append(result.Diffs, diff)
		}
	}

	return result	
}
