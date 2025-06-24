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

// Collection represents a collection of requests
type Collection struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Requests    []RequestItem `json:"requests"`
	CreatedAt   time.Time     `json:"createdAt"`
	UpdatedAt   time.Time     `json:"updatedAt"`
}

// SaveRequest saves a request to a collection
func (c *CollectionService) SaveRequest(ctx context.Context, collectionID string, request RequestItem) error {
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
			ID:        collectionID,
			Name:      "Default Collection",
			Requests:  []RequestItem{},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
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

	return c.saveCollection(collection)
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
