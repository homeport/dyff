[![License](https://img.shields.io/github/license/HeavyWombat/dyff.svg)](https://github.com/HeavyWombat/dyff/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/HeavyWombat/dyff)](https://goreportcard.com/report/github.com/HeavyWombat/dyff)
[![Build Status](https://travis-ci.org/HeavyWombat/dyff.svg?branch=master)](https://travis-ci.org/HeavyWombat/dyff)
[![GoDoc](https://godoc.org/github.com/HeavyWombat/dyff?status.svg)](https://godoc.org/github.com/HeavyWombat/dyff)
[![Release](https://img.shields.io/github/release/HeavyWombat/dyff.svg)](https://github.com/HeavyWombat/dyff/releases/latest)

# δyƒƒ /ˈdʏf/
A diff tool for YAML files, and sometimes JSON

![dyff between example](docs/images/dyff-between-example.png?raw=true "dyff between example of two cf-deployment versions")

## Objectives
- Easy to use **command line interface** with natural language design
- **Human friendly output** through emphasising structures and content specific colours
- Very **compact** and easy to understand differences report
- Emulate **BOSH v1** deployment delta output look and feel
- JSON to YAML and vice versa conversion that **preserves the order** of entries in maps (hashes)
- Usage of **Go** to actually learn the language (and have simple to distribute binaries)

## Description
The `dyff` tool will only show the differences between the files by displaying by their path (location inside the YAML tree) and the respective change. The input location can either be a local file or a URI. Also, reading from STDIN is supported by using `-` as the filename. The default report style is loosely based upon the old [BOSH v1](https://bosh.io/) deployment delta output. The path style is based on the Dot-style that is used by [Spruce](https://github.com/geofffranks/spruce). As an alternative, you can specify to display paths by using the [go-patch](https://github.com/cppforlife/go-patch) style. Have a look in the help section of the respective subcommand to see more options.

As a side effect, `dyff` can convert YAML to JSON and vice versa while preserving the order of keys in hashes from the input document.

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
