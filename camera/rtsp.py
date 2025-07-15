import cv2
import time

class RTSPCamera:
    def __init__(self, url, cam_id="rtsp"):
        self.url = url
        self.id = cam_id
        self.cap = cv2.VideoCapture(url)
        self.is_online = self.cap.isOpened()

        if not self.is_online:
            print(f"[RTSPCamera] {self.id} - Could not open RTSP stream at init: {url}")

    def get_frame(self):
        # If cap is closed, try to reconnect
        if not self.cap.isOpened():
            print(f"[RTSPCamera] {self.id} - Stream closed, attempting reconnect...")
            self.cap.release()
            time.sleep(2)  # Small delay so we don’t hammer the server
            self.cap = cv2.VideoCapture(self.url)
            self.is_online = self.cap.isOpened()
            if not self.is_online:
                print(f"[RTSPCamera] {self.id} - Still offline after reconnect attempt.")
                return None

        ret, frame = self.cap.read()
        if not ret or frame is None:
            print(f"[RTSPCamera] {self.id} - Frame grab failed, releasing for reconnect.")
            self.cap.release()
            self.is_online = False
            return None

        self.is_online = True
        return frame

    def release(self):
        self.cap.release()

    def __str__(self):
        return self.id
