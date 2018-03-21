package core

import (
	"bytes"

	"github.com/HeavyWombat/yaml"
	"github.com/fatih/color"
)

func pathToString(path Path) string {
	return ToDotStyle(path)
}

func yamlString(input interface{}) string {
	// disable coloring if needed during YAML string generation
	prev := color.NoColor
	color.NoColor = true
	defer func() { color.NoColor = prev }()

	// TODO Write code to detect color sequences in the target string in
	// order to check whether we can write some code to merge different
	// color styles. For exmple: string contains "bold" parts and we want
	// to overwrite with "green" so that the previous "bold" parts remain
	// in bold style, but with additional green.

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
