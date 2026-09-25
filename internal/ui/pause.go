package ui

import (
	"ProjetRED/internal/terminal"
	"fmt"
	"strings"
)

const (
	pauseWidth  = 40
	pauseHeight = 11
	pauseButton = 22

	bgBlack = "\033[48;2;0;0;0m"
	bgGrey  = "\033[48;2;170;170;170m"
	bgWhite = "\033[48;2;255;255;255m"
)

var pauseLabels = []string{"PARAMÈTRES", "CRÉDITS", "QUITTER"}

func pauseOrigin() (top, left int) {
	cols, rows := terminal.TerminalSize()
	return max(1, (rows-pauseHeight)/2+1), max(1, (cols-pauseWidth)/2+1)
}

func drawPause(selected int) {
	top, left := pauseOrigin()
	frame := bgBlack + egaBlue
	inner := pauseWidth - 2

	var sb strings.Builder
	title := " PAUSE "
	side := (inner - len([]rune(title))) / 2
	sb.WriteString(terminal.CursorAt(top, left))
	sb.WriteString(frame)
	sb.WriteString("╔")
	sb.WriteString(strings.Repeat("═", side))
	sb.WriteString(title)
	sb.WriteString(strings.Repeat("═", inner-side-len([]rune(title))))
	sb.WriteString("╗")
	for r := 1; r < pauseHeight-1; r++ {
		sb.WriteString(terminal.CursorAt(top+r, left))
		sb.WriteString(frame)
		sb.WriteString("║")
		sb.WriteString(terminal.CursorAt(top+r, left+pauseWidth-1))
		sb.WriteString(frame)
		sb.WriteString("║")
	}
	sb.WriteString(terminal.CursorAt(top+pauseHeight-1, left))
	sb.WriteString(frame)
	sb.WriteString("╚")
	sb.WriteString(strings.Repeat("═", inner))
	sb.WriteString("╝")

	btnLeft := left + (pauseWidth-pauseButton)/2
	for i, label := range pauseLabels {
		bg, heart := bgGrey, "  "
		if i == selected {
			bg, heart = bgWhite, terminal.ColRed+"❤ "
		}
		pad := pauseButton - 2 - len([]rune(label))
		l, r := pad/2, pad-pad/2
		sb.WriteString(terminal.CursorAt(top+2+2*i, btnLeft))
		sb.WriteString(bg)
		sb.WriteString(heart)
		sb.WriteString(egaBlue)
		sb.WriteString(strings.Repeat(" ", l))
		sb.WriteString(label)
		sb.WriteString(strings.Repeat(" ", r))
	}

	hint := "Entrée : valider · Échap : reprendre"
	hintLeft := left + (pauseWidth-len([]rune(hint)))/2
	sb.WriteString(terminal.CursorAt(top+pauseHeight-2, hintLeft))
	sb.WriteString(frame)
	sb.WriteString(hint)

	sb.WriteString(terminal.ColReset)
	fmt.Fprint(terminal.Out, sb.String())
}

func init() {
	terminal.OnPause = openPause
}

func openPause() {
	terminal.FlushInput()

	ts := &terminal.Session{}
	selected := 0
	for {
		drawPause(selected)
		switch ts.ReadKey() {
		case terminal.KeyUp, terminal.KeyLeft:
			selected = (selected + len(pauseLabels) - 1) % len(pauseLabels)
		case terminal.KeyDown, terminal.KeyRight:
			selected = (selected + 1) % len(pauseLabels)
		case terminal.KeyPause, terminal.KeyBack:
			terminal.Redraw()
			return
		case terminal.KeyQuit:
			return
		case terminal.KeyEnter:
			switch selected {
			case 0:
				RunSettings()
				terminal.Redraw()
			case 1:
				RunCredits()
				terminal.Redraw()
			case 2:
				terminal.RequestQuit()
				return
			}
		}
	}
}
