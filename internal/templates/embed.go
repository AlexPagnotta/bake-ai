// Package templates renders the embedded project scaffold into a new project dir.
package templates

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed all:project
var projectFS embed.FS

// ProjectData is the substitution context for project templates.
type ProjectData struct {
	Name  string
	Model string
}

// RenderProject renders the embedded project template tree into destDir.
func RenderProject(destDir string, data ProjectData) error {
	return fs.WalkDir(projectFS, "project", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel("project", path)
		if err != nil {
			return err
		}
		if d.IsDir() {
			return os.MkdirAll(filepath.Join(destDir, rel), 0o755)
		}
		return renderFile(projectFS, path, filepath.Join(destDir, mapName(rel)), data)
	})
}

func renderFile(srcFS fs.FS, srcPath, outPath string, data ProjectData) error {
	content, err := fs.ReadFile(srcFS, srcPath)
	if err != nil {
		return err
	}
	tmpl, err := template.New(filepath.Base(srcPath)).Parse(string(content))
	if err != nil {
		return fmt.Errorf("parse template %s: %w", srcPath, err)
	}
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return err
	}
	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	if err := tmpl.Execute(f, data); err != nil {
		f.Close()
		return fmt.Errorf("render %s: %w", outPath, err)
	}
	return f.Close()
}

// mapName turns a template path into its on-disk name: strips the .tmpl suffix and
// restores the leading dot for goosehints (embed can't ship dotfiles cleanly).
func mapName(rel string) string {
	dir, base := filepath.Split(rel)
	if base == "goosehints.tmpl" {
		base = ".goosehints"
	} else {
		base = strings.TrimSuffix(base, ".tmpl")
	}
	return filepath.Join(dir, base)
}
