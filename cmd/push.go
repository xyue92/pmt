package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/sunny/pmt/internal/models"
	"github.com/sunny/pmt/internal/storage"
	"github.com/sunny/pmt/internal/utils"
)

var (
	pushType    string
	pushName    string
	pushContext string
	pushTags    []string
)

var pushCmd = &cobra.Command{
	Use:   "push [content]",
	Short: "Save a new prompt",
	Long: `Save a new prompt snippet to your local store.

If no content is provided, an editor will open for you to write a longer prompt.
The prompt will be tagged with the current git project automatically.
You can optionally specify a type and tags.`,
	Example: `  pmt push "Fix memory leak in async handler"
  pmt push "Add OAuth login" -t feature --tags auth,api
  pmt push "Refactor error handling" -t refactor
  pmt push   # Opens editor for longer prompts`,
	RunE: runPush,
}

func init() {
	rootCmd.AddCommand(pushCmd)
	pushCmd.Flags().StringVarP(&pushType, "type", "t", "general", "Type: bugfix, feature, refactor, test, general")
	pushCmd.Flags().StringVarP(&pushName, "name", "n", "", "Custom name/title for the prompt")
	pushCmd.Flags().StringVarP(&pushContext, "context", "c", "", "Context within the project")
	pushCmd.Flags().StringSliceVarP(&pushTags, "tags", "g", []string{}, "Tags (comma-separated)")
}

func runPush(cmd *cobra.Command, args []string) error {
	var content string

	// If no args provided, open editor
	if len(args) == 0 {
		var err error
		content, err = openEditor()
		if err != nil {
			return fmt.Errorf("failed to open editor: %w", err)
		}
	} else {
		content = strings.Join(args, " ")
	}

	// Trim whitespace
	content = strings.TrimSpace(content)

	if content == "" {
		return fmt.Errorf("prompt content cannot be empty")
	}

	// Validate type
	validTypes := map[string]bool{
		"bugfix":   true,
		"feature":  true,
		"refactor": true,
		"test":     true,
		"general":  true,
	}

	if !validTypes[pushType] {
		return fmt.Errorf("invalid type: %s (must be bugfix, feature, refactor, test, or general)", pushType)
	}

	// Create the store
	store, err := storage.NewFileStore()
	if err != nil {
		return fmt.Errorf("failed to create store: %w", err)
	}

	// Create the prompt
	prompt := &models.Prompt{
		ID:        utils.GenerateID(),
		Name:      pushName,
		Content:   content,
		Type:      pushType,
		Project:   utils.DetectGitProject(),
		Context:   pushContext,
		Tags:      pushTags,
		CreatedAt: time.Now(),
	}

	// Save the prompt
	if err := store.Save(prompt); err != nil {
		return fmt.Errorf("failed to save prompt: %w", err)
	}

	fmt.Printf("✓ Saved prompt: %s (%s) in project: %s\n", prompt.ID, prompt.Type, prompt.Project)
	return nil
}

// openEditor opens the user's preferred editor to write a prompt
func openEditor() (string, error) {
	// Get editor from environment, default to vim
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = os.Getenv("VISUAL")
	}
	if editor == "" {
		editor = "vim"
	}

	// Create a temporary file
	tmpDir := os.TempDir()
	tmpFile := filepath.Join(tmpDir, fmt.Sprintf("pmt-prompt-%d.md", time.Now().Unix()))

	// Write initial template
	template := `# Write your prompt below this line
# Lines starting with # will be ignored
# Save and close the editor to save the prompt

`
	if err := os.WriteFile(tmpFile, []byte(template), 0644); err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile)

	// Open editor
	cmd := exec.Command(editor, tmpFile)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("editor exited with error: %w", err)
	}

	// Read the file
	data, err := os.ReadFile(tmpFile)
	if err != nil {
		return "", fmt.Errorf("failed to read temp file: %w", err)
	}

	// Parse content (remove comment lines)
	lines := strings.Split(string(data), "\n")
	var contentLines []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "#") && trimmed != "" {
			contentLines = append(contentLines, line)
		}
	}

	return strings.TrimSpace(strings.Join(contentLines, "\n")), nil
}
