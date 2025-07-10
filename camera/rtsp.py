import cv2

class RTSPCamera:
    def __init__(self, url, cam_id="rtsp"):
        self.cap = cv2.VideoCapture(url)
        if not self.cap.isOpened():
            raise RuntimeError(f"Could not open RTSP stream: {url}")
        self.id = cam_id

    def get_frame(self):
        ret, frame = self.cap.read()
        if not ret:
            raise RuntimeError("Failed to grab frame from RTSP stream")
        return frame

    def release(self):
        self.cap.release()

    def __str__(self):
        return self.id
