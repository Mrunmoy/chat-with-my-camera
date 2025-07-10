# 💡 Chat with my camera — Ideas Parking Lot

### Whats working Now
- [x] Local webcam source → tested!
- [x] YOLOv8 detection loop → tested!
- [x] Add ZeroMQ publisher to send detection results.
- [x] Add `camera/rtsp.py` to handle IP cameras or streams.
- [x] Multiview - mutiple camera streams supported.

## Up Next

### [x] 1. **RTSP Source Support**
- Add `camera/rtsp.py` to handle IP cameras or streams.
- Swap between `Webcam` and `RTSPCamera` via config.

### [x] 2. **Config Loader**
- Add a YAML or JSON config file:
  - Camera source: webcam / RTSP
  - YOLO model path: yolov8n.pt, yolov8s.pt, etc.
  - Publisher port
  - Throttling level
- Load this in `main.py` so the pipeline is fully flexible.

### [] 3. **Logging**
- Add a simple file logger subscriber:
  - Save events to JSONL or SQLite.
  - Optionally store snapshot images with filenames.

### [] 4. **LLM Integration**
- Add local Ollama or other local LLM.
- Feed detection logs to LLM to generate daily summaries:
  - “What did the camera see most often?”
  - “When did a person appear?”

### [] 5. **HA Integration**
- Integrate with Home Asistant and pubnlish events
---
