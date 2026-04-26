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

package cmd_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gonvenience/term"
)

var _ = Describe("using dyff yaml", func() {
	BeforeEach(func() {
		term.FixedTerminalWidth = 250
		term.FixedTerminalHeight = 40
	})

	AfterEach(func() {
		term.FixedTerminalWidth = -1
		term.FixedTerminalHeight = -1
	})

	Context("streaming to StdOut", func() {
		It("should write single document without document start marker in default mode", func() {
			filename := createTestFile(`{"foo": "bar"}`)
			defer os.Remove(filename)

			out, err := dyff("yaml", filename)
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(BeEquivalentTo(`foo: bar

`))
		})

		It("should write single document without document start marker in plain mode", func() {
			filename := createTestFile(`{"foo": "bar"}`)
			defer os.Remove(filename)

			out, err := dyff("yaml", "--plain", filename)
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(BeEquivalentTo(`foo: bar
`))
		})

		It("should write multi document with document start marker in default mode", func() {
			out, err := dyff("yaml", assets("multidocs/simple/file.yaml"))
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(BeEquivalentTo(`---
foo: bar

---
bar: foo

`))
		})

		It("should write multi document with document start marker in plain mode", func() {
			out, err := dyff("yaml", "--plain", assets("multidocs/simple/file.yaml"))
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(BeEquivalentTo(`---
foo: bar
---
bar: foo
`))
		})

		It("should quote all string that have special meaning in YAML", func() {
			filename := createTestFile(`{"a": ",", "foo": {"bar": "*", "dash": "-"}}`)
			defer os.Remove(filename)

			out, err := dyff("yaml", filename)
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(BeEquivalentTo(`a: ","
foo:
  bar: "*"
  dash: "-"

`))
		})

		It("should restructure (reorder) fields", func() {
			filename := createTestFile(`---
list:
- aaa: bbb
  name: one
`)
			defer os.Remove(filename)

			out, err := dyff("yaml", "--restructure", filename)
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(BeEquivalentTo(`list:
- name: one
  aaa: bbb

`))
		})

		It("should restructure (reorder) fields in multi-document YAML", func() {
			out, err := dyff("yaml", "--plain", "--restructure", assets("issues", "issue-133", "input.yml"))
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(BeEquivalentTo(`---
name: one
bar: foo
foo: bar
---
name: two
Foo: Bar
Bar: Foo
---
name: three
foobar: foobar
`))
		})
	})

	Context("writing in-place", func() {
		It("should write a YAML file in place using restructure feature", func() {
			filename := createTestFile(`---
list:
- aaa: bbb
  name: one
`)
			defer os.Remove(filename)

			out, err := dyff("yaml", "--restructure", "--in-place", filename)
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(BeEmpty())

			data, err := os.ReadFile(filename)
			Expect(err).To(BeNil())
			Expect(string(data)).To(BeEquivalentTo(`list:
  - name: one
    aaa: bbb
`))

		})
	})

	It("should fail to write a YAML when in place and STDIN are used at the same time", func() {
		_, err := dyff("yaml", "--in-place", "-")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("incompatible flags: cannot use in-place flag in combination with input"))
	})
})
