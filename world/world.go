// Package world définit la carte du jeu : trois zones, chacune gardée par
// un boss. Un quatrième boss (Pucci) est disponible en tant que boss bonus,
// débloqué une fois les trois zones nettoyées.
package world

import "ProjetRED/boss"

// Zone représente une zone de la carte. La plupart gardent un boss, mais
// une zone peut aussi être une boutique (IsShop) : Boss est alors nil, et
// elle reste toujours visitable (jamais marquée Cleared).
type Zone struct {
	Name        string
	Description string
	Boss        *boss.Boss // nil quand IsShop est vrai
	Cleared     bool
	IsShop      bool
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
			{
				Name:        "Zone 4 - La Boutique",
				Description: "Un carrefour marchand entre deux affrontements : le Marchand et le Forgeron y attendent le client. Rien à craindre ici.",
				IsShop:      true,
			},
		},
		BonusBoss: boss.NewPucci(),
	}
}

// AllCleared indique si les trois zones de boss ont été nettoyées (boss
// vaincus ou épargnés), ce qui débloque le boss bonus. La boutique n'est
// jamais "nettoyée" (ce n'est pas un combat) : elle n'entre pas en compte.
func (w *World) AllCleared() bool {
	for _, z := range w.Zones {
		if z.IsShop {
			continue
		}
		if !z.Cleared {
			return false
		}
	}
	return true
}
