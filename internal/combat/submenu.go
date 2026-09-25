package combat

import (
	"fmt"
	"strconv"
	"strings"

	"ProjetRED/internal/boss"
	"ProjetRED/internal/character"
	"ProjetRED/internal/terminal"
)

func chooseOption(ts *terminal.Session, b *boss.Boss, title string, items []terminal.MenuItem) (idx int, ok bool) {
	selected := 0
	render := func() {
		terminal.ClearScreen()
		bw, bh, aw, ah := terminal.Layout()
		if needed := (len(items)+1)/2 + 2; needed > bh {
			bh = needed
		}
		art := fitArt(b, aw, ah)
		artHeight := len(strings.Split(art, "\n"))
		terminal.PrintPadding(artHeight + bh + 6)
		fmt.Fprintln(terminal.Out)
		fmt.Fprintln(terminal.Out, terminal.ColYellow+"  "+title+terminal.ColReset)
		terminal.DrawArt(bw, art, b.Color)
		terminal.DrawOptionGrid(bw, bh, items, selected)
		desc := ""
		if selected < len(items) {
			desc = items[selected].Description
		}
		fmt.Fprintln(terminal.Out, terminal.ColWhite+"* "+desc+terminal.ColReset)
		fmt.Fprintln(terminal.Out)
		fmt.Fprintln(terminal.Out, terminal.ColWhite+"(← ↑ ↓ → pour naviguer, un chiffre pour choisir direct, Entrée pour valider, X pour annuler)"+terminal.ColReset)
	}
	for {
		terminal.CurrentScreen = render
		render()

		key := ts.ReadKey()
		switch key {
		case terminal.KeyUp:
			if selected-2 >= 0 {
				selected -= 2
			}
		case terminal.KeyDown:
			if selected+2 < len(items) {
				selected += 2
			}
		case terminal.KeyLeft:
			if selected%2 == 1 {
				selected--
			}
		case terminal.KeyRight:
			if selected%2 == 0 && selected+1 < len(items) {
				selected++
			}
		case terminal.KeyBack, terminal.KeyQuit:
			return 0, false
		case terminal.KeyEnter:
			return selected, true
		case terminal.KeyOther:
			if n, err := strconv.Atoi(ts.Raw); err == nil && n >= 1 && n <= len(items) {
				return n - 1, true
			}
		}
	}
}

func fightMenuItems(player *character.Character) []terminal.MenuItem {
	items := make([]terminal.MenuItem, len(player.Moves))
	for i, m := range player.Moves {
		desc := m.Description
		if m.MPCost > 0 {
			desc = fmt.Sprintf("%s (coût : %d MP)", desc, m.MPCost)
		}
		items[i] = terminal.MenuItem{Label: m.Name, Description: desc}
	}
	return items
}

func actMenuItems(options []boss.ActOption) []terminal.MenuItem {
	items := make([]terminal.MenuItem, len(options))
	for i, o := range options {
		items[i] = terminal.MenuItem{Label: o.Label, Description: o.Description}
	}
	return items
}
