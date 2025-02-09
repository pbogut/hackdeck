package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/pbogut/hackdeck/pkg/handlers"
	"github.com/pbogut/hackdeck/pkg/logger"
)

func main() {
	port := flag.Int("port", 8191, "Port to listen on")
	host := flag.String("host", "", "host to listen on")
	printDebug := flag.Bool("debug", false, "Print debug messages")
	flag.Parse()

	addr := fmt.Sprintf("%s:%d", *host, *port)
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
	fmt.Printf("Starting server on %s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		logger.Fatal("Error starting server:", err)
	}
}
