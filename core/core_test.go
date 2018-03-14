package core_test

import (
	. "github.com/HeavyWombat/dyff/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Core", func() {
	Describe("common functions", func() {
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
	})
})
