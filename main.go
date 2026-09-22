package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// ============================================================
// STRUCTURES
// ============================================================

type Character struct {
	Name  string
	Level int
	AP    int
	DP    int
	INV   []string
	LP    int
	MaxLP int
	XP    int
	MP    int
}

func NewCharacter(name string, level, ap, dp, lp int) *Character {
	return &Character{
		Name: name, Level: level, AP: ap, DP: dp,
		INV: []string{}, LP: lp, MaxLP: lp,
	}
}

func (c *Character) IsAlive() bool { return c.LP > 0 }

func (c *Character) TakeDamage(raw int) int {
	dmg := raw - c.DP
	if dmg < 0 {
		dmg = 0
	}
	c.LP -= dmg
	if c.LP < 0 {
		c.LP = 0
	}
	return dmg
}

func (c *Character) Heal(amount int) {
	c.LP += amount
	if c.LP > c.MaxLP {
		c.LP = c.MaxLP
	}
}

// ============================================================
// TERMINAL EN MODE "BRUT" (une touche = une action, sans Entrée)
// Utilise la commande système "stty" (présente sur Linux/macOS).
// ============================================================

func sttyRun(args ...string) {
	cmd := exec.Command("stty", args...)
	cmd.Stdin = os.Stdin
	cmd.Run()
}

func saveTermState() string {
	cmd := exec.Command("stty", "-g")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func enableRawMode() string {
	saved := saveTermState()
	sttyRun("cbreak", "-echo")
	return saved
}

func restoreTerm(saved string) {
	if saved != "" {
		sttyRun(saved)
	} else {
		sttyRun("sane")
	}
	fmt.Print("\033[?25h") // réaffiche le curseur
}

// Touches gérées
const (
	KeyUp    = "UP"
	KeyDown  = "DOWN"
	KeyLeft  = "LEFT"
	KeyRight = "RIGHT"
	KeyEnter = "ENTER"
	KeyQuit  = "QUIT"
	KeyOther = "OTHER"
)

func readKey(r *bufio.Reader) string {
	b, err := r.ReadByte()
	if err != nil {
		return KeyQuit
	}
	switch b {
	case 3: // Ctrl+C
		return KeyQuit
	case 13, 10: // Entrée
		return KeyEnter
	case 27: // séquence d'échappement (flèches)
		b2, _ := r.ReadByte()
		if b2 == '[' {
			b3, _ := r.ReadByte()
			switch b3 {
			case 'A':
				return KeyUp
			case 'B':
				return KeyDown
			case 'C':
				return KeyRight
			case 'D':
				return KeyLeft
			}
		}
		return KeyOther
	}
	return KeyOther
}

// ============================================================
// AFFICHAGE STYLE UNDERTALE
// ============================================================

const (
	colReset  = "\033[0m"
	colWhite  = "\033[97m"
	colYellow = "\033[93m"
	colRed    = "\033[91m"
)

func clearScreen() {
	fmt.Print("\033[2J\033[H\033[?25l") // efface + cache le curseur
}

// DrawBox dessine la boîte de combat (largeur fixe) avec le contenu
// (ton ASCII art de boss/perso) centré dedans.
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

// DrawStatBar affiche la barre "Nom  LV x  HP [■■■□□] cur/max" façon Undertale.
func DrawStatBar(c *Character) {
	barLen := 20
	filled := 0
	if c.MaxLP > 0 {
		filled = (c.LP * barLen) / c.MaxLP
	}
	if filled > barLen {
		filled = barLen
	}
	bar := colYellow + strings.Repeat("■", filled) + colWhite + strings.Repeat("□", barLen-filled) + colReset

	fmt.Printf(" %-12s LV %-3d HP %s %d/%d\n", c.Name, c.Level, bar, c.LP, c.MaxLP)
}

// DrawMenu affiche FIGHT / ACT / ITEM / MERCY avec le cœur devant l'option
// sélectionnée (déplacé avec les flèches gauche/droite).
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
func RenderBattleScreen(enemyArt string, enemy *Character, player *Character, selected int, message string) {
	clearScreen()
	fmt.Println()
	DrawBox(50, 12, enemyArt) // <- ton ASCII art de boss vient ici
	fmt.Println()
	if message != "" {
		fmt.Println(colWhite + "* " + message + colReset)
	}
	DrawStatBar(player)
	DrawMenu(selected)
	fmt.Println(colWhite + "\n(← →  pour choisir, Entrée pour valider, Ctrl+C pour quitter)" + colReset)
}

// ============================================================
// BOUCLE DE COMBAT INTERACTIVE
// ============================================================

func RunBattle(player *Character, enemy *Character, enemyArt string) {
	saved := enableRawMode()
	defer restoreTerm(saved)
	reader := bufio.NewReader(os.Stdin)

	selected := 0
	message := fmt.Sprintf("%s bloque le passage !", enemy.Name)

	for player.IsAlive() && enemy.IsAlive() {
		RenderBattleScreen(enemyArt, enemy, player, selected, message)

		key := readKey(reader)
		switch key {
		case KeyLeft:
			selected = (selected + 3) % 4
		case KeyRight:
			selected = (selected + 1) % 4
		case KeyQuit:
			return
		case KeyEnter:
			switch selected {
			case 0: // FIGHT
				dmg := enemy.TakeDamage(player.AP)
				message = fmt.Sprintf("Tu attaques %s ! %d dégâts.", enemy.Name, dmg)
			case 1: // ACT
				message = fmt.Sprintf("Tu observes %s... (à personnaliser par boss)", enemy.Name)
			case 2: // ITEM
				if len(player.INV) == 0 {
					message = "Ton inventaire est vide !"
				} else {
					item := player.INV[0]
					player.INV = player.INV[1:]
					player.Heal(10)
					message = fmt.Sprintf("Tu utilises %s. LP restauré.", item)
				}
			case 3: // MERCY
				message = "Tu tentes d'épargner l'ennemi..."
				RenderBattleScreen(enemyArt, enemy, player, selected, message)
				restoreTerm(saved)
				return
			}

			if !enemy.IsAlive() {
				message = fmt.Sprintf("%s est terrassé(e) ! Victoire.", enemy.Name)
				RenderBattleScreen(enemyArt, enemy, player, selected, message)
				restoreTerm(saved)
				return
			}

			// Riposte simple du boss (remplace ensuite par la phase d'esquive
			// avec ton ASCII art / patterns propres à chaque boss)
			if selected != 3 {
				dmgTaken := player.TakeDamage(enemy.AP)
				message += fmt.Sprintf("  %s riposte : %d dégâts.", enemy.Name, dmgTaken)
			}

			if !player.IsAlive() {
				message = "Tu es tombé(e) au combat..."
				RenderBattleScreen(enemyArt, enemy, player, selected, message)
				restoreTerm(saved)
				return
			}
		}
	}
}

// ============================================================
// MAIN
// ============================================================

func main() {
	hero := NewCharacter("Personnage 1", 1, 5, 2, 20)
	hero.INV = append(hero.INV, "Potion")

	boss := NewCharacter("Boss Anime", 1, 4, 1, 30)

	// Remplace ceci par ton ASCII art du boss (chaîne multi-lignes)
	enemyArt := "  (ᵔᴥᵔ)\n  ton ASCII\n  art ici"

	RunBattle(hero, boss, enemyArt)
}