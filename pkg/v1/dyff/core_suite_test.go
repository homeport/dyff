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

package dyff_test

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"

	. "github.com/HeavyWombat/dyff/pkg/v1/bunt"
	. "github.com/HeavyWombat/dyff/pkg/v1/neat"

	. "github.com/HeavyWombat/dyff/pkg/v1/dyff"
	yaml "gopkg.in/yaml.v2"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCore(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "dyff suite")
}

var _ = BeforeSuite(func() {
	ColorSetting = OFF
	LoggingLevel = NONE
	FixedTerminalWidth = 80
})

func compareAgainstExpected(fromPath string, toPath string, expectedPath string) {
	from, to, err := LoadFiles(fromPath, toPath)
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
		var content InputFile
		var err error
		if content, err = LoadFile(input); err != nil {
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
	switch doc.(type) {
	case yaml.MapSlice:
		return doc.(yaml.MapSlice)
	}

	Fail(fmt.Sprintf("Failed to use YAML, parsed data is not a YAML MapSlice:\n%s\n", input))
	return nil
}

func list(input string) []interface{} {
	doc := singleDoc(input)

	switch doc.(type) {
	case []interface{}:
		return doc.([]interface{})

	case []yaml.MapSlice:
		return SimplifyList(doc.([]yaml.MapSlice))
	}

	Fail(fmt.Sprintf("Failed to use YAML, parsed data is not a slice of any kind:\n%s\nIt was parsed as: %#v", input, doc))
	return nil
}

func singleDoc(input string) interface{} {
	docs, err := LoadYAMLDocuments([]byte(input))
	if err != nil {
		Fail(fmt.Sprintf("Failed to parse as YAML:\n%s\n\n%v", input, err))
	}

	if len(docs) > 1 {
		Fail(fmt.Sprintf("Failed to use YAML, because it contains multiple documents:\n%s\n", input))
	}

	return docs[0]
}

func multiDoc(input string) []interface{} {
	documents, err := LoadYAMLDocuments([]byte(input))
	if err != nil {
		Fail(err.Error())
	}

	return documents
}

func file(input string) InputFile {
	inputfile, err := LoadFile(input)
	if err != nil {
		Fail(fmt.Sprintf("Failed to load input file from %s: %s", input, err.Error()))
	}

	return inputfile
}

func path(path string) Path {
	// path string looks like: /additions/named-entry-list-using-id/id=new

	if path == "" {
		panic("Implementation issue: Unable to create path using an empty string")
	}

	documentIdx := 0

	result := make([]PathElement, 0)
	for i, section := range strings.Split(path, "/") {
		if i == 0 {
			if section != "" {
				if !strings.HasPrefix(section, "#") {
					panic("Implementation issue: Invalid Go-Patch style path, it cannot start with anything other than a slash, or a document idx using #<number>")
				}

				num, err := strconv.Atoi(section[1:])
				if err != nil {
					panic("Implementation issue: Invalid Go-Patch style path, document idx must be a number")
				}

				documentIdx = num
			}

			continue
		}

		keyNameSplit := strings.Split(section, "=")
		switch len(keyNameSplit) {
		case 1:
			result = append(result, PathElement{Name: keyNameSplit[0]})

		case 2:
			result = append(result, PathElement{Key: keyNameSplit[0], Name: keyNameSplit[1]})

		default:
			panic(fmt.Sprintf("Implementation issue: Invalid Go-Patch style path, path element '%s' cannot contain more than one equal sign", section))
		}
	}

	return Path{DocumentIdx: documentIdx, PathElements: result}
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
		InputFile{Documents: []interface{}{from}},
		InputFile{Documents: []interface{}{to}})

	if err != nil {
		return nil, err
	}

	return report.Diffs, nil
}

func grab(obj interface{}, path string) interface{} {
	value, err := Grab(obj, path)
	if err != nil {
		out, _ := ToYAMLString(obj)
		Fail(fmt.Sprintf("Failed to grab by path %s from %s", path, out))
	}

	return value
}

func grabError(obj interface{}, path string) string {
	value, err := Grab(obj, path)
	Expect(value).To(BeNil())
	return err.Error()
}

func pathFromString(path string, obj interface{}) Path {
	result, err := NewPath(path, obj)
	Expect(err).To(BeNil())

	return result
}
