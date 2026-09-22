// Package ascii embarque les artworks ASCII des boss (dossier ascii/*.txt)
// directement dans le binaire (go:embed), et fournit un utilitaire pour les
// redimensionner afin qu'ils rentrent dans la boîte de combat, quelle que
// soit leur taille d'origine.
package ascii

import (
	_ "embed"
	"strings"
)

//go:embed dio.txt
var dioArt string

//go:embed diavolo.txt
var diavoloArt string

//go:embed kira.txt
var kiraArt string

//go:embed disque_de_pucci.txt
var pucciArt string

//go:embed shop.txt
var shopArt string

// registry associe une clé simple à chaque artwork, pour que le package
// boss puisse aller chercher l'art d'un boss par son nom (le combat package
// fait de même pour l'écran de boutique, via la clé "shop").
var registry = map[string]string{
	"dio":     dioArt,
	"diavolo": diavoloArt,
	"kira":    kiraArt,
	"pucci":   pucciArt,
	"shop":    shopArt,
}

// Get renvoie l'artwork associé à une clé ("dio", "diavolo", "kira", "pucci").
func Get(key string) (string, bool) {
	art, ok := registry[key]
	return art, ok
}

// Fit redimensionne un artwork (par échantillonnage) pour qu'il tienne dans
// une largeur/hauteur maximum, sans dépasser la boîte de combat. Certains
// artworks (comme disque_de_pucci.txt, ~100 colonnes sur 55 lignes) sont
// bien plus grands que la fenêtre de combat (~46x14) : sans ça, l'affichage
// undertale-like serait complètement disloqué.
func Fit(art string, maxW, maxH int) string {
	art = strings.TrimRight(art, "\n")
	lines := strings.Split(art, "\n")
	h := len(lines)
	w := 0
	for _, l := range lines {
		if n := len([]rune(l)); n > w {
			w = n
		}
	}
	if h == 0 || w == 0 {
		return art
	}
	if h <= maxH && w <= maxW {
		return art
	}

	outH := h
	if outH > maxH {
		outH = maxH
	}
	outW := w
	if outW > maxW {
		outW = maxW
	}

	out := make([]string, outH)
	for y := 0; y < outH; y++ {
		srcY := y * h / outH
		row := []rune(lines[srcY])
		line := make([]rune, outW)
		for x := 0; x < outW; x++ {
			srcX := x * w / outW
			if srcX < len(row) {
				line[x] = row[srcX]
			} else {
				line[x] = ' '
			}
		}
		out[y] = string(line)
	}
	return strings.Join(out, "\n")
}
