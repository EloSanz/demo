import structlog
from src.postback.ports.transaction_repository_port import TransactionRepositoryPort
from src.postback.exceptions import DuplicateTransactionError

logger = structlog.get_logger()

class InMemoryTransactionRepository(TransactionRepositoryPort):
    def __init__(self):
        self._processed_transactions = set()

    def lock(self, transaction_id: str) -> None:
        logger.info("locking_transaction", transaction_id=transaction_id)
        if transaction_id in self._processed_transactions:
            raise DuplicateTransactionError(transaction_id=transaction_id)
        self._processed_transactions.add(transaction_id)

    def unlock(self, transaction_id: str) -> None:
        self._processed_transactions.discard(transaction_id)
