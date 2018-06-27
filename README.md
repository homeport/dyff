[![License](https://img.shields.io/github/license/HeavyWombat/dyff.svg)](https://github.com/HeavyWombat/dyff/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/HeavyWombat/dyff)](https://goreportcard.com/report/github.com/HeavyWombat/dyff)
[![Build Status](https://travis-ci.org/HeavyWombat/dyff.svg?branch=master)](https://travis-ci.org/HeavyWombat/dyff)
[![GoDoc](https://godoc.org/github.com/HeavyWombat/dyff/pkg?status.svg)](https://godoc.org/github.com/HeavyWombat/dyff/pkg)
[![Release](https://img.shields.io/github/release/HeavyWombat/dyff.svg)](https://github.com/HeavyWombat/dyff/releases/latest)

# δyƒƒ /ˈdʏf/
A diff tool for YAML files, and sometimes JSON

![dyff between example](docs/images/dyff-between-example.png?raw=true "dyff between example of two cf-deployment versions")

## Description
`dyff` is inspired by the way the old [BOSH v1](https://bosh.io/) deployment output reported changes from one version to another by only showing the parts of a YAML file that change.

Each difference is referenced by its location in the YAML document by using either the [Spruce](https://github.com/geofffranks/spruce) or [go-patch](https://github.com/cppforlife/go-patch) path syntax. The output report aims to be as compact as possible to give a clear and simple overview of the change.

Similar to the standard `diff` tool, it follows the principle of describing the change by going from the `from` input file to the target `to` input file.

Input files can be local files (filesystem path), remote files (URI), or the standard input stream (using `-`).

All orders of keys in hashes are preserved during processing and output to the terminal, most notably in the sub-commands to convert YAML to JSON and vice versa.

## Installation
On macOS, `dyff` is available via Homebrew:
```bash
brew tap HeavyWombat/tap
brew install dyff
```

Prebuilt binaries for a lot of operating systems and architectures can be [downloaded from the releases section](https://github.com/HeavyWombat/dyff/releases/latest).

There is a convenience script to download the latest release for Linux or macOS if you want to keep it simple (you need `curl` and `jq` installed on your machine):
```bash
curl --silent --location https://goo.gl/DRXDVN | bash
```

And of course, you can download and build `dyff` from source using `go`:
```bash
go get github.com/HeavyWombat/dyff/...
```

## Use cases and examples
- Show the differences between two versions of [`cf-deployment`](https://github.com/cloudfoundry/cf-deployment/) YAMLs:
    ```bash
    dyff between \
      https://raw.githubusercontent.com/cloudfoundry/cf-deployment/v1.19.0/cf-deployment.yml \
      https://raw.githubusercontent.com/cloudfoundry/cf-deployment/v1.20.0/cf-deployment.yml
    ```

- Convert a JSON stream to YAML
    ```bash
    sometool --json | jq --raw-output '.data' | dyff yaml -
    ```

- Sometimes you end up with YAML or JSON files, where the order of the keys in maps was sorted alphabetically. With `dyff` you can restructure keys in maps to a more human appealing order:
    ```bash
    sometool --export --json | dyff - yaml --restructure
    ```
    Or, rewrite a file _in place_ with the restructured order of keys.
    ```bash
    dyff yaml --restructure --in-place somefile.yml
    ```

- Just print a YAML (or JSON) file to the terminal to look at it. By default, `dyff` will use a neat output schema which includes different colors and indent helper lines to improve readability. The colors are roughly based on the default [Atom](https://atom.io) schema and work best on dark terminal backgrounds. The neat output is disabled the output of `dyff` is redirected into a pipe, or you can disable it explicitly using the `--plain` flag.
    ```bash
    dyff yaml somefile.yml
    ```

- Convert a YAML file to JSON and vice versa:
    ```bash
    dyff json https://raw.githubusercontent.com/cloudfoundry/cf-deployment/v1.19.0/cf-deployment.yml
    ```
    The `dyff` sub-command (`yaml`, or `json`) defines the output format, the tool automatically detects the input format itself.
    ```bash
    dyff yaml https://raw.githubusercontent.com/HeavyWombat/dyff/develop/assets/bosh-yaml/manifest.json
    ```
