// Package combat contient toute l'interface de combat façon Undertale :
// la boîte de combat avec l'artwork ASCII du boss, la barre de vie, le menu
// FIGHT / ACT / ITEM / MERCY navigable au clavier, et la boucle de combat.
package combat

import (
	"fmt"
	"strings"

	"ProjetRED/boss"
	"ProjetRED/character"
)

// Couleurs ANSI (façon boîte de dialogue Undertale : blanc/jaune/rouge).
const (
	colReset  = "\033[0m"
	colWhite  = "\033[97m"
	colYellow = "\033[93m"
	colRed    = "\033[91m"
)

func clearScreen() {
	fmt.Print("\033[2J\033[H\033[?25l") // efface l'écran + cache le curseur
}

// DrawBox dessine la boîte de combat (largeur/hauteur fixes) avec l'artwork
// ASCII du boss centré dedans.
func DrawBox(width, height int, art string) {
	fmt.Println(colWhite + "┌" + strings.Repeat("─", width-2) + "┐" + colReset)

	lines := strings.Split(art, "\n")
	for len(lines) < height-2 {
		lines = append(lines, "")
	}
	if len(lines) > height-2 {
		lines = lines[:height-2]
	}

	for _, line := range lines {
		pad := (width - 2 - len([]rune(line))) / 2
		if pad < 0 {
			pad = 0
		}
		right := width - 2 - pad - len([]rune(line))
		if right < 0 {
			right = 0
		}
		fmt.Println(colWhite + "│" + colReset +
			strings.Repeat(" ", pad) + line + strings.Repeat(" ", right) +
			colWhite + "│" + colReset)
	}

	fmt.Println(colWhite + "└" + strings.Repeat("─", width-2) + "┘" + colReset)
}

// DrawStatBar affiche "Nom  LV x  HP [■■■□□] cur/max  MP cur/max" façon
// Undertale (le MP en plus sert aux attaques de Stand du menu FIGHT).
func DrawStatBar(c *character.Character) {
	barLen := 20
	filled := 0
	if c.MaxLP > 0 {
		filled = (c.LP * barLen) / c.MaxLP
	}
	if filled > barLen {
		filled = barLen
	}
	bar := colYellow + strings.Repeat("■", filled) + colWhite + strings.Repeat("□", barLen-filled) + colReset

	fmt.Printf(" %-14s LV %-3d HP %s %d/%d   MP %d/%d\n", c.Name, c.Level, bar, c.LP, c.MaxLP, c.MP, c.MaxMP)
}

// DrawMenu affiche FIGHT / ACT / ITEM / MERCY avec le cœur devant l'option
// sélectionnée. Chaque option se déplace avec Q/D (ou Z/S), se valide avec
// Entrée, OU se choisit directement avec son raccourci (1/F, 2/A, 3/I,
// 4/M) en une seule saisie.
func DrawMenu(selected int) {
	options := []string{"FIGHT", "ACT", "ITEM", "MERCY"}
	hotkeys := []string{"1/F", "2/A", "3/I", "4/M"}
	fmt.Println()
	line := ""
	for i, opt := range options {
		if i == selected {
			line += colRed + "❤ " + colYellow + opt + colReset
		} else {
			line += "  " + colWhite + opt + colReset
		}
		line += colWhite + "(" + hotkeys[i] + ")   " + colReset
	}
	fmt.Println(line)
}

// RenderBattleScreen redessine l'écran complet (boîte + menu d'actions) à
// chaque frame où c'est au joueur de choisir son action.
func RenderBattleScreen(b *boss.Boss, player *character.Character, selected int, message string) {
	clearScreen()
	fmt.Println()
	fmt.Println(colYellow + "  " + b.Zone + colReset)
	DrawBox(boss.ArtWidth+4, boss.ArtHeight+2, b.Art)
	fmt.Println()
	if message != "" {
		fmt.Println(colWhite + "* " + message + colReset)
	}
	DrawStatBar(player)
	DrawMenu(selected)
	fmt.Println(colWhite + "\n(Q/D ou raccourci direct pour choisir, Entrée pour valider, X pour quitter)" + colReset)
}

// RenderTurnMessage affiche la boîte de combat et un message SANS le menu
// d'actions : utilisé pour bien séparer visuellement le tour du joueur et
// celui du boss (façon Undertale, où le menu disparaît pendant les
// attaques), plutôt que de résoudre les deux tours d'un coup dans le même
// message comme avant.
func RenderTurnMessage(b *boss.Boss, player *character.Character, message string) {
	clearScreen()
	fmt.Println()
	fmt.Println(colYellow + "  " + b.Zone + colReset)
	DrawBox(boss.ArtWidth+4, boss.ArtHeight+2, b.Art)
	fmt.Println()
	if message != "" {
		fmt.Println(colWhite + "* " + message + colReset)
	}
	DrawStatBar(player)
	fmt.Println(colWhite + "\n(Entrée pour continuer, X pour quitter)" + colReset)
}

// waitContinue attend une validation du joueur pour laisser le temps de
// lire le message affiché par RenderTurnMessage. Renvoie false si le
// joueur quitte (X) pendant cette pause, pour que RunBattle propage la
// sortie du combat.
func waitContinue(ts *terminalSession) bool {
	return ts.readKey() != KeyQuit
}

// Result indique comment un combat s'est terminé.
type Result int

const (
	ResultVictory Result = iota
	ResultDefeat
	ResultSpared
	ResultQuit
)

// mainMenuHotkeys associe un raccourci direct (chiffre ou lettre, en plus
// de la navigation Q/D + Entrée) à chaque option du menu FIGHT/ACT/ITEM/
// MERCY, pour pouvoir sauter directement dessus en une seule saisie
// (ex: taper "a" va droit à ACT sans avoir à naviguer jusque-là).
var mainMenuHotkeys = map[string]int{
	"f": 0, "1": 0,
	"a": 1, "2": 1,
	"i": 2, "3": 2,
	"m": 3, "4": 3,
}

// RunBattle lance un combat interactif façon Undertale entre le joueur et
// un boss, et renvoie comment le combat s'est terminé.
//
// Chaque tour se déroule en deux temps bien séparés à l'écran (au lieu
// d'être résolus d'un coup dans le même message) : d'abord le tour du
// joueur (choix dans FIGHT/ACT/ITEM/MERCY puis résultat affiché seul),
// ensuite - une fois validé - le tour du boss (délégué à son Pattern, voir
// boss/patterns.go), affiché séparément avant de repasser la main au
// joueur.
func RunBattle(player *character.Character, b *boss.Boss) Result {
	ts := newTerminalSession()
	defer ts.restore()

	selected := 0
	turn := 1
	message := fmt.Sprintf("%s bloque le passage !", b.Name)

	for player.IsAlive() && b.IsAlive() {
		RenderBattleScreen(b, player, selected, message)

		key := ts.readKey()
		chosen := -1
		switch key {
		case KeyLeft, KeyUp:
			selected = (selected + 3) % 4
		case KeyRight, KeyDown:
			selected = (selected + 1) % 4
		case KeyQuit:
			return ResultQuit
		case KeyEnter:
			chosen = selected
		case KeyOther:
			if idx, ok := mainMenuHotkeys[ts.raw]; ok {
				selected = idx
				chosen = idx
			}
		}
		if chosen == -1 {
			continue
		}

		acted := true
		switch chosen {
		case 0: // FIGHT : choix de l'attaque (coup de poing, Stand...).
			idx, ok := chooseOption(ts, fmt.Sprintf("%s - choisis ton attaque", player.Name), fightMenuItems(player))
			if !ok {
				acted = false
				break
			}
			move := player.Moves[idx]
			if player.MP < move.MPCost {
				message = "Pas assez de MP pour cette action !"
				acted = false
				break
			}
			player.MP -= move.MPCost
			message = move.Perform(player, b.Character)
		case 1: // ACT : choix de l'action "lore" (parler, observer...).
			options := b.ActOptions(player)
			idx, ok := chooseOption(ts, fmt.Sprintf("%s - que fais-tu ?", b.Name), actMenuItems(options))
			if !ok {
				acted = false
				break
			}
			message = options[idx].Resolve(b, player)
		case 2: // ITEM
			if len(player.Inventory) == 0 {
				message = "Ton inventaire est vide !"
				acted = false
				break
			}
			item := player.Inventory[0]
			player.RemoveItem(item)
			player.Heal(10)
			message = fmt.Sprintf("Tu utilises %s. LP restauré.", item)
		case 3: // MERCY
			message = fmt.Sprintf("Tu épargnes %s...", b.Name)
			RenderBattleScreen(b, player, selected, message)
			return ResultSpared
		}

		if !acted {
			// Choix annulé (X) ou impossible (pas assez de MP, inventaire
			// vide...) : on reste au menu principal, le tour du boss
			// n'est pas déclenché.
			continue
		}

		RenderTurnMessage(b, player, message)
		if !waitContinue(ts) {
			return ResultQuit
		}

		if !b.IsAlive() {
			RenderTurnMessage(b, player, fmt.Sprintf("%s est terrassé(e) ! Victoire.", b.Name))
			waitContinue(ts)
			return ResultVictory
		}

		// Le tour du boss, affiché à part : c'est là que chaque boss aura
		// son propre comportement (voir boss/patterns.go).
		RenderTurnMessage(b, player, fmt.Sprintf("-- Tour de %s --", b.Name))
		if !waitContinue(ts) {
			return ResultQuit
		}
		bossMessage := b.Turn(turn, player)
		turn++
		RenderTurnMessage(b, player, bossMessage)
		if !waitContinue(ts) {
			return ResultQuit
		}

		if !player.IsAlive() {
			RenderTurnMessage(b, player, "Tu es tombé(e) au combat...")
			waitContinue(ts)
			return ResultDefeat
		}

		message = ""
	}
	return ResultDefeat
}
