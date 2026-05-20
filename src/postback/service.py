from src.postback.handlers import logger
from src.postback.ports.queue_port import QueuePort
from src.postback.ports.transaction_repository_port import TransactionRepositoryPort
from src.postback.domain.models import (
    Conversion,
    AttributionResult,
    EACEvent,
    EACEventResult,
    OSDomainEnum,
    AttributionDomainEnum,
)


class ConversionService:
    def __init__(self, queue: QueuePort, transaction_repo: TransactionRepositoryPort):
        self.queue = queue
        self.transaction_repo = transaction_repo

    def process(self, postback: Conversion) -> AttributionResult:
        self.transaction_repo.lock(postback.transaction_id)
        
        try:
            #sleep(3)
            if postback.os == OSDomainEnum.android:
                attribution_type = AttributionDomainEnum.deterministic
                internal_status = "attributed"
            else:
                if postback.device_id is None:
                    attribution_type = AttributionDomainEnum.probabilistic
                    internal_status = "queued_for_probabilistic_attribution"

                    self._enqueue_ios_new(postback.transaction_id, attribution_type.name)
                else:
                    attribution_type = AttributionDomainEnum.deterministic
                    internal_status = "attributed"

            return AttributionResult(
                status="success",
                transaction_id=postback.transaction_id,
                attribution_type=attribution_type,
                internal_status=internal_status
            )
        finally:
            self.transaction_repo.unlock(postback.transaction_id)

    def _enqueue_ios_new(self, transaction_id: str, attribution_type: str):
        self.queue.enqueue(
            {
                "action": "ios_attribution",
                "transaction_id": transaction_id,
                "type": attribution_type
            }
        )


class AttributionService:
    async def process_eac_event(self, event: EACEvent) -> EACEventResult:
        logger.info("process_eac_event", event_id=event.event_id)

        if event.device_id and event.device_id.strip() != "":
            attribution_method = "deterministic"
        else:
            attribution_method = "probabilistic"

        return EACEventResult(
            event_id=event.event_id,
            status="accepted",
            attribution_method=attribution_method
        )