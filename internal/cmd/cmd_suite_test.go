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
	"os"
	"path/filepath"
	"testing"

	. "github.com/gonvenience/bunt"
	"github.com/gonvenience/text"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/homeport/dyff/internal/cmd"
)

func TestCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "dyff commands package suite")
}

func createTestFile(input string) string {
	return createTestFileInDir("", input)
}

func createTestFileInDir(dir string, input string) string {
	file, err := os.CreateTemp(dir, "some-file-name")
	Expect(err).To(BeNil())

	_, err = file.Write([]byte(input))
	Expect(err).To(BeNil())

	err = file.Close()
	Expect(err).To(BeNil())

	return file.Name()
}

func createTestDirectory() string {
	var path = filepath.Join(os.TempDir(), text.RandomString(8))

	err := os.MkdirAll(path, os.FileMode(0755))
	Expect(err).ToNot(HaveOccurred())

	return path
}

func assets(pathElement ...string) string {
	targetPath := filepath.Join(append(
		[]string{"..", "..", "assets"},
		pathElement...,
	)...)

	abs, err := filepath.Abs(targetPath)
	if err != nil {
		return targetPath
	}

	return abs
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
	if _, copyErr := io.Copy(&buf, r); copyErr != nil {
		return "", copyErr
	}

	return buf.String(), err
}

func dyff(args ...string) (out string, err error) {
	SetColorSettings(OFF, OFF)
	defer SetColorSettings(AUTO, AUTO)

	return captureStdout(func() error {
		ResetSettings()
		os.Args = append([]string{"dyff"}, args...)
		return Execute()
	})
}
