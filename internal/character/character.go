package character

import "fmt"

type Character struct {
	Name  string
	Level int

	AP int
	DP int

	LP    int
	MaxLP int

	XP       int
	MP       int
	MaxMP    int
	Fragment int

	Inventory []string

	Moves []Move

	Stand *Stand

	SkipBossNextTurn  bool
	HasResurrectCharm bool
}

type Move struct {
	Name        string
	Description string
	MPCost      int
	Perform     func(user, target *Character) string
}

func DefaultMoves() []Move {
	return []Move{
		{
			Name:        "Coup de poing",
			Description: "Une frappe franche, sans coût de MP.",
			Perform: func(user, target *Character) string {
				dmg := target.TakeDamage(user.AP)
				return fmt.Sprintf("%s frappe %s d'un coup de poing ! %d dégâts.", user.Name, target.Name, dmg)
			},
		},
		{
			Name:        "Attaque de Stand",
			Description: "Une attaque puissante (dégâts doublés).",
			MPCost:      4,
			Perform: func(user, target *Character) string {
				dmg := target.TakeDamage(user.AP * 2)
				return fmt.Sprintf("%s déchaîne son Stand sur %s ! %d dégâts.", user.Name, target.Name, dmg)
			},
		},
	}
}

func New(name string, level, ap, dp, lp int) *Character {
	return &Character{
		Name:      name,
		Level:     level,
		AP:        ap,
		DP:        dp,
		LP:        lp,
		MaxLP:     lp,
		MP:        30,
		MaxMP:     30,
		Fragment:  20,
		Inventory: []string{},
		Moves:     DefaultMoves(),
	}
}

func (c *Character) IsAlive() bool { return c.LP > 0 }

func XPForLevel(level int) int {
	return 40 + (level-1)*30
}

func (c *Character) GainXP(amount int) []string {
	if amount <= 0 {
		return nil
	}
	c.XP += amount

	var messages []string
	for c.XP >= XPForLevel(c.Level) {
		c.XP -= XPForLevel(c.Level)
		c.Level++
		c.AP += 2
		c.DP++
		c.MaxLP += 15
		c.MaxMP += 7
		c.LP = c.MaxLP
		c.MP = c.MaxMP
		messages = append(messages, fmt.Sprintf(
			"%s passe niveau %d ! (AP %d, DP %d, PV max %d, MP max %d)",
			c.Name, c.Level, c.AP, c.DP, c.MaxLP, c.MaxMP,
		))
	}
	return messages
}

func (c *Character) TakeDamage(raw int) int {
	dmg := raw - c.DP
	if dmg < 0 {
		dmg = 0
	}
	c.LP -= dmg
	if c.LP < 0 {
		c.LP = 0
	}
	return dmg
}

func (c *Character) Heal(amount int) {
	c.LP += amount
	if c.LP > c.MaxLP {
		c.LP = c.MaxLP
	}
}

func (c *Character) AddItem(item string) {
	c.Inventory = append(c.Inventory, item)
}

func (c *Character) RemoveItem(item string) bool {
	for i, it := range c.Inventory {
		if it == item {
			c.Inventory = append(c.Inventory[:i], c.Inventory[i+1:]...)
			return true
		}
	}
	return false
}
