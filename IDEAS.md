# 💡 Chat with my camera — Ideas Parking Lot

### Whats working Now
- [x] Local webcam source → tested!
- [x] YOLOv8 detection loop → tested!
- [x] Add ZeroMQ publisher to send detection results.
- [x] Add `camera/rtsp.py` to handle IP cameras or streams.

## Up Next

### [x] 1. **RTSP Source Support**
- Add `camera/rtsp.py` to handle IP cameras or streams.
- Swap between `Webcam` and `RTSPCamera` via config.

### 2. **Config Loader**
- Add a YAML or JSON config file:
  - Camera source: webcam / RTSP
  - YOLO model path: yolov8n.pt, yolov8s.pt, etc.
  - Publisher port
  - Throttling level
- Load this in `main.py` so the pipeline is fully flexible.

### 3. **Logging**
- Add a simple file logger subscriber:
  - Save events to JSONL or SQLite.
  - Optionally store snapshot images with filenames.

### 4. **Headless Mode**
- Add a CLI flag to run without GUI (`cv2.imshow()`).
- Useful for Raspberry Pi or cloud deployment.

### 5. **LLM Integration**
- Add local Ollama or other local LLM.
- Feed detection logs to LLM to generate daily summaries:
  - “What did the camera see most often?”
  - “When did a person appear?”

---

## 🔭 Future Ideas

- Port the detector or camera source to C++ for performance.
- Package pub/sub as a pip-installable library.
- Add MQTT backend as an alternative to ZeroMQ.
- Build a tiny web dashboard to show live detections.
- Maybe add Home Assistant integration for smart home automations.

---
