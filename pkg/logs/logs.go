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
Package logs is a convenience package to cover simply logging functionality
*/
package logs

import (
	"log"
	"os"
)

// LogLevel covers different types of logging severities
type LogLevel int

// List of possile logging levels
const (
	NONE LogLevel = iota
	WARN
	DEBUG
)

// LoggingLevel stores the currently configured logging level
var LoggingLevel = WARN

// ErrorLogger is the error logger definition
var ErrorLogger = log.New(os.Stderr, "Error: ", log.Ldate|log.Ltime|log.Lshortfile)

// WarningLogger is the warning logger definition
var WarningLogger = log.New(os.Stdout, "Warning: ", log.Ldate|log.Ltime|log.Lshortfile)

// DebugLogger is the debugging logger definition
var DebugLogger = log.New(os.Stdout, "Debug: ", log.Ldate|log.Ltime|log.Lshortfile)

// Debug prints a debug statement if logging level matches
func Debug(format string, a ...interface{}) {
	switch LoggingLevel {
	case DEBUG:
		DebugLogger.Printf(format, a...)
	}
}

// Warn prints a warning statement if logging level matches
func Warn(format string, a ...interface{}) {
	switch LoggingLevel {
	case WARN, DEBUG:
		WarningLogger.Printf(format, a...)
	}
}

// Error prints a non-fatal error statement
func Error(format string, a ...interface{}) {
	ErrorLogger.Printf(format, a...)
}
