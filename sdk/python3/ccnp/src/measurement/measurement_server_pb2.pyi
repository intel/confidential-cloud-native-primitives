from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class TYPE(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = []
    PAAS: _ClassVar[TYPE]
    SAAS: _ClassVar[TYPE]

class CATEGORY(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = []
    TEE_REPORT: _ClassVar[CATEGORY]
    TPM: _ClassVar[CATEGORY]
    TDX_RTMR: _ClassVar[CATEGORY]
PAAS: TYPE
SAAS: TYPE
TEE_REPORT: CATEGORY
TPM: CATEGORY
TDX_RTMR: CATEGORY

class GetMeasurementRequest(_message.Message):
    __slots__ = ["measurement_type", "measurement_category", "report_data", "register_index"]
    MEASUREMENT_TYPE_FIELD_NUMBER: _ClassVar[int]
    MEASUREMENT_CATEGORY_FIELD_NUMBER: _ClassVar[int]
    REPORT_DATA_FIELD_NUMBER: _ClassVar[int]
    REGISTER_INDEX_FIELD_NUMBER: _ClassVar[int]
    measurement_type: TYPE
    measurement_category: CATEGORY
    report_data: str
    register_index: int
    def __init__(self, measurement_type: _Optional[_Union[TYPE, str]] = ..., measurement_category: _Optional[_Union[CATEGORY, str]] = ..., report_data: _Optional[str] = ..., register_index: _Optional[int] = ...) -> None: ...

class GetMeasurementReply(_message.Message):
    __slots__ = ["measurement"]
    MEASUREMENT_FIELD_NUMBER: _ClassVar[int]
    measurement: str
    def __init__(self, measurement: _Optional[str] = ...) -> None: ...
