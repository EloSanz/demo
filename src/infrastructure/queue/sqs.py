from typing import Any
import structlog
from src.postback.ports.queue_port import QueuePort

logger = structlog.get_logger()


class SQSQueue(QueuePort):
    def __init__(self, project_id: str = "app_stack", topic_name: str = "", queue_url: str = "") -> None:
        self.project_id = project_id
        self.topic_name = topic_name
        self.queue_url = queue_url

    def enqueue(self, message: Any) -> None:
        # mock enqueue
        logger.info(
            "enqueue_sqs_message",
            message=message,
            project_id=self.project_id,
            topic_name=self.topic_name,
            queue_url=self.queue_url
        )