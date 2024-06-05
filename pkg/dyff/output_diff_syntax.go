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

// DiffSyntaxReport is a reporter with human readable output in mind
type DiffSyntaxReport struct {
	PathPrefix            string
	RootDescriptionPrefix string
	ChangeTypePrefix      string
	HumanReport
}

// WriteReport writes a human readable report to the provided writer
func (report *DiffSyntaxReport) WriteReport(out io.Writer) error {
	writer := bufio.NewWriter(out)
	defer writer.Flush()

	// Only show the document index if there is more than one document to show
	showPathRoot := len(report.From.Documents) > 1

	// Loop over the diff and generate each report into the buffer
	for _, diff := range report.Diffs {
		if err := report.generateDiffSyntaxDiffOutput(writer, diff, report.UseGoPatchPaths, showPathRoot); err != nil {
			return err
		}
	}

	// Finish with one last newline so that we do not end next to the prompt
	_, _ = writer.WriteString("\n")
	return nil
}

// generatedyffSyntaxDiffOutput creates a human readable report of the provided diff and writes this into the given bytes buffer. There is an optional flag to indicate whether the document index (which documents of the input file) should be included in the report of the path of the difference.
func (report *DiffSyntaxReport) generateDiffSyntaxDiffOutput(output stringWriter, diff Diff, useGoPatchPaths bool, showPathRoot bool) error {
	_, _ = output.WriteString(fmt.Sprintf("\n%s ", report.PathPrefix))
	if useGoPatchPaths {
		_, _ = output.WriteString(styledGoPatchPath(diff.Path))

	} else {
		_, _ = output.WriteString(styledDotStylePath(diff.Path))
	}
	// Only @@ also needs a postfix
	if report.PathPrefix == "@@" {
		_, _ = output.WriteString(" @@")
	}
	_, _ = output.WriteString("\n")

	// Write the root description onto its own line
	if diff.Path != nil && showPathRoot {
		_, _ = output.WriteString(fmt.Sprintf("%s %s\n", report.RootDescriptionPrefix, diff.Path.RootDescription()))
	}

	blocks := make([]string, len(diff.Details))
	for i, detail := range diff.Details {
		generatedOutput, err := report.generateDiffSyntaxDetailOutput(detail)
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

// generatedyffSyntaxDetailOutput only serves as a dispatcher to call the correct sub function for the respective type of change
func (report *DiffSyntaxReport) generateDiffSyntaxDetailOutput(detail Detail) (string, error) {
	switch detail.Kind {
	case ADDITION:
		detailOutput, err := report.generateHumanDetailOutputAddition(detail)
		if err != nil {
			return "", err
		}
		detailOutput = report.prefixChangeType(detailOutput)

		return report.prefixChangeBlock(detailOutput, ADDITION), nil

	case REMOVAL:
		detailOutput, err := report.generateHumanDetailOutputRemoval(detail)
		if err != nil {
			return "", err
		}
		detailOutput = report.prefixChangeType(detailOutput)

		return report.prefixChangeBlock(detailOutput, REMOVAL), nil

	case MODIFICATION:
		detailOutput, err := report.generateHumanDetailOutputModification(detail)
		if err != nil {
			return "", err
		}
		return report.prefixChangeType(detailOutput), nil

	case ORDERCHANGE:
		detailOutput, err := report.generateHumanDetailOutputOrderchange(detail)
		if err != nil {
			return "", err
		}
		return report.prefixChangeType(detailOutput), nil
	}

	return "", fmt.Errorf("unsupported detail type %c", detail.Kind)
}

func (report *DiffSyntaxReport) prefixChangeType(detailOutput string) string {
	lines := strings.Split(detailOutput, "\n")
	lines[0] = strings.TrimSpace(report.ChangeTypePrefix + " " + lines[0])

	return strings.TrimSpace(strings.Join(lines, "\n"))
}

func (report *DiffSyntaxReport) prefixChangeBlock(detailOutput string, blockPrefix rune) string {
	// trim newline from the end
	detailOutput = strings.TrimSpace(detailOutput)

	var output bytes.Buffer

	// the first line is the change type, we don't want to prefix that
	firstLine := strings.Split(detailOutput, "\n")[0]
	// Remove the ADDITION rune from the first line
	report.writeTextBlocks(&output, 0, firstLine)

	detailOutput = strings.Replace(detailOutput, firstLine+"\n", "", 1)

	report.writeTextBlocks(&output, 0, createStringWithContinuousPrefix(fmt.Sprintf("%c ", blockPrefix), detailOutput, 0))

	return strings.TrimSpace(output.String())
}
