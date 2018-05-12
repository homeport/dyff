// Copyright © 2018 Matthias Diester
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"

	"github.com/HeavyWombat/dyff/pkg/dyff"
	"github.com/HeavyWombat/dyff/pkg/neat"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

var restructure bool
var plainYAML bool
var omitIndentHelper bool

// yamlCmd represents the yaml command
var yamlCmd = &cobra.Command{
	Use:     "yaml [flags] <file-location> ...",
	Aliases: []string{"yml"},
	Short:   "Converts input document into YAML format",
	Long: `
Converts input document into YAML format while preserving the order of all keys.

`,

	Run: func(cmd *cobra.Command, args []string) {
		for _, argument := range args {
			inputFile, err := dyff.LoadFile(argument)
			if err != nil {
				exitWithError("Failed to load input file", err)
			}

			for _, document := range inputFile.Documents {
				if restructure {
					document = dyff.RestructureObject(document)
				}

				if plainYAML { // Run Go YAML library marshalling if plain mode is enabled
					output, err := yaml.Marshal(document)
					if err != nil {
						exitWithError("Failed to marshal object into YAML", err)
					}

					fmt.Printf("---\n%s\n", string(output))

				} else { // Run neat mode to create colorful YAML string
					output, err := neat.NewOutputProcessor(!omitIndentHelper, true, &neat.DefaultColorSchema).ToString(document)
					if err != nil {
						exitWithError("Failed to neatly marshal object into YAML", err)
					}

					fmt.Printf("---\n%s\n", output)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(yamlCmd)

	yamlCmd.PersistentFlags().BoolVarP(&restructure, "restructure", "r", false, "Try to restructure YAML map keys in reasonable order")
	yamlCmd.PersistentFlags().BoolVarP(&plainYAML, "plain", "p", false, "Output YAML in plain style without highlighting")
	yamlCmd.PersistentFlags().BoolVarP(&omitIndentHelper, "omit-indent-helper", "i", false, "Omit indent helper lines in highlighted output")
}
