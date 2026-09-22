package combat

import (
	"fmt"

	"ProjetRED/world"
)

// SelectZone affiche la carte (3 zones) et laisse le joueur en choisir une
// au clavier. Renvoie nil si le joueur quitte (commande X).
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
			if z.Cleared {
				status = colYellow + "  [zone terminée]" + colReset
			}
			fmt.Printf("%s%s%s\n", marker, z.Name, status)
			fmt.Printf("    %s\n\n", z.Description)
		}
		if note != "" {
			fmt.Println(colWhite + "* " + note + colReset)
			fmt.Println()
		}
		fmt.Println(colWhite + "(Z/S pour choisir, Entrée seule pour valider, X pour quitter)" + colReset)

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
			return z
		}
	}
}

// ConfirmBonusBoss affiche un écran simple demandant au joueur s'il veut
// affronter le boss bonus (débloqué une fois les 3 zones terminées).
func ConfirmBonusBoss(bossName string) bool {
	ts := newTerminalSession()
	defer ts.restore()

	for {
		clearScreen()
		printPadding(5)
		fmt.Println()
		fmt.Println(colYellow + "Les trois zones sont terminées." + colReset)
		fmt.Printf("Un dernier adversaire t'attend : %s%s%s\n\n", colRed, bossName, colReset)
		fmt.Println(colWhite + "L'affronter ? (Entrée seule = oui, X = non)" + colReset)

		key := ts.readKey()
		switch key {
		case KeyEnter:
			return true
		case KeyQuit:
			return false
		}
	}
}
