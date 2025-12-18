# Go API - Scalable Architecture

A simple, scalable REST API built with Go that demonstrates clean architecture principles. Easily extensible to support multiple resources (items, clients, orders, etc.).

## Features

- **Clean Package Structure** - Organized into logical packages
- **Generic Storage Layer** - Type-safe storage with Go generics
- **RESTful API** - Proper HTTP methods and status codes
- **Thread-Safe** - Concurrent request handling with mutex locks
- **Middleware Pipeline** - Logging, CORS, JSON headers
- **API Versioning** - Ready for v2, v3, etc.
- **Easily Extensible** - Add new resources in minutes

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                         main.go                             │
│  (Dependency Injection & Server Bootstrap)                  │
└────────────┬────────────────────────────────────────────────┘
             │
             ├─── Initializes Stores
             ├─── Initializes Handlers  
             └─── Configures Router
                         │
        ┌────────────────┴────────────────┐
        │                                  │
        ▼                                  ▼
┌───────────────┐                 ┌────────────────┐
│  router/      │◄────────────────┤  middleware/   │
│  router.go    │   Applies        │  middleware.go │
│               │   middleware     │                │
│ - /api/v1     │                  │ - Logging      │
│ - Routes      │                  │ - CORS         │
│ - Versioning  │                  │ - JSON headers │
└───────┬───────┘                  └────────────────┘
        │
        │ Routes requests to handlers
        │
        ▼
┌─────────────────────────────────────────┐
│            handlers/                     │
│  ┌─────────────────┐ ┌────────────────┐│
│  │ item_handler.go │ │client_handler.go││
│  │                 │ │                 ││
│  │ - GetAll()      │ │ - GetAll()      ││
│  │ - GetByID()     │ │ - GetByID()     ││
│  │ - Create()      │ │ - Create()      ││
│  │ - Update()      │ │ - Update()      ││
│  │ - Delete()      │ │ - Delete()      ││
│  └────────┬────────┘ └────────┬────────┘│
└───────────┼───────────────────┼─────────┘
            │                   │
            │ Uses Store[Item]  │ Uses Store[Client]
            │                   │
            ▼                   ▼
┌────────────────────────────────────────┐
│         storage/store.go                │
│  ┌──────────────────────────────────┐  │
│  │  Store[T any] Interface          │  │
│  │  (Generic Type-Safe Storage)     │  │
│  ├──────────────────────────────────┤  │
│  │  MemoryStore[T] Implementation   │  │
│  │  - In-memory with sync.RWMutex   │  │
│  │  - Thread-safe operations        │  │
│  └──────────────────────────────────┘  │
└────────────────┬───────────────────────┘
                 │
                 │ Manages
                 │
                 ▼
┌────────────────────────────────────────┐
│            models/                      │
│  ┌──────────────┐  ┌─────────────────┐│
│  │  item.go     │  │  client.go      ││
│  │              │  │                 ││
│  │ - ID         │  │ - ID            ││
│  │ - Name       │  │ - Name          ││
│  │ - Desc       │  │ - Email         ││
│  │ - Timestamps │  │ - Phone         ││
│  │              │  │ - Timestamps    ││
│  └──────────────┘  └─────────────────┘│
└────────────────────────────────────────┘

Flow: Request → Router → Middleware → Handler → Storage → Model
```

**Key Design Patterns:**
- **Dependency Injection** - Handlers receive stores via constructor
- **Repository Pattern** - Storage layer abstracts data access
- **Generics** - Type-safe storage works with any model
- **Clean Architecture** - Clear separation of concerns

## Prerequisites

- Go 1.18+ (requires generics support)

## Installation

1. Install dependencies:
```bash
go mod download
```

## Running the API

Start the server:
```bash
go run main.go
```

The server will start on `http://localhost:8080`

## API Endpoints

All endpoints are under `/api/v1`

### Health Check
```
GET /api/v1/health
```

### Items
```
GET    /api/v1/items         # List all items
POST   /api/v1/items         # Create item
GET    /api/v1/items/{id}    # Get item by ID
PUT    /api/v1/items/{id}    # Update item
DELETE /api/v1/items/{id}    # Delete item
```

### Clients
```
GET    /api/v1/clients       # List all clients
POST   /api/v1/clients       # Create client
GET    /api/v1/clients/{id}  # Get client by ID
PUT    /api/v1/clients/{id}  # Update client
DELETE /api/v1/clients/{id}  # Delete client
```

## Example Usage

### Create an item
```bash
curl -X POST http://localhost:8080/api/v1/items \
  -H "Content-Type: application/json" \
  -d '{"name":"Laptop","description":"MacBook Pro 16-inch"}'
```

### Create a client
```bash
curl -X POST http://localhost:8080/api/v1/clients \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com","phone":"555-1234"}'
```

### Get all items
```bash
curl http://localhost:8080/api/v1/items
```

### Get all clients
```bash
curl http://localhost:8080/api/v1/clients
```

### Update an item
```bash
curl -X PUT http://localhost:8080/api/v1/items/{id} \
  -H "Content-Type: application/json" \
  -d '{"name":"Updated Item","description":"Updated description"}'
```

### Delete a client
```bash
curl -X DELETE http://localhost:8080/api/v1/clients/{id}
```

## Adding New Resources

This architecture makes it easy to add new resources. Here's how:

### 1. Create a Model
Add a new file in `models/` (e.g., `order.go`):
```go
package models

import "time"

type Order struct {
    ID        string    `json:"id"`
    ClientID  string    `json:"client_id"`
    Amount    float64   `json:"amount"`
    Status    string    `json:"status"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### 2. Create a Handler
Add a handler in `handlers/` (e.g., `order_handler.go`):
```go
package handlers

import (
    "go_api/models"
    "go_api/storage"
    // ... similar structure to item_handler.go
)

type OrderHandler struct {
    store storage.Store[models.Order]
}

func NewOrderHandler(store storage.Store[models.Order]) *OrderHandler {
    return &OrderHandler{store: store}
}

// Add GetAll, GetByID, Create, Update, Delete methods...
```

### 3. Register Routes
Update `router/router.go`:
```go
func Setup(itemHandler, clientHandler, orderHandler *handlers...) {
    // ... existing code
    
    // Order routes
    api.HandleFunc("/orders", orderHandler.GetAll).Methods("GET")
    api.HandleFunc("/orders", orderHandler.Create).Methods("POST")
    // ... add other routes
}
```

### 4. Initialize in Main
Update `main.go`:
```go
orderStore := storage.NewMemoryStore[models.Order]()
orderHandler := handlers.NewOrderHandler(orderStore)
r := router.Setup(itemHandler, clientHandler, orderHandler)
```

That's it! Your new resource is ready to use.

## Scalability Features

### Current Implementation
- **Generic Storage** - Type-safe, works with any model
- **In-Memory Store** - Fast for development and testing
- **Thread-Safe** - Handles concurrent requests

### Easy Upgrades
- **Database Backend** - Store interface can be implemented with PostgreSQL, MongoDB, etc.
- **Caching Layer** - Add Redis or similar
- **Authentication** - Add JWT middleware
- **Rate Limiting** - Add rate limiter middleware
- **Validation** - Add validator middleware
- **API Versioning** - Already using `/api/v1` prefix

## Project Structure Explained

- **`models/`** - Business domain models. Add new resource types here.
- **`storage/`** - Data persistence layer. Uses generics for type safety. Swap implementations easily.
- **`handlers/`** - HTTP handlers for each resource. Thin layer, delegates to storage.
- **`middleware/`** - Cross-cutting concerns (logging, CORS, auth, etc.)
- **`router/`** - Centralized route configuration. Single source of truth for all endpoints.
- **`main.go`** - Application bootstrap. Wire dependencies, start server.
