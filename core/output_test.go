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
				content := Diff{Path: path("/some/yaml/structure/string"),
					Kind: MODIFICATION,
					From: "foobar",
					To:   "Foobar"}

				Expect(humanDiff(content)).To(BeEquivalentTo(`some.yaml.structure.string
changed from foobar to Foobar
`))
			})
		})
	})
})
