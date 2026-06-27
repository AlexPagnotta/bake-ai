package main

import (
	"fmt"

	"github.com/alexpagnotta/bake-ai/internal/config"
	"github.com/alexpagnotta/bake-ai/internal/workspace"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new <project>",
	Short: "Scaffold a new project in the workspace",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := config.Load()
		if err != nil {
			return err
		}
		p, err := workspace.Create(c, args[0])
		if err != nil {
			return err
		}
		fmt.Printf("Created project %q at %s\n", p.Name, p.Path)
		fmt.Printf("Next: add context in %s/.goosehints, then run `bake chat %s`\n", p.Path, p.Name)
		return nil
	},
}
