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
	"fmt"

	color "github.com/HeavyWombat/color"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/HeavyWombat/dyff/pkg/dyff"
)

var _ = Describe("Core/Output", func() {
	Describe("Human readable report", func() {
		Context("reporting differences", func() {
			It("should show a nice string difference", func() {
				content := singleDiff("/some/yaml/structure/string", MODIFICATION, "foobar", "Foobar")
				Expect(humanDiff(content)).To(BeEquivalentTo(`
some.yaml.structure.string
  ± value change
    - foobar
    + Foobar
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

				result := CompareDocuments(from, to)
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
					"                                     ↵\n\n"))
			})
		})
	})

	Describe("Column output", func() {
		Context("writing output nicely", func() {
			It("should show a nice table output with simple text", func() {
				stringA := `
#1
#2
#3
#4
`

				stringB := `
Mr. Foobar
Mrs. Foobar
Miss Foobar`

				stringC := `
10
200
3000
40000
500000`

				Expect(Cols("  ", 0, stringA, stringB, stringC)).To(BeEquivalentTo(`
#1  Mr. Foobar   10
#2  Mrs. Foobar  200
#3  Miss Foobar  3000
#4               40000
                 500000`))
			})

			It("should show a nice table output with colored text", func() {
				color.NoColor = false
				defer func() {
					color.NoColor = true
				}()

				stringA := fmt.Sprintf(`
%s
%s
%s
%s
`, Color("#1", color.FgGreen), Color("#2", color.FgBlue), Color("#3", color.FgRed, color.Underline), Color("#4", color.FgYellow, color.Bold, color.Italic))

				stringB := `
Mr. Foobar
Mrs. Foobar
Miss Foobar`

				stringC := `
10
200
3000
40000
500000`

				Expect(Cols("  ", 0, stringA, stringB, stringC)).To(BeEquivalentTo(fmt.Sprintf(`
%s  Mr. Foobar   10
%s  Mrs. Foobar  200
%s  Miss Foobar  3000
%s               40000
                 500000`, Color("#1", color.FgGreen), Color("#2", color.FgBlue), Color("#3", color.FgRed, color.Underline), Color("#4", color.FgYellow, color.Bold, color.Italic))))
			})
		})
	})
})
