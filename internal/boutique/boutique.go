package boutique

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"ProjetRED/internal/ascii"
	"ProjetRED/internal/character"
	"ProjetRED/internal/fond"
	"ProjetRED/internal/terminal"
)

func RunShop(player *character.Character) {
	terminal.SetScene(fond.Boutique)
	ts := terminal.NewSession()
	defer ts.Restore()

	art, _ := ascii.Get("shop")

	for {
		idx, ok := shopChoose(ts, art, "Boutique", []terminal.MenuItem{
			{Label: "Marchand", Description: "Objets consommables, payés en Fragments."},
			{Label: "Forgeron", Description: "Équipements fabriqués, payés en Fragments + ressources."},
			{Label: "Partir", Description: "Retourner à la carte."},
		})
		if !ok || idx == 2 {
			return
		}
		if idx == 0 {
			runMarchand(ts, player, art)
		} else {
			runForgeron(ts, player, art)
		}
	}
}

func runMarchand(ts *terminal.Session, player *character.Character, art string) {
	names := sortedKeys(MarchandObjets)

	for {
		items := make([]terminal.MenuItem, len(names)+1)
		for i, n := range names {
			items[i] = terminal.MenuItem{
				Label:       fmt.Sprintf("%s - %d Fragments", n, MarchandObjets[n]),
				Description: Description(n),
			}
		}
		items[len(names)] = terminal.MenuItem{Label: "Partir"}

		title := fmt.Sprintf("Marchand  (Fragments : %d)", player.Fragment)
		idx, ok := shopChoose(ts, art, title, items)
		if !ok || idx == len(names) {
			return
		}

		Acheter(player, names[idx])
		fmt.Fprintln(terminal.Out)
		fmt.Fprintln(terminal.Out, terminal.ColWhite+"(Entrée pour continuer)"+terminal.ColReset)
		ts.ReadKey()
	}
}

func runForgeron(ts *terminal.Session, player *character.Character, art string) {
	names := sortedKeys(ForgeronObjets)

	for {
		items := make([]terminal.MenuItem, len(names)+1)
		for i, n := range names {
			items[i] = terminal.MenuItem{Label: fmt.Sprintf("%s - %d Fragments", n, ForgeronObjets[n])}
		}
		items[len(names)] = terminal.MenuItem{Label: "Partir"}

		title := fmt.Sprintf("Forgeron  (Fragments : %d)", player.Fragment)
		idx, ok := shopChoose(ts, art, title, items)
		if !ok || idx == len(names) {
			return
		}

		Fabriquer(player, names[idx])
		fmt.Fprintln(terminal.Out)
		fmt.Fprintln(terminal.Out, terminal.ColWhite+"(Entrée pour continuer)"+terminal.ColReset)
		ts.ReadKey()
	}
}

func shopChoose(ts *terminal.Session, art, title string, items []terminal.MenuItem) (idx int, ok bool) {
	selected := 0
	render := func() {
		terminal.ClearScreen()
		bw, bh, aw, ah := terminal.Layout()
		if needed := (len(items)+1)/2 + 2; needed > bh {
			bh = needed
		}
		shown := ascii.Fit(art, aw, ah)
		artHeight := len(strings.Split(shown, "\n"))
		terminal.PrintPadding(artHeight + bh + 6)
		fmt.Fprintln(terminal.Out)
		fmt.Fprintln(terminal.Out, terminal.ColYellow+"  "+title+terminal.ColReset)
		terminal.DrawArt(bw, shown, terminal.ColYellow)
		terminal.DrawOptionGrid(bw, bh, items, selected)
		desc := ""
		if selected < len(items) {
			desc = items[selected].Description
		}
		fmt.Fprintln(terminal.Out, terminal.ColWhite+"* "+desc+terminal.ColReset)
		fmt.Fprintln(terminal.Out)
		fmt.Fprintln(terminal.Out, terminal.ColWhite+"(← ↑ ↓ → pour naviguer, un chiffre pour choisir direct, Entrée pour valider, X pour revenir)"+terminal.ColReset)
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

func sortedKeys(m map[string]int) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func Description(item string) string {
	switch item {
	case "Potion de vie":
		return "Restaure 20 PV."
	case "Potion de MP":
		return "Restaure 15 MP."
	case "Disque de Pucci":
		return "Restaure 50 PV."
	case "Tel Diavolo":
		return "Inflige des dégâts, ignore la Defense Power."
	case "Arrow":
		return "Éveille un Stand permanent : bonus d'attaque et nouvelle action de combat."
	}
	return ""
}
