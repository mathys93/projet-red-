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

// New crée un nouveau personnage avec des PV pleins. Fragment démarre à 20 :
// sans ça, un joueur fraîchement créé n'a jamais de quoi acheter quoi que ce
// soit chez le marchand ou le forgeron.
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
		Fragment:  20,
		Inventory: []string{},
		Moves:     DefaultMoves(),
	}
}

// IsAlive indique si le personnage est encore en vie.
func (c *Character) IsAlive() bool { return c.LP > 0 }

// XPForLevel renvoie l'XP nécessaire pour passer du niveau `level` au
// niveau suivant. La courbe est volontairement croissante (plus on monte
// de niveau, plus il faut d'XP) pour que la progression reste sentie sur
// toute la partie plutôt que de s'aplatir après les premiers combats.
func XPForLevel(level int) int {
	return 40 + (level-1)*30
}

// GainXP ajoute de l'XP au personnage et fait monter son niveau en chaîne
// (plusieurs niveaux d'un coup si l'XP gagnée est suffisante). Chaque
// niveau gagné augmente l'Attack/Defense Power et les PV/MP max, et
// restaure les PV/MP au maximum (une montée de niveau doit se ressentir
// tout de suite en combat). Renvoie un message par niveau gagné, prêt à
// être affiché.
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
		c.MaxMP += 5
		c.LP = c.MaxLP
		c.MP = c.MaxMP
		messages = append(messages, fmt.Sprintf(
			"%s passe niveau %d ! (AP %d, DP %d, PV max %d, MP max %d)",
			c.Name, c.Level, c.AP, c.DP, c.MaxLP, c.MaxMP,
		))
	}
	return messages
}

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
