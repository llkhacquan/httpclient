.PHONY: check test fmt vet lint clean help

# Default target
help:
	@echo "Available targets:"
	@echo "  check    - Run all quality checks (fmt, vet, test, lint)"
	@echo "  test     - Run tests"
	@echo "  fmt      - Format code"
	@echo "  vet      - Run go vet"
	@echo "  lint     - Run golangci-lint (if available)"
	@echo "  clean    - Clean build artifacts"

# Run all quality checks
check: fmt vet test lint
	@echo "✅ All checks passed!"

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Run go vet
vet:
	@echo "Running go vet..."
	go vet ./...

# Run golangci-lint if available
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint not found, skipping lint check"; \
		echo "   Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Clean build artifacts
clean:
	@echo "Cleaning..."
	go clean ./...