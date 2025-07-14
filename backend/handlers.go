package main

import (
	// "database/sql"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
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

// === Structures for Ollama ===
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
}

type ChatChoice struct {
	Message ChatMessage `json:"message"`
}

type ChatResponse struct {
	Choices []ChatChoice `json:"choices"`
}

func (app *App) handleChat(w http.ResponseWriter, r *http.Request) {
	// === CORS ===
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// === Parse request ===
	var req struct {
		CameraID string `json:"camera_id"`
		Message  string `json:"message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	log.Printf("handleChat: camera_id=%s message=%s", req.CameraID, req.Message)

	// === STEP 1: Extract object(s) ===
	extractionPrompt := fmt.Sprintf(
		"User question: %s\n\nExtract the main object(s) or labels the user wants to know about. Return ONLY a JSON array, e.g. [\"car\"] or [\"person\", \"dog\"]. If no object, return [].",
		req.Message)

	extractionReq := ChatRequest{
		Model: "llama3",
		Messages: []ChatMessage{
			{Role: "system", Content: "You extract objects only. No explanation."},
			{Role: "user", Content: extractionPrompt},
		},
	}

	bodyBytes, _ := json.Marshal(extractionReq)
	resp, err := http.Post(
		"http://localhost:11434/v1/chat/completions",
		"application/json",
		bytes.NewBuffer(bodyBytes),
	)
	if err != nil {
		log.Printf("Ollama extract POST failed: %v", err)
		http.Error(w, "LLM extraction failed", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Ollama extract status %d", resp.StatusCode)
		http.Error(w, "LLM extraction non-200", http.StatusInternalServerError)
		return
	}

	var extractResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&extractResp); err != nil {
		log.Printf("Decode extract response failed: %v", err)
		http.Error(w, "Invalid extract response", http.StatusInternalServerError)
		return
	}

	// Parse the extracted JSON array
	extracted := extractResp.Choices[0].Message.Content
	log.Printf("handleChat: extracted raw: %s", extracted)

	var objects []string
	if err := json.Unmarshal([]byte(extracted), &objects); err != nil {
		log.Printf("JSON unmarshal failed: %v", err)
		objects = []string{}
	}
	log.Printf("handleChat: extracted objects: %v", objects)

	// === STEP 2: Query DB for most recent matching detection ===
	var contextString string

	if len(objects) > 0 {
		// Use first extracted object for now
		object := objects[0]
		log.Printf("Searching for object: %s", object)

		row := app.DB.QueryRow(`
			SELECT timestamp, labels FROM detections
			WHERE camera_id = ? AND labels LIKE ?
			ORDER BY timestamp DESC LIMIT 1
		`, req.CameraID, "%"+object+"%")

		var ts float64
		var labels string
		err := row.Scan(&ts, &labels)
		if err == nil {
			t := time.Unix(int64(ts), 0).Format(time.RFC3339)
			contextString = fmt.Sprintf("- Last detection: %s Labels: %s\n", t, labels)
		} else {
			contextString = fmt.Sprintf("No detections found for '%s'.", object)
		}

	} else {
		// No object found → fallback to latest 5
		rows, err := app.DB.Query(`
			SELECT timestamp, labels FROM detections
			WHERE camera_id = ?
			ORDER BY timestamp DESC LIMIT 5
		`, req.CameraID)
		if err != nil {
			log.Printf("Fallback DB query failed: %v", err)
			contextString = "No detection history available."
		} else {
			defer rows.Close()
			for rows.Next() {
				var ts float64
				var labels string
				rows.Scan(&ts, &labels)
				t := time.Unix(int64(ts), 0).Format(time.RFC3339)
				contextString += fmt.Sprintf("- Time: %s Labels: %s\n", t, labels)
			}
			if contextString == "" {
				contextString = "No recent detections found."
			}
		}
	}

	// === STEP 3: Final prompt ===
	finalPrompt := fmt.Sprintf(
		"Camera: %s\n\nDetection context:\n%s\n\nUser question: %s",
		req.CameraID, contextString, req.Message)

	log.Printf("handleChat: final prompt:\n%s", finalPrompt)

	// === STEP 4: Send final prompt to Ollama ===
	finalReq := ChatRequest{
		Model: "llama3",
		Messages: []ChatMessage{
			{Role: "system", Content: "You are a helpful camera assistant."},
			{Role: "user", Content: finalPrompt},
		},
	}

	bodyBytes, _ = json.Marshal(finalReq)
	resp2, err := http.Post(
		"http://localhost:11434/v1/chat/completions",
		"application/json",
		bytes.NewBuffer(bodyBytes),
	)
	if err != nil {
		log.Printf("Ollama final POST failed: %v", err)
		http.Error(w, "LLM final call failed", http.StatusInternalServerError)
		return
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusOK {
		log.Printf("Ollama final returned status %d", resp2.StatusCode)
		http.Error(w, "LLM final non-200", http.StatusInternalServerError)
		return
	}

	var finalResp ChatResponse
	if err := json.NewDecoder(resp2.Body).Decode(&finalResp); err != nil {
		log.Printf("Decode final response failed: %v", err)
		http.Error(w, "Invalid final response", http.StatusInternalServerError)
		return
	}

	answer := finalResp.Choices[0].Message.Content
	log.Printf("handleChat: final answer: %s", answer)

	// === Return to frontend ===
	json.NewEncoder(w).Encode(map[string]string{
		"answer": answer,
	})
}
