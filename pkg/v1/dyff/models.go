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

package dyff

import (
	"io"
)

// Constants to distinguish between the different kinds of differences
const (
	ADDITION     = '+'
	REMOVAL      = '-'
	MODIFICATION = '±'
	ORDERCHANGE  = '⇆'
	// ILLEGAL      = '✕'
	// ATTENTION    = '⚠'
)

// PathElement describes a part of a path, meaning its name. In this case the "Key" string is empty. Named list entries such as "name: one" use both "Key" and "Name" to properly specify the path element.
type PathElement struct {
	Key  string
	Name string
}

// Path describes a position inside a YAML (or JSON) structure by providing a name to each hierarchy level (tree structure).
type Path struct {
	DocumentIdx  int
	PathElements []PathElement
}

// Detail encapsulate the actual details of a change, mainly the kind of difference and the values.
type Detail struct {
	Kind rune
	From interface{}
	To   interface{}
}

// Diff encapsulates everything noteworthy about a difference
type Diff struct {
	Path    Path
	Details []Detail
}

// Report encapsulates the actual end-result of the comparison: The input data and the list of differences.
type Report struct {
	From  InputFile
	To    InputFile
	Diffs []Diff
}

// ReportWriter defines the interface required for types that can write reports
type ReportWriter interface {
	WriteReport(out io.Writer) error
}
