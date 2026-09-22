// Fichier forgeron.go : fabrication d'équipement contre des Fragments et
// des ressources trouvées en jeu.
package utilitaire

import (
	"fmt"

	"ProjetRED/character"
)

// forgeronObjets liste les équipements fabricables et leur coût en Fragments.
var forgeronObjets = map[string]int{
	"Chapeau de l'aventurier": 5,
	"Tunique de l'aventurier": 5,
	"Bottes de l'aventurier":  5,
}

// forgeronRessources liste les matériaux nécessaires à chaque fabrication.
var forgeronRessources = map[string]map[string]int{
	"Chapeau de l'aventurier": {
		"Plume de Corbeau": 1,
		"Cuir de Sanglier": 1,
	},
	"Tunique de l'aventurier": {
		"Fourrure de Loup": 2,
		"Peau de Troll":    1,
	},
	"Bottes de l'aventurier": {
		"Fourrure de Loup": 1,
		"Cuir de Sanglier": 1,
	},
}

// possedeRessources vérifie que le joueur a tous les matériaux nécessaires.
func possedeRessources(c *character.Character, item string) bool {
	for ressource, quantite := range forgeronRessources[item] {
		if c.CountItem(ressource) < quantite {
			fmt.Printf("Il te manque des ressources : %s (%d requis, tu en as %d).\n",
				ressource, quantite, c.CountItem(ressource))
			return false
		}
	}
	return true
}

// PeutFabriquer vérifie que l'objet existe, que le joueur a assez de
// Fragments, les ressources nécessaires, et de la place dans l'inventaire.
func PeutFabriquer(c *character.Character, item string) bool {
	prix, existe := forgeronObjets[item]
	if !existe {
		fmt.Println("Le forgeron ne fabrique pas cet objet.")
		return false
	}

	if c.Fragment < prix {
		fmt.Printf("Tu n'as pas assez de Fragments pour fabriquer %s.\n", item)
		fmt.Printf("Il faut %d Fragments, tu n'en as que %d.\n", prix, c.Fragment)
		return false
	}

	if !possedeRessources(c, item) {
		return false
	}

	if len(c.Inventory) >= 10 {
		fmt.Println("Ton inventaire est plein, tu ne peux rien fabriquer de plus.")
		return false
	}

	return true
}

// Fabriquer gère la fabrication d'un équipement chez le forgeron : consomme
// les Fragments, retire les ressources, puis ajoute l'objet fabriqué.
func Fabriquer(c *character.Character, item string) {
	if !PeutFabriquer(c, item) {
		fmt.Println("Fabrication impossible.")
		return
	}

	for ressource, quantite := range forgeronRessources[item] {
		for i := 0; i < quantite; i++ {
			c.RemoveItem(ressource)
		}
	}

	c.Fragment -= forgeronObjets[item]
	c.AddItem(item)
	fmt.Printf("Tu as fabriqué %s !\n", item)
	fmt.Printf("Fragments restants : %d\n", c.Fragment)
	fmt.Printf("Inventaire du joueur : %v\n", c.Inventory)
}
