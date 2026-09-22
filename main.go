// Point d'entrée du jeu.
//
// Tout le code vivait auparavant dans ce seul fichier (structures, entrée
// clavier, rendu de combat, boucle de jeu...). Il est maintenant réparti en
// plusieurs packages :
//
//	character/  -> le personnage (joueur ET boss) : stats, PV, inventaire
//	ascii/      -> chargement + redimensionnement des artworks ASCII des boss
//	boss/       -> les boss (artwork + pattern de combat, à personnaliser)
//	world/      -> la carte (3 zones de boss + 1 zone boutique + un boss bonus)
//	combat/     -> l'interface de combat façon Undertale (boîte, menu, boucle)
//	utilitaire/ -> forgeron (craft), marchand, effets d'objets (potions...)
//
// main.go ne fait plus que les assembler.
package main

import (
	"fmt"

	"ProjetRED/character"
	"ProjetRED/combat"
	"ProjetRED/world"
)

func main() {
	player := character.New("Personnage 1", 1, 8, 3, 50)
	player.AddItem("Potion de vie")
	player.AddItem("Potion de vie")

	w := world.New()
	// Le nom de la zone est affiché au-dessus de la boîte de combat (la
	// zone boutique n'a pas de boss, donc pas de Zone à renseigner).
	for _, z := range w.Zones {
		if z.Boss != nil {
			z.Boss.Zone = z.Name
		}
	}
	w.BonusBoss.Zone = "Zone Bonus - Le Disque de Pucci"

	fmt.Println("Bienvenue,", player.Name, "!")
	fmt.Println("Appuie sur Entrée pour ouvrir la carte...")
	// On passe par combat.WaitEnter() (et pas fmt.Scanln()) : les deux
	// lisent os.Stdin avec leur propre buffer interne, et le premier à
	// lire risquerait d'engloutir des entrées destinées au second. Voir
	// le commentaire sur stdinScanner dans combat/input.go.
	combat.WaitEnter()

	const fragmentsParBoss = 15

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

		result := combat.RunBattle(player, zone.Boss)
		switch result {
		case combat.ResultVictory, combat.ResultSpared:
			zone.Cleared = true
			player.Fragment += fragmentsParBoss
			fmt.Printf("\nTu gagnes %d Fragments ! (Total : %d)\n", fragmentsParBoss, player.Fragment)
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
					fmt.Println("\nFélicitations, tu as terminé le jeu !")
				case combat.ResultDefeat:
					fmt.Println("\nGame Over.")
				}
				return
			}
		}
	}
}
