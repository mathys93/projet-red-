//go:build windows

package combat

import (
	"os"
	"strings"
	"syscall"
	"unsafe"
)

var (
	procReadConsoleInput  = kernel32.NewProc("ReadConsoleInputW")
	procGetConsoleMode    = kernel32.NewProc("GetConsoleMode")
	procSetConsoleMode    = kernel32.NewProc("SetConsoleMode")
	procGetNumberOfEvents = kernel32.NewProc("GetNumberOfConsoleInputEvents")
	procFlushConsoleInput = kernel32.NewProc("FlushConsoleInputBuffer")
)

type keyEventRecord struct {
	KeyDown         int32
	RepeatCount     uint16
	VirtualKeyCode  uint16
	VirtualScanCode uint16
	UnicodeChar     uint16
	ControlKeyState uint32
}

type inputRecord struct {
	EventType uint16
	_         uint16
	KeyEvent  keyEventRecord
}

const (
	eventKey             = 0x0001
	enableProcessedInput = 0x0001
	enableLineInput      = 0x0002
	enableEchoInput      = 0x0004

	vkReturn = 0x0D
	vkEscape = 0x1B
	vkSpace  = 0x20
	vkLeft   = 0x25
	vkUp     = 0x26
	vkRight  = 0x27
	vkDown   = 0x28
)

var (
	savedInputMode uint32
	rawModeActive  bool
)

func stdinHandle() syscall.Handle {
	return syscall.Handle(os.Stdin.Fd())
}

func enterRawMode() {
	if rawModeActive {
		return
	}
	var mode uint32
	ret, _, _ := procGetConsoleMode.Call(uintptr(stdinHandle()), uintptr(unsafe.Pointer(&mode)))
	if ret == 0 {
		return
	}
	savedInputMode = mode
	newMode := mode&^(enableLineInput|enableEchoInput) | enableProcessedInput
	if ret, _, _ := procSetConsoleMode.Call(uintptr(stdinHandle()), uintptr(newMode)); ret == 0 {
		return
	}
	rawModeActive = true
}

func exitRawMode() {
	if !rawModeActive {
		return
	}
	procSetConsoleMode.Call(uintptr(stdinHandle()), uintptr(savedInputMode))
	rawModeActive = false
}

func consoleAvailable() bool {
	var mode uint32
	ret, _, _ := procGetConsoleMode.Call(uintptr(stdinHandle()), uintptr(unsafe.Pointer(&mode)))
	return ret != 0
}

func keyForVirtualKey(vk uint16) (Key, bool) {
	switch vk {
	case vkReturn, vkSpace:
		return KeyEnter, true
	case vkEscape:
		return KeyQuit, true
	case vkLeft:
		return KeyLeft, true
	case vkRight:
		return KeyRight, true
	case vkUp:
		return KeyUp, true
	case vkDown:
		return KeyDown, true
	}
	return KeyOther, false
}

func readConsoleKey() (Key, rune, bool) {
	var rec inputRecord
	var read uint32
	for {
		ret, _, _ := procReadConsoleInput.Call(
			uintptr(stdinHandle()),
			uintptr(unsafe.Pointer(&rec)),
			1,
			uintptr(unsafe.Pointer(&read)),
		)
		if ret == 0 {
			return KeyQuit, 0, true
		}
		if read == 0 || rec.EventType != eventKey || rec.KeyEvent.KeyDown == 0 {
			continue
		}

		if k, ok := keyForVirtualKey(rec.KeyEvent.VirtualKeyCode); ok {
			return k, 0, true
		}

		if ch := rune(rec.KeyEvent.UnicodeChar); ch != 0 {
			if k, ok := keyForRune(ch); ok {
				return k, ch, true
			}
			return KeyOther, ch, true
		}
	}
}

func canPollInput() bool {
	return consoleAvailable()
}

func pollKey() (Key, bool) {
	var pending uint32
	ret, _, _ := procGetNumberOfEvents.Call(uintptr(stdinHandle()), uintptr(unsafe.Pointer(&pending)))
	if ret == 0 || pending == 0 {
		return KeyOther, false
	}

	var rec inputRecord
	var read uint32
	key := KeyOther
	found := false
	for i := uint32(0); i < pending; i++ {
		r, _, _ := procReadConsoleInput.Call(
			uintptr(stdinHandle()),
			uintptr(unsafe.Pointer(&rec)),
			1,
			uintptr(unsafe.Pointer(&read)),
		)
		if r == 0 || read == 0 {
			break
		}
		if rec.EventType != eventKey || rec.KeyEvent.KeyDown == 0 {
			continue
		}
		if k, ok := keyForVirtualKey(rec.KeyEvent.VirtualKeyCode); ok {
			key, found = k, true
			continue
		}
		if ch := rune(rec.KeyEvent.UnicodeChar); ch != 0 {
			if k, ok := keyForRune(ch); ok {
				key, found = k, true
			}
		}
	}
	return key, found
}

func flushInput() {
	procFlushConsoleInput.Call(uintptr(stdinHandle()))
}

func (ts *terminalSession) readKey() Key {
	if !consoleAvailable() {
		return ts.readLineKey()
	}
	k, ch, _ := readConsoleKey()
	ts.raw = ""
	if ch != 0 {
		ts.raw = strings.ToLower(string(ch))
	}
	return k
}

func WaitEnter() {
	if !consoleAvailable() {
		stdinScanner.Scan()
		return
	}
	enterRawMode()
	defer exitRawMode()
	readConsoleKey()
}
