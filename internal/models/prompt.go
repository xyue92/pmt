package models

import "time"

// Prompt represents a saved prompt snippet
type Prompt struct {
	ID        string    `yaml:"id"`
	Content   string    `yaml:"content"`
	Type      string    `yaml:"type"`      // bugfix, feature, refactor, general
	Project   string    `yaml:"project"`   // from git detection
	Tags      []string  `yaml:"tags"`
	CreatedAt time.Time `yaml:"created_at"`
}

// PromptStore represents the collection of all prompts
type PromptStore struct {
	Prompts []Prompt `yaml:"prompts"`
}
