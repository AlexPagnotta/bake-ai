package main

import (
	"fmt"

	"github.com/alexpagnotta/bake-ai/internal/config"
	"github.com/alexpagnotta/bake-ai/internal/gooserun"
	"github.com/alexpagnotta/bake-ai/internal/tui"
	"github.com/alexpagnotta/bake-ai/internal/workspace"
	"github.com/spf13/cobra"
)

var chatCmd = &cobra.Command{
	Use:   "chat <project>",
	Short: "Start a chat session for a project (hands the terminal to goose)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := config.Load()
		if err != nil {
			return err
		}
		p, err := workspace.Get(c, args[0])
		if err != nil {
			return err
		}
		if !gooserun.Installed() {
			return fmt.Errorf("goose is not installed or not on PATH (brew install block-goose-cli)")
		}
		tui.PrintChatHeader(p)
		return gooserun.LaunchChat(p)
	},
}
