# 💡 Chat with my camera — Ideas Parking Lot

### ✅ What’s working now
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
- [x] **Backend refactored to use `App` struct**
  - Clean separation of DB + config.
  - All handlers and jobs as `App` methods.
  - Zero globals → thread-safe, portable, testable.
- [x] ✅ Retention job: deletes old DB rows AND snapshot files → fully tested!
---

## Up Next

### 🟢 Mini HTML / React Timeline UI
- Small test dashboard to render timeline grid with thumbnails.
- Use `/timeline` + `<img src="snapshot_url">`.

### 🔴 LLM Integration
- Feed logs to Ollama or other local LLM.
- Daily or weekly summaries: “What did my cameras see the most?”

### 🔴 Home Assistant Integration
- Publish smart events: motion/person detected → turn on lights, send notifications.

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
