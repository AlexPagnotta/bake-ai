package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/alexpagnotta/bake-ai/internal/config"
	"github.com/alexpagnotta/bake-ai/internal/gooserun"
	"github.com/alexpagnotta/bake-ai/internal/tui"
	"github.com/alexpagnotta/bake-ai/internal/workspace"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bake",
	Short: "bake — your personalized, project-aware AI assistant",
	Long: `bake is a thin, provider-agnostic wrapper around goose that gives each of
your projects its own prompt, docs, and context — so you never start from zero.

Run with no arguments to open the interactive home screen.`,
	SilenceUsage: true,
	Version:      version,
	Args:         cobra.NoArgs,
	RunE:         runHome,
}

// Execute runs the root command.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.AddCommand(initCmd, newCmd, listCmd, chatCmd)
}

// runHome is the no-arg entry: the Bubble Tea picker, looping until the user
// quits or starts a chat (which hands the terminal to goose).
func runHome(cmd *cobra.Command, args []string) error {
	c, err := config.Load()
	if err != nil {
		return err
	}
	for {
		res, err := tui.RunPicker(c)
		if err != nil {
			return err
		}
		switch res.Action {
		case tui.ActionQuit:
			return nil
		case tui.ActionNew:
			name, model, ok, err := tui.NewProjectForm(c)
			if err != nil {
				return err
			}
			if !ok {
				continue
			}
			p, err := workspace.Create(c, name, model)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
			fmt.Printf("Created project %q at %s\n", p.Name, p.Path)
		case tui.ActionChat:
			p, err := workspace.Get(c, res.Project)
			if err != nil {
				return err
			}
			if !gooserun.Installed() {
				return fmt.Errorf("goose is not installed or not on PATH (brew install block-goose-cli)")
			}
			tui.PrintChatHeader(p)
			if err := gooserun.LaunchChat(p); err != nil {
				// A non-zero exit (e.g. Ctrl+C ending the session) is a normal way
				// to leave goose — just return to the picker. Surface only real
				// launch failures.
				var exitErr *exec.ExitError
				if !errors.As(err, &exitErr) {
					fmt.Fprintln(os.Stderr, "goose:", err)
				}
			}
			// Session ended — fall through and loop back to the picker.
		}
	}
}
