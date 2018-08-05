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
			from, to, err := LoadFiles("../../../assets/examples/from.yml", "../../../assets/examples/to.yml")
			Expect(err).To(BeNil())
			Expect(from).ToNot(BeNil())
			Expect(to).ToNot(BeNil())
		})

		It("should load documents from an input string", func() {
			result, err := LoadDocuments([]byte(`---
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

	Context("creating proper texts", func() {
		It("should return human readable plurals", func() {
			Expect(Plural(0, "foobar")).To(BeEquivalentTo("no foobars"))
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

	Context("Grabing values by path", func() {
		It("should create the same path using Go-patch and Spruce style", func() {
			obj := yml("../../../assets/bosh-yaml/manifest.yml")
			Expect(obj).ToNot(BeNil())

			Expect(pathFromString("/name", obj)).To(
				BeEquivalentTo(pathFromString("name", obj)))

			Expect(pathFromString("/releases/name=concourse", obj)).To(
				BeEquivalentTo(pathFromString("releases.concourse", obj)))

			Expect(pathFromString("/instance_groups/name=web/networks/name=concourse/static_ips/0", obj)).To(
				BeEquivalentTo(pathFromString("instance_groups.web.networks.concourse.static_ips.0", obj)))

			Expect(pathFromString("/networks/name=concourse/subnets/0/cloud_properties/name", obj)).To(
				BeEquivalentTo(pathFromString("networks.concourse.subnets.0.cloud_properties.name", obj)))
		})

		It("should return the value referenced by the path", func() {
			example := yml("../../../assets/examples/from.yml")
			Expect(example).ToNot(BeNil())

			Expect(grab(example, "/yaml/map/before")).To(BeEquivalentTo("after"))
			Expect(grab(example, "/yaml/map/intA")).To(BeEquivalentTo(42))
			Expect(grab(example, "/yaml/map/mapA")).To(BeEquivalentTo(yml(`{ key0: A, key1: A }`)))
			Expect(grab(example, "/yaml/map/listA")).To(BeEquivalentTo(list(`[ A, A, A ]`)))

			Expect(grab(example, "/yaml/named-entry-list-using-name/name=B")).To(BeEquivalentTo(yml(`{ name: B }`)))
			Expect(grab(example, "/yaml/named-entry-list-using-key/key=B")).To(BeEquivalentTo(yml(`{ key: B }`)))
			Expect(grab(example, "/yaml/named-entry-list-using-id/id=B")).To(BeEquivalentTo(yml(`{ id: B }`)))

			Expect(grab(example, "/yaml/simple-list/1")).To(BeEquivalentTo("B"))
			Expect(grab(example, "/yaml/named-entry-list-using-key/3")).To(BeEquivalentTo(yml(`{ key: X }`)))

			// --- --- ---

			example = yml("../../../assets/bosh-yaml/manifest.yml")
			Expect(example).ToNot(BeNil())

			Expect(grab(example, "/instance_groups/name=web/networks/name=concourse/static_ips/0")).To(BeEquivalentTo("XX.XX.XX.XX"))
			Expect(grab(example, "/instance_groups/name=worker/jobs/name=baggageclaim/properties")).To(BeEquivalentTo(yml(`{}`)))
		})

		It("should return the value referenced by the path", func() {
			example := yml("../../../assets/examples/from.yml")
			Expect(example).ToNot(BeNil())

			Expect(grabError(example, "/yaml/simple-list/-1")).To(BeEquivalentTo("failed to traverse tree, provided list index -1 is not in range: 0..4"))
			Expect(grabError(example, "/yaml/does-not-exist")).To(BeEquivalentTo("no key 'does-not-exist' found in map, available keys are: map, simple-list, named-entry-list-using-name, named-entry-list-using-key, named-entry-list-using-id"))
			Expect(grabError(example, "/yaml/0")).To(BeEquivalentTo("failed to traverse tree, expected a list but found type map at /yaml"))
			Expect(grabError(example, "/yaml/simple-list/foobar")).To(BeEquivalentTo("failed to traverse tree, expected a map but found type list at /yaml/simple-list"))
			Expect(grabError(example, "/yaml/map/foobar=0")).To(BeEquivalentTo("failed to traverse tree, expected a list but found type map at /yaml/map"))
			Expect(grabError(example, "/yaml/named-entry-list-using-id/id=0")).To(BeEquivalentTo("there is no entry id: 0 in the list"))
		})
	})
})
