package utilitaire

import (
	"fmt"

	"ProjetRED/character"
)

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

func TakeHealingPotion(c *character.Character) {
	if !c.RemoveItem("Potion de vie") {
		fmt.Println("Vous n'avez pas de potion de vie !")
		return
	}
	c.Heal(50)
	fmt.Printf("Vous utilisez Potion de vie. PV : %d / %d\n", c.LP, c.MaxLP)
}

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
	case "Tel Diavolo":
		ApplyPoisonPotion(c)
	default:
		fmt.Printf("Vous utilisez %s.\n", item)
	}
}
