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
	"fmt"

	"github.com/gonvenience/wrap"
	"github.com/gonvenience/ytbx"
	"github.com/spf13/cobra"
)

type jsonCmdOptions struct {
	plainMode        bool
	restructure      bool
	omitIndentHelper bool
	inplace          bool
}

var jsonCmdSettings jsonCmdOptions

// jsonCmd represents the json command
var jsonCmd = &cobra.Command{
	Use:   "json [flags] <file-location> ...",
	Args:  cobra.MinimumNArgs(1),
	Short: "Converts input documents into JSON format",
	Long: `
Converts input document into JSON format while preserving the order of all keys.
`,

	RunE: func(cmd *cobra.Command, args []string) error {
		writer := &OutputWriter{
			OutputStyle:      "json",
			PlainMode:        jsonCmdSettings.plainMode,
			Restructure:      jsonCmdSettings.restructure,
			OmitIndentHelper: jsonCmdSettings.omitIndentHelper,
		}

		var errors []error
		for _, filename := range args {
			if ytbx.IsStdin(filename) && jsonCmdSettings.inplace {
				return wrap.Error(
					fmt.Errorf("cannot use in-place flag in combination with input from STDIN"),
					"incompatible flags",
				)
			}

			if jsonCmdSettings.inplace {
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
	rootCmd.AddCommand(jsonCmd)

	jsonCmd.Flags().SortFlags = false

	jsonCmd.Flags().BoolVarP(&jsonCmdSettings.plainMode, "plain", "p", false, "output in plain style without any highlighting")
	jsonCmd.Flags().BoolVarP(&jsonCmdSettings.restructure, "restructure", "r", false, "restructure map keys in reasonable order")
	jsonCmd.Flags().BoolVarP(&jsonCmdSettings.omitIndentHelper, "omit-indent-helper", "O", false, "omit indent helper lines in highlighted output")
	jsonCmd.Flags().BoolVarP(&jsonCmdSettings.inplace, "in-place", "i", false, "overwrite input file with output of this command")
}
