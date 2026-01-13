package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/sunny/pmt/internal/models"
	"github.com/sunny/pmt/internal/storage"
)

var (
	mvContext string
)

var mvCmd = &cobra.Command{
	Use:   "mv <id> -c <new-context>",
	Short: "Move a prompt to a different context",
	Long: `Move a prompt to a different context (folder).

This allows you to reorganize your prompts by changing their context.
You can use hierarchical paths like "backend/api/auth".`,
	Example: `  pmt mv a7f -c backend/api
  pmt mv abc123 -c frontend/components
  pmt mv def -c ""  # Remove context`,
	Args: cobra.ExactArgs(1),
	RunE: runMv,
}

func init() {
	rootCmd.AddCommand(mvCmd)
	mvCmd.Flags().StringVarP(&mvContext, "context", "c", "", "New context path (required)")
	mvCmd.MarkFlagRequired("context")
}

func runMv(cmd *cobra.Command, args []string) error {
	id := args[0]

	store, err := storage.NewFileStore()
	if err != nil {
		return fmt.Errorf("failed to create store: %w", err)
	}

	// Find the prompt first to show what we're moving
	prompt, err := store.FindByID(id)
	if err != nil {
		return err
	}

	oldContext := prompt.Context
	if oldContext == "" {
		oldContext = "(no context)"
	}

	newContext := mvContext
	displayNewContext := newContext
	if displayNewContext == "" {
		displayNewContext = "(no context)"
	}

	// Update the prompt's context
	err = store.Update(id, func(p *models.Prompt) {
		p.Context = newContext
	})
	if err != nil {
		return fmt.Errorf("failed to move prompt: %w", err)
	}

	fmt.Printf("✓ Moved prompt %s\n", prompt.ID)
	fmt.Printf("  From: %s\n", oldContext)
	fmt.Printf("  To:   %s\n", displayNewContext)

	return nil
}
