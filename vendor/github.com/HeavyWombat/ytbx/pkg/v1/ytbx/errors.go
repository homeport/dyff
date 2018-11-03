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

package ytbx

import (
	"fmt"
	"strings"
)

type KeyNotFoundInMapError struct {
	MissingKey    string
	AvailableKeys []string
}

func (e *KeyNotFoundInMapError) Error() string {
	return fmt.Sprintf("no key '%s' found in map, available keys: %s",
		e.MissingKey,
		strings.Join(e.AvailableKeys, ", "))
}

type NoNamedEntryListError struct {
}

func (e *NoNamedEntryListError) Error() string {
	return "not a named-entry list, one or more entries are not of type map"
}

type InvalidPathString struct {
	Style       PathStyle
	PathString  string
	Explanation string
}

func (e *InvalidPathString) Error() string {
	return fmt.Sprintf("invalid %v style path %s, %s",
		e.Style,
		e.PathString,
		e.Explanation)
}
