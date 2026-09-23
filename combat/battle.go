package combat

import (
	"fmt"
	"strings"
	"time"

	"ProjetRED/ascii"
	"ProjetRED/boss"
	"ProjetRED/character"
)

const (
	colReset  = "\033[0m"
	colWhite  = "\033[97m"
	colYellow = "\033[93m"
	colRed    = "\033[91m"
)

func clearScreen() {
	fmt.Print("\033[2J\033[H\033[?25l")
}

func MaximizeConsoleWindow() {
	maximizeConsoleWindow()
}

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

const (
	minBoxWidth   = 66
	minBoxHeight  = 8
	minArtHeight  = 6
	battleChrome  = 12
	tallTerminal  = 46
	tallBoxHeight = 10
)

func layout() (boxW, boxH, artW, artH int) {
	cols, rows := terminalSize()

	boxW = cols - 4
	if boxW < minBoxWidth {
		boxW = minBoxWidth
	}

	boxH = minBoxHeight
	if rows >= tallTerminal {
		boxH = tallBoxHeight
	}

	artH = rows - 1 - battleChrome - boxH
	for artH < minArtHeight && boxH > 5 {
		boxH--
		artH++
	}
	if artH < minArtHeight {
		artH = minArtHeight
	}

	return boxW, boxH, boxW, artH
}

func fitArt(b *boss.Boss, artW, artH int) string {
	return ascii.Fit(b.Art, artW, artH)
}

func printPadding(contentHeight int) {
	_, rows := terminalSize()
	pad := (rows - contentHeight) / 2
	for i := 0; i < pad; i++ {
		fmt.Println()
	}
}

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

func drawFrameTop(width int) {
	fmt.Println(colWhite + "┌" + strings.Repeat("─", width-2) + "┐" + colReset)
}

func drawFrameBottom(width int) {
	fmt.Println(colWhite + "└" + strings.Repeat("─", width-2) + "┘" + colReset)
}

func drawFrameLine(width int, content string) {
	pad := width - 2 - visibleWidth(content)
	if pad < 0 {
		pad = 0
	}
	fmt.Println(colWhite + "│" + colReset + content + strings.Repeat(" ", pad) + colWhite + "│" + colReset)
}

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

func DrawEmptyBox(width, height int) {
	drawFrameTop(width)
	for i := 0; i < height-2; i++ {
		drawFrameLine(width, "")
	}
	drawFrameBottom(width)
}

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
	inner += 4

	btnOuter := inner + 2
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

	stand := ""
	if c.Stand != nil {
		stand = fmt.Sprintf("   %s[%s]%s", colYellow, c.Stand.Name, colReset)
	}
	fmt.Printf(" %-14s LV %-3d HP %s %d/%d   MP %d/%d   XP %d/%d%s\n",
		c.Name, c.Level, bar, c.LP, c.MaxLP, c.MP, c.MaxMP, c.XP, character.XPForLevel(c.Level), stand)
}

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

func RenderBattleScreen(b *boss.Boss, player *character.Character, selected int, message string) {
	clearScreen()
	bw, bh, aw, ah := layout()
	art := fitArt(b, aw, ah)
	artHeight := len(strings.Split(art, "\n"))
	contentHeight := artHeight + bh + 11
	if message != "" {
		contentHeight++
	}
	printPadding(contentHeight)
	fmt.Println()
	fmt.Println(colYellow + "  " + b.Zone + colReset)
	DrawArt(bw, art, b.Color)
	DrawEmptyBox(bw, bh)
	fmt.Println()
	DrawEnemyBar(b)
	if message != "" {
		fmt.Println(colWhite + "* " + message + colReset)
	}
	DrawStatBar(player)
	DrawActionButtons(bw, selected)
	fmt.Println(colWhite + "\n(← → pour choisir, Entrée pour valider, F/A/I/M en raccourci, X pour quitter)" + colReset)
}

func RenderTurnMessage(b *boss.Boss, player *character.Character, message string) {
	clearScreen()
	bw, bh, aw, ah := layout()
	art := fitArt(b, aw, ah)
	artHeight := len(strings.Split(art, "\n"))
	contentHeight := artHeight + bh + 7
	if message != "" {
		contentHeight++
	}
	printPadding(contentHeight)
	fmt.Println()
	fmt.Println(colYellow + "  " + b.Zone + colReset)
	DrawArt(bw, art, b.Color)
	DrawEmptyBox(bw, bh)
	fmt.Println()
	DrawEnemyBar(b)
	if message != "" {
		fmt.Println(colWhite + "* " + message + colReset)
	}
	DrawStatBar(player)
	fmt.Println(colWhite + "\n(Entrée pour continuer, X pour quitter)" + colReset)
}

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

func waitContinue(ts *terminalSession) bool {
	return ts.readKey() != KeyQuit
}

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

type Result int

const (
	ResultVictory Result = iota
	ResultDefeat
	ResultSpared
	ResultQuit
)

var mainMenuHotkeys = map[string]int{
	"f": 0, "1": 0,
	"a": 1, "2": 1,
	"i": 2, "3": 2,
	"m": 3, "4": 3,
}

func RunBattle(player *character.Character, b *boss.Boss) Result {
	ts := newTerminalSession()
	defer ts.restore()

	selected := 0
	turn := 1
	message := fmt.Sprintf("%s bloque le passage !", b.Name)

	for player.IsAlive() && b.IsAlive() {
		if player.StunnedTurns > 0 {
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
			case 0:
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
			case 1:
				options := b.ActOptions(player)
				idx, ok := chooseOption(ts, b, fmt.Sprintf("%s - que fais-tu ?", b.Name), actMenuItems(options))
				if !ok {
					acted = false
					break
				}
				message = options[idx].Resolve(b, player)
			case 2:
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
			case 3:
				message = fmt.Sprintf("Tu épargnes %s...", b.Name)
				RenderBattleScreen(b, player, selected, message)
				return ResultSpared
			}

			if !acted {
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
			if tryResurrect(player) {
				RenderTurnMessage(b, player, fmt.Sprintf("Le futur où tu meurs a été effacé ! Tu reviens avec %d PV.", player.LP))
				waitContinue(ts)
			} else {
				RenderTurnMessage(b, player, "Tu es tombé(e) au combat...")
				waitContinue(ts)
				return ResultDefeat
			}
		}

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
			struck, quit := RunDodge(b, player)
			if quit {
				return ResultQuit
			}
			if struck {
				bossMessage = b.Turn(turn, player)
			} else {
				bossMessage = fmt.Sprintf("Tu esquives l'attaque de %s sans une égratignure !", b.Name)
			}
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
