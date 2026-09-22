package ascii

import (
	_ "embed"
	"math"
	"sort"
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

var registry = map[string]string{
	"dio":     dioArt,
	"diavolo": diavoloArt,
	"kira":    kiraArt,
	"pucci":   pucciArt,
	"shop":    shopArt,
}

func Get(key string) (string, bool) {
	art, ok := registry[key]
	return art, ok
}

const (
	asciiRamp = " .:-=+*#%@"
	blockRamp = " ░▒▓█"
)

var brailleBits = [4][2]uint8{
	{0x01, 0x08},
	{0x02, 0x10},
	{0x04, 0x20},
	{0x40, 0x80},
}

func rampFor(art string) []rune {
	if strings.ContainsAny(art, blockRamp) {
		return []rune(blockRamp)
	}
	return []rune(asciiRamp)
}

func isBraille(r rune) bool {
	return r >= 0x2800 && r <= 0x28FF
}

func hasBraille(art string) bool {
	for _, r := range art {
		if isBraille(r) {
			return true
		}
	}
	return false
}

func densityOf(ramp []rune, r rune) float64 {
	for i, rr := range ramp {
		if rr == r {
			return float64(i) / float64(len(ramp)-1)
		}
	}
	return 0.5
}

func charFor(ramp []rune, v float64) rune {
	i := int(v*float64(len(ramp)-1) + 0.5)
	if i < 0 {
		i = 0
	}
	if i >= len(ramp) {
		i = len(ramp) - 1
	}
	return ramp[i]
}

func splitArt(art string) (rows [][]rune, w, h int) {
	art = strings.TrimRight(art, "\n")
	lines := strings.Split(art, "\n")
	h = len(lines)
	rows = make([][]rune, h)
	for y, l := range lines {
		rows[y] = []rune(l)
		if n := len(rows[y]); n > w {
			w = n
		}
	}
	return rows, w, h
}

func Fit(art string, maxW, maxH int) string {
	rows, w, h := splitArt(art)
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

	if hasBraille(art) {
		return fitBraille(rows, w, h, outW, outH)
	}
	return fitRamp(rows, w, h, outW, outH, rampFor(art))
}

func fitRamp(rows [][]rune, w, h, outW, outH int, ramp []rune) string {
	useAverage := string(ramp) == blockRamp

	grid := make([]float64, outW*outH)
	for y := 0; y < outH; y++ {
		y0 := y * h / outH
		y1 := (y + 1) * h / outH
		if y1 <= y0 {
			y1 = y0 + 1
		}
		for x := 0; x < outW; x++ {
			x0 := x * w / outW
			x1 := (x + 1) * w / outW
			if x1 <= x0 {
				x1 = x0 + 1
			}

			var sum float64
			var n int
			best := 0.5
			bestDist := -1.0
			for sy := y0; sy < y1 && sy < h; sy++ {
				row := rows[sy]
				for sx := x0; sx < x1; sx++ {
					if sx >= len(row) {
						continue
					}
					v := densityOf(ramp, row[sx])
					sum += v
					n++
					if d := math.Abs(v - 0.5); d > bestDist {
						bestDist = d
						best = v
					}
				}
			}
			if useAverage {
				if n == 0 {
					n = 1
				}
				grid[y*outW+x] = sum / float64(n)
			} else {
				grid[y*outW+x] = best
			}
		}
	}

	percentileStretch(grid, 0.02)

	out := make([]string, outH)
	for y := 0; y < outH; y++ {
		line := make([]rune, outW)
		for x := 0; x < outW; x++ {
			line[x] = charFor(ramp, grid[y*outW+x])
		}
		out[y] = string(line)
	}
	return strings.Join(out, "\n")
}

func fitBraille(rows [][]rune, w, h, outW, outH int) string {
	srcW, srcH := w*2, h*4
	src := make([]float64, srcW*srcH)
	for cy, line := range rows {
		for cx := 0; cx < w; cx++ {
			var mask uint8
			if cx < len(line) && isBraille(line[cx]) {
				mask = uint8(line[cx] - 0x2800)
			}
			for dy := 0; dy < 4; dy++ {
				for dx := 0; dx < 2; dx++ {
					if mask&brailleBits[dy][dx] != 0 {
						src[(cy*4+dy)*srcW+cx*2+dx] = 1
					}
				}
			}
		}
	}

	dstW, dstH := outW*2, outH*4
	grid := make([]float64, dstW*dstH)
	for y := 0; y < dstH; y++ {
		y0 := y * srcH / dstH
		y1 := (y + 1) * srcH / dstH
		if y1 <= y0 {
			y1 = y0 + 1
		}
		for x := 0; x < dstW; x++ {
			x0 := x * srcW / dstW
			x1 := (x + 1) * srcW / dstW
			if x1 <= x0 {
				x1 = x0 + 1
			}

			var sum float64
			var n int
			for sy := y0; sy < y1 && sy < srcH; sy++ {
				for sx := x0; sx < x1 && sx < srcW; sx++ {
					sum += src[sy*srcW+sx]
					n++
				}
			}
			if n == 0 {
				n = 1
			}
			grid[y*dstW+x] = sum / float64(n)
		}
	}

	ditherBinary(grid, dstW, dstH)

	out := make([]string, outH)
	for cy := 0; cy < outH; cy++ {
		line := make([]rune, outW)
		for cx := 0; cx < outW; cx++ {
			var mask rune
			for dy := 0; dy < 4; dy++ {
				for dx := 0; dx < 2; dx++ {
					if grid[(cy*4+dy)*dstW+cx*2+dx] >= 0.5 {
						mask |= rune(brailleBits[dy][dx])
					}
				}
			}
			line[cx] = 0x2800 + mask
		}
		out[cy] = string(line)
	}
	return strings.Join(out, "\n")
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

func clamp01(v float64) float64 {
	return math.Max(0, math.Min(1, v))
}
