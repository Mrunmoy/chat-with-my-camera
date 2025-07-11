# 💡 Chat with my camera — Ideas Parking Lot

### ✅ Whats working Now
- [x] Local webcam source -> tested!
- [x] RTSP camera source -> tested with Unifi Protect!
- [x] YOLOv8 detection loop -> tested!
- [x] ZeroMQ publisher for detection events -> working!
- [x] Config loader (YAML) -> flexible webcam/RTSP.
- [x] Multi-camera support -> multiple streams, single grid view.
- [x] Dynamic grid layout -> auto 1x1, 2x2, 3x3, 4x4.
- [x] Robust reconnect logic for webcam & RTSP.
- [x] Per-camera deduplication (labels only) & throttling.
- [x] Offline placeholder overlay -> tested with real unplug/reboot.

---

## Up Next

### 🟢 1. **Logging / Timeline**
- Keep detection events in SQLite.
- Snapshots saved to disk.
- Add retention job to auto-prune old data.

### 🟢 2. **Timeline API**
- Add a REST API to query events:
  - Last 24 hours.
  - Specific camera.
  - Specific labels.

### 🔴 3. **LLM Integration**
- Add local Ollama or other local LLM.
- Feed detection logs to generate daily summaries:
  - “What did the camera see most often?”
  - “When did a person appear?”

### 🔴 4. **Home Assistant Integration**
- Publish smart events to Home Assistant.
- E.g., motion detected -> turn on lights.

---

## Big Picture

This project is my sandbox to learn:
- Real-time computer vision.
- Interface-driven pub/sub.
- Edge-device AI.
- Config-driven design.
- Robust reconnect + dedup logic.
- Multi-cam smart surveillance.
