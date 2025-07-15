"""
Main entry point for chat-with-my-camera.

This script:
- Captures frames from webcams/RTSP
- Runs YOLOv8 object detection
- Publishes detection results to ZeroMQ
- Displays annotated frames in a multiview grid
"""

from camera.webcam import Webcam
from camera.rtsp import RTSPCamera

from detection.yolov8 import YOLODetector
from publisher.zeromq_pub import ZeroMQPublisher
from utils.config_loader import load_config

import numpy as np
import cv2
import math
import time
import base64


def make_multi_view_grid(frames, cell_size=(640, 480)):
    """
    Arrange frames in a square-ish grid.

    Args:
        frames (List[np.ndarray]): List of images.
        cell_size (Tuple[int,int]): Desired (width, height) for each cell.

    Returns:
        np.ndarray: The combined grid image.
    """
    num_frames = len(frames)
    grid_cols = math.ceil(math.sqrt(num_frames))
    grid_rows = math.ceil(num_frames / grid_cols)

    # Pad frames list with black images if needed
    total_cells = grid_rows * grid_cols
    if num_frames < total_cells:
        black_frame = np.zeros((cell_size[1], cell_size[0], 3), dtype=np.uint8)
        frames += [black_frame] * (total_cells - num_frames)

    # Resize all frames
    resized_frames = [
        cv2.resize(frame, cell_size) for frame in frames
    ]

    # Build rows
    rows = []
    for i in range(0, total_cells, grid_cols):
        row = np.hstack(resized_frames[i:i + grid_cols])
        rows.append(row)

    # Stack rows vertically
    grid_image = np.vstack(rows)
    return grid_image


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
        frames = []
        for cam in cameras:
            frame = cam.get_frame()

            if frame is None:
                # Camera offline, show placeholder
                frame = np.zeros((480, 640, 3), dtype=np.uint8)
                cv2.putText(frame, f"Camera Offline: {cam.id}",
                            (50, 240),
                            cv2.FONT_HERSHEY_SIMPLEX,
                            1, (0, 0, 255), 2)
                frames.append(frame)
                continue  # Skip detection & publishing for offline feed

            # If we have a valid frame, run detection
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
                    "snapshot": jpg_as_text
                }

                # Publish detection event to ZeroMQ
                publisher.publish(event)

            frames.append(annotated_frame)

        # Show multiview
        multi_view = make_multi_view_grid(frames, cell_size=(640, 480))
        cv2.imshow("Multi-Camera View", multi_view)

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
