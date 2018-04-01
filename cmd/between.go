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
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/HeavyWombat/color"
	"github.com/HeavyWombat/dyff/core"
	"github.com/spf13/cobra"
)

var style string

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
		// TODO Add helper function to print absolute path in case it is not a URL, or STDIN indicator -
		fromLocation := args[0]
		toLocation := args[1]

		start := time.Now()

		from, to, err := core.LoadYAMLs(fromLocation, toLocation)
		if err != nil {
			panic(err)
		}

		diffs := core.CompareDocuments(from, to)

		elapsed := time.Since(start)

		// TODO Add style Go-Patch
		// TODO Add style Spruce
		// TODO Add style JSON report
		// TODO Add style YAML report
		// TODO Add style one-line report

		switch strings.ToLower(style) {
		case "human", "bosh":
			fmt.Printf(`      _        __  __
    _| |_   _ / _|/ _|  between %s
  / _' | | | | |_| |_       and %s
 | (_| | |_| |  _|  _|
  \__,_|\__, |_| |_|    %s
        |___/           %s

`, niceLocation(fromLocation),
				niceLocation(toLocation),
				fmt.Sprintf("Number of differences found: %d", len(diffs)),
				fmt.Sprintf("Processing time: %s", elapsed))
			fmt.Print(core.DiffsToHumanStyle(diffs))

		default:
			fmt.Printf("Unkown output style %s\n", style)
			cmd.Usage()
		}
	},
}

func niceLocation(location string) string {
	if location == "-" {
		return core.Italic("<stdin>")
	}

	if _, err := url.ParseRequestURI(location); err == nil {
		return core.Color(location, color.FgHiBlue, color.Underline)
	}

	if abs, err := filepath.Abs(location); err == nil {
		return core.Bold(abs)
	}

	return location
}

func init() {
	rootCmd.AddCommand(betweenCmd)

	// TODO Add flag for swap
	// TODO Add flag for filter on path
	betweenCmd.PersistentFlags().StringVarP(&style, "output", "o", "human", "Specify the output style, e.g. 'human' (more to come ...)")
	betweenCmd.PersistentFlags().BoolVarP(&core.NoTableStyle, "no-table-style", "t", false, "Disable the table output")
	betweenCmd.PersistentFlags().BoolVarP(&core.UseGoPatchPaths, "use-go-patch-style", "g", false, "Use Go-Patch style paths instead of Spruce Dot-Style")
}
