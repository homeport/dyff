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

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Path test cases", func() {
	Context("path to string in dot-style", func() {
		It("should print out simple hash paths nicely", func() {
			path := Path{DocumentIdx: 0, PathElements: []PathElement{
				{Name: "some"},
				{Name: "deep"},
				{Name: "yaml"},
				{Name: "structure"}}}

			Expect(path.ToDotStyle(false)).To(BeEquivalentTo("some.deep.yaml.structure"))
		})

		It("should print out just the root if it is just the root", func() {
			path := Path{DocumentIdx: 0, PathElements: []PathElement{{Name: "root"}}}

			Expect(path.ToDotStyle(false)).To(BeEquivalentTo("root"))
		})

		It("should print out paths nicely that include named list entries", func() {
			path := Path{DocumentIdx: 0, PathElements: []PathElement{{Name: "some"},
				{Name: "deep"},
				{Name: "yaml"},
				{Name: "structure"},
				{Key: "name", Name: "one"},
				{Name: "enabled"}}}

			Expect(path.ToDotStyle(false)).To(BeEquivalentTo("some.deep.yaml.structure.one.enabled"))
		})

		It("should print out paths nicely that include named list entries which contain named list entries", func() {
			path := Path{DocumentIdx: 0, PathElements: []PathElement{{Name: "some"},
				{Name: "deep"},
				{Name: "yaml"},
				{Name: "structure"},
				{Key: "name", Name: "one"},
				{Name: "list"},
				{Key: "id", Name: "first"}}}

			Expect(path.ToDotStyle(false)).To(BeEquivalentTo("some.deep.yaml.structure.one.list.first"))
		})
	})

	Context("path to string in gopatch-style", func() {
		It("should print out simple hash paths nicely", func() {
			path := Path{DocumentIdx: 0, PathElements: []PathElement{{Name: "some"},
				{Name: "deep"},
				{Name: "yaml"},
				{Name: "structure"}}}

			Expect(path.ToGoPatchStyle(false)).To(BeEquivalentTo("/some/deep/yaml/structure"))
		})

		It("should print out just the root if it is just the root", func() {
			path := Path{DocumentIdx: 0, PathElements: []PathElement{{Name: "root"}}}

			Expect(path.ToGoPatchStyle(false)).To(BeEquivalentTo("/root"))
		})

		It("should print out paths nicely that include named list entries", func() {
			path := Path{DocumentIdx: 0, PathElements: []PathElement{{Name: "some"},
				{Name: "deep"},
				{Name: "yaml"},
				{Name: "structure"},
				{Key: "name", Name: "one"},
				{Name: "enabled"}}}

			Expect(path.ToGoPatchStyle(false)).To(BeEquivalentTo("/some/deep/yaml/structure/name=one/enabled"))
		})

		It("should print out paths nicely that include named list entries which contain named list entries", func() {
			path := Path{DocumentIdx: 0, PathElements: []PathElement{{Name: "some"},
				{Name: "deep"},
				{Name: "yaml"},
				{Name: "structure"},
				{Key: "name", Name: "one"},
				{Name: "list"},
				{Key: "id", Name: "first"}}}

			Expect(path.ToGoPatchStyle(false)).To(BeEquivalentTo("/some/deep/yaml/structure/name=one/list/id=first"))
		})
	})
})
