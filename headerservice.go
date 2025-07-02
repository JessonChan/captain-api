package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// HeaderService manages header templates and collections
type HeaderService struct {
	headersPath string
}

// NewHeaderService creates a new header service
func NewHeaderService() *HeaderService {
	// Get user's home directory
	homeDir, _ := os.UserHomeDir()
	headersPath := filepath.Join(homeDir, ".captain-api", "headers")

	// Create directory if it doesn't exist
	os.MkdirAll(headersPath, 0755)

	return &HeaderService{
		headersPath: headersPath,
	}
}

// GetHeaderCollections returns all header collections for a given collection
func (h *HeaderService) GetHeaderCollections(ctx context.Context, collectionID string) ([]HeaderCollection, error) {
	files, err := os.ReadDir(h.headersPath)
	if err != nil {
		return []HeaderCollection{}, nil
	}

	var collections []HeaderCollection
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".json" {
			collection, err := h.loadHeaderCollection(file.Name())
			if err != nil {
				continue
			}
			if collection.CollectionID == collectionID {
				collections = append(collections, *collection)
			}
		}
	}

	return collections, nil
}

// GetHeaderCollection returns a specific header collection
func (h *HeaderService) GetHeaderCollection(ctx context.Context, collectionID string) (*HeaderCollection, error) {
	return h.loadHeaderCollection(collectionID + ".json")
}

// CreateHeaderCollection creates a new header collection
func (h *HeaderService) CreateHeaderCollection(ctx context.Context, collection HeaderCollection) (*HeaderCollection, error) {
	if collection.ID == "" {
		collection.ID = fmt.Sprintf("header_col_%d", time.Now().UnixNano())
	}

	collection.CreatedAt = time.Now()
	collection.UpdatedAt = time.Now()

	err := h.saveHeaderCollection(&collection)
	if err != nil {
		return nil, err
	}
	return &collection, nil
}

// UpdateHeaderCollection updates an existing header collection
func (h *HeaderService) UpdateHeaderCollection(ctx context.Context, collection HeaderCollection) (*HeaderCollection, error) {
	existing, err := h.GetHeaderCollection(ctx, collection.ID)
	if err != nil {
		return nil, err
	}
	collection.CollectionID = existing.CollectionID
	collection.UpdatedAt = time.Now()
	err = h.saveHeaderCollection(&collection)
	if err != nil {
		return nil, err
	}
	return &collection, nil
}

// DeleteHeaderCollection deletes a header collection
func (h *HeaderService) DeleteHeaderCollection(ctx context.Context, collectionID string) error {
	filePath := filepath.Join(h.headersPath, collectionID+".json")
	return os.Remove(filePath)
}

// loadHeaderCollection loads a header collection from file
func (h *HeaderService) loadHeaderCollection(filename string) (*HeaderCollection, error) {
	filePath := filepath.Join(h.headersPath, filename)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var collection HeaderCollection
	err = json.Unmarshal(data, &collection)
	if err != nil {
		return nil, err
	}

	return &collection, nil
}

// saveHeaderCollection saves a header collection to file
func (h *HeaderService) saveHeaderCollection(collection *HeaderCollection) error {
	filePath := filepath.Join(h.headersPath, collection.ID+".json")
	data, err := json.MarshalIndent(collection, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}
