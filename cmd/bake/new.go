package main

import (
	"fmt"

	"github.com/alexpagnotta/bake-ai/internal/config"
	"github.com/alexpagnotta/bake-ai/internal/tui"
	"github.com/alexpagnotta/bake-ai/internal/workspace"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new [project]",
	Short: "Scaffold a new project (interactive form if no name is given)",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := config.Load()
		if err != nil {
			return err
		}

		var name, model string
		if len(args) == 1 {
			name = args[0] // headless path stays scriptable
		} else {
			var ok bool
			if name, model, ok, err = tui.NewProjectForm(c); err != nil {
				return err
			} else if !ok {
				return nil // user aborted
			}
		}

		p, err := workspace.Create(c, name, model)
		if err != nil {
			return err
		}
		fmt.Printf("Created project %q at %s\n", p.Name, p.Path)
		fmt.Printf("Next: add context in %s/.goosehints, then run `bake chat %s`\n", p.Path, p.Name)
		return nil
	},
}
