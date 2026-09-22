package main

import (
	"fmt"
	"math/rand"
)

var inventory []string

func accessInventory() {
	if len(inventory) == 0 {
		fmt.Println("Inventaire vide.")
		return
	}

	fmt.Println("Inventaire du joueur :")
	for _, item := range inventory {
		fmt.Printf("- %s\n", item)
	}
}

func addInventory(item string) {
	inventory = append(inventory, item)
	fmt.Printf("%s ajouté à l'inventaire.\n", item)
}

func removeInventory(item string) {
	for i, current := range inventory {
		if current == item {
			inventory = append(inventory[:i], inventory[i+1:]...)
			fmt.Printf("%s retiré de l'inventaire.\n", item)
			return
		}
	}
	fmt.Printf("%s n'est pas dans l'inventaire.\n", item)
}

func marchand() {
	var choix string
	fmt.Println("Bienvenue chez le marchand !")
	fmt.Println("Il vend gratuitement : Potion de vie")
	fmt.Print("Voulez-vous l'ajouter à votre inventaire ? (oui/non) : ")
	fmt.Scanln(&choix)
	if choix == "oui" {
		addInventory("Potion de vie")
	} else {
		fmt.Println("Vous quittez le marchand.")
	}
}

func genererObjetAleatoire() {
	objets := []struct {
		nom, effet string
	}{
		{"Potion de vie", "+20 HP"},
		{"Potion de poison", "-15 HP"},
	}
	objet := objets[rand.Intn(len(objets))]
	fmt.Printf("Objet obtenu : %s\nEffet : %s\n", objet.nom, objet.effet)
	addInventory(objet.nom)
}

func menu() {
	for {
		var choix int
		fmt.Println("\n--- MENU ---")
		fmt.Println("1. Marchand")
		fmt.Println("2. Inventaire")
		fmt.Println("3. Objet aléatoire")
		fmt.Println("4. Quitter")
		fmt.Print("Choisissez une option : ")
		if _, err := fmt.Scanln(&choix); err != nil {
			fmt.Println("Choix invalide.")
			continue
		}

		switch choix {
		case 1:
			marchand()
		case 2:
			accessInventory()
		case 3:
			genererObjetAleatoire()
		case 4:
			fmt.Println("Fin du programme.")
			return
		default:
			fmt.Println("Choix invalide.")
		}
	}
}

func main() {
	menu()
}
