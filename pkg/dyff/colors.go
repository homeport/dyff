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

	"github.com/gonvenience/bunt"
	"github.com/lucasb-eyer/go-colorful"
)

var (
	additionGreen      = color("#58BF38")
	modificationYellow = color("#C7C43F")
	removalRed         = color("#B9311B")
)

func color(hex string) colorful.Color {
	color, _ := colorful.Hex(hex)
	return color
}

func render(format string, a ...interface{}) string {
	if len(a) == 0 {
		return format
	}

	return fmt.Sprintf(format, a...)
}

func green(format string, a ...interface{}) string {
	return colored(additionGreen, render(format, a...))
}

func red(format string, a ...interface{}) string {
	return colored(removalRed, render(format, a...))
}

func yellow(format string, a ...interface{}) string {
	return colored(modificationYellow, render(format, a...))
}

func lightgreen(format string, a ...interface{}) string {
	return colored(bunt.LightGreen, render(format, a...))
}

func lightred(format string, a ...interface{}) string {
	return colored(bunt.LightSalmon, render(format, a...))
}

func dimgray(format string, a ...interface{}) string {
	return colored(bunt.DimGray, render(format, a...))
}

func bold(format string, a ...interface{}) string {
	return bunt.Style(
		fmt.Sprintf(format, a...),
		bunt.EachLine(),
		bunt.Bold(),
	)
}

func italic(format string, a ...interface{}) string {
	return bunt.Style(
		render(format, a...),
		bunt.EachLine(),
		bunt.Italic(),
	)
}

func colored(color colorful.Color, format string, a ...interface{}) string {
	return bunt.Style(
		render(format, a...),
		bunt.EachLine(),
		bunt.Foreground(color),
	)
}

// Theme-aware color methods

// green returns a green-colored string using the theme's addition color
func (theme *ColorTheme) green(format string, a ...interface{}) string {
	if theme == nil {
		theme = DefaultColorTheme()
	}
	return colored(theme.Addition, render(format, a...))
}

// red returns a red-colored string using the theme's removal color
func (theme *ColorTheme) red(format string, a ...interface{}) string {
	if theme == nil {
		theme = DefaultColorTheme()
	}
	return colored(theme.Removal, render(format, a...))
}

// yellow returns a yellow-colored string using the theme's modification color
func (theme *ColorTheme) yellow(format string, a ...interface{}) string {
	if theme == nil {
		theme = DefaultColorTheme()
	}
	return colored(theme.Modification, render(format, a...))
}
