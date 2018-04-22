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
	"strings"

	"github.com/HeavyWombat/dyff/pkg/dyff"
	"github.com/spf13/cobra"
)

var style string
var swap bool

// betweenCmd represents the between command
var betweenCmd = &cobra.Command{
	Use:   "between",
	Short: "Compares differences between documents",
	Long: `
Compares differences between documents and displays the delta. Supported
document types are: YAML (http://yaml.org/) and JSON (http://json.org/).

`,
	Args:    cobra.ExactArgs(2),
	Aliases: []string{"bw"},
	Run: func(cmd *cobra.Command, args []string) {
		var fromLocation, toLocation string
		if swap {
			fromLocation = args[1]
			toLocation = args[0]
		} else {
			fromLocation = args[0]
			toLocation = args[1]
		}

		from, to, err := dyff.LoadFiles(fromLocation, toLocation)
		if err != nil {
			exitWithError("Failed to load input files", err)
		}

		report, err := dyff.CompareInputFiles(from, to)
		if err != nil {
			exitWithError("Failed to compare input files", err)
		}

		// TODO Add style Go-Patch
		// TODO Add style Spruce
		// TODO Add style JSON report
		// TODO Add style YAML report
		// TODO Add style one-line report

		switch strings.ToLower(style) {
		case "human", "bosh":
			fmt.Print(dyff.CreateHumanStyleReport(report, true))

		default:
			fmt.Printf("Unknown output style %s\n", style)
			cmd.Usage()
		}
	},
}

func init() {
	rootCmd.AddCommand(betweenCmd)

	// TODO Add flag for filter on path
	betweenCmd.PersistentFlags().StringVarP(&style, "output", "o", "human", "Specify the output style, e.g. 'human' (more to come ...)")
	betweenCmd.PersistentFlags().BoolVarP(&swap, "swap", "s", false, "Swap `from` and `to` for compare")
	betweenCmd.PersistentFlags().BoolVarP(&dyff.NoTableStyle, "no-table-style", "t", false, "Disable the table output")
	betweenCmd.PersistentFlags().BoolVarP(&dyff.DoNotInspectCerts, "no-cert-inspection", "c", false, "Disable certificate inspection (compare as raw text)")
	betweenCmd.PersistentFlags().BoolVarP(&dyff.UseGoPatchPaths, "use-go-patch-style", "g", false, "Use Go-Patch style paths instead of Spruce Dot-Style")
}
