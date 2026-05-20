from fastapi import Depends

from src.infrastructure.queue.sqs import SQSQueue
from src.postback.ports.queue_port import QueuePort
# from src.infrastructure.queue.mock_queue import MockQueue
# from src.infrastructure.queue.gcp_pubsub import GCPPubSub
from src.postback.ports.transaction_repository_port import TransactionRepositoryPort
from src.infrastructure.repository.in_memory_transaction_repository import InMemoryTransactionRepository
from src.postback.service import ConversionService, AttributionService


def get_queue() -> QueuePort:
    # return MockQueue()
    # return GCPPubSub(project_id="app_stack-gcp", topic_name="ios-postbacks-topic")
    return SQSQueue(topic_name="topic", queue_url="queue_url")

_transaction_repo = InMemoryTransactionRepository()

def get_transaction_repo() -> TransactionRepositoryPort:
    return _transaction_repo

def get_conversion_service(
    queue: QueuePort = Depends(get_queue),
    transaction_repo: TransactionRepositoryPort = Depends(get_transaction_repo)
) -> ConversionService:
    return ConversionService(queue=queue, transaction_repo=transaction_repo)

def get_attribution_service() -> AttributionService:
    return AttributionService()