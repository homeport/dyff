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

## Examples
- Show the differences between two versions of [`cf-deployment`](https://github.com/cloudfoundry/cf-deployment/):
    ```
    dyff between \
      https://raw.githubusercontent.com/cloudfoundry/cf-deployment/v1.19.0/cf-deployment.yml \
      https://raw.githubusercontent.com/cloudfoundry/cf-deployment/v1.20.0/cf-deployment.yml
    ```

- Convert a YAML file to JSON
    ```
    dyff json https://raw.githubusercontent.com/cloudfoundry/cf-deployment/v1.19.0/cf-deployment.yml
    ```

- Convert a JSON file to YAML
    ```
    dyff yaml https://raw.githubusercontent.com/HeavyWombat/dyff/develop/assets/bosh-yaml/manifest.json
    ```
