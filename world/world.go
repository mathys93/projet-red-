package world

import "ProjetRED/boss"

type Zone struct {
	Name        string
	Description string
	Boss        *boss.Boss
	Cleared     bool
	IsShop      bool
}

type World struct {
	Zones     []*Zone
	BonusBoss *boss.Boss
}

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

func (w *World) Locked(z *Zone) bool {
	if z.IsShop {
		return false
	}
	var prev *Zone
	for _, zone := range w.Zones {
		if zone == z {
			break
		}
		if !zone.IsShop {
			prev = zone
		}
	}
	return prev != nil && !prev.Cleared
}

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
