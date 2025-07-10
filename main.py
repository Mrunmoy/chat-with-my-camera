from camera.webcam import Webcam
import cv2

def main():
    cam = Webcam()

    while True:
        frame = cam.get_frame()
        cv2.imshow("Webcam Test", frame)

        # Exit on 'q'
        if cv2.waitKey(1) & 0xFF == ord('q'):
            break

    cam.release()
    cv2.destroyAllWindows()

if __name__ == "__main__":
    main()
