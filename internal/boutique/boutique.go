package boutique

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"ProjetRED/internal/ascii"
	"ProjetRED/internal/character"
	"ProjetRED/internal/ui"
)

func RunShop(player *character.Character) {
	ui.SetScene(ui.SceneBoutique)
	ts := ui.NewSession()
	defer ts.Restore()

	art, _ := ascii.Get("shop")

	for {
		idx, ok := shopChoose(ts, art, "Boutique", []ui.MenuItem{
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

func runMarchand(ts *ui.Session, player *character.Character, art string) {
	names := sortedKeys(MarchandObjets)

	for {
		items := make([]ui.MenuItem, len(names)+1)
		for i, n := range names {
			items[i] = ui.MenuItem{
				Label:       fmt.Sprintf("%s - %d Fragments", n, MarchandObjets[n]),
				Description: Description(n),
			}
		}
		items[len(names)] = ui.MenuItem{Label: "Partir"}

		title := fmt.Sprintf("Marchand  (Fragments : %d)", player.Fragment)
		idx, ok := shopChoose(ts, art, title, items)
		if !ok || idx == len(names) {
			return
		}

		Acheter(player, names[idx])
		fmt.Fprintln(ui.Out)
		fmt.Fprintln(ui.Out, ui.ColWhite+"(Entrée pour continuer)"+ui.ColReset)
		ts.ReadKey()
	}
}

func runForgeron(ts *ui.Session, player *character.Character, art string) {
	names := sortedKeys(ForgeronObjets)

	for {
		items := make([]ui.MenuItem, len(names)+1)
		for i, n := range names {
			items[i] = ui.MenuItem{Label: fmt.Sprintf("%s - %d Fragments", n, ForgeronObjets[n])}
		}
		items[len(names)] = ui.MenuItem{Label: "Partir"}

		title := fmt.Sprintf("Forgeron  (Fragments : %d)", player.Fragment)
		idx, ok := shopChoose(ts, art, title, items)
		if !ok || idx == len(names) {
			return
		}

		Fabriquer(player, names[idx])
		fmt.Fprintln(ui.Out)
		fmt.Fprintln(ui.Out, ui.ColWhite+"(Entrée pour continuer)"+ui.ColReset)
		ts.ReadKey()
	}
}

func shopChoose(ts *ui.Session, art, title string, items []ui.MenuItem) (idx int, ok bool) {
	selected := 0
	render := func() {
		ui.ClearScreen()
		bw, bh, aw, ah := ui.Layout()
		if needed := (len(items)+1)/2 + 2; needed > bh {
			bh = needed
		}
		shown := ascii.Fit(art, aw, ah)
		artHeight := len(strings.Split(shown, "\n"))
		ui.PrintPadding(artHeight + bh + 6)
		fmt.Fprintln(ui.Out)
		fmt.Fprintln(ui.Out, ui.ColYellow+"  "+title+ui.ColReset)
		ui.DrawArt(bw, shown, ui.ColYellow)
		ui.DrawOptionGrid(bw, bh, items, selected)
		desc := ""
		if selected < len(items) {
			desc = items[selected].Description
		}
		fmt.Fprintln(ui.Out, ui.ColWhite+"* "+desc+ui.ColReset)
		fmt.Fprintln(ui.Out)
		fmt.Fprintln(ui.Out, ui.ColWhite+"(← ↑ ↓ → pour naviguer, un chiffre pour choisir direct, Entrée pour valider, X pour revenir)"+ui.ColReset)
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
