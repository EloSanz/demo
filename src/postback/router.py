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
    conversion_domain = conversion_postback_dto.to_domain()
    result = conversion_service.process(conversion_domain)
    return ConversionResponseDTO.from_domain(result)

@eac_router.post("/v1/eac/events", status_code=status.HTTP_200_OK, response_model=EACEventResponseDTO)
async def received_event_postback(
    payload: EACEventPayload,
    attribution_service: AttributionService = Depends(get_attribution_service),
):
    event_domain = payload.to_domain()
    result = await attribution_service.process_eac_event(event_domain)
    return EACEventResponseDTO.from_domain(result)


