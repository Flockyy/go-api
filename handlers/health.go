package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

// HealthCheck handles GET /health
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	})
}
