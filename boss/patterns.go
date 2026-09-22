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
// Pattern de DIO (Stand "The World") : deux tours de mise en jambe, puis un
// troisième tour où le temps s'arrête - un coup lourd qui ignore totalement
// la Defense Power. Le deuxième tour sert de signal clair ("il prépare
// quelque chose") : c'est le moment de se soigner avant l'impact plutôt que
// de continuer à frapper aveuglément.
type DioPattern struct{}

func (p DioPattern) Act(turn int, self *Boss, player *character.Character) string {
	switch turn % 3 {
	case 1:
		dmg := player.TakeDamage(self.AP)
		return fmt.Sprintf("%s assène un coup de poing sec. %d dégâts.", self.Name, dmg)
	case 2:
		dmg := player.TakeDamage(self.AP)
		return fmt.Sprintf("%s recule et sort sa montre à gousset. Tu sens que quelque chose se prépare... %d dégâts.", self.Name, dmg)
	default:
		raw := self.AP * 2
		player.LP -= raw
		if player.LP < 0 {
			player.LP = 0
		}
		return fmt.Sprintf("\"ZA WARUDO ! TOKI YO TOMARE !\" Le temps s'arrête, impossible d'esquiver : %d dégâts bruts (ignore la Defense Power) !", raw)
	}
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
// Pattern de Diavolo (Stand "King Crimson") : tous les 4 tours, il "efface"
// quelques secondes de dégâts en se soignant d'une fraction de ses PV max,
// ce qui punit un joueur qui grignote lentement au lieu de concentrer ses
// coups les plus forts (Attaque de Stand). Sous 1/3 de ses PV, il s'enrage
// et frappe plus fort - le combat devient plus dangereux juste avant la
// victoire, pas plus facile.
type DiavoloPattern struct{}

func (p DiavoloPattern) Act(turn int, self *Boss, player *character.Character) string {
	ap := self.AP
	enraged := self.LP*3 <= self.MaxLP
	if enraged {
		ap += self.AP / 2
	}

	if turn%4 == 0 {
		heal := self.MaxLP / 6
		self.LP += heal
		if self.LP > self.MaxLP {
			self.LP = self.MaxLP
		}
		dmg := player.TakeDamage(ap)
		msg := fmt.Sprintf("\"King Crimson efface le temps...\" %s referme ses blessures (+%d PV) et riposte pour %d dégâts.", self.Name, heal, dmg)
		if enraged {
			msg += " Acculé, il frappe avec une rage décuplée."
		}
		return msg
	}

	dmg := player.TakeDamage(ap)
	msg := fmt.Sprintf("%s frappe avec la précision froide de King Crimson. %d dégâts.", self.Name, dmg)
	if enraged {
		msg += " Acculé, il frappe avec une rage décuplée."
	}
	return msg
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
// Pattern de Kira (Stand "Killer Queen") : ses bombes rendent ses dégâts
// normaux croissants avec la durée du combat, et tous les 4 tours "Sheer
// Heart Attack" fonce en ligne droite pour un coup garanti qui ignore la
// Defense Power. Un joueur qui s'éternise à soigner au lieu de finir le
// combat rapidement se fait rattraper par des dégâts de plus en plus lourds.
type KiraPattern struct{}

func (p KiraPattern) Act(turn int, self *Boss, player *character.Character) string {
	if turn%4 == 0 {
		raw := self.AP + turn
		player.LP -= raw
		if player.LP < 0 {
			player.LP = 0
		}
		return fmt.Sprintf("\"Sheer Heart Attack\" fonce droit sur toi, impossible à esquiver : %d dégâts bruts !", raw)
	}

	bonus := turn / 2
	dmg := player.TakeDamage(self.AP + bonus)
	return fmt.Sprintf("%s plante discrètement une bombe. %d dégâts (Killer Queen s'impatiente : plus le combat dure, plus elle frappe fort).", self.Name, dmg)
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
// Pattern de Pucci (Stand "Made in Heaven") : le plus complet des quatre,
// il cumule le soin périodique de Diavolo (tous les 4 tours), le coup
// garanti ignorant la Defense Power de DIO/Kira (tous les 5 tours, "le
// temps accélère") et l'enrage à PV bas. Gérer les trois mécaniques à la
// fois - sans jamais pouvoir se relâcher - est le vrai test du boss bonus.
type PucciPattern struct{}

func (p PucciPattern) Act(turn int, self *Boss, player *character.Character) string {
	ap := self.AP
	enraged := self.LP*3 <= self.MaxLP
	if enraged {
		ap += self.AP / 2
	}

	if turn%4 == 0 {
		heal := self.MaxLP / 8
		self.LP += heal
		if self.LP > self.MaxLP {
			self.LP = self.MaxLP
		}
		dmg := player.TakeDamage(ap)
		return fmt.Sprintf("%s referme une plaie grâce à son Disque (+%d PV) et riposte pour %d dégâts.", self.Name, heal, dmg)
	}

	if turn%5 == 0 {
		raw := ap * 2
		player.LP -= raw
		if player.LP < 0 {
			player.LP = 0
		}
		return fmt.Sprintf("\"Made in Heaven\" accélère le temps : impossible d'esquiver, %d dégâts bruts !", raw)
	}

	dmg := player.TakeDamage(ap)
	msg := fmt.Sprintf("%s frappe avec la certitude froide de celui qui a vu le Paradis. %d dégâts.", self.Name, dmg)
	if enraged {
		msg += " Le Disque s'affole : ses coups sont de plus en plus violents."
	}
	return msg
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
