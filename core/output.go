package core

import (
	"bytes"
	"fmt"

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

	switch diff.Kind {
	case ADDITION:
		output.WriteString(Green(yamlString(diff.To)))

	case REMOVAL:
		output.WriteString(Red(yamlString(diff.From)))

	case MODIFICATION:
		output.WriteString(fmt.Sprintf("changed from %s to %s", diff.From, diff.To))
	}

	output.WriteString("\n")
}
