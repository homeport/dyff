package core_test

import (
	"fmt"

	color "github.com/HeavyWombat/color"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/HeavyWombat/dyff/core"
)

var _ = Describe("Core/Output", func() {
	Describe("Human readable report", func() {
		Context("reporting differences", func() {
			It("should show a nice string difference", func() {
				content := singleDiff("/some/yaml/structure/string", MODIFICATION, "foobar", "Foobar")
				Expect(humanDiff(content)).To(BeEquivalentTo(`some.yaml.structure.string
  changed value
   - foobar
   + Foobar

`))
			})
		})
	})

	Describe("Column output", func() {
		Context("writing output nicely", func() {
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

				Expect(Cols("  ", 0, stringA, stringB, stringC)).To(BeEquivalentTo(`
#1  Mr. Foobar   10
#2  Mrs. Foobar  200
#3  Miss Foobar  3000
#4               40000
                 500000
`))
			})

			It("should show a nice table output with colored text", func() {
				color.NoColor = false
				defer func() {
					color.NoColor = true
				}()

				stringA := fmt.Sprintf(`
%s
%s
%s
%s
`, Color("#1", color.FgGreen), Color("#2", color.FgBlue), Color("#3", color.FgRed, color.Underline), Color("#4", color.FgYellow, color.Bold, color.Italic))

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

				Expect(Cols("  ", 0, stringA, stringB, stringC)).To(BeEquivalentTo(fmt.Sprintf(`
%s  Mr. Foobar   10
%s  Mrs. Foobar  200
%s  Miss Foobar  3000
%s               40000
                 500000
`, Color("#1", color.FgGreen), Color("#2", color.FgBlue), Color("#3", color.FgRed, color.Underline), Color("#4", color.FgYellow, color.Bold, color.Italic))))
			})
		})
	})
})
