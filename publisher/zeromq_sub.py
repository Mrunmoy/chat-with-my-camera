"""
ZeroMQSubscriber: Concrete implementation of ISubscriber interface
using ZeroMQ SUB sockets.

Features:
- Optional throttling: only process every Nth message
- Optional deduplication: only show when labels change
"""

import zmq
from publisher.isubscriber import ISubscriber

class ZeroMQSubscriber(ISubscriber):
    """
    ZeroMQSubscriber implements ISubscriber using ZeroMQ PUB/SUB.
    """

    def __init__(self, port=5555, throttle_n=0, deduplicate=True):
        """
        Initialize the subscriber.

        Args:
            port (int): TCP port to connect to.
            throttle_n (int): If > 0, only process every Nth message.
                              0 = no throttling.
            deduplicate (bool): If True, only print when labels change.
        """
        context = zmq.Context()
        self.socket = context.socket(zmq.SUB)
        self.socket.connect(f"tcp://localhost:{port}")
        self.socket.setsockopt_string(zmq.SUBSCRIBE, "")

        self.throttle_n = throttle_n
        self.deduplicate = deduplicate

        print(f"[ZeroMQSubscriber] Subscribed to tcp://localhost:{port} | Throttle: {throttle_n} | Deduplicate: {deduplicate}")

    def subscribe(self):
        """
        Listen for incoming JSON messages and handle throttling/deduplication.
        """
        counter = 0
        last_labels = None

        while True:
            msg = self.socket.recv_json()
            counter += 1

            # Throttle if enabled
            if self.throttle_n > 0 and counter % self.throttle_n != 0:
                continue

            labels = tuple(sorted(msg.get("labels", [])))

            # Deduplicate if enabled
            if self.deduplicate and labels == last_labels:
                continue

            timestamp = msg.get("timestamp")
            boxes = msg.get("boxes", [])

            print(f"[ZeroMQSubscriber] Time: {timestamp} | Labels: {labels} | Boxes: {len(boxes)}")

            last_labels = labels
