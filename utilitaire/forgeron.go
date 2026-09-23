package utilitaire

import (
	"fmt"

	"ProjetRED/character"
)

var ForgeronObjets = map[string]int{
	"Chapeau de l'aventurier": 5,
	"Tunique de l'aventurier": 5,
	"Bottes de l'aventurier":  5,
}

func PeutFabriquer(c *character.Character, item string) bool {
	prix, existe := ForgeronObjets[item]
	if !existe {
		fmt.Println("Le forgeron ne fabrique pas cet objet.")
		return false
	}

	if c.Fragment < prix {
		fmt.Printf("Tu n'as pas assez de Fragments pour fabriquer %s.\n", item)
		fmt.Printf("Il faut %d Fragments, tu n'en as que %d.\n", prix, c.Fragment)
		return false
	}


	if len(c.Inventory) >= 10 {
		fmt.Println("Ton inventaire est plein, tu ne peux rien fabriquer de plus.")
		return false
	}

	return true
}

func Fabriquer(c *character.Character, item string) {
	if !PeutFabriquer(c, item) {
		fmt.Println("Fabrication impossible.")
		return
	}

	c.Fragment -= ForgeronObjets[item]
	c.AddItem(item)
	fmt.Printf("Tu as fabriqué %s !\n", item)
	fmt.Printf("Fragments restants : %d\n", c.Fragment)
	fmt.Printf("Inventaire du joueur : %v\n", c.Inventory)
}
