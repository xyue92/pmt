package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sunny/pmt/internal/models"
	"github.com/sunny/pmt/internal/utils"
	"gopkg.in/yaml.v3"
)

// FilterOptions defines the options for filtering prompts
type FilterOptions struct {
	Type    string
	Project string
	Tags    []string
}

// Store interface defines the methods for prompt storage
type Store interface {
	Save(p *models.Prompt) error
	LoadAll() (*models.PromptStore, error)
	FindByID(id string) (*models.Prompt, error)
	Delete(id string) error
	Filter(opts FilterOptions) ([]models.Prompt, error)
}

// FileStore implements the Store interface using YAML files
type FileStore struct {
	filePath string
}

// NewFileStore creates a new FileStore instance
func NewFileStore() (*FileStore, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	pmtDir := filepath.Join(homeDir, ".pmt")
	if err := os.MkdirAll(pmtDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create .pmt directory: %w", err)
	}

	filePath := filepath.Join(pmtDir, "prompts.yaml")
	return &FileStore{filePath: filePath}, nil
}

// Save saves a prompt to the store
func (s *FileStore) Save(p *models.Prompt) error {
	store, err := s.LoadAll()
	if err != nil {
		// If file doesn't exist, create a new store
		store = &models.PromptStore{Prompts: []models.Prompt{}}
	}

	// Check for ID conflicts
	for _, existing := range store.Prompts {
		if existing.ID == p.ID {
			return fmt.Errorf("prompt with ID %s already exists", p.ID)
		}
	}

	store.Prompts = append(store.Prompts, *p)

	data, err := yaml.Marshal(store)
	if err != nil {
		return fmt.Errorf("failed to marshal prompts: %w", err)
	}

	if err := os.WriteFile(s.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write prompts file: %w", err)
	}

	return nil
}

// LoadAll loads all prompts from the store
func (s *FileStore) LoadAll() (*models.PromptStore, error) {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return &models.PromptStore{Prompts: []models.Prompt{}}, nil
		}
		return nil, fmt.Errorf("failed to read prompts file: %w", err)
	}

	var store models.PromptStore
	if err := yaml.Unmarshal(data, &store); err != nil {
		return nil, fmt.Errorf("failed to unmarshal prompts: %w", err)
	}

	return &store, nil
}

// FindByID finds a prompt by its ID or ID prefix
func (s *FileStore) FindByID(id string) (*models.Prompt, error) {
	store, err := s.LoadAll()
	if err != nil {
		return nil, err
	}

	var matches []*models.Prompt
	for i := range store.Prompts {
		if utils.MatchIDPrefix(store.Prompts[i].ID, id) {
			matches = append(matches, &store.Prompts[i])
		}
	}

	if len(matches) == 0 {
		return nil, fmt.Errorf("prompt with ID %s not found", id)
	}

	if len(matches) > 1 {
		return nil, fmt.Errorf("ambiguous ID %s: matches multiple prompts", id)
	}

	return matches[0], nil
}

// Delete deletes a prompt by its ID or ID prefix
func (s *FileStore) Delete(id string) error {
	store, err := s.LoadAll()
	if err != nil {
		return err
	}

	var matchIndex = -1
	var matchCount = 0

	for i := range store.Prompts {
		if utils.MatchIDPrefix(store.Prompts[i].ID, id) {
			matchIndex = i
			matchCount++
		}
	}

	if matchCount == 0 {
		return fmt.Errorf("prompt with ID %s not found", id)
	}

	if matchCount > 1 {
		return fmt.Errorf("ambiguous ID %s: matches multiple prompts", id)
	}

	// Remove the prompt
	store.Prompts = append(store.Prompts[:matchIndex], store.Prompts[matchIndex+1:]...)

	data, err := yaml.Marshal(store)
	if err != nil {
		return fmt.Errorf("failed to marshal prompts: %w", err)
	}

	if err := os.WriteFile(s.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write prompts file: %w", err)
	}

	return nil
}

// Filter filters prompts based on the provided options
func (s *FileStore) Filter(opts FilterOptions) ([]models.Prompt, error) {
	store, err := s.LoadAll()
	if err != nil {
		return nil, err
	}

	var filtered []models.Prompt
	for _, p := range store.Prompts {
		// Filter by type
		if opts.Type != "" && !strings.EqualFold(p.Type, opts.Type) {
			continue
		}

		// Filter by project
		if opts.Project != "" && !strings.EqualFold(p.Project, opts.Project) {
			continue
		}

		// Filter by tags
		if len(opts.Tags) > 0 {
			hasAllTags := true
			for _, filterTag := range opts.Tags {
				found := false
				for _, pTag := range p.Tags {
					if strings.EqualFold(pTag, filterTag) {
						found = true
						break
					}
				}
				if !found {
					hasAllTags = false
					break
				}
			}
			if !hasAllTags {
				continue
			}
		}

		filtered = append(filtered, p)
	}

	return filtered, nil
}
