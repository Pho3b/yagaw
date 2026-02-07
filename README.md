# Yagaw - Yet Another Golang Advanced Webserver

A lightweight, simple, and powerful HTTP web server framework written in Go. Yagaw provides a clean and intuitive API for building HTTP servers with custom route handling, built-in logging, and support for multiple HTTP methods.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Core Components](#core-components)
  - [Server](#server)
  - [Router](#router)
  - [Request Handlers](#request-handlers)
  - [Dynamic Route Parameters](#dynamic-route-parameters-path-parameter-matching)
- [API Reference](#api-reference)
- [Examples](#examples)
- [Logging](#logging)
- [Contributing](#contributing)
- [License](#license)

## Features

- **Lightweight**: Minimal dependencies, focusing on core HTTP functionality
- **Simple API**: Clean and intuitive interface for defining routes and handlers
- **HTTP Methods Support**: Built-in support for GET, POST, and extensible to other methods
- **Routing**: Fast and efficient route matching with method-based organization
- **Logging**: Integrated logging system with configurable log levels
- **Debug Mode**: Built-in request debugging capabilities
- **Net/HTTP Compatible**: Built on Go's standard `net/http` package for reliability

## Installation

To use Yagaw in your Go project:

```bash
go get github.com/Algatux/yagaw
```

## Quick Start

Here's a minimal example to get your server running:

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/Algatux/yagaw"
	"github.com/pho3b/tiny-logger/logs/log_level"
)

func main() {
	// Initialize logger
	yagaw.Log = yagaw.InitLogger(log_level.DebugLvlName)

	// Create a new server on localhost:8080
	server := yagaw.NewServer("localhost", 8080)

	// Get the router
	router := server.GetRouter()

	// Register a GET route
	router.RegisterRoute(yagaw.GET, "/test", func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "text/plain")
		fmt.Fprintln(rw, "Welcome to Yagaw!")
	})

	// Start the server
	server.Run()
}
```

Then run:

```bash
go run main.go
```

Visit `http://localhost:8080/test` in your browser to see the response.

## Core Components

### Server

The `Server` struct is the main entry point for your web application. It encapsulates the HTTP server configuration, routing, and request handling.

#### Creating a Server

```go
server := yagaw.NewServer(address string, port int) *Server
```

**Parameters:**
- `address` (string): The IP address to bind to. Use an empty string `""` to bind to all interfaces (0.0.0.0)
- `port` (int): The port number to listen on

**Example:**
```go
server := yagaw.NewServer("", 8080)      // Bind to all interfaces on port 8080
server := yagaw.NewServer("localhost", 3000) // Bind to localhost on port 3000
```

#### Server Methods

##### Run()

Starts the HTTP server and begins listening for incoming requests. This is a blocking call.

```go
server.Run()
```

**Example:**
```go
server := yagaw.NewServer("", 8080)
server.Run() // Server will listen until an error occurs
```

##### GetRouter()

Retrieves the router instance associated with this server, allowing you to register routes.

```go
router := server.GetRouter() *Router
```

**Example:**
```go
server := yagaw.NewServer("", 8080)
router := server.GetRouter()
router.RegisterRoute(yagaw.GET, "/api/users", handleGetUsers)
```

### Router

The `Router` struct manages HTTP route registration and request dispatching. It implements the `http.Handler` interface, making it compatible with Go's standard HTTP server.

#### Route Registration

Routes are registered by HTTP method and path. When a request arrives, the router matches the method and path to find the appropriate handler.

##### RegisterRoute()

Registers a new route with the specified HTTP method and path.

```go
func (r *Router) RegisterRoute(method Method, path string, handler ReqHandler)
```

**Parameters:**
- `method` (Method): The HTTP method (GET, POST, etc.)
- `path` (string): The URL path (e.g., "/api/users", "/", "/products/123"). Supports dynamic path parameters using curly braces syntax.
- `handler` (ReqHandler): The handler function to execute for this route

**Example:**
```go
router.RegisterRoute(yagaw.GET, "/", func(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(rw, "Home page")
})

router.RegisterRoute(yagaw.POST, "/api/data", func(rw http.ResponseWriter, req *http.Request) {
	// Handle POST request
})
```

#### Dynamic Route Parameters (Path Parameter Matching)

Yagaw supports dynamic route parameters that allow you to capture values from the URL path. This is useful for creating RESTful endpoints where path segments vary.

##### Syntax

Dynamic parameters are defined using curly braces with the parameter name inside:

```
/users/{id}
/posts/{postId}/comments/{commentId}
/files/{filename}
```

##### Parameter Rules

- Parameter names can contain lowercase letters (a-z), numbers (0-9), hyphens (-), and underscores (_)
- Parameters are case-insensitive during matching
- Parameters use regex pattern matching: `[a-z0-9-_]+`
- Multiple parameters can be used in a single path
- Parameters are matched using regex patterns compiled during route registration

##### How It Works

When you register a route with parameters like `/users/{id}`, the router:

1. **Detection**: Identifies parameter placeholders using the pattern `{[a-z0-9-_]+}`
2. **Conversion**: Converts the path to a regex pattern: `/users/([a-z0-9-_]+)`
3. **Registration**: Stores the regex pattern as the route key
4. **Matching**: When a request arrives, the router matches the request path against the registered regex patterns
5. **Fallback**: If an exact path match fails, it tries regex pattern matching to find a parameterized route

##### Examples

**Basic Parameter Matching:**
```go
// Register a route with a parameter
router.RegisterRoute(yagaw.GET, "/users/{id}", func(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(rw, "User path: %s\n", req.URL.Path)
})

// These requests will match:
// GET /users/123
// GET /users/john-doe
// GET /users/user_42
```

**Multiple Parameters:**
```go
router.RegisterRoute(yagaw.GET, "/posts/{postId}/comments/{commentId}", func(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(rw, `{"path":"%s"}\n`, req.URL.Path)
})

// These requests will match:
// GET /posts/123/comments/456
// GET /posts/my-post/comments/comment-1
```

**Dynamic Resource Endpoints:**
```go
// Get a specific resource by ID
router.RegisterRoute(yagaw.GET, "/api/v1/products/{productId}", handleGetProduct)

// Update a resource
router.RegisterRoute(yagaw.POST, "/api/v1/products/{productId}", handleUpdateProduct)

// Nested resources
router.RegisterRoute(yagaw.GET, "/api/v1/users/{userId}/orders/{orderId}", handleGetUserOrder)
```

**Extracting Parameters from Request:**

To extract the parameter values from the request path, you need to parse the URL path:

```go
import (
	"regexp"
	"strings"
)

router.RegisterRoute(yagaw.GET, "/users/{userId}", func(rw http.ResponseWriter, req *http.Request) {
	// Method 1: Simple string manipulation
	parts := strings.Split(req.URL.Path, "/")
	if len(parts) >= 3 {
		userId := parts[2]
		rw.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(rw, "User ID: %s\n", userId)
	}
})

router.RegisterRoute(yagaw.GET, "/posts/{postId}/comments/{commentId}", func(rw http.ResponseWriter, req *http.Request) {
	// Method 2: Regex extraction for multiple parameters
	re := regexp.MustCompile(`/posts/([a-z0-9-_]+)/comments/([a-z0-9-_]+)`)
	matches := re.FindStringSubmatch(req.URL.Path)
	if len(matches) == 3 {
		postId := matches[1]
		commentId := matches[2]
		rw.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(rw, `{"postId":"%s","commentId":"%s"}\n`, postId, commentId)
	}
})
```

##### Case Insensitivity

Parameter matching is case-insensitive by default. This means `/users/{id}` will match:
- `/users/ABC`
- `/users/abc`
- `/users/AbC`

The matching is case-insensitive but the captured value preserves its original case from the request.

##### Priority and Matching Order

1. **Exact Path Match**: The router first attempts to find an exact match for the requested path
2. **Pattern Match**: If no exact match is found, it then tries regex pattern matching against parameterized routes
3. **Not Found**: If neither matches, the 404 handler is invoked

This ensures that exact routes take precedence over parameterized routes.

##### Performance Considerations

- Parameter matching uses regex compilation and matching, which has some overhead
- Register exact paths without parameters when possible for better performance
- The router caches regex patterns during route registration, not during request matching
- For frequently accessed routes, consider using exact paths when possible

##### RegisteredRoutes()

Returns a pointer to the internal route map, allowing inspection of all registered routes.

```go
routes := router.RegisteredRoutes() *ReqHandlerMap
```

**Example:**
```go
routes := router.RegisteredRoutes()
// You can now inspect or modify the routes map if needed
```

#### ServeHTTP()

This method implements the `http.Handler` interface. It processes incoming HTTP requests by finding the appropriate handler and executing it.

```go
func (r *Router) ServeHTTP(rw http.ResponseWriter, req *http.Request)
```

This method is called automatically by Go's HTTP server and typically does not need to be called directly.

### Request Handlers

A `ReqHandler` is a function that processes HTTP requests. It receives a `http.ResponseWriter` to write the response and a `*http.Request` containing the request details.

#### ReqHandler Type

```go
type ReqHandler func(rw http.ResponseWriter, req *http.Request)
```

#### Handler Signature

All handlers must follow this signature:

```go
func(rw http.ResponseWriter, req *http.Request)
```

**Parameters:**
- `rw` (http.ResponseWriter): Used to write the response headers and body
- `req` (*http.Request): Contains request information (method, path, headers, body, etc.)

#### Writing Responses

Within a handler, you can:

1. **Set Response Headers:**
   ```go
   rw.Header().Set("Content-Type", "application/json")
   rw.Header().Set("X-Custom-Header", "value")
   ```

2. **Set Status Code:**
   ```go
   rw.WriteHeader(http.StatusOK)     // 200
   rw.WriteHeader(http.StatusCreated) // 201
   rw.WriteHeader(http.StatusNotFound) // 404
   ```

3. **Write Response Body:**
   ```go
   fmt.Fprintln(rw, "Hello, World!")
   rw.Write([]byte("Response body"))
   ```

#### Example Handlers

**Simple Text Response:**
```go
router.RegisterRoute(yagaw.GET, "/hello", func(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "text/plain")
	fmt.Fprintln(rw, "Hello, World!")
})
```

**JSON Response:**
```go
import "encoding/json"

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

router.RegisterRoute(yagaw.GET, "/api/user", func(rw http.ResponseWriter, req *http.Request) {
	user := User{ID: 1, Name: "John Doe"}
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(user)
})
```

**Handling Query Parameters:**
```go
router.RegisterRoute(yagaw.GET, "/search", func(rw http.ResponseWriter, req *http.Request) {
	query := req.URL.Query().Get("q")
	rw.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(rw, "You searched for: %s\n", query)
})
```

**Handling Request Body (POST):**
```go
router.RegisterRoute(yagaw.POST, "/api/data", func(rw http.ResponseWriter, req *http.Request) {
	var data map[string]interface{}
	json.NewDecoder(req.Body).Decode(&data)
	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(rw, "Received: %v\n", data)
})
```

## API Reference

### Types

#### Method

```go
type Method string

const (
	GET  Method = "GET"
	POST Method = "POST"
)
```

Represents HTTP methods. You can extend this by defining new constants for additional methods (PUT, DELETE, PATCH, etc.).

#### ReqHandler

```go
type ReqHandler func(rw http.ResponseWriter, req *http.Request)
```

Function type for HTTP request handlers.

#### ReqHandlerMap

```go
type ReqHandlerMap map[Method]map[string]ReqHandler
```

Internal map structure for storing routes organized by method and path.

#### Server

```go
type Server struct {
	address string
	port    int
	server  *http.Server
	router  *Router
}
```

Represents the HTTP server configuration and routing.

#### Router

```go
type Router struct {
	routes ReqHandlerMap
}
```

Manages route registration and request dispatching.

### Functions

#### NewServer

```go
func NewServer(addr string, port int) *Server
```

Creates a new Server instance with the specified address and port.

#### NewRouter

```go
func NewRouter() *Router
```

Creates a new Router instance with initialized route maps.

#### InitLogger

```go
func InitLogger(logLevel log_level.LogLvlName) *logs.Logger
```

Initializes and configures the logger with the specified log level.

### Global Variables

#### Log

```go
var Log *logs.Logger
```

Global logger instance used throughout the package. Initialize it at application startup.

## Examples

### Basic GET Endpoint

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/Algatux/yagaw"
	"github.com/pho3b/tiny-logger/logs/log_level"
)

func main() {
	yagaw.Log = yagaw.InitLogger(log_level.DebugLvlName)
	server := yagaw.NewServer("", 8080)
	router := server.GetRouter()

	router.RegisterRoute(yagaw.GET, "/", func(rw http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(rw, "Home")
	})

	server.Run()
}
```

### Multiple Routes

```go
router.RegisterRoute(yagaw.GET, "/users", handleGetUsers)
router.RegisterRoute(yagaw.GET, "/users/123", handleGetUser)
router.RegisterRoute(yagaw.POST, "/users", handleCreateUser)
router.RegisterRoute(yagaw.GET, "/products", handleGetProducts)
```

### Complex Handler with Error Handling

```go
router.RegisterRoute(yagaw.POST, "/api/process", func(rw http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(rw, "Method not allowed")
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	fmt.Fprintln(rw, `{"status":"success"}`)
})
```

### With Dynamic Path Parameters

```go
package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/Algatux/yagaw"
	"github.com/pho3b/tiny-logger/logs/log_level"
)

func main() {
	yagaw.Log = yagaw.InitLogger(log_level.DebugLvlName)
	server := yagaw.NewServer("", 8080)
	router := server.GetRouter()

	// Simple parameter: /users/{id}
	router.RegisterRoute(yagaw.GET, "/users/{id}", func(rw http.ResponseWriter, req *http.Request) {
		parts := strings.Split(req.URL.Path, "/")
		if len(parts) >= 3 {
			userId := parts[2]
			rw.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(rw, `{"userId":"%s"}\n`, userId)
		}
	})

	// Multiple parameters: /posts/{postId}/comments/{commentId}
	router.RegisterRoute(yagaw.GET, "/posts/{postId}/comments/{commentId}", func(rw http.ResponseWriter, req *http.Request) {
		re := regexp.MustCompile(`/posts/([a-z0-9-_]+)/comments/([a-z0-9-_]+)`)
		matches := re.FindStringSubmatch(req.URL.Path)
		if len(matches) == 3 {
			postId := matches[1]
			commentId := matches[2]
			rw.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(rw, `{"postId":"%s","commentId":"%s"}\n`, postId, commentId)
		}
	})

	// File resource with extension
	router.RegisterRoute(yagaw.GET, "/files/{filename}", func(rw http.ResponseWriter, req *http.Request) {
		parts := strings.Split(req.URL.Path, "/")
		filename := parts[len(parts)-1]
		rw.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(rw, "Requested file: %s\n", filename)
	})

	server.Run()
}
```

**Test the endpoints:**
```bash
# GET /users/123 -> {"userId":"123"}
curl http://localhost:8080/users/123

# GET /users/john-doe -> {"userId":"john-doe"}
curl http://localhost:8080/users/john-doe

# GET /posts/42/comments/1 -> {"postId":"42","commentId":"1"}
curl http://localhost:8080/posts/42/comments/1

# GET /files/document-pdf -> Requested file: document-pdf
curl http://localhost:8080/files/document-pdf
```

### With Query Parameters and Path Extraction

```go
router.RegisterRoute(yagaw.GET, "/api/search", func(rw http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	keyword := query.Get("keyword")
	limit := query.Get("limit")

	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(rw, `{"keyword":"%s","limit":"%s"}\n`, keyword, limit)
})
```

## Logging

Yagaw integrates with the `tiny-logger` package for structured logging.

### Initializing the Logger

```go
yagaw.Log = yagaw.InitLogger(log_level.DebugLvlName)
```

### Available Log Levels

- `log_level.DebugLvlName`: Most verbose, includes debug messages
- `log_level.InfoLvlName`: General information messages
- `log_level.WarningLvlName`: Warning messages
- `log_level.ErrorLvlName`: Error messages only

### Using the Logger in Handlers

```go
router.RegisterRoute(yagaw.GET, "/test", func(rw http.ResponseWriter, req *http.Request) {
	yagaw.Log.Debug("Processing request")
	// Your handler code
})
```

## Project Structure

```
yagaw/
├── router.go           # Router implementation
├── server.go          # Server implementation and logging
├── go.mod             # Module definition
├── go.sum             # Dependency checksums
├── LICENSE            # License file
├── README.md          # This file
└── examples/
    └── api.go         # Example usage
```

## Contributing

Contributions are welcome! Please feel free to submit pull requests or open issues for bugs and feature requests.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

**Version:** 0.0.1  
**Author:** Algatux  
**Repository:** https://github.com/Algatux/yagaw
