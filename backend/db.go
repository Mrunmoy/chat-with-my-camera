package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"fmt"
	"log"
	"os"
)

var db *sql.DB

/*
db.go
------

Handles:
- Opening SQLite database
- Creating tables if they do not exist
- Provides functions to insert detection events
- Future: query timeline, delete old rows for retention
*/

// initDB opens (or creates) the SQLite DB and ensures schema exists.
func initDB() {
	// Make sure ./data exists
	if _, err := os.Stat("./data"); os.IsNotExist(err) {
		err = os.MkdirAll("./data", os.ModePerm)
		if err != nil {
			log.Fatalf("Failed to create data directory: %v", err)
		}
	}

	var err error
	db, err = sql.Open("sqlite3", "./data/detections.db")
	if err != nil {
		log.Fatalf("Failed to open SQLite DB: %v", err)
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS detections (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		timestamp REAL,
		camera_id TEXT,
		labels TEXT,
		boxes TEXT,
		snapshot_file TEXT
	);
	`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	fmt.Println("[DB] SQLite initialized and table ready.")
}

// insertDetection inserts a detection event into the DB.
func insertDetection(timestamp float64, cameraID string, labels string, boxes string, snapshotPath string) error {
	stmt := `INSERT INTO detections (timestamp, camera_id, labels, boxes, snapshot_file) VALUES (?, ?, ?, ?, ?)`
	_, err := db.Exec(stmt, timestamp, cameraID, labels, boxes, snapshotPath)
	return err
}
