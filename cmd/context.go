package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/sunny/pmt/internal/models"
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

var contextTreeCmd = &cobra.Command{
	Use:   "tree",
	Short: "Display contexts in a tree structure",
	Long: `Display all contexts in a hierarchical tree structure.

This shows how contexts are organized with their folder-like paths.
For example, "backend/api/auth" will be shown as nested folders.`,
	Example: `  pmt context tree
  pmt ctx tree`,
	RunE: runContextTree,
}

var contextRenameCmd = &cobra.Command{
	Use:   "rename <old-context> <new-context>",
	Short: "Rename a context",
	Long: `Rename a context, updating all prompts that use it.

This command will rename a context and all of its sub-contexts.
For example, renaming "backend" to "server" will also rename:
  - "backend/api" to "server/api"
  - "backend/auth" to "server/auth"`,
	Example: `  pmt context rename backend server
  pmt context rename backend/api backend/rest
  pmt ctx rename old new`,
	Args: cobra.ExactArgs(2),
	RunE: runContextRename,
}

func init() {
	rootCmd.AddCommand(contextCmd)
	contextCmd.AddCommand(contextListCmd)
	contextCmd.AddCommand(contextTreeCmd)
	contextCmd.AddCommand(contextRenameCmd)
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

// TreeNode represents a node in the context tree
type TreeNode struct {
	Name     string
	Prompts  int
	Children map[string]*TreeNode
}

func runContextTree(cmd *cobra.Command, args []string) error {
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

	root := &TreeNode{
		Name:     "",
		Children: make(map[string]*TreeNode),
	}

	// Count prompts without context
	noContextCount := 0

	// Build the tree
	for _, p := range promptStore.Prompts {
		if p.Context == "" {
			noContextCount++
			continue
		}

		parts := p.GetContextParts()
		current := root

		for _, part := range parts {
			if current.Children[part] == nil {
				current.Children[part] = &TreeNode{
					Name:     part,
					Children: make(map[string]*TreeNode),
				}
			}
			current = current.Children[part]
			current.Prompts++
		}
	}

	// Print the tree
	var printTree func(node *TreeNode, prefix string, isLast bool, depth int)
	printTree = func(node *TreeNode, prefix string, isLast bool, depth int) {
		if depth > 0 {
			// Print current node
			connector := "├── "
			if isLast {
				connector = "└── "
			}
			fmt.Printf("%s%s%s (%d prompt%s)\n", prefix, connector, node.Name, node.Prompts, pluralize(node.Prompts))

			// Update prefix for children
			if isLast {
				prefix += "    "
			} else {
				prefix += "│   "
			}
		}

		// Sort children for consistent output
		childNames := make([]string, 0, len(node.Children))
		for name := range node.Children {
			childNames = append(childNames, name)
		}
		sort.Strings(childNames)

		// Print children
		for i, name := range childNames {
			child := node.Children[name]
			isLastChild := i == len(childNames)-1
			printTree(child, prefix, isLastChild, depth+1)
		}
	}

	// Get project name
	project := "(root)"
	if len(promptStore.Prompts) > 0 && promptStore.Prompts[0].Project != "" {
		project = promptStore.Prompts[0].Project
	}

	totalContexts := countContexts(root)
	totalWithContext := len(promptStore.Prompts) - noContextCount

	if totalContexts == 0 {
		fmt.Printf("%s/\n", project)
		if noContextCount > 0 {
			fmt.Printf("  (no context): %d prompt%s\n", noContextCount, pluralize(noContextCount))
		}
		fmt.Println("\nNo contexts found. Use 'pmt push -c <context>' to organize prompts.")
		return nil
	}

	fmt.Printf("%s/\n", project)
	printTree(root, "", false, 0)

	if noContextCount > 0 {
		fmt.Printf("\n(no context): %d prompt%s\n", noContextCount, pluralize(noContextCount))
	}

	fmt.Printf("\nTotal: %d context%s, %d prompt%s\n",
		totalContexts, pluralize(totalContexts),
		totalWithContext, pluralize(totalWithContext))

	return nil
}

// Helper function to count total contexts in tree
func countContexts(node *TreeNode) int {
	if node == nil {
		return 0
	}

	count := 0
	if node.Prompts > 0 && node.Name != "" {
		count = 1
	}

	for _, child := range node.Children {
		count += countContexts(child)
	}

	return count
}

// Helper function to pluralize words
func pluralize(count int) string {
	if count == 1 {
		return ""
	}
	return "s"
}

func runContextRename(cmd *cobra.Command, args []string) error {
	oldContext := args[0]
	newContext := args[1]

	if oldContext == newContext {
		return fmt.Errorf("old and new context names are the same")
	}

	store, err := storage.NewFileStore()
	if err != nil {
		return fmt.Errorf("failed to create store: %w", err)
	}

	// Count how many prompts will be affected
	promptStore, err := store.LoadAll()
	if err != nil {
		return fmt.Errorf("failed to load prompts: %w", err)
	}

	affectedCount := 0
	for _, p := range promptStore.Prompts {
		// Exact match or prefix match with sub-contexts
		if p.Context == oldContext || strings.HasPrefix(p.Context, oldContext+"/") {
			affectedCount++
		}
	}

	if affectedCount == 0 {
		return fmt.Errorf("no prompts found with context '%s'", oldContext)
	}

	// Perform the rename
	err = store.BulkUpdate(func(p *models.Prompt) bool {
		if p.Context == oldContext {
			// Exact match - replace entirely
			p.Context = newContext
			return true
		} else if strings.HasPrefix(p.Context, oldContext+"/") {
			// Sub-context - replace prefix
			remainder := strings.TrimPrefix(p.Context, oldContext+"/")
			if newContext == "" {
				p.Context = remainder
			} else {
				p.Context = newContext + "/" + remainder
			}
			return true
		}
		return false
	})

	if err != nil {
		return fmt.Errorf("failed to rename context: %w", err)
	}

	displayOld := oldContext
	if displayOld == "" {
		displayOld = "(no context)"
	}
	displayNew := newContext
	if displayNew == "" {
		displayNew = "(no context)"
	}

	fmt.Printf("✓ Renamed context: %s → %s\n", displayOld, displayNew)
	fmt.Printf("  Updated %d prompt%s\n", affectedCount, pluralize(affectedCount))

	return nil
}
