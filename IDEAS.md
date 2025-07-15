# Ideas Parking Lot

### What’s working now
- [x] Local webcam source → tested!
- [x] RTSP camera source → tested with Unifi Protect!
- [x] YOLOv8 detection loop → tested!
- [x] ZeroMQ publisher for detection events → working!
- [x] Config loader (YAML) → flexible webcam/RTSP.
- [x] Multi-camera support → multiple streams, single grid view.
- [x] Dynamic grid layout → auto 1x1, 2x2, 3x3, 4x4.
- [x] Robust reconnect logic for webcam & RTSP.
- [x] Per-camera deduplication (labels only) & throttling.
- [x] Offline placeholder overlay → tested with real unplug/reboot.
- [x] Timeline API → `/timeline` with filters & query params.
- [x] Static snapshot server → `/snapshots/` serving direct images.
- [x] Timeline JSON now includes `snapshot_url` for front-end or HA.
- [x] Retention job: deletes old DB rows AND snapshot files → fully tested!
- [x] **LLM chat pipeline (local Ollama)**
  - Extracts objects/labels from plain user questions
  - Looks up last matching detection in SQLite timeline
  - Returns smart context-aware answer: “I last saw a car at 2 PM”
  - Fully local, no cloud — fast and free!

---

## Up Next

### Home Assistant Integration
- Publish smart events: motion/person detected → turn on lights, send notifications.

### AI Pipeline Ideas
- Pass snapshots for richer context
- Multi-camera queries (“Check all cameras for cars”)
- Natural-language filters for timeline ranges (“last week”, “past hour”)
- Stream LLM replies token-by-token for smoother chat UX
- Experiment with embeddings or RAG to boost factual accuracy

---

## Big Picture

Sandbox to learn:
- Real-time computer vision + OpenCV
- YOLOv8 edge detection
- ZeroMQ pub/sub for decoupled pipelines
- SQLite event store with retention & disk cleanup
- Fast Go backend for REST + static files
- Clean `App` struct pattern for microservice style
- Multi-cam smart surveillance
- LLM-powered edge chat assistant
