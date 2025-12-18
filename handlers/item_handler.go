package handlers

import (
	"encoding/json"
	"net/http"

	"go-api/models"
	"go-api/storage"

	"github.com/gorilla/mux"
)

// ItemHandler handles HTTP requests for items
type ItemHandler struct {
	store storage.Store[models.Item]
}

// NewItemHandler creates a new item handler
func NewItemHandler(store storage.Store[models.Item]) *ItemHandler {
	return &ItemHandler{store: store}
}

// GetAll handles GET /items
func (h *ItemHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	items := h.store.GetAll()
	json.NewEncoder(w).Encode(items)
}

// GetByID handles GET /items/{id}
func (h *ItemHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	item, exists := h.store.GetByID(id)

	if !exists {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Item not found"})
		return
	}

	json.NewEncoder(w).Encode(item)
}

// Create handles POST /items
func (h *ItemHandler) Create(w http.ResponseWriter, r *http.Request) {
	var item models.Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request payload"})
		return
	}

	created := h.store.Create(item)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

// Update handles PUT /items/{id}
func (h *ItemHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var item models.Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request payload"})
		return
	}

	updated, exists := h.store.Update(id, item)
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Item not found"})
		return
	}

	json.NewEncoder(w).Encode(updated)
}

// Delete handles DELETE /items/{id}
func (h *ItemHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if !h.store.Delete(id) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Item not found"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
