//go:build windows

package combat

import (
	"os"
	"syscall"
	"unsafe"
)

var (
	kernel32                       = syscall.NewLazyDLL("kernel32.dll")
	user32                         = syscall.NewLazyDLL("user32.dll")
	procGetConsoleScreenBufferInfo = kernel32.NewProc("GetConsoleScreenBufferInfo")
	procGetConsoleWindow           = kernel32.NewProc("GetConsoleWindow")
	procShowWindow                 = user32.NewProc("ShowWindow")
)

type coord struct {
	X, Y int16
}

type smallRect struct {
	Left, Top, Right, Bottom int16
}

type consoleScreenBufferInfo struct {
	Size              coord
	CursorPosition    coord
	Attributes        uint16
	Window            smallRect
	MaximumWindowSize coord
}

// queryTerminalSize lit la taille de la fenêtre console réellement visible
// (pas celle du buffer de défilement, qui peut être bien plus grande) via
// l'API Windows, pour que l'affichage puisse remplir tout l'écran plutôt
// que de rester coincé dans un coin façon boîte figée.
func queryTerminalSize() (cols, rows int) {
	var info consoleScreenBufferInfo
	handle := syscall.Handle(os.Stdout.Fd())
	ret, _, _ := procGetConsoleScreenBufferInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&info)))
	if ret == 0 {
		return 0, 0
	}
	cols = int(info.Window.Right-info.Window.Left) + 1
	rows = int(info.Window.Bottom-info.Window.Top) + 1
	return cols, rows
}

// maximizeConsoleWindow agrandit la fenêtre de la console au maximum
// (façon jeu en plein écran) si le jeu tourne dans une vraie fenêtre
// console (conhost, Windows Terminal...). Ne fait rien si le programme
// tourne dans un terminal sans fenêtre propre (ex: certains terminaux
// intégrés à un éditeur) : GetConsoleWindow renvoie alors un handle nul.
func maximizeConsoleWindow() {
	hwnd, _, _ := procGetConsoleWindow.Call()
	if hwnd == 0 {
		return
	}
	const swMaximize = 3
	procShowWindow.Call(hwnd, swMaximize)
}
