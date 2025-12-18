package handlers

import (
	"encoding/json"
	"net/http"

	"go_api/models"
	"go_api/storage"

	"github.com/gorilla/mux"
)

// ClientHandler handles HTTP requests for clients
type ClientHandler struct {
	store storage.Store[models.Client]
}

// NewClientHandler creates a new client handler
func NewClientHandler(store storage.Store[models.Client]) *ClientHandler {
	return &ClientHandler{store: store}
}

// GetAll handles GET /clients
func (h *ClientHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	clients := h.store.GetAll()
	json.NewEncoder(w).Encode(clients)
}

// GetByID handles GET /clients/{id}
func (h *ClientHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	client, exists := h.store.GetByID(id)

	if !exists {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Client not found"})
		return
	}

	json.NewEncoder(w).Encode(client)
}

// Create handles POST /clients
func (h *ClientHandler) Create(w http.ResponseWriter, r *http.Request) {
	var client models.Client
	if err := json.NewDecoder(r.Body).Decode(&client); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request payload"})
		return
	}

	created := h.store.Create(client)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

// Update handles PUT /clients/{id}
func (h *ClientHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var client models.Client
	if err := json.NewDecoder(r.Body).Decode(&client); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request payload"})
		return
	}

	updated, exists := h.store.Update(id, client)
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Client not found"})
		return
	}

	json.NewEncoder(w).Encode(updated)
}

// Delete handles DELETE /clients/{id}
func (h *ClientHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if !h.store.Delete(id) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Client not found"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
