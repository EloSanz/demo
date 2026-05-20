from fastapi import APIRouter, Depends
from starlette import status

from src.postback.schemas import (
    ConversionPostbackDTO,
    ConversionResponseDTO,
    EACEventPayload,
    EACEventResponseDTO,
)
from src.postback.service import ConversionService, AttributionService
from src.postback.dependencies import get_conversion_service, get_attribution_service

postback_router = APIRouter(prefix="/conversions", tags=["Postback"])
eac_router = APIRouter(tags=["EAC"])

@postback_router.post("", status_code=status.HTTP_202_ACCEPTED, response_model=ConversionResponseDTO)
def received_conversion_postback(
    conversion_postback_dto: ConversionPostbackDTO,
    conversion_service: ConversionService = Depends(get_conversion_service),
):
    return conversion_service.process(conversion_postback_dto)

@eac_router.post("/v1/eac/events", status_code=status.HTTP_200_OK, response_model=EACEventResponseDTO)
async def received_event_postback(
    payload: EACEventPayload,
    attribution_service: AttributionService = Depends(get_attribution_service),
):
    return await attribution_service.process_eac_event(payload)

