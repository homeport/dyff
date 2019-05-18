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

package dyff_test

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"testing"

	"github.com/homeport/ytbx/pkg/v1/ytbx"
	yaml "gopkg.in/yaml.v2"

	. "github.com/homeport/dyff/pkg/v1/dyff"
	. "github.com/homeport/gonvenience/pkg/v1/term"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCore(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "dyff suite")
}

var _ = BeforeSuite(func() {
	LoggingLevel = NONE
	FixedTerminalWidth = 80
})

func compareAgainstExpected(fromPath string, toPath string, expectedPath string) {
	from, to, err := ytbx.LoadFiles(fromPath, toPath)
	Expect(err).To(BeNil())

	expected, err := ioutil.ReadFile(expectedPath)
	Expect(err).To(BeNil())

	report, err := CompareInputFiles(from, to)
	Expect(report).ToNot(BeNil())
	Expect(err).To(BeNil())

	reportWriter := &HumanReport{
		Report:            report,
		DoNotInspectCerts: false,
		NoTableStyle:      false,
		ShowBanner:        false,
	}

	buffer := &bytes.Buffer{}
	writer := bufio.NewWriter(buffer)
	reportWriter.WriteReport(writer)
	writer.Flush()

	Expect(string(expected)).To(BeIdenticalTo(buffer.String()))
}

func yml(input string) yaml.MapSlice {
	// If input is a file loacation, load this as YAML
	if _, err := os.Open(input); err == nil {
		var content ytbx.InputFile
		var err error
		if content, err = ytbx.LoadFile(input); err != nil {
			Fail(fmt.Sprintf("Failed to load YAML MapSlice from '%s': %s", input, err.Error()))
		}

		if len(content.Documents) > 1 {
			Fail(fmt.Sprintf("Failed to load YAML MapSlice from '%s': Provided file contains more than one document", input))
		}

		switch content.Documents[0].(type) {
		case yaml.MapSlice:
			return content.Documents[0].(yaml.MapSlice)
		}

		Fail(fmt.Sprintf("Failed to load YAML MapSlice from '%s': Document #0 in YAML is not of type MapSlice, but is %s", input, reflect.TypeOf(content.Documents[0])))
	}

	// Load YAML by parsing the actual string as YAML if it was not a file location
	doc := singleDoc(input)
	switch mapslice := doc.(type) {
	case yaml.MapSlice:
		return mapslice
	}

	Fail(fmt.Sprintf("Failed to use YAML, parsed data is not a YAML MapSlice:\n%s\n", input))
	return nil
}

func list(input string) []interface{} {
	doc := singleDoc(input)

	switch obj := doc.(type) {
	case []interface{}:
		return obj

	case []yaml.MapSlice:
		return ytbx.SimplifyList(obj)
	}

	Fail(fmt.Sprintf("Failed to use YAML, parsed data is not a slice of any kind:\n%s\nIt was parsed as: %#v", input, doc))
	return nil
}

func singleDoc(input string) interface{} {
	docs, err := ytbx.LoadYAMLDocuments([]byte(input))
	if err != nil {
		Fail(fmt.Sprintf("Failed to parse as YAML:\n%s\n\n%v", input, err))
	}

	if len(docs) > 1 {
		Fail(fmt.Sprintf("Failed to use YAML, because it contains multiple documents:\n%s\n", input))
	}

	return docs[0]
}

func multiDoc(input string) []interface{} {
	documents, err := ytbx.LoadYAMLDocuments([]byte(input))
	if err != nil {
		Fail(err.Error())
	}

	return documents
}

func file(input string) ytbx.InputFile {
	inputfile, err := ytbx.LoadFile(input)
	if err != nil {
		Fail(fmt.Sprintf("Failed to load input file from %s: %s", input, err.Error()))
	}

	return inputfile
}

func path(path string) ytbx.Path {
	re := regexp.MustCompile(`^(#(\d+))?(/.+)$`)

	captures := re.FindStringSubmatch(path)
	docIdxStr, pathString := captures[2], captures[3]

	result, err := ytbx.ParseGoPatchStylePathString(pathString)
	if err != nil {
		Fail(err.Error())
	}

	if len(docIdxStr) > 0 {
		num, err := strconv.Atoi(docIdxStr)
		if err != nil {
			Fail(err.Error())
		}

		result.DocumentIdx = num
	}

	return result
}

func humanDiff(diff Diff) string {
	reporter := HumanReport{
		Report:            Report{Diffs: []Diff{diff}},
		DoNotInspectCerts: false,
		NoTableStyle:      false,
		ShowBanner:        false,
	}

	var buf bytes.Buffer
	if err := reporter.WriteReport(&buf); err != nil {
		Fail(err.Error())
	}

	return buf.String()
}

func singleDiff(p string, change rune, from, to interface{}) Diff {
	return Diff{
		Path: path(p),
		Details: []Detail{{
			Kind: change,
			From: from,
			To:   to,
		}},
	}
}

func doubleDiff(p string, change1 rune, from1, to1 interface{}, change2 rune, from2, to2 interface{}) Diff {
	return Diff{
		Path: path(p),
		Details: []Detail{{
			Kind: change1,
			From: from1,
			To:   to1,
		},
			{
				Kind: change2,
				From: from2,
				To:   to2,
			}},
	}
}

func compare(from interface{}, to interface{}) ([]Diff, error) {
	report, err := CompareInputFiles(
		ytbx.InputFile{Documents: []interface{}{from}},
		ytbx.InputFile{Documents: []interface{}{to}})

	if err != nil {
		return nil, err
	}

	return report.Diffs, nil
}
