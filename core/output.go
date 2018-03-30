package core

import (
	"bytes"
	"fmt"

	"github.com/HeavyWombat/color"
	"github.com/HeavyWombat/yaml"
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
// For this to work, dyff relies on modified versions of the YAML lib and the
// coloring lib we use here. The YAML lib adds ANSI styles to make keys bold.
// But this means the coloring lib needs to be able to apply styles on already
// styled text without making it look ugly.
func DiffsToHumanStyle(diffs []Diff) string {
	var output bytes.Buffer

	for _, diff := range diffs {
		GenerateHumanDiffOutput(&output, diff)
	}

	return output.String()
}

func GenerateHumanDiffOutput(output *bytes.Buffer, diff Diff) {
	output.WriteString(pathToString(diff.Path))
	output.WriteString("\n")

	for _, detail := range diff.Details {
		switch detail.Kind {
		case ADDITION:
			switch detail.To.(type) {
			case []interface{}:
				output.WriteString(Color(fmt.Sprintf("  %d entries added:\n", len(detail.To.([]interface{}))), color.FgYellow))
			}
			output.WriteString(Green(yamlString(detail.To)))

		case REMOVAL:
			switch detail.From.(type) {
			case []interface{}:
				output.WriteString(Color(fmt.Sprintf("  %d entries removed:\n", len(detail.From.([]interface{}))), color.FgYellow))
			}
			output.WriteString(Red(yamlString(detail.From)))

		case MODIFICATION:
			output.WriteString(Yellow("changed value\n"))
			output.WriteString(Red(fmt.Sprintf(" - %v\n", detail.From)))
			output.WriteString(Green(fmt.Sprintf(" + %v\n", detail.To)))
		}
	}

	output.WriteString("\n")
}
