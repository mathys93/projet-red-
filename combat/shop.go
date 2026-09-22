package combat

import (
	"fmt"
	"sort"
	"strconv"

	"ProjetRED/ascii"
	"ProjetRED/boss"
	"ProjetRED/character"
	"ProjetRED/utilitaire"
)

// RunShop affiche l'écran de la boutique (Zone 4 sur la carte, voir
// world.New) : marchand / forgeron / partir, et boucle jusqu'à ce que le
// joueur choisisse de partir (X ou "Partir").
func RunShop(player *character.Character) {
	ts := newTerminalSession()
	defer ts.restore()

	art, _ := ascii.Get("shop")
	art = ascii.Fit(art, boss.ArtWidth, boss.ArtHeight)

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

// runMarchand affiche le catalogue du marchand et gère l'achat, comme un
// magasin façon Undertale : liste des objets et leur prix, confirmation à
// l'écran après achat.
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

// runForgeron affiche le catalogue du forgeron et gère la fabrication.
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

// shopChoose affiche l'artwork de la boutique (teinté en jaune) suivi d'une
// liste d'options navigable, sur le même modèle que chooseOption (voir
// combat/submenu.go) mais avec la boîte ASCII en plus.
func shopChoose(ts *terminalSession, art, title string, items []menuItem) (idx int, ok bool) {
	selected := 0
	for {
		clearScreen()
		fmt.Println()
		fmt.Println(colYellow + "  " + title + colReset)
		DrawBox(boss.ArtWidth+4, boss.ArtHeight+2, art, colYellow)
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
		fmt.Println(colWhite + "(Z/S pour naviguer, un chiffre pour choisir direct, Entrée pour valider, X pour quitter)" + colReset)

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

// sortedKeys renvoie les clés d'une map d'objets, triées par ordre
// alphabétique pour un affichage stable (l'ordre d'itération d'une map Go
// n'est pas garanti).
func sortedKeys(m map[string]int) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
