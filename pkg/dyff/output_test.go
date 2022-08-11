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

package dyff_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/homeport/dyff/pkg/dyff"

	. "github.com/gonvenience/bunt"
)

var _ = Describe("general output tests", func() {
	Context("writing nice column output", func() {
		It("should show a nice table output with simple text", func() {
			stringA := `
#1
#2
#3
#4
`

			stringB := `
Mr. Foobar
Mrs. Foobar
Miss Foobar`

			stringC := `
10
200
3000
40000
500000`

			Expect(dyff.CreateTableStyleString("  ", 0, stringA, stringB, stringC)).To(BeEquivalentTo(`
#1  Mr. Foobar   10
#2  Mrs. Foobar  200
#3  Miss Foobar  3000
#4               40000
                 500000`))
		})

		It("should show a nice table output with colored text", func() {
			SetColorSettings(ON, ON)
			defer SetColorSettings(AUTO, AUTO)

			stringA := fmt.Sprintf(`
%s
%s
%s
%s
`, Sprintf("Lime{#1}"), Sprintf("Blue{#2}"), Sprintf("Aqua{~#3~}"), Sprintf("LemonChiffon{_*#4*_}"))

			stringB := `
Mr. Foobar
Mrs. Foobar
Miss Foobar`

			stringC := `
10
200
3000
40000
500000`

			Expect(dyff.CreateTableStyleString("  ", 0, stringA, stringB, stringC)).To(BeEquivalentTo(fmt.Sprintf(`
%s  Mr. Foobar   10
%s  Mrs. Foobar  200
%s  Miss Foobar  3000
%s               40000
                 500000`, Sprintf("Lime{#1}"), Sprintf("Blue{#2}"), Sprintf("Aqua{~#3~}"), Sprintf("LemonChiffon{_*#4*_}"))))
		})
	})
})
