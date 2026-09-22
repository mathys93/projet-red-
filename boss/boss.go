// Package boss définit les boss affrontables en combat façon Undertale :
// un boss est un Character (voir package character) auquel s'ajoutent son
// artwork ASCII et son "pattern" (comportement au tour par tour).
package boss

import (
	"fmt"

	"ProjetRED/ascii"
	"ProjetRED/character"
)

// Taille utile de la boîte de combat pour l'artwork (voir package combat).
const (
	ArtWidth  = 46
	ArtHeight = 14
)

// Pattern décrit le comportement d'un boss pendant son tour, et ce qui
// s'affiche quand le joueur choisit ACT. C'est CETTE interface que tu
// personnalises boss par boss (voir patterns.go) : c'est toi qui t'occupes
// du pattern de chaque boss, ce fichier ne fait que poser la structure.
type Pattern interface {
	// Act est appelé après l'action du joueur (sauf MERCY) : c'est ici que
	// le boss riposte / esquive / lance son attaque spéciale. `turn` est le
	// numéro du tour (démarre à 1). Renvoie le message affiché sous la
	// boîte de combat.
	Act(turn int, self *Boss, player *character.Character) string

	// ActText est le texte affiché quand le joueur choisit l'option ACT du
	// menu (l'équivalent des sous-menus "Check/Talk/..." d'Undertale).
	ActText(self *Boss, player *character.Character) string
}

// Boss est un ennemi affrontable : ses statistiques (Character embarqué),
// son artwork, et son pattern de combat.
type Boss struct {
	*character.Character
	Art     string
	Zone    string
	Pattern Pattern
}

// New construit un boss, charge son artwork depuis le package ascii (par sa
// clé, ex: "dio") et le redimensionne pour la boîte de combat.
func New(name string, level, ap, dp, lp int, artKey string, pattern Pattern) *Boss {
	raw, ok := ascii.Get(artKey)
	if !ok {
		raw = fmt.Sprintf("(art manquant: %q)", artKey)
	}
	if pattern == nil {
		pattern = DefaultPattern{}
	}
	return &Boss{
		Character: character.New(name, level, ap, dp, lp),
		Art:       ascii.Fit(raw, ArtWidth, ArtHeight),
		Pattern:   pattern,
	}
}

// Turn fait jouer le tour du boss (délègue à son Pattern).
func (b *Boss) Turn(turn int, player *character.Character) string {
	return b.Pattern.Act(turn, b, player)
}

// Act renvoie le texte de l'option ACT (délègue à son Pattern).
func (b *Boss) ActDescription(player *character.Character) string {
	return b.Pattern.ActText(b, player)
}

// ---------------------------------------------------------------
// Les 4 boss du jeu (3 zones + 1 boss bonus/secret) : voir patterns.go
// pour leurs comportements, et world/world.go pour leur placement sur
// la carte.
// ---------------------------------------------------------------

// Les statistiques ci-dessous sont des valeurs de départ raisonnables
// (juste de quoi tester le jeu de bout en bout) : à ajuster une fois que
// tu auras remplacé les patterns par défaut dans patterns.go par le vrai
// comportement de chaque boss.

// NewDio crée le boss de la Zone 1.
func NewDio() *Boss {
	return New("DIO", 5, 6, 2, 40, "dio", DioPattern{})
}

// NewDiavolo crée le boss de la Zone 2.
func NewDiavolo() *Boss {
	return New("Diavolo", 6, 7, 3, 50, "diavolo", DiavoloPattern{})
}

// NewKira crée le boss de la Zone 3.
func NewKira() *Boss {
	return New("Yoshikage Kira", 6, 6, 3, 45, "kira", KiraPattern{})
}

// NewPucci crée le boss bonus (non assigné à une zone par défaut, voir
// world.go) : Enrico Pucci et son Disque, débloqué après les 3 zones.
func NewPucci() *Boss {
	return New("Enrico Pucci", 8, 8, 4, 70, "pucci", PucciPattern{})
}
