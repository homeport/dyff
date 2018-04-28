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

package bunt_test

import (
	"fmt"

	. "github.com/HeavyWombat/dyff/pkg/bunt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Bunt tests", func() {
	Context("Helper functions", func() {
		It("should be able to break up an existing string with colors in its respective parts", func() {
			sample := fmt.Sprintf("[%s] This is a %s with various %s combined in one %s:\n(%s)",
				Colorize("Info", Green),
				Style("text", Italic),
				Style("styles", Bold),
				Colorize(`"string"`, Magenta, Bold),
				Colorize("https://github.com/HeavyWombat/dyff/pkg/bunt", CornflowerBlue, Italic, Underline))

			result, err := BreakUpStringIntoColorSegments(sample)
			Expect(err).To(BeNil())

			expected := String{Segment{Data: "["},
				Segment{Data: "Info", Attributes: []uint32{38, 2, 0, 128, 0}},
				Segment{Data: "] This is a "},
				Segment{Data: "text", Attributes: []uint32{Italic}},
				Segment{Data: " with various "},
				Segment{Data: "styles", Attributes: []uint32{Bold}},
				Segment{Data: " combined in one "},
				Segment{Data: `"string"`, Attributes: []uint32{1, 38, 2, 255, 0, 255}},
				Segment{Data: ":\n"},
				Segment{Data: "("},
				Segment{Data: "https://github.com/HeavyWombat/dyff/pkg/bunt", Attributes: []uint32{3, 4, 38, 2, 100, 149, 237}},
				Segment{Data: ")"},
			}

			Expect(result).To(BeEquivalentTo(expected))
		})
	})
})
