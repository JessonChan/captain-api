# Captain API - Postman-like Client

A modern API testing client built with Wails (Go + Vue3) that provides a Postman-like interface for testing HTTP APIs.

## Features

### 🚀 Core Features
- **HTTP Request Builder**: Support for GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS methods
- **Request Management**: Save, organize, and reuse requests in collections
- **Response Viewer**: Beautiful response display with syntax highlighting
- **Headers Management**: Easy-to-use interface for managing request headers
- **Body Support**: JSON, text, and raw body types
- **Collections**: Organize requests into collections for better management
- **Persistent Storage**: Requests are saved locally for future use

### 🎨 UI Features
- **Modern Interface**: Clean, responsive design inspired by Postman
- **Tabbed Interface**: Separate tabs for headers, body, and response
- **Real-time Response**: Live response display with status codes, timing, and size
- **JSON Formatting**: Automatic JSON formatting and validation
- **Copy to Clipboard**: Easy copying of response data

## Getting Started

### Prerequisites
- Go 1.22+ 
- Node.js 18+
- Bun (package manager)
- Wails v3

### Installation

1. **Clone the repository** (if not already done):
   ```bash
   git clone <your-repo-url>
   cd captain-api
   ```

2. **Install dependencies**:
   ```bash
   # Install frontend dependencies
   cd frontend && bun install && cd ..
   
   # Install Go dependencies
   go mod tidy
   ```

3. **Generate bindings**:
   ```bash
   task darwin:common:generate:bindings
   ```

4. **Build the application**:
   ```bash
   task darwin:build
   ```

### Development

To run the application in development mode with hot reload:

```bash
task dev
```

This will start the application with:
- Hot reload for frontend changes
- Automatic Go code recompilation
- Debug console access

### Production Build

To create a production build:

```bash
task darwin:build
```

The built application will be available in the `bin/` directory.

## Usage Guide

### Making Your First Request

1. **Start the application**:
   ```bash
   ./bin/captain-api
   ```

2. **Build a request**:
   - Select HTTP method (GET, POST, etc.)
   - Enter the URL (e.g., `https://jsonplaceholder.typicode.com/posts/1`)
   - Add headers if needed (Content-Type, Authorization, etc.)
   - Add request body for POST/PUT requests

3. **Send the request**:
   - Click the "Send" button
   - View the response in the response panel

4. **Save the request**:
   - Enter a name for your request
   - Click "Save Request" to add it to your collection

### Managing Collections

1. **Create a new collection**:
   - Click the "+" button in the Collections sidebar
   - Enter collection name and description
   - Click "Create"

2. **Load saved requests**:
   - Expand a collection in the sidebar
   - Click on any saved request to load it
   - The request will populate in the request builder

3. **Delete requests**:
   - Hover over a request in the sidebar
   - Click the "×" button to delete

### Response Analysis

The response viewer provides:
- **Status Code**: Color-coded status (green for 2xx, red for 4xx/5xx)
- **Response Time**: Request duration in milliseconds
- **Response Size**: Size of the response body
- **Headers Tab**: All response headers
- **Body Tab**: Formatted response body with JSON syntax highlighting
- **Raw Tab**: Complete raw HTTP response

## Project Structure

```
captain-api/
├── main.go                    # Application entry point
├── httpservice.go             # HTTP client service
├── collectionservice.go       # Collection management service
├── greetservice.go           # Example service (can be removed)
├── frontend/
│   ├── src/
│   │   ├── App.vue           # Main application component
│   │   └── components/
│   │       ├── RequestBuilder.vue    # Request building interface
│   │       ├── ResponseViewer.vue    # Response display
│   │       └── CollectionSidebar.vue # Collections management
│   ├── bindings/             # Auto-generated Go-JS bindings
│   └── dist/                 # Built frontend assets
├── build/                    # Build configuration
└── bin/                      # Built executables
```

## API Services

### HTTPService
- `SendRequest(request)`: Send HTTP request and return response
- `FormatJSON(jsonString)`: Format JSON for better readability
- `ValidateURL(url)`: Validate URL format
- `GetSupportedMethods()`: Get list of supported HTTP methods

### CollectionService
- `SaveRequest(collectionID, request)`: Save request to collection
- `GetCollection(collectionID)`: Retrieve collection by ID
- `GetAllCollections()`: Get all collections
- `DeleteRequest(collectionID, requestID)`: Delete request from collection
- `CreateCollection(name, description)`: Create new collection

## Data Storage

Collections and requests are stored locally in:
- **macOS**: `~/.captain-api/collections/`
- **Windows**: `%USERPROFILE%\.captain-api\collections\`
- **Linux**: `~/.captain-api/collections/`

Each collection is stored as a JSON file with the collection ID as the filename.

## Customization

### Adding New Features

1. **Backend (Go)**:
   - Add new methods to existing services
   - Create new services in separate files
   - Register services in `main.go`

2. **Frontend (Vue3)**:
   - Create new components in `frontend/src/components/`
   - Add new views or modify existing ones
   - Use the auto-generated bindings to call Go functions

3. **Regenerate bindings** after Go changes:
   ```bash
   task darwin:common:generate:bindings
   ```

### Styling

The application uses:
- **CSS Variables**: For consistent theming
- **Responsive Design**: Mobile-friendly layout
- **Modern CSS**: Flexbox, Grid, and CSS3 features

## Troubleshooting

### Common Issues

1. **Build Errors**:
   - Ensure all dependencies are installed
   - Run `go mod tidy` to resolve Go dependencies
   - Run `bun install` in the frontend directory

2. **Binding Errors**:
   - Regenerate bindings: `task darwin:common:generate:bindings`
   - Ensure service methods are exported (capitalized)

3. **CORS Issues**:
   - The application runs as a desktop app, so CORS shouldn't be an issue
   - If testing local APIs, ensure they accept requests from the app

4. **Request Failures**:
   - Check URL format (include http:// or https://)
   - Verify network connectivity
   - Check if the target API is accessible

### Debug Mode

Run in development mode for debugging:
```bash
task dev
```

This provides:
- Browser developer tools access
- Console logging
- Hot reload for faster development

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Built with [Wails](https://wails.io/) - Go + Web frontend framework
- UI inspired by [Postman](https://www.postman.com/)
- Icons and styling inspired by modern API tools

---

**Happy API Testing! 🚀**