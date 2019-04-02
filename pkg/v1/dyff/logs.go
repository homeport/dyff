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
	"io/ioutil"
	"log"
	"os"
)

// LogLevel covers different types of logging severities
type LogLevel int

// List of possile logging levels
const (
	NONE LogLevel = iota
	ERROR
	WARN
	DEBUG
)

// LoggingLevel stores the currently configured logging level
var LoggingLevel = ERROR

// ErrorLogger is the error logger definition
var ErrorLogger = log.New(os.Stderr, "Error: ", log.Lshortfile)

// WarningLogger is the warning logger definition
var WarningLogger = log.New(ioutil.Discard, "Warning: ", log.Lshortfile)

// DebugLogger is the debugging logger definition
var DebugLogger = log.New(ioutil.Discard, "Debug: ", log.Lshortfile)

// SetLoggingLevel will initialise the logging set-up according to the provided input
func SetLoggingLevel(loggingLevel LogLevel) {
	switch loggingLevel {
	case NONE:
		ErrorLogger.SetOutput(ioutil.Discard)
		WarningLogger.SetOutput(ioutil.Discard)
		DebugLogger.SetOutput(ioutil.Discard)

	case ERROR:
		ErrorLogger.SetOutput(os.Stderr)
		WarningLogger.SetOutput(ioutil.Discard)
		DebugLogger.SetOutput(ioutil.Discard)

	case WARN:
		ErrorLogger.SetOutput(os.Stderr)
		WarningLogger.SetOutput(os.Stdout)
		DebugLogger.SetOutput(ioutil.Discard)

	case DEBUG:
		ErrorLogger.SetOutput(os.Stderr)
		WarningLogger.SetOutput(os.Stdout)
		DebugLogger.SetOutput(os.Stdout)
	}
}
