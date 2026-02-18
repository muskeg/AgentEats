# AgentEats

AI-agent-first restaurant directory, built in Go for maximum performance and minimal infrastructure cost. Structured data served via REST API and MCP (Model Context Protocol) so that LLM-powered agents can search restaurants, browse menus, get recommendations, and make reservations.

## Quick Start

```bash
# Install dependencies
go mod tidy

# Seed demo data (8 restaurants, 130+ menu items)
go run ./cmd/seed

# Start the REST API (http://localhost:8000)
go run ./cmd/api

# In another terminal — start the MCP server (stdio)
go run ./cmd/mcp
```

### Build binaries

```bash
go build -o agenteats-api ./cmd/api
go build -o agenteats-mcp ./cmd/mcp
go build -o agenteats-seed ./cmd/seed
```

### MCP Client Configuration (Claude Desktop, etc.)

```json
{
  "mcpServers": {
    "agenteats": {
      "command": "/path/to/agenteats-mcp"
    }
  }
}
```

## Architecture

```
┌─────────────────────────────────────────────────┐
│                  AI Agents                       │
│  (Claude, GPT, custom agents, chat interfaces)  │
└──────────┬──────────────────┬────────────────────┘
           │                  │
     MCP Protocol        REST API
      (stdio)          (HTTP/JSON)
           │                  │
┌──────────▼──────────────────▼────────────────────┐
│              AgentEats Service (Go)               │
│  ┌────────────┐  ┌───────────┐                    │
│  │ MCP Server │  │ chi       │                    │
│  │ (mcp-go)   │  │ (router)  │                    │
│  └─────┬──────┘  └─────┬─────┘                    │
│        │               │                          │
│  ┌─────▼───────────────▼──────────────────────┐   │
│  │          Service Layer                      │  │
│  │  search · recommend · reserve · manage      │  │
│  └─────────────────┬───────────────────────────┘  │
│                    │                              │
│  ┌─────────────────▼───────────────────────────┐  │
│  │          GORM + SQLite                      │  │
│  └─────────────────────────────────────────────┘  │
└───────────────────────────────────────────────────┘
```

## Tech Stack

| Component | Choice | Why |
|-----------|--------|-----|
| Language | Go 1.23+ | 200-500K req/s, tiny memory, single binary |
| Router | chi | Lightweight, stdlib-compatible, great middleware |
| ORM | GORM + SQLite | Zero-config, portable, upgradeable to Postgres |
| MCP | mcp-go | Most popular Go MCP SDK |
| Config | envconfig | Env-var driven, zero boilerplate |

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/restaurants` | Search & filter restaurants |
| `GET` | `/restaurants/{id}` | Full restaurant details |
| `POST` | `/restaurants` | Register a restaurant (owner) |
| `PUT` | `/restaurants/{id}` | Update restaurant info |
| `GET` | `/restaurants/{id}/menu` | Get structured menu |
| `POST` | `/restaurants/{id}/menu/items` | Add a menu item |
| `GET` | `/restaurants/{id}/availability` | Check reservation slots |
| `POST` | `/restaurants/{id}/reservations` | Make a reservation |
| `GET` | `/restaurants/{id}/reservations` | List reservations |
| `DELETE` | `/reservations/{id}` | Cancel a reservation |
| `GET` | `/recommendations` | AI-friendly recommendations |
| `GET` | `/health` | Service health check |

## MCP Tools (for AI agents)

| Tool | Description |
|------|-------------|
| `search_restaurants` | Find restaurants by cuisine, price, location, dietary needs |
| `get_restaurant_details` | Full info including hours, contact, description |
| `get_menu` | Structured menu with prices, dietary labels, descriptions |
| `get_recommendations` | Personalized restaurant suggestions |
| `make_reservation` | Book a table (date, time, party size) |
| `check_availability` | Check reservation availability |
| `cancel_reservation` | Cancel an existing reservation |

## Data Model

Restaurants expose rich structured data optimized for AI consumption:

- **Restaurant**: name, cuisines, price range ($–$$$$), location, hours, contact, features
- **Menu Items**: name, description, price, category, dietary labels (vegan, gluten-free, etc.)
- **Reservations**: date, time, party size, status, special requests

## Configuration

Environment variables (or `.env` file):

| Variable | Default | Description |
|----------|---------|-------------|
| `DATABASE_URL` | `agenteats.db` | SQLite database file path |
| `HOST` | `0.0.0.0` | API server host |
| `PORT` | `8000` | API server port |
| `DEBUG` | `false` | Enable debug logging |

## Project Structure

```
├── cmd/
│   ├── api/main.go          # REST API server
│   ├── mcp/main.go          # MCP server (stdio)
│   └── seed/main.go         # Database seeder
├── internal/
│   ├── config/config.go     # Environment configuration
│   ├── database/db.go       # GORM database init
│   ├── dto/dto.go           # Request/response schemas
│   ├── handlers/handlers.go # HTTP route handlers
│   ├── mcpserver/server.go  # MCP tool definitions
│   ├── models/models.go     # Database models
│   └── services/services.go # Business logic
├── go.mod
├── go.sum
└── README.md
```

## License

MIT
