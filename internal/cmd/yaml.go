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

	"github.com/HeavyWombat/dyff/pkg/v1/dyff"
	"github.com/spf13/cobra"
)

// yamlCmd represents the yaml command
var yamlCmd = &cobra.Command{
	Use:     "yaml [flags] <file-location> ...",
	Aliases: []string{"yml"},
	Args:    cobra.MinimumNArgs(1),
	Short:   "Converts input documents into YAML format",
	Long: `
Converts input document into YAML format while preserving the order of all keys.
`,

	Run: func(cmd *cobra.Command, args []string) {
		writer := &OutputWriter{
			PlainMode:        plainMode,
			Restructure:      restructure,
			OmitIndentHelper: omitIndentHelper,
			OutputStyle:      "yaml",
		}

		for _, filename := range args {
			if dyff.IsStdin(filename) && inplace {
				exitWithError("incompatible flag", fmt.Errorf("cannot use in-place flag in combination with input from STDIN"))
			}

			if inplace {
				writer.WriteInplace(filename)

			} else {
				writer.WriteToStdout(filename)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(yamlCmd)

	yamlCmd.Flags().SortFlags = false
	yamlCmd.PersistentFlags().SortFlags = false

	yamlCmd.PersistentFlags().BoolVarP(&plainMode, "plain", "p", false, "output in plain style without any highlighting")
	yamlCmd.PersistentFlags().BoolVarP(&restructure, "restructure", "r", false, "restructure map keys in reasonable order")
	yamlCmd.PersistentFlags().BoolVarP(&omitIndentHelper, "omit-indent-helper", "O", false, "omit indent helper lines in highlighted output")
	yamlCmd.PersistentFlags().BoolVarP(&inplace, "in-place", "i", false, "overwrite input file with output of this command")
}
