"""
IPublisher: An abstract base class that defines the Publisher interface
for our pub/sub system.

Any concrete Publisher (e.g., ZeroMQPublisher) must implement this interface.
"""

from abc import ABC, abstractmethod

class IPublisher(ABC):
    """
    The IPublisher interface declares the 'publish' method that any
    concrete publisher must implement.
    """

    @abstractmethod
    def publish(self, data: dict):
        """
        Publish a message to subscribers.

        Args:
            data (dict): The message payload to publish. Should be JSON-serializable.
        """
        pass
