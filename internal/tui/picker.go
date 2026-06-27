package tui

import (
	"github.com/alexpagnotta/bake-ai/internal/workspace"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type itemKind int

const (
	kindProject itemKind = iota
	kindNew
)

// pickerItem is one row in the project picker.
type pickerItem struct {
	kind itemKind
	name string
	desc string
}

func (i pickerItem) Title() string {
	if i.kind == kindNew {
		return "＋ New project"
	}
	return i.name
}
func (i pickerItem) Description() string { return i.desc }
func (i pickerItem) FilterValue() string { return i.name }

type pickerModel struct {
	list   list.Model
	result Result
}

func newPickerModel(projects []workspace.Project) pickerModel {
	items := make([]list.Item, 0, len(projects)+1)
	for _, p := range projects {
		items = append(items, pickerItem{kind: kindProject, name: p.Name, desc: p.Path})
	}
	items = append(items, pickerItem{kind: kindNew, desc: "scaffold a new project"})

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.Foreground(accent).BorderLeftForeground(accent)
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.Foreground(accent2).BorderLeftForeground(accent)

	l := list.New(items, delegate, 0, 0)
	l.Title = "bake — your projects"
	l.Styles.Title = titleStyle
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)

	return pickerModel{list: l}
}

func (m pickerModel) Init() tea.Cmd { return nil }

func (m pickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := lipgloss.NewStyle().Margin(1, 2).GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	case tea.KeyMsg:
		// While the user is typing a filter, let the list own every key.
		if m.list.FilterState() == list.Filtering {
			break
		}
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			m.result = Result{Action: ActionQuit}
			return m, tea.Quit
		case "n":
			m.result = Result{Action: ActionNew}
			return m, tea.Quit
		case "enter":
			if it, ok := m.list.SelectedItem().(pickerItem); ok {
				if it.kind == kindNew {
					m.result = Result{Action: ActionNew}
				} else {
					m.result = Result{Action: ActionChat, Project: it.name}
				}
				return m, tea.Quit
			}
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m pickerModel) View() string {
	return lipgloss.NewStyle().Margin(1, 2).Render(m.list.View())
}
