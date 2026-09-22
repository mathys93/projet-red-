package combat

import (
	"fmt"
	"strconv"

	"ProjetRED/boss"
	"ProjetRED/character"
)

// menuItem est une entrée affichable dans un sous-menu : un libellé court
// et une description affichée en dessous.
type menuItem struct {
	Label       string
	Description string
}

// chooseOption affiche `title` suivi de la liste `items` et laisse le
// joueur en choisir un.
//
// Navigation : Z/S (haut/bas) puis Entrée pour valider, OU directement le
// numéro de l'option (1..9) pour la choisir ET la valider en une seule
// saisie - c'est ce raccourci qui règle le principal problème de
// maniabilité du menu de combat (avant, il fallait naviguer option par
// option avant de pouvoir valider). X annule et revient au menu FIGHT/ACT/
// ITEM/MERCY sans consommer le tour.
func chooseOption(ts *terminalSession, title string, items []menuItem) (idx int, ok bool) {
	selected := 0
	for {
		clearScreen()
		fmt.Println()
		fmt.Println(colYellow + "  " + title + colReset)
		fmt.Println()
		for i, it := range items {
			marker := "   "
			if i == selected {
				marker = colRed + "❤ " + colReset
			}
			fmt.Printf("%s%s%d) %s%s\n", marker, colYellow, i+1, it.Label, colReset)
			if it.Description != "" {
				fmt.Printf("      %s%s%s\n", colWhite, it.Description, colReset)
			}
		}
		fmt.Println()
		fmt.Println(colWhite + "(Z/S pour naviguer, un chiffre pour choisir direct, Entrée pour valider, X pour annuler)" + colReset)

		key := ts.readKey()
		switch key {
		case KeyUp, KeyLeft:
			selected--
			if selected < 0 {
				selected = len(items) - 1
			}
		case KeyDown, KeyRight:
			selected++
			if selected >= len(items) {
				selected = 0
			}
		case KeyQuit:
			return 0, false
		case KeyEnter:
			return selected, true
		case KeyOther:
			if n, err := strconv.Atoi(ts.raw); err == nil && n >= 1 && n <= len(items) {
				return n - 1, true
			}
		}
	}
}

// fightMenuItems convertit les Move du joueur en entrées de sous-menu,
// avec le coût en MP affiché dans la description quand il y en a un.
func fightMenuItems(player *character.Character) []menuItem {
	items := make([]menuItem, len(player.Moves))
	for i, m := range player.Moves {
		desc := m.Description
		if m.MPCost > 0 {
			desc = fmt.Sprintf("%s (coût : %d MP)", desc, m.MPCost)
		}
		items[i] = menuItem{Label: m.Name, Description: desc}
	}
	return items
}

// actMenuItems convertit les ActOption d'un boss en entrées de sous-menu.
func actMenuItems(options []boss.ActOption) []menuItem {
	items := make([]menuItem, len(options))
	for i, o := range options {
		items[i] = menuItem{Label: o.Label, Description: o.Description}
	}
	return items
}
