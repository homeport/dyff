// Copyright © 2018 Matthias Diester
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
	"fmt"
	"io"

	"github.com/HeavyWombat/dyff/pkg/v1/bunt"
)

const (
	oneline = "%s detected between %s and %s\n"
	twoline = "%s detected between %s\nand %s\n"
)

// BriefReport is a reporter that only prints a summary
type BriefReport struct {
	Report
}

// WriteReport writes a brief summary to the provided writer
func (report *BriefReport) WriteReport(out io.Writer) error {
	writer := bufio.NewWriter(out)
	defer writer.Flush()

	noOfChanges := bunt.BoldText(Plural(len(report.Diffs), "change"))
	niceFrom := HumanReadableLocationInformation(report.From)
	niceTo := HumanReadableLocationInformation(report.To)

	var template string
	switch {
	case len(oneline)-6+plainTextLength(noOfChanges)+plainTextLength(niceFrom)+plainTextLength(niceTo) < getTerminalWidth():
		template = oneline

	default:
		template = twoline
	}

	writer.WriteString(fmt.Sprintf(template,
		noOfChanges,
		niceFrom,
		niceTo,
	))

	// Finish with one last newline so that we do not end next to the prompt
	writer.WriteString("\n")
	return nil
}
