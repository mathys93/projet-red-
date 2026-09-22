// Package character regroupe tout ce qui concerne le personnage joueur
// (et sert aussi de base aux boss, voir le package boss) : statistiques,
// points de vie, inventaire, fragments (monnaie du forgeron).
//
// Avant cette réorganisation, la struct Character était redéfinie
// plusieurs fois (main.go, utilitaire/disque_de_pucci.go, utilitaire/poison POT)
// avec des champs différents et incompatibles (LP/MaxLP vs HpCurrent/HpMax...),
// ce qui provoquait des conflits de compilation ("Character redeclared")
// dès que ces fichiers se retrouvaient dans le même package. Il n'existe
// maintenant plus qu'UNE seule définition, utilisée partout.
package character

import "fmt"

// Character représente un personnage (joueur ou boss).
type Character struct {
	Name    string
	Joestar string // lignée / titre JoJo du personnage (ex: "Joestar", "??")
	Level   int

	AP int // Attack Power (dégâts infligés)
	DP int // Defense Power (réduction des dégâts subis)

	LP    int // Life Points actuels
	MaxLP int // Life Points maximum

	XP       int
	MP       int // Magic/Stand Points actuels (coût des attaques de Stand en FIGHT)
	MaxMP    int
	Fragment int // monnaie utilisée chez le forgeron

	Inventory []string

	// Moves sont les actions proposées dans le sous-menu FIGHT (coup de
	// poing, attaque de Stand...). Personnalise cette liste après New()
	// pour donner des attaques propres à un personnage.
	Moves []Move
}

// Move est une action de combat sélectionnable en FIGHT. Perform applique
// l'effet (dégâts...) sur la cible et renvoie le message affiché comme
// résultat du tour du joueur ; le coût en MP est vérifié et déduit par
// l'appelant (voir combat.RunBattle) avant d'appeler Perform.
type Move struct {
	Name        string
	Description string
	MPCost      int
	Perform     func(user, target *Character) string
}

// DefaultMoves renvoie les actions de FIGHT de base communes à tout
// personnage : un coup de poing sans coût, et une attaque de Stand plus
// puissante qui consomme du MP.
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
			MPCost:      5,
			Perform: func(user, target *Character) string {
				dmg := target.TakeDamage(user.AP * 2)
				return fmt.Sprintf("%s déchaîne son Stand sur %s ! %d dégâts.", user.Name, target.Name, dmg)
			},
		},
	}
}

// New crée un nouveau personnage avec des PV pleins.
func New(name string, level, ap, dp, lp int) *Character {
	return &Character{
		Name:      name,
		Level:     level,
		AP:        ap,
		DP:        dp,
		LP:        lp,
		MaxLP:     lp,
		MP:        20,
		MaxMP:     20,
		Inventory: []string{},
		Moves:     DefaultMoves(),
	}
}

// IsAlive indique si le personnage est encore en vie.
func (c *Character) IsAlive() bool { return c.LP > 0 }

// TakeDamage applique des dégâts bruts (réduits par la Defense Power) et
// renvoie les dégâts réellement subis.
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

// Heal restaure des PV, plafonnés à MaxLP.
func (c *Character) Heal(amount int) {
	c.LP += amount
	if c.LP > c.MaxLP {
		c.LP = c.MaxLP
	}
}

// AddItem ajoute un objet à l'inventaire.
func (c *Character) AddItem(item string) {
	c.Inventory = append(c.Inventory, item)
}

// RemoveItem retire une occurrence de l'objet donné de l'inventaire.
// Renvoie false si l'objet n'était pas présent.
func (c *Character) RemoveItem(item string) bool {
	for i, it := range c.Inventory {
		if it == item {
			c.Inventory = append(c.Inventory[:i], c.Inventory[i+1:]...)
			return true
		}
	}
	return false
}

// CountItem compte le nombre d'occurrences d'un objet dans l'inventaire.
func (c *Character) CountItem(item string) int {
	count := 0
	for _, it := range c.Inventory {
		if it == item {
			count++
		}
	}
	return count
}

// String permet un affichage rapide en debug (fmt.Println(perso) etc.)
func (c *Character) String() string {
	return fmt.Sprintf("%s (LV %d) HP %d/%d AP %d DP %d", c.Name, c.Level, c.LP, c.MaxLP, c.AP, c.DP)
}
