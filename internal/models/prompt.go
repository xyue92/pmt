package models

import (
	"strings"
	"time"
)

// Prompt represents a saved prompt snippet
type Prompt struct {
	ID        string    `yaml:"id"`
	Name      string    `yaml:"name"`      // user-defined title/name for the prompt
	Content   string    `yaml:"content"`
	Type      string    `yaml:"type"`      // bugfix, feature, refactor, test, general
	Project   string    `yaml:"project"`   // from git detection
	Context   string    `yaml:"context"`   // user-defined context within a project (supports hierarchical paths like "backend/api/auth")
	Tags      []string  `yaml:"tags"`
	CreatedAt time.Time `yaml:"created_at"`
}

// PromptStore represents the collection of all prompts
type PromptStore struct {
	Prompts []Prompt `yaml:"prompts"`
}

// GetContextParts returns the context split into hierarchical parts
// Example: "backend/api/auth" -> ["backend", "api", "auth"]
func (p *Prompt) GetContextParts() []string {
	if p.Context == "" {
		return []string{}
	}
	return strings.Split(p.Context, "/")
}

// MatchesContextPrefix checks if the prompt's context matches the given prefix
// Example: prompt.Context="backend/api/auth" matches prefix "backend" and "backend/api"
func (p *Prompt) MatchesContextPrefix(prefix string) bool {
	if prefix == "" {
		return p.Context == ""
	}
	if p.Context == "" {
		return false
	}
	// Exact match
	if p.Context == prefix {
		return true
	}
	// Prefix match with path separator
	return strings.HasPrefix(p.Context, prefix+"/")
}

// GetContextDepth returns the depth of the context hierarchy
// Example: "backend/api/auth" returns 3, "" returns 0
func (p *Prompt) GetContextDepth() int {
	if p.Context == "" {
		return 0
	}
	return len(p.GetContextParts())
}
