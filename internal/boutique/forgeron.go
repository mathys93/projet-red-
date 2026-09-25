package boutique

import (
	"fmt"

	"ProjetRED/internal/character"
	"ProjetRED/internal/ui"
)

var ForgeronObjets = map[string]int{
	"Chapeau de l'aventurier": 5,
	"Tunique de l'aventurier": 5,
	"Bottes de l'aventurier":  5,
}

func PeutFabriquer(c *character.Character, item string) bool {
	prix, existe := ForgeronObjets[item]
	if !existe {
		fmt.Fprintln(ui.Out, "Le forgeron ne fabrique pas cet objet.")
		return false
	}

	if c.Fragment < prix {
		fmt.Fprintf(ui.Out, "Tu n'as pas assez de Fragments pour fabriquer %s.\n", item)
		fmt.Fprintf(ui.Out, "Il faut %d Fragments, tu n'en as que %d.\n", prix, c.Fragment)
		return false
	}

	if len(c.Inventory) >= 10 {
		fmt.Fprintln(ui.Out, "Ton inventaire est plein, tu ne peux rien fabriquer de plus.")
		return false
	}

	return true
}

func Fabriquer(c *character.Character, item string) {
	if !PeutFabriquer(c, item) {
		fmt.Fprintln(ui.Out, "Fabrication impossible.")
		return
	}

	c.Fragment -= ForgeronObjets[item]
	c.AddItem(item)
	fmt.Fprintf(ui.Out, "Tu as fabriqué %s !\n", item)
	fmt.Fprintf(ui.Out, "Fragments restants : %d\n", c.Fragment)
	fmt.Fprintf(ui.Out, "Inventaire du joueur : %v\n", c.Inventory)
}
