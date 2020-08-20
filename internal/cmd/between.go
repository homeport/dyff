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
	"math"
	"os"
	"strings"

	"github.com/gonvenience/wrap"
	"github.com/gonvenience/ytbx"
	"github.com/homeport/dyff/pkg/dyff"
	"github.com/spf13/cobra"
)

const defaultOutputStyle = "human"

type betweenCmdOptions struct {
	style                    string
	swap                     bool
	noTableStyle             bool
	doNotInspectCerts        bool
	exitWithCount            bool
	translateListToDocuments bool
	omitHeader               bool
	useGoPatchPaths          bool
	ignoreOrderChanges       bool
	chroot                   string
	chrootFrom               string
	chrootTo                 string
}

var betweenCmdSettings betweenCmdOptions

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
	RunE: func(cmd *cobra.Command, args []string) error {
		var fromLocation, toLocation string
		if betweenCmdSettings.swap {
			fromLocation = args[1]
			toLocation = args[0]
		} else {
			fromLocation = args[0]
			toLocation = args[1]
		}

		from, to, err := ytbx.LoadFiles(fromLocation, toLocation)
		if err != nil {
			return wrap.Errorf(err, "failed to load input files")
		}

		// If the main change root flag is set, this (re-)sets the individual change roots of the two input files
		if betweenCmdSettings.chroot != "" {
			betweenCmdSettings.chrootFrom = betweenCmdSettings.chroot
			betweenCmdSettings.chrootTo = betweenCmdSettings.chroot
		}

		// Change root of 'from' input file if change root flag for 'from' is set
		if betweenCmdSettings.chrootFrom != "" {
			if err = dyff.ChangeRoot(&from, betweenCmdSettings.chrootFrom, betweenCmdSettings.useGoPatchPaths, betweenCmdSettings.translateListToDocuments); err != nil {
				return wrap.Errorf(err, "failed to change root of %s to path %s", from.Location, betweenCmdSettings.chrootFrom)
			}
		}

		// Change root of 'to' input file if change root flag for 'to' is set
		if betweenCmdSettings.chrootTo != "" {
			if err = dyff.ChangeRoot(&to, betweenCmdSettings.chrootTo, betweenCmdSettings.useGoPatchPaths, betweenCmdSettings.translateListToDocuments); err != nil {
				return wrap.Errorf(err, "failed to change root of %s to path %s", to.Location, betweenCmdSettings.chrootTo)
			}
		}

		report, err := dyff.CompareInputFiles(from, to, dyff.IgnoreOrderChanges(betweenCmdSettings.ignoreOrderChanges))
		if err != nil {
			return wrap.Errorf(err, "failed to compare input files")
		}

		var reportWriter dyff.ReportWriter
		switch strings.ToLower(betweenCmdSettings.style) {
		case "human", "bosh":
			reportWriter = &dyff.HumanReport{
				Report:               report,
				DoNotInspectCerts:    betweenCmdSettings.doNotInspectCerts,
				NoTableStyle:         betweenCmdSettings.noTableStyle,
				ShowBanner:           !betweenCmdSettings.omitHeader,
				UseGoPatchPaths:      betweenCmdSettings.useGoPatchPaths,
				MinorChangeThreshold: 0.1,
			}

		case "brief", "short", "summary":
			reportWriter = &dyff.BriefReport{
				Report: report,
			}

		default:
			return wrap.Errorf(
				fmt.Errorf(cmd.UsageString()),
				"unknown output style %s", betweenCmdSettings.style,
			)
		}

		err = reportWriter.WriteReport(os.Stdout)
		if err != nil {
			return wrap.Errorf(err, "failed to print report")
		}

		// If configured, make sure `dyff` exists with an exit status
		if betweenCmdSettings.exitWithCount {
			return ExitCode{
				Value: int(math.Min(float64(len(report.Diffs)), 255.0)),
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(betweenCmd)

	betweenCmd.Flags().SortFlags = false
	betweenCmd.PersistentFlags().SortFlags = false

	// Main output preferences
	betweenCmd.PersistentFlags().StringVarP(&betweenCmdSettings.style, "output", "o", defaultOutputStyle, "specify the output style, supported styles: human, or brief")
	betweenCmd.PersistentFlags().BoolVarP(&betweenCmdSettings.omitHeader, "omit-header", "b", false, "omit the dyff summary header")
	betweenCmd.PersistentFlags().BoolVarP(&betweenCmdSettings.exitWithCount, "set-exit-status", "s", false, "set exit status to number of diff (capped at 255)")

	// Human/BOSH output related flags
	betweenCmd.PersistentFlags().BoolVarP(&betweenCmdSettings.noTableStyle, "no-table-style", "l", false, "do not place blocks next to each other, always use one row per text block")
	betweenCmd.PersistentFlags().BoolVarP(&betweenCmdSettings.doNotInspectCerts, "no-cert-inspection", "x", false, "disable x509 certificate inspection, compare as raw text")
	betweenCmd.PersistentFlags().BoolVarP(&betweenCmdSettings.useGoPatchPaths, "use-go-patch-style", "g", false, "use Go-Patch style paths in outputs")

	// Compare options
	betweenCmd.PersistentFlags().BoolVarP(&betweenCmdSettings.ignoreOrderChanges, "ignore-order-changes", "i", false, "ignore order changes in lists")

	// Input documents modification flags
	betweenCmd.PersistentFlags().BoolVar(&betweenCmdSettings.swap, "swap", false, "Swap 'from' and 'to' for comparison")
	betweenCmd.PersistentFlags().StringVar(&betweenCmdSettings.chroot, "chroot", "", "change the root level of the input file to another point in the document")
	betweenCmd.PersistentFlags().StringVar(&betweenCmdSettings.chrootFrom, "chroot-of-from", "", "only change the root level of the from input file")
	betweenCmd.PersistentFlags().StringVar(&betweenCmdSettings.chrootTo, "chroot-of-to", "", "only change the root level of the to input file")
	betweenCmd.PersistentFlags().BoolVar(&betweenCmdSettings.translateListToDocuments, "chroot-list-to-documents", false, "in case the change root points to a list, treat this list as a set of documents and not as the list itself")
}
