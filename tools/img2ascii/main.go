package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"

	_ "golang.org/x/image/webp"
)

const ramp = " .:-=+*#%@"

func main() {
	in := flag.String("in", "", "image source (png/jpg)")
	out := flag.String("out", "", "fichier .txt de sortie (art ASCII)")
	cols := flag.Int("cols", 60, "largeur en caractères")
	rows := flag.Int("rows", 0, "hauteur en caractères (0 = calculée automatiquement)")
	charAspect := flag.Float64("charAspect", 0.5, "ratio largeur/hauteur d'un caractère de terminal")
	edgeWeight := flag.Float64("edge", 0.5, "poids des contours (détection de bords) dans le rendu final")
	toneWeight := flag.Float64("tone", 1.0, "poids de la luminance/du volume (tons) dans le rendu final")
	dither := flag.Bool("dither", true, "diffusion d'erreur (Floyd-Steinberg) au lieu d'un arrondi simple par caractère")
	braille := flag.Bool("braille", true, "rendu en braille Unicode (2x4 points par caractère)")
	clipPct := flag.Float64("clip", 0.01, "pourcentage de valeurs extrêmes ignorées lors de l'étirement de contraste")
	contrastBoost := flag.Float64("contrast", 1.0, "renforcement du contraste autour du gris moyen avant tramage, 1 = inchangé")
	invert := flag.Bool("invert", false, "sujet sombre sur fond clair : inverse les tons pour que l'encre s'allume au lieu du fond")
	cropFlag := flag.String("crop", "", "recadrage en pixels \"x,y,w,h\" avant conversion")
	preview := flag.String("preview", "", "écrit une prévisualisation PNG agrandie du rendu")
	flag.Parse()

	if *in == "" || *out == "" {
		fmt.Fprintln(os.Stderr, "usage: img2ascii -in image.png -out assets/art/nom.txt -cols 60 [-crop x,y,w,h] [-rows N] [-edge 1.2] [-tone 0.8] [-contrast 2.0] [-preview debug.png]")
		os.Exit(2)
	}

	f, err := os.Open(*in)
	if err != nil {
		fatal(err)
	}
	img, _, err := image.Decode(f)
	f.Close()
	if err != nil {
		fatal(err)
	}

	bounds := img.Bounds()
	if *cropFlag != "" {
		r, err := parseCrop(*cropFlag)
		if err != nil {
			fatal(err)
		}
		r = r.Add(bounds.Min).Intersect(bounds)
		bounds = r
	}

	srcW := bounds.Dx()
	srcH := bounds.Dy()
	if srcW <= 0 || srcH <= 0 {
		fatal(fmt.Errorf("recadrage vide (image %dx%d)", bounds.Dx(), bounds.Dy()))
	}

	outCols := *cols
	outRows := *rows
	if outRows <= 0 {
		outRows = int(math.Round(float64(outCols) * (float64(srcH) / float64(srcW)) * (*charAspect)))
		if outRows < 1 {
			outRows = 1
		}
	}

	var lines []string
	if *braille {
		lines = renderBraille(img, bounds, outCols, outRows, *clipPct, *edgeWeight, *toneWeight, *contrastBoost, *invert)
	} else {
		lum, srcGridW, srcGridH := luminanceGrid(img, bounds)
		edges := sobelMagnitude(lum, srcGridW, srcGridH)

		edgeGrid := maxPoolGrid(edges, srcGridW, srcGridH, outCols, outRows)
		percentileStretch(edgeGrid, *clipPct)

		toneGrid := sampleGrid(img, bounds, outCols, outRows)
		percentileStretch(toneGrid, *clipPct)
		if *invert {
			for i, v := range toneGrid {
				toneGrid[i] = 1 - v
			}
		}

		grid := make([]float64, outCols*outRows)
		for i := range grid {
			grid[i] = clamp01(edgeGrid[i]**edgeWeight + toneGrid[i]**toneWeight)
		}

		if *dither {
			lines = ditherToLines(grid, outCols, outRows)
		} else {
			lines = make([]string, outRows)
			for y := 0; y < outRows; y++ {
				var sb strings.Builder
				for x := 0; x < outCols; x++ {
					v := grid[y*outCols+x]
					idx := int(math.Round(v * float64(len(ramp)-1)))
					sb.WriteByte(ramp[idx])
				}
				lines[y] = sb.String()
			}
		}
	}

	if err := os.WriteFile(*out, []byte(strings.Join(lines, "\n")+"\n"), 0o644); err != nil {
		fatal(err)
	}
	fmt.Printf("%s écrit (%dx%d caractères, source recadrée %dx%d)\n", *out, outCols, outRows, srcW, srcH)

	if *preview != "" {
		if *braille {
			if err := writeDotPreview(*preview, lines); err != nil {
				fatal(err)
			}
		} else if err := writePreview(*preview, lines); err != nil {
			fatal(err)
		}
		fmt.Printf("%s écrit (prévisualisation)\n", *preview)
	}
}

func writeDotPreview(path string, lines []string) error {
	rows := len(lines)
	cols := 0
	runeLines := make([][]rune, rows)
	for i, l := range lines {
		runeLines[i] = []rune(l)
		if n := len(runeLines[i]); n > cols {
			cols = n
		}
	}

	bitAt := [4][2]uint8{
		{0x01, 0x08},
		{0x02, 0x10},
		{0x04, 0x20},
		{0x40, 0x80},
	}

	const block = 6
	dotsW, dotsH := cols*2, rows*4
	img := image.NewGray(image.Rect(0, 0, dotsW*block, dotsH*block))
	for cy, runeLine := range runeLines {
		for cx := 0; cx < cols; cx++ {
			var mask uint8
			if cx < len(runeLine) && runeLine[cx] >= 0x2800 && runeLine[cx] <= 0x28FF {
				mask = uint8(runeLine[cx] - 0x2800)
			}
			for dy := 0; dy < 4; dy++ {
				for dx := 0; dx < 2; dx++ {
					on := mask&bitAt[dy][dx] != 0
					gray := color.Gray{Y: 0}
					if on {
						gray = color.Gray{Y: 255}
					}
					px, py := (cx*2+dx)*block, (cy*4+dy)*block
					for by := 0; by < block; by++ {
						for bx := 0; bx < block; bx++ {
							img.SetGray(px+bx, py+by, gray)
						}
					}
				}
			}
		}
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

func sampleGrid(img image.Image, bounds image.Rectangle, cols, rows int) []float64 {
	grid := make([]float64, cols*rows)
	srcW := bounds.Dx()
	srcH := bounds.Dy()

	for oy := 0; oy < rows; oy++ {
		y0 := bounds.Min.Y + oy*srcH/rows
		y1 := bounds.Min.Y + (oy+1)*srcH/rows
		if y1 <= y0 {
			y1 = y0 + 1
		}
		for ox := 0; ox < cols; ox++ {
			x0 := bounds.Min.X + ox*srcW/cols
			x1 := bounds.Min.X + (ox+1)*srcW/cols
			if x1 <= x0 {
				x1 = x0 + 1
			}

			var sum float64
			var n int
			for y := y0; y < y1; y++ {
				for x := x0; x < x1; x++ {
					sum += luminance(img.At(x, y))
					n++
				}
			}
			if n == 0 {
				n = 1
			}
			grid[oy*cols+ox] = sum / float64(n)
		}
	}
	return grid
}

func luminance(c color.Color) float64 {
	r, g, b, a := c.RGBA()
	if a == 0 {
		return 0
	}
	rf := float64(r) / 65535
	gf := float64(g) / 65535
	bf := float64(b) / 65535
	af := float64(a) / 65535
	lum := 0.299*rf + 0.587*gf + 0.114*bf
	return lum / af
}

func luminanceGrid(img image.Image, bounds image.Rectangle) (grid []float64, w, h int) {
	w, h = bounds.Dx(), bounds.Dy()
	grid = make([]float64, w*h)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			grid[y*w+x] = luminance(img.At(bounds.Min.X+x, bounds.Min.Y+y))
		}
	}
	return grid, w, h
}

func sobelMagnitude(lum []float64, w, h int) []float64 {
	at := func(x, y int) float64 {
		if x < 0 {
			x = 0
		}
		if x >= w {
			x = w - 1
		}
		if y < 0 {
			y = 0
		}
		if y >= h {
			y = h - 1
		}
		return lum[y*w+x]
	}

	out := make([]float64, w*h)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			gx := -at(x-1, y-1) + at(x+1, y-1) +
				-2*at(x-1, y) + 2*at(x+1, y) +
				-at(x-1, y+1) + at(x+1, y+1)
			gy := -at(x-1, y-1) - 2*at(x, y-1) - at(x+1, y-1) +
				at(x-1, y+1) + 2*at(x, y+1) + at(x+1, y+1)
			out[y*w+x] = math.Hypot(gx, gy)
		}
	}
	return out
}

func maxPoolGrid(src []float64, srcW, srcH, cols, rows int) []float64 {
	out := make([]float64, cols*rows)
	for oy := 0; oy < rows; oy++ {
		y0 := oy * srcH / rows
		y1 := (oy + 1) * srcH / rows
		if y1 <= y0 {
			y1 = y0 + 1
		}
		for ox := 0; ox < cols; ox++ {
			x0 := ox * srcW / cols
			x1 := (ox + 1) * srcW / cols
			if x1 <= x0 {
				x1 = x0 + 1
			}

			max := 0.0
			for y := y0; y < y1 && y < srcH; y++ {
				for x := x0; x < x1 && x < srcW; x++ {
					if v := src[y*srcW+x]; v > max {
						max = v
					}
				}
			}
			out[oy*cols+ox] = max
		}
	}
	return out
}

func percentileStretch(grid []float64, clipPct float64) {
	sorted := append([]float64(nil), grid...)
	sort.Float64s(sorted)
	n := len(sorted)
	lo := sorted[int(clipPct*float64(n-1))]
	hi := sorted[int((1-clipPct)*float64(n-1))]
	span := hi - lo
	if span <= 0 {
		span = 1
	}
	for i, v := range grid {
		grid[i] = clamp01((v - lo) / span)
	}
}

func renderBraille(img image.Image, bounds image.Rectangle, cols, rows int, clipPct, edgeWeight, toneWeight, contrastBoost float64, invert bool) []string {
	dotsW, dotsH := cols*2, rows*4

	tone := sampleGrid(img, bounds, dotsW, dotsH)
	percentileStretch(tone, clipPct)
	if invert {
		for i, v := range tone {
			tone[i] = 1 - v
		}
	}

	lum, srcW, srcH := luminanceGrid(img, bounds)
	edges := sobelMagnitude(lum, srcW, srcH)
	edgeGrid := maxPoolGrid(edges, srcW, srcH, dotsW, dotsH)
	percentileStretch(edgeGrid, clipPct)

	dots := make([]float64, dotsW*dotsH)
	for i := range dots {
		v := clamp01(edgeGrid[i]*edgeWeight + tone[i]*toneWeight)
		v = clamp01(0.5 + (v-0.5)*contrastBoost)
		dots[i] = v
	}
	ditherBinary(dots, dotsW, dotsH)

	bit := [4][2]uint8{
		{0x01, 0x08},
		{0x02, 0x10},
		{0x04, 0x20},
		{0x40, 0x80},
	}

	lines := make([]string, rows)
	for cy := 0; cy < rows; cy++ {
		runes := make([]rune, cols)
		for cx := 0; cx < cols; cx++ {
			var mask rune
			for dy := 0; dy < 4; dy++ {
				for dx := 0; dx < 2; dx++ {
					px, py := cx*2+dx, cy*4+dy
					if dots[py*dotsW+px] >= 0.5 {
						mask |= rune(bit[dy][dx])
					}
				}
			}
			runes[cx] = 0x2800 + mask
		}
		lines[cy] = string(runes)
	}
	return lines
}

func ditherBinary(grid []float64, cols, rows int) {
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			old := clamp01(grid[y*cols+x])
			newV := 0.0
			if old >= 0.5 {
				newV = 1.0
			}
			grid[y*cols+x] = newV
			errv := old - newV

			distribute := func(dx, dy int, frac float64) {
				nx, ny := x+dx, y+dy
				if nx < 0 || nx >= cols || ny < 0 || ny >= rows {
					return
				}
				grid[ny*cols+nx] += errv * frac
			}
			distribute(1, 0, 7.0/16)
			distribute(-1, 1, 3.0/16)
			distribute(0, 1, 5.0/16)
			distribute(1, 1, 1.0/16)
		}
	}
}

func ditherToLines(grid []float64, cols, rows int) []string {
	work := append([]float64(nil), grid...)
	lines := make([]string, rows)
	for y := 0; y < rows; y++ {
		line := make([]byte, cols)
		for x := 0; x < cols; x++ {
			old := clamp01(work[y*cols+x])
			idx := int(math.Round(old * float64(len(ramp)-1)))
			line[x] = ramp[idx]
			quant := float64(idx) / float64(len(ramp)-1)
			errv := old - quant

			distribute := func(dx, dy int, frac float64) {
				nx, ny := x+dx, y+dy
				if nx < 0 || nx >= cols || ny < 0 || ny >= rows {
					return
				}
				work[ny*cols+nx] += errv * frac
			}
			distribute(1, 0, 7.0/16)
			distribute(-1, 1, 3.0/16)
			distribute(0, 1, 5.0/16)
			distribute(1, 1, 1.0/16)
		}
		lines[y] = string(line)
	}
	return lines
}

func popcount(b uint8) int {
	n := 0
	for b != 0 {
		n += int(b & 1)
		b >>= 1
	}
	return n
}

func clamp01(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

func writePreview(path string, lines []string) error {
	runeLines := make([][]rune, len(lines))
	rows := len(lines)
	cols := 0
	for i, l := range lines {
		runeLines[i] = []rune(l)
		if n := len(runeLines[i]); n > cols {
			cols = n
		}
	}
	const block = 12
	img := image.NewGray(image.Rect(0, 0, cols*block, rows*block))
	for y, runeLine := range runeLines {
		for x := 0; x < cols; x++ {
			ch := ' '
			if x < len(runeLine) {
				ch = runeLine[x]
			}
			v := 0.5
			switch {
			case ch >= 0x2800 && ch <= 0x28FF:
				v = float64(popcount(uint8(ch-0x2800))) / 8
			default:
				if i := strings.IndexRune(ramp, ch); i >= 0 {
					v = float64(i) / float64(len([]rune(ramp))-1)
				}
			}
			gray := color.Gray{Y: uint8(math.Round(v * 255))}
			for by := 0; by < block; by++ {
				for bx := 0; bx < block; bx++ {
					img.SetGray(x*block+bx, y*block+by, gray)
				}
			}
		}
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

func parseCrop(s string) (image.Rectangle, error) {
	parts := strings.Split(s, ",")
	if len(parts) != 4 {
		return image.Rectangle{}, fmt.Errorf("crop invalide %q, attendu x,y,w,h", s)
	}
	vals := make([]int, 4)
	for i, p := range parts {
		n, err := strconv.Atoi(strings.TrimSpace(p))
		if err != nil {
			return image.Rectangle{}, fmt.Errorf("crop invalide %q : %w", s, err)
		}
		vals[i] = n
	}
	return image.Rect(vals[0], vals[1], vals[0]+vals[2], vals[1]+vals[3]), nil
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "erreur:", err)
	os.Exit(1)
}
