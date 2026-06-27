package main

import (
	"fmt"

	"github.com/alexpagnotta/bake-ai/internal/config"
	"github.com/alexpagnotta/bake-ai/internal/workspace"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List projects in the workspace",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := config.Load()
		if err != nil {
			return err
		}
		ps, err := workspace.List(c)
		if err != nil {
			return err
		}
		if len(ps) == 0 {
			fmt.Println("No projects yet. Create one with: bake new <name>")
			return nil
		}
		for _, p := range ps {
			fmt.Println(p.Name)
		}
		return nil
	},
}
