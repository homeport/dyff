package dyff_test

import (
	"github.com/gonvenience/ytbx"
	"github.com/homeport/dyff/pkg/dyff"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v3"
)

type testdoc struct {
	Foo    map[string]int
	Seq    []string
	SeqMap []map[string]int
}

var a = `---
foo:
  bar: 34
  bak: 27
seq:
  - a
  - b
  - c
  - d
  - 1
  - 2
  - 3
seqmap:
  - asdf: 99
    zxcvb: 10
  - qwerty: 89
`
var b = `---
foo:
  bar: 75
  baz: 0
seq:
  - a
  - c
  - b
  - d
  - 1
  - 2
  - 3
seqmap:
  - asdf: 98
`

var patchYAML = `---
- op: remove
  fromkind: mapping
  tokind: ""
  path: /foo
  tovalue: null
  fromvalue:
    bak: 27
- op: add
  fromkind: ""
  tokind: mapping
  path: /foo
  tovalue:
    baz: 0
  fromvalue: null
- op: replace
  fromkind: scalar
  tokind: scalar
  path: /foo/bar
  tovalue: 75
  fromvalue: 34
- op: reorder
  fromkind: sequence
  tokind: sequence
  path: /seq
  tovalue:
    - a
    - c
    - b
    - d
    - 1
    - 2
    - 3
  fromvalue:
    - a
    - b
    - c
    - d
    - 1
    - 2
    - 3
- op: remove
  fromkind: sequence
  tokind: ""
  path: /seqmap
  tovalue: null
  fromvalue:
    - asdf: 99
      zxcvb: 10
    - qwerty: 89
- op: add
  fromkind: ""
  tokind: sequence
  path: /seqmap
  tovalue:
    - asdf: 98
  fromvalue: null
`

var _ = Describe("patch tests", func() {
	Context("generate patch", func() {
		It("should result in the correct number of patch operations", func() {
			adocs, err := ytbx.LoadYAMLDocuments([]byte(a))
			Expect(err).ToNot(HaveOccurred())

			inputA := ytbx.InputFile{
				Documents: adocs,
			}

			bdocs, err := ytbx.LoadYAMLDocuments([]byte(b))
			Expect(err).ToNot(HaveOccurred())

			inputB := ytbx.InputFile{
				Documents: bdocs,
			}

			report, err := dyff.CompareInputFiles(inputA, inputB)
			Expect(err).ToNot(HaveOccurred())

			p, err := dyff.GeneratePatch(&report)
			Expect(err).ToNot(HaveOccurred())

			_, err = yaml.Marshal(&p)
			Expect(err).ToNot(HaveOccurred())

			Expect(len(p)).To(BeEquivalentTo(6))
		})
	})

	Context("apply patch", func() {
		It("should apply and result in the original YAML", func() {
			adocs, err := ytbx.LoadYAMLDocuments([]byte(a))
			Expect(err).ToNot(HaveOccurred())

			var patch []dyff.PatchOp
			err = yaml.Unmarshal([]byte(patchYAML), &patch)
			Expect(err).ToNot(HaveOccurred())

			err = dyff.ApplyPatch(adocs[0], patch)
			Expect(err).ToNot(HaveOccurred())

			out, err := yaml.Marshal(adocs[0])
			Expect(err).ToNot(HaveOccurred())

			var td testdoc
			err = yaml.Unmarshal(out, &td)
			Expect(err).ToNot(HaveOccurred())

			Expect(td.Foo).To(HaveKey("bar"))
			Expect(td.Foo["bar"]).To(Equal(75))
			Expect(td.Foo).To(HaveKey("baz"))
			Expect(td.Foo["baz"]).To(Equal(0))
			Expect(td.Foo).NotTo(HaveKey("bak"))
			Expect(td.Seq).Should(HaveLen(7))
			Expect(td.Seq[1]).To(Equal("c"))
			Expect(td.Seq[2]).To(Equal("b"))
			Expect(td.SeqMap).Should(HaveLen(1))
			Expect(td.SeqMap[0]).Should(HaveKey("asdf"))
			Expect(td.SeqMap[0]["asdf"]).To(Equal(98))
		})
	})
})
