//go:build windows

package audio

import (
	"fmt"
	"runtime"
	"syscall"
	"unsafe"
)

var (
	winmm             = syscall.NewLazyDLL("winmm.dll")
	procMCISendString = winmm.NewProc("mciSendStringW")
)

const alias = "projetredbgm"

type mciCall struct {
	command string
	result  chan bool
}

var mciQueue = make(chan mciCall)

func init() {
	go func() {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
		for call := range mciQueue {
			call.result <- rawSend(call.command)
		}
	}()
}

func rawSend(command string) bool {
	ptr, err := syscall.UTF16PtrFromString(command)
	if err != nil {
		return false
	}
	ret, _, _ := procMCISendString.Call(uintptr(unsafe.Pointer(ptr)), 0, 0, 0)
	return ret == 0
}

func send(command string) bool {
	result := make(chan bool, 1)
	mciQueue <- mciCall{command: command, result: result}
	return <-result
}

var current string

func clampedVolume() int {
	return max(0, min(1000, Volume))
}

func closeDevice() {
	send("stop " + alias)
	send("close " + alias)
}

func start(path string) {
	closeDevice()
	path = resolve(path)
	if path == "" {
		return
	}
	if !send(fmt.Sprintf(`open "%s" type mpegvideo alias %s`, path, alias)) {
		if !send(fmt.Sprintf(`open "%s" alias %s`, path, alias)) {
			return
		}
	}
	send(fmt.Sprintf("setaudio %s volume to %d", alias, clampedVolume()))
	if !send(fmt.Sprintf("play %s repeat", alias)) {
		send(fmt.Sprintf("play %s", alias))
	}
}

func Play(path string) {
	current = path
	if !Enabled {
		closeDevice()
		return
	}
	start(path)
}

func Stop() {
	current = ""
	closeDevice()
}

func SetVolume(v int) {
	Volume = max(0, min(1000, v))
	send(fmt.Sprintf("setaudio %s volume to %d", alias, clampedVolume()))
}

func SetEnabled(on bool) {
	Enabled = on
	if !on {
		closeDevice()
		return
	}
	if current != "" {
		start(current)
	}
}
