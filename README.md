# AgentEats

AI-agent-first restaurant directory, built in Go for maximum performance and minimal infrastructure cost. Structured data served via REST API and MCP (Model Context Protocol) so that LLM-powered agents can search restaurants, browse menus, get recommendations, and make reservations.

**Live demo:** [agenteats.fly.dev](https://agenteats.fly.dev/health)

## Documentation

| Audience | Guide |
|----------|-------|
| **AI agent developers** | [docs/agents/](docs/agents/) — REST API reference, MCP integration, data types, examples |
| **Restaurant owners** | [docs/owners/](docs/owners/) — Registration, authentication, menu management, bulk import |

## Quick Start

```bash
# Install dependencies
go mod tidy

# Seed demo data (8 restaurants, ~100 menu items)
go run ./cmd/seed

# Start the REST API (http://localhost:8000)
go run ./cmd/api

# In another terminal — start the MCP server (stdio)
go run ./cmd/mcp
```

### Build

```bash
make build         # Build all binaries
make test          # Run tests
make release       # Optimized release binaries
```

### Docker

```bash
docker build -t agenteats .
docker run --rm -p 8000:8000 agenteats

# With Postgres
docker run --rm -p 8000:8000 \
  -e DATABASE_URL="postgres://user:pass@host:5432/agenteats" \
  agenteats
```

### MCP Client Configuration

**Remote** (ChatGPT, Claude, VS Code, Cursor — no install needed):
```json
{
  "mcpServers": {
    "agenteats": {
      "url": "https://agenteats.fly.dev/mcp"
    }
  }
}
```

**Stdio** (local agents, offline use):
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

```mermaid
graph TD
    A["AI Agents<br/>(Claude, ChatGPT, VS Code, Cursor)"] -->|"MCP<br/>(Streamable HTTP)"| MCP["/mcp endpoint<br/>mcp-go"]
    A -->|"REST API<br/>(HTTP/JSON)"| API["HTTP Router<br/>chi"]

    subgraph AgentEats Service
        MCP --> SVC["Service Layer<br/>search · recommend · reserve · manage"]
        API --> SVC
        SVC --> DB["GORM<br/>SQLite / PostgreSQL"]
    end
```

## Tech Stack

| Component | Choice | Why |
|-----------|--------|-----|
| Language | Go 1.23+ | 200-500K req/s, tiny memory, single binary |
| Router | [chi](https://github.com/go-chi/chi) | Lightweight, stdlib-compatible, great middleware |
| ORM | [GORM](https://gorm.io) | SQLite (dev) or PostgreSQL (prod), auto-detected from `DATABASE_URL` |
| MCP | [mcp-go](https://github.com/mark3labs/mcp-go) | Go MCP SDK — stdio + Streamable HTTP transports |
| Config | [envconfig](https://github.com/kelseyhightower/envconfig) | Env-var driven, zero boilerplate |
| Hosting | [Fly.io](https://fly.io) | Scale-to-zero, global edge, ~$0/mo idle |
| Database | [Neon](https://neon.tech) | Serverless Postgres, free tier, auto-suspend |

## API Endpoints

### Public (no auth required)

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/health` | Service health check |
| `GET` | `/restaurants` | Search & filter restaurants |
| `GET` | `/restaurants/{id}` | Full restaurant details |
| `GET` | `/restaurants/{id}/menu` | Get structured menu |
| `GET` | `/restaurants/{id}/availability` | Check reservation slots |
| `POST` | `/restaurants/{id}/reservations` | Make a reservation |
| `GET` | `/restaurants/{id}/reservations` | List reservations |
| `DELETE` | `/reservations/{id}` | Cancel a reservation |
| `GET` | `/recommendations` | AI-friendly recommendations |
| `POST` | `/owners/register` | Register a restaurant owner account |

### Authenticated (`Authorization: Bearer <api-key>`)

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/owners/rotate-key` | Rotate API key (invalidates old key) |
| `POST` | `/restaurants` | Create a restaurant (assigned to owner) |
| `PUT` | `/restaurants/{id}` | Update restaurant (ownership enforced) |
| `POST` | `/restaurants/{id}/menu/items` | Add a menu item (ownership enforced) |
| `POST` | `/restaurants/{id}/menu/import` | Bulk import menu (`replace` or `merge`) |

**Query parameters** for `GET /restaurants`:

| Param | Example | Description |
|-------|---------|-------------|
| `query` | `sushi` | Free-text search (name, description, cuisines) |
| `cuisine` | `Italian` | Filter by cuisine type |
| `city` | `New York` | Filter by city |
| `price_range` | `$$$` | Filter by price level (`$` to `$$$$`) |
| `features` | `outdoor_seating,wifi` | Comma-separated feature filters |

**Query parameters** for `GET /recommendations`:

| Param | Example | Description |
|-------|---------|-------------|
| `cuisine` | `Japanese` | Preferred cuisine |
| `city` | `New York` | City to search in |
| `price_range` | `$$` | Budget level |
| `features` | `delivery` | Desired features |
| `dietary_needs` | `vegan,gluten_free` | Dietary requirements |
| `occasion` | `date_night` | Type of occasion |

## MCP Tools

| Tool | Description |
|------|-------------|
| `search_restaurants` | Find restaurants by cuisine, price, location, dietary needs |
| `get_restaurant_details` | Full info including hours, contact, description |
| `get_menu` | Structured menu with prices, dietary labels, descriptions |
| `get_recommendations` | Personalized restaurant suggestions with match scoring |
| `check_availability` | Check available reservation time slots |
| `make_reservation` | Book a table (date, time, party size) |
| `cancel_reservation` | Cancel an existing reservation |

**Resource:** `agenteats://info` — service metadata and capabilities summary.

## Data Model

```mermaid
erDiagram
    Owner ||--o{ Restaurant : owns
    Restaurant ||--o{ OperatingHours : has
    Restaurant ||--o{ MenuItem : offers
    Restaurant ||--o{ Reservation : accepts

    Owner {
        string id PK
        string name
        string email UK
        string api_key_hash UK
        bool is_active
    }

    Restaurant {
        string id PK
        string owner_id FK
        string name
        string cuisines
        string price_range
        string city
        float rating
        int review_count
        string features
    }

    MenuItem {
        string id PK
        string restaurant_id FK
        string name
        string category
        float price
        string dietary_labels
        bool is_popular
        int calories
    }

    Reservation {
        string id PK
        string restaurant_id FK
        string customer_name
        string date
        string time
        int party_size
        string status
    }
```

## Configuration

Environment variables (or `.env` file):

| Variable | Default | Description |
|----------|---------|-------------|
| `DATABASE_URL` | `agenteats.db` | SQLite file path, or `postgres://...` for Postgres |
| `HOST` | `0.0.0.0` | Server bind address |
| `PORT` | `8000` | REST API server port |
| `MCP_TRANSPORT` | `stdio` | MCP transport for standalone binary: `stdio` or `http` |
| `MCP_PORT` | `8001` | MCP HTTP server port (when `MCP_TRANSPORT=http`) |
| `DEBUG` | `false` | Enable verbose query logging |

## Deployment

The app runs on [Fly.io](https://fly.io) with [Neon](https://neon.tech) serverless Postgres:

```bash
# First-time setup
fly launch --no-deploy
fly secrets set DATABASE_URL="postgres://..."
fly deploy

# Seed production database
fly ssh console -C "agenteats-seed"
```

CI/CD is handled by GitHub Actions:

| Workflow | Trigger | What it does |
|----------|---------|-------------|
| `ci.yml` | Push / PR to `main` | Build, vet, test with race detector |
| `deploy.yml` | Push to `main` | Build + test gate, then deploy to Fly.io |
| `release.yml` | Push to `main` | release-please PR, cross-platform binaries on release |

## Project Structure

```
├── cmd/
│   ├── api/main.go              # REST API server
│   ├── mcp/main.go              # Standalone MCP server (stdio / HTTP)
│   └── seed/main.go             # Database seeder
├── docs/
│   ├── agents/README.md         # Agent & consumer API guide
│   └── owners/README.md         # Restaurant owner guide
├── internal/
│   ├── config/config.go         # Environment configuration
│   ├── database/db.go           # GORM init (SQLite / Postgres auto-detect)
│   ├── dto/dto.go               # Request/response DTOs
│   ├── handlers/handlers.go     # HTTP route handlers
│   ├── mcpserver/server.go      # MCP tool & resource definitions
│   ├── middleware/auth.go       # API key auth middleware
│   ├── models/models.go         # Database models (Owner, Restaurant, MenuItem, etc.)
│   └── services/services.go     # Business logic
├── .github/workflows/
│   ├── ci.yml                   # Build & test
│   ├── deploy.yml               # CD to Fly.io
│   └── release.yml              # release-please + binary releases
├── Dockerfile                   # Multi-stage production build
├── fly.toml                     # Fly.io app configuration
├── release-please-config.json   # Release automation config
├── Makefile                     # Build, test, run shortcuts
└── go.mod
```

## License

All rights reserved. This software is proprietary and not licensed for redistribution.
