package combat

import (
	"fmt"

	"ProjetRED/world"
)

func SelectZone(w *world.World) *world.Zone {
	ts := newTerminalSession()
	defer ts.restore()

	selected := 0
	note := ""

	for {
		clearScreen()
		contentHeight := 4 + len(w.Zones)*3
		if note != "" {
			contentHeight += 2
		}
		printPadding(contentHeight)
		fmt.Println()
		fmt.Println(colYellow + "======================  CARTE  ======================" + colReset)
		fmt.Println()
		for i, z := range w.Zones {
			marker := "  "
			if i == selected {
				marker = colRed + "❤ " + colReset
			}
			status := ""
			switch {
			case z.Cleared:
				status = colYellow + "  [zone terminée]" + colReset
			case w.Locked(z):
				status = colWhite + "  [verrouillée]" + colReset
			}
			fmt.Printf("%s%s%s\n", marker, z.Name, status)
			fmt.Printf("    %s\n\n", z.Description)
		}
		if note != "" {
			fmt.Println(colWhite + "* " + note + colReset)
			fmt.Println()
		}
		fmt.Println(colWhite + "(↑ ↓ pour choisir, Entrée pour valider, X pour quitter)" + colReset)

		key := ts.readKey()
		last := len(w.Zones) - 1
		switch key {
		case KeyUp, KeyLeft:
			selected--
			if selected < 0 {
				selected = last
			}
		case KeyDown, KeyRight:
			selected++
			if selected > last {
				selected = 0
			}
		case KeyQuit:
			return nil
		case KeyEnter:
			z := w.Zones[selected]
			if z.Cleared {
				note = fmt.Sprintf("%s est déjà terminée.", z.Name)
				continue
			}
			if w.Locked(z) {
				note = fmt.Sprintf("%s est verrouillée : termine d'abord la zone précédente.", z.Name)
				continue
			}
			return z
		}
	}
}

func ConfirmBonusBoss(bossName string) bool {
	ts := newTerminalSession()
	defer ts.restore()

	for {
		clearScreen()
		printPadding(5)
		fmt.Println()
		fmt.Println(colYellow + "Les trois zones sont terminées." + colReset)
		fmt.Printf("Un dernier adversaire t'attend : %s%s%s\n\n", colRed, bossName, colReset)
		fmt.Println(colWhite + "L'affronter ? (Entrée = oui, X = non)" + colReset)

		key := ts.readKey()
		switch key {
		case KeyEnter:
			return true
		case KeyQuit:
			return false
		}
	}
}
