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
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gonvenience/bunt"
	"github.com/gonvenience/neat"
	"github.com/gonvenience/wrap"
	"github.com/gonvenience/ytbx"
	"github.com/spf13/cobra"
	yamlv3 "gopkg.in/yaml.v3"

	"github.com/homeport/dyff/pkg/dyff"
)

const defaultOutputStyle = "human"

type reportConfig struct {
	style                     string
	ignoreOrderChanges        bool
	kubernetesEntityDetection bool
	noTableStyle              bool
	doNotInspectCerts         bool
	exitWithCode              bool
	omitHeader                bool
	useGoPatchPaths           bool
	filters                   []string
}

var reportOptions reportConfig

func applyReportOptionsFlags(cmd *cobra.Command) {
	// Compare options
	cmd.Flags().BoolVarP(&reportOptions.ignoreOrderChanges, "ignore-order-changes", "i", false, "ignore order changes in lists")
	cmd.Flags().BoolVarP(&reportOptions.kubernetesEntityDetection, "detect-kubernetes", "", true, "detect kubernetes entities")
	cmd.Flags().StringSliceVar(&reportOptions.filters, "filter", nil, "filter reports to a subset of differences based on supplied arguments")

	// Main output preferences
	cmd.Flags().StringVarP(&reportOptions.style, "output", "o", defaultOutputStyle, "specify the output style, supported styles: human, or brief")
	cmd.Flags().BoolVarP(&reportOptions.omitHeader, "omit-header", "b", false, "omit the dyff summary header")
	cmd.Flags().BoolVarP(&reportOptions.exitWithCode, "set-exit-code", "s", false, "set program exit code, with 0 meaning no difference, 1 for differences detected, and 255 for program error")

	// Human/BOSH output related flags
	cmd.Flags().BoolVarP(&reportOptions.noTableStyle, "no-table-style", "l", false, "do not place blocks next to each other, always use one row per text block")
	cmd.Flags().BoolVarP(&reportOptions.doNotInspectCerts, "no-cert-inspection", "x", false, "disable x509 certificate inspection, compare as raw text")
	cmd.Flags().BoolVarP(&reportOptions.useGoPatchPaths, "use-go-patch-style", "g", false, "use Go-Patch style paths in outputs")

	// Deprecated
	cmd.Flags().BoolVar(&reportOptions.exitWithCode, "set-exit-status", false, "set program exit code, with 0 meaning no difference, 1 for differences detected, and 255 for program error")
	_ = cmd.Flags().MarkDeprecated("set-exit-status", "use --set-exit-code instead")
}

// OutputWriter encapsulates the required fields to define the look and feel of
// the output
type OutputWriter struct {
	PlainMode        bool
	Restructure      bool
	OmitIndentHelper bool
	OutputStyle      string
}

func humanReadableFilename(filename string) string {
	if ytbx.IsStdin(filename) {
		return bunt.Sprint("_*stdin*_")
	}

	return bunt.Sprintf("_*%s*_", filename)
}

// WriteToStdout is a convenience function to write the content of the documents
// stored in the provided input file to the standard output
func (w *OutputWriter) WriteToStdout(filename string) error {
	if err := w.write(os.Stdout, filename); err != nil {
		return wrap.Error(err, bunt.Sprint("failed to write output to _*stdout*_"))
	}

	return nil
}

// WriteInplace writes the content of the documents stored in the provided input
// file to the file itself overwriting the content in place.
func (w *OutputWriter) WriteInplace(filename string) error {
	var buf bytes.Buffer
	bufWriter := bufio.NewWriter(&buf)

	// Force plain mode to make sure there are no ANSI sequences
	w.PlainMode = true
	if err := w.write(bufWriter, filename); err != nil {
		return wrap.Errorf(err, "failed to write output to %s", humanReadableFilename(filename))
	}

	// Write the buffered output to the provided input file (override in place)
	bufWriter.Flush()
	if err := ioutil.WriteFile(filename, buf.Bytes(), 0644); err != nil {
		return wrap.Errorf(err, "failed to overwrite %s in place", humanReadableFilename(filename))
	}

	return nil
}

func (w *OutputWriter) write(writer io.Writer, filename string) error {
	inputFile, err := ytbx.LoadFile(filename)
	if err != nil {
		return wrap.Errorf(err, "failed to load input from %s", humanReadableFilename(filename))
	}

	for _, document := range inputFile.Documents {
		if w.Restructure {
			ytbx.RestructureObject(document)
		}

		switch {
		case w.PlainMode && w.OutputStyle == "json":
			output, err := neat.NewOutputProcessor(false, false, &neat.DefaultColorSchema).ToCompactJSON(document)
			if err != nil {
				return err
			}
			fmt.Fprintf(writer, "%s\n", output)

		case w.PlainMode && w.OutputStyle == "yaml":
			fmt.Fprintln(writer, "---")
			encoder := yamlv3.NewEncoder(writer)
			encoder.SetIndent(2)

			if err := encoder.Encode(document); err != nil {
				return err
			}

			if err := encoder.Close(); err != nil {
				return err
			}

		case w.OutputStyle == "json":
			output, err := neat.NewOutputProcessor(!w.OmitIndentHelper, true, &neat.DefaultColorSchema).ToJSON(document)
			if err != nil {
				return err
			}
			fmt.Fprintf(writer, "%s\n", output)

		case w.OutputStyle == "yaml":
			output, err := neat.NewOutputProcessor(!w.OmitIndentHelper, true, &neat.DefaultColorSchema).ToYAML(document)
			if err != nil {
				return err
			}
			fmt.Fprintf(writer, "%s\n", output)
		}
	}

	return nil
}

func writeReport(cmd *cobra.Command, report dyff.Report) error {
	var reportWriter dyff.ReportWriter
	switch strings.ToLower(reportOptions.style) {
	case "human", "bosh":
		reportWriter = &dyff.HumanReport{
			Report:               report,
			DoNotInspectCerts:    reportOptions.doNotInspectCerts,
			NoTableStyle:         reportOptions.noTableStyle,
			OmitHeader:           reportOptions.omitHeader,
			UseGoPatchPaths:      reportOptions.useGoPatchPaths,
			MinorChangeThreshold: 0.1,
		}

	case "brief", "short", "summary":
		reportWriter = &dyff.BriefReport{
			Report: report,
		}

	default:
		return wrap.Errorf(
			fmt.Errorf(cmd.UsageString()),
			"unknown output style %s", reportOptions.style,
		)
	}

	if err := reportWriter.WriteReport(os.Stdout); err != nil {
		return wrap.Errorf(err, "failed to print report")
	}

	// If configured, make sure `dyff` exists with an exit status
	if reportOptions.exitWithCode {
		switch len(report.Diffs) {
		case 0:
			return ExitCode{Value: 0}

		default:
			return ExitCode{Value: 1}
		}
	}

	return nil
}
