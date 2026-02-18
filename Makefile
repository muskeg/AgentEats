.PHONY: build seed api mcp mcp-sse clean test docker docker-run deps release

# Build all binaries
build:
	go build -o agenteats-api ./cmd/api
	go build -o agenteats-mcp ./cmd/mcp
	go build -o agenteats-seed ./cmd/seed

# Seed the database with demo data
seed:
	go run ./cmd/seed

# Start the REST API server
api:
	go run ./cmd/api

# Start the MCP server (stdio)
mcp:
	go run ./cmd/mcp

# Start the MCP server (SSE on port 8001)
mcp-sse:
	MCP_TRANSPORT=sse MCP_PORT=8001 go run ./cmd/mcp

# Run tests
test:
	go test ./... -v

# Clean build artifacts and database
clean:
	rm -f agenteats-api agenteats-mcp agenteats-seed agenteats.db

# Download dependencies
deps:
	go mod tidy

# Build optimized release binaries
release:
	CGO_ENABLED=1 go build -ldflags="-s -w" -o agenteats-api ./cmd/api
	CGO_ENABLED=1 go build -ldflags="-s -w" -o agenteats-mcp ./cmd/mcp
	CGO_ENABLED=1 go build -ldflags="-s -w" -o agenteats-seed ./cmd/seed

# Build Docker image
docker:
	docker build -t agenteats .

# Run with Docker
docker-run:
	docker run --rm -p 8000:8000 agenteats
