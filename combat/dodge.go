package combat

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"ProjetRED/boss"
	"ProjetRED/character"
)

type bullet struct {
	x, y   float64
	dx, dy float64
}

const (
	dodgeTick      = 55 * time.Millisecond
	dodgeDuration  = 95
	maxArenaWidth  = 64
	minArenaWidth  = 24
	minArenaHeight = 4
)

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

	hx, hy := w/2, h/2
	var bullets []bullet
	block := boxH + 5
	flushInput()

	for tick := 0; tick < dodgeDuration; tick++ {
		for {
			k, ok := pollKey()
			if !ok {
				break
			}
			switch k {
			case KeyLeft:
				if hx > 0 {
					hx--
				}
			case KeyRight:
				if hx < w-1 {
					hx++
				}
			case KeyUp:
				if hy > 0 {
					hy--
				}
			case KeyDown:
				if hy < h-1 {
					hy++
				}
			case KeyQuit:
				return false, true
			}
		}

		bullets = append(bullets, spawnBullets(b.Style, tick, w, h)...)

		alive := bullets[:0]
		struck := false
		for _, bl := range bullets {
			bl.x += bl.dx
			bl.y += bl.dy
			if bl.x < -2 || bl.x > float64(w)+2 || bl.y < -2 || bl.y > float64(h)+2 {
				continue
			}
			if int(bl.x+0.5) == hx && int(bl.y+0.5) == hy {
				struck = true
			}
			alive = append(alive, bl)
		}
		bullets = alive

		if tick > 0 {
			fmt.Printf("\033[%dA", block)
		}
		drawArena(boxW, w, h, hx, hy, bullets)
		fmt.Println()
		DrawEnemyBar(b)
		DrawStatBar(player)
		fmt.Println(colWhite + "\n(← ↑ ↓ → pour esquiver)" + colReset)

		if struck {
			return true, false
		}
		time.Sleep(dodgeTick)
	}
	return false, false
}

func drawArena(boxW, w, h, hx, hy int, bullets []bullet) {
	pad := strings.Repeat(" ", (boxW-w-2)/2)
	grid := make([][]bool, h)
	for y := range grid {
		grid[y] = make([]bool, w)
	}
	for _, bl := range bullets {
		x, y := int(bl.x+0.5), int(bl.y+0.5)
		if x >= 0 && x < w && y >= 0 && y < h {
			grid[y][x] = true
		}
	}

	fmt.Println(pad + colWhite + "┌" + strings.Repeat("─", w) + "┐" + colReset)
	for y := 0; y < h; y++ {
		var sb strings.Builder
		for x := 0; x < w; x++ {
			switch {
			case x == hx && y == hy:
				sb.WriteString(colRed + "❤" + colReset)
			case grid[y][x]:
				sb.WriteString(colYellow + "◆" + colReset)
			default:
				sb.WriteByte(' ')
			}
		}
		fmt.Println(pad + colWhite + "│" + colReset + sb.String() + colWhite + "│" + colReset)
	}
	fmt.Println(pad + colWhite + "└" + strings.Repeat("─", w) + "┘" + colReset)
}

func spawnBullets(style boss.AttackStyle, tick, w, h int) []bullet {
	switch style {
	case boss.AttackConverge:
		if tick%7 != 0 {
			return nil
		}
		y := float64(rand.Intn(h))
		if rand.Intn(2) == 0 {
			return []bullet{{x: 0, y: y, dx: 1.1, dy: 0}}
		}
		return []bullet{{x: float64(w - 1), y: y, dx: -1.1, dy: 0}}

	case boss.AttackRain:
		if tick%3 != 0 {
			return nil
		}
		return []bullet{{x: float64(rand.Intn(w)), y: 0, dx: 0, dy: 0.45}}

	case boss.AttackStorm:
		if tick%4 != 0 {
			return nil
		}
		out := []bullet{{x: float64(rand.Intn(w)), y: 0, dx: 0, dy: 0.5}}
		y := float64(rand.Intn(h))
		if rand.Intn(2) == 0 {
			out = append(out, bullet{x: 0, y: y, dx: 1.2, dy: 0})
		} else {
			out = append(out, bullet{x: float64(w - 1), y: y, dx: -1.2, dy: 0})
		}
		return out

	default:
		if tick%14 != 0 {
			return nil
		}
		gap := rand.Intn(h)
		var out []bullet
		for y := 0; y < h; y++ {
			if y == gap || y == gap+1 {
				continue
			}
			out = append(out, bullet{x: float64(w - 1), y: float64(y), dx: -0.8, dy: 0})
		}
		return out
	}
}
