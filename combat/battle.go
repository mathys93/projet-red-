package combat

import (
	"fmt"
	"strings"
	"time"

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

// MaximizeConsoleWindow agrandit la fenêtre de la console au maximum, si
// possible (voir combat/termsize_windows.go). Appelée une fois au tout
// début du programme (voir main.go).
func MaximizeConsoleWindow() {
	maximizeConsoleWindow()
}

// terminalSize renvoie la taille (colonnes, lignes) du terminal actuel,
// avec un repli raisonnable quand elle ne peut pas être détectée (terminal
// sans fenêtre propre, sortie redirigée...), pour que l'affichage puisse
// toujours calculer une mise en page plutôt que planter.
func terminalSize() (cols, rows int) {
	cols, rows = queryTerminalSize()
	if cols < 40 {
		cols = 80
	}
	if rows < 20 {
		rows = 24
	}
	return cols, rows
}

// boxWidth calcule la largeur de la boîte de combat/boutique pour qu'elle
// remplisse (quasi) toute la largeur du terminal, façon Undertale en plein
// écran, plutôt que de rester à sa taille minimale dans un coin de l'écran.
// Elle ne descend jamais sous minWidth (la taille requise par l'artwork).
func boxWidth(minWidth int) int {
	cols, _ := terminalSize()
	w := cols - 4
	if w < minWidth {
		w = minWidth
	}
	return w
}

// printPadding centre verticalement le contenu d'un écran (contentHeight
// lignes) en imprimant des lignes vides au-dessus, plutôt que de le laisser
// collé en haut du terminal avec tout le reste de l'écran vide en dessous.
func printPadding(contentHeight int) {
	_, rows := terminalSize()
	pad := (rows - contentHeight) / 2
	for i := 0; i < pad; i++ {
		fmt.Println()
	}
}

// visibleWidth renvoie la largeur "à l'écran" d'une chaîne pouvant contenir
// des codes couleur ANSI (\033[...m), en ignorant ces codes pour le calcul.
// Sans ça, tout padding calculé sur une chaîne déjà colorée serait faussé
// par les octets de couleur (invisibles mais comptés par len/[]rune).
func visibleWidth(s string) int {
	n := 0
	inEsc := false
	for _, r := range s {
		switch {
		case inEsc:
			if r == 'm' {
				inEsc = false
			}
		case r == '\033':
			inEsc = true
		default:
			n++
		}
	}
	return n
}

// drawFrameTop/drawFrameLine/drawFrameBottom dessinent le cadre de combat,
// brique par brique : un cadre vide (voir DrawEmptyBox) ou rempli d'une
// grille d'options (voir DrawOptionGrid) sont tous les deux construits à
// partir de ces trois fonctions.
func drawFrameTop(width int) {
	fmt.Println(colWhite + "┌" + strings.Repeat("─", width-2) + "┐" + colReset)
}

func drawFrameBottom(width int) {
	fmt.Println(colWhite + "└" + strings.Repeat("─", width-2) + "┘" + colReset)
}

// drawFrameLine imprime une ligne de contenu (déjà colorée) à l'intérieur
// du cadre, complétée par des espaces jusqu'à occuper toute la largeur.
func drawFrameLine(width int, content string) {
	pad := width - 2 - visibleWidth(content)
	if pad < 0 {
		pad = 0
	}
	fmt.Println(colWhite + "│" + colReset + content + strings.Repeat(" ", pad) + colWhite + "│" + colReset)
}

// DrawArt affiche l'artwork du boss centré, SEUL, au-dessus du cadre de
// combat (voir DrawEmptyBox/DrawOptionGrid pour le cadre lui-même). Avant,
// l'artwork était dessiné à l'intérieur du cadre, qui restait donc
// inutilisable pour autre chose : il est maintenant sorti au-dessus, et le
// cadre sert de zone d'affichage pour les sous-menus FIGHT/ACT/ITEM.
func DrawArt(width int, art string, artColor string) {
	if artColor == "" {
		artColor = colWhite
	}
	for _, line := range strings.Split(art, "\n") {
		pad := (width - len([]rune(line))) / 2
		if pad < 0 {
			pad = 0
		}
		fmt.Println(strings.Repeat(" ", pad) + artColor + line + colReset)
	}
}

// DrawEmptyBox dessine un cadre vide de la taille donnée : c'est ce que
// devient la boîte de combat en dehors des sous-menus (voir DrawArt pour
// où est passé l'artwork).
func DrawEmptyBox(width, height int) {
	drawFrameTop(width)
	for i := 0; i < height-2; i++ {
		drawFrameLine(width, "")
	}
	drawFrameBottom(width)
}

// DrawOptionGrid dessine les entrées d'un sous-menu (coups de FIGHT,
// options ACT, objets d'ITEM...) en grille à deux colonnes À L'INTÉRIEUR
// du cadre de combat, façon inventaire Undertale (ex : Potion de vie en
// haut à droite, Disque de Pucci en haut à gauche...), plutôt qu'une
// simple liste verticale. La description de l'entrée sélectionnée est
// affichée par l'appelant, sous le cadre (voir combat/submenu.go).
func DrawOptionGrid(width, height int, items []menuItem, selected int) {
	drawFrameTop(width)

	const cols = 2
	colWidth := (width - 2 - (cols + 1)) / cols
	rows := (len(items) + cols - 1) / cols

	interior := height - 2
	top := (interior - rows) / 2
	if top < 0 {
		top = 0
	}

	for i := 0; i < top; i++ {
		drawFrameLine(width, "")
	}
	for r := 0; r < rows; r++ {
		line := " "
		for c := 0; c < cols; c++ {
			idx := r*cols + c
			cell := strings.Repeat(" ", colWidth)
			if idx < len(items) {
				marker := "  "
				col := colYellow
				if idx == selected {
					marker = colRed + "❤ " + colReset
					col = colRed
				}
				text := marker + col + items[idx].Label + colReset
				pad := colWidth - visibleWidth(text)
				if pad < 0 {
					pad = 0
				}
				cell = text + strings.Repeat(" ", pad)
			}
			line += cell + " "
		}
		drawFrameLine(width, line)
	}
	for i := 0; i < interior-top-rows; i++ {
		drawFrameLine(width, "")
	}

	drawFrameBottom(width)
}

// DrawActionButtons affiche FIGHT / ACT / ITEM / MERCY sous forme de gros
// boutons encadrés, répartis sur toute la largeur du cadre (façon boutons
// d'action d'Undertale), plutôt qu'une simple ligne de texte compacte.
// Le bouton sélectionné est mis en évidence en rouge avec un cœur.
func DrawActionButtons(width int, selected int) {
	options := []string{"FIGHT", "ACT", "ITEM", "MERCY"}
	hotkeys := []string{"1/F", "2/A", "3/I", "4/M"}

	labels := make([]string, len(options))
	inner := 0
	for i, opt := range options {
		labels[i] = fmt.Sprintf("%s (%s)", opt, hotkeys[i])
		if n := len([]rune(labels[i])); n > inner {
			inner = n
		}
	}
	inner += 4 // marge intérieure du bouton (cœur/marqueur + espacement)

	btnOuter := inner + 2 // + les deux bordures verticales du bouton
	gap := (width - len(options)*btnOuter) / (len(options) + 1)
	if gap < 2 {
		gap = 2
	}
	margin := strings.Repeat(" ", gap)

	var top, mid, bot strings.Builder
	for i, label := range labels {
		col := colWhite
		heart := "  "
		if i == selected {
			col = colRed
			heart = colRed + "❤ " + colReset
		}
		pad := inner - 2 - len([]rune(label))
		if pad < 0 {
			pad = 0
		}
		left := pad / 2
		right := pad - left

		top.WriteString(margin + col + "┌" + strings.Repeat("─", inner) + "┐" + colReset)
		mid.WriteString(margin + col + "│" + colReset + heart +
			strings.Repeat(" ", left) + col + label + colReset + strings.Repeat(" ", right) +
			col + "│" + colReset)
		bot.WriteString(margin + col + "└" + strings.Repeat("─", inner) + "┘" + colReset)
	}
	top.WriteString(margin)
	mid.WriteString(margin)
	bot.WriteString(margin)

	fmt.Println()
	fmt.Println(top.String())
	fmt.Println(mid.String())
	fmt.Println(bot.String())
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

	fmt.Printf(" %-14s LV %-3d HP %s %d/%d   MP %d/%d   XP %d/%d\n",
		c.Name, c.Level, bar, c.LP, c.MaxLP, c.MP, c.MaxMP, c.XP, character.XPForLevel(c.Level))
}

// DrawEnemyBar affiche la barre de PV du boss (façon barre de vie
// d'adversaire), pour qu'on voie clairement ses PV baisser pendant le
// combat, comme pour le joueur avec DrawStatBar.
func DrawEnemyBar(b *boss.Boss) {
	barLen := 20
	filled := 0
	if b.MaxLP > 0 {
		filled = (b.LP * barLen) / b.MaxLP
	}
	if filled > barLen {
		filled = barLen
	}
	if filled < 0 {
		filled = 0
	}
	bar := colRed + strings.Repeat("■", filled) + colWhite + strings.Repeat("□", barLen-filled) + colReset

	fmt.Printf(" %-14s LV %-3d HP %s %d/%d\n", b.Name, b.Level, bar, b.LP, b.MaxLP)
}

// RenderBattleScreen redessine l'écran complet (artwork + cadre vide +
// gros boutons d'action) à chaque frame où c'est au joueur de choisir son
// action.
func RenderBattleScreen(b *boss.Boss, player *character.Character, selected int, message string) {
	clearScreen()
	bw := boxWidth(boss.ArtWidth + 4)
	bh := boss.ArtHeight + 2
	artHeight := len(strings.Split(b.Art, "\n"))
	contentHeight := artHeight + bh + 11
	if message != "" {
		contentHeight++
	}
	printPadding(contentHeight)
	fmt.Println()
	fmt.Println(colYellow + "  " + b.Zone + colReset)
	DrawArt(bw, b.Art, b.Color)
	DrawEmptyBox(bw, bh)
	fmt.Println()
	DrawEnemyBar(b)
	if message != "" {
		fmt.Println(colWhite + "* " + message + colReset)
	}
	DrawStatBar(player)
	DrawActionButtons(bw, selected)
	fmt.Println(colWhite + "\n(Q/D ou raccourci direct pour choisir, Entrée pour valider, X pour quitter)" + colReset)
}

// RenderTurnMessage affiche l'artwork, le cadre vide et un message SANS le
// menu d'actions : utilisé pour bien séparer visuellement le tour du
// joueur et celui du boss (façon Undertale, où le menu disparaît pendant
// les attaques), plutôt que de résoudre les deux tours d'un coup dans le
// même message comme avant.
func RenderTurnMessage(b *boss.Boss, player *character.Character, message string) {
	clearScreen()
	bw := boxWidth(boss.ArtWidth + 4)
	bh := boss.ArtHeight + 2
	artHeight := len(strings.Split(b.Art, "\n"))
	contentHeight := artHeight + bh + 7
	if message != "" {
		contentHeight++
	}
	printPadding(contentHeight)
	fmt.Println()
	fmt.Println(colYellow + "  " + b.Zone + colReset)
	DrawArt(bw, b.Art, b.Color)
	DrawEmptyBox(bw, bh)
	fmt.Println()
	DrawEnemyBar(b)
	if message != "" {
		fmt.Println(colWhite + "* " + message + colReset)
	}
	DrawStatBar(player)
	fmt.Println(colWhite + "\n(Entrée pour continuer, X pour quitter)" + colReset)
}

// animateHPChange anime la barre de PV du joueur (isPlayer = true) ou du
// boss entre `before` et `after`, en redessinant l'écran de tour par petits
// pas espacés d'un court délai, pour qu'on voie vraiment les PV baisser (ou
// remonter) au lieu de sauter directement à la valeur finale.
func animateHPChange(b *boss.Boss, player *character.Character, isPlayer bool, before, after int, message string) {
	if before == after {
		RenderTurnMessage(b, player, message)
		return
	}

	const steps = 8
	const frameDelay = 60 * time.Millisecond

	for i := 1; i <= steps; i++ {
		cur := before + (after-before)*i/steps
		if isPlayer {
			player.LP = cur
		} else {
			b.LP = cur
		}
		RenderTurnMessage(b, player, message)
		time.Sleep(frameDelay)
	}

	if isPlayer {
		player.LP = after
	} else {
		b.LP = after
	}
	RenderTurnMessage(b, player, message)
}

// waitContinue attend une validation du joueur pour laisser le temps de
// lire le message affiché par RenderTurnMessage. Renvoie false si le
// joueur quitte (X) pendant cette pause, pour que RunBattle propage la
// sortie du combat.
func waitContinue(ts *terminalSession) bool {
	return ts.readKey() != KeyQuit
}

// tryResurrect annule une mort du joueur s'il porte le charme de
// résurrection posé par la Flèche (voir combat/arrow.go, King Crimson) : il
// revient avec un tiers de ses PV max, le charme consommé. Renvoie false
// (rien à faire) si le joueur est encore en vie ou ne porte pas le charme.
func tryResurrect(player *character.Character) bool {
	if player.IsAlive() || !player.HasResurrectCharm {
		return false
	}
	player.HasResurrectCharm = false
	player.LP = player.MaxLP / 3
	if player.LP < 1 {
		player.LP = 1
	}
	return true
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
		if player.StunnedTurns > 0 {
			// Enraciné (voir combat/arrow.go, Hermit Purple) : le tour est
			// perdu sans passer par le menu, mais le boss agit quand même.
			player.StunnedTurns--
			message = "Tu es enraciné(e), impossible d'agir ce tour-ci !"
		} else {
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
				idx, ok := chooseOption(ts, b, fmt.Sprintf("%s - choisis ton attaque", player.Name), fightMenuItems(player))
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
				before := b.LP
				message = move.Perform(player, b.Character)
				animateHPChange(b, player, false, before, b.LP, message)
			case 1: // ACT : choix de l'action "lore" (parler, observer...).
				options := b.ActOptions(player)
				idx, ok := chooseOption(ts, b, fmt.Sprintf("%s - que fais-tu ?", b.Name), actMenuItems(options))
				if !ok {
					acted = false
					break
				}
				message = options[idx].Resolve(b, player)
			case 2: // ITEM : choix de l'objet à utiliser dans l'inventaire.
				if len(player.Inventory) == 0 {
					message = "Ton inventaire est vide !"
					acted = false
					break
				}
				items, names := itemMenuItems(player)
				idx, ok := chooseOption(ts, b, fmt.Sprintf("%s - choisis un objet", player.Name), items)
				if !ok {
					acted = false
					break
				}
				used, msg := useItem(b, player, names[idx])
				message = msg
				if !used {
					acted = false
				}
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

		if !player.IsAlive() {
			// Un objet utilisé au tour du joueur (ex: Potion de poison) peut
			// l'achever avant même le tour du boss.
			if tryResurrect(player) {
				RenderTurnMessage(b, player, fmt.Sprintf("Le futur où tu meurs a été effacé ! Tu reviens avec %d PV.", player.LP))
				waitContinue(ts)
			} else {
				RenderTurnMessage(b, player, "Tu es tombé(e) au combat...")
				waitContinue(ts)
				return ResultDefeat
			}
		}

		// Le tour du boss, affiché à part : c'est là que chaque boss aura
		// son propre comportement (voir boss/patterns.go).
		RenderTurnMessage(b, player, fmt.Sprintf("-- Tour de %s --", b.Name))
		if !waitContinue(ts) {
			return ResultQuit
		}
		beforePlayerLP := player.LP
		var bossMessage string
		if player.SkipBossNextTurn {
			player.SkipBossNextTurn = false
			bossMessage = fmt.Sprintf("Le temps reste figé un instant de plus : %s ne peut pas agir !", b.Name)
		} else {
			bossMessage = b.Turn(turn, player)
			turn++
		}
		animateHPChange(b, player, true, beforePlayerLP, player.LP, bossMessage)
		if !waitContinue(ts) {
			return ResultQuit
		}

		if !player.IsAlive() {
			if tryResurrect(player) {
				RenderTurnMessage(b, player, fmt.Sprintf("Le futur où tu meurs a été effacé ! Tu reviens avec %d PV.", player.LP))
				waitContinue(ts)
			} else {
				RenderTurnMessage(b, player, "Tu es tombé(e) au combat...")
				waitContinue(ts)
				return ResultDefeat
			}
		}

		message = ""
	}
	return ResultDefeat
}
