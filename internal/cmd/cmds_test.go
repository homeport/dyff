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
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gonvenience/bunt"
	"github.com/gonvenience/term"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("command line tool flags", func() {
	BeforeEach(func() {
		bunt.ColorSetting = bunt.OFF
		bunt.TrueColorSetting = bunt.OFF
		term.FixedTerminalWidth = 250
		term.FixedTerminalHeight = 40
	})

	AfterEach(func() {
		bunt.ColorSetting = bunt.AUTO
		bunt.TrueColorSetting = bunt.AUTO
		term.FixedTerminalWidth = -1
		term.FixedTerminalHeight = -1
	})

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

		It("should fail to write a YAML when in place and STDIN are used at the same time", func() {
			_, err := dyff("yaml", "--in-place", "-")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(BeEquivalentTo("incompatible flags: cannot use in-place flag in combination with input from STDIN"))
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

		It("should fail to write a JSON when in place and STDIN are used at the same time", func() {
			_, err := dyff("json", "--in-place", "-")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(BeEquivalentTo("incompatible flags: cannot use in-place flag in combination with input from STDIN"))
		})
	})

	Context("between command", func() {
		It("should create the default report when there are no flags specified", func() {
			from := createTestFile(`{"list":[{"aaa":"bbb","name":"one"}]}`)
			defer os.Remove(from)

			to := createTestFile(`{"list":[{"aaa":"bbb","name":"two"}]}`)
			defer os.Remove(to)

			out, err := dyff("between", from, to)
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(BeEquivalentTo(fmt.Sprintf(`     _        __  __
   _| |_   _ / _|/ _|  between %s
 / _' | | | | |_| |_       and %s
| (_| | |_| |  _|  _|
 \__,_|\__, |_| |_|   returned one difference
        |___/

list
  - one list entry removed:     + one list entry added:
    - name: one                   - name: two
      aaa: bbb                      aaa: bbb

`, from, to)))
		})

		It("should create the same default report when swap flag is used", func() {
			from := createTestFile(`{"list":[{"aaa":"bbb","name":"one"}]}`)
			defer os.Remove(from)

			to := createTestFile(`{"list":[{"aaa":"bbb","name":"two"}]}`)
			defer os.Remove(to)

			out, err := dyff("between", "--swap", to, from)
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(BeEquivalentTo(fmt.Sprintf(`     _        __  __
   _| |_   _ / _|/ _|  between %s
 / _' | | | | |_| |_       and %s
| (_| | |_| |  _|  _|
 \__,_|\__, |_| |_|   returned one difference
        |___/

list
  - one list entry removed:     + one list entry added:
    - name: one                   - name: two
      aaa: bbb                      aaa: bbb

`, from, to)))
		})

		It("should create the oneline report", func() {
			from := createTestFile(`{"list":[{"aaa":"bbb","name":"one"}]}`)
			defer os.Remove(from)

			to := createTestFile(`{"list":[{"aaa":"bbb","name":"two"}]}`)
			defer os.Remove(to)

			out, err := dyff("between", "--output=brief", from, to)
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(BeEquivalentTo(fmt.Sprintf("one change detected between %s and %s\n\n", from, to)))
		})

		It("should return an exit code with the number of differences if respective flag is used", func() {
			from := createTestFile(`{"list":[{"aaa":"bbb","name":"one"}]}`)
			defer os.Remove(from)

			to := createTestFile(`{"list":[{"aaa":"bbb","name":"two"}]}`)
			defer os.Remove(to)

			out, err := dyff("between", "--output=brief", "--set-exit-status", from, to)
			Expect(err).To(HaveOccurred())
			Expect(out).To(BeEquivalentTo(fmt.Sprintf("one change detected between %s and %s\n\n", from, to)))
		})

		It("should fail when input files cannot be read", func() {
			_, err := dyff("between", "/does/not/exist/from.yml", "/does/not/exist/to.yml")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to load input files: unable to load data from /does/not/exist/from.yml"))
		})

		It("should fail when an unsupported output style is defined", func() {
			_, err := dyff("between", "--output", "unknown", "/dev/null", "/dev/null")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("unknown output style unknown"))
		})
	})
})
