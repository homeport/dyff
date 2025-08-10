## dyff between

Compare differences between input files from and to

### Synopsis


Compares differences between files and displays the delta. Supported input file
types are: YAML (http://yaml.org/) and JSON (http://json.org/).


```
dyff between [flags] <from> <to>
```

### Options

```
  -i, --ignore-order-changes                ignore order changes in lists
      --ignore-whitespace-changes           ignore leading or trailing whitespace changes
      --detect-kubernetes                   detect kubernetes entities (default true)
      --additional-identifier stringArray   use additional identifier candidates in named entry lists
      --filter strings                      filter reports to a subset of differences based on supplied arguments
      --exclude strings                     exclude reports from a set of differences based on supplied arguments
      --filter-regexp strings               filter reports to a subset of differences based on supplied regular expressions
      --exclude-regexp strings              exclude reports from a set of differences based on supplied regular expressions
  -v, --ignore-value-changes                exclude changes in values
      --detect-renames                      enable detection for renames (document level for Kubernetes resources) (default true)
  -o, --output string                       specify the output style, supported styles: human, brief, github, gitlab, gitea (default "human")
      --use-indent-lines                    use indent lines in the output (default true)
  -b, --omit-header                         omit the dyff summary header
  -s, --set-exit-code                       set program exit code, with 0 meaning no difference, 1 for differences detected, and 255 for program error
  -l, --no-table-style                      do not place blocks next to each other, always use one row per text block
  -x, --no-cert-inspection                  disable x509 certificate inspection, compare as raw text
  -g, --use-go-patch-style                  use Go-Patch style paths in outputs
      --minor-change-threshold float        minor change threshold (default 0.1)
      --multi-line-context-lines int        multi-line context lines (default 4)
      --swap                                Swap 'from' and 'to' for comparison
      --chroot string                       change the root level of the input file to another point in the document
      --chroot-of-from string               only change the root level of the from input file
      --chroot-of-to string                 only change the root level of the to input file
      --chroot-list-to-documents            in case the change root points to a list, treat this list as a set of documents and not as the list itself
  -h, --help                                help for between
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

