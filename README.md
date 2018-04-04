# δyƒƒ /ˈdʏf/
A diff tool for YAML files, and sometimes JSON

[![Build Status](https://travis-ci.org/HeavyWombat/dyff.svg?branch=master)](https://travis-ci.org/HeavyWombat/dyff) [![GoDoc](https://godoc.org/github.com/HeavyWombat/dyff?status.svg)](https://godoc.org/github.com/HeavyWombat/dyff)

![dyff between example](doc/images/dyff-between-example.png?raw=true "dyff between example of two cf-deployment versions")

Only the differences between the files will be displayed by their path (location inside the YAML tree) and the respective change. The input location can either be a local file or a URI. A location named `-` is supported to read from STDIN. The default report style is loosely based upon the old [BOSH v1](https://bosh.io/) deployment delta output. The path style is based on the Dot-style that is used by [Spruce](https://github.com/geofffranks/spruce). As an alternative, you can specify to have paths use [go-patch](https://github.com/cppforlife/go-patch) style. Have a look in the help section of the respective subcommand to see more options.

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
