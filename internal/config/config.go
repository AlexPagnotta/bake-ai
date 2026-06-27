// Package config handles bake's non-secret configuration. Secrets (the
// OpenRouter key) are never stored here — goose owns those in the OS keyring.
package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Config holds bake's non-secret settings.
type Config struct {
	// WorkspacePath is the directory holding the user's private projects.
	WorkspacePath string `toml:"workspace_path"`
	// DefaultModel is the goose/OpenRouter model new projects use unless overridden.
	DefaultModel string `toml:"default_model"`
}

// Default returns the baseline config used on first run.
func Default() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	return &Config{
		WorkspacePath: filepath.Join(home, "bake"),
		DefaultModel:  "google/gemini-2.5-flash",
	}, nil
}

// Dir returns bake's config directory, honoring XDG_CONFIG_HOME, else ~/.config/bake.
func Dir() (string, error) {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "bake"), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "bake"), nil
}

// FilePath returns the full path to config.toml.
func FilePath() (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.toml"), nil
}

// Exists reports whether a config file is already present.
func Exists() (bool, error) {
	path, err := FilePath()
	if err != nil {
		return false, err
	}
	_, err = os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return err == nil, err
}

// Load reads config.toml. Returns a helpful error if bake isn't initialized.
func Load() (*Config, error) {
	path, err := FilePath()
	if err != nil {
		return nil, err
	}
	var c Config
	if _, err := toml.DecodeFile(path, &c); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("bake is not initialized — run `bake init` first")
		}
		return nil, fmt.Errorf("read config %s: %w", path, err)
	}
	return &c, nil
}

// Save writes the config to disk (0600 — non-secret, but no reason to be lax).
func Save(c *Config) error {
	dir, err := Dir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create config dir %s: %w", dir, err)
	}
	path, err := FilePath()
	if err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("open config %s: %w", path, err)
	}
	defer f.Close()
	if err := toml.NewEncoder(f).Encode(c); err != nil {
		return fmt.Errorf("write config %s: %w", path, err)
	}
	return nil
}
