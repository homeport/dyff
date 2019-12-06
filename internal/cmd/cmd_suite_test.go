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
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"testing"

	. "github.com/homeport/dyff/internal/cmd"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "dyff commands package suite")
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

func captureStdout(f func() error) (string, error) {
	r, w, err := os.Pipe()
	Expect(err).ToNot(HaveOccurred())

	tmp := os.Stdout
	defer func() {
		os.Stdout = tmp
	}()

	os.Stdout = w
	err = f()
	w.Close()

	var buf bytes.Buffer
	io.Copy(&buf, r)

	return buf.String(), err
}

func dyff(args ...string) (out string, err error) {
	return captureStdout(func() error {
		ResetSettings()
		os.Args = append([]string{"dyff"}, args...)
		return Execute()
	})
}
