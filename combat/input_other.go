//go:build !windows

package combat

func enterRawMode() {}

func exitRawMode() {}

func canPollInput() bool { return false }

func pollKey() (Key, bool) { return KeyOther, false }

func flushInput() {}

func (ts *terminalSession) readKey() Key {
	return ts.readLineKey()
}

func WaitEnter() {
	stdinScanner.Scan()
}
