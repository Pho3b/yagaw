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
- `path` (string): The URL path (e.g., "/api/users", "/", "/products/123")
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
