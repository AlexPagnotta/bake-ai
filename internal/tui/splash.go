package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// splashDuration is how long the full-screen animated banner shows on startup.
const splashDuration = 1500 * time.Millisecond

// RunSplash shows the animated BAKE banner centered full-screen for a few seconds
// (or until a key is pressed), then returns. quit is true if the user pressed
// Ctrl+C to abort straight out of the app.
func RunSplash() (quit bool, err error) {
	final, err := tea.NewProgram(splashModel{}, tea.WithAltScreen()).Run()
	if err != nil {
		return false, err
	}
	return final.(splashModel).quit, nil
}

// splashDoneMsg fires once, when the splash timer elapses.
type splashDoneMsg struct{}

type splashModel struct {
	phase         int // gradient animation frame
	width, height int
	quit          bool
}

func (m splashModel) Init() tea.Cmd {
	return tea.Batch(
		tick(),
		tea.Tick(splashDuration, func(time.Time) tea.Msg { return splashDoneMsg{} }),
	)
}

func (m splashModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		m.phase++
		return m, tick()
	case splashDoneMsg:
		return m, tea.Quit
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tea.KeyMsg:
		// Any key skips the splash; Ctrl+C skips *and* quits the app.
		if msg.String() == "ctrl+c" {
			m.quit = true
		}
		return m, tea.Quit
	}
	return m, nil
}

func (m splashModel) View() string {
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, renderBanner(m.phase))
}
