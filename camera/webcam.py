import cv2

class Webcam:
    def __init__(self, index=0, cam_id="webcam"):
        self.cap = cv2.VideoCapture(index)
        if not self.cap.isOpened():
            raise RuntimeError(f"Could not open webcam index {index}")
        self.id = cam_id  # <--- store ID

    def get_frame(self):
        ret, frame = self.cap.read()
        if not ret:
            raise RuntimeError("Failed to grab frame from webcam")
        return frame

    def release(self):
        self.cap.release()

    def __str__(self):
        return self.id  # For window name etc.
