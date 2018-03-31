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
	"os"

	"github.com/HeavyWombat/color"
	"github.com/HeavyWombat/yaml"
	"github.com/spf13/cobra"
)

// NoColor is the gobal switch to decide whether strings should be colored in the output
var NoColor = false

// Debug is the global switch to enable debug output
var Debug = false

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dyff",
	Short: "A diff tool for YAMLs",
	Long: `
A diff tool for YAMLs, and sometimes JSONs. It also comes with conversion
capabilities to transform YAML to JSON, or JSON to YAML.
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initSettings)

	// Here you will define your flags and configuration settings. Cobra supports
	// persistent flags, which, if defined here, will be global for your
	// application.
	rootCmd.PersistentFlags().BoolVar(&NoColor, "no-color", false, "Disable colors in output")
	rootCmd.PersistentFlags().BoolVar(&Debug, "debug", false, "Disable colors in output")
}

func initSettings() {
	if NoColor {
		color.NoColor = true
		yaml.HighlightKeys = false
	}
}
