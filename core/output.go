package core

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"unicode/utf8"

	"github.com/HeavyWombat/color"
	"github.com/HeavyWombat/yaml"
)

// NoTableStyle disables output in table style
var NoTableStyle = false

// UseGoPatchPaths style paths instead of Spruce Dot-Style
var UseGoPatchPaths = false

func pathToString(path Path) string {
	if UseGoPatchPaths {
		return ToGoPatchStyle(path)
	}

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

	if NoTableStyle {
		for _, detail := range diff.Details {
			output.WriteString(GenerateHumanDetailOutput(detail))
			output.WriteString("\n")
		}

	} else {
		cols := make([]string, 0)
		for _, detail := range diff.Details {
			cols = append(cols, GenerateHumanDetailOutput(detail))
		}

		output.WriteString(Cols("    ", 2, cols...))
	}
}

func GenerateHumanDetailOutput(detail Detail) string {
	var output bytes.Buffer

	switch detail.Kind {
	case ADDITION:
		switch detail.To.(type) {
		case []interface{}:
			output.WriteString(Color(fmt.Sprintf("%d entries added:\n", len(detail.To.([]interface{}))), color.FgYellow))
		case yaml.MapSlice:
			output.WriteString(Color(fmt.Sprintf("%d entries added:\n", len(detail.To.(yaml.MapSlice))), color.FgYellow))
		}
		output.WriteString(Green(yamlString(detail.To)))

	case REMOVAL:
		switch detail.From.(type) {
		case []interface{}:
			output.WriteString(Color(fmt.Sprintf("%d entries removed:\n", len(detail.From.([]interface{}))), color.FgYellow))
		case yaml.MapSlice:
			output.WriteString(Color(fmt.Sprintf("%d entries removed:\n", len(detail.From.(yaml.MapSlice))), color.FgYellow))

		}
		output.WriteString(Red(yamlString(detail.From)))

	case MODIFICATION:
		fromType := reflect.TypeOf(detail.From)
		toType := reflect.TypeOf(detail.To)
		if fromType != toType {
			output.WriteString(Yellow(fmt.Sprintf("changed type from %s to %s\n", Italic(fromType.String()), Italic(toType.String()))))

		} else {
			output.WriteString(Yellow("changed value\n"))
		}
		output.WriteString(Red(fmt.Sprintf(" - %v\n", detail.From)))
		output.WriteString(Green(fmt.Sprintf(" + %v\n", detail.To)))
	}

	return output.String()
}

func plainTextLength(text string) int {
	return utf8.RuneCountInString(color.RemoveAllEscapeSequences(text))
}

func Cols(separator string, intend int, columns ...string) string {
	cols := len(columns)
	rows := -1
	max := make([]int, cols)

	for i, col := range columns {
		lines := strings.Split(col, "\n")
		if noOfLines := len(lines); noOfLines > rows {
			rows = noOfLines
		}

		for _, line := range lines {
			if length := plainTextLength(line); length > max[i] {
				max[i] = length
			}
		}
	}

	mtrx := make([][]string, 0)
	for x := 0; x < rows; x++ {
		mtrx = append(mtrx, make([]string, cols))
		for y := 0; y < cols; y++ {
			mtrx[x][y] = strings.Repeat(" ", max[y]+intend)
		}
	}

	for i, col := range columns {
		for j, line := range strings.Split(col, "\n") {
			mtrx[j][i] = strings.Repeat(" ", intend) + line + strings.Repeat(" ", max[i]-plainTextLength(line))
		}
	}

	var buf bytes.Buffer
	for _, row := range mtrx {
		buf.WriteString(strings.TrimRight(strings.Join(row, separator), " "))
		buf.WriteString("\n")
	}

	return buf.String()
}
