package tui

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/alexpagnotta/bake-ai/internal/config"
	"github.com/alexpagnotta/bake-ai/internal/workspace"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// Action is what the user chose in the home screen.
type Action int

const (
	ActionQuit Action = iota
	ActionChat
	ActionNew
)

// Result is the outcome of the picker.
type Result struct {
	Action  Action
	Project string
}

// RunPicker shows the project picker (alt screen) and returns the chosen action.
func RunPicker(c *config.Config) (Result, error) {
	projects, err := workspace.List(c)
	if err != nil {
		return Result{}, err
	}
	final, err := tea.NewProgram(newPickerModel(projects), tea.WithAltScreen()).Run()
	if err != nil {
		return Result{}, err
	}
	return final.(pickerModel).result, nil
}

// NewProjectForm collects a new project's name and model via a Huh form.
// ok is false if the user aborted.
func NewProjectForm(c *config.Config) (name, model string, ok bool, err error) {
	model = c.DefaultModel
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Project name").
				Description("lowercase letters, digits, '-' or '_'").
				Value(&name).
				Validate(func(s string) error {
					if !workspace.ValidName(s) {
						return errors.New("use lowercase letters, digits, '-' or '_'")
					}
					return nil
				}),
			huh.NewSelect[string]().
				Title("Default model").
				Options(
					huh.NewOption(c.DefaultModel+"  (workspace default)", c.DefaultModel),
					huh.NewOption("anthropic/claude-3.5-sonnet", "anthropic/claude-3.5-sonnet"),
					huh.NewOption("openai/gpt-4o-mini", "openai/gpt-4o-mini"),
					huh.NewOption("google/gemini-2.5-flash", "google/gemini-2.5-flash"),
				).
				Value(&model),
		),
	).WithTheme(formTheme())
	if err := form.Run(); err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return "", "", false, nil
		}
		return "", "", false, err
	}
	return name, model, true, nil
}

// PrintChatHeader prints a styled banner (with a Glamour-rendered vault preview)
// before bake hands the terminal to goose.
func PrintChatHeader(p *workspace.Project) {
	info := lipgloss.JoinVertical(lipgloss.Left,
		labelStyle.Render("project ")+valueStyle.Render(p.Name),
		labelStyle.Render("path    ")+valueStyle.Render(p.Path),
		labelStyle.Render("model   ")+valueStyle.Render(recipeModel(p)),
	)
	fmt.Println()
	fmt.Println(titleStyle.Render(" bake · " + p.Name + " "))
	fmt.Println(panelStyle.Render(info))
	fmt.Println(helpStyle.Render("Starting goose… type your message; Ctrl+C ends the session."))
	fmt.Println()
}

var modelRe = regexp.MustCompile(`goose_model:\s*"?([^"\n]+)"?`)

// recipeModel best-effort reads the model from the project's recipe.yaml.
func recipeModel(p *workspace.Project) string {
	b, err := os.ReadFile(filepath.Join(p.Path, "recipe.yaml"))
	if err != nil {
		return "(unknown)"
	}
	if m := modelRe.FindSubmatch(b); m != nil {
		return strings.TrimSpace(string(m[1]))
	}
	return "(unknown)"
}
