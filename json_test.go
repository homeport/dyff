package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/HeavyWombat/dyff/core"
)

var _ = Describe("JSON", func() {
	Describe("Getting YAML input", func() {
		Context("Processing valid YAML input", func() {
			It("should convert YAML to JSON", func() {
				content, err := getYamlFromString(`---
name: foobar
list:
- A
- B
- C
`)
				Expect(err).To(BeNil())

				var result string
				result, err = ToJSONString(content)
				Expect(err).To(BeNil())

				Expect(result).To(Equal(`{"name": "foobar", "list": ["A", "B", "C"]}`))
			})
		})
	})
})
