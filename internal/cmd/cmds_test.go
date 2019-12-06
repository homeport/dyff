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
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("command line tool flags", func() {
	Context("version command", func() {
		It("should print the version", func() {
			out, err := dyff("version")
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(BeEquivalentTo("dyff version (development)\n"))
		})
	})

	Context("yaml command", func() {
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

			data, err := ioutil.ReadFile(filename)
			Expect(err).To(BeNil())
			Expect(string(data)).To(BeEquivalentTo(`---
list:
- name: one
  aaa: bbb

`))
		})

		It("should write a YAML file to STDOUT using restructure feature", func() {
			filename := createTestFile(`---
list:
- aaa: bbb
  name: one
`)
			defer os.Remove(filename)

			out, err := dyff("yaml", "--restructure", filename)
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(BeEquivalentTo(`---
list:
- name: one
  aaa: bbb

`))
		})
	})

	Context("json command", func() {
		It("should write a JSON file in place using restructure feature", func() {
			filename := createTestFile(`{"list":[{"aaa":"bbb","name":"one"}]}`)
			defer os.Remove(filename)

			out, err := dyff("json", "--restructure", "--in-place", filename)
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(BeEmpty())

			data, err := ioutil.ReadFile(filename)
			Expect(err).To(BeNil())
			Expect(string(data)).To(BeEquivalentTo(`{"list": [{"name": "one", "aaa": "bbb"}]}
`))
		})

		It("should write a JSON file to STDOUT using restructure feature", func() {
			filename := createTestFile(`{"list":[{"aaa":"bbb","name":"one"}]}`)
			defer os.Remove(filename)

			out, err := dyff("json", "--restructure", "--plain", filename)
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(BeEquivalentTo(`{"list": [{"name": "one", "aaa": "bbb"}]}
`))
		})
	})
})
