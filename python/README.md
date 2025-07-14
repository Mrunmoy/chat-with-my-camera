# MQ Pub Sub & Camera Handling

## Overview

This folder covers the **Python parts** for `chat-with-my-camera`:

- Camera capture (OpenCV)
- YOLOv8 detection using Ultralytics
- ZeroMQ Publisher to send detections to the Go backend

Everything runs locally on Linux, Pi, or Jetson.

---

## How It Works

- **Camera Capture**

- Uses OpenCV (`cv2.VideoCapture`).
- Supports multiple webcam or RTSP sources.
- Defined in `config.yaml`.

- **YOLOv8 Detection**

- Uses Ultralytics PyTorch model (`yolov8n.pt` or similar).
- Captures frames, resizes, runs inference.
- Saves snapshot images for detections.
- Publishes detection data over ZeroMQ.

- **ZeroMQ Publisher**

- Publishes JSON events: timestamp, camera\_id, labels, bounding boxes, snapshot filename.
- Other apps (Go backend) subscribe and process the events.

---

## Directory Layout

```
python/
├── camera/           # Camera modules (OpenCV wrappers)
├── detection/        # YOLOv8 detection logic
├── publisher/        # ZeroMQPublisher, IPublisher interface
├── utils/            # Shared helpers
├── main.py           # Entry point: camera + detection + pub/sub
├── zmq_subscriber_example.py  # Test subscriber
├── yolov8n.pt        # Model weights
├── config/           # Config YAML with camera sources
```

---

## Example Event JSON

A single detection looks like:

```json
{
  "timestamp": 1720518700.123,
  "camera_id": "garage_webcam",
  "labels": ["person", "car"],
  "boxes": [[x1, y1, x2, y2]],
  "snapshot_file": "./snapshots/garage_webcam_1720518700.jpg"
}
```

---

## How to Run

```bash
python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt

python main.py
```

Optionally run the test subscriber:

```bash
python zmq_subscriber_example.py
```

---

