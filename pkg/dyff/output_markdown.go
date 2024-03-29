// Copyright © 2019 The Homeport Team
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package dyff

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

// MarkdownReport is a reporter with human readable output in mind
type MarkdownReport struct {
	HumanReport
}

// WriteReport writes a human readable report to the provided writer
func (report *MarkdownReport) WriteReport(out io.Writer) error {
	writer := bufio.NewWriter(out)
	defer writer.Flush()

	// Only show the document index if there is more than one document to show
	showPathRoot := len(report.From.Documents) > 1

	// Loop over the diff and generate each report into the buffer
	for _, diff := range report.Diffs {
		if err := report.generateMarkdownDiffOutput(writer, diff, report.UseGoPatchPaths, showPathRoot); err != nil {
			return err
		}
	}

	// Finish with one last newline so that we do not end next to the prompt
	_, _ = writer.WriteString("\n")
	return nil
}

// generateMarkdownDiffOutput creates a human readable report of the provided diff and writes this into the given bytes buffer. There is an optional flag to indicate whether the document index (which documents of the input file) should be included in the report of the path of the difference.
func (report *MarkdownReport) generateMarkdownDiffOutput(output stringWriter, diff Diff, useGoPatchPaths bool, showPathRoot bool) error {
	_, _ = output.WriteString("@@ ")
	if useGoPatchPaths {
		_, _ = output.WriteString(styledGoPatchPath(diff.Path))

	} else {
		_, _ = output.WriteString(styledDotStylePath(diff.Path))
	}
	_, _ = output.WriteString(" @@\n")
	// Write the root description onto its own line
	if diff.Path != nil && showPathRoot {
		_, _ = output.WriteString(fmt.Sprintf("# %s\n", diff.Path.RootDescription()))
	}

	blocks := make([]string, len(diff.Details))
	for i, detail := range diff.Details {
		generatedOutput, err := report.generateMarkdownDetailOutput(detail)
		if err != nil {
			return err
		}

		blocks[i] = generatedOutput
	}

	// For the use case in which only a path-less diff is suppose to be printed,
	// omit the indent in this case since there is only one element to show
	indent := 0
	if diff.Path != nil && len(diff.Path.PathElements) == 0 {
		indent = 0
	}

	report.writeTextBlocks(output, indent, blocks...)
	return nil
}

// generateMarkdownDetailOutput only serves as a dispatcher to call the correct sub function for the respective type of change
func (report *MarkdownReport) generateMarkdownDetailOutput(detail Detail) (string, error) {
	switch detail.Kind {
	case ADDITION:
		return report.generateMarkdownDetailOutputAddition(detail)

	case REMOVAL:
		return report.generateMarkdownDetailOutputRemoval(detail)

	case MODIFICATION:
		return report.generateHumanDetailOutputModification(detail)

	case ORDERCHANGE:
		return report.generateHumanDetailOutputOrderchange(detail)
	}

	return "", fmt.Errorf("unsupported detail type %c", detail.Kind)
}

func (report *MarkdownReport) generateMarkdownDetailOutputAddition(detail Detail) (string, error) {

	yamlOutput, err := report.HumanReport.generateHumanDetailOutputAddition(detail)
	if err != nil {
		return "", err
	}

	var output bytes.Buffer

	// the first line is the change type, we don't want to prefix that
	firstLine := strings.Split(yamlOutput, "\n")[0]
	// Remove the ADDITION rune from the first line
	report.writeTextBlocks(&output, 0, strings.Replace(firstLine, fmt.Sprintf("%c ", ADDITION), "", 1))

	yamlOutput = strings.Replace(yamlOutput, firstLine+"\n", "", 1)

	report.writeTextBlocks(&output, 0, createStringWithContinuousPrefix("+ ", yamlOutput, 0))

	return output.String(), nil
}

func (report *MarkdownReport) generateMarkdownDetailOutputRemoval(detail Detail) (string, error) {

	yamlOutput, err := report.HumanReport.generateHumanDetailOutputRemoval(detail)
	if err != nil {
		return "", err
	}

	var output bytes.Buffer

	// the first line is the change type, we don't want to prefix that
	firstLine := strings.Split(yamlOutput, "\n")[0]
	// Replace REMOVAL rune with a markdown comment, to not highlight the change type
	report.writeTextBlocks(&output, 0, strings.Replace(firstLine, fmt.Sprintf("%c ", REMOVAL), "", 1))

	yamlOutput = strings.Replace(yamlOutput, firstLine+"\n", "", 1)

	report.writeTextBlocks(&output, 0, createStringWithContinuousPrefix("- ", yamlOutput, 0))

	return output.String(), nil
}
