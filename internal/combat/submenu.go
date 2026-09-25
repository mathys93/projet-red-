package combat

import (
	"fmt"
	"strconv"
	"strings"

	"ProjetRED/internal/boss"
	"ProjetRED/internal/character"
	"ProjetRED/internal/ui"
)

func chooseOption(ts *ui.Session, b *boss.Boss, title string, items []ui.MenuItem) (idx int, ok bool) {
	selected := 0
	render := func() {
		ui.ClearScreen()
		bw, bh, aw, ah := ui.Layout()
		if needed := (len(items)+1)/2 + 2; needed > bh {
			bh = needed
		}
		art := fitArt(b, aw, ah)
		artHeight := len(strings.Split(art, "\n"))
		ui.PrintPadding(artHeight + bh + 6)
		fmt.Fprintln(ui.Out)
		fmt.Fprintln(ui.Out, ui.ColYellow+"  "+title+ui.ColReset)
		ui.DrawArt(bw, art, b.Color)
		ui.DrawOptionGrid(bw, bh, items, selected)
		desc := ""
		if selected < len(items) {
			desc = items[selected].Description
		}
		fmt.Fprintln(ui.Out, ui.ColWhite+"* "+desc+ui.ColReset)
		fmt.Fprintln(ui.Out)
		fmt.Fprintln(ui.Out, ui.ColWhite+"(← ↑ ↓ → pour naviguer, un chiffre pour choisir direct, Entrée pour valider, X pour annuler)"+ui.ColReset)
	}
	for {
		ui.CurrentScreen = render
		render()

		key := ts.ReadKey()
		switch key {
		case ui.KeyUp:
			if selected-2 >= 0 {
				selected -= 2
			}
		case ui.KeyDown:
			if selected+2 < len(items) {
				selected += 2
			}
		case ui.KeyLeft:
			if selected%2 == 1 {
				selected--
			}
		case ui.KeyRight:
			if selected%2 == 0 && selected+1 < len(items) {
				selected++
			}
		case ui.KeyBack, ui.KeyQuit:
			return 0, false
		case ui.KeyEnter:
			return selected, true
		case ui.KeyOther:
			if n, err := strconv.Atoi(ts.Raw); err == nil && n >= 1 && n <= len(items) {
				return n - 1, true
			}
		}
	}
}

func fightMenuItems(player *character.Character) []ui.MenuItem {
	items := make([]ui.MenuItem, len(player.Moves))
	for i, m := range player.Moves {
		desc := m.Description
		if m.MPCost > 0 {
			desc = fmt.Sprintf("%s (coût : %d MP)", desc, m.MPCost)
		}
		items[i] = ui.MenuItem{Label: m.Name, Description: desc}
	}
	return items
}

func actMenuItems(options []boss.ActOption) []ui.MenuItem {
	items := make([]ui.MenuItem, len(options))
	for i, o := range options {
		items[i] = ui.MenuItem{Label: o.Label, Description: o.Description}
	}
	return items
}
