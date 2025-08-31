package main

import (
	"fmt"
	"log"

	"github.com/gonvenience/ytbx"
	yamlv3 "gopkg.in/yaml.v3"

	"github.com/homeport/dyff/pkg/dyff"
)

func yml(input string) *yamlv3.Node {
	var node yamlv3.Node
	if err := yamlv3.Unmarshal([]byte(input), &node); err != nil {
		log.Fatal(err)
	}
	return node.Content[0]
}

func main() {
	fromYAML := `---
files:
  simple:
    content: "test"
  newline:
    content: "test"
  complex:
    content: "test"
`
	
	toYAML := `---
files:
  simple:
    content: "modified"
  newline:
    content: "modified"
  complex:
    content: "modified"
`
	
	from := ytbx.InputFile{
		Location:  "from.yml",
		Documents: []*yamlv3.Node{yml(fromYAML)},
	}
	
	to := ytbx.InputFile{
		Location:  "to.yml", 
		Documents: []*yamlv3.Node{yml(toYAML)},
	}
	
	report, err := dyff.CompareInputFiles(from, to)
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("Found %d diffs:\n", len(report.Diffs))
	for i, diff := range report.Diffs {
		pathStr := ""
		if diff.Path != nil {
			sections := []string{}
			for _, element := range diff.Path.PathElements {
				switch {
				case element.Key == "" && element.Name != "":
					sections = append(sections, element.Name)
				case element.Key != "" && element.Name != "":
					sections = append(sections, element.Name)
				case element.Idx >= 0:
					sections = append(sections, fmt.Sprintf("%d", element.Idx))
				default:
					sections = append(sections, element.Key)
				}
			}
			pathStr = fmt.Sprintf("files.%s.content", sections[len(sections)-2])
		}
		fmt.Printf("Diff %d: %s\n", i, pathStr)
	}
}
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

	result, err := dyff.CompareInputFiles(
		ytbx.InputFile{Documents: []*yamlv3.Node{from}},
		ytbx.InputFile{Documents: []*yamlv3.Node{to}},
		dyff.KubernetesEntityDetection(false),
		dyff.DetailedListDiff(true),
	)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Number of diffs: %d\n", len(result.Diffs))
	for i, diff := range result.Diffs {
		fmt.Printf("Diff %d: %s (details: %d)\n", i, diff.Path.String(), len(diff.Details))
		for j, detail := range diff.Details {
			fmt.Printf("  Detail %d: %c\n", j, detail.Kind)
		}
	}
}
