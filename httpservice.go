package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// HTTPService handles HTTP requests for the Postman-like client
type HTTPService struct {
	client *http.Client
	envService *EnvironmentService
}

// NewHTTPService creates a new HTTP service
func NewHTTPService() *HTTPService {
	return &HTTPService{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		envService: NewEnvironmentService(),
	}
}

// NewHTTPServiceWithEnv creates a new HTTP service with a shared environment service
func NewHTTPServiceWithEnv(envService *EnvironmentService) *HTTPService {
	return &HTTPService{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		envService: envService,
	}
}

// HTTPRequest represents an HTTP request structure
type HTTPRequest struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

// HTTPResponse represents an HTTP response structure
type HTTPResponse struct {
	StatusCode int               `json:"statusCode"`
	Status     string            `json:"status"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
	Duration   int64             `json:"duration"` // in milliseconds
	Size       int64             `json:"size"`     // response size in bytes
}

// SendRequest sends an HTTP request and returns the response
func (h *HTTPService) SendRequest(ctx context.Context, req HTTPRequest) (*HTTPResponse, error) {
	start := time.Now()

	// Resolve URL with environment base URL if needed
	resolvedURL, err := h.resolveURL(ctx, req.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve URL: %w", err)
	}
	req.URL = resolvedURL

	// Validate and clean JSON body if content-type is JSON
	if req.Body != "" {
		contentType := req.Headers["Content-Type"]
		if contentType == "" {
			contentType = req.Headers["content-type"]
		}
		
		if strings.Contains(strings.ToLower(contentType), "application/json") {
			// Validate JSON format
			var jsonTest interface{}
			if err := json.Unmarshal([]byte(req.Body), &jsonTest); err != nil {
				return nil, fmt.Errorf("invalid JSON body: %w", err)
			}
			// Re-marshal to ensure clean JSON
			cleanJSON, err := json.Marshal(jsonTest)
			if err != nil {
				return nil, fmt.Errorf("failed to clean JSON: %w", err)
			}
			req.Body = string(cleanJSON)
		}
	}

	// Create HTTP request
	var bodyReader io.Reader
	if req.Body != "" {
		bodyReader = strings.NewReader(req.Body)
	}

	httpReq, err := http.NewRequestWithContext(ctx, req.Method, req.URL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	// Send request
	resp, err := h.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Calculate duration
	duration := time.Since(start).Milliseconds()

	// Convert headers to map
	headers := make(map[string]string)
	for key, values := range resp.Header {
		headers[key] = strings.Join(values, ", ")
	}

	return &HTTPResponse{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Headers:    headers,
		Body:       string(bodyBytes),
		Duration:   duration,
		Size:       int64(len(bodyBytes)),
	}, nil
}

// FormatJSON formats JSON string for better readability
func (h *HTTPService) FormatJSON(jsonStr string) (string, error) {
	var obj interface{}
	err := json.Unmarshal([]byte(jsonStr), &obj)
	if err != nil {
		return jsonStr, err // Return original if not valid JSON
	}

	formatted, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return jsonStr, err
	}

	return string(formatted), nil
}

// ValidateURL checks if the URL is valid
func (h *HTTPService) ValidateURL(url string) bool {
	if url == "" {
		return false
	}
	
	// Add http:// if no protocol specified
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}
	
	req, err := http.NewRequest("GET", url, nil)
	return err == nil && req.URL != nil
}

// GetSupportedMethods returns a list of supported HTTP methods
func (h *HTTPService) GetSupportedMethods(ctx context.Context) []string {
	return []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}
}

// resolveURL resolves a URL using the active environment's base URL if needed
func (h *HTTPService) resolveURL(ctx context.Context, url string) (string, error) {
	// If URL is already absolute, return it as is
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		return url, nil
	}
	
	// Get active environment
	activeEnv, err := h.envService.GetActiveEnvironment(ctx)
	if err != nil {
		return "", err
	}
	
	// Ensure base URL doesn't end with a slash
	baseURL := strings.TrimSuffix(activeEnv.BaseURL, "/")
	
	// Ensure URL doesn't start with a slash
	relativeURL := strings.TrimPrefix(url, "/")
	
	// Combine base URL and relative URL
	return baseURL + "/" + relativeURL, nil
}