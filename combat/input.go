package combat

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Key string

const (
	KeyUp    Key = "UP"
	KeyDown  Key = "DOWN"
	KeyLeft  Key = "LEFT"
	KeyRight Key = "RIGHT"
	KeyEnter Key = "ENTER"
	KeyQuit  Key = "QUIT"
	KeyOther Key = "OTHER"
)

var stdinScanner = bufio.NewScanner(os.Stdin)

type terminalSession struct {
	raw string
}

func newTerminalSession() *terminalSession {
	enterRawMode()
	return &terminalSession{}
}

func (ts *terminalSession) restore() {
	exitRawMode()
	fmt.Print("\033[?25h")
}

func keyForRune(r rune) (Key, bool) {
	switch r {
	case 'z', 'Z', 'w', 'W':
		return KeyUp, true
	case 's', 'S':
		return KeyDown, true
	case 'q', 'Q':
		return KeyLeft, true
	case 'd', 'D':
		return KeyRight, true
	case 'x', 'X':
		return KeyQuit, true
	case '\r', '\n':
		return KeyEnter, true
	}
	return KeyOther, false
}

func (ts *terminalSession) readLineKey() Key {
	if !stdinScanner.Scan() {
		return KeyQuit
	}
	line := strings.ToLower(strings.TrimSpace(stdinScanner.Text()))
	ts.raw = line
	if line == "" {
		return KeyEnter
	}
	switch line {
	case "up", "haut":
		return KeyUp
	case "down", "bas":
		return KeyDown
	case "left", "gauche":
		return KeyLeft
	case "right", "droite":
		return KeyRight
	case "quit", "exit", "quitter":
		return KeyQuit
	}
	if r := []rune(line); len(r) == 1 {
		if k, ok := keyForRune(r[0]); ok {
			return k
		}
	}
	return KeyOther
}
