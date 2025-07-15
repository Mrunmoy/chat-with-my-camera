# Backend

![Go](https://img.shields.io/badge/Go-1.21-blue?logo=go&logoColor=white)

## Overview

This folder contains the **Go backend** for `chat-with-my-camera`.

- Receive detection events via ZeroMQ subscriber.
- Store detections in a **SQLite** database.
- Save snapshot image files for each detection.
- Serve a **REST API**: `/timeline`, `/snapshots`, `/cameras`.
- Provide **LLM chat** endpoint: smart 2-step object extraction + timeline lookup.
- Periodically run retention cleanup to keep DB and disk small.

## Key Files

- `main.go` — entry point: spins up ZeroMQ subscriber + HTTP server.
- `db.go` — database connection, schema, and CRUD helpers.
- `handlers.go` — REST API routes: `/timeline`, `/snapshots`, `/cameras`, `/chat`.
- `retention.go` — deletes old rows & images past retention window.
- `go.mod`, `go.sum` — Go dependencies.

## How to Run

```bash
cd backend
go run main.go
```

Default port: **8080**. Snapshot images must be accessible under `./snapshots` relative to the backend.

## API Endpoints

- `GET /timeline?camera_id=...&start_time=...&end_time=...` → JSON of detections.
- `GET /snapshots/...` → serve saved JPEGs.
- `GET /cameras` → all configured cameras.
- `POST /chat` → JSON `{ camera_id, message }` → auto-extract objects → query timeline → call local Ollama → return `{ answer }`.


## API Responses — Example JSON

### `/cameras`

```json
[
  {
    "id": "garage_webcam",
    "type": "webcam",
    "thumbnail": "thumbnails/webcam.png"
  }
]
```

### `/timeline`

```json
[
  {
    "timestamp": 1752205052.2,
    "camera_id": "lounge_rtsp",
    "labels": ["car"],
    "boxes": [[100, 200, 300, 400]],
    "snapshot_url": "/snapshots/lounge_rtsp_1752205052.jpg"
  }
]
```

### `/chat` (POST)

**Request:**

```json
{
  "camera_id": "garage_webcam",
  "message": "When did you last see a car?"
}
```

**Response:**

```json
{
  "answer": "The last car was seen at 2 PM yesterday."
}
```

## How the LLM Works

- Uses **2-step pipeline**: first extract relevant objects/labels → then query timeline for real data → then build final prompt → then answer.
- Ollama runs locally at `localhost:11434`.
- All LLM calls stay local, no cloud API.

### LLM Chat Pipeline

- The backend uses a **2-step LLM flow**:

  1. **Extraction step:** The user’s question goes to Ollama with a system prompt to extract keywords (e.g., "car").
     ```go
     extractionPrompt := fmt.Sprintf("Extract object: %s", req.Message)
     // POST to http://localhost:11434/v1/chat/completions
     ```
  2. **Timeline lookup:** The backend runs a SQLite query for the latest matching detection:
     ```sql
     SELECT timestamp, labels FROM detections WHERE camera_id=? AND labels LIKE '%car%' ORDER BY timestamp DESC LIMIT 1
     ```
  3. **Final step:** Builds a new prompt that includes the detection context and sends it back to Ollama to generate a natural answer.
     ```go
     finalPrompt := fmt.Sprintf("Camera: %s\n\nContext: %s\n\nUser question: %s", req.CameraID, context, req.Message)
     ```

- The `/chat` endpoint returns `{ "answer": "..." }`.

This ensures the LLM only talks about real detections. No cloud API, everything stays local!


## Retention Job

Run automatically at startup. Cleans up:

- Old DB rows older than N days.
- Deletes matching snapshot files from disk.

## Tips

- Use `log.Printf` for debugging timeline queries. 
- Always run Ollama `serve` before calling `/chat`. 
- Place snapshots under `./snapshots`.
