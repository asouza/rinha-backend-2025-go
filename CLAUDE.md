# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go-based backend project for "Rinha de Backend" (Backend Challenge), implementing a simple HTTP API server using Chi router. The project follows a clean architecture pattern with separate layers for API handlers and main application entry point.

## Architecture

- **cmd/main.go**: Application entry point that sets up the Chi router and starts the HTTP server on port 8080
- **internal/api/handler.go**: Contains HTTP handlers for API endpoints
- **go.mod**: Go module definition with Chi router dependency

The project uses:
- Chi router (github.com/go-chi/chi/v5) for HTTP routing
- Standard Go HTTP server
- Clean architecture with internal package structure

## Common Commands

### Build and Run
```bash
# Build the application
go build -o bin/server cmd/main.go

# Run directly
go run cmd/main.go

# Run with Go modules
go mod tidy  # Update dependencies
go run .
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests for specific package
go test ./internal/api
```

### Code Quality
```bash
# Format code
go fmt ./...

# Check for potential issues
go vet ./...

# Run golint (if installed)
golint ./...
```

### Development
```bash
# Download dependencies
go mod download

# Clean module cache
go clean -modcache

# Check module dependencies
go mod verify
```

## Project Structure

The project follows Go's standard project layout:
- `cmd/`: Application entry points
- `internal/`: Private application code that shouldn't be imported by other applications
- `go.mod`/`go.sum`: Go module files for dependency management

## Key Implementation Details

- Server runs on port 8080
- Uses Chi router for HTTP handling
- Root endpoint ("/") returns "Olá, Rinha de Backend!"
- Error handling includes server startup failures
- Logging is implemented for server status