// +build !windows

package core

import (
	"syscall"
	"unsafe"
)

// GetTerminalWidth return the current terminal width
func GetTerminalWidth() int {
	const defaultWidth = 80

	type winsize struct {
		Row    uint16
		Col    uint16
		Xpixel uint16
		Ypixel uint16
	}

	ws := &winsize{}
	retCode, _, _ := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)))

	if int(retCode) == -1 {
		// TODO Add debug output that default terminal sizings were returned due
		return defaultWidth
	}

	return int(ws.Col)
}
