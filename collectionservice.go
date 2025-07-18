package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

// CollectionService manages request collections
type CollectionService struct {
	collectionsPath string
}

// NewCollectionService creates a new collection service
func NewCollectionService() *CollectionService {
	// Get user's home directory
	homeDir, _ := os.UserHomeDir()
	collectionsPath := filepath.Join(homeDir, ".captain-api", "collections")

	// Create directory if it doesn't exist
	os.MkdirAll(collectionsPath, 0755)

	return &CollectionService{
		collectionsPath: collectionsPath,
	}
}

// RequestItem represents a saved request
type RequestItem struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Method      string            `json:"method"`
	URL         string            `json:"url"`
	Headers     map[string]string `json:"headers"`
	Body        string            `json:"body"`
	Description string            `json:"description"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
}

// CollectionEnvironment represents an environment within a collection
type CollectionEnvironment struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	BaseURL     string `json:"baseUrl"`
	Description string `json:"description"`
	IsActive    bool   `json:"isActive"`
}

// HeaderCollection represents a collection of header templates
type HeaderCollection struct {
	ID           string            `json:"id"`
	CollectionID string            `json:"collectionId"`
	Name         string            `json:"name"`
	Description  string            `json:"description"`
	Headers      map[string]string `json:"headers"`
	CreatedAt    time.Time         `json:"createdAt"`
	UpdatedAt    time.Time         `json:"updatedAt"`
}

// Collection represents a collection of requests
type Collection struct {
	ID                       string                  `json:"id"`
	Name                     string                  `json:"name"`
	Description              string                  `json:"description"`
	ActiveHeaderCollectionID string                  `json:"activeHeaderCollectionId,omitempty"`
	Environments             []CollectionEnvironment `json:"environments"`
	HeaderCollections        []HeaderCollection      `json:"headerCollections"`
	Requests                 []RequestItem           `json:"requests"`
	CreatedAt                time.Time               `json:"createdAt"`
	UpdatedAt                time.Time               `json:"updatedAt"`
}

// SaveRequest saves a request to a collection
func (c *CollectionService) SaveRequest(ctx context.Context, collectionID string, request RequestItem) (RequestItem, error) {
	// Generate ID if not provided
	if request.ID == "" {
		request.ID = fmt.Sprintf("req_%d", time.Now().UnixNano())
	}

	request.UpdatedAt = time.Now()
	if request.CreatedAt.IsZero() {
		request.CreatedAt = request.UpdatedAt
	}

	// Load existing collection or create new one
	collection, err := c.GetCollection(ctx, collectionID)
	if err != nil {
		// Create new collection if it doesn't exist
		collection = &Collection{
			ID:           collectionID,
			Name:         "Default Collection",
			Environments: c.createDefaultEnvironments(),
			Requests:     []RequestItem{},
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
	}

	// Check if request already exists (update) or add new
	found := false
	for i, req := range collection.Requests {
		if req.ID == request.ID {
			collection.Requests[i] = request
			found = true
			break
		}
	}

	if !found {
		collection.Requests = append(collection.Requests, request)
	}

	collection.UpdatedAt = time.Now()

	if err := c.saveCollection(collection); err != nil {
		return RequestItem{}, fmt.Errorf("failed to save collection: %w", err)
	}

	// Return the saved request
	for _, req := range collection.Requests {
		if req.ID == request.ID {
			return req, nil
		}
	}

	return request, nil
}

// UpdateRequest updates an existing request in a collection
func (c *CollectionService) UpdateRequest(ctx context.Context, collectionID string, requestID string, request RequestItem) (RequestItem, error) {
	// Load the collection
	collection, err := c.GetCollection(ctx, collectionID)
	if err != nil {
		return RequestItem{}, fmt.Errorf("failed to get collection: %w", err)
	}

	// Find the request to update
	found := false
	for i, req := range collection.Requests {
		if req.ID == requestID {
			// Update the request while preserving the created_at timestamp
			createdAt := collection.Requests[i].CreatedAt
			// Ensure we don't change the request ID
			request.ID = requestID
			collection.Requests[i] = request
			collection.Requests[i].CreatedAt = createdAt
			collection.Requests[i].UpdatedAt = time.Now()
			found = true
			break
		}
	}

	if !found {
		return RequestItem{}, fmt.Errorf("request with ID %s not found in collection %s", requestID, collectionID)
	}

	// Save the updated collection
	collection.UpdatedAt = time.Now()
	if err := c.saveCollection(collection); err != nil {
		return RequestItem{}, fmt.Errorf("failed to save collection: %w", err)
	}

	// Return the updated request
	for _, req := range collection.Requests {
		if req.ID == requestID {
			return req, nil
		}
	}

	// This should never happen since we just updated the request
	return request, nil
}

// GetCollection retrieves a collection by ID
func (c *CollectionService) GetCollection(ctx context.Context, collectionID string) (*Collection, error) {
	filePath := filepath.Join(c.collectionsPath, collectionID+".json")

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("collection not found: %w", err)
	}

	var collection Collection
	err = json.Unmarshal(data, &collection)
	if err != nil {
		return nil, fmt.Errorf("failed to parse collection: %w", err)
	}

	return &collection, nil
}

// GetAllCollections retrieves all collections
func (c *CollectionService) GetAllCollections(ctx context.Context) ([]Collection, error) {
	files, err := os.ReadDir(c.collectionsPath)
	if err != nil {
		return []Collection{}, nil // Return empty slice if directory doesn't exist
	}

	var collections []Collection

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			collectionID := strings.TrimSuffix(file.Name(), ".json")
			collection, err := c.GetCollection(ctx, collectionID)
			if err == nil {
				collections = append(collections, *collection)
			}
		}
	}

	return collections, nil
}

// DeleteRequest removes a request from a collection
func (c *CollectionService) DeleteRequest(ctx context.Context, collectionID, requestID string) error {
	collection, err := c.GetCollection(ctx, collectionID)
	if err != nil {
		return err
	}

	// Find and remove the request
	for i, req := range collection.Requests {
		if req.ID == requestID {
			collection.Requests = append(collection.Requests[:i], collection.Requests[i+1:]...)
			collection.UpdatedAt = time.Now()
			return c.saveCollection(collection)
		}
	}

	return fmt.Errorf("request not found")
}

// CreateCollection creates a new collection
func (c *CollectionService) CreateCollection(ctx context.Context, name, description string) (*Collection, error) {
	collection := &Collection{
		ID:          fmt.Sprintf("col_%d", time.Now().UnixNano()),
		Name:        name,
		Description: description,
		Requests:    []RequestItem{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := c.saveCollection(collection)
	if err != nil {
		return nil, err
	}

	return collection, nil
}

// saveCollection saves a collection to disk
func (c *CollectionService) saveCollection(collection *Collection) error {
	slices.SortFunc(collection.Requests, func(a, b RequestItem) int {
		return b.CreatedAt.Compare(a.CreatedAt) // Sort in descending order (newest first)
	})
	data, err := json.MarshalIndent(collection, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal collection: %w", err)
	}

	filePath := filepath.Join(c.collectionsPath, collection.ID+".json")
	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to save collection: %w", err)
	}

	return nil
}

// createDefaultEnvironments creates default environments for a new collection
func (c *CollectionService) createDefaultEnvironments() []CollectionEnvironment {
	return []CollectionEnvironment{
		{
			ID:          "dev",
			Name:        "Development",
			BaseURL:     "http://localhost:3000",
			Description: "Local development environment",
			IsActive:    true,
		},
		{
			ID:          "staging",
			Name:        "Staging",
			BaseURL:     "https://staging.api.example.com",
			Description: "Staging environment for testing",
			IsActive:    false,
		},
		{
			ID:          "prod",
			Name:        "Production",
			BaseURL:     "https://api.example.com",
			Description: "Production environment",
			IsActive:    false,
		},
	}
}

// GetCollectionEnvironments returns all environments for a collection
func (c *CollectionService) GetCollectionEnvironments(ctx context.Context, collectionID string) ([]CollectionEnvironment, error) {
	collection, err := c.GetCollection(ctx, collectionID)
	if err != nil {
		return nil, err
	}
	return collection.Environments, nil
}

// GetActiveCollectionEnvironment returns the active environment for a collection
func (c *CollectionService) GetActiveCollectionEnvironment(ctx context.Context, collectionID string) (*CollectionEnvironment, error) {
	collection, err := c.GetCollection(ctx, collectionID)
	if err != nil {
		return nil, err
	}

	for _, env := range collection.Environments {
		if env.IsActive {
			return &env, nil
		}
	}
	return nil, fmt.Errorf("no active environment found for collection %s", collectionID)
}

// SetActiveCollectionEnvironment sets the active environment for a collection
func (c *CollectionService) SetActiveCollectionEnvironment(ctx context.Context, collectionID, envID string) error {
	collection, err := c.GetCollection(ctx, collectionID)
	if err != nil {
		return err
	}

	found := false
	for i := range collection.Environments {
		if collection.Environments[i].ID == envID {
			collection.Environments[i].IsActive = true
			found = true
		} else {
			collection.Environments[i].IsActive = false
		}
	}

	if !found {
		return fmt.Errorf("environment with ID %s not found in collection %s", envID, collectionID)
	}

	collection.UpdatedAt = time.Now()
	return c.saveCollection(collection)
}

// AddCollectionEnvironment adds a new environment to a collection
func (c *CollectionService) AddCollectionEnvironment(ctx context.Context, collectionID string, env CollectionEnvironment) error {
	collection, err := c.GetCollection(ctx, collectionID)
	if err != nil {
		return err
	}

	// Generate ID if not provided
	if env.ID == "" {
		env.ID = fmt.Sprintf("env_%d", time.Now().UnixNano())
	}

	// Check if ID already exists
	for _, existing := range collection.Environments {
		if existing.ID == env.ID {
			return fmt.Errorf("environment with ID %s already exists in collection %s", env.ID, collectionID)
		}
	}

	// If this is the first environment, make it active
	if len(collection.Environments) == 0 {
		env.IsActive = true
	}

	collection.Environments = append(collection.Environments, env)
	collection.UpdatedAt = time.Now()
	return c.saveCollection(collection)
}

// UpdateCollectionEnvironment updates an existing environment in a collection
func (c *CollectionService) UpdateCollectionEnvironment(ctx context.Context, collectionID string, env CollectionEnvironment) error {
	collection, err := c.GetCollection(ctx, collectionID)
	if err != nil {
		return err
	}

	for i, existing := range collection.Environments {
		if existing.ID == env.ID {
			// Preserve the IsActive state
			env.IsActive = existing.IsActive
			collection.Environments[i] = env
			collection.UpdatedAt = time.Now()
			return c.saveCollection(collection)
		}
	}

	return fmt.Errorf("environment with ID %s not found in collection %s", env.ID, collectionID)
}

// DeleteCollectionEnvironment deletes an environment from a collection
func (c *CollectionService) DeleteCollectionEnvironment(ctx context.Context, collectionID, envID string) error {
	collection, err := c.GetCollection(ctx, collectionID)
	if err != nil {
		return err
	}

	// Don't allow deleting the last environment
	if len(collection.Environments) == 1 {
		return fmt.Errorf("cannot delete the last environment in collection %s", collectionID)
	}

	for i, env := range collection.Environments {
		if env.ID == envID {
			wasActive := env.IsActive

			// Remove environment
			collection.Environments = append(collection.Environments[:i], collection.Environments[i+1:]...)

			// If deleted environment was active, make the first one active
			if wasActive && len(collection.Environments) > 0 {
				collection.Environments[0].IsActive = true
			}

			collection.UpdatedAt = time.Now()
			return c.saveCollection(collection)
		}
	}

	return fmt.Errorf("environment with ID %s not found in collection %s", envID, collectionID)
}

// DeleteCollection deletes a collection
func (c *CollectionService) DeleteCollection(ctx context.Context, collectionID string) error {
	filePath := filepath.Join(c.collectionsPath, collectionID+".json")
	err := os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("failed to delete collection: %w", err)
	}
	return nil
}

// UpdateCollection updates a collection's top-level properties (e.g., name, description)
func (c *CollectionService) UpdateCollection(ctx context.Context, collectionID string, updatedCollection Collection) (*Collection, error) {
	collection, err := c.GetCollection(ctx, collectionID)
	if err != nil {
		return nil, err
	}

	// Update fields
	collection.Name = updatedCollection.Name
	collection.Description = updatedCollection.Description
	collection.UpdatedAt = time.Now()

	// Save the updated collection
	err = c.saveCollection(collection)
	if err != nil {
		return nil, err
	}

	return collection, nil
}

// SetActiveHeaderCollection sets the active header collection for a main collection
func (c *CollectionService) SetActiveHeaderCollection(ctx context.Context, collectionID string, headerCollectionID string) error {
	collection, err := c.GetCollection(ctx, collectionID)
	if err != nil {
		return err
	}

	collection.ActiveHeaderCollectionID = headerCollectionID
	collection.UpdatedAt = time.Now()

	return c.saveCollection(collection)
}

// CreateHeaderCollection creates a new header collection
func (c *CollectionService) CreateHeaderCollection(ctx context.Context, collectionID string, headerCollection HeaderCollection) (*HeaderCollection, error) {
	collection, err := c.GetCollection(ctx, collectionID)
	if err != nil {
		return nil, err
	}

	if headerCollection.ID == "" {
		headerCollection.ID = fmt.Sprintf("header_col_%d", time.Now().UnixNano())
	}

	headerCollection.CreatedAt = time.Now()
	headerCollection.UpdatedAt = time.Now()

	collection.HeaderCollections = append(collection.HeaderCollections, headerCollection)
	collection.UpdatedAt = time.Now()

	err = c.saveCollection(collection)
	if err != nil {
		return nil, err
	}
	return &headerCollection, nil
}

// UpdateHeaderCollection updates an existing header collection
func (c *CollectionService) UpdateHeaderCollection(ctx context.Context, collectionID string, headerCollection HeaderCollection) (*HeaderCollection, error) {
	collection, err := c.GetCollection(ctx, collectionID)
	if err != nil {
		return nil, err
	}

	for i, hc := range collection.HeaderCollections {
		if hc.ID == headerCollection.ID {
			headerCollection.UpdatedAt = time.Now()
			collection.HeaderCollections[i] = headerCollection
			collection.UpdatedAt = time.Now()
			err = c.saveCollection(collection)
			if err != nil {
				return nil, err
			}
			return &headerCollection, nil
		}
	}

	return nil, fmt.Errorf("header collection not found")
}

// DeleteHeaderCollection deletes a header collection
func (c *CollectionService) DeleteHeaderCollection(ctx context.Context, collectionID, headerCollectionID string) error {
	collection, err := c.GetCollection(ctx, collectionID)
	if err != nil {
		return err
	}

	for i, hc := range collection.HeaderCollections {
		if hc.ID == headerCollectionID {
			collection.HeaderCollections = append(collection.HeaderCollections[:i], collection.HeaderCollections[i+1:]...)
			collection.UpdatedAt = time.Now()
			return c.saveCollection(collection)
		}
	}

	return fmt.Errorf("header collection not found")
}
