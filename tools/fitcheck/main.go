package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"strconv"
	"strings"

	"ProjetRED/ascii"
)

func main() {
	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	w, h := 116, 36
	if len(os.Args) > 4 {
		w, _ = strconv.Atoi(os.Args[3])
		h, _ = strconv.Atoi(os.Args[4])
	}
	fitted := string(data)
	if !(len(os.Args) > 3 && os.Args[3] == "raw") {
		fitted = ascii.Fit(string(data), w, h)
	}
	fmt.Println(fitted)

	if len(os.Args) > 2 {
		writePreview(os.Args[2], fitted)
	}
}

const (
	asciiRamp = " .:-=+*#%@"
	blockRamp = " ░▒▓█"
)

func rampFor(art string) []rune {
	if strings.ContainsAny(art, blockRamp) {
		return []rune(blockRamp)
	}
	return []rune(asciiRamp)
}

var brailleBits = [4][2]uint8{
	{0x01, 0x08},
	{0x02, 0x10},
	{0x04, 0x20},
	{0x40, 0x80},
}

func writeDotPreview(path, art string) {
	lines := strings.Split(strings.TrimRight(art, "\n"), "\n")
	rows := len(lines)
	runeLines := make([][]rune, rows)
	cols := 0
	for i, l := range lines {
		runeLines[i] = []rune(l)
		if n := len(runeLines[i]); n > cols {
			cols = n
		}
	}
	const block = 4
	img := image.NewGray(image.Rect(0, 0, cols*2*block, rows*4*block))
	for cy, line := range runeLines {
		for cx := 0; cx < cols; cx++ {
			var mask uint8
			if cx < len(line) && line[cx] >= 0x2800 && line[cx] <= 0x28FF {
				mask = uint8(line[cx] - 0x2800)
			}
			for dy := 0; dy < 4; dy++ {
				for dx := 0; dx < 2; dx++ {
					g := color.Gray{Y: 0}
					if mask&brailleBits[dy][dx] != 0 {
						g = color.Gray{Y: 255}
					}
					px, py := (cx*2+dx)*block, (cy*4+dy)*block
					for by := 0; by < block; by++ {
						for bx := 0; bx < block; bx++ {
							img.SetGray(px+bx, py+by, g)
						}
					}
				}
			}
		}
	}
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	png.Encode(f, img)
}

func writePreview(path, art string) {
	for _, r := range art {
		if r >= 0x2800 && r <= 0x28FF {
			writeDotPreview(path, art)
			return
		}
	}
	ramp := rampFor(art)
	lines := strings.Split(art, "\n")
	runeLines := make([][]rune, len(lines))
	rows := len(lines)
	cols := 0
	for i, l := range lines {
		runeLines[i] = []rune(l)
		if n := len(runeLines[i]); n > cols {
			cols = n
		}
	}
	const block = 16
	img := image.NewGray(image.Rect(0, 0, cols*block, rows*block))
	for y, l := range runeLines {
		for x := 0; x < cols; x++ {
			ch := ' '
			if x < len(l) {
				ch = l[x]
			}
			v := 0.5
			for i, rr := range ramp {
				if rr == ch {
					v = float64(i) / float64(len(ramp)-1)
					break
				}
			}
			gray := color.Gray{Y: uint8(v * 255)}
			for by := 0; by < block; by++ {
				for bx := 0; bx < block; bx++ {
					img.SetGray(x*block+bx, y*block+by, gray)
				}
			}
		}
	}
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	png.Encode(f, img)
}
