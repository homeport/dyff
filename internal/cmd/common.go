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
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/caarlos0/env/v11"
	"github.com/gonvenience/bunt"
	"github.com/gonvenience/neat"
	"github.com/gonvenience/ytbx"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	yamlv3 "go.yaml.in/yaml/v3"

	"github.com/homeport/dyff/pkg/dyff"
)

type reportConfig struct {
	Style          string `envDefault:"human"`
	UseIndentLines bool   `envDefault:"true"`

	IgnoreOrderChanges      bool `envDefault:"false"`
	IgnoreWhitespaceChanges bool `envDefault:"false"`
	IgnoreValueChanges      bool `envDefault:"false"`
	FormatStrings           bool `envDefault:"true"`
	DetectRenames           bool `envDefault:"true"`

	NoTableStyle          bool    `envDefault:"false"`
	DoNotInspectCerts     bool    `envDefault:"false"`
	UseGoPatchPaths       bool    `envDefault:"false"`
	MinorChangeThreshold  float64 `envDefault:"0.1"`
	MultilineContextLines int     `envDefault:"4"`

	KubernetesEntityDetection bool `envDefault:"true"`
	AdditionalIdentifiers     []string
	Filters                   []string
	Excludes                  []string
	FilterRegexps             []string
	ExcludeRegexps            []string

	ExitWithCode bool `envDefault:"false"`
	OmitHeader   bool `envDefault:"false"`
}

func initReportConfig() reportConfig {
	return env.Must(env.ParseAsWithOptions[reportConfig](env.Options{
		Prefix:                "DYFF_",
		UseFieldNameByDefault: true,
	}))
}

var reportOptions = initReportConfig()

func flagSet(name string, f ...func(*pflag.FlagSet)) *pflag.FlagSet {
	var flatSet = pflag.NewFlagSet(name, pflag.ExitOnError)
	flatSet.SortFlags = false

	for _, fn := range f {
		fn(flatSet)
	}

	return flatSet
}

func reportOptionsFlags() []*pflag.FlagSet {
	return []*pflag.FlagSet{
		flagSet("Output Preferences", func(fs *pflag.FlagSet) {
			fs.StringVarP(&reportOptions.Style, "output", "o", reportOptions.Style, "specify the output style, supported styles: human, brief, github, gitlab, gitea")
			fs.BoolVar(&reportOptions.UseIndentLines, "use-indent-lines", reportOptions.UseIndentLines, "use indent lines in the output")
		}),

		flagSet("Compare Options", func(fs *pflag.FlagSet) {
			fs.BoolVarP(&reportOptions.IgnoreOrderChanges, "ignore-order-changes", "i", reportOptions.IgnoreOrderChanges, "ignore order changes in lists")
			fs.BoolVar(&reportOptions.IgnoreWhitespaceChanges, "ignore-whitespace-changes", reportOptions.IgnoreWhitespaceChanges, "ignore leading or trailing whitespace changes")
			fs.BoolVarP(&reportOptions.IgnoreValueChanges, "ignore-value-changes", "v", reportOptions.IgnoreValueChanges, "exclude changes in values")
			fs.BoolVar(&reportOptions.DetectRenames, "detect-renames", reportOptions.DetectRenames, "enable detection for renames (document level for Kubernetes resources)")
			fs.BoolVar(&reportOptions.FormatStrings, "format-strings", reportOptions.FormatStrings, "format strings (i.e. inline JSON) before comparison to avoid formatting differences")
		}),

		flagSet("Human Output Preferences", func(fs *pflag.FlagSet) {
			fs.BoolVarP(&reportOptions.NoTableStyle, "no-table-style", "l", reportOptions.NoTableStyle, "do not place blocks next to each other, always use one row per text block")
			fs.BoolVarP(&reportOptions.DoNotInspectCerts, "no-cert-inspection", "x", reportOptions.DoNotInspectCerts, "disable x509 certificate inspection, compare as raw text")
			fs.BoolVarP(&reportOptions.UseGoPatchPaths, "use-go-patch-style", "g", reportOptions.UseGoPatchPaths, "use Go-Patch style paths in outputs")
			fs.Float64VarP(&reportOptions.MinorChangeThreshold, "minor-change-threshold", "", reportOptions.MinorChangeThreshold, "minor change threshold")
			fs.IntVarP(&reportOptions.MultilineContextLines, "multi-line-context-lines", "", reportOptions.MultilineContextLines, "multi-line context lines")
		}),

		flagSet("Filter Options", func(fs *pflag.FlagSet) {
			fs.BoolVarP(&reportOptions.KubernetesEntityDetection, "detect-kubernetes", "", reportOptions.KubernetesEntityDetection, "detect kubernetes entities")
			fs.StringArrayVar(&reportOptions.AdditionalIdentifiers, "additional-identifier", reportOptions.AdditionalIdentifiers, "use additional identifier candidates in named entry lists")
			fs.StringSliceVar(&reportOptions.Filters, "filter", reportOptions.Filters, "filter reports to a subset of differences based on supplied arguments")
			fs.StringSliceVar(&reportOptions.Excludes, "exclude", reportOptions.Excludes, "exclude reports from a set of differences based on supplied arguments")
			fs.StringSliceVar(&reportOptions.FilterRegexps, "filter-regexp", reportOptions.FilterRegexps, "filter reports to a subset of differences based on supplied regular expressions")
			fs.StringSliceVar(&reportOptions.ExcludeRegexps, "exclude-regexp", reportOptions.ExcludeRegexps, "exclude reports from a set of differences based on supplied regular expressions")
		}),

		flagSet("General Options", func(fs *pflag.FlagSet) {
			fs.BoolVarP(&reportOptions.OmitHeader, "omit-header", "b", reportOptions.OmitHeader, "omit the dyff summary header")
			fs.BoolVarP(&reportOptions.ExitWithCode, "set-exit-code", "s", reportOptions.ExitWithCode, "set program exit code, with 0 meaning no difference, 1 for differences detected, and 255 for program error")

			// Deprecated
			fs.BoolVar(&reportOptions.ExitWithCode, "set-exit-status", reportOptions.ExitWithCode, "set program exit code, with 0 meaning no difference, 1 for differences detected, and 255 for program error")
			_ = fs.MarkDeprecated("set-exit-status", "use --set-exit-code instead")
		}),
	}
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
		return bunt.Errorf("failed to write output to _*stdout*_: %w", err)
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
		return fmt.Errorf("failed to write output to %s: %w", humanReadableFilename(filename), err)
	}

	// Write the buffered output to the provided input file (override in place)
	_ = bufWriter.Flush()
	if err := os.WriteFile(filename, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to overwrite %s in place: %w", humanReadableFilename(filename), err)
	}

	return nil
}

func (w *OutputWriter) write(writer io.Writer, filename string) error {
	inputFile, err := ytbx.LoadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to load input from %s: %w", humanReadableFilename(filename), err)
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
			_, _ = fmt.Fprintln(writer, output)

		case w.PlainMode && w.OutputStyle == "yaml":
			_, _ = fmt.Fprintln(writer, "---")
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
			_, _ = fmt.Fprintln(writer, output)

		case w.OutputStyle == "yaml":
			output, err := neat.NewOutputProcessor(!w.OmitIndentHelper, true, &neat.DefaultColorSchema).ToYAML(document)
			if err != nil {
				return err
			}
			_, _ = fmt.Fprintln(writer, output)
		}
	}

	return nil
}

func writeReport(cmd *cobra.Command, report dyff.Report) error {
	var reportWriter dyff.ReportWriter
	switch strings.ToLower(reportOptions.Style) {
	case "human", "bosh":
		reportWriter = &dyff.HumanReport{
			Report:                report,
			Indent:                2,
			UseIndentLines:        reportOptions.UseIndentLines,
			DoNotInspectCerts:     reportOptions.DoNotInspectCerts,
			NoTableStyle:          reportOptions.NoTableStyle,
			OmitHeader:            reportOptions.OmitHeader,
			UseGoPatchPaths:       reportOptions.UseGoPatchPaths,
			MinorChangeThreshold:  reportOptions.MinorChangeThreshold,
			MultilineContextLines: reportOptions.MultilineContextLines,
			PrefixMultiline:       false,
		}

	case "github", "linguist":
		reportWriter = &dyff.DiffSyntaxReport{
			PathPrefix:            "@@",
			RootDescriptionPrefix: "#",
			ChangeTypePrefix:      "!",
			HumanReport: dyff.HumanReport{
				Report:                report,
				Indent:                0,
				UseIndentLines:        reportOptions.UseIndentLines,
				DoNotInspectCerts:     reportOptions.DoNotInspectCerts,
				NoTableStyle:          true,
				OmitHeader:            true,
				UseGoPatchPaths:       reportOptions.UseGoPatchPaths,
				MinorChangeThreshold:  reportOptions.MinorChangeThreshold,
				MultilineContextLines: reportOptions.MultilineContextLines,
				PrefixMultiline:       true,
			},
		}

	case "gitlab", "rogue":
		reportWriter = &dyff.DiffSyntaxReport{
			PathPrefix:            "=",
			RootDescriptionPrefix: "=",
			ChangeTypePrefix:      "#",
			HumanReport: dyff.HumanReport{
				Report:                report,
				Indent:                0,
				UseIndentLines:        reportOptions.UseIndentLines,
				DoNotInspectCerts:     reportOptions.DoNotInspectCerts,
				NoTableStyle:          true,
				OmitHeader:            true,
				UseGoPatchPaths:       reportOptions.UseGoPatchPaths,
				MinorChangeThreshold:  reportOptions.MinorChangeThreshold,
				MultilineContextLines: reportOptions.MultilineContextLines,
				PrefixMultiline:       true,
			},
		}

	case "gitea", "forgejo":
		reportWriter = &dyff.DiffSyntaxReport{
			PathPrefix:            "@@",
			RootDescriptionPrefix: "=",
			ChangeTypePrefix:      "!",
			HumanReport: dyff.HumanReport{
				Report:                report,
				Indent:                0,
				UseIndentLines:        reportOptions.UseIndentLines,
				DoNotInspectCerts:     reportOptions.DoNotInspectCerts,
				NoTableStyle:          true,
				OmitHeader:            true,
				UseGoPatchPaths:       reportOptions.UseGoPatchPaths,
				MinorChangeThreshold:  reportOptions.MinorChangeThreshold,
				MultilineContextLines: reportOptions.MultilineContextLines,
				PrefixMultiline:       true,
			},
		}

	case "brief", "short", "summary":
		reportWriter = &dyff.BriefReport{
			Report: report,
		}

	default:
		return fmt.Errorf("unknown output style %s: %w", reportOptions.Style, errors.New(cmd.UsageString()))
	}

	if err := reportWriter.WriteReport(os.Stdout); err != nil {
		return fmt.Errorf("failed to print report: %w", err)
	}

	// If configured, make sure `dyff` exists with an exit status
	if reportOptions.ExitWithCode {
		switch len(report.Diffs) {
		case 0:
			return errorWithExitCode{value: 0}

		default:
			return errorWithExitCode{value: 1}
		}
	}

	return nil
}
