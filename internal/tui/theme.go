// Package tui is bake's Charm/Bubble Tea front-end. It is a thin presentation
// layer over the internal engine (config, workspace, gooserun) — no core logic
// lives here, so the CLI and TUI share the same behavior.
package tui

import (
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// Brand palette — hot pink + sky cyan are the two primary colors; white and
// lilac are the only extras (for text/muted accents). Also recorded in AGENTS.md
// "Design" so future features reuse the same colors and nothing else.
const (
	brandPink  = "#FE5283"
	brandCyan  = "#6AD5FD"
	brandWhite = "#FFFFFF"
	brandLilac = "#EEADEE"
)

// Palette — everything maps to a brand color; no neutral grays.
var (
	accent  = lipgloss.Color(brandPink)  // hot pink — primary accent
	accent2 = lipgloss.Color(brandCyan)  // sky cyan — secondary accent
	faint   = lipgloss.Color(brandLilac) // lilac — muted / secondary text
	fg      = lipgloss.Color(brandWhite) // white — primary text
)

// Shared styles, reused across screens.
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(fg).
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

// formTheme is a Huh form theme built purely from the brand palette — it strips
// out Huh's default green/indigo/fuchsia (e.g. the green selected option) so the
// new-project form matches the rest of the TUI.
func formTheme() *huh.Theme {
	t := huh.ThemeBase()
	f := &t.Focused

	f.Base = f.Base.BorderForeground(accent2)
	f.Card = f.Base
	f.Title = f.Title.Foreground(accent).Bold(true)
	f.NoteTitle = f.NoteTitle.Foreground(accent).Bold(true)
	f.Description = f.Description.Foreground(faint)
	f.ErrorIndicator = f.ErrorIndicator.Foreground(accent)
	f.ErrorMessage = f.ErrorMessage.Foreground(accent)
	f.SelectSelector = f.SelectSelector.Foreground(accent2)
	f.NextIndicator = f.NextIndicator.Foreground(accent2)
	f.PrevIndicator = f.PrevIndicator.Foreground(accent2)
	f.Option = f.Option.Foreground(fg)
	f.MultiSelectSelector = f.MultiSelectSelector.Foreground(accent2)
	f.SelectedOption = f.SelectedOption.Foreground(accent)
	f.SelectedPrefix = lipgloss.NewStyle().Foreground(accent2).SetString("✓ ")
	f.UnselectedPrefix = lipgloss.NewStyle().Foreground(faint).SetString("• ")
	f.UnselectedOption = f.UnselectedOption.Foreground(fg)
	f.FocusedButton = f.FocusedButton.Foreground(fg).Background(accent)
	f.Next = f.FocusedButton
	f.BlurredButton = f.BlurredButton.Foreground(faint).UnsetBackground()

	f.TextInput.Cursor = f.TextInput.Cursor.Foreground(accent)
	f.TextInput.Placeholder = f.TextInput.Placeholder.Foreground(faint)
	f.TextInput.Prompt = f.TextInput.Prompt.Foreground(accent2)

	t.Blurred = t.Focused
	t.Blurred.Base = t.Blurred.Base.BorderStyle(lipgloss.HiddenBorder())
	t.Blurred.Card = t.Blurred.Base
	t.Blurred.NextIndicator = lipgloss.NewStyle()
	t.Blurred.PrevIndicator = lipgloss.NewStyle()

	t.Group.Title = t.Focused.Title
	t.Group.Description = t.Focused.Description
	return t
}
