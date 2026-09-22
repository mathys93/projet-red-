package combat

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"ProjetRED/ascii"
	"ProjetRED/character"
	"ProjetRED/utilitaire"
)

func RunShop(player *character.Character) {
	ts := newTerminalSession()
	defer ts.restore()

	art, _ := ascii.Get("shop")

	for {
		idx, ok := shopChoose(ts, art, "Boutique", []menuItem{
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

func runMarchand(ts *terminalSession, player *character.Character, art string) {
	names := sortedKeys(utilitaire.MarchandObjets)

	for {
		items := make([]menuItem, len(names)+1)
		for i, n := range names {
			items[i] = menuItem{
				Label:       fmt.Sprintf("%s - %d Fragments", n, utilitaire.MarchandObjets[n]),
				Description: itemDescription(n),
			}
		}
		items[len(names)] = menuItem{Label: "Partir"}

		title := fmt.Sprintf("Marchand  (Fragments : %d)", player.Fragment)
		idx, ok := shopChoose(ts, art, title, items)
		if !ok || idx == len(names) {
			return
		}

		utilitaire.Acheter(player, names[idx])
		fmt.Println()
		fmt.Println(colWhite + "(Entrée pour continuer)" + colReset)
		ts.readKey()
	}
}

func runForgeron(ts *terminalSession, player *character.Character, art string) {
	names := sortedKeys(utilitaire.ForgeronObjets)

	for {
		items := make([]menuItem, len(names)+1)
		for i, n := range names {
			items[i] = menuItem{Label: fmt.Sprintf("%s - %d Fragments", n, utilitaire.ForgeronObjets[n])}
		}
		items[len(names)] = menuItem{Label: "Partir"}

		title := fmt.Sprintf("Forgeron  (Fragments : %d)", player.Fragment)
		idx, ok := shopChoose(ts, art, title, items)
		if !ok || idx == len(names) {
			return
		}

		utilitaire.Fabriquer(player, names[idx])
		fmt.Println()
		fmt.Println(colWhite + "(Entrée pour continuer)" + colReset)
		ts.readKey()
	}
}

func shopChoose(ts *terminalSession, art, title string, items []menuItem) (idx int, ok bool) {
	selected := 0
	for {
		clearScreen()
		bw, bh, aw, ah := layout()
		if needed := (len(items)+1)/2 + 2; needed > bh {
			bh = needed
		}
		shown := ascii.Fit(art, aw, ah)
		artHeight := len(strings.Split(shown, "\n"))
		printPadding(artHeight + bh + 6)
		fmt.Println()
		fmt.Println(colYellow + "  " + title + colReset)
		DrawArt(bw, shown, colYellow)
		DrawOptionGrid(bw, bh, items, selected)
		desc := ""
		if selected < len(items) {
			desc = items[selected].Description
		}
		fmt.Println(colWhite + "* " + desc + colReset)
		fmt.Println()
		fmt.Println(colWhite + "(← ↑ ↓ → pour naviguer, un chiffre pour choisir direct, Entrée pour valider, X pour quitter)" + colReset)

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

func sortedKeys(m map[string]int) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
