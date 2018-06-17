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

package cmd_test

import (
	"io/ioutil"
	"os"
	"testing"

	. "github.com/HeavyWombat/dyff/internal/cmd"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "cmds suite")
}

func createTestFile(input string) string {
	file, err := ioutil.TempFile("", "some-file-name")
	Expect(err).To(BeNil())

	_, err = file.Write([]byte(input))
	Expect(err).To(BeNil())

	err = file.Close()
	Expect(err).To(BeNil())

	return file.Name()
}

var _ = Describe("flag tests", func() {
	Context("in place feature", func() {
		It("should write a YAML file in place using restructure feature", func() {
			filename := createTestFile(`---
list:
- foo: bar
  name: one
`)
			defer os.Remove(filename)

			writer := &OutputWriter{Restructure: true, OutputStyle: "yaml"}
			writer.WriteInplace(filename)

			data, err := ioutil.ReadFile(filename)
			Expect(err).To(BeNil())
			Expect(string(data)).To(BeEquivalentTo(`---
list:
- name: one
  foo: bar

`))
		})

		It("should write a JSON file in place using restructure feature", func() {
			filename := createTestFile(`{"list":[{"foo":"bar","name":"one"}]}`)
			defer os.Remove(filename)

			writer := &OutputWriter{Restructure: true, OutputStyle: "json"}
			writer.WriteInplace(filename)

			data, err := ioutil.ReadFile(filename)
			Expect(err).To(BeNil())
			Expect(string(data)).To(BeEquivalentTo(`{"list": [{"name": "one", "foo": "bar"}]}
`))
		})
	})
})
