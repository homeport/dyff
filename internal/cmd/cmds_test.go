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
		Context("using restructure", func() {
			Context("to write the file to STDOUT", func() {
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

			Context("to write the file in-place", func() {
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
					Expect(string(data)).To(BeEquivalentTo(`list:
  - name: one
    aaa: bbb
`))

				})
			})

			Context("incorrect usage", func() {
				It("should fail to write a YAML when in place and STDIN are used at the same time", func() {
					_, err := dyff("yaml", "--in-place", "-")
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(BeEquivalentTo("incompatible flags: cannot use in-place flag in combination with input from STDIN"))
				})
			})
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

		It("should write a plain JSON file to STDOUT using restructure feature", func() {
			filename := createTestFile(`{"list":[{"aaa":"bbb","name":"one"}]}`)
			defer os.Remove(filename)

			out, err := dyff("json", "--restructure", "--plain", filename)
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(BeEquivalentTo(`{"list": [{"name": "one", "aaa": "bbb"}]}
`))
		})

		It("should write a JSON file to STDOUT using restructure feature", func() {
			filename := createTestFile(`{"list":[{"aaa":"bbb","name":"one"}]}`)
			defer os.Remove(filename)

			out, err := dyff("json", "--restructure", filename)
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(BeEquivalentTo(`{
  "list": [
    {
      "name": "one",
      "aaa": "bbb"
    }
  ]
}
`))
		})

		It("should fail to write a JSON when in place and STDIN are used at the same time", func() {
			_, err := dyff("json", "--in-place", "-")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(BeEquivalentTo("incompatible flags: cannot use in-place flag in combination with input from STDIN"))
		})

		It("should write timestamps with proper quotes in plain mode", func() {
			out, err := dyff("json", "--plain", assets("issues", "issue-120", "buildpack.toml"))
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(Equal(`{"metadata": {"dependencies": [{"deprecation_date": "2021-08-21T00:00:00Z"}], "dependency_deprecation_dates": [{"date": "2021-08-21T13:37:00Z"}]}}
`))
		})

		It("should write timestamps with proper quotes in default mode", func() {
			out, err := dyff("json", assets("issues", "issue-120", "buildpack.toml"))
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(Equal(`{
  "metadata": {
    "dependencies": [
      {
        "deprecation_date": "2021-08-21T00:00:00Z"
      }
    ],
    "dependency_deprecation_dates": [
      {
        "date": "2021-08-21T13:37:00Z"
      }
    ]
  }
}
`))
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

		It("should create a report using a custom root in the files", func() {
			from, to := assets("examples", "from.yml"), assets("examples", "to.yml")
			expected := fmt.Sprintf(`     _        __  __
   _| |_   _ / _|/ _|  between %s, YAML root was changed to yaml.map
 / _' | | | | |_| |_       and %s, YAML root was changed to yaml.map
| (_| | |_| |  _|  _|
 \__,_|\__, |_| |_|   returned four differences
        |___/

(root level)
- six map entries removed:   + six map entries added:
  stringB: fOObAr              stringY: YAML!
  intB: 10                     intY: 147
  floatB: 2.71                 floatY: 24.0
  boolB: false                 boolY: true
  mapB:                        mapY:
    key0: B                      key0: Y
    key1: B                      key1: Y
  listB:                       listY:
  - B                          - Yo
  - B                          - Yo
  - B                          - Yo

type-change-1
  ± type change from string to int
    - string
    + 147

type-change-2
  ± type change from string to int
    - 12
    + 12

whitespaces
  ± whitespace only change
    - Strings·can··have·whitespaces.     + Strings·can··have·whitespaces.↵
                                           ↵
                                           ↵


`, from, to)

			out, err := dyff("between", from, to, "--chroot", "yaml.map")
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(BeEquivalentTo(expected))
		})

		It("should fail when change root is used with files containing multiple documents", func() {
			from, to := assets("testbed", "from.yml"), assets("testbed", "to.yml")
			_, err := dyff("between", from, to, "--chroot", "orderchanges")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(BeEquivalentTo(fmt.Sprintf("failed to change root of %s to path orderchanges: change root for an input file is only possible if there is only one document, but %s contains two documents", from, from)))
		})

		It("should fail when change root is used with files that do not have the specified path", func() {
			from, to := assets("examples", "from.yml"), assets("binary", "to.yml")
			_, err := dyff("between", from, to, "--chroot", "yaml.map")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(BeEquivalentTo(fmt.Sprintf("failed to change root of %s to path yaml.map: no key 'map' found in map, available keys: data", to)))
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

		It("should omit the dyff banner header if respective flag is set", func() {
			from := createTestFile(`{"list":[{"aaa":"bbb","name":"one"}]}`)
			defer os.Remove(from)

			to := createTestFile(`{"list":[{"aaa":"bbb","name":"two"}]}`)
			defer os.Remove(to)

			out, err := dyff("between", "--omit-header", from, to)
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(BeEquivalentTo(`
list
  - one list entry removed:     + one list entry added:
    - name: one                   - name: two
      aaa: bbb                      aaa: bbb

`))
		})

		It("should ignore order changes if respective flag is set", func() {
			from := createTestFile(`{"list":[{"name":"one"},{"name":"two"},{"name":"three"}]}`)
			defer os.Remove(from)

			to := createTestFile(`{"list":[{"name":"one"},{"name":"three"},{"name":"two"}]}`)
			defer os.Remove(to)

			out, err := dyff("between", "--omit-header", "--ignore-order-changes", from, to)
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(BeEquivalentTo("\n"))
		})

		It("should not panic when timestamps need to reported", func() {
			out, err := dyff("between", "--omit-header", "../../assets/issues/issue-111/from.yml", "../../assets/issues/issue-111/to.yml")
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(BeEquivalentTo(`
AWSTemplateFormatVersion
  ± type change from timestamp to string
    - 2010-09-09
    + 2010-09-09

`))
		})

		It("should not try to evaluate variables in the user-provided strings", func() {
			out, err := dyff("between", "--omit-header", assets("issues", "issue-132", "from.yml"), assets("issues", "issue-132", "to.yml"))
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(BeEquivalentTo(`
example_one
  ± value change
    - %{one}
    + one

example_two
  ± value change
    - two
    + %{two}

`))
		})
	})

	Context("last-applied command", func() {
		It("should create the default report when there are no flags specified", func() {
			kubeYAML := createTestFile(`---
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      { "metadata": { "annotations": {} }, "yaml": { "foo": "bat" } }
yaml:
  foo: bar
`)
			defer os.Remove(kubeYAML)

			out, err := dyff("last-applied", "--omit-header", kubeYAML)
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(BeEquivalentTo(`
yaml.foo
  ± value change
    - bat
    + bar

`))
		})

		It("should fail on an input file with multiple documents", func() {
			kubeYAML := createTestFile(`---
foo: bar
--
foo: bar
`)
			defer os.Remove(kubeYAML)

			_, err := dyff("last-applied", kubeYAML)
			Expect(err).To(HaveOccurred())
		})

		It("should fail on an input file when the last applied configuration is not set", func() {
			kubeYAML := createTestFile(`foo: bar`)
			defer os.Remove(kubeYAML)

			_, err := dyff("last-applied", kubeYAML)
			Expect(err).To(HaveOccurred())
		})
	})
})
