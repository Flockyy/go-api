package router

import (
	"go-api/handlers"
	"go-api/middleware"

	"github.com/gorilla/mux"
)

// Setup configures all routes and middleware
func Setup(itemHandler *handlers.ItemHandler, clientHandler *handlers.ClientHandler) *mux.Router {
	router := mux.NewRouter()

	// API v1 routes
	api := router.PathPrefix("/api/v1").Subrouter()

	// Health check
	api.HandleFunc("/health", handlers.HealthCheck).Methods("GET")

	// Item routes
	api.HandleFunc("/items", itemHandler.GetAll).Methods("GET")
	api.HandleFunc("/items", itemHandler.Create).Methods("POST")
	api.HandleFunc("/items/{id}", itemHandler.GetByID).Methods("GET")
	api.HandleFunc("/items/{id}", itemHandler.Update).Methods("PUT")
	api.HandleFunc("/items/{id}", itemHandler.Delete).Methods("DELETE")

	// Client routes
	api.HandleFunc("/clients", clientHandler.GetAll).Methods("GET")
	api.HandleFunc("/clients", clientHandler.Create).Methods("POST")
	api.HandleFunc("/clients/{id}", clientHandler.GetByID).Methods("GET")
	api.HandleFunc("/clients/{id}", clientHandler.Update).Methods("PUT")
	api.HandleFunc("/clients/{id}", clientHandler.Delete).Methods("DELETE")

	// Global middleware
	router.Use(middleware.Logging)
	router.Use(middleware.JSON)
	router.Use(middleware.CORS)

	return router
}
