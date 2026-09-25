package ui

import (
	"fmt"

	"ProjetRED/internal/fond"
	"ProjetRED/internal/terminal"
	"ProjetRED/internal/world"
)

func SelectZone(w *world.World) *world.Zone {
	terminal.SetScene(fond.Carte)
	ts := terminal.NewSession()
	defer ts.Restore()

	selected := 0
	note := ""

	render := func() {
		terminal.ClearScreen()
		contentHeight := 4 + len(w.Zones)*3
		if note != "" {
			contentHeight += 2
		}
		terminal.PrintPadding(contentHeight)
		fmt.Fprintln(terminal.Out)
		fmt.Fprintln(terminal.Out, terminal.ColYellow+"======================  CARTE  ======================"+terminal.ColReset)
		fmt.Fprintln(terminal.Out)
		for i, z := range w.Zones {
			marker := "  "
			if i == selected {
				marker = terminal.ColRed + "❤ " + terminal.ColReset
			}
			status := ""
			switch {
			case z.IsFarm:
				status = terminal.ColWhite + "  [toujours ouverte]" + terminal.ColReset
			case z.Cleared:
				status = terminal.ColYellow + "  [zone terminée]" + terminal.ColReset
			case w.Locked(z):
				status = terminal.ColWhite + "  [verrouillée]" + terminal.ColReset
			}
			fmt.Fprintf(terminal.Out, "%s%s%s\n", marker, z.Name, status)
			fmt.Fprintf(terminal.Out, "    %s\n\n", z.Description)
		}
		if note != "" {
			fmt.Fprintln(terminal.Out, terminal.ColWhite+"* "+note+terminal.ColReset)
			fmt.Fprintln(terminal.Out)
		}
		fmt.Fprintln(terminal.Out, terminal.ColWhite+"(↑ ↓ pour choisir, Entrée pour valider, Échap pour le menu pause)"+terminal.ColReset)
	}

	for {
		terminal.CurrentScreen = render
		render()

		last := len(w.Zones) - 1
		switch ts.ReadKey() {
		case terminal.KeyUp, terminal.KeyLeft:
			selected--
			if selected < 0 {
				selected = last
			}
		case terminal.KeyDown, terminal.KeyRight:
			selected++
			if selected > last {
				selected = 0
			}
		case terminal.KeyQuit:
			return nil
		case terminal.KeyEnter:
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
	terminal.SetScene(fond.Carte)
	ts := terminal.NewSession()
	defer ts.Restore()

	render := func() {
		terminal.ClearScreen()
		terminal.PrintPadding(5)
		fmt.Fprintln(terminal.Out)
		fmt.Fprintln(terminal.Out, terminal.ColYellow+"Les trois zones sont terminées."+terminal.ColReset)
		fmt.Fprintf(terminal.Out, "Un dernier adversaire t'attend : %s%s%s\n\n", terminal.ColRed, bossName, terminal.ColReset)
		fmt.Fprintln(terminal.Out, terminal.ColWhite+"L'affronter ? (Entrée = oui, X = non)"+terminal.ColReset)
	}

	for {
		terminal.CurrentScreen = render
		render()

		switch ts.ReadKey() {
		case terminal.KeyEnter:
			return true
		case terminal.KeyBack, terminal.KeyQuit:
			return false
		}
	}
}
