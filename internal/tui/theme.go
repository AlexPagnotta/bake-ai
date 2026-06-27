// Package tui is bake's Charm/Bubble Tea front-end. It is a thin presentation
// layer over the internal engine (config, workspace, gooserun) — no core logic
// lives here, so the CLI and TUI share the same behavior.
package tui

import "github.com/charmbracelet/lipgloss"

// Palette — a warm "bake" accent (toasted amber) over neutral grays.
var (
	accent  = lipgloss.Color("214") // amber
	accent2 = lipgloss.Color("130") // darker toast
	faint   = lipgloss.Color("240")
	fg      = lipgloss.Color("252")
)

// Shared styles, reused across screens.
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("0")).
			Background(accent).
			Padding(0, 1)

	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(accent2).
			Padding(0, 1)

	labelStyle = lipgloss.NewStyle().Foreground(accent).Bold(true)
	valueStyle = lipgloss.NewStyle().Foreground(fg)
	helpStyle  = lipgloss.NewStyle().Foreground(faint)
)
