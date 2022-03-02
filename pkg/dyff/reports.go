package dyff

import "regexp"

// Filter accepts YAML paths as input and returns a new report with differences for those paths only
func (r Report) Filter(paths ...string) (result Report) {
	if len(paths) == 0 {
		return r
	}

	result = Report{
		From: r.From,
		To:   r.To,
	}

	regexps := make([]*regexp.Regexp, len(paths))
	for i := range paths {
		regexps[i] = regexp.MustCompile(paths[i])
	}

	for _, diff := range r.Diffs {
		diffPathString := diff.Path.String()
		for _, regexp := range regexps {
			if regexp.MatchString(diffPathString) {
				result.Diffs = append(result.Diffs, diff)
				break
			}
		}
	}

	return result
}

// Exclude accepts YAML paths as input and returns a new report with differences without those paths
func (r Report) Exclude(paths ...string) (result Report) {
	if len(paths) == 0 {
		return r
	}

	result = Report{
		From: r.From,
		To:   r.To,
	}

	regexps := make([]*regexp.Regexp, len(paths))
	for i := range paths {
		regexps[i] = regexp.MustCompile(paths[i])
	}

	for _, diff := range r.Diffs {
		diffPathString := diff.Path.String()
		var any bool
		for _, regexp := range regexps {
			if regexp.MatchString(diffPathString) {
				any = true
				break
			}
		}

		if !any {
			result.Diffs = append(result.Diffs, diff)
		}
	}

	return result
}
