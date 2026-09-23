package combat

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	"ProjetRED/boss"
	"ProjetRED/character"
)

type bullet struct {
	x, y   float64
	dx, dy float64
	fuse   int
	homing float64
}

const (
	dodgeTick      = 55 * time.Millisecond
	dodgeDuration  = 95
	maxArenaWidth  = 64
	minArenaWidth  = 24
	minArenaHeight = 4

	tapStepX   = 1.6
	tapStepY   = 1.0
	baseSpeedX = 0.55
	baseSpeedY = 0.34
	maxSpeedX  = 1.7
	maxSpeedY  = 1.05
	accelTicks = 14
	holdExpiry = 36

	aspect     = 0.6
	homingCapX = 1.25
	homingCapY = 0.78
	safeRadius = 4.5
)

const (
	dirUp = iota
	dirDown
	dirLeft
	dirRight
)

type heldDir struct {
	down     bool
	since    int
	lastSeen int
}

func dirIndex(k Key) (int, bool) {
	switch k {
	case KeyUp:
		return dirUp, true
	case KeyDown:
		return dirDown, true
	case KeyLeft:
		return dirLeft, true
	case KeyRight:
		return dirRight, true
	}
	return 0, false
}

func dodgeLayout() (boxW, boxH, artW, artH int) {
	boxW, boxH, artW, artH = layout()
	grow := 8
	if artH-grow < minArtHeight {
		grow = artH - minArtHeight
	}
	if grow < 0 {
		grow = 0
	}
	return boxW, boxH + grow, artW, artH - grow
}

func RunDodge(b *boss.Boss, player *character.Character) (hit bool, quit bool) {
	if !canPollInput() {
		return true, false
	}

	boxW, boxH, artW, artH := dodgeLayout()
	arenaW := boxW
	if arenaW > maxArenaWidth {
		arenaW = maxArenaWidth
	}
	w, h := arenaW-2, boxH-2
	if w < minArenaWidth || h < minArenaHeight {
		return true, false
	}

	art := fitArt(b, artW, artH)
	artHeight := len(strings.Split(art, "\n"))

	clearScreen()
	printPadding(artHeight + boxH + 7)
	fmt.Println()
	fmt.Println(colYellow + "  " + b.Zone + colReset)
	DrawArt(boxW, art, b.Color)

	hx, hy := float64(w/2), float64(h/2)
	var held [4]heldDir
	var bullets []bullet
	atk := attack{style: b.Style}
	block := boxH + 5
	flushInput()

	for tick := 0; tick < dodgeDuration; tick++ {
		for _, ev := range pollKeyEvents() {
			d, ok := dirIndex(ev.key)
			if !ok {
				if ev.key == KeyQuit && ev.down {
					return false, true
				}
				continue
			}
			if !ev.down {
				held[d].down = false
				continue
			}
			if !held[d].down {
				held[d].down = true
				held[d].since = tick
				switch d {
				case dirUp:
					hy -= tapStepY
				case dirDown:
					hy += tapStepY
				case dirLeft:
					hx -= tapStepX
				case dirRight:
					hx += tapStepX
				}
			}
			held[d].lastSeen = tick
		}

		for d := range held {
			if held[d].down && tick-held[d].lastSeen > holdExpiry {
				held[d].down = false
			}
		}

		for d := range held {
			if !held[d].down {
				continue
			}
			ramp := float64(tick-held[d].since) / accelTicks
			if ramp > 1 {
				ramp = 1
			}
			switch d {
			case dirUp:
				hy -= baseSpeedY + (maxSpeedY-baseSpeedY)*ramp
			case dirDown:
				hy += baseSpeedY + (maxSpeedY-baseSpeedY)*ramp
			case dirLeft:
				hx -= baseSpeedX + (maxSpeedX-baseSpeedX)*ramp
			case dirRight:
				hx += baseSpeedX + (maxSpeedX-baseSpeedX)*ramp
			}
		}
		hx, hy = clampPos(hx, hy, w, h)

		cx, cy := cell(hx), cell(hy)
		bullets = append(bullets, atk.spawn(tick, w, h, hx, hy)...)

		scale := atk.speedScale(tick)
		var spawned []bullet
		alive := bullets[:0]
		struck := false
		for _, bl := range bullets {
			if bl.homing > 0 {
				bl = steer(bl, hx, hy)
			}
			if bl.fuse > 0 {
				bl.fuse--
				if bl.fuse == 0 {
					spawned = append(spawned, burst(bl)...)
					continue
				}
			}
			bl.x += bl.dx * scale
			bl.y += bl.dy * scale
			if bl.x < -2 || bl.x > float64(w)+2 || bl.y < -2 || bl.y > float64(h)+2 {
				continue
			}
			if cell(bl.x) == cx && cell(bl.y) == cy {
				struck = true
			}
			alive = append(alive, bl)
		}
		bullets = append(alive, spawned...)

		if tick > 0 {
			fmt.Printf("\033[%dA", block)
		}
		drawArena(boxW, w, h, cx, cy, bullets, atk.hidden(tick))
		fmt.Println()
		DrawEnemyBar(b)
		DrawStatBar(player)
		fmt.Println(colWhite + "\n" + atk.hint() + colReset)

		if struck {
			return true, false
		}
		time.Sleep(dodgeTick)
	}
	return false, false
}

func cell(v float64) int {
	return int(math.Floor(v + 0.5))
}

func clampPos(x, y float64, w, h int) (float64, float64) {
	if x < 0 {
		x = 0
	}
	if x > float64(w-1) {
		x = float64(w - 1)
	}
	if y < 0 {
		y = 0
	}
	if y > float64(h-1) {
		y = float64(h - 1)
	}
	return x, y
}

func steer(bl bullet, hx, hy float64) bullet {
	tx, ty := hx-bl.x, (hy-bl.y)/aspect
	n := math.Hypot(tx, ty)
	if n < 0.001 {
		return bl
	}
	bl.dx += tx / n * bl.homing
	bl.dy += ty / n * bl.homing * aspect
	if bl.dx > homingCapX {
		bl.dx = homingCapX
	}
	if bl.dx < -homingCapX {
		bl.dx = -homingCapX
	}
	if bl.dy > homingCapY {
		bl.dy = homingCapY
	}
	if bl.dy < -homingCapY {
		bl.dy = -homingCapY
	}
	return bl
}

func burst(b bullet) []bullet {
	dirs := [8][2]float64{
		{1, 0}, {-1, 0}, {0, 1}, {0, -1},
		{0.7, 0.7}, {-0.7, 0.7}, {0.7, -0.7}, {-0.7, -0.7},
	}
	out := make([]bullet, 0, len(dirs))
	for _, d := range dirs {
		out = append(out, bullet{x: b.x, y: b.y, dx: d[0] * 0.95, dy: d[1] * 0.95 * aspect})
	}
	return out
}

func awayFrom(hx, hy, w, h float64) (float64, float64) {
	for i := 0; i < 24; i++ {
		x := rand.Float64() * w
		y := rand.Float64() * h
		if math.Hypot(x-hx, (y-hy)/aspect) >= safeRadius {
			return x, y
		}
	}
	if hx > w/2 {
		return 0, rand.Float64() * h
	}
	return w - 1, rand.Float64() * h
}

func drawArena(boxW, w, h, hx, hy int, bullets []bullet, hidden bool) {
	pad := strings.Repeat(" ", (boxW-w-2)/2)
	grid := make([][]int, h)
	for y := range grid {
		grid[y] = make([]int, w)
	}
	if !hidden {
		for _, bl := range bullets {
			x, y := cell(bl.x), cell(bl.y)
			if x < 0 || x >= w || y < 0 || y >= h {
				continue
			}
			if bl.fuse > 0 || bl.homing > 0 {
				grid[y][x] = 2
				continue
			}
			if grid[y][x] == 0 {
				grid[y][x] = 1
			}
		}
	}

	fmt.Println(pad + colWhite + "┌" + strings.Repeat("─", w) + "┐" + colReset)
	for y := 0; y < h; y++ {
		var sb strings.Builder
		for x := 0; x < w; x++ {
			switch {
			case x == hx && y == hy:
				sb.WriteString(colRed + "❤" + colReset)
			case grid[y][x] == 2:
				sb.WriteString(colRed + "◉" + colReset)
			case grid[y][x] == 1:
				sb.WriteString(colYellow + "◆" + colReset)
			default:
				sb.WriteByte(' ')
			}
		}
		fmt.Println(pad + colWhite + "│" + colReset + sb.String() + colWhite + "│" + colReset)
	}
	fmt.Println(pad + colWhite + "└" + strings.Repeat("─", w) + "┘" + colReset)
}

type attack struct {
	style boss.AttackStyle
}

func (a attack) hint() string {
	switch a.style {
	case boss.AttackTimeStop:
		return "(← ↑ ↓ → pour esquiver — quand le temps s'arrête, place-toi avant la reprise)"
	case boss.AttackErase:
		return "(← ↑ ↓ → pour esquiver — le temps effacé rend les coups invisibles, souviens-toi)"
	case boss.AttackBombs:
		return "(← ↑ ↓ → pour esquiver — tout ce qu'il touche explose, ne reste pas à côté)"
	case boss.AttackAccelerate:
		return "(← ↑ ↓ → pour esquiver — le temps accélère, ça ne fera qu'empirer)"
	}
	return "(← ↑ ↓ → pour esquiver, maintenir pour aller plus vite)"
}

func (a attack) speedScale(tick int) float64 {
	switch a.style {
	case boss.AttackTimeStop:
		switch phase := tick % 34; {
		case phase >= 24 && phase < 30:
			return 0
		case phase >= 30:
			return 2.1
		}
		return 1

	case boss.AttackAccelerate:
		return 1 + 1.3*float64(tick)/dodgeDuration
	}
	return 1
}

func (a attack) hidden(tick int) bool {
	return a.style == boss.AttackErase && tick%30 >= 22
}

func (a attack) spawn(tick, w, h int, hx, hy float64) []bullet {
	fw, fh := float64(w), float64(h)

	switch a.style {
	case boss.AttackTimeStop:
		phase := tick % 34
		if phase == 24 {
			out := make([]bullet, 0, 7)
			for i := 0; i < 7; i++ {
				x, y := awayFrom(hx, hy, fw, fh)
				dx, dy := aim(x, y, hx, hy, 1.15)
				out = append(out, bullet{x: x, y: y, dx: dx, dy: dy})
			}
			return out
		}
		if phase < 24 && phase%4 == 0 {
			return []bullet{{x: fw, y: rand.Float64() * fh, dx: -1.15, dy: 0}}
		}
		return nil

	case boss.AttackErase:
		if tick%5 != 0 {
			return nil
		}
		x, y := edgePoint(fw, fh)
		dx, dy := aim(x, y, hx, hy, 1.0)
		return []bullet{{x: x, y: y, dx: dx, dy: dy}}

	case boss.AttackBombs:
		var out []bullet
		if tick%9 == 0 {
			x, y := awayFrom(hx, hy, fw, fh)
			out = append(out, bullet{x: x, y: y, dx: 0, dy: 0, fuse: 13})
		}
		if tick%40 == 20 {
			x, y := awayFrom(hx, hy, fw, fh)
			out = append(out, bullet{x: x, y: y, dx: 0, dy: 0, homing: 0.055})
		}
		return out

	case boss.AttackAccelerate:
		if tick%4 != 0 {
			return nil
		}
		out := []bullet{{x: rand.Float64() * fw, y: -1, dx: 0, dy: 0.5}}
		y := rand.Float64() * fh
		if rand.Intn(2) == 0 {
			out = append(out, bullet{x: -1, y: y, dx: 1.2, dy: 0})
		} else {
			out = append(out, bullet{x: fw, y: y, dx: -1.2, dy: 0})
		}
		return out
	}

	if tick%9 != 0 {
		return nil
	}
	return []bullet{{x: fw, y: rand.Float64() * fh, dx: -0.9, dy: 0}}
}

func aim(x, y, hx, hy, speed float64) (float64, float64) {
	tx, ty := hx-x, (hy-y)/aspect
	n := math.Hypot(tx, ty)
	if n < 0.001 {
		return -speed, 0
	}
	return tx / n * speed, ty / n * speed * aspect
}

func edgePoint(w, h float64) (float64, float64) {
	switch rand.Intn(4) {
	case 0:
		return -1, rand.Float64() * h
	case 1:
		return w, rand.Float64() * h
	case 2:
		return rand.Float64() * w, -1
	}
	return rand.Float64() * w, h
}
