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
)

// jsonCmd represents the json command
var jsonCmd = &cobra.Command{
	Use:   "json [flags] <file-location> ...",
	Args:  cobra.MinimumNArgs(1),
	Short: "Converts input documents into JSON format",
	Long: `
Converts input document into JSON format while preserving the order of all keys.
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

				if plainMode {
					output, err := neat.NewOutputProcessor(false, false, &neat.DefaultColorSchema).ToCompactJSON(document)
					if err != nil {
						exitWithError("Failed to marshal object into JSON", err)
					}

					fmt.Printf("%s\n", string(output))

				} else {
					output, err := neat.NewOutputProcessor(!omitIndentHelper, true, &neat.DefaultColorSchema).ToJSON(document)
					if err != nil {
						exitWithError("Failed to marshal object into JSON", err)
					}

					fmt.Printf("%s\n", output)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(jsonCmd)

	jsonCmd.Flags().SortFlags = false
	jsonCmd.PersistentFlags().SortFlags = false

	jsonCmd.PersistentFlags().BoolVarP(&plainMode, "plain", "p", false, "output in plain style without any highlighting")
	jsonCmd.PersistentFlags().BoolVarP(&restructure, "restructure", "r", false, "restructure map keys in reasonable order")
	jsonCmd.PersistentFlags().BoolVarP(&omitIndentHelper, "omit-indent-helper", "i", false, "omit indent helper lines in highlighted output")
}
