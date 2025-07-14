package main

import (
	// "database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// timelineHandler handles GET requests to /timeline
// It queries the SQLite 'detections' table and returns matching events as JSON.
func (app *App) handleTimeline(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// === Parse query params ===
	// Example: /timeline?camera_id=garage_webcam&label=person&start_time=...&end_time=...
	cameraID := r.URL.Query().Get("camera_id")
	label := r.URL.Query().Get("label")
	startTimeStr := r.URL.Query().Get("start_time")
	endTimeStr := r.URL.Query().Get("end_time")

	log.Printf("[TimelineHandler] camera_id: %s, label: %s, start_time: %s, end_time: %s\n", cameraID, label, startTimeStr, endTimeStr)

	// === Build WHERE conditions and arguments ===
	var conditions []string
	var args []interface{}

	if cameraID != "" {
		conditions = append(conditions, "camera_id = ?")
		args = append(args, cameraID)
	}

	if label != "" {
		// The 'labels' column is JSON text, so we use LIKE to match substrings.
		// Because labels are stored as JSON text (["person"]), so we match substrings with LIKE '%person%'.
		conditions = append(conditions, "labels LIKE ?")
		args = append(args, "%"+label+"%")
	}

	if startTimeStr != "" {
		conditions = append(conditions, "timestamp >= ?")
		startTime, err := strconv.ParseFloat(startTimeStr, 64)
		if err == nil {
			args = append(args, startTime)
		}
	}

	if endTimeStr != "" {
		conditions = append(conditions, "timestamp <= ?")
		endTime, err := strconv.ParseFloat(endTimeStr, 64)
		if err == nil {
			args = append(args, endTime)
		}
	}

	// === Final SQL query ===
	query := "SELECT timestamp, camera_id, labels, boxes, snapshot_file FROM detections"
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	// So we don’t dump thousands of rows on a single query. Can add ?limit later.
	query += " ORDER BY timestamp DESC LIMIT 100" // Limit to 100 results for now

	log.Printf("Timeline query: %s\n", query)

	// === Execute query ===
	rows, err := app.DB.Query(query, args...)
	if err != nil {
		http.Error(w, "Query failed", http.StatusInternalServerError)
		log.Printf("Timeline query error: %v", err)
		return
	}
	defer rows.Close()

	// === Collect rows ===
	var results []map[string]interface{}

	log.Printf("results: %v\n", results)
	for rows.Next() {
		var ts float64
		var cid, labels, boxes, snapshotFile string

		if err := rows.Scan(&ts, &cid, &labels, &boxes, &snapshotFile); err != nil {
			log.Printf("Timeline row scan failed: %v", err)
			continue
		}

		// Build safe snapshot URL (strip ./snapshots/)
		snapshotURL := ""
		if snapshotFile != "" {
			snapshotURL = fmt.Sprintf("/snapshots/%s", filepath.Base(snapshotFile))
		}

		results = append(results, map[string]interface{}{
			"timestamp":     ts,
			"camera_id":     cid,
			"labels":        labels,
			"boxes":         boxes,
			"snapshot_file": snapshotFile, // raw path, for debug
			"snapshot_url":  snapshotURL,  // public URL via static file server
		})
	}

	// === Return JSON response ===
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// handleSnapshot serves snapshot image files from your ./snapshots directory.
// Example: GET /snapshot?file=garage_webcam_1752205055.jpg
func handleSnapshot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	fileName := r.URL.Query().Get("file")
	if fileName == "" {
		http.Error(w, "Missing 'file' query parameter", http.StatusBadRequest)
		return
	}

	// Secure the path to avoid path traversal
	safePath := filepath.Join("./snapshots", filepath.Base(fileName))

	file, err := os.Open(safePath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		log.Printf("Snapshot file not found: %s", safePath)
		return
	}
	defer file.Close()

	// Serve as image/jpeg
	w.Header().Set("Content-Type", "image/jpeg")
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "Failed to serve snapshot", http.StatusInternalServerError)
		log.Printf("Snapshot serve failed: %v", err)
	}
}

// handleLatest returns the latest detection for a camera
func (app *App) handleLatest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	cameraID := r.URL.Query().Get("camera_id")
	if cameraID == "" {
		http.Error(w, "Missing camera_id", http.StatusBadRequest)
		return
	}

	row := app.DB.QueryRow(
		"SELECT timestamp, snapshot_file FROM detections WHERE camera_id = ? ORDER BY timestamp DESC LIMIT 1",
		cameraID,
	)

	var timestamp float64
	var snapshotFile string

	err := row.Scan(&timestamp, &snapshotFile)
	if err != nil {
		http.Error(w, "No detections found", http.StatusNotFound)
		return
	}

	result := map[string]interface{}{
		"timestamp":     timestamp,
		"snapshot_file": snapshotFile,
	}

	json.NewEncoder(w).Encode(result)
}

// camerasHandler returns all cameras from your config.
func (app *App) camerasHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	cameras := []CameraInfo{}
	for i, cam := range app.Config.Cameras {
		cameras = append(cameras, CameraInfo{
			ID:        cam.ID,
			Number:    i + 1,         // 1-based serial number
			Thumbnail: cam.Thumbnail, // from config.json
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(cameras); err != nil {
		http.Error(w, "Failed to encode cameras", http.StatusInternalServerError)
		log.Printf("Failed to write /cameras response: %v", err)
	}
}

// handleChat handles POST /chat requests.
// It receives { camera_id, message } JSON and returns { answer: "..." } JSON.// handleChat handles POST /chat for LLM queries.
// It returns a fake answer for now, with full CORS handling.
func (app *App) handleChat(w http.ResponseWriter, r *http.Request) {
	// === CORS headers ===
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	// === Handle preflight ===
	if r.Method == http.MethodOptions {
		return
	}

	// === Only allow POST ===
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// === Parse request body ===
	var req struct {
		CameraID string `json:"camera_id"`
		Message  string `json:"message"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Invalid request body: %v", err)
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	log.Printf("handleChat: camera_id=%s message=%s", req.CameraID, req.Message)

	// === Fake LLM reply ===
	fakeAnswer := "Pretend I'm your LLM! You asked: '" +
		req.Message + "' about camera '" + req.CameraID + "'."

	// === Send response ===
	resp := map[string]string{
		"answer": fakeAnswer,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Printf("handleChat: replied with fake answer")
}
