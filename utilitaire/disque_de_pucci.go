package utilitaire

import (
	"fmt"

	"ProjetRED/character"
)

// ApplyDisqueDePucci consomme un "Disque de Pucci" de l'inventaire (si
// présent) et restaure 50 PV.
func ApplyDisqueDePucci(c *character.Character) {
	if !c.RemoveItem("Disque de Pucci") {
		fmt.Println("Vous n'avez pas de Disque de Pucci !")
		return
	}
	c.Heal(50)
	fmt.Printf("Vous utilisez Disque de Pucci. PV : %d / %d\n", c.LP, c.MaxLP)
}

// BuyDisqueDePucci simule l'achat d'un Disque de Pucci chez le marchand.
func BuyDisqueDePucci(c *character.Character) {
	c.AddItem("Disque de Pucci")
	fmt.Println("Vous avez acheté un Disque de Pucci et l'avez ajouté à votre inventaire.")
}
