// Package world définit la carte du jeu : trois zones, chacune gardée par
// un boss. Un quatrième boss (Pucci) est disponible en tant que boss bonus,
// débloqué une fois les trois zones nettoyées.
package world

import "ProjetRED/boss"

// Zone représente une zone de la carte, avec le boss qui la garde.
type Zone struct {
	Name        string
	Description string
	Boss        *boss.Boss
	Cleared     bool
}

// World contient les zones de la carte et le boss bonus.
type World struct {
	Zones     []*Zone
	BonusBoss *boss.Boss // débloqué quand toutes les zones sont Cleared
}

// New construit la carte à trois zones utilisée par le jeu.
func New() *World {
	return &World{
		Zones: []*Zone{
			{
				Name:        "Zone 1 - Le Manoir de DIO",
				Description: "Un manoir figé hors du temps. DIO t'y attend.",
				Boss:        boss.NewDio(),
			},
			{
				Name:        "Zone 2 - La Villa de Diavolo",
				Description: "Personne ne connaît son visage. Diavolo rôde ici.",
				Boss:        boss.NewDiavolo(),
			},
			{
				Name:        "Zone 3 - Morioh",
				Description: "Une ville tranquille... en apparence. Kira s'y cache.",
				Boss:        boss.NewKira(),
			},
		},
		BonusBoss: boss.NewPucci(),
	}
}

// AllCleared indique si les trois zones ont été nettoyées (boss vaincus ou
// épargnés), ce qui débloque le boss bonus.
func (w *World) AllCleared() bool {
	for _, z := range w.Zones {
		if !z.Cleared {
			return false
		}
	}
	return true
}
