import cv2
import time

class Webcam:
    def __init__(self, index=0, cam_id="webcam"):
        self.index = index
        self.id = cam_id
        self.cap = cv2.VideoCapture(index)
        self.is_online = self.cap.isOpened()

        if not self.is_online:
            print(f"[Webcam] {self.id} - Could not open webcam index {index}")

    def get_frame(self):
        if not self.cap.isOpened():
            print(f"[Webcam] {self.id} - Stream closed, attempting reconnect...")
            self.cap.release()
            time.sleep(2)
            self.cap = cv2.VideoCapture(self.index)
            self.is_online = self.cap.isOpened()
            if not self.is_online:
                print(f"[Webcam] {self.id} - Still offline after reconnect attempt.")
                return None

        ret, frame = self.cap.read()
        if not ret or frame is None:
            print(f"[Webcam] {self.id} - Frame grab failed, releasing for reconnect.")
            self.cap.release()
            self.is_online = False
            return None

        self.is_online = True
        return frame

    def release(self):
        self.cap.release()

    def __str__(self):
        return self.id
