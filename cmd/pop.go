package cmd

import (
	"fmt"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
	"github.com/sunny/pmt/internal/storage"
	"github.com/sunny/pmt/internal/ui"
)

var (
	popContext string
)

var popCmd = &cobra.Command{
	Use:   "pop",
	Short: "Select, copy, and delete a prompt",
	Long: `Interactively select a prompt from your saved prompts.

The selected prompt will be copied to your clipboard and then deleted from storage.
Similar to 'git stash pop' - use this when you want to consume the prompt.`,
	Example: `  pmt pop
  pmt pop -c backend`,
	RunE: runPop,
}

func init() {
	rootCmd.AddCommand(popCmd)
	popCmd.Flags().StringVarP(&popContext, "context", "c", "", "Filter by context")
}

func runPop(cmd *cobra.Command, args []string) error {
	store, err := storage.NewFileStore()
	if err != nil {
		return fmt.Errorf("failed to create store: %w", err)
	}

	// Apply filters
	filterOpts := storage.FilterOptions{
		Context: popContext,
	}

	prompts, err := store.Filter(filterOpts)
	if err != nil {
		return fmt.Errorf("failed to load prompts: %w", err)
	}

	if len(prompts) == 0 {
		return fmt.Errorf("no prompts available. Use 'pmt push' to add prompts")
	}

	// Show interactive selector
	selected, err := ui.SelectPrompt(prompts)
	if err != nil {
		return fmt.Errorf("selection cancelled or failed: %w", err)
	}

	// Copy to clipboard
	if err := clipboard.WriteAll(selected.Content); err != nil {
		return fmt.Errorf("failed to copy to clipboard: %w", err)
	}

	// Delete the prompt
	if err := store.Delete(selected.ID); err != nil {
		return fmt.Errorf("failed to delete prompt: %w", err)
	}

	fmt.Printf("\n✓ Copied and removed: %s\n", selected.ID)
	fmt.Println("💡 Now paste (Ctrl+V) into Copilot!")

	return nil
}
