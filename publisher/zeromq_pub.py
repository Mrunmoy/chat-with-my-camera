"""
ZeroMQPublisher: Concrete implementation of the IPublisher interface
using ZeroMQ PUB sockets.

This publisher broadcasts JSON messages to any connected subscribers.
"""

import zmq
from publisher.ipublisher import IPublisher

class ZeroMQPublisher(IPublisher):
    """
    ZeroMQPublisher implements the IPublisher interface using ZeroMQ.
    """

    def __init__(self, port=5555):
        """
        Initialize the ZeroMQ publisher socket.

        Args:
            port (int): The TCP port to bind the PUB socket on.
        """
        context = zmq.Context()
        self.socket = context.socket(zmq.PUB)
        self.socket.bind(f"tcp://*:{port}")
        print(f"[ZeroMQPublisher] Publishing on tcp://*:{port}")

    def publish(self, data: dict):
        """
        Publish a JSON-serializable message to all subscribers.

        Args:
            data (dict): The message payload.
        """
        self.socket.send_json(data)
