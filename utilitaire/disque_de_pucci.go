package utilitaire

import (
	"fmt"

	"ProjetRED/character"
)

func ApplyDisqueDePucci(c *character.Character) {
	if !c.RemoveItem("Disque de Pucci") {
		fmt.Println("Vous n'avez pas de Disque de Pucci !")
		return
	}
	c.Heal(50)
	fmt.Printf("Vous utilisez Disque de Pucci. PV : %d / %d\n", c.LP, c.MaxLP)
}

func BuyDisqueDePucci(c *character.Character) {
	c.AddItem("Disque de Pucci")
	fmt.Println("Vous avez acheté un Disque de Pucci et l'avez ajouté à votre inventaire.")
}
