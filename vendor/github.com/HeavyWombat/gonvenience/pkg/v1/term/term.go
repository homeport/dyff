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

/*
Package term contains convenience functions for terminal related tasks, for
example to get the current terminal width or a reasonable fallback.
*/
package term

import (
	"os"

	isatty "github.com/mattn/go-isatty"
	"golang.org/x/crypto/ssh/terminal"
)

// DefaultTerminalWidth is the default fallback terminal width.
const DefaultTerminalWidth = 80

// DefaultTerminalHeight is the default fallback terminal height.
const DefaultTerminalHeight = 25

// FixedTerminalWidth allows a manual fixed override of the terminal width.
var FixedTerminalWidth = -1

// FixedTerminalHeight allows a manual fixed override of the terminal height.
var FixedTerminalHeight = -1

// isDumbTerminal points to true if the current terminal has limited support for
// escape sequences, false otherwise, or nil if uninitialised.
var isDumbTerminal *bool

// isTerminal points to true if the current STDOUT stream writes to a terminal
// (not a redirect), false otherwise, or nil if uninitialised.
var isTerminal *bool

// isTrueColor points to true if the current terminal reports to support 24-bit
// colors, false otherwise, or nil if uninitialised.
var isTrueColor *bool

// GetTerminalWidth return the terminal width (available characters per line)
func GetTerminalWidth() int {
	width, _ := GetTerminalSize()
	return width
}

// GetTerminalHeight return the terminal height (available lines).
func GetTerminalHeight() int {
	_, height := GetTerminalSize()
	return height
}

// GetTerminalSize return the terminal size as a width and height tuple. In
// case the terminal size cannot be determined, a reasonable default is
// used: 80x25. A manual override is possible using FixedTerminalWidth
// and FixedTerminalHeight.
func GetTerminalSize() (int, int) {
	// Return user preference (explicit overwrite) of both width and height
	if FixedTerminalWidth > 0 && FixedTerminalHeight > 0 {
		return FixedTerminalWidth, FixedTerminalHeight
	}

	width, height, err := terminal.GetSize(int(os.Stdout.Fd()))

	switch {
	// Return default fallback value
	case err != nil:
		return DefaultTerminalWidth, DefaultTerminalHeight

	// Return user preference of width, actual value for height
	case FixedTerminalWidth > 0:
		return FixedTerminalWidth, height

	// Return user preference of height, actual value for width
	case FixedTerminalHeight > 0:
		return width, FixedTerminalHeight

	// Return actual determined values
	default:
		return width, height
	}
}

// IsDumbTerminal returns whether the current terminal has a limited feature set
func IsDumbTerminal() bool {
	if isDumbTerminal == nil {
		isTermDumbCheck := os.Getenv("TERM") == "dumb"
		isDumbTerminal = &isTermDumbCheck
	}

	return *isDumbTerminal
}

// IsTerminal returns whether this program runs in a terminal and not in a pipe
func IsTerminal() bool {
	if isTerminal == nil {
		isTerminalCheck := isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())
		isTerminal = &isTerminalCheck
	}

	return *isTerminal
}

// IsTrueColor returns whether the current terminal supports 24 bit colors
func IsTrueColor() bool {
	if isTrueColor == nil {
		var isTrueColorCheck = false
		switch os.Getenv("COLORTERM") {
		case "truecolor", "24bit":
			isTrueColorCheck = true
		}

		isTrueColor = &isTrueColorCheck
	}

	return *isTrueColor
}
