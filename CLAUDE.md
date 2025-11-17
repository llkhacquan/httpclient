# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

A Go HTTP client library that wraps `net/http` with automatic JSON marshaling/unmarshaling and flexible request configuration. Designed for simplicity and customization.

## Development Commands

```bash
make check           # Run all quality checks (fmt, vet, test, lint)
make test            # Run tests only
go test -v           # Run tests with verbose output
make fmt             # Format code
make vet             # Run go vet
make lint            # Run golangci-lint (if available)
```

## Architecture

The package provides two usage patterns:

1. **Default client** (`default_client.go`): Package-level functions (Get, Post, Patch, Delete) that use a singleton client for simple use cases
2. **Custom client** (`client.go`): Configurable `Client` struct for advanced scenarios with custom HTTP clients, marshal/unmarshal functions

### Core Components

- `Client` struct: Wraps `http.Client` with customizable `MarshalFunc` and `UnmarshalFunc`
- Options pattern: Functional options for headers, status code handling, and per-request HTTP client override
- Context support: All methods accept `context.Context` for cancellation and timeouts

### Key Design Patterns

- **Byte array pass-through**: Request bodies of type `[]byte` bypass JSON marshaling and are sent directly
- **Status code handling**: `WithStatus(&statusVar)` option allows non-2xx responses without errors and captures the HTTP status code
- **Error wrapping**: All errors include context about the failed operation (e.g., "failed to make POST request")
- **Graceful unmarshal failures**: When using `WithStatus()`, unmarshal errors on 4xx/5xx responses are ignored to handle non-JSON error responses

## Testing

Tests use real external APIs:
- Pokemon API (pokeapi.co) for GET requests
- httpbin.org for POST/PATCH/DELETE and status code testing

All HTTP methods support context cancellation, custom headers, and status code capture via the options pattern.
