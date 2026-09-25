package main

import (
	"fmt"

	"ProjetRED/internal/boutique"
	"ProjetRED/internal/character"
	"ProjetRED/internal/combat"
	"ProjetRED/internal/ui"
	"ProjetRED/internal/world"
)

func main() {
	ui.MaximizeConsoleWindow()
	defer ui.RestoreTerminal()

	for {
		switch ui.MainMenu() {
		case ui.MenuStart:
			name, ok := ui.PromptName()
			if !ok {
				continue
			}
			play(name)
			if ui.QuitRequested() {
				fmt.Fprintln(ui.Out, "\nÀ bientôt !")
				return
			}
		case ui.MenuSettings:
			ui.RunSettings()
		case ui.MenuCredits:
			ui.RunCredits()
		case ui.MenuQuit:
			fmt.Fprintln(ui.Out, "\nÀ bientôt !")
			return
		}
	}
}

func pause(msg string) {
	fmt.Fprintln(ui.Out, msg)
	fmt.Fprintln(ui.Out, "Appuie sur Entrée pour continuer...")
	ui.WaitEnter()
}

func play(name string) {
	ui.EnterGame()
	defer ui.LeaveGame()

	player := character.New(name, 1, 8, 3, 50)
	player.AddItem("Potion de vie")
	player.AddItem("Potion de vie")

	w := world.New()
	for _, z := range w.Zones {
		if z.Boss != nil {
			z.Boss.Zone = z.Name
		}
	}
	w.BonusBoss.Zone = "Zone Bonus - Le Disque de Pucci"

	for {
		zone := ui.SelectZone(w)
		if zone == nil {
			return
		}

		if zone.IsShop {
			boutique.RunShop(player)
			continue
		}

		if zone.IsFarm {
			zone.Boss.LP = zone.Boss.MaxLP
		}

		switch combat.RunBattle(player, zone.Boss) {
		case combat.ResultVictory, combat.ResultSpared:
			if !zone.IsFarm {
				zone.Cleared = true
			}
			gain := zone.Boss.FragmentReward
			player.Fragment += gain
			fmt.Fprintf(ui.Out, "\nTu gagnes %d Fragments ! (Total : %d)\n", gain, player.Fragment)
			for _, msg := range player.GainXP(zone.Boss.XPReward) {
				fmt.Fprintln(ui.Out, msg)
			}
			pause("")
		case combat.ResultDefeat:
			pause("\nGame Over.")
			return
		case combat.ResultQuit:
			return
		}

		if w.AllCleared() && ui.ConfirmBonusBoss(w.BonusBoss.Name) {
			switch combat.RunBattle(player, w.BonusBoss) {
			case combat.ResultVictory, combat.ResultSpared:
				for _, msg := range player.GainXP(w.BonusBoss.XPReward) {
					fmt.Fprintln(ui.Out, msg)
				}
				pause(fmt.Sprintf("\nFélicitations %s, tu as terminé le jeu !", player.Name))
			case combat.ResultDefeat:
				pause("\nGame Over.")
			}
			return
		}
	}
}
