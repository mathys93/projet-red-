package combat

import (
	"fmt"
	"time"

	"ProjetRED/boss"
	"ProjetRED/character"
)

func itemDescription(item string) string {
	switch item {
	case "Potion de vie":
		return "Restaure 20 PV."
	case "Potion de MP":
		return "Restaure 15 MP."
	case "Disque de Pucci":
		return "Restaure 50 PV."
	case "Potion de poison":
		return "Inflige des dégâts, ignore la Defense Power."
	case "Arrow":
		return "Éveille un Stand permanent : bonus d'attaque et nouvelle action de combat."
	}
	return ""
}

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

func useItem(b *boss.Boss, player *character.Character, item string) (bool, string) {
	switch item {
	case "Potion de vie", "Disque de Pucci":
		soin := 20
		if item == "Disque de Pucci" {
			soin = 50
		}
		player.RemoveItem(item)
		before := player.LP
		player.Heal(soin)
		msg := fmt.Sprintf("Tu utilises %s. +%d PV.", item, player.LP-before)
		animateHPChange(b, player, true, before, player.LP, msg)
		return true, msg

	case "Potion de MP":
		player.RemoveItem(item)
		before := player.MP
		player.MP += 15
		if player.MP > player.MaxMP {
			player.MP = player.MaxMP
		}
		return true, fmt.Sprintf("Tu utilises %s. +%d MP (%d/%d).", item, player.MP-before, player.MP, player.MaxMP)

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
