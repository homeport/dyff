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
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gonvenience/bunt"
	"github.com/gonvenience/neat"
	"github.com/gonvenience/term"
	"github.com/gonvenience/wrap"
	"github.com/homeport/dyff/pkg/v1/dyff"
	"github.com/homeport/ytbx/pkg/v1/ytbx"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

// colormode is used by CLI parser to store user input for further internal processing into the proper value
var colormode string

// truecolormode is used by the CLI flag processing routines to store the user preference for true color usage
var truecolormode string

// debugMode set to true will set-up the logging package to use the debug logger
var debugMode bool

// OutputWriter encapsulates the required fields to define the look and feel of
// the output
type OutputWriter struct {
	PlainMode        bool
	Restructure      bool
	OmitIndentHelper bool
	OutputStyle      string
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "dyff",
	Long: `
δyƒƒ /ˈdʏf/ - a diff tool for YAML files, and sometimes JSON. Also, It
can transform YAML to JSON, and vice versa. The order of keys in hashes
is preserved during the conversion.
`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		if bunt.ColorSetting, err = parseSetting(colormode); err != nil {
			return wrap.Errorf(err, "invalid color setting '%s'", colormode)
		}

		if bunt.TrueColorSetting, err = parseSetting(truecolormode); err != nil {
			return wrap.Errorf(err, "invalid true color setting '%s'", truecolormode)
		}

		return nil
	},
}

// ResetSettings resets command settings to default. This is only required by
// the test suite to make sure that the flag parsing works correctly.
func ResetSettings() {
	betweenCmdSettings = struct {
		style                    string
		swap                     bool
		noTableStyle             bool
		doNotInspectCerts        bool
		exitWithCount            bool
		translateListToDocuments bool
		chroot                   string
		chrootFrom               string
		chrootTo                 string
	}{
		style: defaultOutputStyle,
	}

	yamlCmdSettings = struct {
		plainMode        bool
		restructure      bool
		omitIndentHelper bool
		inplace          bool
	}{}

	jsonCmdSettings = struct {
		plainMode        bool
		restructure      bool
		omitIndentHelper bool
		inplace          bool
	}{}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initSettings)

	rootCmd.SilenceErrors = true
	rootCmd.SilenceUsage = true
	rootCmd.Flags().SortFlags = false
	rootCmd.PersistentFlags().SortFlags = false

	rootCmd.PersistentFlags().StringVarP(&colormode, "color", "c", "auto", "specify color usage: on, off, or auto")
	rootCmd.PersistentFlags().StringVarP(&truecolormode, "truecolor", "t", "auto", "specify true color usage: on, off, or auto")
	rootCmd.PersistentFlags().IntVarP(&term.FixedTerminalWidth, "fixed-width", "w", -1, "disable terminal width detection and use provided fixed value")
	rootCmd.PersistentFlags().BoolVarP(&ytbx.PreserveKeyOrderInJSON, "preserve-key-order-in-json", "k", false, "use ordered keys during JSON decoding (non standard behavior)")
	rootCmd.PersistentFlags().BoolVarP(&debugMode, "debug", "d", false, "enable debug mode")
}

func parseSetting(setting string) (bunt.SwitchState, error) {
	switch strings.ToLower(setting) {
	case "auto":
		return bunt.AUTO, nil

	case "off", "no", "false":
		return bunt.OFF, nil

	case "on", "yes", "true":
		return bunt.ON, nil

	default:
		return bunt.OFF, fmt.Errorf("invalid state '%s' used, supported modes are: auto, on, or off", setting)
	}
}

func initSettings() {
	if debugMode {
		dyff.SetLoggingLevel(dyff.DEBUG)
	}
}

// WriteToStdout is a convenience function to write the content of the documents
// stored in the provided input file to the standard output
func (w *OutputWriter) WriteToStdout(filename string) error {
	if err := w.write(os.Stdout, filename); err != nil {
		return wrap.Errorf(err, "failed to write output _%s_", filename)
	}

	return nil
}

// WriteInplace writes the content of the documents stored in the provided input
// file to the file itself overwriting the conent in place.
func (w *OutputWriter) WriteInplace(filename string) error {
	var buf bytes.Buffer
	bufWriter := bufio.NewWriter(&buf)

	// Force plain mode to make sure there are no ANSI sequences
	w.PlainMode = true
	if err := w.write(bufWriter, filename); err != nil {
		return wrap.Errorf(err, "failed to write output _%s_", filename)
	}

	// Write the buffered output to the provided input file (override in place)
	bufWriter.Flush()
	if err := ioutil.WriteFile(filename, buf.Bytes(), 0644); err != nil {
		return wrap.Errorf(err, "failed to overwrite file _%s_ in place", filename)
	}

	return nil
}

func (w *OutputWriter) write(writer io.Writer, filename string) error {
	inputFile, err := ytbx.LoadFile(filename)
	if err != nil {
		return wrap.Errorf(err, "failed to load input file _%s_", filename)
	}

	for _, document := range inputFile.Documents {
		if w.Restructure {
			document = ytbx.RestructureObject(document)
		}

		switch {
		case w.PlainMode && w.OutputStyle == "json":
			output, err := neat.NewOutputProcessor(false, false, &neat.DefaultColorSchema).ToCompactJSON(document)
			if err != nil {
				return err
			}
			fmt.Fprintf(writer, "%s\n", output)

		case w.PlainMode && w.OutputStyle == "yaml":
			output, err := yaml.Marshal(document)
			if err != nil {
				return err
			}
			fmt.Fprintf(writer, "---\n%s\n", string(output))

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
			fmt.Fprintf(writer, "---\n%s\n", output)
		}
	}

	return nil
}
