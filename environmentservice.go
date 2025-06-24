package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Environment represents a development environment configuration
type Environment struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	BaseURL     string `json:"baseUrl"`
	Description string `json:"description"`
	IsActive    bool   `json:"isActive"`
}

// EnvironmentService handles environment management
type EnvironmentService struct {
	environments []Environment
	configPath   string
}

// NewEnvironmentService creates a new environment service
func NewEnvironmentService() *EnvironmentService {
	homeDir, _ := os.UserHomeDir()
	configPath := filepath.Join(homeDir, ".captain-api", "environments.json")

	service := &EnvironmentService{
		environments: []Environment{},
		configPath:   configPath,
	}

	// Load existing environments or create defaults
	service.loadEnvironments()
	return service
}

// GetEnvironments returns all environments
func (e *EnvironmentService) GetEnvironments(ctx context.Context) ([]Environment, error) {
	return e.environments, nil
}

// GetActiveEnvironment returns the currently active environment
func (e *EnvironmentService) GetActiveEnvironment(ctx context.Context) (*Environment, error) {
	for _, env := range e.environments {
		if env.IsActive {
			return &env, nil
		}
	}
	return nil, fmt.Errorf("no active environment found")
}

// SetActiveEnvironment sets the active environment
func (e *EnvironmentService) SetActiveEnvironment(ctx context.Context, envID string) error {
	found := false
	for i := range e.environments {
		if e.environments[i].ID == envID {
			e.environments[i].IsActive = true
			found = true
		} else {
			e.environments[i].IsActive = false
		}
	}

	if !found {
		return fmt.Errorf("environment with ID %s not found", envID)
	}

	return e.saveEnvironments()
}

// AddEnvironment adds a new environment
func (e *EnvironmentService) AddEnvironment(ctx context.Context, env Environment) error {
	// Generate ID if not provided
	if env.ID == "" {
		env.ID = fmt.Sprintf("env_%d", len(e.environments)+1)
	}

	// Check if ID already exists
	for _, existing := range e.environments {
		if existing.ID == env.ID {
			return fmt.Errorf("environment with ID %s already exists", env.ID)
		}
	}

	// If this is the first environment, make it active
	if len(e.environments) == 0 {
		env.IsActive = true
	}

	e.environments = append(e.environments, env)
	return e.saveEnvironments()
}

// UpdateEnvironment updates an existing environment
func (e *EnvironmentService) UpdateEnvironment(ctx context.Context, env Environment) error {
	for i, existing := range e.environments {
		if existing.ID == env.ID {
			// Preserve active status if not explicitly set
			if !env.IsActive && existing.IsActive {
				env.IsActive = true
			}
			e.environments[i] = env
			return e.saveEnvironments()
		}
	}
	return fmt.Errorf("environment with ID %s not found", env.ID)
}

// DeleteEnvironment deletes an environment
func (e *EnvironmentService) DeleteEnvironment(ctx context.Context, envID string) error {
	for i, env := range e.environments {
		if env.ID == envID {
			// Don't allow deleting the last environment
			if len(e.environments) == 1 {
				return fmt.Errorf("cannot delete the last environment")
			}

			wasActive := env.IsActive

			// Remove environment
			e.environments = append(e.environments[:i], e.environments[i+1:]...)

			// If deleted environment was active, make the first one active
			if wasActive && len(e.environments) > 0 {
				e.environments[0].IsActive = true
			}

			return e.saveEnvironments()
		}
	}
	return fmt.Errorf("environment with ID %s not found", envID)
}

// loadEnvironments loads environments from config file
func (e *EnvironmentService) loadEnvironments() {
	// Create config directory if it doesn't exist
	configDir := filepath.Dir(e.configPath)
	os.MkdirAll(configDir, 0755)

	// Try to load existing config
	data, err := os.ReadFile(e.configPath)
	if err != nil {
		// Create default environments if file doesn't exist
		e.createDefaultEnvironments()
		return
	}

	var environments []Environment
	if err := json.Unmarshal(data, &environments); err != nil {
		// Create default environments if JSON is invalid
		e.createDefaultEnvironments()
		return
	}

	e.environments = environments

	// Ensure at least one environment is active
	e.ensureActiveEnvironment()
}

// saveEnvironments saves environments to config file
func (e *EnvironmentService) saveEnvironments() error {
	data, err := json.MarshalIndent(e.environments, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(e.configPath, data, 0644)
}

// createDefaultEnvironments creates default development environments
func (e *EnvironmentService) createDefaultEnvironments() {
	e.environments = []Environment{
		{
			ID:          "local",
			Name:        "Local Development",
			BaseURL:     "http://localhost:3000",
			Description: "Local development server",
			IsActive:    true,
		},
		{
			ID:          "staging",
			Name:        "Staging",
			BaseURL:     "https://api-staging.example.com",
			Description: "Staging environment for testing",
			IsActive:    false,
		},
		{
			ID:          "production",
			Name:        "Production",
			BaseURL:     "https://api.example.com",
			Description: "Production environment",
			IsActive:    false,
		},
	}
	e.saveEnvironments()
}

// ensureActiveEnvironment ensures at least one environment is active
func (e *EnvironmentService) ensureActiveEnvironment() {
	hasActive := false
	for _, env := range e.environments {
		if env.IsActive {
			hasActive = true
			break
		}
	}

	if !hasActive && len(e.environments) > 0 {
		e.environments[0].IsActive = true
		e.saveEnvironments()
	}
}
