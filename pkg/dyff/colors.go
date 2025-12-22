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

func green(text string) string {
	return colored(additionGreen, text)
}

func greenf(format string, a ...interface{}) string {
	return coloredf(additionGreen, format, a...)
}

func red(text string) string {
	return colored(removalRed, text)
}

func redf(format string, a ...interface{}) string {
	return coloredf(removalRed, format, a...)
}

func yellowf(format string, a ...interface{}) string {
	return coloredf(modificationYellow, format, a...)
}

func lightgreen(text string) string {
	return colored(bunt.LightGreen, text)
}

func lightred(text string) string {
	return colored(bunt.LightSalmon, text)
}

func dimgray(text string) string {
	return colored(bunt.DimGray, text)
}

func bold(text string) string {
	return bunt.Style(text,
		bunt.EachLine(),
		bunt.Bold(),
	)
}

func italic(text string) string {
	return bunt.Style(text,
		bunt.EachLine(),
		bunt.Italic(),
	)
}

func colored(color colorful.Color, text string) string {
	return bunt.Style(text,
		bunt.EachLine(),
		bunt.Foreground(color),
	)
}

func coloredf(color colorful.Color, format string, a ...interface{}) string {
	return bunt.Style(render(format, a...),
		bunt.EachLine(),
		bunt.Foreground(color),
	)
}
