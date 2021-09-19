// Copyright © 2019 The Homeport Team
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
	"github.com/gonvenience/bunt"
	"github.com/gonvenience/wrap"
	"github.com/gonvenience/ytbx"
	"github.com/spf13/cobra"
)

type yamlCmdOptions struct {
	plainMode        bool
	restructure      bool
	omitIndentHelper bool
	inplace          bool
}

var yamlCmdSettings yamlCmdOptions

// yamlCmd represents the yaml command
var yamlCmd = &cobra.Command{
	Use:     "yaml [flags] <file-location> ...",
	Aliases: []string{"yml"},
	Args:    cobra.MinimumNArgs(1),
	Short:   "Converts input documents into YAML format",
	Long: `
Converts input document into YAML format while preserving the order of all keys.
`,

	RunE: func(cmd *cobra.Command, args []string) error {
		writer := &OutputWriter{
			OutputStyle:      "yaml",
			PlainMode:        yamlCmdSettings.plainMode,
			Restructure:      yamlCmdSettings.restructure,
			OmitIndentHelper: yamlCmdSettings.omitIndentHelper,
		}

		var errors []error
		for _, filename := range args {
			if ytbx.IsStdin(filename) && yamlCmdSettings.inplace {
				return wrap.Error(
					bunt.Errorf("cannot use in-place flag in combination with input from _*stdin*_"),
					"incompatible flags",
				)
			}

			if yamlCmdSettings.inplace {
				if err := writer.WriteInplace(filename); err != nil {
					errors = append(errors, err)
				}
			} else {
				if err := writer.WriteToStdout(filename); err != nil {
					errors = append(errors, err)
				}
			}
		}

		if len(errors) > 0 {
			return wrap.Errors(errors, "failed to process input files")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(yamlCmd)

	yamlCmd.Flags().SortFlags = false

	yamlCmd.Flags().BoolVarP(&yamlCmdSettings.plainMode, "plain", "p", false, "output in plain style without any highlighting")
	yamlCmd.Flags().BoolVarP(&yamlCmdSettings.restructure, "restructure", "r", false, "restructure map keys in reasonable order")
	yamlCmd.Flags().BoolVarP(&yamlCmdSettings.omitIndentHelper, "omit-indent-helper", "O", false, "omit indent helper lines in highlighted output")
	yamlCmd.Flags().BoolVarP(&yamlCmdSettings.inplace, "in-place", "i", false, "overwrite input file with output of this command")
}
