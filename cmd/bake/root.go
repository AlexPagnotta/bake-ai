package main

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "bake",
	Short: "bake — your personalized, project-aware AI assistant",
	Long: `bake is a thin, provider-agnostic wrapper around goose that gives each of
your projects its own prompt, docs, and context — so you never start from zero.`,
	SilenceUsage: true,
	Version:      version,
}

// Execute runs the root command.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.AddCommand(initCmd, newCmd, listCmd, chatCmd)
}
