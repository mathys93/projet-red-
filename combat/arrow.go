package combat

import (
	"fmt"
	"math/rand"

	"ProjetRED/character"
)

type standRoll struct {
	Name  string
	Apply func(c *character.Character) string
}

var rouletteStand = []standRoll{
	{
		Name: "Crazy Diamond : soin complet",
		Apply: func(c *character.Character) string {
			c.LP = c.MaxLP
			return "Tes blessures se referment instantanément. PV au maximum !"
		},
	},
	{
		Name: "Gold Experience : vitalité renforcée",
		Apply: func(c *character.Character) string {
			c.MaxLP += 20
			c.LP += 20
			if c.LP > c.MaxLP {
				c.LP = c.MaxLP
			}
			return "Une vie nouvelle circule en toi. PV maximum +20, définitivement."
		},
	},
	{
		Name: "Star Platinum : puissance à l'état pur",
		Apply: func(c *character.Character) string {
			c.AP += 4
			return "Tu sens une force nouvelle dans tes poings. Attack Power +4, définitivement."
		},
	},
	{
		Name: "The World : temps arrêté",
		Apply: func(c *character.Character) string {
			c.SkipBossNextTurn = true
			return "\"ZA WARUDO !\" Le temps s'arrête un instant : l'adversaire ne pourra pas agir à son prochain tour."
		},
	},
	{
		Name: "King Crimson : le futur est effacé",
		Apply: func(c *character.Character) string {
			c.HasResurrectCharm = true
			return "Le futur où tu meurs a été effacé. Si tu tombes à 0 PV, tu reviendras une fois."
		},
	},
	{
		Name: "Hermit Purple : racines",
		Apply: func(c *character.Character) string {
			c.StunnedTurns += 2
			return "Des ronces jaillissent du sol et t'enracinent : tu ne pourras pas agir pendant 2 tours."
		},
	},
	{
		Name: "Whitesnake : objet volé",
		Apply: func(c *character.Character) string {
			if len(c.Inventory) == 0 {
				return "Une main spectrale te fouille, mais ton inventaire est déjà vide."
			}
			perdu := c.Inventory[0]
			c.RemoveItem(perdu)
			return fmt.Sprintf("Une main spectrale te dérobe %s !", perdu)
		},
	},
	{
		Name: "Rejet du Stand",
		Apply: func(c *character.Character) string {
			degats := c.LP / 5
			if degats < 1 {
				degats = 1
			}
			c.LP -= degats
			if c.LP < 1 {
				c.LP = 1
			}
			return fmt.Sprintf("Ton corps rejette le Stand ! %d dégâts, mais tu tiens encore debout.", degats)
		},
	},
}

func useArrow(player *character.Character) (bool, string) {
	tirage := rouletteStand[rand.Intn(len(rouletteStand))]
	effet := tirage.Apply(player)
	msg := fmt.Sprintf("La Flèche te transperce... Stand éveillé : %s ! %s", tirage.Name, effet)
	return true, msg
}
