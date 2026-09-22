package combat

import (
	"fmt"
	"strconv"
	"strings"

	"ProjetRED/boss"
	"ProjetRED/character"
)

type menuItem struct {
	Label       string
	Description string
}

func chooseOption(ts *terminalSession, b *boss.Boss, title string, items []menuItem) (idx int, ok bool) {
	selected := 0
	for {
		clearScreen()
		bw, bh, aw, ah := layout()
		if needed := (len(items)+1)/2 + 2; needed > bh {
			bh = needed
		}
		art := fitArt(b, aw, ah)
		artHeight := len(strings.Split(art, "\n"))
		printPadding(artHeight + bh + 6)
		fmt.Println()
		fmt.Println(colYellow + "  " + title + colReset)
		DrawArt(bw, art, b.Color)
		DrawOptionGrid(bw, bh, items, selected)
		desc := ""
		if selected < len(items) {
			desc = items[selected].Description
		}
		fmt.Println(colWhite + "* " + desc + colReset)
		fmt.Println()
		fmt.Println(colWhite + "(← ↑ ↓ → pour naviguer, un chiffre pour choisir direct, Entrée pour valider, X pour annuler)" + colReset)

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

func actMenuItems(options []boss.ActOption) []menuItem {
	items := make([]menuItem, len(options))
	for i, o := range options {
		items[i] = menuItem{Label: o.Label, Description: o.Description}
	}
	return items
}
