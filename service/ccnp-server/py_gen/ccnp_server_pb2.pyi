from google.protobuf.internal import containers as _containers
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class LEVEL(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    UNKNOWN: _ClassVar[LEVEL]
    PAAS: _ClassVar[LEVEL]
    SAAS: _ClassVar[LEVEL]
UNKNOWN: LEVEL
PAAS: LEVEL
SAAS: LEVEL

class HealthCheckRequest(_message.Message):
    __slots__ = ("service",)
    SERVICE_FIELD_NUMBER: _ClassVar[int]
    service: str
    def __init__(self, service: _Optional[str] = ...) -> None: ...

class HealthCheckResponse(_message.Message):
    __slots__ = ("status",)
    class ServingStatus(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
        __slots__ = ()
        UNKNOWN: _ClassVar[HealthCheckResponse.ServingStatus]
        SERVING: _ClassVar[HealthCheckResponse.ServingStatus]
        NOT_SERVING: _ClassVar[HealthCheckResponse.ServingStatus]
        SERVICE_UNKNOWN: _ClassVar[HealthCheckResponse.ServingStatus]
    UNKNOWN: HealthCheckResponse.ServingStatus
    SERVING: HealthCheckResponse.ServingStatus
    NOT_SERVING: HealthCheckResponse.ServingStatus
    SERVICE_UNKNOWN: HealthCheckResponse.ServingStatus
    STATUS_FIELD_NUMBER: _ClassVar[int]
    status: HealthCheckResponse.ServingStatus
    def __init__(self, status: _Optional[_Union[HealthCheckResponse.ServingStatus, str]] = ...) -> None: ...

class GetReportRequest(_message.Message):
    __slots__ = ("level", "user_data", "nonce")
    LEVEL_FIELD_NUMBER: _ClassVar[int]
    USER_DATA_FIELD_NUMBER: _ClassVar[int]
    NONCE_FIELD_NUMBER: _ClassVar[int]
    level: LEVEL
    user_data: str
    nonce: str
    def __init__(self, level: _Optional[_Union[LEVEL, str]] = ..., user_data: _Optional[str] = ..., nonce: _Optional[str] = ...) -> None: ...

class GetReportResponse(_message.Message):
    __slots__ = ("report",)
    REPORT_FIELD_NUMBER: _ClassVar[int]
    report: bytes
    def __init__(self, report: _Optional[bytes] = ...) -> None: ...

class GetMeasurementRequest(_message.Message):
    __slots__ = ("level", "index")
    LEVEL_FIELD_NUMBER: _ClassVar[int]
    INDEX_FIELD_NUMBER: _ClassVar[int]
    level: LEVEL
    index: int
    def __init__(self, level: _Optional[_Union[LEVEL, str]] = ..., index: _Optional[int] = ...) -> None: ...

class GetMeasurementResponse(_message.Message):
    __slots__ = ("measurement",)
    MEASUREMENT_FIELD_NUMBER: _ClassVar[int]
    measurement: bytes
    def __init__(self, measurement: _Optional[bytes] = ...) -> None: ...

class GetEventlogRequest(_message.Message):
    __slots__ = ("level", "start", "count")
    LEVEL_FIELD_NUMBER: _ClassVar[int]
    START_FIELD_NUMBER: _ClassVar[int]
    COUNT_FIELD_NUMBER: _ClassVar[int]
    level: LEVEL
    start: int
    count: int
    def __init__(self, level: _Optional[_Union[LEVEL, str]] = ..., start: _Optional[int] = ..., count: _Optional[int] = ...) -> None: ...

class TcgDigest(_message.Message):
    __slots__ = ("algo_id", "hash")
    ALGO_ID_FIELD_NUMBER: _ClassVar[int]
    HASH_FIELD_NUMBER: _ClassVar[int]
    algo_id: int
    hash: bytes
    def __init__(self, algo_id: _Optional[int] = ..., hash: _Optional[bytes] = ...) -> None: ...

class TcgEvent(_message.Message):
    __slots__ = ("imr_index", "event_type", "event_size", "event", "digest", "digests")
    IMR_INDEX_FIELD_NUMBER: _ClassVar[int]
    EVENT_TYPE_FIELD_NUMBER: _ClassVar[int]
    EVENT_SIZE_FIELD_NUMBER: _ClassVar[int]
    EVENT_FIELD_NUMBER: _ClassVar[int]
    DIGEST_FIELD_NUMBER: _ClassVar[int]
    DIGESTS_FIELD_NUMBER: _ClassVar[int]
    imr_index: int
    event_type: int
    event_size: int
    event: bytes
    digest: bytes
    digests: _containers.RepeatedCompositeFieldContainer[TcgDigest]
    def __init__(self, imr_index: _Optional[int] = ..., event_type: _Optional[int] = ..., event_size: _Optional[int] = ..., event: _Optional[bytes] = ..., digest: _Optional[bytes] = ..., digests: _Optional[_Iterable[_Union[TcgDigest, _Mapping]]] = ...) -> None: ...

class GetEventlogResponse(_message.Message):
    __slots__ = ("events",)
    EVENTS_FIELD_NUMBER: _ClassVar[int]
    events: _containers.RepeatedCompositeFieldContainer[TcgEvent]
    def __init__(self, events: _Optional[_Iterable[_Union[TcgEvent, _Mapping]]] = ...) -> None: ...
