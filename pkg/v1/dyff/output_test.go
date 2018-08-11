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

package dyff_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/HeavyWombat/dyff/pkg/v1/bunt"
	. "github.com/HeavyWombat/dyff/pkg/v1/dyff"
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

			Expect(CreateTableStyleString("  ", 0, stringA, stringB, stringC)).To(BeEquivalentTo(`
#1  Mr. Foobar   10
#2  Mrs. Foobar  200
#3  Miss Foobar  3000
#4               40000
                 500000`))
		})

		It("should show a nice table output with colored text", func() {
			ColorSetting = ON
			defer func() {
				ColorSetting = OFF
			}()

			stringA := fmt.Sprintf(`
%s
%s
%s
%s
`, Colorize("#1", Lime), Colorize("#2", Blue), Colorize("#3", Aqua, Underline), Colorize("#4", LemonChiffon, Bold, Italic))

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

			Expect(CreateTableStyleString("  ", 0, stringA, stringB, stringC)).To(BeEquivalentTo(fmt.Sprintf(`
%s  Mr. Foobar   10
%s  Mrs. Foobar  200
%s  Miss Foobar  3000
%s               40000
                 500000`, Colorize("#1", Lime), Colorize("#2", Blue), Colorize("#3", Aqua, Underline), Colorize("#4", LemonChiffon, Bold, Italic))))
		})
	})
})
