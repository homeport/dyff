// Copyright © 2018 Matthias Diester
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package dyff_test

import (
	. "github.com/HeavyWombat/dyff/pkg/v1/dyff"
	yaml "gopkg.in/yaml.v2"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Restructure", func() {
	Context("YAML MapSlice key reorderings of the MapSlice itself", func() {
		It("should restructure Concourse root level keys", func() {
			input := yml("{ groups: [], jobs: [], resources: [], resource_types: [] }")
			output := RestructureObject(input).(yaml.MapSlice)

			keys, err := ListStringKeys(output)
			Expect(err).To(BeNil())
			Expect(keys).To(BeEquivalentTo([]string{"jobs", "resources", "resource_types", "groups"}))
		})

		It("should restructure Concourse resource and resource_type keys", func() {
			input := yml("{ source: {}, name: {}, type: {}, privileged: {} }")
			output := RestructureObject(input).(yaml.MapSlice)

			keys, err := ListStringKeys(output)
			Expect(err).To(BeNil())
			Expect(keys).To(BeEquivalentTo([]string{"name", "type", "source", "privileged"}))
		})
	})

	Context("YAML MapSlice key reorderings of the MapSlice values", func() {
		It("should restructure Concourse resource keys as part as part of a MapSlice value", func() {
			input := yml("{ resources: [ { privileged: false, source: { branch: foo, paths: [] }, name: myname, type: mytype } ] }")
			output := RestructureObject(input).(yaml.MapSlice)

			value := output[0].Value.([]interface{})
			obj := value[0].(yaml.MapSlice)

			keys, err := ListStringKeys(obj)
			Expect(err).To(BeNil())
			Expect(keys).To(BeEquivalentTo([]string{"name", "type", "source", "privileged"}))
		})
	})
})
