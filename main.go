package main

import (
	"flag"
	"net/http"

	"github.com/pbogut/hackdeck/pkg/handlers"
	"github.com/pbogut/hackdeck/pkg/logger"
)

func main() {
	printDebug := flag.Bool("debug", false, "Print debug messages")
	flag.Parse()

	if *printDebug {
		logger.Init(logger.DEBUG)
	} else {
		logger.Init(logger.INFO)
	}

	handlers.Init()

	http.HandleFunc("/", handlers.WsHandler)
	http.HandleFunc("/ping", handlers.PingHandler)
	http.HandleFunc("/reload", handlers.ReloadHandler)
	// Start the HTTP server on port 8191
	logger.Info("Starting server on :8191")
	if err := http.ListenAndServe(":8191", nil); err != nil {
		logger.Fatal("Error starting server:", err)
	}
}
