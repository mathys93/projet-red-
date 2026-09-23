package world

import "ProjetRED/boss"

type Zone struct {
	Name        string
	Description string
	Boss        *boss.Boss
	Cleared     bool
	IsShop      bool
	IsFarm      bool
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
			{
				Name:        "Zone 5 - Le Terrain Vague",
				Description: "Un chien errant garde ce bout de désert. Il revient toujours : de quoi s'entraîner autant que tu veux pour de l'XP et des Fragments.",
				Boss:        boss.NewIggy(),
				IsFarm:      true,
			},
		},
		BonusBoss: boss.NewPucci(),
	}
}

func (w *World) Locked(z *Zone) bool {
	if z.IsShop || z.IsFarm {
		return false
	}
	var prev *Zone
	for _, zone := range w.Zones {
		if zone == z {
			break
		}
		if !zone.IsShop && !zone.IsFarm {
			prev = zone
		}
	}
	return prev != nil && !prev.Cleared
}

func (w *World) AllCleared() bool {
	for _, z := range w.Zones {
		if z.IsShop || z.IsFarm {
			continue
		}
		if !z.Cleared {
			return false
		}
	}
	return true
}
