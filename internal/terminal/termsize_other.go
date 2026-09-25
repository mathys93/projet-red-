//go:build !windows

package terminal

func queryTerminalSize() (cols, rows int) {
	return 0, 0
}

func maximizeConsoleWindow() {}
