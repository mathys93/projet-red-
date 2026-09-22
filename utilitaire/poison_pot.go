<<<<<<< HEAD
// Fichier poison_pot.go : effet de la Potion de poison + achat chez le
// marchand.
//
// Ce fichier remplace l'ancien "utilitaire/poison POT" : ce nom de fichier
// (avec une espace, et sans extension .go) n'était pas reconnu par le
// compilateur Go, qui l'ignorait purement et simplement — tout son code
// était donc mort, jamais compilé ni exécuté. Il définissait aussi son
// propre type Character (encore un !) et son propre func main(), ce qui
// aurait de toute façon provoqué des conflits avec le reste du paquet une
// fois renommé en .go. Tout est maintenant réécrit avec le character.Character
// partagé, sans func main() (ce n'est pas le rôle d'un paquet utilitaire).
package main
=======
package utilitaire
>>>>>>> b6405f91968060dd6b975508bd490f4401c96579

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
