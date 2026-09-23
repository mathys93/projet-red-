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
	// "close all" est un filet de sécurité : si le "close" ci-dessus a
	// échoué (l'alias restait alors verrouillé sur l'ancien device MCI),
	// le Play() suivant retombait sur l'ancienne piste toujours en train
	// de jouer EN PLUS de la nouvelle - d'où les deux musiques superposées
	// en passant d'un combat à l'autre. "close all" force la fermeture de
	// tout device MCI encore ouvert, peu importe son alias.
	send("close all")
}
