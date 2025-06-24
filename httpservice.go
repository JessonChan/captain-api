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
}

// NewHTTPService creates a new HTTP service
func NewHTTPService() *HTTPService {
	return &HTTPService{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
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

// GetSupportedMethods returns list of supported HTTP methods
func (h *HTTPService) GetSupportedMethods() []string {
	return []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}
}