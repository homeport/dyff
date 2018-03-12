package main_test

import (
	. "github.com/HeavyWombat/dyff/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	yaml "gopkg.in/yaml.v2"
)

func getYamlFromString(input string) (yaml.MapSlice, error) {
	content := yaml.MapSlice{}
	err := yaml.UnmarshalStrict([]byte(input), &content)
	if err != nil {
		return nil, err
	}

	return content, nil
}

var _ = Describe("YAML", func() {
	Describe("Getting JSON input", func() {
		Context("Processing valid JSON input", func() {
			It("should convert JSON to YAML", func() {
				content, err := getYamlFromString(`{ "name": "foobar", "list": [A, B, C] }`)
				Expect(err).To(BeNil())

				var result string
				result, err = ToYAMLString(content)
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
				content, err := getYamlFromString(`{ "list": [C, B, A], "name": "foobar" }`)
				Expect(err).To(BeNil())

				var result string
				result, err = ToYAMLString(content)
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
