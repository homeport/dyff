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
	"os"
	"strings"

	"github.com/gonvenience/wrap"
	"github.com/gonvenience/ytbx"
	"github.com/spf13/cobra"

	. "github.com/gonvenience/bunt"
	"github.com/homeport/dyff/pkg/dyff"
)

const (
	// Arbitrary length of the git SHA1 in short format
	// It does not calculate git-sha1
	GIT_SHA1_LENGTH = 7
)

// betweenCmd represents the between command
var gitDiffCmd = &cobra.Command{
	Use:   "git-diff name infile1 infile1-sha1 infile1-mode infile2 infile2-sha1 infile2-mode [ rename-to ]",
	Short: "Supposed to be used as a git diff command for yaml files",
	Long: `Supposed to be used as a git diff command for yaml files
Example gitconfig:

[diff "dyff"]
  command =  dyff git-diff

Example in gitattributes file:
*.yaml diff=dyff
*.yml diff=dyff

https://git-scm.com/docs/git-difftool for more information about gif diff
`,
	Args: cobra.RangeArgs(7, 9),
	RunE: func(cmd *cobra.Command, args []string) error {

		fmt.Print(generateGitDiffContext(args))
		from, to, err := ytbx.LoadFiles(args[1], args[4])
		if err != nil {
			return wrap.Errorf(err, "failed to load input files")
		}

		report, err := dyff.CompareInputFiles(from, to,
			dyff.IgnoreOrderChanges(reportOptions.ignoreOrderChanges),
		)
		if err != nil {
			return wrap.Errorf(err, "failed to compare input files")
		}

		reportOptions.omitHeader = true
		return writeReport(cmd, report)
	},
}

func init() {
	rootCmd.AddCommand(gitDiffCmd)

	GIT_CONFIG_PARAMETERS := os.Getenv("GIT_CONFIG_PARAMETERS")
	enableColor := strings.Contains(GIT_CONFIG_PARAMETERS, "'color.ui=always'")
	enableColor = enableColor || strings.Contains(GIT_CONFIG_PARAMETERS, "'color.diff=always'")
	// color.diff=auto|true|false|xxx will disable color and default to dyff auto
	// by checking <stdout-is-tty>
	if !strings.Contains(GIT_CONFIG_PARAMETERS, "'color.diff=always'") && strings.Contains(GIT_CONFIG_PARAMETERS, "'color.diff=") {
		enableColor = false
	}
	if enableColor {
		SetColorSettings(ON, ON)
	}
}

// https://github.com/git/git/blob/6cd33dceed60949e2dbc32e3f0f5e67c4c882e1e/diff.c#L6221
// https://github.com/git/git/blob/a68dfadae5e95c7f255cf38c9efdcbc2e36d1931/diff.c#L4244-L4250
/* An external diff command takes:
 *
 * diff-cmd name infile1 infile1-sha1 infile1-mode \
 *               infile2 infile2-sha1 infile2-mode [ rename-to ]
 *
 */
func generateGitDiffContext(args []string) string {

	// Example output format from git diff:
	/*
		--------------------------------------------------
		diff --git a/.github/dependabot.yml b/.github/dependabot.yml
		new file mode 100644
		index 0000000..72e0f35
		--- /dev/null
		+++ b/.github/dependabot.yml

		--------------------------------------------------
		diff --git a/.github/dependabot.yml b/.github/dependabot.yml
		deleted file mode 100644
		index 72e0f35..0000000
		--- a/.github/dependabot.yml
		+++ /dev/null

		--------------------------------------------------
		diff --git a/.github/dependabot.yml b/.github/dependabot.yml
		index 72e0f35..8f2b7dc 100644
		--- a/.github/dependabot.yml
		+++ b/.github/dependabot.yml

		--------------------------------------------------
		diff --git a/.github/dependabot.yml b/.github/dependabot.yml
		old mode 100644
		new mode 100755
		index 72e0f35..8f2b7dc
		--- a/.github/dependabot.yml
		+++ b/.github/dependabot.yml
	*/

	var sbGitHashCtx strings.Builder
	var gitname, infile1, infile1sha1, infile1mode, infile2, infile2sha2, infile2mode string = args[0], args[1], args[2], args[3], args[4], args[5], args[6]
	sbGitHashCtx.WriteString(fmt.Sprintf("dyff --git a/%s b/%s\n", gitname, gitname))
	if infile1 == os.DevNull {
		sbGitHashCtx.WriteString(fmt.Sprintf("new file mode %s\n", infile2mode))
		// caculate git object hash of the file
		infile1sha1 = "0000000"
	} else if infile2 == os.DevNull {
		sbGitHashCtx.WriteString(fmt.Sprintf("deleted file mode %s\n", infile1mode))
		infile2sha2 = "0000000"
	} else if infile1mode != infile2mode {
		sbGitHashCtx.WriteString(fmt.Sprintf("old mode %s\n", infile1mode))
		sbGitHashCtx.WriteString(fmt.Sprintf("new mode %s\n", infile2mode))
	}

	// file renamed ? Add rename info passed from git diff
	if len(args) >= 9 {
		sbGitHashCtx.WriteString(args[8])
		// renamed and no changes
		if strings.Contains(args[8], "similarity index 100") {
			return sbGitHashCtx.String()
		}
	} else {
		sbGitHashCtx.WriteString(fmt.Sprintf("index %s..%s\n", formatGitSHA1(infile1sha1), formatGitSHA1(infile2sha2)))
	}

	if infile1 == os.DevNull {
		sbGitHashCtx.WriteString("--- /dev/null\n")
	} else {
		sbGitHashCtx.WriteString(fmt.Sprintf("--- a/%s\n", gitname))
	}
	if infile2 == os.DevNull {
		sbGitHashCtx.WriteString("+++ /dev/null\n")
	} else if len(args) >= 9 { // file renamed
		sbGitHashCtx.WriteString(fmt.Sprintf("+++ b/%s\n", args[7]))
	} else {
		sbGitHashCtx.WriteString(fmt.Sprintf("+++ b/%s\n", gitname))
	}
	return sbGitHashCtx.String()
}

func formatGitSHA1(gitSHA1 string) string {
	if len(gitSHA1) > GIT_SHA1_LENGTH {
		return gitSHA1[:GIT_SHA1_LENGTH]
	}
	return "0000000"
}
