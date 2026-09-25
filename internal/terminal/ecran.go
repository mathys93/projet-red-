package terminal

import (
	"fmt"
	"strings"
)

const (
	ColWhite  = "\033[97m"
	ColYellow = "\033[93m"
	ColRed    = "\033[91m"

	bgNeutralCode = 234
	ColReset      = "\033[0m\033[48;5;234m"
)

func ClearScreen() {
	Out.clear()
}

func RestoreTerminal() {
	SetScene(nil)
	fmt.Fprint(Out, "\033[0m\033[?25h")
}

func MaximizeConsoleWindow() {
	maximizeConsoleWindow()
}

func TerminalSize() (cols, rows int) {
	cols, rows = queryTerminalSize()
	if cols < 40 {
		cols = 80
	}
	if rows < 20 {
		rows = 24
	}
	return cols, rows
}

const (
	minBoxWidth   = 66
	minBoxHeight  = 8
	MinArtHeight  = 6
	battleChrome  = 12
	tallTerminal  = 46
	tallBoxHeight = 10
)

func Layout() (boxW, boxH, artW, artH int) {
	cols, rows := TerminalSize()

	boxW = cols - 4
	if boxW < minBoxWidth {
		boxW = minBoxWidth
	}

	boxH = minBoxHeight
	if rows >= tallTerminal {
		boxH = tallBoxHeight
	}

	artH = rows - 1 - battleChrome - boxH
	for artH < MinArtHeight && boxH > 5 {
		boxH--
		artH++
	}
	if artH < MinArtHeight {
		artH = MinArtHeight
	}

	return boxW, boxH, boxW, artH
}

func PrintPadding(contentHeight int) {
	_, rows := TerminalSize()
	pad := (rows - contentHeight) / 2
	for i := 0; i < pad; i++ {
		fmt.Fprintln(Out)
	}
}

func VisibleWidth(s string) int {
	n := 0
	inEsc := false
	for _, r := range s {
		switch {
		case inEsc:
			if r == 'm' {
				inEsc = false
			}
		case r == '\033':
			inEsc = true
		default:
			n++
		}
	}
	return n
}

func DrawFrameTop(width int) {
	fmt.Fprintln(Out, ColWhite+"┌"+strings.Repeat("─", width-2)+"┐"+ColReset)
}

func DrawFrameBottom(width int) {
	fmt.Fprintln(Out, ColWhite+"└"+strings.Repeat("─", width-2)+"┘"+ColReset)
}

func DrawFrameLine(width int, content string) {
	pad := width - 2 - VisibleWidth(content)
	if pad < 0 {
		pad = 0
	}
	fmt.Fprintln(Out, ColWhite+"│"+ColReset+content+strings.Repeat(" ", pad)+ColWhite+"│"+ColReset)
}

func DrawArt(width int, art string, artColor string) {
	if artColor == "" {
		artColor = ColWhite
	}
	for _, line := range strings.Split(art, "\n") {
		pad := (width - len([]rune(line))) / 2
		if pad < 0 {
			pad = 0
		}
		fmt.Fprintln(Out, strings.Repeat(" ", pad)+artColor+line+ColReset)
	}
}

func DrawEmptyBox(width, height int) {
	DrawFrameTop(width)
	for i := 0; i < height-2; i++ {
		DrawFrameLine(width, "")
	}
	DrawFrameBottom(width)
}

func DrawOptionGrid(width, height int, items []MenuItem, selected int) {
	DrawFrameTop(width)

	const cols = 2
	colWidth := (width - 2 - (cols + 1)) / cols
	rows := (len(items) + cols - 1) / cols

	interior := height - 2
	top := (interior - rows) / 2
	if top < 0 {
		top = 0
	}

	for i := 0; i < top; i++ {
		DrawFrameLine(width, "")
	}
	for r := 0; r < rows; r++ {
		line := " "
		for c := 0; c < cols; c++ {
			idx := r*cols + c
			cell := strings.Repeat(" ", colWidth)
			if idx < len(items) {
				marker := "  "
				col := ColYellow
				if idx == selected {
					marker = ColRed + "❤ " + ColReset
					col = ColRed
				}
				text := marker + col + items[idx].Label + ColReset
				pad := colWidth - VisibleWidth(text)
				if pad < 0 {
					pad = 0
				}
				cell = text + strings.Repeat(" ", pad)
			}
			line += cell + " "
		}
		DrawFrameLine(width, line)
	}
	for i := 0; i < interior-top-rows; i++ {
		DrawFrameLine(width, "")
	}

	DrawFrameBottom(width)
}

type MenuItem struct {
	Label       string
	Description string
}

func CursorAt(row, col int) string {
	return fmt.Sprintf("\033[%d;%dH", row, col)
}
