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

				result := CompareDocuments(from, to)
				Expect(result).NotTo(BeNil())
				Expect(len(result)).To(BeEquivalentTo(1))

				Expect(result[0]).To(BeEquivalentTo(Diff{
					Kind: ADDITION,
					Path: path("/some/yaml/structure"),
					From: nil,
					To:   yml(`version: v1`)}))
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

				result := CompareDocuments(from, to)
				Expect(result).NotTo(BeNil())
				Expect(len(result)).To(BeEquivalentTo(1))

				Expect(result[0]).To(BeEquivalentTo(Diff{
					Kind: REMOVAL,
					Path: path("/some/yaml/structure"),
					From: yml(`version: v1`),
					To:   nil}))
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

				result := CompareDocuments(from, to)
				Expect(result).NotTo(BeNil())
				Expect(len(result)).To(BeEquivalentTo(2))

				Expect(result[0]).To(BeEquivalentTo(Diff{
					Kind: REMOVAL,
					Path: path("/some/yaml/structure"),
					From: yml(`version: v1`),
					To:   nil}))

				Expect(result[1]).To(BeEquivalentTo(Diff{
					Kind: ADDITION,
					Path: path("/some/yaml/structure"),
					From: nil,
					To:   yml("release: v1")}))
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

				result := CompareDocuments(from, to)
				Expect(result).NotTo(BeNil())
				Expect(len(result)).To(BeEquivalentTo(1))

				Expect(result[0]).To(BeEquivalentTo(Diff{
					Kind: ADDITION,
					Path: path("/some/yaml/structure/list"),
					To:   yml(`list: [ three ]`)[0].Value}))
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

				result := CompareDocuments(from, to)
				Expect(result).NotTo(BeNil())
				Expect(len(result)).To(BeEquivalentTo(1))

				Expect(result[0]).To(BeEquivalentTo(Diff{
					Kind: ADDITION,
					Path: path("/some/yaml/structure/list"),
					To:   yml(`list: [ 3 ]`)[0].Value}))
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

				result := CompareDocuments(from, to)
				Expect(result).NotTo(BeNil())
				Expect(len(result)).To(BeEquivalentTo(1))

				Expect(result[0]).To(BeEquivalentTo(Diff{
					Kind: REMOVAL,
					Path: path("/some/yaml/structure/list"),
					From: yml(`list: [ three ]`)[0].Value}))
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

				result := CompareDocuments(from, to)
				Expect(result).NotTo(BeNil())
				Expect(len(result)).To(BeEquivalentTo(1))

				Expect(result[0]).To(BeEquivalentTo(Diff{
					Kind: REMOVAL,
					Path: path("/some/yaml/structure/list"),
					From: yml(`list: [ 3 ]`)[0].Value}))
			})
		})

		Context("Given two YAML structures with complext content", func() {
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

				result := CompareDocuments(from, to)
				Expect(result).NotTo(BeNil())
				Expect(len(result)).To(BeEquivalentTo(7))

				Expect(result[0].Kind).To(BeIdenticalTo(MODIFICATION))
				Expect(result[0].Path.String()).To(BeIdenticalTo("/instance_groups/name=web/networks/name=concourse/static_ips"))
				Expect(result[0].From).To(BeIdenticalTo("192.168.1.1"))
				Expect(result[0].To).To(BeIdenticalTo("192.168.0.1"))

				Expect(result[1].Kind).To(BeIdenticalTo(MODIFICATION))
				Expect(result[1].Path.String()).To(BeIdenticalTo("/instance_groups/name=web/jobs/name=atc/properties/external_url"))
				Expect(result[1].From).To(BeIdenticalTo("http://192.168.1.100:8080"))
				Expect(result[1].To).To(BeIdenticalTo("http://192.168.0.100:8080"))

				Expect(result[2].Kind).To(BeIdenticalTo(MODIFICATION))
				Expect(result[2].Path.String()).To(BeIdenticalTo("/instance_groups/name=web/jobs/name=atc/properties/development_mode"))
				Expect(result[2].From).To(BeIdenticalTo(true))
				Expect(result[2].To).To(BeIdenticalTo(false))

				Expect(result[3].Kind).To(BeIdenticalTo(MODIFICATION))
				Expect(result[3].Path.String()).To(BeIdenticalTo("/instance_groups/name=web/jobs/name=db/instances"))
				Expect(result[3].From).To(BeIdenticalTo(1))
				Expect(result[3].To).To(BeIdenticalTo(2))

				Expect(result[4]).To(BeEquivalentTo(Diff{
					Kind: REMOVAL,
					Path: path("/instance_groups/name=web/jobs/name=db/networks"),
					From: yml(`list: [ { name: testnet } ]`)[0].Value}))

				Expect(result[5].Kind).To(BeIdenticalTo(MODIFICATION))
				Expect(result[5].Path.String()).To(BeIdenticalTo("/instance_groups/name=web/jobs/name=db/jobs/name=postgresql/properties/databases/name=atc/password"))
				Expect(result[5].From).To(BeIdenticalTo("supersecret"))
				Expect(result[5].To).To(BeIdenticalTo("zwX#(;P=%hTfFzM["))

				Expect(result[6]).To(BeEquivalentTo(Diff{
					Kind: ADDITION,
					Path: path("/instance_groups/name=web/jobs"),
					To:   yml(`list: [ { release: custom, name: logger } ]`)[0].Value}))
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

				result := CompareDocuments(from, to)
				Expect(result).NotTo(BeNil())

				Expect(len(result)).To(BeEquivalentTo(1))

				Expect(result[0].Kind).To(BeIdenticalTo(MODIFICATION))
				Expect(result[0].Path.String()).To(BeIdenticalTo("/resource_pools/name=concourse_resource_pool/cloud_properties/datacenters/0/clusters/0/CLS_PAAS_SFT_035/resource_pool"))
				Expect(result[0].From).To(BeIdenticalTo("other-vsphere-res-pool"))
				Expect(result[0].To).To(BeIdenticalTo("new-vsphere-res-pool"))
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

				result := CompareDocuments(from, to)
				Expect(result).NotTo(BeNil())

				Expect(len(result)).To(BeEquivalentTo(2))

				Expect(result[0]).To(BeEquivalentTo(Diff{
					Kind: REMOVAL,
					Path: path("/resource_pools/name=concourse_resource_pool/cloud_properties/datacenters/0/clusters"),
					From: yml(`list: [ {CLS_PAAS_SFT_035: {resource_pool: 35-vsphere-res-pool}}, {CLS_PAAS_SFT_036: {resource_pool: 36-vsphere-res-pool}} ]`)[0].Value}))

				Expect(result[1]).To(BeEquivalentTo(Diff{
					Kind: ADDITION,
					Path: path("/resource_pools/name=concourse_resource_pool/cloud_properties/datacenters/0/clusters"),
					To:   yml(`list: [ {CLS_PAAS_SFT_035: {resource_pool: 35a-vsphere-res-pool}}, {CLS_PAAS_SFT_036: {resource_pool: 36a-vsphere-res-pool}} ]`)[0].Value}))
			})
		})

		Context("Given two YAML files", func() {
			It("should return all differences in there", func() {
				results := CompareDocuments(yml("../assets/examples/from.yml"), yml("../assets/examples/to.yml"))
				expected := []Diff{
					Diff{
						Kind: REMOVAL,
						Path: path("/yaml/map"),
						From: yml(`---
stringB: fOObAr
intB: 10
floatB: 2.71
boolB: false
mapB: { key0: B, key1: B }
listB: [ B, B, B ]
`)},

					Diff{
						Kind: ADDITION,
						Path: path("/yaml/map"),
						To: yml(`---
stringY: YAML!
intY: 147
floatY: 24.0
boolY: true
mapY: { key0: Y, key1: Y }
listY: [ Y, Y, Y ]
`)},

					Diff{
						Kind: REMOVAL,
						Path: path("/yaml/simple-list"),
						From: yml(`list: [ X, Z ]`)[0].Value},

					Diff{
						Kind: ADDITION,
						Path: path("/yaml/simple-list"),
						To:   yml(`list: [ D, E ]`)[0].Value},

					Diff{
						Kind: REMOVAL,
						Path: path("/yaml/named-entry-list-using-name"),
						From: yml(`list: [ {name: X}, {name: Z} ]`)[0].Value},

					Diff{
						Kind: ADDITION,
						Path: path("/yaml/named-entry-list-using-name"),
						To:   yml(`list: [ {name: D}, {name: E} ]`)[0].Value},

					Diff{
						Kind: REMOVAL,
						Path: path("/yaml/named-entry-list-using-key"),
						From: yml(`list: [ {key: X}, {key: Z} ]`)[0].Value},

					Diff{
						Kind: ADDITION,
						Path: path("/yaml/named-entry-list-using-key"),
						To:   yml(`list: [ {key: D}, {key: E} ]`)[0].Value},

					Diff{
						Kind: REMOVAL,
						Path: path("/yaml/named-entry-list-using-id"),
						From: yml(`list: [ {id: X}, {id: Z} ]`)[0].Value},

					Diff{
						Kind: ADDITION,
						Path: path("/yaml/named-entry-list-using-id"),
						To:   yml(`list: [ {id: D}, {id: E} ]`)[0].Value},
				}

				Expect(results).NotTo(BeNil())
				Expect(len(results)).To(BeEquivalentTo(len(expected)))

				for i, result := range results {
					Expect(result).To(BeEquivalentTo(expected[i]))
				}
			})
		})
	})
})
