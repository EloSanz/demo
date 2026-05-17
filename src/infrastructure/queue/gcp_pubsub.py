from typing import Any
import structlog
from src.postback.ports.queue_port import QueuePort

logger = structlog.get_logger()

class GCPPubSub(QueuePort):
    def __init__(self, project_id: str, topic_name: str):
        self.project_id = project_id
        self.topic_name = topic_name

    def enqueue(self, message: Any) -> None:
        logger.info(
            "publish_gcp_pubsub_event",
            project_id=self.project_id,
            topic_name=self.topic_name,
            endpoint=f"projects/{self.project_id}/topics/{self.topic_name}",
            payload=message
        )
