package utilitaire

import "fmt"

type Character struct {
	string
	joestar   string
	Name      string
	Joestar   string
	Level     int
	MaxHP     int
	CurrentHP int
	Inventory []string // une liste de noms d'items
}

func addInventory(c *Character, item string) {
	c.Inventory = append(c.Inventory, item)
}

func removeInventory(c *Character, item string) bool {
	for i, it := range c.Inventory {
		if it == item {
			// on recolle le slice sans l'élément à l'index i
			c.Inventory = append(c.Inventory[:i], c.Inventory[i+1:]...)
			return true
		}
	}
	return false // l'item n'était pas dans l'inventaire
}
func takePot(c *Character) {
	// 1. on retire la potion (et on vérifie qu'on en avait une)
	if !removeInventory(c, "Potion de vie") {
		fmt.Println("Vous n'avez pas de potion de vie !")
		return
	}

	// 2. on soigne de 50 PV
	c.CurrentHP += 50

	// 3. on plafonne aux PV max
	if c.CurrentHP > c.MaxHP {
		c.CurrentHP = c.MaxHP
	}

	// 4. on affiche PV actuels / PV max
	fmt.Printf("Vous utilisez Potion de vie. PV : %d / %d\n", c.CurrentHP, c.MaxHP)
}

func displayInventory(c *Character) {
	for i, item := range c.Inventory {
		fmt.Printf("%d. %s\n", i+1, item)
	}
	fmt.Println("Tapez le nom de l'item à utiliser (ou 0 pour retour) :")
	// lis l'entrée avec bufio.Scanner ou fmt.Scanln, puis :
	// switch choix { case "Potion de vie": takePot(c) ... }
}

func accessInventory(c *Character) {
	for i, item := range c.Inventory {
		fmt.Printf("%d. %s\n", i+1, item)
	}
	fmt.Println("Tapez le nom de l'item à utiliser (ou 0 pour retour) :")
	// lis l'entrée avec bufio.Scanner ou fmt.Scanln, puis :
	// switch choix { case "Potion de vie": takePot(c) ... }
}
