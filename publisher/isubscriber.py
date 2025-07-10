"""
ISubscriber: An abstract base class that defines the Subscriber interface
for our pub/sub system.

Any concrete Subscriber (e.g., ZeroMQSubscriber) must implement this interface.
"""

from abc import ABC, abstractmethod

class ISubscriber(ABC):
    """
    The ISubscriber interface declares the 'subscribe' method that any
    concrete subscriber must implement.
    """

    @abstractmethod
    def subscribe(self):
        """
        Subscribe to a publisher and receive messages.

        This method should contain the main loop that listens for incoming messages.
        """
        pass
