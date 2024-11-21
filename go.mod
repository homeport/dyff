module github.com/homeport/dyff

go 1.22.0

toolchain go1.23.2

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/gonvenience/bunt v1.3.5
	github.com/gonvenience/neat v1.3.13
	github.com/gonvenience/term v1.0.2
	github.com/gonvenience/text v1.0.7
	github.com/gonvenience/wrap v1.2.0
	github.com/gonvenience/ytbx v1.4.4
	github.com/lucasb-eyer/go-colorful v1.2.0
	github.com/mitchellh/hashstructure v1.1.0
	github.com/onsi/ginkgo/v2 v2.20.1
	github.com/onsi/gomega v1.34.1
	github.com/spf13/cobra v1.8.1
	github.com/texttheater/golang-levenshtein v1.0.1
	gopkg.in/yaml.v3 v3.0.1
)

// usage untagged version of go-diff
// cause https://github.com/sergi/go-diff/issues/123
// fixed in https://github.com/sergi/go-diff/pull/136
// but currently not tagged
require github.com/sergi/go-diff v1.3.2-0.20230802210424-5b0b94c5c0d3

require (
	github.com/BurntSushi/toml v1.4.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-task/slim-sprig/v3 v3.0.0 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/pprof v0.0.0-20240727154555-813a5fbdbec8 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/mattn/go-ciede2000 v0.0.0-20170301095244-782e8c62fec3 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mitchellh/go-ps v1.0.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/virtuald/go-ordered-json v0.0.0-20170621173500-b18e6e673d74 // indirect
	golang.org/x/exp v0.0.0-20240719175910-8a7402abbf56 // indirect
	golang.org/x/net v0.29.0 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
	golang.org/x/term v0.24.0 // indirect
	golang.org/x/text v0.18.0 // indirect
	golang.org/x/tools v0.25.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
