package main

import (
	"fmt"
	"log"
	"net/http"
)

/*
main.go
--------

This is the entry point for the Go backend service.

Responsibilities:
- Start ZeroMQ subscriber (receives detection events from Python YOLO)
- Initialize SQLite database connection
- Start retention cleanup loop (rolling window)
- Serve HTTP API endpoints (health check, timeline)
*/

func main() {
	config := loadConfig() // ✅ Load YAML once

	// Init DB once
	initDB()

	// Create your app instance
	app := NewApp(db, &config)

	// Start background jobs
	fmt.Println("[Go Backend] Starting ZeroMQ subscriber...")
	go app.runSubscriber()

	fmt.Println("[Go Backend] Starting retention job...")
	go app.runRetention()

	// Use your own ServeMux
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// API endpoints use the App methods
	mux.HandleFunc("/timeline", app.handleTimeline)
	mux.HandleFunc("/cameras", app.camerasHandler)
	mux.HandleFunc("/snapshot", handleSnapshot)
	mux.HandleFunc("/latest", app.handleLatest)

	// Static file servers
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	mux.Handle("/snapshots/", http.StripPrefix("/snapshots/", http.FileServer(http.Dir("./snapshots"))))

	// Start server with your mux
	fmt.Println("[Go Backend] HTTP server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
