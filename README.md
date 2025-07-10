# 📸 chat-with-my-camera
![Project Status](https://img.shields.io/badge/status-under--construction-yellow) 


![Under Construction Cat](https://media.giphy.com/media/VbnUQpnihPSIgIXuZv/giphy.gif)


🚧 **This project is under progress!** 🚧


A modular real-time object detection pipeline for Linux boxes, Raspberry Pi, or Jetson — built to run YOLOv8 and talk to you about what it sees. 

## Features
- Modular camera source (webcam now, RTSP next!)
- YOLOv8 inference with PyTorch/Ultralytics
- Real-time bounding boxes drawn on live video
- ZeroMQ publisher planned for detection events
- Future: integrate local LLM (Ollama) to chat with your detection logs

## How to run
```bash
# Create and activate a venv
python -m venv venv
source venv/bin/activate

# Install dependencies
pip install opencv-python ultralytics

# Run it!
python main.py
```

## 🗺️ Roadmap
See [IDEAS.md](IDEAS.md) for current and future plans.

---

Built with ❤️ by Mrunmoy — keep watching!