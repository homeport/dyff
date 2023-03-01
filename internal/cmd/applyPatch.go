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
	"github.com/gonvenience/ytbx"
	"github.com/homeport/dyff/pkg/dyff"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	yamlv3 "gopkg.in/yaml.v3"
	"io"
	"os"
)

type applyPatchCmdOptions struct {
	output string
	indent uint
}

var applyPatchCmdSettings applyPatchCmdOptions

// yamlCmd represents the yaml command
var applyPatchCmd = &cobra.Command{
	Use:     "apply-patch [flags] <patch-file-location> <input-file or '-'>",
	Aliases: []string{"apply", "ap"},
	Args:    cobra.MinimumNArgs(2),
	Short:   "Apply a YAML patch produced by 'between' to a file (use '-' for stdin)",
	Long:    `Applies a YAML patch created by the 'between' command (see between help for details) to a YAML file or stdin and writes to an output file or stdout.`,

	RunE: func(cmd *cobra.Command, args []string) error {
		patchb, err := os.ReadFile(args[0])
		if err != nil {
			return errors.Wrap(err, "error reading patch file")
		}

		var patch []dyff.PatchOp
		err = yamlv3.Unmarshal(patchb, &patch)
		if err != nil {
			return errors.Wrap(err, "error unmarshaling patch (is patch file corrupted?)")
		}

		var input []byte
		if args[1] == "-" {
			in, err := io.ReadAll(os.Stdin)
			if err != nil {
				return errors.Wrap(err, "error reading from stdin")
			}
			input = in
		} else {
			in, err := os.ReadFile(args[1])
			if err != nil {
				return errors.Wrap(err, "error reading input file")
			}
			input = in
		}

		inputdocs, err := ytbx.LoadYAMLDocuments(input)
		if err != nil {
			return errors.Wrap(err, "error unmarshaling input (is it valid YAML?)")
		}

		if len(inputdocs) == 0 {
			return fmt.Errorf("no YAML documents found in input")
		}

		if len(inputdocs) > 1 {
			fmt.Fprintf(os.Stderr, "warning: multiple yaml docs found in input, applying patch to the first doc only")
		}

		err = dyff.ApplyPatch(inputdocs[0], patch)
		if err != nil {
			return errors.Wrap(err, "error applying patch")
		}

		var b bytes.Buffer
		bw := bufio.NewWriter(&b)
		yenc := yamlv3.NewEncoder(bw)
		yenc.SetIndent(int(applyPatchCmdSettings.indent))

		err = yenc.Encode(inputdocs[0])
		if err != nil {
			return errors.Wrap(err, "error marshaling yaml output")
		}
		err = yenc.Close()
		if err != nil {
			return errors.Wrap(err, "error flushing yaml encoder")
		}
		err = bw.Flush()
		if err != nil {
			return errors.Wrap(err, "error flushing output buffer")
		}

		out := b.Bytes()

		if applyPatchCmdSettings.output == "-" {
			fmt.Printf("%s\v", out)
		} else {
			err := os.WriteFile(applyPatchCmdSettings.output, out, 0644)
			if err != nil {
				return errors.Wrap(err, "error writing output file")
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(applyPatchCmd)

	applyPatchCmd.Flags().SortFlags = false

	applyPatchCmd.Flags().StringVarP(&applyPatchCmdSettings.output, "output", "o", "-", "output path (use - for stdout)")
	applyPatchCmd.Flags().UintVarP(&applyPatchCmdSettings.indent, "indent", "i", 2, "YAML output indent size (in spaces)")
}
