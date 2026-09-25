package ui

import (
	"fmt"
	"image"
	"strings"
	"unicode"

	"ProjetRED/internal/ascii"
	"ProjetRED/internal/audio"
)

type MenuChoice int

const (
	MenuStart MenuChoice = iota
	MenuSettings
	MenuCredits
	MenuQuit
)

const (
	egaBlue    = "\033[38;2;0;0;170m"
	egaGrey    = "\033[38;2;170;170;170m"
	egaMagenta = "\033[38;2;170;0;170m"
	egaWhite   = "\033[38;2;255;255;255m"

	maxNameLen  = 14
	defaultName = "Personnage 1"
	titleMargin = 8
)

type creditSection struct {
	heading string
	lines   []string
}

var Credits = []creditSection{
	{"Développement", []string{"Mhenri", "Ouanis Benabbou", "bmathys", "mathys93"}},
	{"Inspiré de", []string{"UNDERTALE — Toby Fox", "JoJo's Bizarre Adventure — Hirohiko Araki"}},
	{"Musique", []string{"Fichiers du dossier assets/musique/"}},
}

type titlePixel struct {
	r, g, b uint8
	on      bool
}

func samplePixel(img image.Image, x0, y0, x1, y1 int) titlePixel {
	var sr, sg, sb, lit, total int
	for y := y0; y < y1; y++ {
		for x := x0; x < x1; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			r8, g8, b8 := int(r>>8), int(g>>8), int(b>>8)
			total++
			if r8 < 40 && g8 < 40 && b8 < 40 {
				continue
			}
			sr, sg, sb = sr+r8, sg+g8, sb+b8
			lit++
		}
	}
	if total == 0 || lit*2 < total {
		return titlePixel{}
	}
	return titlePixel{uint8(sr / lit), uint8(sg / lit), uint8(sb / lit), true}
}

func titlePixels(img image.Image, cols, rows int) [][]titlePixel {
	bd := img.Bounds()
	iw, ih := bd.Dx(), bd.Dy()
	pxRows := rows * 2
	grid := make([][]titlePixel, pxRows)
	for py := range pxRows {
		grid[py] = make([]titlePixel, cols)
		y0 := bd.Min.Y + py*ih/pxRows
		y1 := bd.Min.Y + (py+1)*ih/pxRows
		if y1 <= y0 {
			y1 = y0 + 1
		}
		for cx := range cols {
			x0 := bd.Min.X + cx*iw/cols
			x1 := bd.Min.X + (cx+1)*iw/cols
			if x1 <= x0 {
				x1 = x0 + 1
			}
			grid[py][cx] = samplePixel(img, x0, y0, x1, y1)
		}
	}
	return grid
}

type logoLine struct {
	text  string
	width int
}

var (
	logoCacheKey   [2]int
	logoCacheLines []logoLine
)

func renderLogo(maxCols, maxRows int) []logoLine {
	title, sub := ascii.Title(), ascii.Subtitle()
	if title == nil || maxCols < 1 || maxRows < 3 {
		return []logoLine{{egaWhite + "undertale jom" + ColReset, 13}}
	}
	key := [2]int{maxCols, maxRows}
	if key == logoCacheKey && logoCacheLines != nil {
		return logoCacheLines
	}

	tw, th := title.Bounds().Dx(), title.Bounds().Dy()
	sh := 0
	if sub != nil {
		sh = sub.Bounds().Dy()
	}
	scale := float64(maxCols) / float64(tw)
	gap := 0
	if sub != nil {
		gap = 1
	}
	for rowsFor(th, scale)+rowsFor(sh, scale)+gap > maxRows {
		scale *= 0.97
	}

	var lines []logoLine
	add := func(img image.Image) {
		cols := max(1, int(float64(img.Bounds().Dx())*scale))
		rows := max(1, rowsFor(img.Bounds().Dy(), scale))
		for _, l := range renderImage(img, cols, rows) {
			lines = append(lines, logoLine{l, cols})
		}
	}
	add(title)
	if sub != nil {
		lines = append(lines, logoLine{"", 0})
		add(sub)
	}

	logoCacheKey, logoCacheLines = key, lines
	return lines
}

func rowsFor(pixelHeight int, scale float64) int {
	if pixelHeight == 0 {
		return 0
	}
	return int(float64(pixelHeight)*scale/2 + 0.5)
}

func renderImage(img image.Image, cols, rows int) []string {
	grid := titlePixels(img, cols, rows)
	lines := make([]string, rows)
	for cy := range rows {
		var sb strings.Builder
		for cx := range cols {
			top, bot := grid[2*cy][cx], grid[2*cy+1][cx]
			switch {
			case top.on && bot.on:
				fmt.Fprintf(&sb, "\033[38;2;%d;%d;%dm\033[48;2;%d;%d;%dm▀", top.r, top.g, top.b, bot.r, bot.g, bot.b)
			case top.on:
				fmt.Fprintf(&sb, "%s\033[38;2;%d;%d;%dm▀", ColReset, top.r, top.g, top.b)
			case bot.on:
				fmt.Fprintf(&sb, "%s\033[38;2;%d;%d;%dm▄", ColReset, bot.r, bot.g, bot.b)
			default:
				sb.WriteString(ColReset)
				sb.WriteString(" ")
			}
		}
		sb.WriteString(ColReset)
		lines[cy] = sb.String()
	}
	return lines
}

func printCentered(cols int, s string) {
	pad := max((cols-VisibleWidth(s))/2, 0)
	fmt.Fprintln(Out, strings.Repeat(" ", pad)+s)
}

func printLogo(cols int, lines []logoLine) {
	for _, l := range lines {
		fmt.Fprintln(Out, strings.Repeat(" ", max(0, (cols-l.width)/2))+l.text)
	}
}

func drawMenuButtons(width int, labels []string, selected int) {
	inner := 0
	for _, l := range labels {
		if n := len([]rune(l)); n > inner {
			inner = n
		}
	}
	inner += 6
	btnOuter := inner + 2
	gap := max((width-len(labels)*btnOuter)/(len(labels)+1), 2)
	margin := strings.Repeat(" ", gap)

	var top, mid, bot, shadow strings.Builder
	for i, label := range labels {
		border, text, heart := egaBlue, egaGrey, "  "
		if i == selected {
			border, text, heart = egaMagenta, egaWhite, ColRed+"❤ "+ColReset
		}
		pad := inner - 2 - len([]rune(label))
		left := pad / 2
		right := pad - left

		top.WriteString(margin)
		top.WriteString(border)
		top.WriteString("╔")
		top.WriteString(strings.Repeat("═", inner))
		top.WriteString("╗")
		top.WriteString(ColReset)
		mid.WriteString(margin)
		mid.WriteString(border)
		mid.WriteString("║")
		mid.WriteString(ColReset)
		mid.WriteString(heart)
		mid.WriteString(strings.Repeat(" ", left))
		mid.WriteString(text)
		mid.WriteString(label)
		mid.WriteString(ColReset)
		mid.WriteString(strings.Repeat(" ", right))
		mid.WriteString(border)
		mid.WriteString("║")
		mid.WriteString(ColReset)
		bot.WriteString(margin)
		bot.WriteString(border)
		bot.WriteString("╚")
		bot.WriteString(strings.Repeat("═", inner))
		bot.WriteString("╝")
		bot.WriteString(ColReset)
		shadow.WriteString(margin)
		shadow.WriteString(" ")
		shadow.WriteString(egaMagenta)
		shadow.WriteString(strings.Repeat("▀", btnOuter-1))
		shadow.WriteString(ColReset)
	}

	fmt.Fprintln(Out)
	fmt.Fprintln(Out, top.String())
	fmt.Fprintln(Out, mid.String())
	fmt.Fprintln(Out, bot.String())
	fmt.Fprintln(Out, shadow.String())
}

func MainMenu() MenuChoice {
	SetScene(SceneMenu)
	ts := NewSession()
	defer ts.Restore()
	fmt.Fprint(Out, "\033]0;UNDERTALE\007")

	labels := []string{"COMMENCER", "PARAMÈTRES", "CRÉDITS"}
	selected := 0
	for {
		cols, rows := TerminalSize()
		ClearScreen()

		logo := renderLogo(cols-titleMargin, rows-12)
		top := max((rows-len(logo))/2-3, 1)
		gap := max(rows-8-top-len(logo), 1)

		for range top {
			fmt.Fprintln(Out)
		}
		printLogo(cols, logo)
		for range gap {
			fmt.Fprintln(Out)
		}
		drawMenuButtons(cols, labels, selected)
		fmt.Fprintln(Out)
		printCentered(cols, egaGrey+"← → pour choisir · Entrée pour valider · Échap pour quitter"+ColReset)

		switch ts.ReadKey() {
		case KeyLeft, KeyUp:
			selected = (selected + len(labels) - 1) % len(labels)
		case KeyRight, KeyDown:
			selected = (selected + 1) % len(labels)
		case KeyEnter:
			return MenuChoice(selected)
		case KeyPause, KeyQuit:
			return MenuQuit
		case KeyOther:
			switch ts.Raw {
			case "1":
				return MenuStart
			case "2":
				return MenuSettings
			case "3":
				return MenuCredits
			}
		}
	}
}

func cleanName(s string) string {
	var out []rune
	for _, r := range strings.TrimSpace(s) {
		if unicode.IsPrint(r) && len(out) < maxNameLen {
			out = append(out, r)
		}
	}
	return string(out)
}

func drawScreenHeader(cols, rows, contentHeight int) {
	logo := renderLogo(cols-titleMargin, max(3, rows/4))
	top := (rows - len(logo) - contentHeight - 2) / 2
	if top < 1 {
		top = 1
	}
	for i := 0; i < top; i++ {
		fmt.Fprintln(Out)
	}
	printLogo(cols, logo)
	fmt.Fprintln(Out)
	fmt.Fprintln(Out)
}

func PromptName() (string, bool) {
	SetScene(SceneMenu)
	if !CanPollInput() {
		fmt.Fprint(Out, "Ton nom : ")
		if !stdinScanner.Scan() {
			return defaultName, true
		}
		if name := cleanName(stdinScanner.Text()); name != "" {
			return name, true
		}
		return defaultName, true
	}

	ts := NewSession()
	defer ts.Restore()

	var name []rune
	for {
		cols, rows := TerminalSize()
		ClearScreen()
		drawScreenHeader(cols, rows, 9)

		boxInner := maxNameLen + 6
		field := string(name) + egaMagenta + "_" + ColReset
		pad := boxInner - 2 - len(name) - 1
		if pad < 0 {
			pad = 0
		}
		printCentered(cols, egaWhite+"Quel est ton nom ?"+ColReset)
		fmt.Fprintln(Out)
		printCentered(cols, egaBlue+"╔"+strings.Repeat("═", boxInner)+"╗"+ColReset)
		printCentered(cols, egaBlue+"║"+ColReset+"  "+egaWhite+field+strings.Repeat(" ", pad)+egaBlue+"║"+ColReset)
		printCentered(cols, egaBlue+"╚"+strings.Repeat("═", boxInner)+"╝"+ColReset)
		printCentered(cols, " "+egaMagenta+strings.Repeat("▀", boxInner+1)+ColReset)
		fmt.Fprintln(Out)
		printCentered(cols, egaGrey+fmt.Sprintf("%d/%d caractères", len(name), maxNameLen)+ColReset)
		printCentered(cols, egaGrey+"Entrée pour valider · Retour arrière pour effacer · Échap pour revenir"+ColReset)

		k, ch := readTextKey()
		switch k {
		case textEnter:
			if n := cleanName(string(name)); n != "" {
				return n, true
			}
		case textCancel:
			return "", false
		case textBackspace:
			if len(name) > 0 {
				name = name[:len(name)-1]
			}
		case textChar:
			if len(name) < maxNameLen && unicode.IsPrint(ch) {
				name = append(name, ch)
			}
		}
	}
}

func volumeBar(v int) string {
	const slots = 20
	filled := v * slots / 1000
	return egaMagenta + strings.Repeat("■", filled) + egaBlue + strings.Repeat("□", slots-filled) + ColReset
}

func RunSettings() {
	defer SetScene(SetScene(SceneMenu))
	ts := NewSession()
	defer ts.Restore()

	const options = 3
	selected := 0
	for {
		cols, rows := TerminalSize()
		ClearScreen()
		drawScreenHeader(cols, rows, 10)

		etat := "Désactivée"
		if audio.Enabled {
			etat = "Activée"
		}
		entries := []string{
			fmt.Sprintf("Volume de la musique   ◄ %s ► %3d%%", volumeBar(audio.Volume), audio.Volume/10),
			fmt.Sprintf("Musique                %s", etat),
			"Retour",
		}

		printCentered(cols, egaWhite+"PARAMÈTRES"+ColReset)
		fmt.Fprintln(Out)
		width := 0
		for _, e := range entries {
			if n := VisibleWidth(e); n > width {
				width = n
			}
		}
		left := strings.Repeat(" ", max(0, (cols-width-4)/2))
		for i, e := range entries {
			marker, text := "  ", egaGrey
			if i == selected {
				marker, text = ColRed+"❤ "+ColReset, egaWhite
			}
			fmt.Fprintln(Out, left+marker+text+e+ColReset)
			fmt.Fprintln(Out)
		}
		printCentered(cols, egaGrey+"↑ ↓ pour choisir · ← → pour le volume · Entrée pour valider · Échap pour revenir"+ColReset)

		switch ts.ReadKey() {
		case KeyUp:
			selected = (selected + options - 1) % options
		case KeyDown:
			selected = (selected + 1) % options
		case KeyLeft:
			if selected == 0 {
				audio.SetVolume(audio.Volume - 50)
			}
		case KeyRight:
			if selected == 0 {
				audio.SetVolume(audio.Volume + 50)
			}
		case KeyEnter:
			switch selected {
			case 1:
				audio.SetEnabled(!audio.Enabled)
			case 2:
				return
			}
		case KeyBack, KeyPause, KeyQuit:
			return
		}
	}
}

func RunCredits() {
	defer SetScene(SetScene(SceneMenu))
	ts := NewSession()
	defer ts.Restore()

	height := 0
	for _, s := range Credits {
		height += len(s.lines) + 2
	}

	cols, rows := TerminalSize()
	ClearScreen()
	drawScreenHeader(cols, rows, height+3)
	printCentered(cols, egaWhite+"CRÉDITS"+ColReset)
	fmt.Fprintln(Out)
	for _, s := range Credits {
		printCentered(cols, egaMagenta+s.heading+ColReset)
		for _, l := range s.lines {
			printCentered(cols, egaGrey+l+ColReset)
		}
		fmt.Fprintln(Out)
	}
	printCentered(cols, egaGrey+"Entrée ou Échap pour revenir"+ColReset)

	for {
		switch ts.ReadKey() {
		case KeyEnter, KeyBack, KeyPause, KeyQuit:
			return
		}
	}
}
