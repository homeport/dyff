package core_test

import (
	. "github.com/HeavyWombat/dyff/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("JSON", func() {
	Describe("Getting YAML input", func() {
		Context("Processing valid YAML input", func() {
			It("should convert YAML to JSON", func() {
				content := getYamlFromString(`---
name: foobar
list:
- A
- B
- C
`)

				result, err := ToJSONString(content)
				Expect(err).To(BeNil())

				Expect(result).To(Equal(`{"name": "foobar", "list": ["A", "B", "C"]}`))
			})
		})
	})
})
