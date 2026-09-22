package utilitaire

import (
	"fmt"

	"ProjetRED/character"
)

// MarchandObjets liste les objets vendus par le marchand et leur coût en
// Fragments.
var MarchandObjets = map[string]int{
	"Potion de vie":    4,
	"Disque de Pucci":  6,
	"Potion de poison": 3,
	"Arrow":            15,
}

// PeutAcheter vérifie que l'objet existe chez le marchand et que le joueur a
// assez de Fragments pour l'acheter.
func PeutAcheter(c *character.Character, item string) bool {
	prix, existe := MarchandObjets[item]
	if !existe {
		fmt.Println("Le marchand ne vend pas cet objet.")
		return false
	}

	if c.Fragment < prix {
		fmt.Printf("Tu n'as pas assez de Fragments pour acheter %s.\n", item)
		fmt.Printf("Il faut %d Fragments, tu n'en as que %d.\n", prix, c.Fragment)
		return false
	}

	return true
}

// Acheter fait acheter un objet au joueur chez le marchand : consomme les
// Fragments puis ajoute l'objet à l'inventaire.
func Acheter(c *character.Character, item string) {
	if !PeutAcheter(c, item) {
		fmt.Println("Achat impossible.")
		return
	}

	c.Fragment -= MarchandObjets[item]
	c.AddItem(item)
	fmt.Printf("Tu as acheté %s !\n", item)
	fmt.Printf("Fragments restants : %d\n", c.Fragment)
}
