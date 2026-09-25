package character

import (
	"fmt"
	"math/rand"
)

type Rarity int

const (
	RarityCommun Rarity = iota
	RarityRare
	RarityLegendaire
)

func (r Rarity) String() string {
	switch r {
	case RarityRare:
		return "Rare"
	case RarityLegendaire:
		return "Légendaire"
	}
	return "Commun"
}

func (r Rarity) Bonus() int {
	switch r {
	case RarityRare:
		return 5
	case RarityLegendaire:
		return 9
	}
	return 2
}

type Stand struct {
	Name   string
	Rarity Rarity
	Move   Move
}

var standPool = []Stand{
	{
		Name:   "Hermit Purple",
		Rarity: RarityCommun,
		Move: Move{
			Name:        "Hermit Purple : Ronces",
			Description: "Coup de grâce : des ronces cinglantes.",
			MPCost:      7,
			Perform: func(user, target *Character) string {
				dmg := target.TakeDamage(user.AP*2 + 6)
				return fmt.Sprintf("Des ronces jaillissent et fouettent %s ! %d dégâts.", target.Name, dmg)
			},
		},
	},
	{
		Name:   "Silver Chariot",
		Rarity: RarityCommun,
		Move: Move{
			Name:        "Silver Chariot : Estoc",
			Description: "Coup de grâce : une rafale de coups d'épée.",
			MPCost:      8,
			Perform: func(user, target *Character) string {
				dmg := target.TakeDamage(user.AP*2 + 9)
				return fmt.Sprintf("La rapière de Silver Chariot transperce %s ! %d dégâts.", target.Name, dmg)
			},
		},
	},
	{
		Name:   "Whitesnake",
		Rarity: RarityCommun,
		Move: Move{
			Name:        "Whitesnake : Vol de Disque",
			Description: "Coup de grâce : frappe et récupère du MP.",
			MPCost:      6,
			Perform: func(user, target *Character) string {
				dmg := target.TakeDamage(user.AP*2 + 4)
				user.MP += 10
				if user.MP > user.MaxMP {
					user.MP = user.MaxMP
				}
				return fmt.Sprintf("Whitesnake arrache un disque à %s ! %d dégâts, +10 MP.", target.Name, dmg)
			},
		},
	},
	{
		Name:   "Crazy Diamond",
		Rarity: RarityRare,
		Move: Move{
			Name:        "Crazy Diamond : Réparation",
			Description: "Restaure les deux tiers de tes PV maximum.",
			MPCost:      10,
			Perform: func(user, target *Character) string {
				soin := user.MaxLP * 2 / 3
				user.Heal(soin)
				return fmt.Sprintf("Crazy Diamond recompose tes blessures ! +%d PV.", soin)
			},
		},
	},
	{
		Name:   "Killer Queen",
		Rarity: RarityRare,
		Move: Move{
			Name:        "Killer Queen : Détonation",
			Description: "Coup de grâce : explosion qui ignore la Defense Power.",
			MPCost:      12,
			Perform: func(user, target *Character) string {
				raw := user.AP*3 + 12
				target.LP -= raw
				if target.LP < 0 {
					target.LP = 0
				}
				return fmt.Sprintf("Killer Queen fait détoner %s ! %d dégâts bruts.", target.Name, raw)
			},
		},
	},
	{
		Name:   "Sticky Fingers",
		Rarity: RarityRare,
		Move: Move{
			Name:        "Sticky Fingers : Fermeture",
			Description: "Coup de grâce : frappe et te soigne de la moitié des dégâts.",
			MPCost:      11,
			Perform: func(user, target *Character) string {
				dmg := target.TakeDamage(user.AP*2 + 14)
				user.Heal(dmg / 2)
				return fmt.Sprintf("Sticky Fingers ouvre une fermeture dans %s ! %d dégâts, tu récupères %d PV.", target.Name, dmg, dmg/2)
			},
		},
	},
	{
		Name:   "Star Platinum",
		Rarity: RarityLegendaire,
		Move: Move{
			Name:        "Star Platinum : ORA ORA ORA",
			Description: "Coup de grâce : la rafale la plus dévastatrice du jeu.",
			MPCost:      14,
			Perform: func(user, target *Character) string {
				dmg := target.TakeDamage(user.AP*3 + 16)
				return fmt.Sprintf("ORA ORA ORA ORA ! Star Platinum écrase %s : %d dégâts.", target.Name, dmg)
			},
		},
	},
	{
		Name:   "The World",
		Rarity: RarityLegendaire,
		Move: Move{
			Name:        "The World : Za Warudo",
			Description: "Coup de grâce : fige le temps, l'adversaire saute son tour.",
			MPCost:      15,
			Perform: func(user, target *Character) string {
				dmg := target.TakeDamage(user.AP*2 + 14)
				user.SkipBossNextTurn = true
				return fmt.Sprintf("ZA WARUDO ! Le temps s'arrête : %d dégâts, %s ne jouera pas.", dmg, target.Name)
			},
		},
	},
	{
		Name:   "King Crimson",
		Rarity: RarityLegendaire,
		Move: Move{
			Name:        "King Crimson : Epitaph",
			Description: "Coup de grâce : efface ta mort, tu reviendras une fois.",
			MPCost:      14,
			Perform: func(user, target *Character) string {
				dmg := target.TakeDamage(user.AP*2 + 10)
				user.HasResurrectCharm = true
				return fmt.Sprintf("King Crimson efface le futur : %d dégâts, et ta mort est annulée une fois.", dmg)
			},
		},
	},
}

func RollStand() Stand {
	var wanted Rarity
	switch n := rand.Intn(100); {
	case n < 10:
		wanted = RarityLegendaire
	case n < 40:
		wanted = RarityRare
	default:
		wanted = RarityCommun
	}

	var pool []Stand
	for _, s := range standPool {
		if s.Rarity == wanted {
			pool = append(pool, s)
		}
	}
	if len(pool) == 0 {
		pool = standPool
	}
	return pool[rand.Intn(len(pool))]
}

func (c *Character) EquipStand(s Stand) string {
	previous := ""
	if c.Stand != nil {
		c.AP -= c.Stand.Rarity.Bonus()
		previous = c.Stand.Name
	}

	equipped := s
	c.Stand = &equipped
	c.AP += s.Rarity.Bonus()
	c.Moves = append(DefaultMoves(), s.Move)

	if previous != "" {
		return fmt.Sprintf("La Flèche te transperce... %s remplace %s ! [%s] Attack Power +%d, nouvelle action : %s.",
			s.Name, previous, s.Rarity, s.Rarity.Bonus(), s.Move.Name)
	}
	return fmt.Sprintf("La Flèche te transperce... %s s'éveille en toi ! [%s] Attack Power +%d, nouvelle action : %s.",
		s.Name, s.Rarity, s.Rarity.Bonus(), s.Move.Name)
}
