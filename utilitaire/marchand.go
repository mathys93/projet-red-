package utilitaire

import (
	"fmt"

	"ProjetRED/character"
)

var MarchandObjets = map[string]int{
	"Potion de vie":   4,
	"Potion de MP":    4,
	"Disque de Pucci": 6,
	"Tel Diavolo":     3,
	"Arrow":           15,
}

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
