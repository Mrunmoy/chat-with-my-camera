from camera.webcam import Webcam
from detection.yolov8 import YOLODetector
import cv2

def main():
    cam = Webcam()
    detector = YOLODetector()

    while True:
        frame = cam.get_frame()
        results = detector.process(frame)

        # Draw boxes on frame
        for result in results:
            annotated_frame = result.plot()

            cv2.imshow("YOLO Detection", annotated_frame)

        # Quit on 'q'
        if cv2.waitKey(1) & 0xFF == ord('q'):
            break

    cam.release()
    cv2.destroyAllWindows()

if __name__ == "__main__":
    main()
