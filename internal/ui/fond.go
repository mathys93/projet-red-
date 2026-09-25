package ui

import (
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

const fondTick = 110 * time.Millisecond

type glyph struct {
	ch rune
	fg int
}

type screen struct {
	mu          sync.Mutex
	w           io.Writer
	cols, rows  int
	row, col    int
	wrapPending bool
	started     bool
	skip        int
	bgDefault   bool
	opaque      []bool
	shown       map[int]glyph
	esc         []byte
	scene       Scene
	sceneName   string
	tick        int
	anim        sync.Once
}

var Out = &screen{w: os.Stdout, bgDefault: true, shown: map[int]glyph{}}

func (s *screen) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var out []byte
	for i := 0; i < len(p); {
		if s.esc != nil {
			b := p[i]
			i++
			s.esc = append(s.esc, b)
			if s.escDone() {
				out = append(out, s.esc...)
				s.applyEscape(string(s.esc))
				s.esc = nil
			}
			continue
		}
		r, size := utf8.DecodeRune(p[i:])
		chunk := p[i : i+size]
		i += size
		switch {
		case r == '\033':
			s.esc = []byte{'\033'}
		case r == '\n':
			s.skip = 0
			s.newline()
			out = append(out, '\n')
		case r == '\r':
			s.skip = 0
			s.col = 0
			s.wrapPending = false
			out = append(out, '\r')
		case r == '\b':
			if s.col > 0 {
				s.col--
			}
			s.wrapPending = false
			out = append(out, '\b')
		case r == ' ' && s.bgDefault && !s.started:
			s.skip++
			s.col++
		default:
			if s.skip > 0 {
				out = append(out, "\033["+strconv.Itoa(s.skip)+"C"...)
				s.skip = 0
			}
			if s.wrapPending {
				s.newline()
			}
			wd := runeWidth(r)
			for k := 0; k < wd; k++ {
				s.markOpaque(s.row, s.col+k)
			}
			s.col += wd
			if s.col >= s.cols {
				s.col = s.cols - 1
				s.wrapPending = true
			}
			s.started = true
			out = append(out, chunk...)
		}
	}
	_, err := s.w.Write(out)
	return len(p), err
}

func (s *screen) escDone() bool {
	e := s.esc
	if len(e) < 2 {
		return false
	}
	switch e[1] {
	case '[':
		if len(e) < 3 {
			return false
		}
		last := e[len(e)-1]
		return last >= 0x40 && last <= 0x7e
	case ']':
		return e[len(e)-1] == '\a' || (len(e) >= 2 && e[len(e)-2] == '\033' && e[len(e)-1] == '\\')
	}
	return true
}

func (s *screen) applyEscape(seq string) {
	if len(seq) < 3 || seq[1] != '[' {
		return
	}
	body := seq[2 : len(seq)-1]
	final := seq[len(seq)-1]
	if strings.HasPrefix(body, "?") {
		return
	}
	nums := func(def int) []int {
		var v []int
		for _, part := range strings.Split(body, ";") {
			n, err := strconv.Atoi(part)
			if err != nil {
				n = def
			}
			v = append(v, n)
		}
		return v
	}
	switch final {
	case 'm':
		s.applySGR(nums(0))
	case 'H', 'f':
		v := nums(1)
		r, c := 1, 1
		if len(v) > 0 && v[0] > 0 {
			r = v[0]
		}
		if len(v) > 1 && v[1] > 0 {
			c = v[1]
		}
		s.moveTo(r-1, c-1)
	case 'A':
		s.moveTo(s.row-max(nums(1)[0], 1), s.col)
	case 'B':
		s.moveTo(s.row+max(nums(1)[0], 1), s.col)
	case 'C':
		s.moveTo(s.row, s.col+max(nums(1)[0], 1))
	case 'D':
		s.moveTo(s.row, s.col-max(nums(1)[0], 1))
	case 'J':
		if body == "2" {
			s.resetGrid()
		}
	}
}

func (s *screen) applySGR(v []int) {
	for i := 0; i < len(v); i++ {
		switch v[i] {
		case 0, 49:
			s.bgDefault = true
		case 38:
			if i+1 < len(v) && v[i+1] == 5 {
				i += 2
			} else if i+1 < len(v) && v[i+1] == 2 {
				i += 4
			}
		case 48:
			if i+2 < len(v) && v[i+1] == 5 {
				s.bgDefault = v[i+2] == bgNeutralCode
				i += 2
			} else if i+1 < len(v) && v[i+1] == 2 {
				s.bgDefault = false
				i += 4
			}
		default:
			if v[i] >= 40 && v[i] <= 47 || v[i] >= 100 && v[i] <= 107 {
				s.bgDefault = false
			}
		}
	}
}

func (s *screen) moveTo(r, c int) {
	s.row = min(max(r, 0), max(s.rows-1, 0))
	s.col = min(max(c, 0), max(s.cols-1, 0))
	s.wrapPending = false
	s.started = false
	s.skip = 0
}

func (s *screen) newline() {
	s.col = 0
	s.wrapPending = false
	s.started = false
	if s.row < s.rows-1 {
		s.row++
		return
	}
	copy(s.opaque, s.opaque[s.cols:])
	for i := len(s.opaque) - s.cols; i < len(s.opaque); i++ {
		s.opaque[i] = false
	}
	shifted := map[int]glyph{}
	for k, g := range s.shown {
		if k >= s.cols {
			shifted[k-s.cols] = g
		}
	}
	s.shown = shifted
}

func (s *screen) markOpaque(r, c int) {
	if r < 0 || r >= s.rows || c < 0 || c >= s.cols {
		return
	}
	k := r*s.cols + c
	s.opaque[k] = true
	delete(s.shown, k)
}

func (s *screen) resetGrid() {
	s.cols, s.rows = TerminalSize()
	s.opaque = make([]bool, s.cols*s.rows)
	s.shown = map[int]glyph{}
	s.row, s.col = 0, 0
	s.wrapPending = false
	s.started = false
	s.skip = 0
}

func runeWidth(r rune) int {
	switch {
	case r >= 0x1100 && r <= 0x115f,
		r >= 0x2e80 && r <= 0xa4cf,
		r >= 0xac00 && r <= 0xd7a3,
		r >= 0xf900 && r <= 0xfaff,
		r >= 0xff00 && r <= 0xff60,
		r >= 0xffe0 && r <= 0xffe6:
		return 2
	}
	return 1
}

func (s *screen) clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.resetGrid()
	s.bgDefault = true
	io.WriteString(s.w, ColReset+"\033[2J\033[H\033[?25l")
	s.paint()
	io.WriteString(s.w, "\033[H")
}

func (s *screen) frame() map[int]glyph {
	next := map[int]glyph{}
	if s.scene == nil || s.cols < 2 || s.rows < 1 {
		return next
	}
	w, h := s.cols-1, s.rows
	s.scene(s.tick, w, h, func(x, y int, ch rune, fg int) {
		if x < 0 || x >= w || y < 0 || y >= h {
			return
		}
		next[y*s.cols+x] = glyph{ch, fg}
	})
	return next
}

func (s *screen) paint() {
	next := s.frame()
	var b strings.Builder
	for k, g := range s.shown {
		if n, ok := next[k]; ok && n == g {
			continue
		}
		if s.opaque[k] {
			delete(s.shown, k)
			continue
		}
		b.WriteString(CursorAt(k/s.cols+1, k%s.cols+1))
		b.WriteString(ColReset)
		b.WriteByte(' ')
		delete(s.shown, k)
	}
	for k, g := range next {
		if s.opaque[k] {
			continue
		}
		if old, ok := s.shown[k]; ok && old == g {
			continue
		}
		b.WriteString(CursorAt(k/s.cols+1, k%s.cols+1))
		b.WriteString(ColReset)
		b.WriteString("\033[38;5;")
		b.WriteString(strconv.Itoa(g.fg))
		b.WriteByte('m')
		b.WriteRune(g.ch)
		s.shown[k] = g
	}
	if b.Len() == 0 {
		return
	}
	io.WriteString(s.w, "\0337"+b.String()+"\0338")
}

func (s *screen) animate() {
	for {
		time.Sleep(fondTick)
		s.mu.Lock()
		s.tick++
		if s.scene != nil && s.opaque != nil && s.esc == nil {
			s.paint()
		}
		s.mu.Unlock()
	}
}

func SetScene(name string) string {
	Out.mu.Lock()
	defer Out.mu.Unlock()
	prev := Out.sceneName
	Out.sceneName = name
	Out.scene = scenes[name]
	if CanPollInput() {
		Out.anim.Do(func() { go Out.animate() })
	}
	return prev
}
