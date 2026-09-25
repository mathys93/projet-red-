package combat

import (
	"fmt"
	"strings"
	"time"

	"ProjetRED/internal/ascii"
	"ProjetRED/internal/audio"
	"ProjetRED/internal/boss"
	"ProjetRED/internal/character"
	"ProjetRED/internal/ui"
)

func fitArt(b *boss.Boss, artW, artH int) string {
	return ascii.Fit(b.Art, artW, artH)
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
	gap := max((width-len(options)*btnOuter)/(len(options)+1), 2)
	margin := strings.Repeat(" ", gap)

	var top, mid, bot strings.Builder
	for i, label := range labels {
		col := ui.ColWhite
		heart := "  "
		if i == selected {
			col = ui.ColRed
			heart = ui.ColRed + "❤ " + ui.ColReset
		}
		pad := max(inner-2-len([]rune(label)), 0)
		left := pad / 2
		right := pad - left

		top.WriteString(margin)
		top.WriteString(col)
		top.WriteString("┌")
		top.WriteString(strings.Repeat("─", inner))
		top.WriteString("┐")
		top.WriteString(ui.ColReset)
		mid.WriteString(margin)
		mid.WriteString(col)
		mid.WriteString("│")
		mid.WriteString(ui.ColReset)
		mid.WriteString(heart)
		mid.WriteString(strings.Repeat(" ", left))
		mid.WriteString(col)
		mid.WriteString(label)
		mid.WriteString(ui.ColReset)
		mid.WriteString(strings.Repeat(" ", right))
		mid.WriteString(col)
		mid.WriteString("│")
		mid.WriteString(ui.ColReset)
		bot.WriteString(margin)
		bot.WriteString(col)
		bot.WriteString("└")
		bot.WriteString(strings.Repeat("─", inner))
		bot.WriteString("┘")
		bot.WriteString(ui.ColReset)
	}
	top.WriteString(margin)
	mid.WriteString(margin)
	bot.WriteString(margin)

	fmt.Fprintln(ui.Out)
	fmt.Fprintln(ui.Out, top.String())
	fmt.Fprintln(ui.Out, mid.String())
	fmt.Fprintln(ui.Out, bot.String())
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
	bar := ui.ColYellow + strings.Repeat("■", filled) + ui.ColWhite + strings.Repeat("□", barLen-filled) + ui.ColReset

	stand := ""
	if c.Stand != nil {
		stand = fmt.Sprintf("   %s[%s]%s", ui.ColYellow, c.Stand.Name, ui.ColReset)
	}
	fmt.Fprintf(ui.Out, " %-14s LV %-3d HP %s %d/%d   MP %d/%d   XP %d/%d%s\n",
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
	bar := ui.ColRed + strings.Repeat("■", filled) + ui.ColWhite + strings.Repeat("□", barLen-filled) + ui.ColReset

	fmt.Fprintf(ui.Out, " %-14s LV %-3d HP %s %d/%d\n", b.Name, b.Level, bar, b.LP, b.MaxLP)
}

func RenderBattleScreen(b *boss.Boss, player *character.Character, selected int, message string) {
	ui.CurrentScreen = func() { RenderBattleScreen(b, player, selected, message) }
	ui.ClearScreen()
	bw, bh, aw, ah := ui.Layout()
	art := fitArt(b, aw, ah)
	artHeight := len(strings.Split(art, "\n"))
	contentHeight := artHeight + bh + 11
	if message != "" {
		contentHeight++
	}
	ui.PrintPadding(contentHeight)
	fmt.Fprintln(ui.Out)
	fmt.Fprintln(ui.Out, ui.ColYellow+"  "+b.Zone+ui.ColReset)
	ui.DrawArt(bw, art, b.Color)
	ui.DrawEmptyBox(bw, bh)
	fmt.Fprintln(ui.Out)
	DrawEnemyBar(b)
	if message != "" {
		fmt.Fprintln(ui.Out, ui.ColWhite+"* "+message+ui.ColReset)
	}
	DrawStatBar(player)
	DrawActionButtons(bw, selected)
	fmt.Fprintln(ui.Out, ui.ColWhite+"\n(← → pour choisir, Entrée pour valider, F/A/I/M en raccourci, Échap pour le menu pause)"+ui.ColReset)
}

func RenderTurnMessage(b *boss.Boss, player *character.Character, message string) {
	ui.CurrentScreen = func() { RenderTurnMessage(b, player, message) }
	ui.ClearScreen()
	bw, bh, aw, ah := ui.Layout()
	art := fitArt(b, aw, ah)
	artHeight := len(strings.Split(art, "\n"))
	contentHeight := artHeight + bh + 7
	if message != "" {
		contentHeight++
	}
	ui.PrintPadding(contentHeight)
	fmt.Fprintln(ui.Out)
	fmt.Fprintln(ui.Out, ui.ColYellow+"  "+b.Zone+ui.ColReset)
	ui.DrawArt(bw, art, b.Color)
	ui.DrawEmptyBox(bw, bh)
	fmt.Fprintln(ui.Out)
	DrawEnemyBar(b)
	if message != "" {
		fmt.Fprintln(ui.Out, ui.ColWhite+"* "+message+ui.ColReset)
	}
	DrawStatBar(player)
	fmt.Fprintln(ui.Out, ui.ColWhite+"\n(Entrée pour continuer, Échap pour le menu pause)"+ui.ColReset)
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

func waitContinue(ts *ui.Session) bool {
	return ts.ReadKey() != ui.KeyQuit
}

func tryResurrect(player *character.Character) bool {
	if player.IsAlive() || !player.HasResurrectCharm {
		return false
	}
	player.HasResurrectCharm = false
	player.LP = max(player.MaxLP/3, 1)
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

func sceneFor(b *boss.Boss) string {
	switch b.Style {
	case boss.AttackTimeStop:
		return ui.SceneDio
	case boss.AttackErase:
		return ui.SceneDiavolo
	case boss.AttackBombs:
		return ui.SceneKira
	case boss.AttackAccelerate:
		return ui.ScenePucci
	}
	return ui.SceneIggy
}

func RunBattle(player *character.Character, b *boss.Boss) Result {
	ts := ui.NewSession()
	defer ts.Restore()

	audio.Play(audio.TrackFor(b.Name))
	defer audio.Stop()

	prevScene := ui.SetScene(sceneFor(b))
	defer func() {
		ui.SetScene(prevScene)
		ui.ClearScreen()
	}()

	selected := 0
	turn := 1
	message := fmt.Sprintf("%s bloque le passage !", b.Name)

	for player.IsAlive() && b.IsAlive() {
		RenderBattleScreen(b, player, selected, message)

		key := ts.ReadKey()
		chosen := -1
		switch key {
		case ui.KeyLeft, ui.KeyUp:
			selected = (selected + 3) % 4
		case ui.KeyRight, ui.KeyDown:
			selected = (selected + 1) % 4
		case ui.KeyQuit:
			return ResultQuit
		case ui.KeyEnter:
			chosen = selected
		case ui.KeyOther:
			if idx, ok := mainMenuHotkeys[ts.Raw]; ok {
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
