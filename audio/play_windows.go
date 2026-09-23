//go:build windows

package audio

import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	winmm             = syscall.NewLazyDLL("winmm.dll")
	procMCISendString = winmm.NewProc("mciSendStringW")
)

const alias = "projetredbgm"

func send(command string) bool {
	ptr, err := syscall.UTF16PtrFromString(command)
	if err != nil {
		return false
	}
	ret, _, _ := procMCISendString.Call(uintptr(unsafe.Pointer(ptr)), 0, 0, 0)
	return ret == 0
}

func Play(path string) {
	Stop()
	path = resolve(path)
	if path == "" {
		return
	}
	if !send(fmt.Sprintf(`open "%s" type mpegvideo alias %s`, path, alias)) {
		if !send(fmt.Sprintf(`open "%s" alias %s`, path, alias)) {
			return
		}
	}
	v := Volume
	if v < 0 {
		v = 0
	}
	if v > 1000 {
		v = 1000
	}
	send(fmt.Sprintf("setaudio %s volume to %d", alias, v))
	if !send(fmt.Sprintf("play %s repeat", alias)) {
		send(fmt.Sprintf("play %s", alias))
	}
}

func Stop() {
	send("stop " + alias)
	send("close " + alias)
}
