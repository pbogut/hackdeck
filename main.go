package main

import (
	"fmt"
	"net/http"

	"github.com/pbogut/hackdeck/pkg/handlers"
)

func main() {
	http.HandleFunc("/", handlers.WsHandler)
	http.HandleFunc("/ping", handlers.PingHandler)
	// Start the HTTP server on port 8191
	fmt.Println("Starting server on :8191")
	if err := http.ListenAndServe(":8191", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
