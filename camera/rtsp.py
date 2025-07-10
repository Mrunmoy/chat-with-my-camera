import cv2

class RTSPCamera:
    def __init__(self, url):
        self.cap = cv2.VideoCapture(url)
        if not self.cap.isOpened():
            raise RuntimeError(f"Could not open RTSP stream: {url}")

    def get_frame(self):
        ret, frame = self.cap.read()
        if not ret:
            raise RuntimeError("Failed to grab frame from RTSP stream")
        return frame

    def release(self):
        self.cap.release()
