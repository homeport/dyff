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
	"os"
	"path/filepath"

	"github.com/gonvenience/bunt"
	"github.com/gonvenience/term"
	"github.com/gonvenience/ytbx"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// ExitCode is just a way to transport the exit code to the main package
type ExitCode struct {
	Value int
	Cause error
}

func (e ExitCode) Error() string {
	if e.Cause != nil {
		return e.Cause.Error()
	}

	return ""
}

var name = func() string {
	ep, err := os.Executable()
	if err != nil {
		return "dyff"
	}

	return filepath.Base(ep)
}()

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:           name,
	SilenceErrors: true,
	SilenceUsage:  true,
	Long: `
δyƒƒ /ˈdʏf/ - a diff tool for YAML files, and sometimes JSON. Also, It
can transform YAML to JSON, and vice versa. The order of keys in hashes
is preserved during the conversion.
`,
}

// ResetSettings resets command settings to default. This is only required by
// the test suite to make sure that the flag parsing works correctly.
func ResetSettings() {
	reportOptions = reportConfig{style: defaultOutputStyle}
	betweenCmdSettings = betweenCmdOptions{}
	yamlCmdSettings = yamlCmdOptions{}
	jsonCmdSettings = jsonCmdOptions{}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	// Run with `--kubectl-external-diff` explicitly to trigger argument reordering.
	// Example when running:
	// KUBECTL_EXTERNAL_DIFF="dyff between --kubectl-external-diff" kubectl diff -f xxx/yyy/zzz.yaml
	// kubectl will invoke dyff with the arguments: (injecting two directores as arg[1] and arg[2])
	// dyff tmp/fileA tmp/fileB between --kubectl-external-diff

	if indexSlice(os.Args, "--kubectl-external-diff") > 1 && indexSlice(os.Args, "between") > 1 {
		ytbx.PreserveKeyOrderInJSON = true
		// Rearrange the arguments to match `dyff between --flags from to` to
		// mitigate an issue in `kubectl`, which puts the `from` and `to` at
		// the second and third position in the command arguments.

		// only check arg[1] and arg[2], in case `dyff` or `between` are directories
		for i := 1; i <= 2; i++ {
			if info, err := os.Stat(os.Args[i]); err == nil && info.IsDir() {
			} else {
				return ExitCode{
					Value: 255,
					Cause: errors.Errorf("%s is not a directory? is this invoked by kubectl diff?", os.Args[i]),
				}
			}
		}
		var args []string
		args = append(args, os.Args[0])
		args = append(args, os.Args[3:]...)
		args = append(args, os.Args[1], os.Args[2])
		os.Args = args
	}

	if err := rootCmd.Execute(); err != nil {
		// Special case ExitCode, which means that we will exit immediately
		// with the given exit code
		switch err.(type) {
		case ExitCode:
			return err
		}

		// In any other case, create a default ExitCode with `error` value
		return ExitCode{
			Value: 255,
			Cause: err,
		}
	}

	return nil
}

func init() {
	rootCmd.Flags().SortFlags = false
	rootCmd.PersistentFlags().SortFlags = false

	rootCmd.PersistentFlags().VarP(&bunt.ColorSetting, "color", "c", "specify color usage: on, off, or auto")
	rootCmd.PersistentFlags().VarP(&bunt.TrueColorSetting, "truecolor", "t", "specify true color usage: on, off, or auto")
	rootCmd.PersistentFlags().IntVarP(&term.FixedTerminalWidth, "fixed-width", "w", -1, "disable terminal width detection and use provided fixed value")
	rootCmd.PersistentFlags().BoolVarP(&ytbx.PreserveKeyOrderInJSON, "preserve-key-order-in-json", "k", false, "use ordered keys during JSON decoding (non standard behavior)")
}

func indexSlice(s []string, str string) int {
	for i, v := range s {
		if v == str {
			return i
		}
	}
	return -1
}
