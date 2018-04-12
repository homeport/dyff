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

	"github.com/HeavyWombat/dyff/core"
	"github.com/HeavyWombat/yaml"
	"github.com/spf13/cobra"
)

var restructure bool

// yamlCmd represents the yaml command
var yamlCmd = &cobra.Command{
	Use:     "yaml",
	Aliases: []string{"yml"},
	Short:   "Converts input document into YAML format",
	Long: `
Converts input document into YAML format while preserving the order of all keys.

`,

	Run: func(cmd *cobra.Command, args []string) {
		for _, x := range args {
			obj, err := core.LoadFile(x)
			if err != nil {
				core.ExitWithError("Failed to load input file", err)
			}

			switch obj.(type) {
			case yaml.MapSlice:
				mapslice := obj.(yaml.MapSlice)

				if restructure {
					mapslice = core.RestructureMapSlice(mapslice)
				}

				output, yamlerr := core.ToYAMLString(mapslice)
				if yamlerr != nil {
					core.ExitWithError("Failed to marshal object into YAML", err)
				}

				fmt.Print(output)

			default:
				core.ExitWithError("Failed to process file",
					fmt.Errorf("Provided input file is not YAML compatible"))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(yamlCmd)

	yamlCmd.PersistentFlags().BoolVarP(&restructure, "restructure", "r", false, "Try to restructure YAML map keys in reasonable order")
}
