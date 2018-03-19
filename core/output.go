package core

import (
	"bytes"

	yaml "gopkg.in/yaml.v2"
)

func pathToString(path Path) string {
	return ToDotStyle(path)
}

func yamlString(input interface{}) string {
	output, err := yaml.Marshal(input)
	if err != nil {
		panic(err)
	}

	return string(output)
}

// DiffsToHumanStyle creates a string with human readable report of the differences
func DiffsToHumanStyle(diffs []Diff) string {
	var output bytes.Buffer

	for _, diff := range diffs {
		output.WriteString(pathToString(diff.Path))
		output.WriteString("\n")

		switch diff.Kind {
		case ADDITION:
			output.WriteString(Green(yamlString(diff.To)))

		case REMOVAL:
			output.WriteString(Red(yamlString(diff.From)))
		}

		output.WriteString("\n")
	}

	return output.String()
}
