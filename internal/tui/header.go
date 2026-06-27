package tui

import (
	"fmt"
	"math"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// bakeArt is the big, bold "BAKE" banner on the home screen. All lines are the
// same width so the animated gradient lines up cleanly column-to-column.
var bakeArt = []string{
	` ___    _    _  _  ___ `,
	`| _ )  /_\  | |/ // __|`,
	`| _ \ / _ \ |   < | _| `,
	`|___//_/ \_\|_|\_\|___|`,
}

// tickMsg drives the gradient animation; one arrives on every frame.
type tickMsg time.Time

// bannerFPS is the animation cadence. ~12 fps is smooth enough for a sweeping
// gradient while staying cheap on the terminal.
const bannerFPS = 12

// tick schedules the next animation frame.
func tick() tea.Cmd {
	return tea.Tick(time.Second/bannerFPS, func(t time.Time) tea.Msg { return tickMsg(t) })
}

// renderBanner draws bakeArt with a pink→cyan gradient that sweeps left→right.
// phase advances each frame, shifting the gradient so it appears to flow.
func renderBanner(phase int) string {
	var b strings.Builder
	for _, line := range bakeArt {
		for x, r := range line {
			if r == ' ' {
				b.WriteByte(' ')
				continue
			}
			// A sine wave across the columns, shifted by phase, gives a smooth
			// looping sweep. Subtracting phase moves the highlight to the right.
			t := 0.5 + 0.5*math.Sin(float64(x)*0.30-float64(phase)*0.45)
			st := lipgloss.NewStyle().Bold(true).Foreground(lerpColor(brandPink, brandCyan, t))
			b.WriteString(st.Render(string(r)))
		}
		b.WriteByte('\n')
	}
	return strings.TrimRight(b.String(), "\n")
}

// bannerHeight is the rendered banner's line count (constant — art is fixed).
func bannerHeight() int { return len(bakeArt) }

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
