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
	"math"
	"os"
	"strings"

	"github.com/HeavyWombat/dyff/pkg/v1/dyff"
	"github.com/spf13/cobra"
)

var style string
var swap bool
var noTableStyle bool
var doNotInspectCerts bool
var exitWithCount bool
var translateListToDocuments bool
var chroot string
var chrootFrom string
var chrootTo string

// betweenCmd represents the between command
var betweenCmd = &cobra.Command{
	Use:   "between [flags] <from> <to>",
	Short: "Compare differences between input files from and to",
	Long: `
Compares differences between files and displays the delta. Supported input file
types are: YAML (http://yaml.org/) and JSON (http://json.org/).
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

		// If the main change root flag is set, this (re-)sets the individual change roots of the two input files
		if chroot != "" {
			chrootFrom = chroot
			chrootTo = chroot
		}

		// Change root of from input file if change root flag for form is set
		if chrootFrom != "" {
			if err = dyff.ChangeRoot(&from, chrootFrom, translateListToDocuments); err != nil {
				exitWithError(fmt.Sprintf("Failed to change root of %s to path %s", from.Location, chrootFrom), err)
			}
		}

		// Change root of to input file if change root flag for to is set
		if chrootTo != "" {
			if err = dyff.ChangeRoot(&to, chrootTo, translateListToDocuments); err != nil {
				exitWithError(fmt.Sprintf("Failed to change root of %s to path %s", to.Location, chrootTo), err)
			}
		}

		report, err := dyff.CompareInputFiles(from, to)
		if err != nil {
			exitWithError("Failed to compare input files", err)
		}

		// If configured, make sure `dyff` exists with an exit status
		if exitWithCount {
			defer os.Exit(int(math.Min(
				float64(len(report.Diffs)),
				255.0)))
		}

		// TODO Add style Go-Patch
		// TODO Add style Spruce
		// TODO Add style JSON report
		// TODO Add style YAML report

		var reportWriter dyff.ReportWriter
		switch strings.ToLower(style) {
		case "human", "bosh":
			reportWriter = &dyff.HumanReport{
				Report:            report,
				DoNotInspectCerts: doNotInspectCerts,
				NoTableStyle:      noTableStyle,
				ShowBanner:        true,
			}

		case "brief", "short", "summary":
			reportWriter = &dyff.BriefReport{
				Report: report,
			}

		default:
			fmt.Printf("Unknown output style %s\n", style)
			cmd.Usage()
		}

		reportWriter.WriteReport(os.Stdout)
	},
}

func init() {
	rootCmd.AddCommand(betweenCmd)

	betweenCmd.Flags().SortFlags = false
	betweenCmd.PersistentFlags().SortFlags = false

	// Main output preferences
	betweenCmd.PersistentFlags().StringVarP(&style, "output", "o", "human", "specify the output style, supported style: human")
	betweenCmd.PersistentFlags().BoolVarP(&exitWithCount, "set-exit-status", "s", false, "set exit status to number of diff (capped at 255)")

	// Human/BOSH output related flags
	betweenCmd.PersistentFlags().BoolVarP(&noTableStyle, "no-table-style", "l", false, "do not place blocks next to each other, always use one row per text block")
	betweenCmd.PersistentFlags().BoolVarP(&doNotInspectCerts, "no-cert-inspection", "x", false, "disable x509 certificate inspection, compare as raw text")

	// General `dyff` package related preferences
	betweenCmd.PersistentFlags().BoolVarP(&dyff.UseGoPatchPaths, "use-go-patch-style", "g", false, "use Go-Patch style paths in outputs")

	// Input documents modification flags
	betweenCmd.PersistentFlags().BoolVar(&swap, "swap", false, "Swap 'from' and 'to' for comparison")
	betweenCmd.PersistentFlags().StringVar(&chroot, "chroot", "", "change the root level of the input file to another point in the document")
	betweenCmd.PersistentFlags().StringVar(&chrootFrom, "chroot-of-from", "", "only change the root level of the from input file")
	betweenCmd.PersistentFlags().StringVar(&chrootTo, "chroot-of-to", "", "only change the root level of the to input file")
	betweenCmd.PersistentFlags().BoolVar(&translateListToDocuments, "chroot-list-to-documents", false, "in case the change root points to a list, treat this list as a set of documents and not as the list itself")
}
