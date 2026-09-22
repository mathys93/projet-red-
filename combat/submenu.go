package combat

import (
	"fmt"
	"strconv"
	"strings"

	"ProjetRED/boss"
	"ProjetRED/character"
)

// menuItem est une entrée affichable dans un sous-menu : un libellé court
// et une description affichée en dessous du cadre quand elle est survolée.
type menuItem struct {
	Label       string
	Description string
}

// chooseOption affiche l'artwork du boss au-dessus du cadre de combat, la
// liste `items` en grille à deux colonnes DANS le cadre (voir
// combat.DrawOptionGrid), et laisse le joueur en choisir un. La
// description de l'entrée survolée s'affiche sous le cadre.
//
// Navigation : Z/S/Q/D (haut/bas/gauche/droite) dans la grille puis Entrée
// pour valider, OU directement le numéro de l'option (1..9) pour la
// choisir ET la valider en une seule saisie - c'est ce raccourci qui règle
// le principal problème de maniabilité du menu de combat (avant, il
// fallait naviguer option par option avant de pouvoir valider). X annule
// et revient au menu FIGHT/ACT/ITEM/MERCY sans consommer le tour.
func chooseOption(ts *terminalSession, b *boss.Boss, title string, items []menuItem) (idx int, ok bool) {
	selected := 0
	for {
		clearScreen()
		bw := boxWidth(boss.ArtWidth + 4)
		bh := boss.ArtHeight + 2
		artHeight := len(strings.Split(b.Art, "\n"))
		printPadding(artHeight + bh + 6)
		fmt.Println()
		fmt.Println(colYellow + "  " + title + colReset)
		DrawArt(bw, b.Art, b.Color)
		DrawOptionGrid(bw, bh, items, selected)
		desc := ""
		if selected < len(items) {
			desc = items[selected].Description
		}
		fmt.Println(colWhite + "* " + desc + colReset)
		fmt.Println()
		fmt.Println(colWhite + "(Z/S/Q/D pour naviguer, un chiffre pour choisir direct, Entrée pour valider, X pour annuler)" + colReset)

		key := ts.readKey()
		switch key {
		case KeyUp:
			if selected-2 >= 0 {
				selected -= 2
			}
		case KeyDown:
			if selected+2 < len(items) {
				selected += 2
			}
		case KeyLeft:
			if selected%2 == 1 {
				selected--
			}
		case KeyRight:
			if selected%2 == 0 && selected+1 < len(items) {
				selected++
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
