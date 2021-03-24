# δyƒƒ /ˈdʏf/ [![License](https://img.shields.io/github/license/homeport/dyff.svg)](https://github.com/homeport/dyff/blob/main/LICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/homeport/dyff)](https://goreportcard.com/report/github.com/homeport/dyff) [![Build and Tests](https://github.com/homeport/dyff/workflows/Build%20and%20Tests/badge.svg)](https://github.com/homeport/dyff/actions?query=workflow%3A%22Build+and+Tests%22) [![Codecov](https://img.shields.io/codecov/c/github/homeport/dyff/main.svg)](https://codecov.io/gh/homeport/dyff) [![Go Reference](https://pkg.go.dev/badge/github.com/homeport/dyff.svg)](https://pkg.go.dev/github.com/homeport/dyff) [![Release](https://img.shields.io/github/release/homeport/dyff.svg)](https://github.com/homeport/dyff/releases/latest)

![dyff](.docs/logo.png?raw=true "dyff logo - the letters d, y, and f in the colors green, yellow and red")

A diff tool for YAML files, and sometimes JSON

## Description

`dyff` is inspired by the way the old [BOSH v1](https://bosh.io/) deployment output reported changes from one version to another by only showing the parts of a YAML file that change.

Each difference is referenced by its location in the YAML document by using either the [Spruce](https://github.com/geofffranks/spruce) or [go-patch](https://github.com/cppforlife/go-patch) path syntax. The output report aims to be as compact as possible to give a clear and simple overview of the change.

Similar to the standard `diff` tool, it follows the principle of describing the change by going from the `from` input file to the target `to` input file.

Input files can be local files (filesystem path), remote files (URI), or the standard input stream (using `-`).

All orders of keys in hashes are preserved during processing and output to the terminal, most notably in the sub-commands to convert YAML to JSON and vice versa.

## Installation

### Homebrew

The `homeport/tap` has macOS and GNU/Linux pre-built binaries available:

```bash
brew install homeport/tap/dyff
```

### Snap

It is [available in the `snapcraft` store](https://snapcraft.io/dyff) in the Productivity section.

```bash
snap install dyff
```

### Pre-built binaries in GitHub

Prebuilt binaries can be [downloaded from the GitHub Releases section](https://github.com/homeport/dyff/releases/latest).

### Curl To Shell Convenience Script

There is a convenience script to download the latest release for Linux or macOS if you want to need it simple (you need `curl` and `jq` installed on your machine):

```bash
curl --silent --location https://git.io/JYfAY | bash
```

### Build from Source

You can download and build `dyff` from source using `go get`:

```bash
GO111MODULE=on go get github.com/homeport/dyff/cmd/dyff
```

## Use cases and examples

- Show differences between the live configuration of Kubernetes resources and what would be applied:

  ```bash
  # Setup
  export KUBECTL_EXTERNAL_DIFF="dyff between --omit-header --set-exit-code"

  # Usage
  kubectl diff [...]
  ```

  ![dyff between example with kubectl diff](.docs/dyff-between-kubectl-diff.png?raw=true "dyff in kubectl diff example")

  The `--set-exit-code` flag is required so that the `dyff` exit code matches `kubectl` expectations. An exit code `0` refers to no differences, `1` in case differences are detected. Other exit codes are treated as program issues.

- Show the differences between two versions of [`cf-deployment`](https://github.com/cloudfoundry/cf-deployment/) YAMLs:

    ```bash
    dyff between \
      https://raw.githubusercontent.com/cloudfoundry/cf-deployment/v1.10.0/cf-deployment.yml \
      https://raw.githubusercontent.com/cloudfoundry/cf-deployment/v1.20.0/cf-deployment.yml
    ```

    ![dyff between example](.docs/dyff-between-deployment-manifest-example.png?raw=true "dyff between example of two cf-deployment versions")

- Embed `dyff` into **Git** for better understandable differences

    ```bash
    # Setup...
    git config --local diff.dyff.command 'dyff_between() { dyff --color on between --omit-header "$2" "$5"; }; dyff_between'
    echo '*.yml diff=dyff' >> .gitattributes

    # And have fun, e.g.:
    git log --ext-diff -u
    git show --ext-diff HEAD
    ```

    ![dyff between example of a Git commit](.docs/dyff-between-git-commits-example.png?raw=true "dyff in Git example of an example commit")

- Convert a JSON stream to YAML

    ```bash
    sometool --json | jq --raw-output '.data' | dyff yaml -
    ```

- Sometimes you end up with YAML or JSON files, where the order of the keys in maps was sorted alphabetically. With `dyff` you can restructure keys in maps to a more human appealing order:

    ```bash
    sometool --export --json | dyff yaml --restructure -
    ```

    Or, rewrite a file _in place_ with the restructured order of keys.

    ```bash
    dyff yaml --restructure --in-place somefile.yml
    ```

- Just print a YAML (or JSON) file to the terminal to look at it. By default, `dyff` will use a neat output schema which includes different colors and indent helper lines to improve readability. The colors are roughly based on the default [Atom](https://atom.io) schema and work best on dark terminal backgrounds. The neat output is disabled if the output of `dyff` is redirected into a pipe, or you can disable it explicitly using the `--plain` flag.

    ```bash
    dyff yaml somefile.yml
    ```

- Convert a YAML file to JSON and vice versa:

    ```bash
    dyff json https://raw.githubusercontent.com/cloudfoundry/cf-deployment/v1.19.0/cf-deployment.yml
    ```

    The `dyff` sub-command (`yaml`, or `json`) defines the output format, the tool automatically detects the input format itself.

    ```bash
    dyff yaml https://raw.githubusercontent.com/homeport/dyff/develop/assets/bosh-yaml/manifest.json
    ```
