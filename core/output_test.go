package core_test

import (
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
			It("should show a nice table output", func() {
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

				Expect(Cols(stringA, stringB, stringC)).To(BeEquivalentTo(`
#1  Mr. Foobar   10
#2  Mrs. Foobar  200
#3  Miss Foobar  3000
#4               40000
                 500000
`))
			})
		})
	})
})
