package terminal

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
	KeyBack  Key = "BACK"
	KeyPause Key = "PAUSE"
	KeyQuit  Key = "QUIT"
	KeyOther Key = "OTHER"
)

type KeyEvent struct {
	Key  Key
	Down bool
}

type TextKey int

const (
	TextChar TextKey = iota
	TextBackspace
	TextEnter
	TextCancel
)

var stdinScanner = bufio.NewScanner(os.Stdin)

var (
	gameActive    bool
	inPause       bool
	quitRequested bool
	CurrentScreen func()
)

func EnterGame() {
	gameActive = true
	quitRequested = false
	CurrentScreen = nil
}

func LeaveGame() {
	gameActive = false
	CurrentScreen = nil
}

func QuitRequested() bool {
	return quitRequested
}

var OnPause func()

func Pause() {
	if OnPause == nil || inPause {
		return
	}
	inPause = true
	defer func() { inPause = false }()
	OnPause()
}

func RequestQuit() {
	quitRequested = true
}

func Redraw() {
	if CurrentScreen != nil {
		CurrentScreen()
	}
}

type Session struct {
	Raw   string
	owner bool
}

func NewSession() *Session {
	ts := &Session{owner: !rawModeOn()}
	enterRawMode()
	return ts
}

func (ts *Session) Restore() {
	if !ts.owner {
		return
	}
	exitRawMode()
	fmt.Fprint(Out, "\033[?25h")
}

func (ts *Session) ReadKey() Key {
	for {
		if quitRequested {
			return KeyQuit
		}
		k := ts.readKeyRaw()
		if k == KeyQuit {
			quitRequested = true
			return KeyQuit
		}
		if k == KeyPause && gameActive && !inPause {
			Pause()
			continue
		}
		return k
	}
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
		return KeyBack, true
	case '\r', '\n':
		return KeyEnter, true
	}
	return KeyOther, false
}

func (ts *Session) readLineKey() Key {
	if !stdinScanner.Scan() {
		return KeyQuit
	}
	line := strings.ToLower(strings.TrimSpace(stdinScanner.Text()))
	ts.Raw = line
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
	case "esc", "echap", "échap", "pause":
		return KeyPause
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

func ReadLine() (string, bool) {
	if !stdinScanner.Scan() {
		return "", false
	}
	return stdinScanner.Text(), true
}
