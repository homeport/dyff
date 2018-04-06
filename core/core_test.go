package core_test

import (
	. "github.com/HeavyWombat/dyff/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Core/Functions", func() {
	Describe("common functions", func() {
		Context("loading input data", func() {
			It("should load input files from disk", func() {
				from, to, err := LoadFiles("../assets/examples/from.yml", "../assets/examples/to.yml")
				Expect(err).To(BeNil())
				Expect(from).ToNot(BeNil())
				Expect(to).ToNot(BeNil())
			})
		})

		Context("creating proper texts", func() {
			It("should return human readable plurals", func() {
				Expect(Plural(0, "foobar")).To(BeEquivalentTo("zero foobars"))
				Expect(Plural(1, "foobar")).To(BeEquivalentTo("one foobar"))
				Expect(Plural(2, "foobar")).To(BeEquivalentTo("two foobars"))
				Expect(Plural(3, "foobar")).To(BeEquivalentTo("three foobars"))
				Expect(Plural(4, "foobar")).To(BeEquivalentTo("four foobars"))
				Expect(Plural(5, "foobar")).To(BeEquivalentTo("five foobars"))
				Expect(Plural(6, "foobar")).To(BeEquivalentTo("six foobars"))
				Expect(Plural(7, "foobar")).To(BeEquivalentTo("seven foobars"))
				Expect(Plural(8, "foobar")).To(BeEquivalentTo("eight foobars"))
				Expect(Plural(9, "foobar")).To(BeEquivalentTo("nine foobars"))
				Expect(Plural(10, "foobar")).To(BeEquivalentTo("ten foobars"))
				Expect(Plural(11, "foobar")).To(BeEquivalentTo("eleven foobars"))
				Expect(Plural(12, "foobar")).To(BeEquivalentTo("twelve foobars"))
				Expect(Plural(13, "foobar")).To(BeEquivalentTo("13 foobars"))
				Expect(Plural(147, "foobar")).To(BeEquivalentTo("147 foobars"))

				Expect(Plural(1, "basis", "bases")).To(BeEquivalentTo("one basis"))
				Expect(Plural(2, "basis", "bases")).To(BeEquivalentTo("two bases"))
			})
		})

		Context("path to string in dot-style", func() {
			It("should print out simple hash paths nicely", func() {
				path := Path{PathElement{Name: "some"},
					PathElement{Name: "deep"},
					PathElement{Name: "yaml"},
					PathElement{Name: "structure"}}

				output := ToDotStyle(path)
				Expect(output).To(BeEquivalentTo("some.deep.yaml.structure"))
			})

			It("should print out just the root if it is just the root", func() {
				path := Path{PathElement{Name: "root"}}

				output := ToDotStyle(path)
				Expect(output).To(BeEquivalentTo("root"))
			})

			It("should print out paths nicely that include named list entries", func() {
				path := Path{PathElement{Name: "some"},
					PathElement{Name: "deep"},
					PathElement{Name: "yaml"},
					PathElement{Name: "structure"},
					PathElement{Key: "name", Name: "one"},
					PathElement{Name: "enabled"}}

				output := ToDotStyle(path)
				Expect(output).To(BeEquivalentTo("some.deep.yaml.structure.one.enabled"))
			})

			It("should print out paths nicely that include named list entries which contain named list entries", func() {
				path := Path{PathElement{Name: "some"},
					PathElement{Name: "deep"},
					PathElement{Name: "yaml"},
					PathElement{Name: "structure"},
					PathElement{Key: "name", Name: "one"},
					PathElement{Name: "list"},
					PathElement{Key: "id", Name: "first"}}

				output := ToDotStyle(path)
				Expect(output).To(BeEquivalentTo("some.deep.yaml.structure.one.list.first"))
			})
		})

		Context("path to string in gopatch-style", func() {
			It("should print out simple hash paths nicely", func() {
				path := Path{PathElement{Name: "some"},
					PathElement{Name: "deep"},
					PathElement{Name: "yaml"},
					PathElement{Name: "structure"}}

				output := ToGoPatchStyle(path)
				Expect(output).To(BeEquivalentTo("/some/deep/yaml/structure"))
			})

			It("should print out just the root if it is just the root", func() {
				path := Path{PathElement{Name: "root"}}

				output := ToGoPatchStyle(path)
				Expect(output).To(BeEquivalentTo("/root"))
			})

			It("should print out paths nicely that include named list entries", func() {
				path := Path{PathElement{Name: "some"},
					PathElement{Name: "deep"},
					PathElement{Name: "yaml"},
					PathElement{Name: "structure"},
					PathElement{Key: "name", Name: "one"},
					PathElement{Name: "enabled"}}

				output := ToGoPatchStyle(path)
				Expect(output).To(BeEquivalentTo("/some/deep/yaml/structure/name=one/enabled"))
			})

			It("should print out paths nicely that include named list entries which contain named list entries", func() {
				path := Path{PathElement{Name: "some"},
					PathElement{Name: "deep"},
					PathElement{Name: "yaml"},
					PathElement{Name: "structure"},
					PathElement{Key: "name", Name: "one"},
					PathElement{Name: "list"},
					PathElement{Key: "id", Name: "first"}}

				output := ToGoPatchStyle(path)
				Expect(output).To(BeEquivalentTo("/some/deep/yaml/structure/name=one/list/id=first"))
			})
		})

		Context("identify the main identifier key in named lists", func() {
			It("should return 'name' as the main identifier if list uses 'name'", func() {
				sample := yml(`---
list:
- name: one
  version: v1
- name: two
  version: v1
`)

				output := GetIdentifierFromNamedList(sample[0].Value.([]interface{}))
				Expect(output).To(BeEquivalentTo("name"))
			})

			It("should return 'key' as the main identifier if list uses 'key'", func() {
				sample := yml(`---
list:
- key: one
  version: v1
- key: two
  version: v1
`)

				output := GetIdentifierFromNamedList(sample[0].Value.([]interface{}))
				Expect(output).To(BeEquivalentTo("key"))
			})

			It("should return 'id' as the main identifier if list uses 'id'", func() {
				sample := yml(`---
list:
- id: one
  version: v1
- id: two
  version: v1
`)

				output := GetIdentifierFromNamedList(sample[0].Value.([]interface{}))
				Expect(output).To(BeEquivalentTo("id"))
			})

			It("should return nothing as the main identifier if there is no common identifier", func() {
				sample := yml(`---
list:
- name: one
  version: v1
- id: two
  version: v1
`)

				output := GetIdentifierFromNamedList(sample[0].Value.([]interface{}))
				Expect(output).To(BeEmpty())
			})
		})
	})
})
