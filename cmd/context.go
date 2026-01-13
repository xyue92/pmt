package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/sunny/pmt/internal/storage"
)

var contextCmd = &cobra.Command{
	Use:     "context",
	Aliases: []string{"ctx"},
	Short:   "Manage contexts",
	Long: `Manage contexts within a project.

Contexts help you organize prompts by different areas, features, or workflows
within the same git repository.`,
	Example: `  pmt context list
  pmt ctx ls`,
}

var contextListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all contexts",
	Long:    `List all contexts that have prompts associated with them.`,
	Example: `  pmt context list
  pmt ctx ls`,
	RunE: runContextList,
}

func init() {
	rootCmd.AddCommand(contextCmd)
	contextCmd.AddCommand(contextListCmd)
}

func runContextList(cmd *cobra.Command, args []string) error {
	store, err := storage.NewFileStore()
	if err != nil {
		return fmt.Errorf("failed to create store: %w", err)
	}

	// Load all prompts
	promptStore, err := store.LoadAll()
	if err != nil {
		return fmt.Errorf("failed to load prompts: %w", err)
	}

	if len(promptStore.Prompts) == 0 {
		fmt.Println("No prompts found. Use 'pmt push' to add prompts.")
		return nil
	}

	// Collect unique contexts with counts
	contextCounts := make(map[string]int)
	for _, p := range promptStore.Prompts {
		context := p.Context
		if context == "" {
			context = "(default)"
		}
		contextCounts[context]++
	}

	if len(contextCounts) == 0 {
		fmt.Println("No contexts found.")
		return nil
	}

	// Sort contexts alphabetically
	contexts := make([]string, 0, len(contextCounts))
	for ctx := range contextCounts {
		contexts = append(contexts, ctx)
	}
	sort.Strings(contexts)

	// Print header
	fmt.Printf("%-20s %s\n", "Context", "Prompts")
	fmt.Println(strings.Repeat("-", 35))

	// Print each context
	for _, ctx := range contexts {
		count := contextCounts[ctx]
		fmt.Printf("%-20s %d\n", ctx, count)
	}

	fmt.Printf("\nTotal: %d context(s)\n", len(contexts))
	return nil
}
