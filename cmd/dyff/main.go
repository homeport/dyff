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

package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/gonvenience/bunt"
	"github.com/gonvenience/neat"

	"github.com/homeport/dyff/internal/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		switch err := err.(type) {
		case cmd.ExitCode:
			var headline, content string

			if unwrapped := errors.Unwrap(err.Cause()); unwrapped != nil {
				headline = strings.Split(err.Error(), ":")[0]
				content = unwrapped.Error()

			} else {
				headline = "Error occurred"
				content = err.Cause().Error()
			}

			fmt.Fprint(
				os.Stderr,
				neat.ContentBox(
					headline,
					content,
					neat.HeadlineColor(bunt.Coral),
					neat.ContentColor(bunt.DimGray),
				),
			)

			os.Exit(err.Value())

		default: // fail safe for somehow an non exit code error slips through
			fmt.Fprint(os.Stderr, err.Error())
			os.Exit(1)
		}
	}
}
