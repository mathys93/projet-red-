//go:build !windows

package combat

func enterRawMode() {}

func exitRawMode() {}

func (ts *terminalSession) readKey() Key {
	return ts.readLineKey()
}

func WaitEnter() {
	stdinScanner.Scan()
}
