import structlog
from fastapi import FastAPI, Request
from starlette import status
from starlette.responses import JSONResponse

from src.postback.exceptions import AndroidMissingIDError, CampaignMissingIDError, DuplicateTransactionError

logger = structlog.get_logger()

def setup_postback_exception_handlers(app: FastAPI):
    @app.exception_handler(DuplicateTransactionError)
    async def duplicate_transaction_exception_handler(_: Request, exc: DuplicateTransactionError):
        logger.warning("duplicate_transaction_attempt", transaction_id=exc.transaction_id)
        return JSONResponse(
            status_code=423,
            content={
                "error": "Locked",
                "message": f"Transaction {exc.transaction_id} is already processed or being processed."
            }
        )

    @app.exception_handler(AndroidMissingIDError)
    async def android_missing_id_exception_handler(_: Request, exc: AndroidMissingIDError):
        logger.warning("android_missing_device_id", transaction_id=exc.transaction_id)

        return JSONResponse(
            status_code=status.HTTP_400_BAD_REQUEST,
            content={
                "error": "Bad Request",
                "message": "device_id is mandatory for Android deterministic attribution."
            }
        )

    @app.exception_handler(CampaignMissingIDError)
    async def campaign_missing_id_exception_handler(_: Request, exc: CampaignMissingIDError):
        logger.warning("campaign_missing_id", transaction_id=exc.transaction_id)

        return JSONResponse(
            status_code=status.HTTP_400_BAD_REQUEST,
            content={
                "error": "Bad Request",
                "message": "campaign id must be provided."
            }
        )
