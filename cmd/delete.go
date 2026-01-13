package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/sunny/pmt/internal/storage"
)

var deleteForce bool

var deleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a prompt",
	Long: `Delete a specific prompt by its ID.

You can use the full ID or just a prefix (e.g., 'a7f' instead of 'a7f3c2b').
By default, you will be asked to confirm the deletion.`,
	Example: `  pmt delete a7f3c2b
  pmt delete a7f
  pmt delete a7f -f  # Force delete without confirmation`,
	Aliases: []string{"rm"},
	Args:    cobra.ExactArgs(1),
	RunE:    runDelete,
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().BoolVarP(&deleteForce, "force", "f", false, "Force deletion without confirmation")
}

func runDelete(cmd *cobra.Command, args []string) error {
	id := args[0]

	store, err := storage.NewFileStore()
	if err != nil {
		return fmt.Errorf("failed to create store: %w", err)
	}

	// Find the prompt first to show what will be deleted
	prompt, err := store.FindByID(id)
	if err != nil {
		return err
	}

	// Ask for confirmation unless force flag is set
	if !deleteForce {
		fmt.Printf("Delete prompt %s? (%s)\n", prompt.ID, truncateString(prompt.Content, 50))
		fmt.Print("Type 'yes' to confirm: ")

		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read confirmation: %w", err)
		}

		response = strings.TrimSpace(strings.ToLower(response))
		if response != "yes" && response != "y" {
			fmt.Println("Deletion cancelled.")
			return nil
		}
	}

	// Delete the prompt
	if err := store.Delete(prompt.ID); err != nil {
		return fmt.Errorf("failed to delete prompt: %w", err)
	}

	fmt.Printf("✓ Deleted prompt: %s\n", prompt.ID)
	return nil
}
