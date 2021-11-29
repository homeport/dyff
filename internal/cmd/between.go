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
	"github.com/gonvenience/wrap"
	"github.com/gonvenience/ytbx"
	"github.com/spf13/cobra"

	"github.com/homeport/dyff/pkg/dyff"
)

type betweenCmdOptions struct {
	swap                     bool
	translateListToDocuments bool
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
			if err = dyff.ChangeRoot(&from, betweenCmdSettings.chrootFrom, reportOptions.useGoPatchPaths, betweenCmdSettings.translateListToDocuments); err != nil {
				return wrap.Errorf(err, "failed to change root of %s to path %s", from.Location, betweenCmdSettings.chrootFrom)
			}
		}

		// Change root of 'to' input file if change root flag for 'to' is set
		if betweenCmdSettings.chrootTo != "" {
			if err = dyff.ChangeRoot(&to, betweenCmdSettings.chrootTo, reportOptions.useGoPatchPaths, betweenCmdSettings.translateListToDocuments); err != nil {
				return wrap.Errorf(err, "failed to change root of %s to path %s", to.Location, betweenCmdSettings.chrootTo)
			}
		}

		report, err := dyff.CompareInputFiles(from, to,
			dyff.IgnoreOrderChanges(reportOptions.ignoreOrderChanges),
			dyff.KubernetesEntityDetection(reportOptions.kubernetesEntityDetection),
		)
		if err != nil {
			return wrap.Errorf(err, "failed to compare input files")
		}

		if reportOptions.filters != nil {
			var filterPaths []*ytbx.Path
			for _, pathString := range reportOptions.filters {
				path, err := ytbx.ParsePathStringUnsafe(pathString)
				if err != nil {
					return wrap.Errorf(err, "failed to set path filter, because path %s cannot be parsed", pathString)
				}

				filterPaths = append(filterPaths, &path)
			}

			report = report.Filter(filterPaths...)
		}

		return writeReport(cmd, report)
	},
}

func init() {
	rootCmd.AddCommand(betweenCmd)

	betweenCmd.Flags().SortFlags = false

	applyReportOptionsFlags(betweenCmd)

	// Input documents modification flags
	betweenCmd.Flags().BoolVar(&betweenCmdSettings.swap, "swap", false, "Swap 'from' and 'to' for comparison")
	betweenCmd.Flags().StringVar(&betweenCmdSettings.chroot, "chroot", "", "change the root level of the input file to another point in the document")
	betweenCmd.Flags().StringVar(&betweenCmdSettings.chrootFrom, "chroot-of-from", "", "only change the root level of the from input file")
	betweenCmd.Flags().StringVar(&betweenCmdSettings.chrootTo, "chroot-of-to", "", "only change the root level of the to input file")
	betweenCmd.Flags().BoolVar(&betweenCmdSettings.translateListToDocuments, "chroot-list-to-documents", false, "in case the change root points to a list, treat this list as a set of documents and not as the list itself")
}
