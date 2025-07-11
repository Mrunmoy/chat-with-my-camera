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
	loadConfig() // read config.yaml at startup

	// Initialize SQLite DB
	initDB()

	fmt.Println("[Go Backend] Starting ZeroMQ subscriber...")
	go runSubscriber()

	fmt.Println("[Go Backend] Starting retention job...")
	go runRetentionJob()

	fmt.Println("[Go Backend] Starting HTTP server on :8080...")
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	http.HandleFunc("/timeline", handleTimeline)

	http.Handle("/snapshots/", http.StripPrefix("/snapshots/", http.FileServer(http.Dir("./snapshots"))))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
