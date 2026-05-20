from enum import Enum
from typing import Optional

from pydantic import BaseModel, Field, model_validator
from src.postback.exceptions import (
    AndroidMissingIDError,
    CampaignMissingIDError,
    AndroidDeviceIDMissingError,
    iOSAttributionDataMissingError,
)
from src.postback.domain.models import (
    Conversion,
    AttributionResult,
    EACEvent,
    EACEventResult,
    OSDomainEnum,
    AttributionDomainEnum,
)


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

    def to_domain(self) -> Conversion:
        return Conversion(
            transaction_id=self.transaction_id,
            campaign_id=self.campaign_id or "",
            device_id=self.device_id,
            os=OSDomainEnum(self.os.value),
            revenue=self.revenue
        )

class ConversionResponseDTO(BaseModel):
    status: str
    transaction_id: str
    attribution_type: AttributionTypeEnum
    internal_status: str

    @classmethod
    def from_domain(cls, result: AttributionResult) -> "ConversionResponseDTO":
        return cls(
            status=result.status,
            transaction_id=result.transaction_id,
            attribution_type=AttributionTypeEnum(result.attribution_type.value),
            internal_status=result.internal_status
        )

class EACEventPayload(BaseModel):
    event_id: str = Field(min_length=1)
    event_type: str = Field(min_length=1)
    os: OSEnum
    device_id: Optional[str] = None
    ip_address: Optional[str] = None
    user_agent: Optional[str] = None
    revenue: float = Field(gt=0.0)

    @model_validator(mode='after')
    def validate_eac_attribution_logic(self) -> 'EACEventPayload':
        if self.os == OSEnum.android:
            if not self.device_id or self.device_id.strip() == '':
                raise AndroidDeviceIDMissingError()
        elif self.os == OSEnum.ios:
            if not self.device_id or self.device_id.strip() == '':
                if (not self.ip_address or self.ip_address.strip() == '') or \
                   (not self.user_agent or self.user_agent.strip() == ''):
                    raise iOSAttributionDataMissingError()
        return self

    def to_domain(self) -> EACEvent:
        return EACEvent(
            event_id=self.event_id,
            event_type=self.event_type,
            os=OSDomainEnum(self.os.value),
            device_id=self.device_id,
            ip_address=self.ip_address,
            user_agent=self.user_agent,
            revenue=self.revenue
        )

class EACEventResponseDTO(BaseModel):
    event_id: str
    status: str
    attribution_method: str

    @classmethod
    def from_domain(cls, result: EACEventResult) -> "EACEventResponseDTO":
        return cls(
            event_id=result.event_id,
            status=result.status,
            attribution_method=result.attribution_method
        )