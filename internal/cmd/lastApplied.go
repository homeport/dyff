// Copyright © 2020 The Homeport Team
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

	"github.com/gonvenience/ytbx"
	"github.com/spf13/cobra"
	yamlv3 "go.yaml.in/yaml/v3"

	"github.com/homeport/dyff/pkg/dyff"
)

// lastAppliedCmd represents the lastApplied command
var lastAppliedCmd = &cobra.Command{
	Use:   "last-applied",
	Short: "Compare differences between the current state and the one stored in Kubernetes last-applied configuration",
	Long: `
Kubernetes resource YAML (or JSON) contain the previously used configuration of
that resource in the metadata. For convenience, the respective metadata is used
to compare it against the current configuration.
`,
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"la"},
	RunE: func(cmd *cobra.Command, args []string) error {
		inputFile, err := ytbx.LoadFile(args[0])
		if err != nil {
			return err
		}

		if len(inputFile.Documents) != 1 {
			return fmt.Errorf("failed to compare, because the input contains more than one document")
		}

		lastConfiguration, err := lookUpLastAppliedConfiguration(inputFile)
		if err != nil {
			return err
		}

		purgeWellKnownMetadataEntries(inputFile.Documents[0])

		report, err := dyff.CompareInputFiles(lastConfiguration, inputFile, dyff.IgnoreOrderChanges(reportOptions.ignoreOrderChanges))
		if err != nil {
			return fmt.Errorf("failed to compare input files: %w", err)
		}

		return writeReport(cmd, report)
	},
}

func init() {
	rootCmd.AddCommand(lastAppliedCmd)

	lastAppliedCmd.Flags().SortFlags = false

	applyReportOptionsFlags(lastAppliedCmd)
}

func lookUpLastAppliedConfiguration(inputFile ytbx.InputFile) (ytbx.InputFile, error) {
	kubectlLastApplied, err := ytbx.Grab(inputFile.Documents[0], "/metadata/annotations/kubectl.kubernetes.io\\/last-applied-configuration")
	if err != nil {
		return ytbx.InputFile{}, fmt.Errorf("provided input file does not contain the last applied configuration metadata")
	}

	documents, err := ytbx.LoadDocuments([]byte(kubectlLastApplied.Value))
	if err != nil {
		return ytbx.InputFile{}, err
	}

	return ytbx.InputFile{
		Documents: documents,
		Location:  "/metadata/annotations/kubectl.kubernetes.io/last-applied-configuration",
	}, nil
}

func purgeWellKnownMetadataEntries(document *yamlv3.Node) {
	_, _ = ytbx.Delete(document, "/metadata/annotations/kubectl.kubernetes.io\\/last-applied-configuration")
}
