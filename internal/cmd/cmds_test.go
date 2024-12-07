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
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/homeport/dyff/internal/cmd"

	"github.com/gonvenience/term"
)

var _ = Describe("command line tool flags", func() {
	BeforeEach(func() {
		term.FixedTerminalWidth = 250
		term.FixedTerminalHeight = 40
	})

	AfterEach(func() {
		term.FixedTerminalWidth = -1
		term.FixedTerminalHeight = -1
	})

	Context("version command", func() {
		It("should print the version", func() {
			out, err := dyff("version")
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(ContainSubstring("version (development)"))
		})
	})

	Context("yaml command", func() {
		Context("creating yaml output", func() {
			It("should not create YAML output that is not valid", func() {
				filename := createTestFile(`{"a": ",", "foo": {"bar": "*", "dash": "-"}}`)
				defer os.Remove(filename)

				out, err := dyff("yaml", filename)
				Expect(err).ToNot(HaveOccurred())
				Expect(out).To(BeEquivalentTo(`---
a: ","
foo:
  bar: "*"
  dash: "-"

`))
			})
		})

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

				It("should write a YAML file with multiple documents to STDOUT using restructure feature", func() {
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

					data, err := os.ReadFile(filename)
					Expect(err).To(BeNil())
					Expect(string(data)).To(BeEquivalentTo(`---
list:
  - name: one
    aaa: bbb
`))

				})
			})

			Context("incorrect usage", func() {
				It("should fail to write a YAML when in place and STDIN are used at the same time", func() {
					_, err := dyff("yaml", "--in-place", "-")
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(BeEquivalentTo("incompatible flags: cannot use in-place flag in combination with input from stdin"))
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

			data, err := os.ReadFile(filename)
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

		It("should print various timestamp formats and strings that look like timestamps", func() {
			out, err := dyff("json", assets("issues", "issue-217", "datestring.yml"))
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(Equal(`{
  "Datestring": "2033-12-20T00:00:00Z",
  "ThirteenthMonth": "2033-13-20",
  "FortyDays": "2033-13-40",
  "TheYear9999": "9999-11-20T00:00:00Z",
  "OneShortDatestring": "999-99-99",
  "ExtDatestring": "2021-01-01-04-05",
  "DatestringFake": "9999-99-99",
  "DatestringNonHyphenated": 99999999,
  "DatestringOneHyphen": "9999-9999",
  "DatestringSlashes": "2022/01/01"
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

			out, err := dyff("between", "--output=brief", "--set-exit-code", from, to)
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

		It("should be understandingly when the arguments are in an incorrect order when executed by kubectl diff", func() {
			from := createTestDirectory()
			defer os.RemoveAll(from)

			to := createTestDirectory()
			defer os.RemoveAll(to)

			createTestFileInDir(from, `{"list":[{"name":"one", "version":"v1"}]}`)
			createTestFileInDir(to, `{"list":[{"name":"two", "version":"v2"}]}`)

			_, err := dyff(from, to, "between", "--omit-header")
			Expect(err).To(HaveOccurred())

			// Usually, the environment variable would be like `dyff between`,
			// but the binary name during testing is `cmd.test` and therefore
			// the variable needs to adjusted for the internal program logic to
			// correctly accept this as the kubectl context it looks for.
			var tmp = os.Getenv("KUBECTL_EXTERNAL_DIFF")
			os.Setenv("KUBECTL_EXTERNAL_DIFF", "cmd.test between --omit-header")
			defer os.Setenv("KUBECTL_EXTERNAL_DIFF", tmp)

			_, err = dyff(from, to, "between", "--omit-header")
			Expect(err).ToNot(HaveOccurred())
		})

		It("should create exit code zero if there are no changes", func() {
			from := createTestFile(`{"foo": "bar"}`)
			defer os.Remove(from)

			to := createTestFile(`{"foo": "bar"}`)
			defer os.Remove(to)

			_, err := dyff("between", "--set-exit-code", from, to)
			Expect(err).To(HaveOccurred())

			exitCode, ok := err.(ExitCode)
			Expect(ok).To(BeTrue())
			Expect(exitCode.Value).To(Equal(0))
		})

		It("should create exit code one if there are changes", func() {
			from := createTestFile(`{"foo": "bar"}`)
			defer os.Remove(from)

			to := createTestFile(`{"foo": "BAR"}`)
			defer os.Remove(to)

			_, err := dyff("between", "--set-exit-code", from, to)
			Expect(err).To(HaveOccurred())

			exitCode, ok := err.(ExitCode)
			Expect(ok).To(BeTrue())
			Expect(exitCode.Value).To(Equal(1))
		})

		It("should fail with an exit code other than zero or one in case of an error", func() {
			_, err := dyff("between", "--set-exit-code", "from", "to")
			Expect(err).To(HaveOccurred())

			exitCode, ok := err.(ExitCode)
			Expect(ok).To(BeTrue())
			Expect(exitCode.Value).To(Equal(255))
		})

		It("should accept a list of paths and filter the report based on these", func() {
			expected := `
yaml.map.type-change-1
  ± type change from string to int
    - string
    + 147

yaml.map.whitespaces
  ± whitespace only change
    - Strings·can··have·whitespaces.     + Strings·can··have·whitespaces.↵
                                           ↵
                                           ↵


`

			By("using GoPatch style path", func() {
				out, err := dyff("between", "--omit-header", "--filter", "/yaml/map/whitespaces,/yaml/map/type-change-1", assets("examples", "from.yml"), assets("examples", "to.yml"))
				Expect(err).ToNot(HaveOccurred())
				Expect(out).To(BeEquivalentTo(expected))
			})

			By("using DotStyle paths", func() {
				out, err := dyff("between", "--omit-header", "--filter", "yaml.map.whitespaces,yaml.map.type-change-1", assets("examples", "from.yml"), assets("examples", "to.yml"))
				Expect(err).ToNot(HaveOccurred())
				Expect(out).To(BeEquivalentTo(expected))
			})
		})

		It("should show a report when filtering and a document has been removed between inputs", func() {
			expected := `
spec.replicas  (apps/v1/Deployment/test)
  ± value change
    - 2
    + 3

`
			By("using the --filter flag", func() {
				out, err := dyff("between", "--omit-header", "--filter", "/spec/replicas", assets("issues", "issue-232", "from.yml"), assets("issues", "issue-232", "to.yml"))
				Expect(err).ToNot(HaveOccurred())
				Expect(out).To(BeEquivalentTo(expected))
			})
		})

		It("should properly print multi-line strings (https://github.com/homeport/dyff/issues/180)", func() {
			out, err := dyff("between", "--omit-header", assets("issues", "issue-180", "old.yml"), assets("issues", "issue-180", "new.yml"))
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(BeEquivalentTo(`
(root level)
+ four map entries added:
  kind: ConfigMap
  apiVersion: v1
  metadata:
    name: atlantis-repo-config
    namespace: default
    labels:
      app: atlantis
      chart: atlantis-3.14.0
      heritage: Tiller
      release: default
  data:
    repos.yaml: |
      repos:
      - apply_requirements:
        - approved
        - mergeable
        id: /.*/

`))
		})

		It("should accept additional identifier candidates for named-entry list processing", func() {
			out, err := dyff("between", "--additional-identifier=branch", "--ignore-order-changes", "--omit-header", assets("issues", "issue-243", "from.yml"), assets("issues", "issue-243", "to.yml"))
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(BeEquivalentTo("\n"))
		})

		It("should ignore the 'whitespace only' changes", func() {
			out, err := dyff("between",
				"--omit-header",
				"--ignore-whitespace-changes",
				assets("issues", "issue-222", "from.yml"),
				assets("issues", "issue-222", "to.yml"))

			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(BeEquivalentTo("\n"))
		})

		It("should ignore the changes in values", func() {
			expected := `
(file level)
  - one document removed:
    ---
    apiVersion: v1
    kind: Namespace
    metadata:
      name: test

`
			By("using the --ignore-value-changes", func() {
				out, err := dyff("between", "--omit-header", "--ignore-value-changes", assets("issues", "issue-232", "from.yml"), assets("issues", "issue-232", "to.yml"))
				Expect(err).ToNot(HaveOccurred())
				Expect(out).To(BeEquivalentTo(expected))
			})

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
