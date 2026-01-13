package ui

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/sunny/pmt/internal/models"
)

// SelectPrompt displays an interactive prompt selector and returns the selected prompt
func SelectPrompt(prompts []models.Prompt) (*models.Prompt, error) {
	if len(prompts) == 0 {
		return nil, fmt.Errorf("no prompts available")
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "▸ {{ .ID | cyan }} ({{ .Type | yellow }}) {{ .Content | truncate 60 }}",
		Inactive: "  {{ .ID | cyan }} ({{ .Type | yellow }}) {{ .Content | truncate 60 }}",
		Selected: "✓ Selected: {{ .ID | cyan }}",
		Details: `
--------- Details ----------
ID:       {{ .ID }}
Type:     {{ .Type }}
Project:  {{ .Project }}
Created:  {{ .CreatedAt.Format "2006-01-02 15:04" }}
Tags:     {{ joinTags .Tags }}
Content:
{{ .Content }}`,
	}

	// Custom function map for templates
	templates.FuncMap = promptui.FuncMap
	templates.FuncMap["truncate"] = func(length int, s string) string {
		if len(s) <= length {
			return s
		}
		return s[:length] + "..."
	}
	templates.FuncMap["joinTags"] = func(tags []string) string {
		if len(tags) == 0 {
			return "(none)"
		}
		return strings.Join(tags, ", ")
	}

	searcher := func(input string, index int) bool {
		prompt := prompts[index]
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		// Search in ID, type, content, and tags
		searchStr := strings.ToLower(prompt.ID + prompt.Type + prompt.Content + strings.Join(prompt.Tags, " "))
		searchStr = strings.Replace(searchStr, " ", "", -1)
		return strings.Contains(searchStr, input)
	}

	prompt := promptui.Select{
		Label:     "Select Prompt",
		Items:     prompts,
		Templates: templates,
		Size:      10,
		Searcher:  searcher,
	}

	i, _, err := prompt.Run()
	if err != nil {
		return nil, err
	}

	return &prompts[i], nil
}
