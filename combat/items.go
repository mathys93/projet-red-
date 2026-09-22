package combat

import (
	"fmt"
	"time"

	"ProjetRED/boss"
	"ProjetRED/character"
)

// itemDescription renvoie une courte description affichée dans le sous-menu
// ITEM pour un objet connu (vide si l'objet n'a pas de description dédiée).
func itemDescription(item string) string {
	switch item {
	case "Potion de vie", "Disque de Pucci":
		return "Restaure 50 PV."
	case "Potion de poison":
		return "Inflige des dégâts, ignore la Defense Power."
	case "Arrow":
		return "Pari risqué : réveille un Stand aléatoire, bénéfique ou non."
	}
	return ""
}

// itemMenuItems construit les entrées du sous-menu ITEM à partir de
// l'inventaire du joueur, en regroupant les exemplaires identiques
// ("Potion de vie x2") pour ne pas afficher une ligne par exemplaire.
// Le deuxième slice renvoyé donne, pour chaque entrée du menu, le nom de
// l'objet correspondant (même index).
func itemMenuItems(player *character.Character) ([]menuItem, []string) {
	var names []string
	counts := map[string]int{}
	for _, it := range player.Inventory {
		if counts[it] == 0 {
			names = append(names, it)
		}
		counts[it]++
	}

	items := make([]menuItem, len(names))
	for i, name := range names {
		label := name
		if counts[name] > 1 {
			label = fmt.Sprintf("%s x%d", name, counts[name])
		}
		items[i] = menuItem{Label: label, Description: itemDescription(name)}
	}
	return items, names
}

// useItem consomme un objet de l'inventaire pendant un combat et applique
// son effet, en animant la barre de PV concernée (voir animateHPChange) au
// lieu de sauter directement au résultat final. Renvoie false (sans
// consommer l'objet ni le tour) si l'objet ne peut pas être utilisé ici.
func useItem(b *boss.Boss, player *character.Character, item string) (bool, string) {
	switch item {
	case "Potion de vie", "Disque de Pucci":
		player.RemoveItem(item)
		before := player.LP
		player.Heal(50)
		msg := fmt.Sprintf("Tu utilises %s. PV restaurés.", item)
		animateHPChange(b, player, true, before, player.LP, msg)
		return true, msg

	case "Potion de poison":
		player.RemoveItem(item)
		RenderTurnMessage(b, player, fmt.Sprintf("Tu bois %s... mauvaise idée.", item))
		time.Sleep(400 * time.Millisecond)
		for i := 1; i <= 3; i++ {
			before := player.LP
			player.LP -= 10
			if player.LP < 0 {
				player.LP = 0
			}
			tickMsg := fmt.Sprintf("Le poison te ronge... (%d/3)", i)
			animateHPChange(b, player, true, before, player.LP, tickMsg)
			time.Sleep(300 * time.Millisecond)
		}
		return true, "Le poison se dissipe."

	case "Arrow":
		player.RemoveItem(item)
		return useArrow(player)
	}

	return false, fmt.Sprintf("%s ne peut pas être utilisé maintenant.", item)
}
