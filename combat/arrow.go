package combat

import (
	"fmt"
	"math/rand"

	"ProjetRED/character"
)

// standRoll est un résultat possible de la Flèche (voir rouletteStand
// ci-dessous) : un nom façon "réveil de Stand", et l'effet appliqué au
// joueur qui l'utilise.
type standRoll struct {
	Name  string
	Apply func(c *character.Character) string
}

// rouletteStand est le tirage au sort de la Flèche : comme dans JoJo, elle
// réveille un Stand aléatoire chez celui qu'elle transperce — parfois un
// don puissant, parfois un fardeau, et parfois un rejet pur et simple du
// corps qui la reçoit. L'utiliser est donc un vrai pari, pas un simple soin
// garanti.
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

// useArrow consomme une Flèche (voir combat/items.go) et applique un effet
// aléatoire de la roulette de Stand ci-dessus.
func useArrow(player *character.Character) (bool, string) {
	tirage := rouletteStand[rand.Intn(len(rouletteStand))]
	effet := tirage.Apply(player)
	msg := fmt.Sprintf("La Flèche te transperce... Stand éveillé : %s ! %s", tirage.Name, effet)
	return true, msg
}
