//go:build !windows

package terminal

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

func ReadTextKey() (TextKey, rune) { return TextCancel, 0 }
