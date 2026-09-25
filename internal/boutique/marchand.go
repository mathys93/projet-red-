package boutique

import (
	"fmt"

	"ProjetRED/internal/character"
	"ProjetRED/internal/ui"
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
		fmt.Fprintln(ui.Out, "Le marchand ne vend pas cet objet.")
		return false
	}

	if c.Fragment < prix {
		fmt.Fprintf(ui.Out, "Tu n'as pas assez de Fragments pour acheter %s.\n", item)
		fmt.Fprintf(ui.Out, "Il faut %d Fragments, tu n'en as que %d.\n", prix, c.Fragment)
		return false
	}

	return true
}

func Acheter(c *character.Character, item string) {
	if !PeutAcheter(c, item) {
		fmt.Fprintln(ui.Out, "Achat impossible.")
		return
	}

	c.Fragment -= MarchandObjets[item]
	c.AddItem(item)
	fmt.Fprintf(ui.Out, "Tu as acheté %s !\n", item)
	fmt.Fprintf(ui.Out, "Fragments restants : %d\n", c.Fragment)
}
