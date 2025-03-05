## dyff json

Converts input documents into JSON format

### Synopsis


Converts input document into JSON format while preserving the order of all keys.


```
dyff json [flags] <file-location> ...
```

### Options

```
  -p, --plain                output in plain style without any highlighting
  -r, --restructure          restructure map keys in reasonable order
  -O, --omit-indent-helper   omit indent helper lines in highlighted output
  -i, --in-place             overwrite input file with output of this command
  -h, --help                 help for json
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

