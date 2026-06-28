package tui

import (
	"strings"

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

// appStyle frames the whole home screen (margin around header + list).
var appStyle = lipgloss.NewStyle().Margin(1, 2)

// dividerStyle colors the vertical line between the logo and the title block.
var dividerStyle = lipgloss.NewStyle().Foreground(accent2)

// homeHeader is the home screen header: the animated cake logo on the left, a
// vertical divider, then the tool name, description, and version on the right.
// phase drives the logo's gradient animation.
func homeHeader(phase int, version string) string {
	logo := renderLogo(phase)

	info := lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render(" BAKE AI "),
		"",
		helpStyle.Render("Your personalized, project-aware AI assistant"),
		labelStyle.Render("version ")+valueStyle.Render(version),
	)

	height := lipgloss.Height(logo)
	divider := dividerStyle.Render(strings.TrimRight(strings.Repeat("│\n", height), "\n"))

	return lipgloss.JoinHorizontal(lipgloss.Top,
		logo, "  ", divider, "  ", info,
	)
}

// listHeader is the title + one-line description shown directly above the
// project list.
func listHeader() string {
	return lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render("Your projects"),
		helpStyle.Render("Pick a project to open, or start a new one"),
	)
}

type pickerModel struct {
	list    list.Model
	result  Result
	phase   int    // logo gradient animation frame
	version string // shown in the header
}

func newPickerModel(projects []workspace.Project, version string) pickerModel {
	items := make([]list.Item, 0, len(projects)+1)
	for _, p := range projects {
		desc := p.Description
		if desc == "" {
			desc = p.Path
		}
		items = append(items, pickerItem{kind: kindProject, name: p.Name, desc: desc})
	}
	items = append(items, pickerItem{kind: kindNew, desc: "scaffold a new project"})

	delegate := list.NewDefaultDelegate()
	delegate.Styles.NormalTitle = delegate.Styles.NormalTitle.Foreground(fg)
	delegate.Styles.NormalDesc = delegate.Styles.NormalDesc.Foreground(faint)
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.Foreground(accent).BorderLeftForeground(accent)
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.Foreground(accent2).BorderLeftForeground(accent)

	l := list.New(items, delegate, 0, 0)
	l.SetShowTitle(false) // we render our own title + description above the list
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)

	return pickerModel{list: l, version: version}
}

func (m pickerModel) Init() tea.Cmd { return tick() }

func (m pickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		m.phase++
		return m, tick()
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		// Leave room for the home header, the list header (title + desc), and a
		// one-line spacer above and below the list header — but never let the
		// list height go negative on short terminals.
		listH := msg.Height - v - lipgloss.Height(m.header()) - lipgloss.Height(listHeader()) - 2
		if listH < 1 {
			listH = 1
		}
		m.list.SetSize(msg.Width-h, listH)
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

// header renders this model's home header at its current animation phase.
func (m pickerModel) header() string { return homeHeader(m.phase, m.version) }

func (m pickerModel) View() string {
	body := lipgloss.JoinVertical(lipgloss.Left,
		m.header(),
		"",
		listHeader(),
		"",
		m.list.View(),
	)
	return appStyle.Render(body)
}
