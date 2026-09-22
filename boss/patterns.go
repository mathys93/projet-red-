package boss

import (
	"fmt"

	"ProjetRED/character"
)

// ===================================================================
//  PATTERNS DES BOSS
//
//  C'est ICI que tu personnalises le comportement de chaque boss
//  (attaques spéciales, esquives, dialogues, phases...). Chaque type
//  ci-dessous implémente l'interface Pattern (voir boss.go) :
//
//    Act(turn, self, player)  -> ce qui se passe au tour du boss
//    ActText(self, player)    -> le texte affiché avec le menu ACT
//
//  Pour l'instant chaque pattern se contente d'une riposte simple
//  (DefaultPattern) : remplace le corps de Act()/ActText() par la vraie
//  logique du boss (patterns d'esquive, projectiles ASCII qui bougent,
//  attaques qui changent selon les PV restants, etc.) sans toucher au
//  reste du programme (combat, map...).
// ===================================================================

// DefaultPattern : comportement de base utilisé si un boss n'a pas encore
// de pattern personnalisé, ou en secours depuis les patterns ci-dessous.
type DefaultPattern struct{}

func (DefaultPattern) Act(turn int, self *Boss, player *character.Character) string {
	dmg := player.TakeDamage(self.AP)
	return fmt.Sprintf("%s riposte : %d dégâts.", self.Name, dmg)
}

func (DefaultPattern) ActText(self *Boss, player *character.Character) string {
	return fmt.Sprintf("Tu observes %s...", self.Name)
}

// --- Zone 1 : DIO ---------------------------------------------------
// TODO: pattern de DIO (Stand "The World" : arrêt du temps, rafale de
// coups de poing "MUDA MUDA MUDA", esquive du joueur limitée pendant le
// tour où le temps est arrêté, etc.)
type DioPattern struct{}

func (p DioPattern) Act(turn int, self *Boss, player *character.Character) string {
	// TODO: remplace par le vrai pattern de DIO.
	return DefaultPattern{}.Act(turn, self, player)
}

func (p DioPattern) ActText(self *Boss, player *character.Character) string {
	// TODO: texte d'observation propre à DIO.
	return fmt.Sprintf("%s te regarde avec mépris. \"Za Warudo...\"", self.Name)
}

// --- Zone 2 : Diavolo -------------------------------------------------
// TODO: pattern de Diavolo (Stand "King Crimson" : suppression de temps,
// esquive impossible sur un tour donné, contre-attaque garantie, etc.)
type DiavoloPattern struct{}

func (p DiavoloPattern) Act(turn int, self *Boss, player *character.Character) string {
	// TODO: remplace par le vrai pattern de Diavolo.
	return DefaultPattern{}.Act(turn, self, player)
}

func (p DiavoloPattern) ActText(self *Boss, player *character.Character) string {
	// TODO: texte d'observation propre à Diavolo.
	return fmt.Sprintf("%s ne supporte pas d'être observé.", self.Name)
}

// --- Zone 3 : Yoshikage Kira -------------------------------------------
// TODO: pattern de Kira (Stand "Killer Queen" : bombes à retardement sur
// les objets touchés, "Sheer Heart Attack" qui traque le joueur, etc.)
type KiraPattern struct{}

func (p KiraPattern) Act(turn int, self *Boss, player *character.Character) string {
	// TODO: remplace par le vrai pattern de Kira.
	return DefaultPattern{}.Act(turn, self, player)
}

func (p KiraPattern) ActText(self *Boss, player *character.Character) string {
	// TODO: texte d'observation propre à Kira.
	return fmt.Sprintf("%s a l'air d'un parfait inconnu. Trop parfait.", self.Name)
}

// --- Boss bonus : Enrico Pucci (Disque) --------------------------------
// TODO: pattern de Pucci (Stand "Made in Heaven" / Disque de Pucci :
// vol de Stand, accélération du temps en fin de combat, etc.)
type PucciPattern struct{}

func (p PucciPattern) Act(turn int, self *Boss, player *character.Character) string {
	// TODO: remplace par le vrai pattern de Pucci.
	return DefaultPattern{}.Act(turn, self, player)
}

func (p PucciPattern) ActText(self *Boss, player *character.Character) string {
	// TODO: texte d'observation propre à Pucci.
	return fmt.Sprintf("%s tient son Disque entre ses mains.", self.Name)
}
