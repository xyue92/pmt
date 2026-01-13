package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/sunny/pmt/internal/storage"
)

var showCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show detailed information about a prompt",
	Long: `Display the complete details of a specific prompt.

You can use the full ID or just a prefix (e.g., 'a7f' instead of 'a7f3c2b').`,
	Example: `  pmt show a7f3c2b
  pmt show a7f`,
	Args: cobra.ExactArgs(1),
	RunE: runShow,
}

func init() {
	rootCmd.AddCommand(showCmd)
}

func runShow(cmd *cobra.Command, args []string) error {
	id := args[0]

	store, err := storage.NewFileStore()
	if err != nil {
		return fmt.Errorf("failed to create store: %w", err)
	}

	// Find the prompt
	prompt, err := store.FindByID(id)
	if err != nil {
		return err
	}

	// Display the prompt details
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("ID:        %s\n", prompt.ID)
	fmt.Printf("Type:      %s\n", prompt.Type)
	fmt.Printf("Project:   %s\n", prompt.Project)

	if prompt.Context != "" {
		fmt.Printf("Context:   %s\n", prompt.Context)
	}

	if len(prompt.Tags) > 0 {
		fmt.Printf("Tags:      %s\n", strings.Join(prompt.Tags, ", "))
	} else {
		fmt.Println("Tags:      (none)")
	}

	fmt.Printf("Created:   %s\n", prompt.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("\nContent:")
	fmt.Println(prompt.Content)
	fmt.Println()

	return nil
}
