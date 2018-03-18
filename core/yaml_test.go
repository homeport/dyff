package core_test

import (
	. "github.com/HeavyWombat/dyff/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("YAML", func() {
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
