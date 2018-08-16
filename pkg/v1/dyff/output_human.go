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
	"bytes"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"io"
	"reflect"
	"strings"
	"unicode/utf8"

	"github.com/HeavyWombat/dyff/pkg/v1/bunt"
	"github.com/HeavyWombat/dyff/pkg/v1/neat"
	colorful "github.com/lucasb-eyer/go-colorful"
	"github.com/sergi/go-diff/diffmatchpatch"
	yaml "gopkg.in/yaml.v2"
)

const banner = `     _        __  __
   _| |_   _ / _|/ _|
 / _' | | | | |_| |_
| (_| | |_| |  _|  _|
 \__,_|\__, |_| |_|
        |___/
`

// stringWriter is the interface that wraps the WriteString method.
type stringWriter interface {
	WriteString(s string) (int, error)
}

// HumanReport is a reporter with human readable output in mind
type HumanReport struct {
	NoTableStyle      bool
	DoNotInspectCerts bool
	ShowBanner        bool

	Report
}

// WriteReport writes a human readable report to the provided writer
func (report *HumanReport) WriteReport(out io.Writer) error {
	writer := bufio.NewWriter(out)
	defer writer.Flush()

	// Only show the document index if there is more than one document to show
	showDocumentIdx := len(report.From.Documents) > 1

	// Show banner if enabled
	if report.ShowBanner {
		var stats bytes.Buffer
		stats.WriteString("\n")
		stats.WriteString(fmt.Sprintf(" between %s\n", HumanReadableLocationInformation(report.From)))
		stats.WriteString(fmt.Sprintf("     and %s\n", HumanReadableLocationInformation(report.To)))
		stats.WriteString("\n")
		stats.WriteString(fmt.Sprintf("returned %s\n", bunt.Style(Plural(len(report.Diffs), "difference"), bunt.Bold)))

		writer.WriteString(CreateTableStyleString(" ", 0, bunt.Style(banner, bunt.Bold), stats.String()))
	}

	// Loop over the diff and generate each report into the buffer
	for _, diff := range report.Diffs {
		if err := report.generateHumanDiffOutput(writer, diff, showDocumentIdx); err != nil {
			return err
		}
	}

	// Finish with one last newline so that we do not end next to the prompt
	writer.WriteString("\n")
	return nil
}

// generateHumanDiffOutput creates a human readable report of the provided diff and writes this into the given bytes buffer. There is an optional flag to indicate whether the document index (which documents of the input file) should be included in the report of the path of the difference.
func (report *HumanReport) generateHumanDiffOutput(output stringWriter, diff Diff, showDocumentIdx bool) error {
	output.WriteString("\n")
	output.WriteString(diff.Path.ToString(showDocumentIdx))
	output.WriteString("\n")

	blocks := make([]string, len(diff.Details))
	for i, detail := range diff.Details {
		generatedOutput, err := report.generateHumanDetailOutput(detail)
		if err != nil {
			return err
		}

		blocks[i] = generatedOutput
	}

	// For the use case in which only a path-less diff is suppose to be printed,
	// omit the indent in this case since there is only one element to show
	indent := 2
	if len(diff.Path.PathElements) == 0 {
		indent = 0
	}

	report.writeTextBlocks(output, indent, blocks...)
	return nil
}

// generateHumanDetailOutput only serves as a dispatcher to call the correct sub function for the respective type of change
func (report *HumanReport) generateHumanDetailOutput(detail Detail) (string, error) {
	switch detail.Kind {
	case ADDITION:
		return report.generateHumanDetailOutputAddition(detail)

	case REMOVAL:
		return report.generateHumanDetailOutputRemoval(detail)

	case MODIFICATION:
		return report.generateHumanDetailOutputModification(detail)

	case ORDERCHANGE:
		return report.generateHumanDetailOutputOrderchange(detail)
	}

	return "", fmt.Errorf("Unsupported detail type %c", detail.Kind)
}

func (report *HumanReport) generateHumanDetailOutputAddition(detail Detail) (string, error) {
	var output bytes.Buffer

	switch detail.To.(type) {
	case []interface{}:
		output.WriteString(bunt.Colorize(fmt.Sprintf("%c %s added:\n", ADDITION, Plural(len(detail.To.([]interface{})), "list entry", "list entries")), bunt.ModificationYellow))

	case yaml.MapSlice:
		output.WriteString(bunt.Colorize(fmt.Sprintf("%c %s added:\n", ADDITION, Plural(len(detail.To.(yaml.MapSlice)), "map entry", "map entries")), bunt.ModificationYellow))
	}

	yamlOutput, err := yamlStringInGreenishColors(RestructureObject(detail.To))
	if err != nil {
		return "", err
	}

	report.writeTextBlocks(&output, 2, yamlOutput)

	return output.String(), nil
}

func (report *HumanReport) generateHumanDetailOutputRemoval(detail Detail) (string, error) {
	var output bytes.Buffer

	switch detail.From.(type) {
	case []interface{}:
		output.WriteString(bunt.Colorize(fmt.Sprintf("%c %s removed:\n", REMOVAL, Plural(len(detail.From.([]interface{})), "list entry", "list entries")), bunt.ModificationYellow))

	case yaml.MapSlice:
		output.WriteString(bunt.Colorize(fmt.Sprintf("%c %s removed:\n", REMOVAL, Plural(len(detail.From.(yaml.MapSlice)), "map entry", "map entries")), bunt.ModificationYellow))
	}

	yamlOutput, err := yamlStringInRedishColors(RestructureObject(detail.From))
	if err != nil {
		return "", err
	}

	report.writeTextBlocks(&output, 2, yamlOutput)

	return output.String(), nil
}

func (report *HumanReport) generateHumanDetailOutputModification(detail Detail) (string, error) {
	var output bytes.Buffer

	fromType := determineReflectKind(detail.From)
	toType := determineReflectKind(detail.To)
	if fromType == reflect.String && toType == reflect.String {
		// delegate to special string output
		report.writeStringDiff(&output, detail.From.(string), detail.To.(string))

	} else {
		// default output
		if fromType != toType {
			output.WriteString(yellow(fmt.Sprintf("%c type change from %s to %s\n", MODIFICATION, italic(reflectKindToString(fromType)), italic(reflectKindToString(toType)))))

		} else {
			output.WriteString(yellow(fmt.Sprintf("%c value change\n", MODIFICATION)))
		}

		output.WriteString(red(fmt.Sprintf("  - %v\n", detail.From)))
		output.WriteString(green(fmt.Sprintf("  + %v\n", detail.To)))
	}

	return output.String(), nil
}

func (report *HumanReport) generateHumanDetailOutputOrderchange(detail Detail) (string, error) {
	var output bytes.Buffer

	output.WriteString(yellow(fmt.Sprintf("%c order changed\n", ORDERCHANGE)))
	switch detail.From.(type) {
	case []string:
		from := detail.From.([]string)
		to := detail.To.([]string)
		const singleLineSeparator = ", "

		threshold := getTerminalWidth() / 2
		fromSingleLineLength := stringArrayLen(from) + ((len(from) - 1) * plainTextLength(singleLineSeparator))
		toStringleLineLength := stringArrayLen(to) + ((len(to) - 1) * plainTextLength(singleLineSeparator))
		if estimatedLength := max(fromSingleLineLength, toStringleLineLength); estimatedLength < threshold {
			output.WriteString(red(fmt.Sprintf("  - %s\n", strings.Join(from, singleLineSeparator))))
			output.WriteString(green(fmt.Sprintf("  + %s\n", strings.Join(to, singleLineSeparator))))

		} else {
			output.WriteString(CreateTableStyleString(" ", 2,
				red(strings.Join(from, "\n")),
				green(strings.Join(to, "\n"))))
		}

	case []interface{}:
		fromOutput, err := yamlString(detail.From.([]interface{}))
		if err != nil {
			return "", err
		}

		toOutput, err := yamlString(detail.To.([]interface{}))
		if err != nil {
			return "", err
		}

		output.WriteString(CreateTableStyleString(" ", 2, red(fromOutput), green(toOutput)))
	}

	return output.String(), nil
}

func (report *HumanReport) writeStringDiff(output stringWriter, from string, to string) {
	// TODO Simplify code by only writing the output code once and set-up the respective strings in each if block.

	if fromCertText, toCertText, err := report.LoadX509Certs(from, to); err == nil {
		output.WriteString(yellow(fmt.Sprintf("%c certificate change\n", MODIFICATION)))
		output.WriteString(report.highlightByLine(fromCertText, toCertText))

	} else if !isValidUTF8String(from, to) {
		output.WriteString(yellow(fmt.Sprintf("%c content change\n", MODIFICATION)))
		report.writeTextBlocks(output, 0,
			createStringWithPrefix("  - ", hex.Dump([]byte(from)), bunt.RemovalRed),
			createStringWithPrefix("  + ", hex.Dump([]byte(to)), bunt.AdditionGreen))

	} else if isWhitespaceOnlyChange(from, to) {
		output.WriteString(yellow(fmt.Sprintf("%c whitespace only change\n", MODIFICATION)))
		report.writeTextBlocks(output, 0,
			createStringWithPrefix("  - ", showWhitespaceCharacters(from), bunt.RemovalRed),
			createStringWithPrefix("  + ", showWhitespaceCharacters(to), bunt.AdditionGreen))

	} else if isMultiLine(from, to) {
		output.WriteString(yellow(fmt.Sprintf("%c value change\n", MODIFICATION)))
		report.writeTextBlocks(output, 0,
			createStringWithPrefix("  - ", from, bunt.RemovalRed),
			createStringWithPrefix("  + ", to, bunt.AdditionGreen))

	} else if isMinorChange(from, to) {
		output.WriteString(yellow(fmt.Sprintf("%c value change\n", MODIFICATION)))
		diffs := diffmatchpatch.New().DiffMain(from, to, false)
		output.WriteString(highlightRemovals(diffs))
		output.WriteString(highlightAdditions(diffs))

	} else {
		output.WriteString(yellow(fmt.Sprintf("%c value change\n", MODIFICATION)))
		output.WriteString(createStringWithPrefix("  - ", from, bunt.RemovalRed))
		output.WriteString(createStringWithPrefix("  + ", to, bunt.AdditionGreen))
	}
}

func (report *HumanReport) highlightByLine(from, to string) string {
	fromLines := strings.Split(from, "\n")
	toLines := strings.Split(to, "\n")

	var buf bytes.Buffer

	if len(fromLines) == len(toLines) {
		for i := range fromLines {
			if fromLines[i] != toLines[i] {
				fromLines[i] = bunt.Colorize(fromLines[i], bunt.RemovalRed, bunt.Bold)
				toLines[i] = bunt.Colorize(toLines[i], bunt.AdditionGreen, bunt.Bold)

			} else {
				fromLines[i] = bunt.Colorize(fromLines[i], bunt.LightSalmon)
				toLines[i] = bunt.Colorize(toLines[i], bunt.LightGreen)
			}
		}

		report.writeTextBlocks(&buf, 0,
			createStringWithPrefix(bunt.Colorize("  - ", bunt.RemovalRed), strings.Join(fromLines, "\n")),
			createStringWithPrefix(bunt.Colorize("  + ", bunt.AdditionGreen), strings.Join(toLines, "\n")))

	} else {
		report.writeTextBlocks(&buf, 0,
			createStringWithPrefix("  - ", from, bunt.RemovalRed),
			createStringWithPrefix("  + ", to, bunt.AdditionGreen))
	}

	return buf.String()
}

func determineReflectKind(obj interface{}) reflect.Kind {
	if obj == nil {
		return reflect.Invalid
	}

	return reflect.TypeOf(obj).Kind()
}

func reflectKindToString(kind reflect.Kind) string {
	switch kind {
	case reflect.Invalid:
		return "<nil>"

	default:
		return kind.String()
	}
}

func highlightRemovals(diffs []diffmatchpatch.Diff) string {
	var buf bytes.Buffer

	buf.WriteString(bunt.Colorize("  - ", bunt.RemovalRed))
	for _, part := range diffs {
		switch part.Type {
		case diffmatchpatch.DiffEqual:
			buf.WriteString(bunt.Colorize(part.Text, bunt.LightSalmon))

		case diffmatchpatch.DiffDelete:
			buf.WriteString(bunt.Colorize(part.Text, bunt.RemovalRed, bunt.Bold))
		}
	}

	buf.WriteString("\n")
	return buf.String()
}

func highlightAdditions(diffs []diffmatchpatch.Diff) string {
	var buf bytes.Buffer

	buf.WriteString(bunt.Colorize("  + ", bunt.AdditionGreen))
	for _, part := range diffs {
		switch part.Type {
		case diffmatchpatch.DiffEqual:
			buf.WriteString(bunt.Colorize(part.Text, bunt.LightGreen))

		case diffmatchpatch.DiffInsert:
			buf.WriteString(bunt.Colorize(part.Text, bunt.AdditionGreen, bunt.Bold))
		}
	}

	buf.WriteString("\n")
	return buf.String()
}

// LoadX509Certs tries to load the provided strings as a cert each and returns a textual representation of the certs, or an error if the strings are not X509 certs
func (report *HumanReport) LoadX509Certs(from, to string) (string, string, error) {
	// Back out quickly if cert inspection is disabled
	if report.DoNotInspectCerts {
		return "", "", fmt.Errorf("Certificate inspection is disabled")
	}

	fromDecoded, _ := pem.Decode([]byte(from))
	if fromDecoded == nil {
		return "", "", fmt.Errorf("string '%s' is no PEM string", from)
	}

	toDecoded, _ := pem.Decode([]byte(to))
	if toDecoded == nil {
		return "", "", fmt.Errorf("string '%s' is no PEM string", to)
	}

	fromCert, err := x509.ParseCertificate(fromDecoded.Bytes)
	if err != nil {
		return "", "", err
	}

	toCert, err := x509.ParseCertificate(toDecoded.Bytes)
	if err != nil {
		return "", "", err
	}

	fromCertText := certificateSummaryAsYAML(fromCert)
	toCertText := certificateSummaryAsYAML(toCert)

	yamlStringFrom, err := yamlString(fromCertText)
	if err != nil {
		return "", "", err
	}

	yamlStringTo, err := yamlString(toCertText)
	if err != nil {
		return "", "", err
	}

	return yamlStringFrom, yamlStringTo, nil
}

// Create a YAML (hash with key/value) from a certificate to only display a few important fields (https://www.sslshopper.com/certificate-decoder.html):
//   Common Name: www.example.com
//   Organization: Company Name
//   Organization Unit: Org
//   Locality: Portland
//   State: Oregon
//   Country: US
//   Valid From: April 2, 2018
//   Valid To: April 2, 2019
//   Issuer: www.example.com, Company Name
//   Serial Number: 14581103526614300972 (0xca5a7c67490a792c)
func certificateSummaryAsYAML(cert *x509.Certificate) yaml.MapSlice {
	result := yaml.MapSlice{}

	result = append(result, yaml.MapItem{Key: "Subject", Value: yaml.MapSlice{
		yaml.MapItem{Key: "Common Name", Value: cert.Subject.CommonName},
		yaml.MapItem{Key: "Organization", Value: strings.Join(cert.Subject.Organization, " ")},
		yaml.MapItem{Key: "Organization Unit", Value: strings.Join(cert.Subject.OrganizationalUnit, " ")},
		yaml.MapItem{Key: "Locality", Value: strings.Join(cert.Subject.Locality, " ")},
		yaml.MapItem{Key: "State", Value: strings.Join(cert.Subject.Province, " ")},
		yaml.MapItem{Key: "Country", Value: strings.Join(cert.Subject.Country, " ")},
	}})

	result = append(result, yaml.MapItem{Key: "Validity Period", Value: yaml.MapSlice{
		yaml.MapItem{Key: "NotBefore", Value: cert.NotBefore.Format("Jan 2 15:04:05 2006 MST")},
		yaml.MapItem{Key: "NotAfter", Value: cert.NotAfter.Format("Jan 2 15:04:05 2006 MST")},
	}})

	result = append(result, yaml.MapItem{Key: "Issuer", Value: fmt.Sprintf("%s, %s", cert.Issuer.CommonName, strings.Join(cert.Issuer.Organization, " "))})
	result = append(result, yaml.MapItem{Key: "Serial Number", Value: fmt.Sprintf("%d (%#x)", cert.SerialNumber, cert.SerialNumber)})

	return result
}

func yamlString(input interface{}) (string, error) {
	return neat.NewOutputProcessor(false, true, nil).ToYAML(input)
}

func isWhitespaceOnlyChange(from string, to string) bool {
	return strings.Trim(from, " \n") == strings.Trim(to, " \n")
}

func showWhitespaceCharacters(text string) string {
	return strings.Replace(strings.Replace(text, "\n", bold("↵\n"), -1), " ", bold("·"), -1)
}

func createStringWithPrefix(prefix string, obj interface{}, color ...colorful.Color) string {
	var buf bytes.Buffer
	var lines = strings.Split(fmt.Sprintf("%v", obj), "\n")
	for i, line := range lines {
		if i == 0 {
			buf.WriteString(prefix)

		} else {
			buf.WriteString(strings.Repeat(" ", plainTextLength(prefix)))
		}

		buf.WriteString(line)
		buf.WriteString("\n")
	}

	switch len(color) {
	case 1:
		return bunt.Colorize(buf.String(), color[0])

	case 2:
		return bunt.ColorizeFgBg(buf.String(), color[0], color[1])

	default:
		return buf.String()
	}
}

func plainTextLength(text string) int {
	return utf8.RuneCountInString(bunt.RemoveAllEscapeSequences(text))
}

func stringArrayLen(list []string) int {
	result := 0
	for _, entry := range list {
		result += plainTextLength(entry)
	}

	return result
}

// writeTextBlocks writes strings into the provided buffer in either a table style (each string a column) or list style (each string a row)
func (report *HumanReport) writeTextBlocks(buf stringWriter, indent int, blocks ...string) {
	const separator = "   "

	// Calcuclate the theoretical maximum line length if blocks would be rendered next to each other
	theoreticalMaxLineLength := indent + ((len(blocks) - 1) * plainTextLength(separator))
	for _, block := range blocks {
		maxLineLengthInBlock := 0
		for _, line := range strings.Split(block, "\n") {
			if lineLength := plainTextLength(line); maxLineLengthInBlock < lineLength {
				maxLineLengthInBlock = lineLength
			}
		}

		theoreticalMaxLineLength += maxLineLengthInBlock
	}

	// In case the line with blocks next to each other would surpass the terminal width, fall back to the no-table-style
	if report.NoTableStyle || theoreticalMaxLineLength > getTerminalWidth() {
		for _, block := range blocks {
			lines := strings.Split(block, "\n")
			for _, line := range lines {
				buf.WriteString(strings.Repeat(" ", indent))
				buf.WriteString(line)
				buf.WriteString("\n")
			}
		}

	} else {
		buf.WriteString(CreateTableStyleString(separator, indent, blocks...))
	}
}

// CreateTableStyleString takes the multi-line input strings as columns and arranges an output string to create a table-style output format with proper padding so that the text blocks can be arranged next to each other.
func CreateTableStyleString(separator string, indent int, columns ...string) string {
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
			mtrx[x][y] = strings.Repeat(" ", max[y]+indent)
		}
	}

	for i, col := range columns {
		for j, line := range strings.Split(col, "\n") {
			mtrx[j][i] = strings.Repeat(" ", indent) + line + strings.Repeat(" ", max[i]-plainTextLength(line))
		}
	}

	var buf bytes.Buffer
	for i, row := range mtrx {
		buf.WriteString(strings.TrimRight(strings.Join(row, separator), " "))

		if i < len(mtrx)-1 {
			buf.WriteString("\n")
		}
	}

	return buf.String()
}
