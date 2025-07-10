# 💡 Chat with my camera — Ideas Parking Lot

### Whats working Now
- [x] Local webcam source → tested!
- [x] YOLOv8 detection loop → tested!

### Next Steps
- Add ZeroMQ publisher to send detection results.
- Test RTSP camera input.
- Add config file (YAML or JSON) for camera source & backend.
- Log detections to file or DB (for LLM queries later).

### Future Ideas
- Integrate local LLM (Ollama) to generate human-friendly summaries.
- Build simple CLI to query detection logs.
- Port detection pipeline to C++ with ONNX Runtime or TensorRT.
- Run on Raspberry Pi with TFLite or Edge TPU.
- Add Home Assistant integration (MQTT or webhook).
- Headless deployment mode as a service.
