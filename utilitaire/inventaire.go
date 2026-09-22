package utilitaire

import (
	"fmt"

	"ProjetRED/character"
)

// PrintInventory affiche le contenu de l'inventaire d'un personnage.
func PrintInventory(c *character.Character) {
	if len(c.Inventory) == 0 {
		fmt.Println("Inventaire vide.")
		return
	}
	fmt.Println("\n--- Inventaire ---")
	for i, item := range c.Inventory {
		fmt.Printf("%d. %s\n", i+1, item)
	}
}

// TakeHealingPotion consomme une "Potion de vie" de l'inventaire (si
// présente) et restaure 50 PV.
func TakeHealingPotion(c *character.Character) {
	if !c.RemoveItem("Potion de vie") {
		fmt.Println("Vous n'avez pas de potion de vie !")
		return
	}
	c.Heal(50)
	fmt.Printf("Vous utilisez Potion de vie. PV : %d / %d\n", c.LP, c.MaxLP)
}

// UseInventory ouvre un menu texte (hors combat) pour choisir un objet à
// utiliser dans l'inventaire, et déclenche son effet.
//
// Note : cette fonction utilise fmt.Scanln (son propre buffer sur
// os.Stdin), séparé de celui du package combat (stdinScanner). Ça ne pose
// aucun problème tant qu'elle n'est pas appelée entre deux lectures du
// package combat (carte / combat) dans la même exécution. Si tu l'intègres
// un jour dans la boucle de jeu principale (entre deux zones, par
// exemple), fais-la lire au clavier via combat.WaitEnter()/le mécanisme du
// package combat plutôt que fmt.Scanln, pour la même raison que documentée
// dans combat/input.go (deux buffers différents sur le même os.Stdin =
// entrées perdues).
func UseInventory(c *character.Character) {
	if len(c.Inventory) == 0 {
		fmt.Println("Votre inventaire est vide.")
		return
	}

	PrintInventory(c)
	fmt.Print("Choisissez le numéro de l'objet à utiliser (0 pour quitter) : ")
	var choice int
	fmt.Scanln(&choice)

	if choice == 0 {
		fmt.Println("Vous fermez l'inventaire.")
		return
	}
	if choice < 0 || choice > len(c.Inventory) {
		fmt.Println("Choix invalide.")
		return
	}

	item := c.Inventory[choice-1]
	c.RemoveItem(item)

	switch item {
	case "Potion de vie":
		c.Heal(50)
		fmt.Printf("Vous utilisez Potion de vie. PV : %d / %d\n", c.LP, c.MaxLP)
	case "Potion de poison":
		ApplyPoisonPotion(c)
	default:
		fmt.Printf("Vous utilisez %s.\n", item)
	}
}
