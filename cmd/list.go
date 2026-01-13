package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/sunny/pmt/internal/storage"
)

var (
	listType    string
	listProject string
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all prompts",
	Long: `List all saved prompts in a table format.

You can filter by type or project using flags.`,
	Example: `  pmt list
  pmt list -t bugfix
  pmt list -p my-api
  pmt list -t feature -p my-api`,
	Aliases: []string{"ls"},
	RunE:    runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&listType, "type", "t", "", "Filter by type")
	listCmd.Flags().StringVarP(&listProject, "project", "p", "", "Filter by project")
}

func runList(cmd *cobra.Command, args []string) error {
	store, err := storage.NewFileStore()
	if err != nil {
		return fmt.Errorf("failed to create store: %w", err)
	}

	// Apply filters
	filterOpts := storage.FilterOptions{
		Type:    listType,
		Project: listProject,
	}

	prompts, err := store.Filter(filterOpts)
	if err != nil {
		return fmt.Errorf("failed to load prompts: %w", err)
	}

	if len(prompts) == 0 {
		fmt.Println("No prompts found.")
		return nil
	}

	// Print header
	fmt.Printf("%-9s %-10s %-15s %-40s %s\n", "ID", "Type", "Project", "Content", "Created")
	fmt.Println(strings.Repeat("-", 100))

	// Print each prompt
	for _, p := range prompts {
		content := p.Content
		if len(content) > 40 {
			content = content[:37] + "..."
		}

		createdStr := p.CreatedAt.Format("2006-01-02 15:04")
		fmt.Printf("%-9s %-10s %-15s %-40s %s\n",
			p.ID,
			p.Type,
			truncateString(p.Project, 15),
			content,
			createdStr,
		)
	}

	fmt.Printf("\nTotal: %d prompt(s)\n", len(prompts))
	return nil
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
