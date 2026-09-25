//go:build !windows

package ui

func queryTerminalSize() (cols, rows int) {
	return 0, 0
}

func maximizeConsoleWindow() {}
