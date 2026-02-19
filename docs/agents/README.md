# AgentEats — Agent & Consumer Guide

This guide is for developers building AI agents, chatbots, or applications that consume the AgentEats restaurant directory. All read endpoints are **public** — no authentication required.

**Base URL:** `https://agenteats.fly.dev`

---

## Table of Contents

- [Quick Start](#quick-start)
- [REST API Reference](#rest-api-reference)
  - [Health Check](#health-check)
  - [Search Restaurants](#search-restaurants)
  - [Get Restaurant Details](#get-restaurant-details)
  - [Get Menu](#get-menu)
  - [Get Recommendations](#get-recommendations)
  - [Check Availability](#check-availability)
  - [Make Reservation](#make-reservation)
  - [List Reservations](#list-reservations)
  - [Cancel Reservation](#cancel-reservation)
- [MCP Integration](#mcp-integration)
  - [Stdio Transport](#stdio-transport-local)
  - [Remote (Streamable HTTP)](#remote-streamable-http)
  - [MCP Tools Reference](#mcp-tools-reference)
  - [MCP Resource](#mcp-resource)
- [Data Types](#data-types)
- [Error Handling](#error-handling)
- [Rate Limits & Best Practices](#rate-limits--best-practices)

---

## Quick Start

### REST API (any HTTP client)

```bash
# Search for Italian restaurants in New York
curl "https://agenteats.fly.dev/restaurants?cuisine=Italian&city=New%20York"

# Get a specific restaurant's menu
curl "https://agenteats.fly.dev/restaurants/{id}/menu"

# Get AI-friendly recommendations
curl "https://agenteats.fly.dev/recommendations?cuisine=Japanese&city=New%20York&occasion=date_night"
```

### MCP (Claude Desktop, Cursor, etc.)

Add to your MCP client config:

```json
{
  "mcpServers": {
    "agenteats": {
      "command": "/path/to/agenteats-mcp"
    }
  }
}
```

Download the latest binary from [GitHub Releases](https://github.com/agenteats/agenteats/releases).

---

## REST API Reference

### Health Check

```
GET /health
```

**Response:**

```json
{
  "status": "ok",
  "version": "0.1.0",
  "service": "AgentEats"
}
```

---

### Search Restaurants

```
GET /restaurants
```

Searches and filters the restaurant directory. All parameters are optional.

**Query Parameters:**

| Parameter | Type | Example | Description |
|-----------|------|---------|-------------|
| `q` | string | `sushi` | Free-text search across name, description, and cuisines |
| `cuisine` | string | `Italian` | Filter by cuisine type |
| `city` | string | `New York` | Filter by city |
| `price_range` | string | `$$$` | Filter by price level: `$`, `$$`, `$$$`, `$$$$` |
| `features` | string | `outdoor_seating,wifi` | Comma-separated feature filter |
| `limit` | int | `10` | Max results (1–100, default 20) |
| `offset` | int | `0` | Pagination offset (default 0) |

**Available features:** `outdoor_seating`, `wifi`, `live_music`, `parking`, `delivery`, `takeout`, `wheelchair_accessible`, `pet_friendly`, `private_dining`, `bar`, `brunch`

**Response:** Array of `RestaurantSummary`

```json
[
  {
    "id": "abc-123-...",
    "name": "Bella Notte",
    "cuisines": ["Italian", "Mediterranean"],
    "price_range": "$$$",
    "city": "New York",
    "rating": 4.7,
    "review_count": 342,
    "address": "142 Thompson St",
    "features": ["outdoor_seating", "wifi", "live_music", "wheelchair_accessible"]
  }
]
```

---

### Get Restaurant Details

```
GET /restaurants/{id}
```

Returns full details for a single restaurant, including operating hours.

**Response:** `RestaurantDetail`

```json
{
  "id": "abc-123-...",
  "name": "Bella Notte",
  "description": "Authentic Italian trattoria with handmade pasta...",
  "cuisines": ["Italian", "Mediterranean"],
  "price_range": "$$$",
  "address": "142 Thompson St",
  "city": "New York",
  "state": "NY",
  "zip_code": "10012",
  "country": "US",
  "latitude": 40.727,
  "longitude": -73.999,
  "phone": "+1-212-555-0142",
  "email": "reservations@bellanotte.example.com",
  "website": "https://bellanotte.example.com",
  "features": ["outdoor_seating", "wifi", "live_music", "wheelchair_accessible"],
  "total_seats": 80,
  "rating": 4.7,
  "review_count": 342,
  "is_active": true,
  "hours": [
    { "day": "monday", "open_time": "17:00", "close_time": "23:00", "is_closed": false },
    { "day": "tuesday", "open_time": "17:00", "close_time": "23:00", "is_closed": false },
    { "day": "sunday", "open_time": "12:00", "close_time": "00:00", "is_closed": false }
  ]
}
```

---

### Get Menu

```
GET /restaurants/{id}/menu
```

Returns the full menu organized by category. Each item includes dietary labels and pricing.

**Response:** `MenuOut`

```json
{
  "restaurant_id": "abc-123-...",
  "restaurant_name": "Bella Notte",
  "currency": "USD",
  "categories": {
    "Appetizer": [
      {
        "id": "item-456-...",
        "category": "Appetizer",
        "name": "Bruschetta Trio",
        "description": "Tomato basil, mushroom truffle, and nduja spread on grilled sourdough",
        "price": 16.0,
        "currency": "USD",
        "dietary_labels": ["vegetarian"],
        "is_available": true,
        "is_popular": true,
        "calories": 420
      }
    ],
    "Main": [ ... ],
    "Dessert": [ ... ]
  }
}
```

**Dietary label values:** `vegetarian`, `vegan`, `gluten_free`, `dairy_free`, `nut_free`, `halal`, `kosher`, `spicy`, `raw`

---

### Get Recommendations

```
GET /recommendations
```

Returns personalized restaurant recommendations with match scoring. Ideal for agent-driven "where should I eat?" flows.

**Query Parameters:**

| Parameter | Type | Example | Description |
|-----------|------|---------|-------------|
| `cuisine` | string | `Japanese` | Preferred cuisine type |
| `city` | string | `New York` | City to search in |
| `price_range` | string | `$$` | Budget level |
| `features` | string | `delivery,wifi` | Desired features (comma-separated) |
| `dietary_needs` | string | `vegan,gluten_free` | Dietary requirements (comma-separated) |
| `occasion` | string | `date_night` | Type of occasion |
| `limit` | int | `5` | Number of results (1–20, default 5) |

**Occasion values:** `date_night`, `business`, `family`, `casual`, `celebration`

**Response:** Array of `RecommendationOut`

```json
[
  {
    "restaurant": {
      "id": "abc-123-...",
      "name": "Bella Notte",
      "cuisines": ["Italian", "Mediterranean"],
      "price_range": "$$$",
      "city": "New York",
      "rating": 4.7,
      "review_count": 342,
      "address": "142 Thompson St",
      "features": ["outdoor_seating", "wifi", "live_music"]
    },
    "match_reasons": ["Matches cuisine preference: Italian", "Has outdoor_seating"],
    "relevance_score": 0.85
  }
]
```

---

### Check Availability

```
GET /restaurants/{id}/availability?date=2026-03-15&party_size=4
```

Returns available 30-minute reservation slots for a given date.

**Query Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `date` | string | Yes | Date in `YYYY-MM-DD` format |
| `party_size` | int | No | Number of guests (default 2) |

**Response:**

```json
{
  "restaurant_id": "abc-123-...",
  "restaurant_name": "Bella Notte",
  "date": "2026-03-15",
  "available_times": ["11:00", "11:30", "12:00", "12:30", "18:00", "18:30", "19:00"],
  "max_party_size": 80
}
```

---

### Make Reservation

```
POST /restaurants/{id}/reservations
Content-Type: application/json
```

**Request body:**

```json
{
  "customer_name": "Alice Johnson",
  "customer_email": "alice@example.com",
  "customer_phone": "+1-555-0101",
  "party_size": 2,
  "date": "2026-03-15",
  "time": "19:00",
  "special_requests": "Window table if possible"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `customer_name` | string | Yes | Full name for the reservation |
| `party_size` | int | Yes | Number of guests |
| `date` | string | Yes | `YYYY-MM-DD` format |
| `time` | string | Yes | `HH:MM` 24-hour format |
| `customer_email` | string | No | Email for confirmation |
| `customer_phone` | string | No | Phone number |
| `special_requests` | string | No | Notes (allergies, birthday, high chair, etc.) |

**Response:** `201 Created`

```json
{
  "id": "res-789-...",
  "restaurant_id": "abc-123-...",
  "restaurant_name": "Bella Notte",
  "customer_name": "Alice Johnson",
  "customer_email": "alice@example.com",
  "party_size": 2,
  "date": "2026-03-15",
  "time": "19:00",
  "status": "confirmed",
  "special_requests": "Window table if possible",
  "created_at": "2026-02-18T22:30:00Z"
}
```

---

### List Reservations

```
GET /restaurants/{id}/reservations
GET /restaurants/{id}/reservations?date=2026-03-15
```

Returns reservations for a restaurant, optionally filtered by date.

---

### Cancel Reservation

```
DELETE /reservations/{id}
```

Cancels an existing reservation. Returns the updated reservation with `status: "cancelled"`.

---

## MCP Integration

AgentEats exposes a full [Model Context Protocol](https://modelcontextprotocol.io) server, enabling LLM agents to interact with the restaurant directory using structured tools.

### Stdio Transport (Local)

For local agents — Claude Desktop, Cursor, Cline, etc. Download the binary from [GitHub Releases](https://github.com/agenteats/agenteats/releases).

**Claude Desktop config** (`~/.config/claude/claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "agenteats": {
      "command": "/usr/local/bin/agenteats-mcp"
    }
  }
}
```

**With a local database:**

```json
{
  "mcpServers": {
    "agenteats": {
      "command": "/usr/local/bin/agenteats-mcp",
      "env": {
        "DATABASE_URL": "postgres://user:pass@host:5432/agenteats"
      }
    }
  }
}
```

### Remote (Streamable HTTP)

The hosted AgentEats API exposes a remote MCP endpoint — no local binary required:

```json
{
  "mcpServers": {
    "agenteats": {
      "url": "https://agenteats.fly.dev/mcp"
    }
  }
}
```

The endpoint uses the [Streamable HTTP](https://modelcontextprotocol.io/specification/2025-03-26/basic/transports#streamable-http) transport (stateless POST/response per tool call), which is scale-to-zero friendly and works with any MCP-compatible client.

### MCP Tools Reference

| Tool | Description | Key Parameters |
|------|-------------|----------------|
| `search_restaurants` | Find restaurants by cuisine, price, city, features | `query`, `city`, `cuisine`, `price_range`, `features`, `limit` |
| `get_restaurant_details` | Full info including hours, contact, description | `restaurant_id` (required) |
| `get_menu` | Structured menu with prices, dietary labels | `restaurant_id` (required) |
| `get_recommendations` | Personalized suggestions with match scoring | `cuisine`, `city`, `price_range`, `features`, `dietary_needs`, `occasion`, `limit` |
| `check_availability` | Check available reservation slots | `restaurant_id` (required), `date` (required), `party_size` |
| `make_reservation` | Book a table | `restaurant_id`, `customer_name`, `party_size`, `date`, `time` (all required) |
| `cancel_reservation` | Cancel an existing reservation | `reservation_id` (required) |

### MCP Resource

| URI | Description |
|-----|-------------|
| `agenteats://info` | Service metadata and capabilities summary (JSON) |

---

## Data Types

### Price Range

| Value | Meaning |
|-------|---------|
| `$` | Budget (~$10–15/person) |
| `$$` | Moderate (~$15–30/person) |
| `$$$` | Upscale (~$30–60/person) |
| `$$$$` | Fine dining ($60+/person) |

### Reservation Status

| Value | Meaning |
|-------|---------|
| `confirmed` | Active reservation |
| `cancelled` | Cancelled by customer |
| `completed` | Successfully completed |
| `no_show` | Customer did not arrive |

### Dietary Labels

`vegetarian`, `vegan`, `gluten_free`, `dairy_free`, `nut_free`, `halal`, `kosher`, `spicy`, `raw`

### Features

`outdoor_seating`, `wifi`, `live_music`, `parking`, `delivery`, `takeout`, `wheelchair_accessible`, `pet_friendly`, `private_dining`, `bar`, `brunch`

---

## Error Handling

All errors return a JSON body with an `error` field:

```json
{
  "error": "Restaurant not found"
}
```

**HTTP status codes:**

| Code | Meaning |
|------|---------|
| `200` | Success |
| `201` | Created (reservations, restaurants) |
| `400` | Bad request (missing/invalid parameters) |
| `404` | Resource not found |
| `500` | Internal server error |

---

## Rate Limits & Best Practices

- **No rate limits** are currently enforced. Be respectful — batch requests where possible.
- Use `limit` and `offset` for pagination instead of fetching all records.
- Cache restaurant details and menus when appropriate — they change infrequently.
- Always use `check_availability` before `make_reservation` to ensure the slot is open.
- The `/recommendations` endpoint does server-side scoring — prefer it over client-side filtering.
- All times are in **24-hour format** (`HH:MM`). All dates are **`YYYY-MM-DD`**.
