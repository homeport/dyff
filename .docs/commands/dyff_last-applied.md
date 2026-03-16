## dyff last-applied

Compare differences between the current state and the one stored in Kubernetes last-applied configuration

### Synopsis


Kubernetes resource YAML (or JSON) contain the previously used configuration of
that resource in the metadata. For convenience, the respective metadata is used
to compare it against the current configuration.


```
dyff last-applied [flags]
```

### Options

```
  -o, --output string                       specify the output style, supported styles: human, brief, github, gitlab, gitea (default "human")
      --use-indent-lines                    use indent lines in the output
  -i, --ignore-order-changes                ignore order changes in lists
      --ignore-whitespace-changes           ignore leading or trailing whitespace changes
  -v, --ignore-value-changes                exclude changes in values
      --detect-renames                      enable detection for renames (document level for Kubernetes resources) (default true)
      --format-strings                      format strings (i.e. inline JSON) before comparison to avoid formatting differences (default true)
  -l, --no-table-style                      do not place blocks next to each other, always use one row per text block
  -x, --no-cert-inspection                  disable x509 certificate inspection, compare as raw text
  -g, --use-go-patch-style                  use Go-Patch style paths in outputs
      --minor-change-threshold float        minor change threshold (default 0.1)
      --multi-line-context-lines int        multi-line context lines (default 10)
      --detect-kubernetes                   detect kubernetes entities (default true)
      --additional-identifier stringArray   use additional identifier candidates in named entry lists
      --filter strings                      filter reports to a subset of differences based on supplied arguments
      --exclude strings                     exclude reports from a set of differences based on supplied arguments
      --filter-regexp strings               filter reports to a subset of differences based on supplied regular expressions
      --exclude-regexp strings              exclude reports from a set of differences based on supplied regular expressions
  -b, --omit-header                         omit the dyff summary header
  -s, --set-exit-code                       set program exit code, with 0 meaning no difference, 1 for differences detected, and 255 for program error
  -h, --help                                help for last-applied
```

### Options inherited from parent commands

```
  -c, --color                        specify color usage: on, off, or auto (default auto)
  -w, --fixed-width int              disable terminal width detection and use provided fixed value (default -1)
  -k, --preserve-key-order-in-json   use ordered keys during JSON decoding (non standard behavior)
  -t, --truecolor                    specify true color usage: on, off, or auto (default auto)
```

### SEE ALSO

* [dyff](dyff.md)	 - dyff

