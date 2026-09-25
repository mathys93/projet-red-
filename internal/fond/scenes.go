package fond

import "math"

func hash(a, b, c int) int {
	h := uint32(a)*374761393 + uint32(b)*668265263 + uint32(c)*2246822519 + 3266489917
	h = (h ^ (h >> 13)) * 1274126177
	h ^= h >> 16
	return int(h & 0x7fffffff)
}

func wrap(v, n int) int {
	if n <= 0 {
		return 0
	}
	return ((v % n) + n) % n
}

func Etoiles(t, w, h int, put func(int, int, rune, int)) {
	for y := range h {
		for x := range w {
			v := hash(x, y, 1)
			if v%41 != 0 {
				continue
			}
			switch phase := (t + v>>8) % 56; {
			case phase < 2:
				put(x, y, '✦', 231)
			case phase < 6:
				put(x, y, '+', 252)
			case v>>20%3 == 0:
				put(x, y, '·', 243)
			default:
				put(x, y, '·', 238)
			}
		}
	}
	cycle := t / 60
	p := t % 60
	if p < 14 {
		sx := hash(cycle, 7, 7) % max(w, 1)
		sy := hash(cycle, 8, 8) % max(h/3, 1)
		for k := range 4 {
			x, y := sx-(p-k)*3, sy+(p-k)
			ch, fg := '·', 245
			if k == 0 {
				ch, fg = '✦', 231
			}
			if p-k >= 0 {
				put(x, y, ch, fg)
			}
		}
	}
}

func Carte(t, w, h int, put func(int, int, rune, int)) {
	for y := range h {
		for x := range w {
			v := hash(x, y, 2)
			if v%23 != 0 {
				continue
			}
			switch v >> 8 % 5 {
			case 0:
				put(x, y, '▲', 101)
			case 1:
				put(x, y, '♣', 28)
			case 2:
				if (t/5+x)%2 == 0 {
					put(x, y, '≈', 31)
				} else {
					put(x, y, '~', 38)
				}
			case 3:
				put(x, y, '·', 58)
			}
		}
	}
	for i := range max(h/6, 2) {
		cy := hash(i, 3, 3) % max(h, 1)
		span := w + 16
		cx := wrap(hash(i, 4, 4)+t/(3+i%3), span) - 8
		for dx := range 8 {
			ch := '░'
			if dx >= 2 && dx < 6 {
				ch = '▒'
			}
			put(cx+dx, cy, ch, 250)
		}
		for dx := 2; dx < 6; dx++ {
			put(cx+dx, cy-1, '░', 248)
		}
	}
}

func Boutique(t, w, h int, put func(int, int, rune, int)) {
	n := max(w*h/70, 8)
	for i := range n {
		x0 := hash(i, 1, 5) % max(w, 1)
		speed := 1 + hash(i, 2, 5)%3
		y := h - 1 - wrap(t*speed/5+hash(i, 3, 5), h+4)
		x := x0 + int(math.Round(math.Sin(float64(t+i*7)/6)))
		switch i % 4 {
		case 0:
			put(x, y, '◆', 178)
		case 1:
			put(x, y, '◇', 222)
		case 2:
			put(x, y, '•', 214)
		default:
			put(x, y, '·', 136)
		}
	}
	for y := range h {
		for x := range w {
			v := hash(x, y, 6)
			if v%97 == 0 && (t+v>>9)%30 < 3 {
				put(x, y, '✦', 229)
			}
		}
	}
}

func Dio(t, w, h int, put func(int, int, rune, int)) {
	const cycle, run = 90, 64
	stopped := t%cycle >= run
	te := (t/cycle)*run + min(t%cycle, run)
	knife, clock, dust := 223, 178, 94
	if stopped {
		knife, clock, dust = 153, 111, 60
	}
	n := max(w*h/90, 8)
	for i := range n {
		x0 := hash(i, 1, 9) % max(w, 1)
		v := 2 + hash(i, 2, 9)%3
		y := wrap(te*v/4+hash(i, 3, 9), h+3) - 2
		x := x0 + y/3
		put(x, y, '†', knife)
		if !stopped {
			put(x, y-1, '╎', dust)
		}
	}
	for y := range h {
		for x := range w {
			v := hash(x, y, 10)
			if v%173 == 0 {
				put(x, y, []rune("◴◷◶◵")[(te/3+v>>8)%4], clock)
			} else if v%61 == 0 {
				put(x, y, '·', dust)
			}
		}
	}
}

func Diavolo(t, w, h int, put func(int, int, rune, int)) {
	const cycle = 56
	p := t % cycle
	if p >= 42 && p < 49 {
		return
	}
	te := t + (t/cycle)*30
	if p >= 49 {
		te += 30
	}
	n := max(w*h/80, 8)
	for i := range n {
		x0 := hash(i, 1, 11) % max(w, 1)
		y := h - 1 - wrap(te/2+hash(i, 2, 11), h+2)
		x := x0 + (h-y)/4
		put(x, y, '╱', 161)
		put(x+1, y-1, '╱', 125)
	}
	for y := range h {
		for x := range w {
			v := hash(x, y, 12)
			if v%67 == 0 {
				if (te+v>>8)%20 < 10 {
					put(x, y, '▚', 89)
				} else {
					put(x, y, '▞', 53+36)
				}
			}
		}
	}
	col := hash(te/9, 13, 13) % max(w, 1)
	if te%9 < 3 {
		for y := range h {
			if hash(col, y, te/9)%3 != 0 {
				put(col, y, '│', 197)
			}
		}
	}
}

func Kira(t, w, h int, put func(int, int, rune, int)) {
	const cycle = 44
	n := max(w*h/220, 5)
	for i := range n {
		off := hash(i, 1, 14) % cycle
		p := (t + off) % cycle
		round := (t + off) / cycle
		x := hash(i, round, 15) % max(w, 1)
		y := hash(i, round, 16) % max(h, 1)
		switch {
		case p < 20:
			ch, fg := '●', 211
			if p > 12 && p%2 == 0 {
				ch, fg = '○', 224
			}
			put(x, y, ch, fg)
		case p < 27:
			r := float64(p - 19)
			for a := 0; a < 16; a++ {
				ang := float64(a) * math.Pi / 8
				dx := int(math.Round(math.Cos(ang) * r * 2))
				dy := int(math.Round(math.Sin(ang) * r))
				ch, fg := '*', 214
				if r > 3 {
					ch, fg = '·', 208
				}
				put(x+dx, y+dy, ch, fg)
			}
		}
	}
	for y := range h {
		for x := range w {
			v := hash(x, y, 17)
			if v%89 == 0 {
				put(x, y, '♪', 30+int(v>>8)%2*6)
			} else if v%53 == 0 && (t+v>>7)%24 < 12 {
				put(x, y, '✶', 73)
			}
		}
	}
}

func Pucci(t, w, h int, put func(int, int, rune, int)) {
	const cycle = 170.0
	tc := float64(t % int(cycle))
	speed := 0.15 + 3.2*tc*tc/(cycle*cycle)
	dist := 0.15*tc + 3.2*tc*tc*tc/(3*cycle*cycle) + float64(t/int(cycle))*97
	n := max(w*h/50, 10)
	for i := range n {
		y := hash(i, 1, 18) % max(h, 1)
		lane := 0.6 + float64(hash(i, 2, 18)%5)/5
		x := w - 1 - wrap(int(dist*lane)+hash(i, 3, 18), w+8)
		switch {
		case speed*lane < 0.6:
			put(x, y, '·', 217)
		case speed*lane < 1.6:
			put(x, y, '-', 210)
			put(x+1, y, '·', 131)
		default:
			put(x, y, '━', 231)
			for k := 1; k <= min(int(speed*lane*2), 6); k++ {
				put(x+k, y, '─', 203)
			}
		}
	}
	for y := range h {
		for x := range w {
			if v := hash(x, y, 19); v%211 == 0 {
				put(x, y, '†', 95)
			}
		}
	}
}

func Iggy(t, w, h int, put func(int, int, rune, int)) {
	n := max(w*h/45, 10)
	for i := range n {
		v := 1 + hash(i, 1, 20)%3
		x := wrap(hash(i, 2, 20)+t*v, w+4) - 2
		y0 := hash(i, 3, 20) % max(h, 1)
		y := y0 + int(math.Round(math.Sin(float64(x+t)/7)*1.4))
		switch i % 4 {
		case 0:
			put(x, y, '░', 137)
		case 1:
			put(x, y, '∙', 180)
		default:
			put(x, y, '·', 223)
		}
	}
	for x := range w {
		base := h - 2 - int(math.Round((math.Sin(float64(x)/9)+1)*1.5))
		for y := max(base, 0); y < h; y++ {
			if y == base {
				put(x, y, '▁', 137)
			} else if (x+y)%3 == 0 {
				put(x, y, '░', 94)
			}
		}
	}
	dx := wrap(t/2, w+10) - 5
	for k := range 5 {
		put(dx+int(math.Round(math.Sin(float64(t+k*3)/2)*float64(k)/2)), h-4-k, '@', 180-k%2*43)
	}
}
