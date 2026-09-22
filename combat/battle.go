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

// DrawStatBar affiche "Nom  LV x  HP [■■■□□] cur/max" façon Undertale.
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

	fmt.Printf(" %-14s LV %-3d HP %s %d/%d\n", c.Name, c.Level, bar, c.LP, c.MaxLP)
}

// DrawMenu affiche FIGHT / ACT / ITEM / MERCY avec le cœur devant l'option
// sélectionnée (déplacé avec Q/D ou A/D, tape la commande puis Entrée).
func DrawMenu(selected int) {
	options := []string{"FIGHT", "ACT", "ITEM", "MERCY"}
	fmt.Println()
	line := ""
	for i, opt := range options {
		if i == selected {
			line += colRed + "❤ " + colYellow + opt + colReset + "   "
		} else {
			line += "  " + colWhite + opt + colReset + "   "
		}
	}
	fmt.Println(line)
}

// RenderBattleScreen redessine l'écran complet à chaque frame.
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
	fmt.Println(colWhite + "\n(Q/D pour choisir, Entrée seule pour valider, X pour quitter)" + colReset)
}

// Result indique comment un combat s'est terminé.
type Result int

const (
	ResultVictory Result = iota
	ResultDefeat
	ResultSpared
	ResultQuit
)

// RunBattle lance un combat interactif façon Undertale entre le joueur et
// un boss, et renvoie comment le combat s'est terminé.
func RunBattle(player *character.Character, b *boss.Boss) Result {
	ts := newTerminalSession()
	defer ts.restore()

	selected := 0
	turn := 1
	message := fmt.Sprintf("%s bloque le passage !", b.Name)

	for player.IsAlive() && b.IsAlive() {
		RenderBattleScreen(b, player, selected, message)

		key := ts.readKey()
		switch key {
		case KeyLeft, KeyUp:
			selected = (selected + 3) % 4
		case KeyRight, KeyDown:
			selected = (selected + 1) % 4
		case KeyQuit:
			return ResultQuit
		case KeyEnter:
			switch selected {
			case 0: // FIGHT
				dmg := b.TakeDamage(player.AP)
				message = fmt.Sprintf("Tu attaques %s ! %d dégâts.", b.Name, dmg)
			case 1: // ACT
				message = b.ActDescription(player)
			case 2: // ITEM
				if len(player.Inventory) == 0 {
					message = "Ton inventaire est vide !"
				} else {
					item := player.Inventory[0]
					player.RemoveItem(item)
					player.Heal(10)
					message = fmt.Sprintf("Tu utilises %s. LP restauré.", item)
				}
			case 3: // MERCY
				message = fmt.Sprintf("Tu épargnes %s...", b.Name)
				RenderBattleScreen(b, player, selected, message)
				return ResultSpared
			}

			if !b.IsAlive() {
				message = fmt.Sprintf("%s est terrassé(e) ! Victoire.", b.Name)
				RenderBattleScreen(b, player, selected, message)
				return ResultVictory
			}

			// Le tour du boss est délégué à son Pattern (voir boss/patterns.go) :
			// c'est là que chaque boss aura son propre comportement.
			if selected != 3 {
				message += "  " + b.Turn(turn, player)
				turn++
			}

			if !player.IsAlive() {
				message = "Tu es tombé(e) au combat..."
				RenderBattleScreen(b, player, selected, message)
				return ResultDefeat
			}
		}
	}
	return ResultDefeat
}
