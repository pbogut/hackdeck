package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/pbogut/hackdeck/pkg/types"
)

// pingHandler handles the /ping route
func ReloadHandler(w http.ResponseWriter, r *http.Request) {
	ReloadConfig()
	// Create the response object
	response := types.ReloadResponse{
		ConfigReloaded: true,
	}
	// Set the content type to application/json
	w.Header().Set("Content-Type", "application/json")
	// Encode the response as JSON and write it to the response writer
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
	}
}

// pingHandler handles the /ping route
func PingHandler(w http.ResponseWriter, r *http.Request) {
	// Get the hostname
	hostname, err := os.Hostname()
	if err != nil {
		http.Error(w, "Unable to get hostname", http.StatusInternalServerError)
		return
	}
	// Create the response object
	response := types.PingResponse{
		MachineName: hostname,
	}
	// Set the content type to application/json
	w.Header().Set("Content-Type", "application/json")
	// Encode the response as JSON and write it to the response writer
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
	}
}
