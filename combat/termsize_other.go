//go:build !windows

package combat

// queryTerminalSize n'a pas d'implémentation spécifique hors Windows : le
// jeu se replie sur la taille de terminal par défaut (voir terminalSize
// dans battle.go).
func queryTerminalSize() (cols, rows int) {
	return 0, 0
}

// maximizeConsoleWindow n'a pas d'équivalent générique hors Windows.
func maximizeConsoleWindow() {}
