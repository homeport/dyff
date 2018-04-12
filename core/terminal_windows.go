// +build windows

package core

// GetTerminalWidth return the current terminal width
func GetTerminalWidth() int {
	// TODO Implement Windows terminal width if this is possible at all.
	return 80
}
