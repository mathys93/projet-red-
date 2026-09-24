package main

import (
	"fmt"

	"ProjetRED/character"
	"ProjetRED/combat"
	"ProjetRED/world"
)

func main() {
	combat.MaximizeConsoleWindow()

	player := character.New("Personnage 1", 1, 8, 3, 50)
	player.AddItem("Potion de vie")
	player.AddItem("Potion de vie")

	w := world.New()
	for _, z := range w.Zones {
		if z.Boss != nil {
			z.Boss.Zone = z.Name
		}
	}
	w.BonusBoss.Zone = "Zone Bonus - Le Disque de Pucci"

	fmt.Println("Bienvenue,", player.Name, "!")
	fmt.Println("Appuie sur Entrée pour ouvrir la carte...")
	combat.WaitEnter()

	for {
		zone := combat.SelectZone(w)
		if zone == nil {
			fmt.Println("\nÀ bientôt !")
			return
		}

		if zone.IsShop {
			combat.RunShop(player)
			continue
		}

		if zone.IsFarm {
			zone.Boss.LP = zone.Boss.MaxLP
		}

		result := combat.RunBattle(player, zone.Boss)
		switch result {
		case combat.ResultVictory, combat.ResultSpared:
			if !zone.IsFarm {
				zone.Cleared = true
			}
			gain := zone.Boss.FragmentReward
			player.Fragment += gain
			fmt.Printf("\nTu gagnes %d Fragments ! (Total : %d)\n", gain, player.Fragment)
			for _, msg := range player.GainXP(zone.Boss.XPReward) {
				fmt.Println(msg)
			}
			fmt.Println("Appuie sur Entrée pour continuer...")
			combat.WaitEnter()
		case combat.ResultDefeat:
			fmt.Println("\nGame Over.")
			return
		case combat.ResultQuit:
			fmt.Println("\nÀ bientôt !")
			return
		}

		if w.AllCleared() {
			if combat.ConfirmBonusBoss(w.BonusBoss.Name) {
				result := combat.RunBattle(player, w.BonusBoss)
				switch result {
				case combat.ResultVictory, combat.ResultSpared:
					for _, msg := range player.GainXP(w.BonusBoss.XPReward) {
						fmt.Println(msg)
					}
					fmt.Println("\nFélicitations, tu as terminé le jeu !")
					fmt.Println("Appuie sur Entrée pour continuer...")
					combat.WaitEnter()
					combat.ShowCredits()
				case combat.ResultDefeat:
					fmt.Println("\nGame Over.")
				}
				return
			}
		}
	}
}
