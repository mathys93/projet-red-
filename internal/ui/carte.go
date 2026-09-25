package ui

import (
	"fmt"

	"ProjetRED/internal/world"
)

func SelectZone(w *world.World) *world.Zone {
	SetScene(SceneCarte)
	ts := NewSession()
	defer ts.Restore()

	selected := 0
	note := ""

	render := func() {
		ClearScreen()
		contentHeight := 4 + len(w.Zones)*3
		if note != "" {
			contentHeight += 2
		}
		PrintPadding(contentHeight)
		fmt.Fprintln(Out)
		fmt.Fprintln(Out, ColYellow+"======================  CARTE  ======================"+ColReset)
		fmt.Fprintln(Out)
		for i, z := range w.Zones {
			marker := "  "
			if i == selected {
				marker = ColRed + "❤ " + ColReset
			}
			status := ""
			switch {
			case z.IsFarm:
				status = ColWhite + "  [toujours ouverte]" + ColReset
			case z.Cleared:
				status = ColYellow + "  [zone terminée]" + ColReset
			case w.Locked(z):
				status = ColWhite + "  [verrouillée]" + ColReset
			}
			fmt.Fprintf(Out, "%s%s%s\n", marker, z.Name, status)
			fmt.Fprintf(Out, "    %s\n\n", z.Description)
		}
		if note != "" {
			fmt.Fprintln(Out, ColWhite+"* "+note+ColReset)
			fmt.Fprintln(Out)
		}
		fmt.Fprintln(Out, ColWhite+"(↑ ↓ pour choisir, Entrée pour valider, Échap pour le menu pause)"+ColReset)
	}

	for {
		CurrentScreen = render
		render()

		last := len(w.Zones) - 1
		switch ts.ReadKey() {
		case KeyUp, KeyLeft:
			selected--
			if selected < 0 {
				selected = last
			}
		case KeyDown, KeyRight:
			selected++
			if selected > last {
				selected = 0
			}
		case KeyQuit:
			return nil
		case KeyEnter:
			z := w.Zones[selected]
			if z.Cleared {
				note = fmt.Sprintf("%s est déjà terminée.", z.Name)
				continue
			}
			if w.Locked(z) {
				note = fmt.Sprintf("%s est verrouillée : termine d'abord la zone précédente.", z.Name)
				continue
			}
			return z
		}
	}
}

func ConfirmBonusBoss(bossName string) bool {
	SetScene(SceneCarte)
	ts := NewSession()
	defer ts.Restore()

	render := func() {
		ClearScreen()
		PrintPadding(5)
		fmt.Fprintln(Out)
		fmt.Fprintln(Out, ColYellow+"Les trois zones sont terminées."+ColReset)
		fmt.Fprintf(Out, "Un dernier adversaire t'attend : %s%s%s\n\n", ColRed, bossName, ColReset)
		fmt.Fprintln(Out, ColWhite+"L'affronter ? (Entrée = oui, X = non)"+ColReset)
	}

	for {
		CurrentScreen = render
		render()

		switch ts.ReadKey() {
		case KeyEnter:
			return true
		case KeyBack, KeyQuit:
			return false
		}
	}
}
