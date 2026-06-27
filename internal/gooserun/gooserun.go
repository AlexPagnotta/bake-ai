// Package gooserun launches goose and inspects its (non-secret) configuration.
// bake never reads, stores, or logs secrets — goose owns the API key in the keyring.
package gooserun

import (
	"bufio"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"

	"github.com/alexpagnotta/bake-ai/internal/workspace"
)

// LaunchChat hands the terminal to goose for an interactive session driven by the
// project's recipe, with the working directory set to the project so .goosehints
// and the docs/ + vault/ files are in scope.
//
// goose shares our terminal and process group, so a Ctrl+C delivers SIGINT to
// both processes. We ignore it in bake for the duration of the session — goose
// owns the interrupt and uses it to end the chat — so that ending the session
// returns control to bake (e.g. back to the picker) instead of killing it.
func LaunchChat(p *workspace.Project) error {
	signal.Ignore(os.Interrupt)
	defer signal.Reset(os.Interrupt)

	cmd := exec.Command("goose", "run", "--recipe", filepath.Join(p.Path, "recipe.yaml"), "--interactive")
	cmd.Dir = p.Path
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Installed reports whether the goose binary is on PATH.
func Installed() bool {
	_, err := exec.LookPath("goose")
	return err == nil
}

// gooseConfigPath returns goose's config.yaml location, honoring XDG_CONFIG_HOME
// (which goose itself respects) before falling back to ~/.config.
func gooseConfigPath() (string, error) {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "goose", "config.yaml"), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "goose", "config.yaml"), nil
}

// HasOpenRouter does a best-effort, secret-free check of goose's config for an
// OpenRouter provider. Returns false if it can't tell — callers treat it as a hint.
func HasOpenRouter() bool {
	path, err := gooseConfigPath()
	if err != nil {
		return false
	}
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		if strings.Contains(strings.ToLower(sc.Text()), "openrouter") {
			return true
		}
	}
	return false
}
