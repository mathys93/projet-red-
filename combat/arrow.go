package combat

import "ProjetRED/character"

func useArrow(player *character.Character) (bool, string) {
	return true, player.EquipStand(character.RollStand())
}
