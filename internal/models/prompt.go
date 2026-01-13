package models

import "time"

// Prompt represents a saved prompt snippet
type Prompt struct {
	ID        string    `yaml:"id"`
	Content   string    `yaml:"content"`
	Type      string    `yaml:"type"`      // bugfix, feature, refactor, test, general
	Project   string    `yaml:"project"`   // from git detection
	Context   string    `yaml:"context"`   // user-defined context within a project
	Tags      []string  `yaml:"tags"`
	CreatedAt time.Time `yaml:"created_at"`
}

// PromptStore represents the collection of all prompts
type PromptStore struct {
	Prompts []Prompt `yaml:"prompts"`
}
