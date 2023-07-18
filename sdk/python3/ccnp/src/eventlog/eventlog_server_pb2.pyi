from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class CATEGORY(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = []
    TDX_EVENTLOG: _ClassVar[CATEGORY]
    TPM_EVENTLOG: _ClassVar[CATEGORY]

class LEVEL(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = []
    PAAS: _ClassVar[LEVEL]
    SAAS: _ClassVar[LEVEL]
TDX_EVENTLOG: CATEGORY
TPM_EVENTLOG: CATEGORY
PAAS: LEVEL
SAAS: LEVEL

class GetEventlogRequest(_message.Message):
    __slots__ = ["eventlog_level", "eventlog_category", "start_position", "count"]
    EVENTLOG_LEVEL_FIELD_NUMBER: _ClassVar[int]
    EVENTLOG_CATEGORY_FIELD_NUMBER: _ClassVar[int]
    START_POSITION_FIELD_NUMBER: _ClassVar[int]
    COUNT_FIELD_NUMBER: _ClassVar[int]
    eventlog_level: LEVEL
    eventlog_category: CATEGORY
    start_position: int
    count: int
    def __init__(self, eventlog_level: _Optional[_Union[LEVEL, str]] = ..., eventlog_category: _Optional[_Union[CATEGORY, str]] = ..., start_position: _Optional[int] = ..., count: _Optional[int] = ...) -> None: ...

class GetEventlogReply(_message.Message):
    __slots__ = ["eventlog_data_loc"]
    EVENTLOG_DATA_LOC_FIELD_NUMBER: _ClassVar[int]
    eventlog_data_loc: str
    def __init__(self, eventlog_data_loc: _Optional[str] = ...) -> None: ...
