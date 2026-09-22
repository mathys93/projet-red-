//go:build !windows

package combat

func queryTerminalSize() (cols, rows int) {
	return 0, 0
}

func maximizeConsoleWindow() {}
