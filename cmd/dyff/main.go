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
	"fmt"
	"os"
	"strings"

	"github.com/gonvenience/bunt"
	"github.com/gonvenience/neat"
	"github.com/gonvenience/wrap"

	"github.com/homeport/dyff/internal/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		var (
			headline string
			content  string
			code     int
		)

		switch typed := err.(type) {
		case wrap.ContextError:
			headline = bunt.Sprintf("*Error:* _%s_", typed.Context())
			content = typed.Cause().Error()
			code = 1

		case cmd.ExitCode:
			code = typed.Value

		case error:
			headline = "Error occurred"
			content = err.Error()
			code = 1

		default:
			headline = "Error occurred"
			content = fmt.Sprint(err)
			code = 1
		}

		if len(headline) > 0 {
			neat.Box(os.Stderr,
				headline, strings.NewReader(content),
				neat.HeadlineColor(bunt.Coral),
				neat.ContentColor(bunt.DimGray),
			)
		}

		os.Exit(code)
	}
}
