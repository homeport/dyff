package core_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/HeavyWombat/dyff/core"

	"github.com/HeavyWombat/yaml"
)

var _ = Describe("Restructure", func() {
	Context("YAML MapSlice key reorderings of the MapSlice itself", func() {
		It("should restructure Concourse root level keys", func() {
			input := yml("{ groups: [], jobs: [], resources: [], resource_types: [] }")
			output := RestructureMapSlice(input)

			keys, err := ListStringKeys(output)
			Expect(err).To(BeNil())
			Expect(keys).To(BeEquivalentTo([]string{"jobs", "resources", "resource_types", "groups"}))
		})

		It("should restructure Concourse resource and resource_type keys", func() {
			input := yml("{ source: {}, name: {}, type: {}, privileged: {} }")
			output := RestructureMapSlice(input)

			keys, err := ListStringKeys(output)
			Expect(err).To(BeNil())
			Expect(keys).To(BeEquivalentTo([]string{"name", "type", "source", "privileged"}))
		})
	})

	Context("YAML MapSlice key reorderings of the MapSlice values", func() {
		It("should restructure Concourse resource keys as part as part of a MapSlice value", func() {
			input := yml("{ resources: [ { privileged: false, source: { branch: foo, paths: [] }, name: myname, type: mytype } ] }")
			output := RestructureMapSlice(input)

			value := output[0].Value.([]interface{})
			obj := value[0].(yaml.MapSlice)

			keys, err := ListStringKeys(obj)
			Expect(err).To(BeNil())
			Expect(keys).To(BeEquivalentTo([]string{"name", "type", "source", "privileged"}))
		})
	})
})
