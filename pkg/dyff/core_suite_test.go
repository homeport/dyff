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
	"path/filepath"
	"regexp"
	"strconv"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"

	. "github.com/homeport/dyff/pkg/dyff"

	"github.com/davecgh/go-spew/spew"
	"github.com/gonvenience/term"
	"github.com/gonvenience/ytbx"
	yamlv3 "gopkg.in/yaml.v3"
)

func TestCore(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "dyff core package suite")
}

var _ = BeforeSuite(func() {
	term.FixedTerminalWidth = 80
})

func loadFiles(fromPath string, toPath string) (ytbx.InputFile, ytbx.InputFile) {
	from, to, err := ytbx.LoadFiles(fromPath, toPath)
	Expect(err).To(BeNil())
	Expect(from).ToNot(BeNil())
	Expect(to).ToNot(BeNil())

	return from, to
}

func assets(pathElement ...string) string {
	targetPath := filepath.Join(append(
		[]string{"..", "..", "assets"},
		pathElement...,
	)...)

	abs, err := filepath.Abs(targetPath)
	if err != nil {
		return targetPath
	}

	return abs
}

func BeLike(expected interface{}) types.GomegaMatcher {
	return &extendedStringMatcher{
		expected: expected,
	}
}

type extendedStringMatcher struct {
	expected interface{}
}

func (matcher *extendedStringMatcher) Match(actual interface{}) (success bool, err error) {
	actualString, ok := actual.(string)
	if !ok {
		return false, fmt.Errorf("BeLike matcher expected a string, not %T", actual)
	}

	expectedString, ok := matcher.expected.(string)
	if !ok {
		return false, fmt.Errorf("BeLike matcher expected a string, not %T", actual)
	}

	return actualString == expectedString, nil
}

func (matcher *extendedStringMatcher) FailureMessage(actual interface{}) string {
	return fmt.Sprintf("Expected\n\t%#v\nto be like\n\t%#v",
		actual,
		matcher.expected)
}

func (matcher *extendedStringMatcher) NegatedFailureMessage(actual interface{}) string {
	return fmt.Sprintf("Expected\n\t%#v\nnot to be like\n\t%#v",
		actual,
		matcher.expected,
	)
}

func compareAgainstExpected(fromPath string, toPath string, expectedPath string, useGoPatchPaths bool, compareOptions ...CompareOption) {
	from, to, err := ytbx.LoadFiles(fromPath, toPath)
	Expect(err).To(BeNil())

	rawBytes, err := ioutil.ReadFile(expectedPath)
	Expect(err).To(BeNil())

	report, err := CompareInputFiles(from, to, compareOptions...)
	Expect(report).ToNot(BeNil())
	Expect(err).To(BeNil())

	reportWriter := &HumanReport{
		Report:            report,
		DoNotInspectCerts: false,
		NoTableStyle:      false,
		OmitHeader:        true,
		UseGoPatchPaths:   useGoPatchPaths,
	}

	buffer := &bytes.Buffer{}
	writer := bufio.NewWriter(buffer)
	Expect(reportWriter.WriteReport(writer)).To(BeNil())
	Expect(writer.Flush()).To(BeNil())

	expected := fmt.Sprintf("%#v", string(rawBytes))
	actual := fmt.Sprintf("%#v", buffer.String())
	Expect(expected).To(BeLike(actual))
}

func yml(input string) *yamlv3.Node {
	// If input is a file loacation, load this as YAML
	if _, err := os.Open(input); err == nil {
		var content ytbx.InputFile
		var err error
		if content, err = ytbx.LoadFile(input); err != nil {
			Fail(fmt.Sprintf("Failed to load YAML from '%s': %s", input, err.Error()))
		}

		if len(content.Documents) > 1 {
			Fail(fmt.Sprintf("Failed to load YAML from '%s': Provided file contains more than one document", input))
		}

		return content.Documents[0].Content[0]
	}

	// Load by parsing the actual string as YAML if it was not a file location
	return singleDoc(input)
}

func list(input string) *yamlv3.Node {
	return singleDoc(input)
}

func singleDoc(input string) *yamlv3.Node {
	docs, err := ytbx.LoadYAMLDocuments([]byte(input))
	if err != nil {
		Fail(fmt.Sprintf("Failed to parse as YAML:\n%s\n\n%v", input, err))
	}

	if len(docs) > 1 {
		Fail(fmt.Sprintf("Failed to use YAML, because it contains multiple documents:\n%s\n", input))
	}

	return docs[0].Content[0]
}

func multiDoc(input string) []*yamlv3.Node {
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

func path(path string) *ytbx.Path {
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

	return &result
}

func humanDiff(diff Diff) string {
	reporter := HumanReport{
		Report:            Report{Diffs: []Diff{diff}},
		DoNotInspectCerts: false,
		NoTableStyle:      false,
		OmitHeader:        true,
	}

	var buf bytes.Buffer
	if err := reporter.WriteReport(&buf); err != nil {
		Fail(err.Error())
	}

	return buf.String()
}

func nodify(obj interface{}) *yamlv3.Node {
	if obj == nil {
		return nil
	}

	switch tobj := obj.(type) {
	case *yamlv3.Node:
		return tobj

	case []string:
		return AsSequenceNode(tobj)

	case string:
		return &yamlv3.Node{
			Kind:  yamlv3.ScalarNode,
			Tag:   "!!str",
			Value: tobj,
		}

	case int:
		return &yamlv3.Node{
			Kind:  yamlv3.ScalarNode,
			Tag:   "!!int",
			Value: strconv.Itoa(tobj),
		}

	case float64:
		return &yamlv3.Node{
			Kind:  yamlv3.ScalarNode,
			Tag:   "!!float",
			Value: strconv.FormatFloat(tobj, 'f', -1, 64),
		}

	case bool:
		return &yamlv3.Node{
			Kind:  yamlv3.ScalarNode,
			Tag:   "!!bool",
			Value: fmt.Sprintf("%v", tobj),
		}
	}

	Fail(fmt.Sprintf("Unable to translate %v (%T) into a YAML v3 Node", obj, obj))
	return nil
}

func singleDiff(p string, change rune, from, to interface{}) Diff {
	return Diff{
		Path: path(p),
		Details: []Detail{
			{
				Kind: change,
				From: nodify(from),
				To:   nodify(to),
			},
		},
	}
}

func doubleDiff(p string, change1 rune, from1, to1 interface{}, change2 rune, from2, to2 interface{}) Diff {
	return Diff{
		Path: path(p),
		Details: []Detail{
			{
				Kind: change1,
				From: nodify(from1),
				To:   nodify(to1),
			},
			{
				Kind: change2,
				From: nodify(from2),
				To:   nodify(to2),
			},
		},
	}
}

func compare(from *yamlv3.Node, to *yamlv3.Node, compareOptions ...CompareOption) ([]Diff, error) {
	report, err := CompareInputFiles(
		ytbx.InputFile{Documents: []*yamlv3.Node{from}},
		ytbx.InputFile{Documents: []*yamlv3.Node{to}},
		compareOptions...,
	)

	if err != nil {
		return nil, err
	}

	return report.Diffs, nil
}

func BeSameDiffAs(expected Diff) types.GomegaMatcher {
	return &diffMatcher{
		expected: expected,
	}
}

type diffMatcher struct {
	expected Diff
}

func (matcher *diffMatcher) Match(actual interface{}) (success bool, err error) {
	actualDiff, ok := actual.(Diff)
	if !ok {
		return false, fmt.Errorf("BeSameDiffAs matcher expected a object of type Diff, not %T", actual)
	}

	return isSameDiff(actualDiff, matcher.expected)
}

func (matcher *diffMatcher) FailureMessage(actual interface{}) string {
	return fmt.Sprintf("Expected\n\t%s\nto be same as\n\t%s",
		spew.Sdump(actual),
		spew.Sdump(matcher.expected))
}

func (matcher *diffMatcher) NegatedFailureMessage(actual interface{}) string {
	return fmt.Sprintf("Expected\n\t%s\nnot to be same as\n\t%s",
		spew.Sdump(actual),
		spew.Sdump(matcher.expected),
	)
}

func isSameDiff(a, b Diff) (bool, error) {
	if a.Path.ToGoPatchStyle() != b.Path.ToGoPatchStyle() {
		return false, nil
	}

	if len(a.Details) != len(b.Details) {
		return false, nil
	}

	for i := range a.Details {
		if sameDetail, err := isSameDetail(a.Details[i], b.Details[i]); !sameDetail {
			return sameDetail, err
		}
	}

	return true, nil
}

func isSameDetail(a, b Detail) (bool, error) {
	if a.Kind != b.Kind {
		return false, nil
	}

	if sameNode, err := isSameNode(a.From, b.From); !sameNode {
		return sameNode, err
	}

	if sameNode, err := isSameNode(a.To, b.To); !sameNode {
		return sameNode, err
	}

	return true, nil
}

func isSameNode(a, b *yamlv3.Node) (bool, error) {
	if a == nil && b == nil {
		return true, nil
	}

	if (a == nil && b != nil) || (a != nil && b == nil) {
		return false, nil
	}

	if a.Kind != b.Kind {
		return false, nil
	}

	if a.Tag != b.Tag {
		return false, nil
	}

	if a.Value != b.Value {
		return false, nil
	}

	if len(a.Content) != len(b.Content) {
		return false, nil
	}

	for i := range a.Content {
		if same, err := isSameNode(a.Content[i], b.Content[i]); !same {
			return same, err
		}
	}

	return true, nil
}
