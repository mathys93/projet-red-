package utilitaire

import (
	"fmt"
	"time"

	"ProjetRED/character"
)

func ApplyPoisonPotion(c *character.Character) {
	fmt.Println("--- Effet de la potion de poison ---")
	for i := 1; i <= 3; i++ {
		time.Sleep(1 * time.Second)
		c.LP -= 10
		if c.LP < 0 {
			c.LP = 0
		}
		fmt.Printf("Seconde %d : PV %d / %d\n", i, c.LP, c.MaxLP)
	}
}

func BuyPoisonPotion(c *character.Character) {
	c.AddItem("Tel Diavolo")
	fmt.Println("Vous avez acheté un Tel Diavolo et l'avez ajouté à votre inventaire.")
}
