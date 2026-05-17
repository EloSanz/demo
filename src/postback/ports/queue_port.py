from typing import Protocol, Any

class QueuePort(Protocol):
    def enqueue(self, message: Any) -> None: ...
