package tui

import (
	_ "embed"
	"fmt"
	"math"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// bannerArt is the big "BAKE" banner shown on the home screen. Cells are filled
// (any non-'.', non-space rune) or empty ('.' / space). Filled cells get the
// animated gradient; empty cells render as blank, so the letters show through.
//
//go:embed banner.txt
var bannerSrc string

var bakeArt = strings.Split(strings.TrimRight(bannerSrc, "\n"), "\n")

// tickMsg drives the gradient animation; one arrives on every frame.
type tickMsg time.Time

// bannerFPS is the animation cadence. ~12 fps is smooth enough for a flowing
// gradient while staying cheap on the terminal.
const bannerFPS = 12

// tick schedules the next animation frame.
func tick() tea.Cmd {
	return tea.Tick(time.Second/bannerFPS, func(t time.Time) tea.Msg { return tickMsg(t) })
}

// renderBanner draws bakeArt with a lilac→pink gradient that ripples diagonally.
// The color of each filled cell depends on (x+y), so equal-color bands run along
// the anti-diagonal; subtracting phase makes those bands drift over time, giving
// the wavy motion. Empty cells ('.'/space) are left blank.
func renderBanner(phase int) string {
	var b strings.Builder
	for y, line := range bakeArt {
		for x, r := range line {
			if r == '.' || r == ' ' {
				b.WriteByte(' ')
				continue
			}
			t := 0.5 + 0.5*math.Sin((float64(x)+float64(y))*0.14-float64(phase)*0.30)
			st := lipgloss.NewStyle().Foreground(lerpColor(brandLilac, brandPink, t))
			b.WriteString(st.Render(string(r)))
		}
		b.WriteByte('\n')
	}
	return strings.TrimRight(b.String(), "\n")
}

// lerpColor linearly interpolates between two #RRGGBB colors. t is clamped 0..1.
func lerpColor(from, to string, t float64) lipgloss.Color {
	if t < 0 {
		t = 0
	} else if t > 1 {
		t = 1
	}
	fr, fg, fb := hexToRGB(from)
	tr, tg, tb := hexToRGB(to)
	r := int(float64(fr) + (float64(tr)-float64(fr))*t)
	g := int(float64(fg) + (float64(tg)-float64(fg))*t)
	bl := int(float64(fb) + (float64(tb)-float64(fb))*t)
	return lipgloss.Color(fmt.Sprintf("#%02X%02X%02X", r, g, bl))
}

// hexToRGB parses a #RRGGBB string into its components.
func hexToRGB(h string) (r, g, b int) {
	fmt.Sscanf(strings.TrimPrefix(h, "#"), "%02x%02x%02x", &r, &g, &b)
	return r, g, b
}
