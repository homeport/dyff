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
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/gonvenience/bunt"

	"github.com/gonvenience/ytbx"

	"github.com/homeport/dyff/pkg/dyff"
)

var _ = Describe("human readable report", func() {
	Context("reporting differences", func() {
		BeforeEach(func() {
			SetColorSettings(OFF, OFF)
		})

		AfterEach(func() {
			SetColorSettings(AUTO, AUTO)
		})

		It("should show a nice string difference", func() {
			content := singleDiff("/some/yaml/structure/string", dyff.MODIFICATION, "fOObar?", "Foobar!")
			Expect(humanDiff(content)).To(BeEquivalentTo(`
some.yaml.structure.string
  ± value change
    - fOObar?
    + Foobar!

`))
		})

		It("should show a nice integer difference", func() {
			content := singleDiff("/some/yaml/structure/int", dyff.MODIFICATION, 12, 147)
			Expect(humanDiff(content)).To(BeEquivalentTo(`
some.yaml.structure.int
  ± value change
    - 12
    + 147

`))
		})

		It("should show a type difference", func() {
			content := singleDiff("/some/yaml/structure/test", dyff.MODIFICATION, 12, 12.0)
			Expect(humanDiff(content)).To(BeEquivalentTo(`
some.yaml.structure.test
  ± type change from int to float
    - 12
    + 12

`))
		})

		It("should return that strings are identical if only the usage of whitespaces differs", func() {
			from := yml(`---
input: |+
  This is a text with
  newlines and stuff
  to show case whitespace
  issues.
`)

			to := yml(`---
input: |+
  This is a text with
  newlines and stuff
  to show case whitespace
  issues.

`)

			result, err := compare(from, to)
			Expect(err).To(BeNil())
			Expect(result).NotTo(BeNil())
			Expect(len(result)).To(BeEquivalentTo(1))
			Expect(result[0]).To(BeSameDiffAs(singleDiff("/input",
				dyff.MODIFICATION,
				"This is a text with"+"\n"+"newlines and stuff"+"\n"+"to show case whitespace"+"\n"+"issues.\n",
				"This is a text with"+"\n"+"newlines and stuff"+"\n"+"to show case whitespace"+"\n"+"issues.\n\n")))

			Expect(humanDiff(result[0])).To(BeEquivalentTo("\ninput\n  ± whitespace only change\n    - This·is·a·text·with↵         + This·is·a·text·with↵\n" +
				"      newlines·and·stuff↵            newlines·and·stuff↵\n" +
				"      to·show·case·whitespace↵       to·show·case·whitespace↵\n" +
				"      issues.↵                       issues.↵\n" +
				"                                     ↵\n\n\n"))
		})

		It("should show a binary data difference in hex dump style", func() {
			compareAgainstExpectedHuman("../../assets/binary/from.yml",
				"../../assets/binary/to.yml",
				"../../assets/binary/dyff.expected",
				false,
			)
		})

		It("should show the testbed results as expected", func() {
			compareAgainstExpectedHuman("../../assets/testbed/from.yml",
				"../../assets/testbed/to.yml",
				"../../assets/testbed/expected-dyff-spruce.human",
				false,
			)

			compareAgainstExpectedHuman("../../assets/testbed/from.yml",
				"../../assets/testbed/to.yml",
				"../../assets/testbed/expected-dyff-gopatch.human",
				true,
			)
		})

		It("should report each file level change separately for better readability", func() {
			compareAgainstExpectedHuman(
				assets("kubernetes/multi-docs-file-level/from.yaml"),
				assets("kubernetes/multi-docs-file-level/to.yaml"),
				assets("kubernetes/multi-docs-file-level/expected-dyff.human"),
				false,
			)
		})

		It("should detect renames for kubernetes documents", func() {
			compareAgainstExpectedHuman(
				assets("kubernetes/rename/from.yaml"),
				assets("kubernetes/rename/to.yaml"),
				assets("kubernetes/rename/expected-dyff.human"),
				false,
			)
		})
	})

	Context("nicely colored human readable differences", func() {
		BeforeEach(func() {
			SetColorSettings(ON, ON)
		})

		AfterEach(func() {
			SetColorSettings(AUTO, AUTO)
		})

		It("should show nicely colored difference output", func() {
			compareAgainstExpectedHuman(
				"../../assets/colors/from.yml",
				"../../assets/colors/to.yml",
				"../../assets/colors/dyff.expected",
				false,
			)
		})
	})

	Context("human path rendering", func() {
		BeforeEach(func() {
			SetColorSettings(ON, ON)
		})

		AfterEach(func() {
			SetColorSettings(AUTO, AUTO)
		})

		It("should render path with underscores correctly (https://github.com/homeport/dyff/issues/33)", func() {
			// Please note: The actual error is in the gonvenience package, this test
			// case exists to verify the issue from with dyff.

			path, err := ytbx.ParseGoPatchStylePathString("/variables/name=ROUTER_TLS_PEM/options")
			Expect(err).ToNot(HaveOccurred())

			content := singleDiff(path.String(), dyff.MODIFICATION, 12, 12.0)
			actual := humanDiff(content)

			Expect(fmt.Sprintf("%#v", RemoveAllEscapeSequences(actual))).To(
				BeEquivalentTo(fmt.Sprintf("%#v", `
variables.ROUTER_TLS_PEM.options
  ± type change from int to float
    - 12
    + 12

`)))
		})
	})

	Context("human friendly multiline text differences", func() {
		BeforeEach(func() {
			SetColorSettings(ON, ON)
		})

		AfterEach(func() {
			SetColorSettings(AUTO, AUTO)
		})

		It("should use human friendly compact diff of multiline text differences", func() {
			compareAgainstExpectedHuman(
				assets("kubernetes/configmaps/from.yml"),
				assets("kubernetes/configmaps/to.yml"),
				assets("kubernetes/configmaps/expected-dyff-spruce.human"),
				false,
			)
		})

		It("should use human friendly compact diff of multiline text differences for complex files", func() {
			compareAgainstExpectedHuman(
				assets("multiline/from.yml"),
				assets("multiline/to.yml"),
				assets("multiline/expected-dyff-spruce.human"),
				false,
			)
		})
	})

	Context("reported output issues (without colors)", func() {
		BeforeEach(func() {
			SetColorSettings(OFF, OFF)
		})

		AfterEach(func() {
			SetColorSettings(AUTO, AUTO)
		})

		It("should use correct indentions for all kind of changes", func() {
			compareAgainstExpectedHuman("../../assets/issues/issue-89/from.yml",
				"../../assets/issues/issue-89/to.yml",
				"../../assets/issues/issue-89/expected-dyff-spruce.human",
				false,
			)
		})
	})

	Context("reported output issues (with colors)", func() {
		BeforeEach(func() {
			SetColorSettings(ON, ON)
		})

		AfterEach(func() {
			SetColorSettings(AUTO, AUTO)
		})

		It("should not parse markdown in multiline text diff output", func() {
			compareAgainstExpectedHuman(
				assets("issues", "issue-225", "from.yml"),
				assets("issues", "issue-225", "to.yml"),
				assets("issues", "issue-225", "expected-dyff-spruce.human"),
				false,
			)
		})
	})
})
