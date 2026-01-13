package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pmt",
	Short: "Prompt Manager Tool - Manage your AI prompt snippets",
	Long: `pmt is a CLI tool for managing AI prompt snippets.
Save, organize, and quickly apply your commonly used prompts for GitHub Copilot and other AI assistants.

Similar to git stash, but for your AI prompts.`,
	Version: "1.0.0",
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Global flags can be added here if needed
}
