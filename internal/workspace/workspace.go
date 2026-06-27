// Package workspace contains bake's core project operations — presentation-free,
// so both the CLI (Phase 1) and the TUI (Phase 2) call the same functions.
package workspace

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"

	"github.com/alexpagnotta/bake-ai/internal/config"
	"github.com/alexpagnotta/bake-ai/internal/templates"
)

var nameRe = regexp.MustCompile(`^[a-z0-9][a-z0-9_-]*$`)

// ValidName reports whether name is a usable project name.
func ValidName(name string) bool { return nameRe.MatchString(name) }

// Project is a single project in the workspace.
type Project struct {
	Name string
	Path string
}

// ProjectsDir returns the directory holding all projects.
func ProjectsDir(c *config.Config) string {
	return filepath.Join(c.WorkspacePath, "projects")
}

// Create scaffolds a new project from the embedded templates. If model is empty,
// the workspace default is used.
func Create(c *config.Config, name, model string) (*Project, error) {
	if !nameRe.MatchString(name) {
		return nil, fmt.Errorf("invalid project name %q: use lowercase letters, digits, '-' or '_'", name)
	}
	if model == "" {
		model = c.DefaultModel
	}
	dir := filepath.Join(ProjectsDir(c), name)
	switch _, err := os.Stat(dir); {
	case err == nil:
		return nil, fmt.Errorf("project %q already exists at %s", name, dir)
	case !os.IsNotExist(err):
		return nil, err
	}
	if err := templates.RenderProject(dir, templates.ProjectData{Name: name, Model: model}); err != nil {
		return nil, err
	}
	return &Project{Name: name, Path: dir}, nil
}

// List returns all projects in the workspace, sorted by name.
func List(c *config.Config) ([]Project, error) {
	root := ProjectsDir(c)
	entries, err := os.ReadDir(root)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var projects []Project
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		if _, err := os.Stat(filepath.Join(root, e.Name(), "recipe.yaml")); err != nil {
			continue
		}
		projects = append(projects, Project{Name: e.Name(), Path: filepath.Join(root, e.Name())})
	}
	sort.Slice(projects, func(i, j int) bool { return projects[i].Name < projects[j].Name })
	return projects, nil
}

// Get returns a project by name, or an error if it doesn't exist.
func Get(c *config.Config, name string) (*Project, error) {
	dir := filepath.Join(ProjectsDir(c), name)
	if _, err := os.Stat(filepath.Join(dir, "recipe.yaml")); err != nil {
		return nil, fmt.Errorf("project %q not found (looked in %s)", name, dir)
	}
	return &Project{Name: name, Path: dir}, nil
}
