package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// RequestLog represents a logged HTTP request/response pair
type RequestLog struct {
	ID        string         `json:"id"`
	Method    string         `json:"method"`
	URL       string         `json:"url"`
	Status    int            `json:"status"`
	Timestamp time.Time      `json:"timestamp"`
	Duration  int64          `json:"duration"` // in milliseconds
	Request   LoggedRequest  `json:"request"`
	Response  LoggedResponse `json:"response"`
}

// LoggedRequest represents the request part of a log entry
type LoggedRequest struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

// LoggedResponse represents the response part of a log entry
type LoggedResponse struct {
	StatusCode int               `json:"statusCode"`
	Status     string            `json:"status"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
	Size       int64             `json:"size"`
}

// LogService manages request/response logs
type LogService struct {
	mu      sync.RWMutex
	logs    []RequestLog
	logsDir string
}

// NewLogService creates a new log service
func NewLogService() *LogService {
	// Get user's home directory
	homeDir, _ := os.UserHomeDir()
	logsDir := filepath.Join(homeDir, ".captain-api", "logs")

	// Create directory if it doesn't exist
	err := os.MkdirAll(logsDir, 0755)
	if err != nil {
	}

	service := &LogService{
		logs:    make([]RequestLog, 0),
		logsDir: logsDir,
	}

	// Load existing logs from disk
	service.loadLogsFromDisk()

	return service
}

// LogRequest adds a new request/response log entry
func (l *LogService) LogRequest(ctx context.Context, req HTTPRequest, resp *HTTPResponse, duration int64) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Generate a simple ID based on timestamp
	id := time.Now().Format("20060102150405.000")

	log := RequestLog{
		ID:        id,
		Method:    req.Method,
		URL:       req.URL,
		Status:    resp.StatusCode,
		Timestamp: time.Now(),
		Duration:  duration,
		Request: LoggedRequest{
			Method:  req.Method,
			URL:     req.URL,
			Headers: copyHeaders(req.Headers),
			Body:    req.Body,
		},
		Response: LoggedResponse{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
			Headers:    copyHeaders(resp.Headers),
			Body:       resp.Body,
			Size:       resp.Size,
		},
	}

	// Add to the beginning of the slice (most recent first)
	l.logs = append([]RequestLog{log}, l.logs...)

	// Keep only the last 100 logs to prevent memory issues
	if len(l.logs) > 100 {
		l.logs = l.logs[:100]
	}

	// Save logs to disk
	l.saveLogsToDisk()

	return nil
}

// GetAllLogs returns all logged requests
func (l *LogService) GetAllLogs(ctx context.Context) ([]RequestLog, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	// Return a copy of the logs to prevent external modification
	logsCopy := make([]RequestLog, len(l.logs))
	copy(logsCopy, l.logs)

	return logsCopy, nil
}

// GetLogByID returns a specific log by ID
func (l *LogService) GetLogByID(ctx context.Context, id string) (*RequestLog, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	for _, log := range l.logs {
		if log.ID == id {
			// Return a copy
			logCopy := log
			return &logCopy, nil
		}
	}

	return nil, nil // Not found
}

// ClearLogs removes all logged requests
func (l *LogService) ClearLogs(ctx context.Context) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.logs = make([]RequestLog, 0)
	// Clear logs from disk
	l.saveLogsToDisk()
	return nil
}

// GetLogsCount returns the number of logged requests
func (l *LogService) GetLogsCount(ctx context.Context) (int, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return len(l.logs), nil
}

// ExportLogsAsJSON exports all logs as JSON string
func (l *LogService) ExportLogsAsJSON(ctx context.Context) (string, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	jsonData, err := json.MarshalIndent(l.logs, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

// copyHeaders creates a deep copy of headers map
func copyHeaders(headers map[string]string) map[string]string {
	if headers == nil {
		return nil
	}

	copy := make(map[string]string, len(headers))
	for k, v := range headers {
		copy[k] = v
	}
	return copy
}

// loadLogsFromDisk loads logs from the disk storage
func (l *LogService) loadLogsFromDisk() {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Path to the logs file
	logsFile := filepath.Join(l.logsDir, "request_logs.json")

	// Check if file exists
	if _, err := os.Stat(logsFile); os.IsNotExist(err) {
		// File doesn't exist, nothing to load
		return
	}

	// Read file
	data, err := os.ReadFile(logsFile)
	if err != nil {
		fmt.Printf("Error reading logs file: %v\n", err)
		return
	}

	// Parse JSON
	var logs []RequestLog
	if err := json.Unmarshal(data, &logs); err != nil {
		fmt.Printf("Error parsing logs file: %v\n", err)
		return
	}

	// Set logs
	l.logs = logs
}

// saveLogsToDisk saves logs to the disk storage
func (l *LogService) saveLogsToDisk() {
	// Path to the logs file
	logsFile := filepath.Join(l.logsDir, "request_logs.json")

	// Convert logs to JSON
	data, err := json.MarshalIndent(l.logs, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling logs: %v\n", err)
		return
	}

	// Write to file
	if err := os.WriteFile(logsFile, data, 0644); err != nil {
		fmt.Printf("Error writing logs file: %v\n", err)
	}
}
