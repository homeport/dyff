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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/homeport/dyff/pkg/dyff"

	"github.com/gonvenience/ytbx"
	yamlv3 "gopkg.in/yaml.v3"
)

var nullNode = &yamlv3.Node{
	Kind:  yamlv3.ScalarNode,
	Tag:   "!!null",
	Value: "null",
}

var _ = Describe("Core/Compare", func() {
	Describe("Difference between YAMLs", func() {
		Context("Given two simple YAML structures", func() {
			It("should return that a string value was modified", func() {
				from := yml(`---
some:
  yaml:
    structure:
      name: foobar
      version: v1
`)

				to := yml(`---
some:
  yaml:
    structure:
      name: fOObAr
      version: v1
`)

				result, err := compare(from, to)
				Expect(err).To(BeNil())
				Expect(result).NotTo(BeNil())
				Expect(result).To(HaveLen(1))
				Expect(result[0].Details).To(HaveLen(1))
				Expect(result[0]).To(BeSameDiffAs(singleDiff("/some/yaml/structure/name", dyff.MODIFICATION, "foobar", "fOObAr")))
			})

			It("should return that an integer was modified", func() {
				from := yml(`---
some:
  yaml:
    structure:
      name: 1
      version: v1
`)

				to := yml(`---
some:
  yaml:
    structure:
      name: 2
      version: v1
`)

				result, err := compare(from, to)
				Expect(err).To(BeNil())
				Expect(result).NotTo(BeNil())
				Expect(result).To(HaveLen(1))
				Expect(result[0]).To(BeSameDiffAs(singleDiff("/some/yaml/structure/name", dyff.MODIFICATION, 1, 2)))
			})

			It("should return that different representations of true are treated as the same", func() {
				from := yml("---\nkey: true")

				to := yml("---\nkey: True")

				result, err := compare(from, to)
				Expect(err).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should return that a float was modified", func() {
				from := yml(`---
some:
  yaml:
    structure:
      name: foobar
      level: 3.14159265359
`)

				to := yml(`---
some:
  yaml:
    structure:
      name: foobar
      level: 2.7182818284
`)

				result, err := compare(from, to)
				Expect(err).To(BeNil())
				Expect(result).NotTo(BeNil())
				Expect(result).To(HaveLen(1))
				Expect(result[0]).To(BeSameDiffAs(singleDiff("/some/yaml/structure/level", dyff.MODIFICATION, 3.14159265359, 2.7182818284)))
			})

			It("should return that a boolean was modified", func() {
				from := yml(`---
some:
  yaml:
    structure:
      name: foobar
      enabled: false
`)

				to := yml(`---
some:
  yaml:
    structure:
      name: foobar
      enabled: true
`)

				result, err := compare(from, to)
				Expect(err).To(BeNil())
				Expect(result).NotTo(BeNil())
				Expect(result).To(HaveLen(1))
				Expect(result[0]).To(BeSameDiffAs(singleDiff("/some/yaml/structure/enabled", dyff.MODIFICATION, false, true)))
			})

			It("should return that one value was added", func() {
				from := yml(`---
some:
  yaml:
    structure:
      name: foobar
`)

				to := yml(`---
some:
  yaml:
    structure:
      name: foobar
      version: v1
`)

				result, err := compare(from, to)
				Expect(err).To(BeNil())
				Expect(result).NotTo(BeNil())
				Expect(result).To(HaveLen(1))
				Expect(result[0]).To(BeSameDiffAs(singleDiff("/some/yaml/structure", dyff.ADDITION, nil, yml(`version: v1`))))
			})

			It("should return that one value was removed", func() {
				from := yml(`---
some:
  yaml:
    structure:
      name: foobar
      version: v1
`)

				to := yml(`---
some:
  yaml:
    structure:
      name: foobar
`)

				result, err := compare(from, to)
				Expect(err).To(BeNil())
				Expect(result).NotTo(BeNil())
				Expect(result).To(HaveLen(1))
				Expect(result[0]).To(BeSameDiffAs(singleDiff("/some/yaml/structure", dyff.REMOVAL, yml(`version: v1`), nil)))
			})

			It("should return two diffs if one value was removed and another one added", func() {
				from := yml(`---
some:
  yaml:
    structure:
      name: foobar
      version: v1
`)

				to := yml(`---
some:
  yaml:
    structure:
      name: foobar
      release: v1
`)

				result, err := compare(from, to)
				Expect(err).To(BeNil())
				Expect(result).NotTo(BeNil())
				Expect(result).To(HaveLen(1))
				Expect(result[0]).To(BeSameDiffAs(doubleDiff("/some/yaml/structure",
					dyff.REMOVAL, yml(`version: v1`), nil,
					dyff.ADDITION, nil, yml(`release: v1`))))
			})

			It("should ignore leading and trailing whitespace changes if configured", func() {
				from := yml(`{"foo": "bar"}`)
				to := yml(`{"foo": "bar "}`)

				diffs, err := compare(from, to, dyff.IgnoreWhitespaceChanges(true))
				Expect(err).To(BeNil())
				Expect(diffs).To(BeNil())
			})
		})

		Context("Given two YAML structures with simple lists", func() {
			It("should return that a string list entry was added", func() {
				from := yml(`---
some:
  yaml:
    structure:
      list:
      - one
      - two
`)

				to := yml(`---
some:
  yaml:
    structure:
      list:
      - one
      - two
      - three
`)

				result, err := compare(from, to)
				Expect(err).To(BeNil())
				Expect(result).NotTo(BeNil())
				Expect(result).To(HaveLen(1))
				Expect(result[0]).To(BeSameDiffAs(singleDiff("/some/yaml/structure/list", dyff.ADDITION, nil, list(`[ three ]`))))
			})

			It("should return that an integer list entry was added", func() {
				from := yml(`---
some:
  yaml:
    structure:
      list:
      - 1
      - 2
`)

				to := yml(`---
some:
  yaml:
    structure:
      list:
      - 1
      - 2
      - 3
`)

				result, err := compare(from, to)
				Expect(err).To(BeNil())
				Expect(result).NotTo(BeNil())
				Expect(result).To(HaveLen(1))
				Expect(result[0]).To(BeSameDiffAs(singleDiff("/some/yaml/structure/list", dyff.ADDITION, nil, list(`[ 3 ]`))))
			})

			It("should return that a string list entry was removed", func() {
				from := yml(`---
some:
  yaml:
    structure:
      list:
      - one
      - two
      - three
`)

				to := yml(`---
some:
  yaml:
    structure:
      list:
      - one
      - two
`)

				result, err := compare(from, to)
				Expect(err).To(BeNil())
				Expect(result).NotTo(BeNil())
				Expect(result).To(HaveLen(1))
				Expect(result[0]).To(BeSameDiffAs(singleDiff("/some/yaml/structure/list", dyff.REMOVAL, list(`[ three ]`), nil)))
			})

			It("should return that an integer list entry was removed", func() {
				from := yml(`---
some:
  yaml:
    structure:
      list:
      - 1
      - 2
      - 3
`)

				to := yml(`---
some:
  yaml:
    structure:
      list:
      - 1
      - 2
`)

				result, err := compare(from, to)
				Expect(err).To(BeNil())
				Expect(result).NotTo(BeNil())
				Expect(result).To(HaveLen(1))
				Expect(result[0]).To(BeSameDiffAs(singleDiff("/some/yaml/structure/list", dyff.REMOVAL, list(`[ 3 ]`), nil)))
			})

			It("should not return a change if only the order in a hash was changed", func() {
				from := yml(`---
list:
- enabled: true
- foo: bar
  version: 1
`)

				to := yml(`---
list:
- enabled: true
- version: 1
  foo: bar
`)

				result, err := compare(from, to)
				Expect(err).To(BeNil())
				Expect(result).To(HaveLen(0))
			})
		})

		Context("Given two YAML structures with complex content", func() {
			It("should return all differences in there", func() {
				from := yml(`---
instance_groups:
- name: web
  instances: 1
  resource_pool: concourse_resource_pool
  networks:
  - name: concourse
    static_ips: 192.168.1.1
  jobs:
  - release: concourse
    name: atc
    properties:
      postgresql_database: &atc-db atc
      external_url: http://192.168.1.100:8080
      development_mode: true
  - release: concourse
    name: tsa
    properties: {}

  - name: db
    instances: 1
    resource_pool: concourse_resource_pool
    networks: [{name: concourse}, {name: testnet}]
    persistent_disk: 10240
    jobs:
    - release: concourse
      name: postgresql
      properties:
        databases:
        - name: *atc-db
          role: atc
          password: supersecret
`)

				to := yml(`---
instance_groups:
- name: web
  instances: 1
  resource_pool: concourse_resource_pool
  networks:
  - name: concourse
    static_ips: 192.168.0.1
  jobs:
  - release: concourse
    name: atc
    properties:
      postgresql_database: &atc-db atc
      external_url: http://192.168.0.100:8080
      development_mode: false
  - release: concourse
    name: tsa
    properties: {}
  - release: custom
    name: logger

  - name: db
    instances: 2
    resource_pool: concourse_resource_pool
    networks: [{name: concourse}]
    persistent_disk: 10240
    jobs:
    - release: concourse
      name: postgresql
      properties:
        databases:
        - name: *atc-db
          role: atc
          password: "zwX#(;P=%hTfFzM["
`)

				result, err := compare(from, to)
				Expect(err).To(BeNil())
				Expect(result).NotTo(BeNil())
				Expect(result).To(HaveLen(7))
				Expect(result[0]).To(BeSameDiffAs(singleDiff("/instance_groups/name=web/networks/name=concourse/static_ips", dyff.MODIFICATION, "192.168.1.1", "192.168.0.1")))
				Expect(result[1]).To(BeSameDiffAs(singleDiff("/instance_groups/name=web/jobs", dyff.ADDITION, nil, list(`[ { release: custom, name: logger } ]`))))
				Expect(result[2]).To(BeSameDiffAs(singleDiff("/instance_groups/name=web/jobs/name=atc/properties/external_url", dyff.MODIFICATION, "http://192.168.1.100:8080", "http://192.168.0.100:8080")))
				Expect(result[3]).To(BeSameDiffAs(singleDiff("/instance_groups/name=web/jobs/name=atc/properties/development_mode", dyff.MODIFICATION, true, false)))
				Expect(result[4]).To(BeSameDiffAs(singleDiff("/instance_groups/name=web/jobs/name=db/instances", dyff.MODIFICATION, 1, 2)))
				Expect(result[5]).To(BeSameDiffAs(singleDiff("/instance_groups/name=web/jobs/name=db/networks", dyff.REMOVAL, list(`[ { name: testnet } ]`), nil)))
				Expect(result[6]).To(BeSameDiffAs(singleDiff("/instance_groups/name=web/jobs/name=db/jobs/name=postgresql/properties/databases/name=atc/password", dyff.MODIFICATION, "supersecret", "zwX#(;P=%hTfFzM[")))
			})

			It("should return even difficult ones", func() {
				from := yml(`---
resource_pools:
- name: concourse_resource_pool
  stemcell:
    name: bosh-vsphere-esxi-ubuntu-trusty-go_agent
    version: '3232.2'
  network: concourse
  cloud_properties:
    ram: 4096
    disk: 32768
    cpu: 2
    datacenters:
    - clusters:
      - CLS_PAAS_SFT_035: {resource_pool: other-vsphere-res-pool}
`)

				to := yml(`---
resource_pools:
- name: concourse_resource_pool
  stemcell:
    name: bosh-vsphere-esxi-ubuntu-trusty-go_agent
    version: '3232.2'
  network: concourse
  cloud_properties:
    ram: 4096
    disk: 32768
    cpu: 2
    datacenters:
    - clusters:
      - CLS_PAAS_SFT_035:
          resource_pool: new-vsphere-res-pool
`)

				result, err := compare(from, to)
				Expect(err).To(BeNil())
				Expect(result).NotTo(BeNil())
				Expect(result).To(HaveLen(1))
				Expect(result[0]).To(BeSameDiffAs(singleDiff("/resource_pools/name=concourse_resource_pool/cloud_properties/datacenters/0/clusters/0/CLS_PAAS_SFT_035/resource_pool", dyff.MODIFICATION, "other-vsphere-res-pool", "new-vsphere-res-pool")))
			})

			It("should return even difficult ones that cannot simply be compared", func() {
				from := yml(`---
resource_pools:
- name: concourse_resource_pool
  cloud_properties:
    datacenters:
    - clusters:
      - CLS_PAAS_SFT_035: {resource_pool: 35-vsphere-res-pool}
      - CLS_PAAS_SFT_036: {resource_pool: 36-vsphere-res-pool}
`)

				to := yml(`---
resource_pools:
- name: concourse_resource_pool
  cloud_properties:
    datacenters:
    - clusters:
      - CLS_PAAS_SFT_035: {resource_pool: 35a-vsphere-res-pool}
      - CLS_PAAS_SFT_036: {resource_pool: 36a-vsphere-res-pool}
`)

				result, err := compare(from, to)
				Expect(err).To(BeNil())
				Expect(result).NotTo(BeNil())
				Expect(result).To(HaveLen(1))
				Expect(result[0]).To(BeSameDiffAs(doubleDiff("/resource_pools/name=concourse_resource_pool/cloud_properties/datacenters/0/clusters",
					dyff.REMOVAL, list(`[ {CLS_PAAS_SFT_035: {resource_pool: 35-vsphere-res-pool}}, {CLS_PAAS_SFT_036: {resource_pool: 36-vsphere-res-pool}} ]`), nil,
					dyff.ADDITION, nil, list(`[ {CLS_PAAS_SFT_035: {resource_pool: 35a-vsphere-res-pool}}, {CLS_PAAS_SFT_036: {resource_pool: 36a-vsphere-res-pool}} ]`))))
			})
		})

		Context("Given two YAMLs with a list as the root", func() {
			It("should return the differences the same way", func() {
				from := list(`---
- name: one
  version: 1

- name: two
  version: 2

- name: three
  version: 4
`)

				to := list(`---
- name: one
  version: 1

- name: two
  version: 2

- name: three
  version: 3
`)

				result, err := compare(from, to)
				Expect(err).To(BeNil())
				Expect(result).NotTo(BeNil())
				Expect(result).To(HaveLen(1))
				Expect(result[0]).To(BeSameDiffAs(singleDiff("/name=three/version",
					dyff.MODIFICATION, 4, 3)))
			})
		})

		Context("Given two files", func() {
			It("should return differences in raw texts", func() {
				from := file(assets("raw-text/from.txt"))
				to := file(assets("raw-text/to.txt"))
				Expect(from.Documents).To(HaveLen(1))
				Expect(to.Documents).To(HaveLen(1))

				results, err := compare(from.Documents[0], to.Documents[0])
				Expect(err).To(BeNil())
				Expect(results).NotTo(BeNil())
				Expect(results).To(HaveLen(1))
			})
		})

		Context("Given two YAML files", func() {
			It("should return all differences between the files", func() {
				results, err := compare(yml("../../assets/examples/from.yml"), yml("../../assets/examples/to.yml"))
				Expect(err).To(BeNil())
				expected := []dyff.Diff{
					doubleDiff("/yaml/map",
						dyff.REMOVAL, yml(`---
stringB: fOObAr
intB: 10
floatB: 2.71
boolB: false
mapB: { key0: B, key1: B }
listB: [ B, B, B ]
`), nil,
						dyff.ADDITION, nil, yml(`---
stringY: YAML!
intY: 147
floatY: 24.0
boolY: true
mapY: { key0: Y, key1: Y }
listY: [ Yo, Yo, Yo ]
`)),

					singleDiff("/yaml/map/type-change-1", dyff.MODIFICATION, "string", 147),

					singleDiff("/yaml/map/type-change-2", dyff.MODIFICATION, "12", 12),

					singleDiff("/yaml/map/whitespaces", dyff.MODIFICATION, "Strings can  have whitespaces.", "Strings can  have whitespaces.\n\n\n"),

					doubleDiff("/yaml/simple-list",
						dyff.REMOVAL, list(`[ X, Z ]`), nil,
						dyff.ADDITION, nil, list(`[ D, E ]`)),

					doubleDiff("/yaml/named-entry-list-using-name",
						dyff.REMOVAL, list(`[ {name: X}, {name: Z} ]`), nil,
						dyff.ADDITION, nil, list(`[ {name: D}, {name: E} ]`)),

					doubleDiff("/yaml/named-entry-list-using-key",
						dyff.REMOVAL, list(`[ {key: X}, {key: Z} ]`), nil,
						dyff.ADDITION, nil, list(`[ {key: D}, {key: E} ]`)),

					doubleDiff("/yaml/named-entry-list-using-id",
						dyff.REMOVAL, list(`[ {id: X}, {id: Z} ]`), nil,
						dyff.ADDITION, nil, list(`[ {id: D}, {id: E} ]`)),
				}

				Expect(results).NotTo(BeNil())
				Expect(results).To(HaveLen(len(expected)))

				for i, result := range results {
					Expect(result).To(BeSameDiffAs(expected[i]))
				}
			})

			It("should return order changes in named entry lists (ignoring dyff.additions and dyff.removals)", func() {
				from := yml(`list: [ {name: A}, {name: C}, {name: B}, {name: D}, {name: E} ]`)
				to := yml(`list: [ {name: A}, {name: X1}, {name: B}, {name: C}, {name: D}, {name: X2} ]`)
				results, err := compare(from, to)
				Expect(err).To(BeNil())
				Expect(results).NotTo(BeNil())
				Expect(results).To(HaveLen(1))
				Expect(results[0].Details).To(HaveLen(3))
				Expect(results[0].Details[0]).To(BeEquivalentTo(dyff.Detail{
					Kind: dyff.ORDERCHANGE,
					From: dyff.AsSequenceNode("A", "C", "B", "D"),
					To:   dyff.AsSequenceNode("A", "B", "C", "D"),
				}))
			})

			It("should return changes where one of the items is null", func() {
				from := yml(`foo: null`)
				to := yml(`foo: "bar"`)
				results, err := compare(from, to)
				Expect(err).To(BeNil())
				Expect(results).NotTo(BeNil())
				Expect(results).To(HaveLen(1))
				Expect(results[0]).To(BeSameDiffAs(singleDiff(
					"/foo",
					dyff.MODIFICATION,
					nullNode,
					"bar",
				)))
			})

			It("should not return order changes in named entry lists in case the ignore option is enabled", func() {
				results, err := compare(
					yml(`list: [ {name: A}, {name: C}, {name: B}, {name: D}, {name: E} ]`),
					yml(`list: [ {name: A}, {name: B}, {name: C}, {name: D}, {name: E} ]`),
					dyff.IgnoreOrderChanges(true),
				)

				Expect(err).To(BeNil())
				Expect(results).To(HaveLen(0))
			})

			It("should return order changes in simple lists (ignoring dyff.additions and dyff.removals)", func() {
				from := yml(`list: [ A, C, B, D, E ]`)
				to := yml(`list: [ A, X1, B, C, D, X2 ]`)
				results, err := compare(from, to)
				Expect(err).To(BeNil())
				Expect(results).NotTo(BeNil())
				Expect(results).To(HaveLen(1))
				Expect(results[0].Details).To(HaveLen(3))

				actual := results[0].Details[0]
				expected := dyff.Detail{
					Kind: dyff.ORDERCHANGE,
					From: dyff.AsSequenceNode("A", "C", "B", "D"),
					To:   dyff.AsSequenceNode("A", "B", "C", "D"),
				}

				Expect(isSameDetail(actual, expected)).To(BeTrue())
			})

			It("should process simple lists independent of other compare settings", func() {
				from := yml(`list: [ A, B, C, D, E ]`)
				to := yml(`list: [ A, B, C, D, E ]`)

				results, err := compare(from, to, dyff.KubernetesEntityDetection(false))
				Expect(err).To(BeNil())
				Expect(results).To(BeNil())
			})

			It("should return all differences between the files with multiple documents", func() {
				expected := []dyff.Diff{
					singleDiff("#0/spec/template/spec/containers/name=registry/resources/limits/cpu", dyff.MODIFICATION, "100m", "1000m"),
					singleDiff("#0/spec/template/spec/containers/name=registry/resources/limits/memory", dyff.MODIFICATION, "100Mi", "10Gi"),
					singleDiff("#0/spec/template/spec/containers/name=registry/ports", dyff.ADDITION, nil, list(`[ {containerPort: 5001, name: backdoor, protocol: TCP} ]`)),
					singleDiff("#1/spec/ports", dyff.ADDITION, nil, list(`[ {name: backdoor, port: 5001, protocol: TCP} ]`)),
				}

				results, err := dyff.CompareInputFiles(file(assets("kubernetes/multi-docs/from.yml")), file(assets("kubernetes/multi-docs/to.yml")))
				Expect(err).To(BeNil())
				Expect(results).NotTo(BeNil())
				Expect(results.Diffs).NotTo(BeNil())
				Expect(results.Diffs).To(HaveLen(len(expected)))

				for i, result := range results.Diffs {
					Expect(result).To(BeSameDiffAs(expected[i]))
				}
			})

			It("should fail to compare files with different number of documents", func() {
				from := ytbx.InputFile{Location: "/ginkgo/compare/test/from", Documents: multiDoc("foo: bar", "dead: beef")}
				to := ytbx.InputFile{Location: "/ginkgo/compare/test/to", Documents: multiDoc("bar: foo")}

				_, err := dyff.CompareInputFiles(from, to)
				Expect(err).To(HaveOccurred())
			})

			It("should return differences in named lists even if no standard identifier is used", func() {
				results, err := dyff.CompareInputFiles(
					file(assets("prometheus/from.yml")),
					file(assets("prometheus/to.yml")),
				)

				Expect(err).To(BeNil())
				Expect(results).NotTo(BeNil())
				Expect(results.Diffs).NotTo(BeNil())

				expected := []dyff.Diff{
					singleDiff("/scrape_configs", dyff.ORDERCHANGE,
						[]string{
							"kubernetes-nodes",
							"kubernetes-apiservers",
							"kubernetes-cadvisor",
							"kubernetes-service-endpoints",
							"kubernetes-services",
							"kubernetes-ingresses",
							"kubernetes-pods",
						},
						[]string{
							"kubernetes-apiservers",
							"kubernetes-nodes",
							"kubernetes-cadvisor",
							"kubernetes-service-endpoints",
							"kubernetes-services",
							"kubernetes-ingresses",
							"kubernetes-pods",
						},
					),

					singleDiff("/scrape_configs/job_name=kubernetes-apiservers/scheme", dyff.MODIFICATION,
						"http",
						"https",
					),

					singleDiff("/scrape_configs/job_name=kubernetes-apiservers/relabel_configs/0/regex", dyff.MODIFICATION,
						"default;kubernetes;http",
						"default;kubernetes;https",
					),
				}

				Expect(results.Diffs).To(HaveLen(len(expected)))
				for i, diff := range results.Diffs {
					Expect(diff).To(BeSameDiffAs(expected[i]))
				}
			})

			It("should fail to find the non-standard identifier if the threshold is too high", func() {
				report, err := dyff.CompareInputFiles(
					file(assets("prometheus/from.yml")),
					file(assets("prometheus/to.yml")),
					dyff.NonStandardIdentifierGuessCountThreshold(8),
				)

				Expect(err).To(BeNil())
				Expect(report).NotTo(BeNil())
				Expect(report.Diffs).NotTo(BeNil())

				var orderChangeDiffs = 0
				for _, diff := range report.Diffs {
					for _, detail := range diff.Details {
						if detail.Kind == dyff.ORDERCHANGE {
							orderChangeDiffs++
						}
					}
				}

				Expect(orderChangeDiffs).To(BeEquivalentTo(0))
			})

			It("should filter my report based on set of paths", func() {
				pathString := "/yaml/map/foobar"

				report := dyff.Report{Diffs: []dyff.Diff{
					singleDiff(pathString, dyff.ADDITION, nil, "foobar"),
					singleDiff("/yaml/map/barfoo", dyff.ADDITION, nil, "barfoo"),
				}}

				Expect(report.Filter()).To(BeEquivalentTo(report))
				Expect(report.Filter(pathString)).To(BeEquivalentTo(dyff.Report{Diffs: []dyff.Diff{
					singleDiff(pathString, dyff.ADDITION, nil, "foobar"),
				}}))

				Expect(report.Filter("/does/not/exist")).To(BeEquivalentTo(dyff.Report{}))
			})

			It("should filter my report based on set of regular expressions", func() {
				pathString := "/yaml/map/foobar"

				report := dyff.Report{Diffs: []dyff.Diff{
					singleDiff(pathString, dyff.ADDITION, nil, "foobar"),
					singleDiff("/yaml/map/barfoo", dyff.ADDITION, nil, "barfoo"),
				}}

				Expect(report.FilterRegexp()).To(BeEquivalentTo(report))
				Expect(report.FilterRegexp("foobar")).To(BeEquivalentTo(dyff.Report{Diffs: []dyff.Diff{
					singleDiff(pathString, dyff.ADDITION, nil, "foobar"),
				}}))

				Expect(report.FilterRegexp("/does/not/exist")).To(BeEquivalentTo(dyff.Report{}))
			})

			It("should exclude my report based on regular expressions", func() {
				pathString := "/yaml/map/foobar"

				report := dyff.Report{Diffs: []dyff.Diff{
					singleDiff(pathString, dyff.ADDITION, nil, "foobar"),
					singleDiff("/yaml/map/barfoo", dyff.ADDITION, nil, "barfoo"),
				}}

				Expect(report.ExcludeRegexp()).To(BeEquivalentTo(report))
				Expect(report.ExcludeRegexp("barfoo")).To(BeEquivalentTo(dyff.Report{Diffs: []dyff.Diff{
					singleDiff(pathString, dyff.ADDITION, nil, "foobar"),
				}}))

				Expect(report.ExcludeRegexp("/does/not/exist")).To(BeEquivalentTo(report))
			})

			It("should ignore changes in values", func() {
				report := dyff.Report{Diffs: []dyff.Diff{
					singleDiff("/yaml/map/add", dyff.ADDITION, nil, "added"),
					singleDiff("/yaml/map/removed", dyff.REMOVAL, nil, "removed"),
					singleDiff("/yaml/map/changed", dyff.MODIFICATION, "foobar", "barfoo"),
				}}

				Expect(report.IgnoreValueChanges()).To(BeEquivalentTo(dyff.Report{Diffs: []dyff.Diff{
					singleDiff("/yaml/map/add", dyff.ADDITION, nil, "added"),
					singleDiff("/yaml/map/removed", dyff.REMOVAL, nil, "removed"),
				}}))
			})
		})

		Context("change root for comparison", func() {
			It("should change the root of an input file", func() {
				from := ytbx.InputFile{Location: "/ginkgo/compare/test/from", Documents: multiDoc(`---
a: foo
---
b: bar
`)}

				to := ytbx.InputFile{Location: "/ginkgo/compare/test/to", Documents: multiDoc(`{
"items": [
  {"a": "Foo"},
  {"b": "Bar"}
]
}`)}

				err := dyff.ChangeRoot(&to, "/items", false, true)
				if err != nil {
					Fail(err.Error())
				}

				results, err := dyff.CompareInputFiles(from, to)
				Expect(err).To(BeNil())

				expected := []dyff.Diff{
					singleDiff("#0/a", dyff.MODIFICATION, "foo", "Foo"),
					singleDiff("#1/b", dyff.MODIFICATION, "bar", "Bar"),
				}

				Expect(results.Diffs).To(HaveLen(len(expected)))
				for i, result := range results.Diffs {
					Expect(result).To(BeSameDiffAs(expected[i]))
				}
			})
		})

		Context("two YAML structures with Kubernetes lists", func() {
			It("should identify individual list entries based on the nested name field in the respective entry metadata", func() {
				from, to := loadFiles(
					assets("kubernetes", "lists", "from.yml"),
					assets("kubernetes", "lists", "to.yml"),
				)

				results, err := dyff.CompareInputFiles(from, to, dyff.KubernetesEntityDetection(true))
				Expect(err).ToNot(HaveOccurred())

				Expect(results.Diffs).To(HaveLen(2))
				Expect(results.Diffs[0]).To(BeSameDiffAs(singleDiff(
					"#0/items",
					dyff.ORDERCHANGE,
					dyff.AsSequenceNode("v1/Pod/foobar/foo-2", "v1/Pod/foobar/foo-1"),
					dyff.AsSequenceNode("v1/Pod/foobar/foo-1", "v1/Pod/foobar/foo-2"),
				)))
				Expect(results.Diffs[1]).To(BeSameDiffAs(singleDiff(
					"/items/resource=v1\\/Pod\\/foobar\\/foo-1/metadata/labels/foo",
					dyff.MODIFICATION,
					"bAr",
					"bar",
				)))
			})

			It("should detect Kubernetes items by their apiVersion, kind, and name", func() {
				from, to := loadFiles(
					assets("kubernetes", "lists-versioned-crds", "from.yml"),
					assets("kubernetes", "lists-versioned-crds", "to.yml"),
				)

				results, err := dyff.CompareInputFiles(from, to, dyff.KubernetesEntityDetection(true))
				Expect(err).ToNot(HaveOccurred())

				Expect(results.Diffs).To(HaveLen(1))
				Expect(results.Diffs[0]).To(BeSameDiffAs(singleDiff(
					"#0/items",
					dyff.ORDERCHANGE,
					dyff.AsSequenceNode("v2/kindv2/item", "v1/kindv1/item"),
					dyff.AsSequenceNode("v1/kindv1/item", "v2/kindv2/item"),
				)))
			})
		})

		Context("checking known issues of compare", func() {
			It("should not return order change differences in case the named-entry list does not have unique identifiers", func() {
				from, to, err := ytbx.LoadFiles("../../assets/issues/issue-38/from.yml", "../../assets/issues/issue-38/to.yml")
				Expect(err).To(BeNil())
				Expect(from).ToNot(BeNil())
				Expect(to).ToNot(BeNil())

				results, err := dyff.CompareInputFiles(from, to)
				Expect(err).ToNot(HaveOccurred())
				Expect(results).ToNot(BeNil())

				Expect(results.Diffs).To(HaveLen(0))
			})

			It("should process files with YAML anchors", func() {
				from, to, err := ytbx.LoadFiles("../../assets/issues/issue-76/from.yml", "../../assets/issues/issue-76/to.yml")
				Expect(err).To(BeNil())
				Expect(from).ToNot(BeNil())
				Expect(to).ToNot(BeNil())

				results, err := dyff.CompareInputFiles(from, to)
				Expect(err).ToNot(HaveOccurred())
				Expect(results).ToNot(BeNil())

				Expect(results.Diffs).To(HaveLen(2))
			})

			It("should treat the string content as-is with no evaluation", func() {
				from, to, err := ytbx.LoadFiles("../../assets/issues/issue-132/from.yml", "../../assets/issues/issue-132/to.yml")
				Expect(err).To(BeNil())
				Expect(from).ToNot(BeNil())
				Expect(to).ToNot(BeNil())

				results, err := dyff.CompareInputFiles(from, to)
				Expect(err).ToNot(HaveOccurred())
				Expect(results).ToNot(BeNil())

				Expect(results.Diffs[0].Details[0].From.Value).To(Equal("%{one}"))
				Expect(results.Diffs[1].Details[0].To.Value).To(Equal("%{two}"))
			})

			It("should detect change in simple list in case a duplicate entry is introduced or removed", func() {
				from, to, err := ytbx.LoadFiles("../../assets/issues/issue-143/from.json", "../../assets/issues/issue-143/to.json")
				Expect(err).To(BeNil())
				Expect(from).ToNot(BeNil())
				Expect(to).ToNot(BeNil())

				results, err := dyff.CompareInputFiles(from, to)
				Expect(err).ToNot(HaveOccurred())
				Expect(results).ToNot(BeNil())
				Expect(results.Diffs).ToNot(HaveLen(0))
			})

			It("should detect order changes in simple lists with duplicate entries", func() {
				from := ytbx.InputFile{Location: "/ginkgo/compare/test/from", Documents: multiDoc(`{ "keys": [ "value1", "value2", "value1", "value2" ] }`)}
				to := ytbx.InputFile{Location: "/ginkgo/compare/test/to", Documents: multiDoc(`{ "keys": [ "value1", "value1", "value2", "value2" ] }`)}

				results, err := dyff.CompareInputFiles(from, to)
				Expect(err).ToNot(HaveOccurred())
				Expect(results).ToNot(BeNil())

				Expect(results.Diffs).To(HaveLen(1))
				Expect(results.Diffs[0]).To(BeSameDiffAs(singleDiff(
					"/keys",
					dyff.ORDERCHANGE,
					dyff.AsSequenceNode("value1", "value2", "value1", "value2"),
					dyff.AsSequenceNode("value1", "value1", "value2", "value2"),
				)))
			})

			It("should ignore order changes when comparing list entries that are lists itself", func() {
				from, to, err := ytbx.LoadFiles("../../assets/issues/issue-125/from.json", "../../assets/issues/issue-125/to.json")
				Expect(err).To(BeNil())
				Expect(from).ToNot(BeNil())
				Expect(to).ToNot(BeNil())

				results, err := dyff.CompareInputFiles(from, to, dyff.IgnoreOrderChanges(false))
				Expect(err).ToNot(HaveOccurred())
				Expect(results).ToNot(BeNil())
				Expect(results.Diffs).ToNot(HaveLen(0))

				results, err = dyff.CompareInputFiles(from, to, dyff.IgnoreOrderChanges(true))
				Expect(err).ToNot(HaveOccurred())
				Expect(results).ToNot(BeNil())
				Expect(results.Diffs).To(HaveLen(0))
			})
		})

		Context("input files containing Kubernetes resources", func() {
			It("should return order change differences in YAML files with Kubernetes resources", func() {
				from, to, err := ytbx.LoadFiles(assets("issues", "issue-184", "from.yml"), assets("issues", "issue-184", "to.yml"))
				Expect(err).To(BeNil())
				Expect(from).ToNot(BeNil())
				Expect(to).ToNot(BeNil())

				results, err := dyff.CompareInputFiles(from, to)
				Expect(err).ToNot(HaveOccurred())
				Expect(results).ToNot(BeNil())

				Expect(results.Diffs).To(HaveLen(1))
			})

			It("should report that a document was added", func() {
				from := ytbx.InputFile{
					Location: "/ginkgo/compare/test/from",
					Documents: multiDoc(
						`{"apiVersion": "apps/v1", "kind": "Deployment", "metadata": {"name": "x"}}`,
					),
				}

				to := ytbx.InputFile{
					Location: "/ginkgo/compare/test/to",
					Documents: multiDoc(
						`{"apiVersion": "apps/v1", "kind": "Deployment", "metadata": {"name": "x"}}`,
						`{"apiVersion": "v1", "kind": "Service", "metadata": {"name": "y"}}`,
					),
				}

				result, err := dyff.CompareInputFiles(from, to)
				Expect(err).NotTo(HaveOccurred())

				Expect(result.Diffs).To(HaveLen(1))
				Expect(result.Diffs[0].Details).To(HaveLen(1))
				Expect(result.Diffs[0].Details[0].Kind).To(Equal(dyff.ADDITION))
			})

			It("should report that a document was removed", func() {
				from := ytbx.InputFile{
					Location: "/ginkgo/compare/test/from",
					Documents: multiDoc(
						`{"apiVersion": "apps/v1", "kind": "Deployment", "metadata": {"name": "x"}}`,
						`{"apiVersion": "v1", "kind": "Service", "metadata": {"name": "y"}}`,
					),
				}

				to := ytbx.InputFile{
					Location: "/ginkgo/compare/test/to",
					Documents: multiDoc(
						`{"apiVersion": "apps/v1", "kind": "Deployment", "metadata": {"name": "x"}}`,
					),
				}

				result, err := dyff.CompareInputFiles(from, to)
				Expect(err).NotTo(HaveOccurred())

				Expect(result.Diffs).To(HaveLen(1))
				Expect(result.Diffs[0].Details).To(HaveLen(1))
				Expect(result.Diffs[0].Details[0].Kind).To(Equal(dyff.REMOVAL))
			})

			It("should omit nil/empty documents", func() {
				from := ytbx.InputFile{
					Location: "/ginkgo/compare/test/from",
					Documents: multiDoc(
						`{"apiVersion": "apps/v1", "kind": "Deployment", "metadata": {"name": "x"}}`,
						`{"apiVersion": "v1", "kind": "Service", "metadata": {"name": "y"}}`,
					),
				}

				to := ytbx.InputFile{
					Location: "/ginkgo/compare/test/to",
					Documents: multiDoc(
						`{"apiVersion": "apps/v1", "kind": "Deployment", "metadata": {"name": "x"}}`,
						``,
						`{"apiVersion": "v1", "kind": "Service", "metadata": {"name": "y"}}`,
						``,
					),
				}

				result, err := dyff.CompareInputFiles(from, to)
				Expect(err).NotTo(HaveOccurred())
				Expect(result.Diffs).To(HaveLen(0))
			})
		})

		Context("input files containing complex objects with custom keys", func() {
			It("cannot determine the keys through heuristics", func() {
				from, to, err := ytbx.LoadFiles(assets("issues", "issue-243", "to.yml"), assets("issues", "issue-243", "from.yml"))
				Expect(err).To(BeNil())
				Expect(from).ToNot(BeNil())
				Expect(to).ToNot(BeNil())

				results, err := dyff.CompareInputFiles(from, to, dyff.IgnoreOrderChanges(true))
				Expect(err).ToNot(HaveOccurred())
				Expect(results).ToNot(BeNil())
				Expect(results.Diffs).ToNot(HaveLen(0))
			})

			It("accurately reports no differences when keys are given", func() {
				from, to, err := ytbx.LoadFiles(assets("issues", "issue-243", "to.yml"), assets("issues", "issue-243", "from.yml"))
				Expect(err).To(BeNil())
				Expect(from).ToNot(BeNil())
				Expect(to).ToNot(BeNil())

				results, err := dyff.CompareInputFiles(from, to, dyff.IgnoreOrderChanges(true), dyff.AdditionalIdentifiers("branch"))
				Expect(err).NotTo(HaveOccurred())
				Expect(results).ToNot(BeNil())
				Expect(results.Diffs).To(HaveLen(0))
			})
		})
	})
})
