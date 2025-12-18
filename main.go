package main

import (
	"log"
	"net/http"

	"go_api/handlers"
	"go_api/models"
	"go_api/router"
	"go_api/storage"
)

func main() {
	// Initialize stores
	itemStore := storage.NewMemoryStore[models.Item]()
	clientStore := storage.NewMemoryStore[models.Client]()

	// Initialize handlers
	itemHandler := handlers.NewItemHandler(itemStore)
	clientHandler := handlers.NewClientHandler(clientStore)

	// Setup router
	r := router.Setup(itemHandler, clientHandler)

	// Start server
	port := ":8080"
	log.Printf("Server starting on http://localhost%s", port)
	log.Printf("API endpoints:")
	log.Printf("  - GET    /api/v1/health")
	log.Printf("  - GET    /api/v1/items")
	log.Printf("  - POST   /api/v1/items")
	log.Printf("  - GET    /api/v1/items/{id}")
	log.Printf("  - PUT    /api/v1/items/{id}")
	log.Printf("  - DELETE /api/v1/items/{id}")
	log.Printf("  - GET    /api/v1/clients")
	log.Printf("  - POST   /api/v1/clients")
	log.Printf("  - GET    /api/v1/clients/{id}")
	log.Printf("  - PUT    /api/v1/clients/{id}")
	log.Printf("  - DELETE /api/v1/clients/{id}")
	
	log.Fatal(http.ListenAndServe(port, r))
}
