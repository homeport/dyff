module github.com/homeport/dyff

go 1.23.0

toolchain go1.24.1

require (
	github.com/gonvenience/bunt v1.4.2
	github.com/gonvenience/idem v0.0.2
	github.com/gonvenience/neat v1.3.16
	github.com/gonvenience/term v1.0.4
	github.com/gonvenience/text v1.0.9
	github.com/gonvenience/ytbx v1.4.7
	github.com/lucasb-eyer/go-colorful v1.3.0
	github.com/mitchellh/hashstructure v1.1.0
	github.com/onsi/ginkgo/v2 v2.27.1
	github.com/onsi/gomega v1.38.2
	github.com/spf13/cobra v1.10.1
	github.com/texttheater/golang-levenshtein v1.0.1
	gopkg.in/yaml.v3 v3.0.1
)

// usage untagged version of go-diff
// cause https://github.com/sergi/go-diff/issues/123
// fixed in https://github.com/sergi/go-diff/pull/136
// but currently not tagged
require github.com/sergi/go-diff v1.4.0

require (
	github.com/BurntSushi/toml v1.5.0 // indirect
	github.com/Masterminds/semver/v3 v3.4.0 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.7 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-task/slim-sprig/v3 v3.0.0 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/pprof v0.0.0-20250403155104-27863c87afa6 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/mattn/go-ciede2000 v0.0.0-20170301095244-782e8c62fec3 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mitchellh/go-ps v1.0.0 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	github.com/stretchr/testify v1.11.1 // indirect
	github.com/virtuald/go-ordered-json v0.0.0-20170621173500-b18e6e673d74 // indirect
	go.uber.org/automaxprocs v1.6.0 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/mod v0.27.0 // indirect
	golang.org/x/net v0.43.0 // indirect
	golang.org/x/sync v0.16.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
	golang.org/x/term v0.34.0 // indirect
	golang.org/x/text v0.28.0 // indirect
	golang.org/x/tools v0.36.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
