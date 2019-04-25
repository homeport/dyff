// Copyright © 2019 The Homeport Team
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
	. "github.com/homeport/dyff/pkg/v1/dyff"
	"github.com/homeport/ytbx/pkg/v1/ytbx"
	yaml "gopkg.in/yaml.v2"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Common core functions", func() {
	Context("convenience test functions", func() {
		It("should return nice lists", func() {
			Expect(list(`[]`)).To(BeEquivalentTo([]interface{}{}))
			Expect(list(`[1]`)).To(BeEquivalentTo([]interface{}{1}))
			Expect(list(`[1,2,3]`)).To(BeEquivalentTo([]interface{}{1, 2, 3}))
			Expect(list(`[A]`)).To(BeEquivalentTo([]interface{}{"A"}))
			Expect(list(`[A,B]`)).To(BeEquivalentTo([]interface{}{"A", "B"}))

			Expect(list(`[{A: B}]`)).To(BeEquivalentTo([]interface{}{yaml.MapSlice{yaml.MapItem{Key: "A", Value: "B"}}}))
			Expect(list(`[{A: B}, {C: D}]`)).To(BeEquivalentTo([]interface{}{
				yaml.MapSlice{yaml.MapItem{Key: "A", Value: "B"}},
				yaml.MapSlice{yaml.MapItem{Key: "C", Value: "D"}},
			}))
		})
	})

	Context("loading input data", func() {
		It("should load input files from disk", func() {
			from, to, err := ytbx.LoadFiles("../../../assets/examples/from.yml", "../../../assets/examples/to.yml")
			Expect(err).To(BeNil())
			Expect(from).ToNot(BeNil())
			Expect(to).ToNot(BeNil())
		})

		It("should load documents from an input string", func() {
			result, err := ytbx.LoadDocuments([]byte(`---
yaml: yes
foo: bar
---
another: yes
foo: bar
---
- type: something
  action: no
`))
			Expect(err).To(BeNil())
			Expect(result[0]).To(BeEquivalentTo(yaml.MapSlice{
				yaml.MapItem{Key: "yaml", Value: true},
				yaml.MapItem{Key: "foo", Value: "bar"},
			}))
			Expect(result[1]).To(BeEquivalentTo(yaml.MapSlice{
				yaml.MapItem{Key: "another", Value: true},
				yaml.MapItem{Key: "foo", Value: "bar"},
			}))
			Expect(result[2]).To(BeEquivalentTo([]yaml.MapSlice{{
				yaml.MapItem{Key: "type", Value: "something"},
				yaml.MapItem{Key: "action", Value: false},
			}}))
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
