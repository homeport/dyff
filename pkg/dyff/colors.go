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

package dyff

import (
	"fmt"

	"github.com/mgutz/ansi"
)

func render(format string, a ...interface{}) string {
	if len(a) == 0 {
		return format
	}

	return fmt.Sprintf(format, a...)
}

func green(text string) string {
	return ansi.ColorFunc("green")(text)
}

func greenf(format string, a ...interface{}) string {
	return ansi.ColorFunc("green")(render(format, a...))
}

func red(text string) string {
	return ansi.ColorFunc("red")(text)
}

func redf(format string, a ...interface{}) string {
	return ansi.ColorFunc("red")(render(format, a...))
}

func yellowf(format string, a ...interface{}) string {
	return ansi.ColorFunc("yellow")(render(format, a...))
}

func lightgreen(text string) string {
	return ansi.ColorFunc("green+h")(text)
}

func lightred(text string) string {
	return ansi.ColorFunc("red+h")(text)
}

func dimgray(text string) string {
	return ansi.ColorFunc("white+d")(text)
}

func boldgreen(text string) string {
	return ansi.ColorFunc("green+b")(text)
}

func boldred(text string) string {
	return ansi.ColorFunc("red+b")(text)
}

func italic(text string) string {
	return text
}
