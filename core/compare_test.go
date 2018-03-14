package core_test

import (
	. "github.com/HeavyWombat/dyff/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Compare", func() {
	Describe("Difference between YAMLs", func() {
		Context("Given two simple YAML structures", func() {
			It("should return that a string value was modified", func() {
				from := getYamlFromString(`---
some:
  yaml:
    structure:
      name: foobar
      version: v1
`)

				to := getYamlFromString(`---
some:
  yaml:
    structure:
      name: fOObAr
      version: v1
`)

				result := CompareDocuments(from, to)
				Expect(result).NotTo(BeNil())
				Expect(len(result)).To(BeEquivalentTo(1))

				difference := result[0]
				Expect(difference.Kind).To(BeIdenticalTo(MODIFICATION))
				Expect(difference.Path.String()).To(BeIdenticalTo("/some/yaml/structure/name"))
				Expect(difference.From).To(BeIdenticalTo("foobar"))
				Expect(difference.To).To(BeIdenticalTo("fOObAr"))
			})

			It("should return that an integer was modified", func() {
				from := getYamlFromString(`---
some:
  yaml:
    structure:
      name: 1
      version: v1
`)

				to := getYamlFromString(`---
some:
  yaml:
    structure:
      name: 2
      version: v1
`)

				result := CompareDocuments(from, to)
				Expect(result).NotTo(BeNil())
				Expect(len(result)).To(BeEquivalentTo(1))

				difference := result[0]
				Expect(difference.Kind).To(BeIdenticalTo(MODIFICATION))
				Expect(difference.Path.String()).To(BeIdenticalTo("/some/yaml/structure/name"))
				Expect(difference.From).To(BeIdenticalTo(1))
				Expect(difference.To).To(BeIdenticalTo(2))
			})

			It("should return that a float was modified", func() {
				from := getYamlFromString(`---
some:
  yaml:
    structure:
      name: foobar
      level: 3.14159265359
`)

				to := getYamlFromString(`---
some:
  yaml:
    structure:
      name: foobar
      level: 2.7182818284
`)

				result := CompareDocuments(from, to)
				Expect(result).NotTo(BeNil())
				Expect(len(result)).To(BeEquivalentTo(1))

				difference := result[0]
				Expect(difference.Kind).To(BeIdenticalTo(MODIFICATION))
				Expect(difference.Path.String()).To(BeIdenticalTo("/some/yaml/structure/level"))
				Expect(difference.From).To(BeNumerically("~", 3.14, 3.15))
				Expect(difference.To).To(BeNumerically("~", 2.71, 2.72))
			})

			It("should return that a boolean was modified", func() {
				from := getYamlFromString(`---
some:
  yaml:
    structure:
      name: foobar
      enabled: false
`)

				to := getYamlFromString(`---
some:
  yaml:
    structure:
      name: foobar
      enabled: true
`)

				result := CompareDocuments(from, to)
				Expect(result).NotTo(BeNil())
				Expect(len(result)).To(BeEquivalentTo(1))

				difference := result[0]
				Expect(difference.Kind).To(BeIdenticalTo(MODIFICATION))
				Expect(difference.Path.String()).To(BeIdenticalTo("/some/yaml/structure/enabled"))
				Expect(difference.From).To(BeIdenticalTo(false))
				Expect(difference.To).To(BeIdenticalTo(true))
			})

			It("should return that one value was added", func() {
				from := getYamlFromString(`---
some:
  yaml:
    structure:
      name: foobar
`)

				to := getYamlFromString(`---
some:
  yaml:
    structure:
      name: foobar
      version: v1
`)

				result := CompareDocuments(from, to)
				Expect(result).NotTo(BeNil())
				Expect(len(result)).To(BeEquivalentTo(1))

				difference := result[0]
				Expect(difference.Kind).To(BeIdenticalTo(ADDITION))
				Expect(difference.Path.String()).To(BeIdenticalTo("/some/yaml/structure/version"))
				Expect(difference.From).To(BeNil())
				Expect(difference.To).To(BeIdenticalTo("v1"))
			})

			It("should return that one value was removed", func() {
				from := getYamlFromString(`---
some:
  yaml:
    structure:
      name: foobar
      version: v1
`)

				to := getYamlFromString(`---
some:
  yaml:
    structure:
      name: foobar
`)

				result := CompareDocuments(from, to)
				Expect(result).NotTo(BeNil())
				Expect(len(result)).To(BeEquivalentTo(1))

				difference := result[0]
				Expect(difference.Kind).To(BeIdenticalTo(REMOVAL))
				Expect(difference.Path.String()).To(BeIdenticalTo("/some/yaml/structure/version"))
				Expect(difference.From).To(BeIdenticalTo("v1"))
				Expect(difference.To).To(BeNil())
			})

			It("should return two diffs if one value was removed and another one added", func() {
				from := getYamlFromString(`---
some:
  yaml:
    structure:
      name: foobar
      version: v1
`)

				to := getYamlFromString(`---
some:
  yaml:
    structure:
      name: foobar
      release: v1
`)

				result := CompareDocuments(from, to)
				Expect(result).NotTo(BeNil())
				Expect(len(result)).To(BeEquivalentTo(2))

				Expect(result[0].Kind).To(BeIdenticalTo(REMOVAL))
				Expect(result[0].Path.String()).To(BeIdenticalTo("/some/yaml/structure/version"))
				Expect(result[0].From).To(BeIdenticalTo("v1"))
				Expect(result[0].To).To(BeNil())

				Expect(result[1].Kind).To(BeIdenticalTo(ADDITION))
				Expect(result[1].Path.String()).To(BeIdenticalTo("/some/yaml/structure/release"))
				Expect(result[1].From).To(BeNil())
				Expect(result[1].To).To(BeIdenticalTo("v1"))
			})
		})

		Context("Given two YAML structures with simple lists", func() {
			It("should return that a string list entry was added", func() {
				from := getYamlFromString(`---
some:
  yaml:
    structure:
      list:
      - one
      - two
`)

				to := getYamlFromString(`---
some:
  yaml:
    structure:
      list:
      - one
      - two
      - three
`)

				result := CompareDocuments(from, to)
				Expect(result).NotTo(BeNil())
				Expect(len(result)).To(BeEquivalentTo(1))

				Expect(result[0].Kind).To(BeIdenticalTo(ADDITION))
				Expect(result[0].Path.String()).To(BeIdenticalTo("/some/yaml/structure/list"))
				Expect(result[0].From).To(BeNil())
				Expect(result[0].To).To(BeIdenticalTo("three"))
			})

			It("should return that an integer list entry was added", func() {
				from := getYamlFromString(`---
some:
  yaml:
    structure:
      list:
      - 1
      - 2
`)

				to := getYamlFromString(`---
some:
  yaml:
    structure:
      list:
      - 1
      - 2
      - 3
`)

				result := CompareDocuments(from, to)
				Expect(result).NotTo(BeNil())
				Expect(len(result)).To(BeEquivalentTo(1))

				Expect(result[0].Kind).To(BeIdenticalTo(ADDITION))
				Expect(result[0].Path.String()).To(BeIdenticalTo("/some/yaml/structure/list"))
				Expect(result[0].From).To(BeNil())
				Expect(result[0].To).To(BeIdenticalTo(3))
			})

			It("should return that a string list entry was removed", func() {
				from := getYamlFromString(`---
some:
  yaml:
    structure:
      list:
      - one
      - two
      - three
`)

				to := getYamlFromString(`---
some:
  yaml:
    structure:
      list:
      - one
      - two
`)

				result := CompareDocuments(from, to)
				Expect(result).NotTo(BeNil())
				Expect(len(result)).To(BeEquivalentTo(1))

				Expect(result[0].Kind).To(BeIdenticalTo(REMOVAL))
				Expect(result[0].Path.String()).To(BeIdenticalTo("/some/yaml/structure/list"))
				Expect(result[0].From).To(BeIdenticalTo("three"))
				Expect(result[0].To).To(BeNil())
			})

			It("should return that an integer list entry was removed", func() {
				from := getYamlFromString(`---
some:
  yaml:
    structure:
      list:
      - 1
      - 2
      - 3
`)

				to := getYamlFromString(`---
some:
  yaml:
    structure:
      list:
      - 1
      - 2
`)

				result := CompareDocuments(from, to)
				Expect(result).NotTo(BeNil())
				Expect(len(result)).To(BeEquivalentTo(1))

				Expect(result[0].Kind).To(BeIdenticalTo(REMOVAL))
				Expect(result[0].Path.String()).To(BeIdenticalTo("/some/yaml/structure/list"))
				Expect(result[0].From).To(BeIdenticalTo(3))
				Expect(result[0].To).To(BeNil())
			})
		})
	})
})
