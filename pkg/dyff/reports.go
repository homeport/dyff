package dyff

import "regexp"

func (r Report) filter(f func(string) bool) (result Report) {
	result = Report{
		From: r.From,
		To:   r.To,
	}

	for _, diff := range r.Diffs {
		diffPathString := diff.Path.String()
		if f(diffPathString) {
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

	return r.filter(func(s string) bool {
		for _, path := range paths {
			if path == s {
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

	return r.filter(func(s string) bool {
		for _, path := range paths {
			if path == s {
				return false
			}
		}
		return true
	})
}

// FilterRegexp accepts YAML paths as input and returns a new report with differences without those paths
func (r Report) FilterRegexp(pattern ...string) (result Report) {
	if len(pattern) == 0 {
		return r
	}

	regexps := make([]*regexp.Regexp, len(pattern))
	for i := range pattern {
		regexps[i] = regexp.MustCompile(pattern[i])
	}

	return r.filter(func(s string) bool {
		for _, regexp := range regexps {
			if regexp.MatchString(s) {
				return true
			}
		}
		return false
	})
}

// ExcludeRegexp accepts YAML paths as input and returns a new report with differences without those paths
func (r Report) ExcludeRegexp(pattern ...string) (result Report) {
	if len(pattern) == 0 {
		return r
	}

	regexps := make([]*regexp.Regexp, len(pattern))
	for i := range pattern {
		regexps[i] = regexp.MustCompile(pattern[i])
	}

	return r.filter(func(s string) bool {
		for _, regexp := range regexps {
			if regexp.MatchString(s) {
				return false
			}
		}
		return true
	})
}
