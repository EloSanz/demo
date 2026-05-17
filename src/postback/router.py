from fastapi import APIRouter, Depends
from starlette import status

from src.postback.schemas import ConversionPostbackDTO, ConversionResponseDTO
from src.postback.service import ConversionService
from src.postback.dependencies import get_conversion_service

postback_router = APIRouter(prefix="/conversions", tags=["Postback"])

@postback_router.post("", status_code=status.HTTP_202_ACCEPTED, response_model=ConversionResponseDTO)
def received_conversion_postback(
    conversion_postback_dto: ConversionPostbackDTO,
    conversion_service: ConversionService = Depends(get_conversion_service),
):
    return conversion_service.process(conversion_postback_dto)
