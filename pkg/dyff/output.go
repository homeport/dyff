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
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"reflect"
	"strings"
	"unicode/utf8"

	"github.com/HeavyWombat/color"
	"github.com/HeavyWombat/yaml"
)

// TODO Separate code into different output source files: human, and the new stuff

// NoTableStyle disables output in table style
var NoTableStyle = false

// DoNotInspectCerts disables certificates inspection (compare text only)
var DoNotInspectCerts = false

// UseGoPatchPaths style paths instead of Spruce Dot-Style
var UseGoPatchPaths = false

func pathToString(path Path, showDocumentIdx bool) string {
	if UseGoPatchPaths {
		return ToGoPatchStyle(path, showDocumentIdx)
	}

	return ToDotStyle(path, showDocumentIdx)
}

func yamlString(input interface{}) string {
	output, err := yaml.Marshal(input)
	if err != nil {
		ExitWithError("Failed to marshal input object", err)
	}

	return string(output)
}

// DiffsToHumanStyle creates a string with human readable report of the differences
// For this to work, dyff relies on modified versions of the YAML lib and the
// coloring lib we use here. The YAML lib adds ANSI styles to make keys bold.
// But this means the coloring lib needs to be able to apply styles on already
// styled text without making it look ugly.
func DiffsToHumanStyle(diffs []Diff) string {
	// Map the different document indicies listed in the paths of the diffs
	counterMap := make(map[int]struct{}, 0)
	for _, diff := range diffs {
		counterMap[diff.Path.DocumentIdx] = struct{}{}
	}

	// Only show the document index if there is more than one document to show
	showDocumentIdx := len(counterMap) != 1

	// Again, loop over the diff and generate each report into the buffer
	var output bytes.Buffer
	for _, diff := range diffs {
		GenerateHumanDiffOutput(&output, diff, showDocumentIdx)
	}

	// Finish with one last newline so that we do not end next to the prompt
	output.WriteString("\n")
	return output.String()
}

func GenerateHumanDiffOutput(output *bytes.Buffer, diff Diff, showDocumentIdx bool) {
	output.WriteString("\n")
	output.WriteString(pathToString(diff.Path, showDocumentIdx))
	output.WriteString("\n")

	blocks := make([]string, len(diff.Details))
	for i, detail := range diff.Details {
		blocks[i] = GenerateHumanDetailOutput(detail)
	}

	// For the use case in which only a path-less diff is suppose to be printed,
	// omit the indent in this case since there is only one element to show
	indent := 2
	if len(diff.Path.PathElements) == 0 {
		indent = 0
	}

	WriteTextBlocks(output, indent, blocks...)
}

func GenerateHumanDetailOutput(detail Detail) string {
	var output bytes.Buffer

	// TODO Externalise part of code into separate functions for readability
	switch detail.Kind {
	case ADDITION:
		switch detail.To.(type) {
		case []interface{}:
			output.WriteString(Color(fmt.Sprintf("%c %s added:\n", ADDITION, Plural(len(detail.To.([]interface{})), "list entry", "list entries")), color.FgYellow))
		case yaml.MapSlice:
			output.WriteString(Color(fmt.Sprintf("%c %s added:\n", ADDITION, Plural(len(detail.To.(yaml.MapSlice)), "map entry", "map entries")), color.FgYellow))
		}
		WriteTextBlocks(&output, 2, Green(yamlString(RestructureObject(detail.To))))

	case REMOVAL:
		switch detail.From.(type) {
		case []interface{}:
			output.WriteString(Color(fmt.Sprintf("%c %s removed:\n", REMOVAL, Plural(len(detail.From.([]interface{})), "list entry", "list entries")), color.FgYellow))
		case yaml.MapSlice:
			output.WriteString(Color(fmt.Sprintf("%c %s removed:\n", REMOVAL, Plural(len(detail.From.(yaml.MapSlice)), "map entry", "map entries")), color.FgYellow))

		}
		WriteTextBlocks(&output, 2, Red(yamlString(RestructureObject(detail.From))))

	case MODIFICATION:
		fromType := reflect.TypeOf(detail.From).Kind()
		toType := reflect.TypeOf(detail.To).Kind()
		if fromType == reflect.String && toType == reflect.String {
			// delegate to special string output
			writeStringDiff(&output, detail.From.(string), detail.To.(string))

		} else {
			// default output
			if fromType != toType {
				output.WriteString(Yellow(fmt.Sprintf("%c type change from %s to %s\n", MODIFICATION, Italic(fromType.String()), Italic(toType.String()))))

			} else {
				output.WriteString(Yellow(fmt.Sprintf("%c value change\n", MODIFICATION)))
			}

			output.WriteString(Red(fmt.Sprintf("  - %v\n", detail.From)))
			output.WriteString(Green(fmt.Sprintf("  + %v\n", detail.To)))

		}

	case ORDERCHANGE:
		output.WriteString(Yellow(fmt.Sprintf("%c order changed\n", ORDERCHANGE)))
		switch detail.From.(type) {
		case []string:
			from := detail.From.([]string)
			to := detail.To.([]string)
			const singleLineSeparator = ", "

			threshold := getTerminalWidth() / 2
			fromSingleLineLength := stringArrayLen(from) + ((len(from) - 1) * plainTextLength(singleLineSeparator))
			toStringleLineLength := stringArrayLen(to) + ((len(to) - 1) * plainTextLength(singleLineSeparator))
			if estimatedLength := max(fromSingleLineLength, toStringleLineLength); estimatedLength < threshold {
				output.WriteString(Red(fmt.Sprintf("  - %s\n", strings.Join(from, singleLineSeparator))))
				output.WriteString(Green(fmt.Sprintf("  + %s\n", strings.Join(to, singleLineSeparator))))

			} else {
				output.WriteString(Cols(" ", 2,
					Red(fmt.Sprintf("%s", strings.Join(from, "\n"))),
					Green(fmt.Sprintf("%s", strings.Join(to, "\n")))))
			}

		case []interface{}:
			output.WriteString(Cols(" ", 2,
				Red(fmt.Sprintf("%s", yamlString(detail.From.([]interface{})))),
				Green(fmt.Sprintf("%s", yamlString(detail.To.([]interface{}))))))
		}
	}

	return output.String()
}

func writeStringDiff(output *bytes.Buffer, from string, to string) {
	// TODO Simplify code by only writing the output code once and
	// set-up the respective strings in each if block.

	if fromCertText, toCertText, err := LoadX509Certs(from, to); err == nil {
		output.WriteString(Yellow(fmt.Sprintf("%c certificate change\n", MODIFICATION)))
		WriteTextBlocks(output, 0,
			createStringWithPrefix("  - ", fromCertText, color.FgRed),
			createStringWithPrefix("  + ", toCertText, color.FgGreen))

	} else if isWhitespaceOnlyChange(from, to) {
		output.WriteString(Yellow(fmt.Sprintf("%c whitespace only change\n", MODIFICATION)))
		WriteTextBlocks(output, 0,
			createStringWithPrefix("  - ", showWhitespaceCharacters(from), color.FgRed),
			createStringWithPrefix("  + ", showWhitespaceCharacters(to), color.FgGreen))

	} else if isMinorChange(from, to) {
		// TODO Highlight the actual change more than the common part using https://github.com/sergi/go-diff DiffCommonPrefix and DiffCommonSuffix
		output.WriteString(Yellow(fmt.Sprintf("%c value change\n", MODIFICATION)))
		output.WriteString(createStringWithPrefix("  - ", from, color.FgRed))
		output.WriteString(createStringWithPrefix("  + ", to, color.FgGreen))

	} else if isMultiLine(from, to) {
		output.WriteString(Yellow(fmt.Sprintf("%c value change\n", MODIFICATION)))
		WriteTextBlocks(output, 0,
			createStringWithPrefix("  - ", from, color.FgRed),
			createStringWithPrefix("  + ", to, color.FgGreen))

	} else {
		output.WriteString(Yellow(fmt.Sprintf("%c value change\n", MODIFICATION)))
		output.WriteString(createStringWithPrefix("  - ", from, color.FgRed))
		output.WriteString(createStringWithPrefix("  + ", to, color.FgGreen))
	}
}

// LoadX509Certs tries to load the provided strings as a cert each and returns a textual representation of the certs, or an error if the strings are not X509 certs
func LoadX509Certs(from, to string) (string, string, error) {
	// Back out quickly if cert inspection is disabled
	if DoNotInspectCerts {
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

	return yamlString(fromCertText), yamlString(toCertText), nil
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
	result = append(result, yaml.MapItem{Key: "Common Name", Value: cert.Subject.CommonName})
	result = append(result, yaml.MapItem{Key: "Organization", Value: strings.Join(cert.Subject.Organization, " ")})
	result = append(result, yaml.MapItem{Key: "Organization Unit", Value: strings.Join(cert.Subject.OrganizationalUnit, " ")})
	result = append(result, yaml.MapItem{Key: "Locality", Value: strings.Join(cert.Subject.Locality, " ")})
	result = append(result, yaml.MapItem{Key: "State", Value: strings.Join(cert.Subject.Province, " ")})
	result = append(result, yaml.MapItem{Key: "Country", Value: strings.Join(cert.Subject.Country, " ")})
	result = append(result, yaml.MapItem{Key: "Valid From", Value: cert.NotBefore.Format("Jan 2 15:04:05 2006 MST")})
	result = append(result, yaml.MapItem{Key: "Valid To", Value: cert.NotAfter.Format("Jan 2 15:04:05 2006 MST")})
	result = append(result, yaml.MapItem{Key: "Issuer", Value: fmt.Sprintf("%s, %s", cert.Issuer.CommonName, strings.Join(cert.Issuer.Organization, " "))})
	result = append(result, yaml.MapItem{Key: "Serial Number", Value: fmt.Sprintf("%d (%#x)", cert.SerialNumber, cert.SerialNumber)})

	return result
}

func isWhitespaceOnlyChange(from string, to string) bool {
	return strings.Trim(from, " \n") == strings.Trim(to, " \n")
}

func showWhitespaceCharacters(text string) string {
	return strings.Replace(strings.Replace(text, "\n", Bold("↵\n"), -1), " ", Bold("·"), -1)
}

func createStringWithPrefix(prefix string, obj interface{}, attributes ...color.Attribute) string {
	var buf bytes.Buffer
	var lines = strings.Split(fmt.Sprintf("%v", obj), "\n")
	for i, line := range lines {
		if i == 0 {
			buf.WriteString(Color(prefix, color.Bold))

		} else {
			buf.WriteString(strings.Repeat(" ", len(prefix)))
		}

		buf.WriteString(line)
		buf.WriteString("\n")
	}

	return Color(buf.String(), attributes...)
}

func plainTextLength(text string) int {
	return utf8.RuneCountInString(color.RemoveAllEscapeSequences(text))
}

func stringArrayLen(list []string) int {
	result := 0
	for _, entry := range list {
		result += plainTextLength(entry)
	}

	return result
}

// WriteTextBlocks writes strings into the provided buffer in either a table style (each string a column) or list style (each string a row)
func WriteTextBlocks(buf *bytes.Buffer, indent int, blocks ...string) {
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
	if NoTableStyle || theoreticalMaxLineLength > getTerminalWidth() {
		for _, block := range blocks {
			lines := strings.Split(block, "\n")
			for _, line := range lines {
				buf.WriteString(strings.Repeat(" ", indent))
				buf.WriteString(line)
				buf.WriteString("\n")
			}
		}

	} else {
		buf.WriteString(Cols(separator, indent, blocks...))
	}
}

func Cols(separator string, indent int, columns ...string) string {
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
