# Captain API - Code Quality & Maintainability Improvements

## 🐛 JSON Error Fix

### Problem Solved
Fixed the "invalid character 'â' looking for beginning of object key string" error that occurred when sending JSON requests.

### Root Cause
The error was caused by:
1. **Encoding Issues**: Improper handling of UTF-8 characters in JSON strings
2. **Missing Content-Type**: JSON requests weren't properly identified by the backend
3. **No Client-Side Validation**: Malformed JSON was sent to the backend without validation

### Solution Implemented

#### Backend Improvements (`httpservice.go`)
```go
// Added JSON validation and cleaning in SendRequest method
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
```

#### Frontend Improvements (`RequestBuilder.vue`)
1. **Auto Content-Type Header**: Automatically sets `Content-Type: application/json` when JSON body type is selected
2. **Real-time JSON Validation**: Validates JSON syntax as user types
3. **JSON Formatting**: One-click JSON formatting with proper indentation
4. **Visual Feedback**: Red border and error messages for invalid JSON
5. **Pre-send Validation**: Prevents sending requests with invalid JSON

---

## 🚀 Code Quality Enhancement Suggestions

### 1. **Error Handling & Logging**

#### Current State
- Basic error handling with simple error messages
- Console logging in frontend
- No structured logging in backend

#### Recommendations
```go
// Add structured logging
import "log/slog"

// In HTTPService
func (h *HTTPService) SendRequest(ctx context.Context, req HTTPRequest) (*HTTPResponse, error) {
    logger := slog.With("method", req.Method, "url", req.URL)
    logger.Info("Sending HTTP request")
    
    // ... existing code ...
    
    if err != nil {
        logger.Error("Request failed", "error", err)
        return nil, err
    }
    
    logger.Info("Request completed", "status", resp.StatusCode, "duration", duration)
    return response, nil
}
```

### 2. **Input Validation & Sanitization**

#### Add URL Validation
```go
// Enhance ValidateURL method
func (h *HTTPService) ValidateURL(url string) (string, error) {
    if url == "" {
        return "", errors.New("URL cannot be empty")
    }
    
    // Sanitize URL
    url = strings.TrimSpace(url)
    
    // Add protocol if missing
    if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
        url = "https://" + url
    }
    
    // Validate URL format
    parsedURL, err := url.Parse(url)
    if err != nil {
        return "", fmt.Errorf("invalid URL format: %w", err)
    }
    
    // Additional security checks
    if parsedURL.Host == "" {
        return "", errors.New("URL must have a valid host")
    }
    
    return url, nil
}
```

### 3. **Configuration Management**

#### Create Config Service
```go
// config.go
type Config struct {
    HTTPTimeout     time.Duration `json:"httpTimeout"`
    MaxRequestSize  int64         `json:"maxRequestSize"`
    AllowedHosts    []string      `json:"allowedHosts"`
    ProxySettings   ProxyConfig   `json:"proxySettings"`
}

type ConfigService struct {
    config *Config
}

func (c *ConfigService) GetConfig() *Config {
    return c.config
}

func (c *ConfigService) UpdateTimeout(timeout time.Duration) error {
    c.config.HTTPTimeout = timeout
    return c.saveConfig()
}
```

### 4. **Request/Response Middleware**

#### Add Interceptors
```go
// middleware.go
type RequestInterceptor func(*HTTPRequest) error
type ResponseInterceptor func(*HTTPResponse) error

type HTTPService struct {
    client               *http.Client
    requestInterceptors  []RequestInterceptor
    responseInterceptors []ResponseInterceptor
}

// Example: Request size limiter
func RequestSizeLimiter(maxSize int64) RequestInterceptor {
    return func(req *HTTPRequest) error {
        if int64(len(req.Body)) > maxSize {
            return fmt.Errorf("request body too large: %d bytes (max: %d)", len(req.Body), maxSize)
        }
        return nil
    }
}

// Example: Response time tracker
func ResponseTimeTracker() ResponseInterceptor {
    return func(resp *HTTPResponse) error {
        if resp.Duration > 5000 { // 5 seconds
            slog.Warn("Slow response detected", "duration", resp.Duration, "url", resp.Headers["X-Request-URL"])
        }
        return nil
    }
}
```

### 5. **Testing Strategy**

#### Unit Tests
```go
// httpservice_test.go
func TestHTTPService_SendRequest(t *testing.T) {
    tests := []struct {
        name           string
        request        HTTPRequest
        expectedStatus int
        expectedError  string
    }{
        {
            name: "valid JSON request",
            request: HTTPRequest{
                Method: "POST",
                URL:    "https://httpbin.org/post",
                Headers: map[string]string{"Content-Type": "application/json"},
                Body:   `{"test": "value"}`,
            },
            expectedStatus: 200,
        },
        {
            name: "invalid JSON request",
            request: HTTPRequest{
                Method: "POST",
                URL:    "https://httpbin.org/post",
                Headers: map[string]string{"Content-Type": "application/json"},
                Body:   `{"test": }`, // Invalid JSON
            },
            expectedError: "invalid JSON body",
        },
    }
    
    service := NewHTTPService()
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            resp, err := service.SendRequest(context.Background(), tt.request)
            
            if tt.expectedError != "" {
                assert.Error(t, err)
                assert.Contains(t, err.Error(), tt.expectedError)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.expectedStatus, resp.StatusCode)
            }
        })
    }
}
```

### 6. **Performance Optimizations**

#### Connection Pooling
```go
// Optimize HTTP client
func NewHTTPService() *HTTPService {
    transport := &http.Transport{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
        IdleConnTimeout:     90 * time.Second,
        DisableCompression:  false,
    }
    
    client := &http.Client{
        Transport: transport,
        Timeout:   30 * time.Second,
    }
    
    return &HTTPService{client: client}
}
```

#### Request Caching
```go
// Add response caching for GET requests
type CacheEntry struct {
    Response  *HTTPResponse
    Timestamp time.Time
    TTL       time.Duration
}

type HTTPService struct {
    client *http.Client
    cache  map[string]*CacheEntry
    mutex  sync.RWMutex
}
```

### 7. **Security Enhancements**

#### Request Sanitization
```go
// Add security headers and validation
func (h *HTTPService) sanitizeRequest(req *HTTPRequest) error {
    // Validate headers
    for key, value := range req.Headers {
        if strings.Contains(key, "\n") || strings.Contains(value, "\n") {
            return errors.New("headers cannot contain newline characters")
        }
    }
    
    // Add security headers
    if req.Headers == nil {
        req.Headers = make(map[string]string)
    }
    
    req.Headers["User-Agent"] = "Captain-API/1.0"
    
    return nil
}
```

### 8. **Frontend Improvements**

#### Component Composition
```vue
<!-- Create reusable components -->
<template>
  <div class="request-builder">
    <RequestMethodSelector v-model="request.method" />
    <URLInput v-model="request.url" @send="sendRequest" />
    <HeadersEditor v-model="headers" />
    <BodyEditor v-model="request.body" :content-type="contentType" />
    <RequestActions @send="sendRequest" @save="saveRequest" />
  </div>
</template>
```

#### State Management
```javascript
// Use Pinia for state management
import { defineStore } from 'pinia'

export const useRequestStore = defineStore('request', {
  state: () => ({
    currentRequest: null,
    requestHistory: [],
    collections: [],
    settings: {
      timeout: 30000,
      followRedirects: true,
      validateSSL: true
    }
  }),
  
  actions: {
    async sendRequest(request) {
      // Add to history
      this.requestHistory.unshift({
        ...request,
        timestamp: new Date(),
        id: crypto.randomUUID()
      })
      
      // Limit history size
      if (this.requestHistory.length > 100) {
        this.requestHistory = this.requestHistory.slice(0, 100)
      }
    }
  }
})
```

### 9. **Documentation & Code Comments**

#### Add Comprehensive Documentation
```go
// HTTPService provides HTTP client functionality for the Captain API application.
// It handles request/response processing, validation, and error handling.
//
// Example usage:
//   service := NewHTTPService()
//   response, err := service.SendRequest(ctx, HTTPRequest{
//       Method: "GET",
//       URL:    "https://api.example.com/users",
//       Headers: map[string]string{"Authorization": "Bearer token"},
//   })
type HTTPService struct {
    // client is the underlying HTTP client with configured timeouts and transport
    client *http.Client
}
```

### 10. **Monitoring & Metrics**

#### Add Request Metrics
```go
// metrics.go
type Metrics struct {
    RequestCount    int64
    ErrorCount      int64
    AverageLatency  time.Duration
    LastRequestTime time.Time
}

func (h *HTTPService) GetMetrics() *Metrics {
    return h.metrics
}

func (h *HTTPService) recordMetrics(duration time.Duration, err error) {
    h.metrics.RequestCount++
    h.metrics.LastRequestTime = time.Now()
    
    if err != nil {
        h.metrics.ErrorCount++
    }
    
    // Update average latency
    h.metrics.AverageLatency = (h.metrics.AverageLatency + duration) / 2
}
```

---

## 📋 Implementation Priority

### High Priority
1. ✅ **JSON Error Fix** - Already implemented
2. **Error Handling & Logging** - Critical for debugging
3. **Input Validation** - Security and stability
4. **Unit Testing** - Code reliability

### Medium Priority
5. **Configuration Management** - User experience
6. **Performance Optimizations** - Scalability
7. **Security Enhancements** - Data protection

### Low Priority
8. **Advanced Frontend Features** - Enhanced UX
9. **Monitoring & Metrics** - Operational insights
10. **Documentation** - Maintainability

---

## 🔧 Quick Wins

### Immediate Improvements (< 1 hour)
1. Add request timeout configuration
2. Implement basic request/response logging
3. Add URL validation with better error messages
4. Create basic unit tests for core functions

### Short-term Goals (< 1 day)
1. Implement configuration service
2. Add request size limits
3. Create comprehensive error handling
4. Add request history in frontend

### Long-term Goals (< 1 week)
1. Complete test coverage
2. Performance optimizations
3. Advanced security features
4. Monitoring and metrics dashboard

This roadmap ensures the Captain API project maintains high code quality while providing a robust and user-friendly Postman-like experience.