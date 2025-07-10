"""
Main entry point for chat-with-my-camera.

This script:
- Captures frames from the webcam
- Runs YOLOv8 object detection
- Publishes detection results to ZeroMQ
- Displays annotated frames in a window
"""

from camera.webcam import Webcam
from camera.rtsp import RTSPCamera

from detection.yolov8 import YOLODetector
from publisher.zeromq_pub import ZeroMQPublisher
from utils.config_loader import load_config

import cv2
import time
import base64

def create_camera(cam_cfg):
    cam_id = cam_cfg.get("id", "unknown")
    if cam_cfg["type"] == "webcam":
        return Webcam(index=cam_cfg["index"], cam_id=cam_id)
    elif cam_cfg["type"] == "rtsp":
        return RTSPCamera(url=cam_cfg["url"], cam_id=cam_id)
    else:
        raise ValueError(f"Unknown camera type: {cam_cfg['type']}")


def main():
    config = load_config(filename="config/config.yaml")

    cameras = [create_camera(cam_cfg) for cam_cfg in config["cameras"]]

    # Initialize YOLO detector (uses yolov8n.pt by default)
    detector = YOLODetector()

    # Initialize ZeroMQ publisher on port 5555
    publisher = ZeroMQPublisher(port=config["publisher"]["port"])

    print("[Main] Running with config:", config)

    while True:
        for cam in cameras:
            frame = cam.get_frame()
            results = detector.process(frame)

            annotated_frame = frame.copy()

            for result in results:
                boxes = result.boxes.xyxy.cpu().numpy().tolist() if result.boxes else []
                labels = [result.names[i] for i in result.boxes.cls.cpu().numpy().astype(int)] if result.boxes else []

                # Encode frame with drawn boxes
                annotated_frame = result.plot()
                ret, buffer = cv2.imencode('.jpg', annotated_frame)
                jpg_as_text = base64.b64encode(buffer).decode('utf-8')

                # Create event payload
                event = {
                    "timestamp": time.time(),
                    "camera_id": cam.id,
                    "boxes": boxes,
                    "labels": labels,
                    # "snapshot": jpg_as_text
                }

                # Publish detection event to ZeroMQ
                publisher.publish(event)

                # Draw boxes and labels on the frame for display
                annotated_frame = result.plot()

            # Show each camera in its own window
            window_name = f"Feed - {cam.id}"
            cv2.imshow(window_name, annotated_frame)

        # Press 'q' to exit
        if cv2.waitKey(1) & 0xFF == ord('q'):
            print("[Main] Quitting detection loop.")
            break

    # Release resources
    for cam in cameras:
        cam.release()
    cv2.destroyAllWindows()

if __name__ == "__main__":
    main()

