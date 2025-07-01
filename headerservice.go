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

// AddHeaderTemplate adds a header template to a collection
func (h *HeaderService) AddHeaderTemplate(ctx context.Context, collectionID string, template HeaderTemplate) error {
	collection, err := h.GetHeaderCollection(ctx, collectionID)
	if err != nil {
		return err
	}

	if template.ID == "" {
		template.ID = fmt.Sprintf("header_tpl_%d", time.Now().UnixNano())
	}

	template.CreatedAt = time.Now()
	template.UpdatedAt = time.Now()

	collection.HeaderTemplates = append(collection.HeaderTemplates, template)
	collection.UpdatedAt = time.Now()

	return h.saveHeaderCollection(collection)
}

// UpdateHeaderTemplate updates a header template
func (h *HeaderService) UpdateHeaderTemplate(ctx context.Context, collectionID string, template HeaderTemplate) error {
	collection, err := h.GetHeaderCollection(ctx, collectionID)
	if err != nil {
		return err
	}

	for i, t := range collection.HeaderTemplates {
		if t.ID == template.ID {
			template.UpdatedAt = time.Now()
			template.CreatedAt = t.CreatedAt // Preserve original creation time
			collection.HeaderTemplates[i] = template
			collection.UpdatedAt = time.Now()
			return h.saveHeaderCollection(collection)
		}
	}

	return fmt.Errorf("header template not found")
}

// DeleteHeaderTemplate deletes a header template
func (h *HeaderService) DeleteHeaderTemplate(ctx context.Context, collectionID, templateID string) error {
	collection, err := h.GetHeaderCollection(ctx, collectionID)
	if err != nil {
		return err
	}

	for i, t := range collection.HeaderTemplates {
		if t.ID == templateID {
			collection.HeaderTemplates = append(collection.HeaderTemplates[:i], collection.HeaderTemplates[i+1:]...)
			collection.UpdatedAt = time.Now()
			return h.saveHeaderCollection(collection)
		}
	}

	return fmt.Errorf("header template not found")
}

// GetHeaderTemplates returns all header templates from a collection
func (h *HeaderService) GetHeaderTemplates(ctx context.Context, collectionID string) ([]HeaderTemplate, error) {
	collection, err := h.GetHeaderCollection(ctx, collectionID)
	if err != nil {
		return []HeaderTemplate{}, err
	}

	return collection.HeaderTemplates, nil
}

// ApplyHeaderTemplate applies a header template to existing headers
func (h *HeaderService) ApplyHeaderTemplate(ctx context.Context, collectionID, templateID string, existingHeaders map[string]string) (map[string]string, error) {
	collection, err := h.GetHeaderCollection(ctx, collectionID)
	if err != nil {
		return existingHeaders, err
	}

	for _, template := range collection.HeaderTemplates {
		if template.ID == templateID {
			// Merge template headers with existing headers
			// Existing headers take precedence
			result := make(map[string]string)

			// First add template headers
			for key, value := range template.Headers {
				result[key] = value
			}

			// Then add existing headers (overriding template headers)
			for key, value := range existingHeaders {
				result[key] = value
			}

			return result, nil
		}
	}

	return existingHeaders, fmt.Errorf("header template not found")
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
