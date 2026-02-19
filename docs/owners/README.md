# AgentEats — Restaurant Owner Guide

This guide is for restaurant owners and operators who want to list their restaurant, manage their menu, and keep their information up to date on AgentEats.

**Base URL:** `https://agenteats.fly.dev`

---

## Table of Contents

- [Overview](#overview)
- [Getting Started](#getting-started)
  - [1. Register an Account](#1-register-an-account)
  - [2. Create Your Restaurant](#2-create-your-restaurant)
  - [3. Import Your Menu](#3-import-your-menu)
- [Authentication](#authentication)
  - [API Key Format](#api-key-format)
  - [Using Your API Key](#using-your-api-key)
  - [Rotating Your API Key](#rotating-your-api-key)
- [API Reference](#api-reference)
  - [Register Owner](#register-owner)
  - [Rotate API Key](#rotate-api-key)
  - [Create Restaurant](#create-restaurant)
  - [Update Restaurant](#update-restaurant)
  - [Add Menu Item](#add-menu-item)
  - [Bulk Import Menu](#bulk-import-menu)
- [Data Formats](#data-formats)
  - [Restaurant Fields](#restaurant-fields)
  - [Menu Item Fields](#menu-item-fields)
  - [Operating Hours Format](#operating-hours-format)
- [Examples](#examples)
  - [Full Restaurant Setup](#full-restaurant-setup)
  - [Seasonal Menu Update](#seasonal-menu-update)
  - [Bulk Import from CSV](#bulk-import-from-csv)
- [FAQ](#faq)

---

## Overview

AgentEats is an AI-agent-first restaurant directory. When you list your restaurant here, AI assistants like Claude, ChatGPT, and custom agents can find your restaurant, browse your menu, recommend you to hungry users, and handle reservations — all through structured data.

**What you get:**
- Your restaurant discoverable by AI agents worldwide
- Structured menu with dietary labels, pricing, and categories
- Automated reservation handling
- AI-powered recommendation matching (your restaurant surfaces when users ask for your cuisine, price range, or features)

**What you manage:**
- Restaurant details (name, address, hours, features)
- Menu items (categories, prices, descriptions, dietary info)
- Everything else (availability, recommendations, search ranking) is automatic

---

## Getting Started

### 1. Register an Account

```bash
curl -X POST https://agenteats.fly.dev/owners/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Maria Rossi",
    "email": "maria@bellanotte.com"
  }'
```

**Response:**

```json
{
  "id": "owner-uuid-...",
  "name": "Maria Rossi",
  "email": "maria@bellanotte.com",
  "api_key": "ae_7e924bfa8fe1e4190d905cebe864dac7..."
}
```

> **IMPORTANT:** The `api_key` is only shown once. Copy it immediately and store it securely. You'll need it for all management operations. If you lose it, use the [key rotation](#rotating-your-api-key) endpoint (you'll need the current key to do so).

### 2. Create Your Restaurant

```bash
curl -X POST https://agenteats.fly.dev/restaurants \
  -H "Authorization: Bearer ae_7e924bfa8fe1e4190d..." \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Bella Notte",
    "description": "Authentic Italian trattoria with handmade pasta and wood-fired pizza.",
    "cuisines": ["Italian", "Mediterranean"],
    "price_range": "$$$",
    "address": "142 Thompson St",
    "city": "New York",
    "state": "NY",
    "zip_code": "10012",
    "country": "US",
    "phone": "+1-212-555-0142",
    "email": "reservations@bellanotte.com",
    "website": "https://bellanotte.com",
    "features": ["outdoor_seating", "wifi", "live_music", "wheelchair_accessible"],
    "total_seats": 80,
    "hours": [
      {"day": "monday",    "open_time": "17:00", "close_time": "23:00", "is_closed": false},
      {"day": "tuesday",   "open_time": "17:00", "close_time": "23:00", "is_closed": false},
      {"day": "wednesday", "open_time": "17:00", "close_time": "23:00", "is_closed": false},
      {"day": "thursday",  "open_time": "17:00", "close_time": "23:00", "is_closed": false},
      {"day": "friday",    "open_time": "17:00", "close_time": "00:00", "is_closed": false},
      {"day": "saturday",  "open_time": "12:00", "close_time": "00:00", "is_closed": false},
      {"day": "sunday",    "open_time": "12:00", "close_time": "22:00", "is_closed": false}
    ]
  }'
```

Save the returned `id` — you'll need it for menu and update operations.

### 3. Import Your Menu

```bash
curl -X POST https://agenteats.fly.dev/restaurants/{id}/menu/import \
  -H "Authorization: Bearer ae_7e924bfa8fe1e4190d..." \
  -H "Content-Type: application/json" \
  -d '{
    "strategy": "replace",
    "items": [
      {
        "category": "Appetizer",
        "name": "Bruschetta Trio",
        "description": "Tomato basil, mushroom truffle, and nduja spread on grilled sourdough",
        "price": 16.00,
        "dietary_labels": ["vegetarian"],
        "is_available": true,
        "is_popular": true,
        "calories": 420
      },
      {
        "category": "Main",
        "name": "Osso Buco",
        "description": "Braised veal shank with saffron risotto and gremolata",
        "price": 42.00,
        "is_available": true,
        "is_popular": true,
        "calories": 850
      },
      {
        "category": "Dessert",
        "name": "Tiramisu",
        "description": "Classic espresso-soaked ladyfingers with mascarpone cream",
        "price": 14.00,
        "dietary_labels": ["vegetarian"],
        "is_available": true
      }
    ]
  }'
```

That's it — your restaurant is live and discoverable by AI agents.

---

## Authentication

All write operations require API key authentication. Read operations (searching, browsing menus) are always public.

### API Key Format

API keys use the format `ae_<64 hex characters>`:

```
ae_7e924bfa8fe1e4190d905cebe864dac78940896d76676a5fd037ed1b5b248344
```

The `ae_` prefix makes AgentEats keys easy to identify in logs or key managers.

### Using Your API Key

Include it in the `Authorization` header with the `Bearer` scheme:

```
Authorization: Bearer ae_7e924bfa8fe1e4190d905cebe864dac78940896d76676a5fd037ed1b5b248344
```

### Rotating Your API Key

If your key is compromised, rotate it immediately. The old key is invalidated the moment the new one is generated.

```bash
curl -X POST https://agenteats.fly.dev/owners/rotate-key \
  -H "Authorization: Bearer ae_OLD_KEY_HERE"
```

**Response:**

```json
{
  "api_key": "ae_NEW_KEY_HERE..."
}
```

> The old key stops working immediately. Update all your integrations with the new key.

---

## API Reference

### Register Owner

```
POST /owners/register
```

Creates a new owner account. No authentication required.

**Request:**

```json
{
  "name": "Your Name",
  "email": "you@restaurant.com"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | Yes | Owner's full name |
| `email` | string | Yes | Email address (must be unique) |

**Response:** `201 Created`

```json
{
  "id": "uuid",
  "name": "Your Name",
  "email": "you@restaurant.com",
  "api_key": "ae_..."
}
```

---

### Rotate API Key

```
POST /owners/rotate-key
Authorization: Bearer <current-api-key>
```

Generates a new API key and invalidates the current one.

**Response:** `200 OK`

```json
{
  "api_key": "ae_NEW_KEY..."
}
```

---

### Create Restaurant

```
POST /restaurants
Authorization: Bearer <api-key>
Content-Type: application/json
```

Creates a restaurant owned by the authenticated owner. See [Restaurant Fields](#restaurant-fields) for the full schema.

**Response:** `201 Created` — returns `RestaurantDetail`

---

### Update Restaurant

```
PUT /restaurants/{id}
Authorization: Bearer <api-key>
Content-Type: application/json
```

Updates a restaurant you own. You must own the restaurant — attempting to update another owner's restaurant returns `403 Forbidden`.

Same request body as [Create Restaurant](#create-restaurant). All fields are replaced (this is a full update, not a partial patch).

**Response:** `200 OK` — returns updated `RestaurantDetail`

---

### Add Menu Item

```
POST /restaurants/{id}/menu/items
Authorization: Bearer <api-key>
Content-Type: application/json
```

Adds a single menu item. Good for quick additions. For full menu updates, use [Bulk Import](#bulk-import-menu) instead.

**Request:**

```json
{
  "category": "Main",
  "name": "Grilled Salmon",
  "description": "Atlantic salmon with lemon butter sauce and seasonal vegetables",
  "price": 28.00,
  "currency": "USD",
  "dietary_labels": ["gluten_free"],
  "is_available": true,
  "is_popular": false,
  "calories": 520
}
```

**Response:** `201 Created` — returns `MenuItemOut`

---

### Bulk Import Menu

```
POST /restaurants/{id}/menu/import
Authorization: Bearer <api-key>
Content-Type: application/json
```

Import multiple menu items in a single request. Supports two strategies:

| Strategy | Behavior |
|----------|----------|
| `replace` (default) | Deletes **all existing items** then inserts the new ones. Use for full menu refreshes. |
| `merge` | Appends new items to the existing menu. Use for adding seasonal specials. |

**Request:**

```json
{
  "strategy": "replace",
  "items": [
    {
      "category": "Appetizer",
      "name": "Soup du Jour",
      "description": "Daily rotating soup",
      "price": 8.50,
      "is_available": true
    },
    {
      "category": "Main",
      "name": "Steak Frites",
      "description": "Grilled ribeye with hand-cut fries",
      "price": 34.00,
      "dietary_labels": ["gluten_free"],
      "is_available": true,
      "is_popular": true,
      "calories": 920
    }
  ]
}
```

**Response:** `200 OK`

```json
{
  "restaurant_id": "abc-123-...",
  "imported": 2,
  "strategy": "replace"
}
```

> **Tip:** The `replace` strategy is transactional — if any item fails validation, none are imported. Your old menu remains intact.

---

## Data Formats

### Restaurant Fields

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `name` | string | Yes | — | Restaurant name (max 200 chars) |
| `description` | string | No | — | Description for AI agents and customers |
| `cuisines` | string[] | Yes | — | Cuisine types: `["Italian", "Mediterranean"]` |
| `price_range` | string | Yes | `$$` | Price level: `$`, `$$`, `$$$`, `$$$$` |
| `address` | string | Yes | — | Street address |
| `city` | string | Yes | — | City name |
| `state` | string | No | — | State/province |
| `zip_code` | string | No | — | Postal code |
| `country` | string | No | `US` | Country code |
| `latitude` | float | No | — | GPS latitude |
| `longitude` | float | No | — | GPS longitude |
| `phone` | string | No | — | Contact phone |
| `email` | string | No | — | Contact email |
| `website` | string | No | — | Website URL |
| `features` | string[] | No | — | See available features below |
| `total_seats` | int | No | `50` | Total seating capacity (used for availability) |
| `hours` | object[] | No | — | Operating hours per day (see below) |

**Available features:**

`outdoor_seating`, `wifi`, `live_music`, `parking`, `delivery`, `takeout`, `wheelchair_accessible`, `pet_friendly`, `private_dining`, `bar`, `brunch`

**Price range guide:**

| Value | Typical per-person cost |
|-------|------------------------|
| `$` | ~$10–15 |
| `$$` | ~$15–30 |
| `$$$` | ~$30–60 |
| `$$$$` | $60+ |

### Menu Item Fields

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `category` | string | No | `Main` | Menu category: `Appetizer`, `Main`, `Dessert`, `Drink`, `Side`, etc. |
| `name` | string | Yes | — | Dish name |
| `description` | string | No | — | Brief description (helps AI agents recommend dishes) |
| `price` | float | Yes | — | Price in the specified currency |
| `currency` | string | No | `USD` | ISO currency code |
| `dietary_labels` | string[] | No | — | See available labels below |
| `is_available` | bool | No | `true` | Whether the item is currently available |
| `is_popular` | bool | No | `false` | Mark signature/popular dishes |
| `image_url` | string | No | — | URL to a dish photo |
| `calories` | int | No | — | Calorie count |

**Available dietary labels:**

`vegetarian`, `vegan`, `gluten_free`, `dairy_free`, `nut_free`, `halal`, `kosher`, `spicy`, `raw`

> **Tip:** Accurate dietary labels significantly improve recommendation matching. AI agents use these labels when users specify dietary requirements.

### Operating Hours Format

```json
{
  "day": "monday",
  "open_time": "17:00",
  "close_time": "23:00",
  "is_closed": false
}
```

| Field | Type | Description |
|-------|------|-------------|
| `day` | string | Lowercase day: `monday` through `sunday` |
| `open_time` | string | Opening time in `HH:MM` 24-hour format |
| `close_time` | string | Closing time in `HH:MM` 24-hour format |
| `is_closed` | bool | Set to `true` for days the restaurant is closed |

---

## Examples

### Full Restaurant Setup

Complete script to register, create a restaurant, and import a full menu:

```bash
#!/bin/bash
BASE="https://agenteats.fly.dev"

# 1. Register
REGISTER=$(curl -s -X POST "$BASE/owners/register" \
  -H "Content-Type: application/json" \
  -d '{"name":"Chef Marco","email":"marco@pizzaroma.com"}')

API_KEY=$(echo "$REGISTER" | python3 -c "import sys,json;print(json.load(sys.stdin)['api_key'])")
echo "API Key: $API_KEY"

# 2. Create restaurant
RESTAURANT=$(curl -s -X POST "$BASE/restaurants" \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Pizza Roma",
    "description": "Neapolitan-style pizza with imported Italian ingredients",
    "cuisines": ["Italian", "Pizza"],
    "price_range": "$$",
    "address": "456 Oak Ave",
    "city": "San Francisco",
    "state": "CA",
    "features": ["delivery", "takeout", "outdoor_seating"],
    "total_seats": 40,
    "hours": [
      {"day":"monday","open_time":"11:00","close_time":"22:00","is_closed":false},
      {"day":"tuesday","open_time":"11:00","close_time":"22:00","is_closed":false},
      {"day":"wednesday","open_time":"11:00","close_time":"22:00","is_closed":false},
      {"day":"thursday","open_time":"11:00","close_time":"22:00","is_closed":false},
      {"day":"friday","open_time":"11:00","close_time":"23:00","is_closed":false},
      {"day":"saturday","open_time":"11:00","close_time":"23:00","is_closed":false},
      {"day":"sunday","open_time":"12:00","close_time":"21:00","is_closed":false}
    ]
  }')

RESTAURANT_ID=$(echo "$RESTAURANT" | python3 -c "import sys,json;print(json.load(sys.stdin)['id'])")
echo "Restaurant ID: $RESTAURANT_ID"

# 3. Import full menu
curl -s -X POST "$BASE/restaurants/$RESTAURANT_ID/menu/import" \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "strategy": "replace",
    "items": [
      {"category":"Pizza","name":"Margherita","description":"San Marzano tomatoes, fresh mozzarella, basil","price":14.00,"dietary_labels":["vegetarian"],"is_available":true,"is_popular":true,"calories":800},
      {"category":"Pizza","name":"Diavola","description":"Spicy salami, chili flakes, mozzarella","price":16.00,"dietary_labels":["spicy"],"is_available":true,"calories":900},
      {"category":"Pizza","name":"Quattro Formaggi","description":"Mozzarella, gorgonzola, parmesan, fontina","price":17.00,"dietary_labels":["vegetarian"],"is_available":true,"calories":950},
      {"category":"Appetizer","name":"Arancini","description":"Crispy risotto balls with marinara","price":10.00,"dietary_labels":["vegetarian"],"is_available":true},
      {"category":"Dessert","name":"Panna Cotta","description":"Vanilla bean panna cotta with berry compote","price":9.00,"dietary_labels":["vegetarian","gluten_free"],"is_available":true}
    ]
  }'

echo "Done! Your restaurant is live."
```

### Seasonal Menu Update

Add seasonal specials without removing existing items:

```bash
curl -X POST "$BASE/restaurants/$RESTAURANT_ID/menu/import" \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "strategy": "merge",
    "items": [
      {"category":"Seasonal","name":"Truffle Pizza","description":"Black truffle, fontina, arugula","price":24.00,"is_available":true,"is_popular":true},
      {"category":"Seasonal","name":"Pumpkin Ravioli","description":"Roasted pumpkin, sage brown butter","price":19.00,"dietary_labels":["vegetarian"],"is_available":true}
    ]
  }'
```

### Bulk Import from CSV

Convert a CSV menu file to the JSON format using a script:

```bash
# menu.csv
# category,name,description,price,dietary_labels,is_popular
# Appetizer,Bruschetta,Tomato and basil on grilled bread,12.00,vegetarian,true
# Main,Pasta Carbonara,Classic Roman carbonara with guanciale,22.00,,true

python3 -c "
import csv, json, sys

items = []
with open('menu.csv') as f:
    reader = csv.DictReader(f)
    for row in reader:
        item = {
            'category': row['category'],
            'name': row['name'],
            'description': row['description'],
            'price': float(row['price']),
            'is_available': True,
            'is_popular': row.get('is_popular', '').lower() == 'true'
        }
        if row.get('dietary_labels'):
            item['dietary_labels'] = [l.strip() for l in row['dietary_labels'].split(';')]
        items.append(item)

print(json.dumps({'strategy': 'replace', 'items': items}, indent=2))
" > menu.json

curl -X POST "$BASE/restaurants/$RESTAURANT_ID/menu/import" \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -d @menu.json
```

---

## FAQ

### How do AI agents find my restaurant?

AI agents search by cuisine, city, price range, features, and dietary labels. To maximize visibility:
- Use accurate **cuisine tags** (e.g., `["Italian", "Pizza"]` not just `["Food"]`)
- Set correct **price range** — agents filter by budget
- Add all relevant **features** (delivery, outdoor_seating, etc.)
- Include **dietary labels** on menu items — agents use these when users specify requirements like "vegan" or "gluten-free"
- Write a good **description** — agents use it for free-text search

### How does the recommendation engine work?

When a user asks an AI agent "find me a good Italian restaurant for date night," the agent calls the recommendation endpoint. AgentEats scores restaurants on:
- Cuisine match
- Price range match
- Feature match (e.g., `live_music` for date nights)
- Dietary compatibility
- Rating and review count

### Can I manage multiple restaurants?

Yes. Each API key is tied to an owner, and one owner can create multiple restaurants. All restaurants created with your API key are yours to manage.

### Is there a cost?

AgentEats is currently free to list on.

### What happens to my data?

Your restaurant and menu data is stored in a PostgreSQL database. It's used exclusively for serving search results and recommendations to AI agents. We don't sell or share your data.

### Can I delete my restaurant?

This feature is coming soon. For now, contact us to deactivate a listing.

### How quickly do changes go live?

Immediately. There is no review queue — as soon as you create or update your restaurant or menu, AI agents can see the changes.

### What if I lose my API key?

You need your current API key to rotate to a new one. If you've completely lost access, contact us to verify ownership and reissue credentials.
