from dataclasses import dataclass
from enum import Enum

class OSDomainEnum(str, Enum):
    android = "android"
    ios = "ios"

class AttributionDomainEnum(str, Enum):
    probabilistic = "probabilistic"
    deterministic = "deterministic"

@dataclass
class Conversion:
    transaction_id: str
    campaign_id: str
    device_id: str | None
    os: OSDomainEnum
    revenue: float

@dataclass
class AttributionResult:
    status: str
    transaction_id: str
    attribution_type: AttributionDomainEnum
    internal_status: str

@dataclass
class EACEvent:
    event_id: str
    event_type: str
    os: OSDomainEnum
    device_id: str | None
    ip_address: str | None
    user_agent: str | None
    revenue: float

@dataclass
class EACEventResult:
    event_id: str
    status: str
    attribution_method: str
