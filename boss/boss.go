package boss

import (
	"fmt"

	"ProjetRED/ascii"
	"ProjetRED/character"
)

const (
	ColorWhite  = "\033[97m"
	ColorRed    = "\033[91m"
	ColorYellow = "\033[93m"
	ColorPurple = "\033[95m"
	ColorCyan   = "\033[96m"
)

type Pattern interface {
	Act(turn int, self *Boss, player *character.Character) string

	ActOptions(self *Boss, player *character.Character) []ActOption
}

type ActOption struct {
	Label       string
	Description string
	Resolve     func(self *Boss, player *character.Character) string
}

type AttackStyle int

const (
	AttackStroll AttackStyle = iota
	AttackTimeStop
	AttackErase
	AttackBombs
	AttackAccelerate
)

type Boss struct {
	*character.Character
	Art            string
	Zone           string
	Pattern        Pattern
	Color          string
	XPReward       int
	FragmentReward int
	Style          AttackStyle
}

func New(name string, level, ap, dp, lp, xpReward int, artKey string, pattern Pattern, color string) *Boss {
	raw, ok := ascii.Get(artKey)
	if !ok {
		raw = fmt.Sprintf("(art manquant: %q)", artKey)
	}
	if pattern == nil {
		pattern = DefaultPattern{}
	}
	if color == "" {
		color = ColorWhite
	}
	return &Boss{
		Character:      character.New(name, level, ap, dp, lp),
		Art:            raw,
		Pattern:        pattern,
		Color:          color,
		XPReward:       xpReward,
		FragmentReward: 15,
	}
}

func (b *Boss) Turn(turn int, player *character.Character) string {
	return b.Pattern.Act(turn, b, player)
}

func (b *Boss) ActOptions(player *character.Character) []ActOption {
	return b.Pattern.ActOptions(b, player)
}

func NewDio() *Boss {
	b := New("DIO", 5, 10, 4, 70, 50, "dio", DioPattern{}, ColorYellow)
	b.Style = AttackTimeStop
	return b
}

func NewDiavolo() *Boss {
	b := New("Diavolo", 8, 14, 6, 110, 90, "diavolo", DiavoloPattern{}, ColorPurple)
	b.Style = AttackErase
	return b
}

func NewKira() *Boss {
	b := New("Yoshikage Kira", 11, 18, 8, 150, 140, "kira", KiraPattern{}, ColorCyan)
	b.Style = AttackBombs
	return b
}

func NewPucci() *Boss {
	b := New("Enrico Pucci", 15, 24, 10, 220, 260, "pucci", PucciPattern{}, ColorRed)
	b.Style = AttackAccelerate
	return b
}

func NewIggy() *Boss {
	b := New("Iggy", 3, 5, 1, 30, 18, "iggy", IggyPattern{}, ColorWhite)
	b.Style = AttackStroll
	b.FragmentReward = 6
	return b
}
