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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/HeavyWombat/dyff/pkg/v1/dyff"
)

var _ = Describe("human readable report", func() {
	Context("reporting differences", func() {
		It("should show a nice string difference", func() {
			content := singleDiff("/some/yaml/structure/string", MODIFICATION, "fOObar?", "Foobar!")
			Expect(humanDiff(content)).To(BeEquivalentTo(`
some.yaml.structure.string
  ± value change
    - fOObar?
    + Foobar!

`))
		})

		It("should show a nice integer difference", func() {
			content := singleDiff("/some/yaml/structure/int", MODIFICATION, 12, 147)
			Expect(humanDiff(content)).To(BeEquivalentTo(`
some.yaml.structure.int
  ± value change
    - 12
    + 147

`))
		})

		It("should show a type difference", func() {
			content := singleDiff("/some/yaml/structure/test", MODIFICATION, 12, 12.0)
			Expect(humanDiff(content)).To(BeEquivalentTo(`
some.yaml.structure.test
  ± type change from int to float64
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
			Expect(result[0]).To(BeEquivalentTo(singleDiff("/input",
				MODIFICATION,
				"This is a text with"+"\n"+"newlines and stuff"+"\n"+"to show case whitespace"+"\n"+"issues.\n",
				"This is a text with"+"\n"+"newlines and stuff"+"\n"+"to show case whitespace"+"\n"+"issues.\n\n")))

			Expect(humanDiff(result[0])).To(BeEquivalentTo("\ninput\n  ± whitespace only change\n    - This·is·a·text·with↵         + This·is·a·text·with↵\n" +
				"      newlines·and·stuff↵            newlines·and·stuff↵\n" +
				"      to·show·case·whitespace↵       to·show·case·whitespace↵\n" +
				"      issues.↵                       issues.↵\n" +
				"                                     ↵\n\n\n"))
		})

		It("should show a binary data difference in hex dump style", func() {
			compareAgainstExpected("../../../assets/binary/from.yml",
				"../../../assets/binary/to.yml",
				"../../../assets/binary/dyff.expected")
		})

		It("should show the testbed results as expected", func() {
			compareAgainstExpected("../../../assets/testbed/from.yml",
				"../../../assets/testbed/to.yml",
				"../../../assets/testbed/expected-dyff-spruce.human")

			UseGoPatchPaths = true
			compareAgainstExpected("../../../assets/testbed/from.yml",
				"../../../assets/testbed/to.yml",
				"../../../assets/testbed/expected-dyff-gopatch.human")
		})
	})
})
