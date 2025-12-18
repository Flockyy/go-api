# Quick Guide: Adding New Resources

This guide shows you how to add a new resource (e.g., Order) to the API in 5 minutes.

## Step 1: Create the Model

Create `models/order.go`:

```go
package models

import "time"

type Order struct {
	ID        string    `json:"id"`
	ClientID  string    `json:"client_id"`
	ItemID    string    `json:"item_id"`
	Quantity  int       `json:"quantity"`
	Total     float64   `json:"total"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
```

## Step 2: Create the Handler

Create `handlers/order_handler.go`:

```go
package handlers

import (
	"encoding/json"
	"net/http"

	"go_api/models"
	"go_api/storage"

	"github.com/gorilla/mux"
)

type OrderHandler struct {
	store storage.Store[models.Order]
}

func NewOrderHandler(store storage.Store[models.Order]) *OrderHandler {
	return &OrderHandler{store: store}
}

func (h *OrderHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	orders := h.store.GetAll()
	json.NewEncoder(w).Encode(orders)
}

func (h *OrderHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	order, exists := h.store.GetByID(id)

	if !exists {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Order not found"})
		return
	}

	json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request payload"})
		return
	}

	created := h.store.Create(order)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func (h *OrderHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request payload"})
		return
	}

	updated, exists := h.store.Update(id, order)
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Order not found"})
		return
	}

	json.NewEncoder(w).Encode(updated)
}

func (h *OrderHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if !h.store.Delete(id) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Order not found"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
```

## Step 3: Update Storage (if needed)

If your model has different fields, update `storage/store.go` in the `Create` and `Update` methods:

```go
// Add this case in Create method
case *models.Order:
	v.ID = uuid.New().String()
	v.CreatedAt = time.Now()
	v.UpdatedAt = time.Now()
	s.items[v.ID] = any(*v).(T)

// Add this case in Update method  
case *models.Order:
	old := s.items[id]
	oldOrder := any(old).(models.Order)
	v.ID = id
	v.CreatedAt = oldOrder.CreatedAt
	v.UpdatedAt = time.Now()
	s.items[id] = any(*v).(T)
```

## Step 4: Register Routes

Update `router/router.go`:

```go
// Update function signature
func Setup(
	itemHandler *handlers.ItemHandler, 
	clientHandler *handlers.ClientHandler,
	orderHandler *handlers.OrderHandler,  // Add this
) *mux.Router {
	// ... existing code ...

	// Add order routes
	api.HandleFunc("/orders", orderHandler.GetAll).Methods("GET")
	api.HandleFunc("/orders", orderHandler.Create).Methods("POST")
	api.HandleFunc("/orders/{id}", orderHandler.GetByID).Methods("GET")
	api.HandleFunc("/orders/{id}", orderHandler.Update).Methods("PUT")
	api.HandleFunc("/orders/{id}", orderHandler.Delete).Methods("DELETE")

	// ... rest of the code ...
}
```

## Step 5: Initialize in Main

Update `main.go`:

```go
func main() {
	// Initialize stores
	itemStore := storage.NewMemoryStore[models.Item]()
	clientStore := storage.NewMemoryStore[models.Client]()
	orderStore := storage.NewMemoryStore[models.Order]()  // Add this

	// Initialize handlers
	itemHandler := handlers.NewItemHandler(itemStore)
	clientHandler := handlers.NewClientHandler(clientStore)
	orderHandler := handlers.NewOrderHandler(orderStore)  // Add this

	// Setup router
	r := router.Setup(itemHandler, clientHandler, orderHandler)  // Add orderHandler

	// ... rest of the code (add log entries if desired) ...
}
```

## Done! 🎉

Your new resource is now available at:
- `GET    /api/v1/orders`
- `POST   /api/v1/orders`
- `GET    /api/v1/orders/{id}`
- `PUT    /api/v1/orders/{id}`
- `DELETE /api/v1/orders/{id}`

Test it:
```bash
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{"client_id":"abc","item_id":"xyz","quantity":2,"total":199.99,"status":"pending"}'
```

## Key Points

- **No changes needed** to middleware or core infrastructure
- **Type-safe** - Compiler catches errors
- **Consistent** - All resources follow the same pattern
- **Fast** - Takes ~5 minutes to add a new resource
- **Scalable** - Can easily switch storage backend later
