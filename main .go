package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ============================================================
// STRUCTURES
// ============================================================

// Character représente un personnage jouable OU un boss.
// En Go, on utilise un struct pour regrouper des données liées,
// plutôt que des paramètres de fonction séparés.
type Character struct {
	Name  string
	AP    int      // Attack Power (dégâts infligés)
	DP    int      // Defense Power (réduction des dégâts reçus)
	INV   []string // Inventaire (liste d'objets, plus flexible qu'un string unique)
	LP    int      // Life Points actuels (HP)
	MaxLP int      // Life Points maximum
	XP    int      // Expérience accumulée
	MP    int      // Monnaie façon Undertale ("Gold"), ou points de magie selon ton design
}

// ============================================================
// CRÉATION / INITIALISATION
// ============================================================

// NewCharacter crée un personnage avec les stats passées en argument.
// Contrairement à ta version, ici les valeurs sont bien utilisées
// et la fonction RENVOIE le personnage créé (pointeur *Character,
// pour pouvoir le modifier plus tard sans le recopier).
func NewCharacter(name string, ap, dp, lp int) *Character {
	return &Character{
		Name:  name,
		AP:    ap,
		DP:    dp,
		INV:   []string{},
		LP:    lp,
		MaxLP: lp,
		XP:    0,
		MP:    0,
	}
}

// ============================================================
// AFFICHAGE
// ============================================================

// DisplayInfo affiche les stats du personnage dans la console,
// façon "menu stats" d'Undertale.
func (c *Character) DisplayInfo() {
	fmt.Println("========================================")
	fmt.Printf(" %s\n", strings.ToUpper(c.Name))
	fmt.Printf(" LP  %d / %d\n", c.LP, c.MaxLP)
	fmt.Printf(" AP  %d   DP  %d\n", c.AP, c.DP)
	fmt.Printf(" XP  %d   MP  %d\n", c.XP, c.MP)
	if len(c.INV) == 0 {
		fmt.Println(" INV  (vide)")
	} else {
		fmt.Printf(" INV  %s\n", strings.Join(c.INV, ", "))
	}
	fmt.Println("========================================")
}

// ============================================================
// LOGIQUE DE COMBAT (base à étoffer ensuite)
// ============================================================

// IsAlive renvoie true tant que le personnage a des LP > 0.
func (c *Character) IsAlive() bool {
	return c.LP > 0
}

// TakeDamage applique des dégâts en tenant compte de la défense (DP).
// On empêche les dégâts négatifs et on clamp les LP à 0 minimum.
func (c *Character) TakeDamage(rawDamage int) int {
	dmg := rawDamage - c.DP
	if dmg < 0 {
		dmg = 0
	}
	c.LP -= dmg
	if c.LP < 0 {
		c.LP = 0
	}
	return dmg
}

// Attack fait attaquer c sur target, et affiche le résultat.
func (c *Character) Attack(target *Character) {
	dmg := target.TakeDamage(c.AP)
	fmt.Printf("* %s attaque %s et inflige %d dégâts !\n", c.Name, target.Name, dmg)
	if !target.IsAlive() {
		fmt.Printf("* %s est terrassé(e) !\n", target.Name)
	}
}

// AddItem ajoute un objet à l'inventaire.
func (c *Character) AddItem(item string) {
	c.INV = append(c.INV, item)
}

// Heal soigne le personnage sans dépasser MaxLP.
func (c *Character) Heal(amount int) {
	c.LP += amount
	if c.LP > c.MaxLP {
		c.LP = c.MaxLP
	}
}

// ============================================================
// MENU DE COMBAT FIGHT / ACT / ITEM / MERCY
// ============================================================

// CombatMenu affiche le menu façon Undertale et lit le choix du joueur.
func CombatMenu(player *Character, enemy *Character) string {
	fmt.Println()
	fmt.Println("[ FIGHT ]  [ ACT ]  [ ITEM ]  [ MERCY ]")
	fmt.Print("> Ton choix : ")

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.ToUpper(strings.TrimSpace(input))
}

// RunCombatTurn exécute un tour de combat en fonction du choix du joueur.
// Renvoie true si le combat doit continuer, false s'il doit s'arrêter
// (victoire, fuite/mercy, ou mort du joueur).
func RunCombatTurn(player *Character, enemy *Character) bool {
	choice := CombatMenu(player, enemy)

	switch choice {
	case "FIGHT":
		player.Attack(enemy)
	case "ACT":
		fmt.Printf("* Tu observes %s...\n", enemy.Name)
		// TODO: ici tu pourras brancher des actions spécifiques par boss
		// (ex: "Parler", "Flatter", "Provoquer") selon le perso d'anime choisi
	case "ITEM":
		if len(player.INV) == 0 {
			fmt.Println("* Ton inventaire est vide !")
		} else {
			item := player.INV[0]
			player.INV = player.INV[1:]
			fmt.Printf("* Tu utilises %s.\n", item)
			player.Heal(10) // exemple simple, à adapter par objet
		}
	case "MERCY":
		fmt.Println("* Tu tentes d'épargner l'ennemi...")
		return false
	default:
		fmt.Println("* Choix non reconnu.")
		return true
	}

	if !enemy.IsAlive() {
		fmt.Printf("* %s a gagné le combat !\n", player.Name)
		return false
	}

	// Riposte simple du boss (à remplacer plus tard par des patterns
	// d'esquive façon "bullet hell" propres à chaque boss d'anime)
	if enemy.IsAlive() {
		enemy.Attack(player)
	}

	if !player.IsAlive() {
		fmt.Printf("* %s est tombé(e) au combat...\n", player.Name)
		return false
	}

	return true
}

// ============================================================
// MAIN (exemple de démo)
// ============================================================

func main() {
	hero := NewCharacter("Personnage 1", 5, 2, 20)
	boss := NewCharacter("Boss Anime", 4, 1, 30)

	hero.AddItem("Potion")

	hero.DisplayInfo()
	boss.DisplayInfo()

	for hero.IsAlive() && boss.IsAlive() {
		if !RunCombatTurn(hero, boss) {
			break
		}
	}
}
