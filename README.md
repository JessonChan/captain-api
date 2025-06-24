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

1. **Install Wails3** (if not already installed):
   ```bash
   go install github.com/wailsapp/wails/v3/cmd/wails3@latest
   ```

2. **Clone the repository** (if not already done):
   ```bash
   git clone <your-repo-url>
   cd captain-api
   ```

3. **Install dependencies and build**:
   ```bash
   # This will automatically install frontend dependencies, 
   # generate bindings, and build the application
   wails3 build
   ```

   Or for development with auto-reload:
   ```bash
   wails3 dev
   ```

### Development

To run the application in development mode with hot reload:

```bash
wails3 dev
```

This will start the application with:
- Hot reload for frontend changes
- Automatic Go code recompilation
- Debug console access

### Production Build

To create a production build:

```bash
wails3 build
```

The built application will be available in the `build/bin/` directory.

### Generate Bindings

If you modify Go services and need to regenerate frontend bindings:

```bash
wails3 generate bindings
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
- Inspired by [Postman](https://www.postman.com/)
- Icons and styling inspired by modern API tools

---

**Happy API Testing! 🚀**