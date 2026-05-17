from typing import Any
import structlog
from src.postback.ports.queue_port import QueuePort

logger = structlog.get_logger()

class MockQueue(QueuePort):
    def enqueue(self, message: Any) -> None:
        logger.info("enqueue_mock_message", message=message)
