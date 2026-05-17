from enum import Enum
from typing import Optional

from pydantic import BaseModel, Field, model_validator
from src.postback.exceptions import AndroidMissingIDError, CampaignMissingIDError


class OSEnum(str, Enum):
    android = "android"
    ios = "ios"

class AttributionTypeEnum(str, Enum):
    probabilistic = "probabilistic"
    deterministic = "deterministic"

class ConversionPostbackDTO(BaseModel):
    transaction_id: str = Field(min_length=1)
    campaign_id: Optional[str] = None
    device_id: Optional[str] = None
    os: OSEnum
    revenue: float = Field(gt=0.0)

    @model_validator(mode='after')
    def validate_conditional_os_logic(self) -> 'ConversionPostbackDTO':
        if self.os == OSEnum.android:
            if not self.device_id or self.device_id.strip() == '':
                raise AndroidMissingIDError(transaction_id=self.transaction_id)
        return self

    @model_validator(mode='after')
    def validate_campaign_id(self) -> 'ConversionPostbackDTO':
        if not self.campaign_id:
            raise CampaignMissingIDError(transaction_id=self.transaction_id)

        return self

class ConversionResponseDTO(BaseModel):
    status: str
    transaction_id: str
    attribution_type: AttributionTypeEnum
    internal_status: str
