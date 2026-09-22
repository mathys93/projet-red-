package main

import "fmt"

// forgeronObjets liste les équipements fabricables et leur coût en Fragments
var forgeronObjets = map[string]int{
	"Chapeau de l'aventurier": 5,
	"Tunique de l'aventurier": 5,
	"Bottes de l'aventurier":  5,
}

// forgeronRessources liste les matériaux nécessaires à chaque fabrication
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

// compteItem compte le nombre d'occurrences d'un item dans l'inventaire
func compteItem(c *Character, item string) int {
	count := 0
	for _, it := range c.Inventory {
		if it == item {
			count++
		}
	}
	return count
}

// possedeRessources vérifie que le joueur a tous les matériaux nécessaires
func possedeRessources(c *Character, item string) bool {
	for ressource, quantite := range forgeronRessources[item] {
		if compteItem(c, ressource) < quantite {
			fmt.Printf("Il te manque des ressources : %s (%d requis, tu en as %d).\n",
				ressource, quantite, compteItem(c, ressource))
			return false
		}
	}
	return true
}

// peutFabriquer vérifie que l'objet existe, que le joueur a assez de Fragments,
// les ressources nécessaires, et de la place dans l'inventaire.
func peutFabriquer(c *Character, item string) bool {
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

// fabriquer gère la fabrication d'un équipement chez le forgeron :
// consomme les Fragments, retire les ressources, puis ajoute l'objet fabriqué.
func fabriquer(c *Character, item string) {
	if !peutFabriquer(c, item) {
		fmt.Println("Fabrication impossible.")
		return
	}

	// on retire les ressources consommées
	for ressource, quantite := range forgeronRessources[item] {
		for i := 0; i < quantite; i++ {
			removeInventory(c, ressource)
		}
	}

	c.Fragment -= forgeronObjets[item]
	addInventory(c, item)

	fmt.Printf("Tu as fabriqué %s !\n", item)
	fmt.Printf("Fragments restants : %d\n", c.Fragment)
	fmt.Printf("Inventaire du joueur : %v\n", c.Inventory)
}