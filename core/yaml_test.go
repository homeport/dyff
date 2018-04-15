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

package core_test

import (
	. "github.com/HeavyWombat/dyff/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Core/YAML", func() {
	Describe("Getting JSON input", func() {
		Context("Processing valid JSON input", func() {
			It("should convert JSON to YAML", func() {
				content := yml(`{ "name": "foobar", "list": [A, B, C] }`)

				result, err := ToYAMLString(content)
				Expect(err).To(BeNil())

				Expect(result).To(Equal(`---
name: foobar
list:
- A
- B
- C

`))
			})

			It("should preserve the order inside the structure", func() {
				content := yml(`{ "list": [C, B, A], "name": "foobar" }`)

				result, err := ToYAMLString(content)
				Expect(err).To(BeNil())

				Expect(result).To(Equal(`---
list:
- C
- B
- A
name: foobar

`))
			})
		})
	})
})
