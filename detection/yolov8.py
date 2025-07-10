from ultralytics import YOLO

class YOLODetector:
    def __init__(self, model_path="yolov8n.pt"):
        # You can change to yolov8s.pt or any other if you want
        self.model = YOLO(model_path)

    def process(self, frame):
        # Runs detection on a single frame
        results = self.model.predict(
            source=frame,
            save=False,
            conf=0.5,
            verbose=False
        )
        return results
