package core_test

import (
	. "github.com/HeavyWombat/dyff/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Compare", func() {
	Describe("Difference between YAMLs", func() {
		Context("Given two simple YAML structures", func() {
			It("should return that one value was modified", func() {
				from := getYamlFromString(`---
name: foobar
version: v1
`)

				to := getYamlFromString(`---
name: fOObAr
version: v1
`)

				result := CompareObjects(from, to)
				Expect(result).NotTo(BeNil())
				Expect(len(result)).To(BeEquivalentTo(1))

				difference := result[0]
				Expect(difference.Kind).To(BeIdenticalTo(MODIFICATION))
				Expect(difference.From).To(BeIdenticalTo("foobar"))
				Expect(difference.To).To(BeIdenticalTo("fOObAr"))
			})

			It("should return that one value was added", func() {
				from := getYamlFromString(`---
name: foobar
`)

				to := getYamlFromString(`---
name: foobar
version: v1
`)

				result := CompareObjects(from, to)
				Expect(result).NotTo(BeNil())
				Expect(len(result)).To(BeEquivalentTo(1))

				difference := result[0]
				Expect(difference.Kind).To(BeIdenticalTo(ADDITION))
				Expect(difference.From).To(BeNil())
				Expect(difference.To).To(BeIdenticalTo("v1"))
			})

			It("should return that one value was removed", func() {
				from := getYamlFromString(`---
name: foobar
version: v1
`)

				to := getYamlFromString(`---
name: foobar
`)

				result := CompareObjects(from, to)
				Expect(result).NotTo(BeNil())
				Expect(len(result)).To(BeEquivalentTo(1))

				difference := result[0]
				Expect(difference.Kind).To(BeIdenticalTo(REMOVAL))
				Expect(difference.From).To(BeIdenticalTo("v1"))
				Expect(difference.To).To(BeNil())
			})
		})
	})
})
