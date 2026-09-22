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
//    Act(turn, self, player)     -> ce qui se passe au tour du boss
//    ActOptions(self, player)    -> les sous-options du menu ACT (une par
//                                   action "lore" possible : parler,
//                                   observer, chanter...)
//
//  Pour l'instant Act() se contente d'une riposte simple (DefaultPattern) :
//  remplace son corps par la vraie logique du boss (patterns d'esquive,
//  projectiles ASCII qui bougent, attaques qui changent selon les PV
//  restants, etc.) sans toucher au reste du programme (combat, map...).
//  Les ActOptions ci-dessous sont déjà écrites avec du texte propre à
//  chaque boss : ajoute/retire des options si tu veux enrichir le lore.
// ===================================================================

// DefaultPattern : comportement de base utilisé si un boss n'a pas encore
// de pattern personnalisé, ou en secours depuis les patterns ci-dessous.
type DefaultPattern struct{}

func (DefaultPattern) Act(turn int, self *Boss, player *character.Character) string {
	dmg := player.TakeDamage(self.AP)
	return fmt.Sprintf("%s riposte : %d dégâts.", self.Name, dmg)
}

func (DefaultPattern) ActOptions(self *Boss, player *character.Character) []ActOption {
	return []ActOption{
		{
			Label:       "Observer",
			Description: fmt.Sprintf("Regarder %s attentivement.", self.Name),
			Resolve: func(self *Boss, player *character.Character) string {
				return fmt.Sprintf("Tu observes %s...", self.Name)
			},
		},
	}
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

func (p DioPattern) ActOptions(self *Boss, player *character.Character) []ActOption {
	return []ActOption{
		{
			Label:       "Observer",
			Description: "Étudier son Stand, The World.",
			Resolve: func(self *Boss, player *character.Character) string {
				return fmt.Sprintf("%s te regarde avec mépris. \"Za Warudo...\"", self.Name)
			},
		},
		{
			Label:       "Parler de Jonathan",
			Description: "Lui rappeler le corps qu'il a volé.",
			Resolve: func(self *Boss, player *character.Character) string {
				return fmt.Sprintf("%s se crispe. \"Ne prononce plus jamais ce nom, Ningen.\"", self.Name)
			},
		},
		{
			Label:       "Défier",
			Description: "Le provoquer ouvertement.",
			Resolve: func(self *Boss, player *character.Character) string {
				return fmt.Sprintf("%s éclate de rire. \"WRYYYYY ! Tu oses défier un dieu ?\"", self.Name)
			},
		},
	}
}

// --- Zone 2 : Diavolo -------------------------------------------------
// TODO: pattern de Diavolo (Stand "King Crimson" : suppression de temps,
// esquive impossible sur un tour donné, contre-attaque garantie, etc.)
type DiavoloPattern struct{}

func (p DiavoloPattern) Act(turn int, self *Boss, player *character.Character) string {
	// TODO: remplace par le vrai pattern de Diavolo.
	return DefaultPattern{}.Act(turn, self, player)
}

func (p DiavoloPattern) ActOptions(self *Boss, player *character.Character) []ActOption {
	return []ActOption{
		{
			Label:       "Observer",
			Description: "Regarder son Stand, King Crimson.",
			Resolve: func(self *Boss, player *character.Character) string {
				return fmt.Sprintf("%s ne supporte pas d'être observé.", self.Name)
			},
		},
		{
			Label:       "Mentionner Trish",
			Description: "Parler de sa fille, Trish Una.",
			Resolve: func(self *Boss, player *character.Character) string {
				return fmt.Sprintf("%s se fige un instant, puis retrouve son sang-froid glacial.", self.Name)
			},
		},
		{
			Label:       "Fixer son visage",
			Description: "Chercher à percer son identité changeante.",
			Resolve: func(self *Boss, player *character.Character) string {
				return fmt.Sprintf("Le visage de %s semble... instable. Il détourne le regard, furieux.", self.Name)
			},
		},
	}
}

// --- Zone 3 : Yoshikage Kira -------------------------------------------
// TODO: pattern de Kira (Stand "Killer Queen" : bombes à retardement sur
// les objets touchés, "Sheer Heart Attack" qui traque le joueur, etc.)
type KiraPattern struct{}

func (p KiraPattern) Act(turn int, self *Boss, player *character.Character) string {
	// TODO: remplace par le vrai pattern de Kira.
	return DefaultPattern{}.Act(turn, self, player)
}

func (p KiraPattern) ActOptions(self *Boss, player *character.Character) []ActOption {
	return []ActOption{
		{
			Label:       "Observer",
			Description: "L'étudier discrètement.",
			Resolve: func(self *Boss, player *character.Character) string {
				return fmt.Sprintf("%s a l'air d'un parfait inconnu. Trop parfait.", self.Name)
			},
		},
		{
			Label:       "Complimenter ses mains",
			Description: "Mentionner sa fascination pour les mains.",
			Resolve: func(self *Boss, player *character.Character) string {
				return fmt.Sprintf("%s sourit, gêné. \"J'apprécie simplement une vie tranquille, c'est tout.\"", self.Name)
			},
		},
		{
			Label:       "Parler de Reimi",
			Description: "Évoquer Reimi Sugimoto.",
			Resolve: func(self *Boss, player *character.Character) string {
				return fmt.Sprintf("%s perd son sourire poli. Son regard devient glacial.", self.Name)
			},
		},
	}
}

// --- Boss bonus : Enrico Pucci (Disque) --------------------------------
// TODO: pattern de Pucci (Stand "Made in Heaven" / Disque de Pucci :
// vol de Stand, accélération du temps en fin de combat, etc.)
type PucciPattern struct{}

func (p PucciPattern) Act(turn int, self *Boss, player *character.Character) string {
	// TODO: remplace par le vrai pattern de Pucci.
	return DefaultPattern{}.Act(turn, self, player)
}

func (p PucciPattern) ActOptions(self *Boss, player *character.Character) []ActOption {
	return []ActOption{
		{
			Label:       "Observer",
			Description: "Regarder le Disque qu'il tient.",
			Resolve: func(self *Boss, player *character.Character) string {
				return fmt.Sprintf("%s tient son Disque entre ses mains.", self.Name)
			},
		},
		{
			Label:       "Parler de DIO",
			Description: "Évoquer son maître déchu.",
			Resolve: func(self *Boss, player *character.Character) string {
				return fmt.Sprintf("%s ferme les yeux. \"DIO m'a montré le Ciel. Je ne faiblirai pas.\"", self.Name)
			},
		},
		{
			Label:       "Prier",
			Description: "Prier avec lui, ou contre lui.",
			Resolve: func(self *Boss, player *character.Character) string {
				return fmt.Sprintf("%s récite un verset. \"Tout ceci n'est qu'un chemin vers le Paradis.\"", self.Name)
			},
		},
	}
}
