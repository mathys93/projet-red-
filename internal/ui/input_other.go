//go:build !windows

package ui

func enterRawMode() {}

func exitRawMode() {}

func rawModeOn() bool { return false }

func CanPollInput() bool { return false }

func PollKeyEvents() []KeyEvent { return nil }

func FlushInput() {}

func (ts *Session) readKeyRaw() Key {
	return ts.readLineKey()
}

func WaitEnter() {
	stdinScanner.Scan()
}

func readTextKey() (textKey, rune) { return textCancel, 0 }
