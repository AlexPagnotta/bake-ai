package main

import (
	"fmt"
	"os"

	"github.com/alexpagnotta/bake-ai/internal/config"
	"github.com/alexpagnotta/bake-ai/internal/gooserun"
	"github.com/alexpagnotta/bake-ai/internal/workspace"
	"github.com/spf13/cobra"
)

var (
	initWorkspace string
	initModel     string
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize bake (workspace + non-secret config)",
	RunE: func(cmd *cobra.Command, args []string) error {
		exists, err := config.Exists()
		if err != nil {
			return err
		}

		var c *config.Config
		if exists {
			if c, err = config.Load(); err != nil {
				return err
			}
			fmt.Println("bake is already initialized.")
		} else {
			if c, err = config.Default(); err != nil {
				return err
			}
			if initWorkspace != "" {
				c.WorkspacePath = initWorkspace
			}
			if initModel != "" {
				c.DefaultModel = initModel
			}
			if err := os.MkdirAll(workspace.ProjectsDir(c), 0o755); err != nil {
				return err
			}
			if err := config.Save(c); err != nil {
				return err
			}
			path, _ := config.FilePath()
			fmt.Printf("Initialized bake.\n  config:    %s\n", path)
		}
		fmt.Printf("  workspace: %s\n  model:     %s\n", c.WorkspacePath, c.DefaultModel)

		// goose is where secrets live; we only check, never store them.
		switch {
		case !gooserun.Installed():
			fmt.Println("\n⚠  goose is not on your PATH. Install it: brew install block-goose-cli")
		case !gooserun.HasOpenRouter():
			fmt.Println("\n⚠  No OpenRouter provider detected in goose.\n   Run: goose configure  →  Configure Providers  →  OpenRouter")
		default:
			fmt.Println("\n✓ goose + OpenRouter detected. Try: bake new <name>")
		}
		return nil
	},
}

func init() {
	initCmd.Flags().StringVar(&initWorkspace, "workspace", "", "workspace path (default ~/bake)")
	initCmd.Flags().StringVar(&initModel, "model", "", "default model for new projects")
}
