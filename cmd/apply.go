package cmd

import (
	"fmt"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
	"github.com/sunny/pmt/internal/storage"
	"github.com/sunny/pmt/internal/ui"
)

var (
	applyContext string
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Select and copy a prompt to clipboard",
	Long: `Interactively select a prompt from your saved prompts.

The selected prompt will be copied to your clipboard automatically.
Use arrow keys to navigate and press Enter to select.
Press / to search.`,
	Example: `  pmt apply
  pmt apply -c backend`,
	RunE: runApply,
}

func init() {
	rootCmd.AddCommand(applyCmd)
	applyCmd.Flags().StringVarP(&applyContext, "context", "c", "", "Filter by context")
}

func runApply(cmd *cobra.Command, args []string) error {
	store, err := storage.NewFileStore()
	if err != nil {
		return fmt.Errorf("failed to create store: %w", err)
	}

	// Apply filters
	filterOpts := storage.FilterOptions{
		Context: applyContext,
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

	fmt.Printf("\n✓ Copied to clipboard: %s\n", selected.ID)
	fmt.Println("💡 Now paste (Ctrl+V) into Copilot!")

	return nil
}
