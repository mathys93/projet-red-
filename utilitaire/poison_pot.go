package utilitaire

import (
	"fmt"
	"time"

	"ProjetRED/character"
)

// ApplyPoisonPotion inflige 10 dégâts par seconde pendant 3 secondes,
// directement sur les PV (le poison ignore la Defense Power).
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

// BuyPoisonPotion simule l'achat d'une Potion de poison chez le marchand.
func BuyPoisonPotion(c *character.Character) {
	c.AddItem("Potion de poison")
	fmt.Println("Vous avez acheté une Potion de poison et l'avez ajoutée à votre inventaire.")
}
