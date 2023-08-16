# pylint: disable=duplicate-code
# Copyright (c) 2023, Intel Corporation. All rights reserved.<BR>
# SPDX-License-Identifier: Apache-2.0
"""
This package provides the definitions and helper class for the event log of TD
or TPM.

Reference:
1. https://github.com/tpm2-software/tpm2-tcti-uefi/blob/master/src/tcg2-protocol.h
"""

import logging
import os
import json
from typing import List
import grpc

# pylint: disable=E1101
from ccnp.eventlog import eventlog_server_pb2
from ccnp.eventlog import eventlog_server_pb2_grpc

LOG = logging.getLogger(__name__)
TIMEOUT = 5

class CCAlgorithms:
    """
    Algorithms class for confidential computing.

    The definitions are aligning with TCG specification - "TCG Algorithm Registry"
    at https://trustedcomputinggroup.org/wp-content/uploads/TCGAlgorithmRegistry_Rev01.15.pdf
    """

    ALG_SHA1 = 0xA
    ALG_SHA256 = 0xB
    ALG_SHA384 = 0xC
    ALG_SHA512 = 0xD
    ALG_SM3_256 = 0xE
    ALG_INVALID = 0xffffffff

    _algo_dict = None

    _digest_size = {
        ALG_SHA1: 20,
        ALG_SHA256: 32,
        ALG_SHA384: 48,
        ALG_SHA512: 64,
        ALG_SM3_256: 32
    }

    _block_size = {
        ALG_SHA1: 64,
        ALG_SHA256: 64,
        ALG_SHA384: 128,
        ALG_SHA512: 128,
        ALG_SM3_256: 64
    }

    def __init__(self):
        self._algo_id = CCAlgorithms.ALG_INVALID

    @property
    def algo_id(self):
        """
        Property of algorithms ID
        """
        return self._algo_id

    @property
    def digest_size(self):
        """
        Property of digest size
        """
        return self._digest_size[self.algo_id]

    @property
    def block_size(self):
        """
        Property of block size
        """
        return self._block_size[self.algo_id]

    @algo_id.setter
    def algo_id(self, value):
        """
        Setter for the property algorithms ID
        """
        assert value != CCAlgorithms.ALG_INVALID
        self._algo_id = value

    @property
    def is_valid(self):
        """
        Property of algorithm id valid check
        """
        return self._algo_id != CCAlgorithms.ALG_INVALID

    @classmethod
    def algo_dict(cls):
        """
        Class method to construct the algo dict
        """
        if cls._algo_dict is not None:
            return cls._algo_dict

        # first time initialization
        cls._algo_dict = {}
        for key, value in cls.__dict__.items():
            if key.startswith('ALG'):
                # pylint: disable=E1137
                cls._algo_dict[value] = key
        return cls._algo_dict

    def __str__(self):
        """
        Get string of algorithms name
        """
        CCAlgorithms.algo_dict()
        assert CCAlgorithms._algo_dict is not None
        # pylint: disable=E1135
        assert self.algo_id in CCAlgorithms._algo_dict
        # pylint: disable=E1136
        return CCAlgorithms._algo_dict[self.algo_id]


class CCEventLogType:
    """
    Event log type for Confidential Computing
    """

    # TCG PC Client Specific Implementation Specification for Conventional BIOS
    EV_PREBOOT_CERT = 0x0
    EV_POST_CODE = 0x1
    EV_UNUSED = 0x2
    EV_NO_ACTION = 0x3
    EV_SEPARATOR = 0x4
    EV_ACTION = 0x5
    EV_EVENT_TAG = 0x6
    EV_S_CRTM_CONTENTS = 0x7
    EV_S_CRTM_VERSION = 0x8
    EV_CPU_MICROCODE = 0x9
    EV_PLATFORM_CONFIG_FLAGS = 0xa
    EV_TABLE_OF_DEVICES = 0xb
    EV_COMPACT_HASH = 0xc
    EV_IPL = 0xd
    EV_IPL_PARTITION_DATA = 0xe
    EV_NONHOST_CODE = 0xf
    EV_NONHOST_CONFIG = 0x10
    EV_NONHOST_INFO = 0x11
    EV_OMIT_BOOT_DEVICE_EVENTS = 0x12

    # TCG EFI Platform Specification For TPM Family 1.1 or 1.2
    EV_EFI_EVENT_BASE = 0x80000000
    EV_EFI_VARIABLE_DRIVER_CONFIG = EV_EFI_EVENT_BASE + 0x1
    EV_EFI_VARIABLE_BOOT = EV_EFI_EVENT_BASE + 0x2
    EV_EFI_BOOT_SERVICES_APPLICATION = EV_EFI_EVENT_BASE + 0x3
    EV_EFI_BOOT_SERVICES_DRIVER = EV_EFI_EVENT_BASE + 0x4
    EV_EFI_RUNTIME_SERVICES_DRIVER = EV_EFI_EVENT_BASE + 0x5
    EV_EFI_GPT_EVENT = EV_EFI_EVENT_BASE + 0x6
    EV_EFI_ACTION = EV_EFI_EVENT_BASE + 0x7
    EV_EFI_PLATFORM_FIRMWARE_BLOB = EV_EFI_EVENT_BASE + 0x8
    EV_EFI_HANDOFF_TABLES = EV_EFI_EVENT_BASE + 0x9
    EV_EFI_VARIABLE_AUTHORITY = EV_EFI_EVENT_BASE + 0xe0
    EV_UNKNOWN_A = EV_EFI_EVENT_BASE + 0xa
    EV_UNKNOWN_B = EV_EFI_EVENT_BASE + 0xb
    EV_UNKNOWN_C = EV_EFI_EVENT_BASE + 0xc


    EV_INVALID = 0xffffffff

    _type_dict = None

    @classmethod
    def event_log_dict(cls):
        """
        Class method to construct the event log dict
        """
        if cls._type_dict is not None:
            return cls._type_dict

        # first time initialization
        cls._type_dict = {}
        for key, value in cls.__dict__.items():
            if key.startswith('EV_'):
                # pylint: disable=E1137
                cls._type_dict[value] = key
        return cls._type_dict

    def __init__(self):
        """
        Constructor
        """
        self._log_type = CCEventLogType.EV_INVALID

    @property
    def log_type(self):
        """
        Property for type of event log
        """
        return self._log_type

    @log_type.setter
    def log_type(self, value):
        """
        Set the event log type
        """
        CCEventLogType.event_log_dict()
        assert CCEventLogType._type_dict is not None
        # pylint: disable=E1135
        assert value in CCEventLogType._type_dict
        self._log_type = value

    @classmethod
    def log_type_string(cls, value):
        """
        Get string of eventlog type
        """
        assert CCEventLogType._type_dict is not None
        # pylint: disable=E1135
        assert value in CCEventLogType._type_dict
        # pylint: disable=E1136
        return CCEventLogType._type_dict[value]

class CCEventLogEntry:

    INVALID_MEASURE_REGISTER_INDEX = -1
    INVALID_EVENT_TYPE = -1
    INVALID_ALGORITHMS_ID = -1
    UNKNOWN_EVENT_TYPE_NAME = "UNKNOWN"

    def __init__(self) -> None:
        self._reg_idx: int = CCEventLogEntry.INVALID_MEASURE_REGISTER_INDEX
        self._evt_type: int = CCEventLogType()
        self._evt_type_str: str = CCEventLogEntry.UNKNOWN_EVENT_TYPE_NAME
        self._evt_size: int = -1
        self._alg_id: CCAlgorithms = CCAlgorithms()
        self._event: bytearray = None
        self._digest: bytearray = None

    @property
    def reg_idx(self):
        """
        Property for type register index
        """
        return self._reg_idx

    @reg_idx.setter
    def reg_idx(self, value):
        """
        Setter for the property register index
        """
        assert value != CCEventLogEntry.INVALID_MEASURE_REGISTER_INDEX
        self._reg_idx = value

    @property
    def evt_type(self):
        """
        Property for type event type
        """
        return self._evt_type

    @evt_type.setter
    def evt_type(self, value):
        """
        Setter for the property event type
        """
        assert value != CCEventLogEntry.INVALID_EVENT_TYPE
        self._evt_type = value

    @property
    def evt_type_str(self):
        """
        Property for type event type string
        """
        return self._evt_type_str

    @evt_type_str.setter
    def evt_type_str(self, value):
        """
        Setter for the property event type string
        """
        assert value != CCEventLogEntry.UNKNOWN_EVENT_TYPE_NAME
        self._evt_type_str = value

    @property
    def evt_size(self):
        """
        Property for type event size
        """
        return self._evt_size

    @evt_size.setter
    def evt_size(self, value):
        """
        Setter for the property event size
        """
        assert value > 0
        self._evt_size = value

    @property
    def alg_id(self):
        """
        Property for type algorithm id
        """
        return self._alg_id

    @alg_id.setter
    def alg_id(self, value):
        """
        Setter for the property algorithms ID
        """
        assert value != CCEventLogEntry.INVALID_ALGORITHMS_ID
        self._alg_id = value

    @property
    def event(self):
        """
        Property for type event
        """
        return self._event

    @event.setter
    def event(self, value):
        """
        Setter for the property event
        """
        assert value is not None
        self._event = value

    @property
    def digest(self):
        """
        Property for type digest
        """
        return self._digest

    @digest.setter
    def digest(self, value):
        """
        Setter for the property digest
        """
        assert value is not None
        self._digest = value

class EventlogUtility:
    """
    Common utility for eventlog related actions
    """

    def __init__(self, target="unix:/run/ccnp/uds/eventlog.sock"):
        if target[:5] != "unix:":
            raise ValueError("Invalid server path, only unix domain socket supported.")

        if not os.path.exists(target.replace('unix:', '')):
            raise FileNotFoundError('eventlog socket does not exist.')
        self._channel = grpc.insecure_channel(target)
        try:
            grpc.channel_ready_future(self._channel).result(timeout=TIMEOUT)
        except grpc.FutureTimeoutError as err:
            raise ConnectionRefusedError('Connection to eventlog server failed') from err
        self._stub = eventlog_server_pb2_grpc.EventlogStub(self._channel)
        self._request = eventlog_server_pb2.GetEventlogRequest()
        self._raw_eventlogs = ""

    def setup_eventlog_request(self, eventlog_level=0, eventlog_category=0,
                               start_position=None, count=None):
        """
        Generate grpc request to get eventlog

        Args:
            eventlog_level(LEVEL): level of event logs to fetch(platform level
            or container level)
            eventlog_category(CATEGORY): different category of event logs to fetch
            start_position(int): start position of event log to fetch
            count(int): number of event log to fetch
        """
        self._request = eventlog_server_pb2.GetEventlogRequest(
            eventlog_level=eventlog_level,
            eventlog_category=eventlog_category,
            start_position=start_position,
            count=count)

    def cleanup_channel(self):
        """ Close the channel used for grpc """
        self._channel.close()

    def get_raw_eventlogs(self):
        """
        Get raw eventlogs

        Args:
          request (GetEventlogRequest): request data
          stub (EventlogStub): the stub to call server

        Returns:
          string: json string of eventlogs
        """

        e = self._stub.GetEventlog(self._request)
        if e.eventlog_data_loc == "":
            LOG.info("Failed to get eventlog from server.")
            return ""

        LOG.info("Fetch eventlog successfully.")
        with open(e.eventlog_data_loc, 'r', encoding='utf-8') as f:
            self._raw_eventlogs = f.read()
        return ""

    def parse_saas_eventlogs(self, eventlogs) -> List[CCEventLogEntry]:
        """
        Parse SaaS level eventlog into CCEventLogEntry

        Args:
          eventlogs (dict): raw eventlog data

        Returns:
          array: list of CCEventLogEntry
        """
        LOG.info("Not implemented")
        return []

    def parse_eventlogs(self, eventlogs) -> List[CCEventLogEntry]:
        """
        Parse eventlog into CCEventLogEntry

        Args:
          eventlogs (dict): raw eventlog data

        Returns:
          array: list of CCEventLogEntry
        """
        if self._request.eventlog_level == eventlog_server_pb2.LEVEL.SAAS:
            return self.parse_saas_eventlogs(eventlogs)

        event_log_array = []
        eventlog_list = eventlogs['EventLogs']

        etypes = CCEventLogType()
        etypes.event_log_dict()

        CCAlgorithms.algo_dict()
        algs = CCAlgorithms()

        for item in eventlog_list:
            etypes.log_type = item['Etype']
            algs.algo_id = item['AlgorithmId']
            digest_num = item['DigestCount']

            if digest_num < 1:
                LOG.info("No digest available")
                continue
            digests = item['Digests']

            event_log = CCEventLogEntry()
            if self._request.eventlog_category == eventlog_server_pb2.CATEGORY.TDX_EVENTLOG:
                event_log.reg_idx = item['Rtmr']
            else:
                event_log.reg_idx = item['Pcr']
            event_log.evt_type = etypes.log_type
            event_log.evt_type_str = etypes.log_type_string(item['Etype'])
            event_log.evt_size = item['EventSize']
            event_log.alg_id = algs
            event_log.event = item['Event']
            event_log.digest = digests[digest_num-1]
            event_log_array.append(event_log)

        return event_log_array

    def get_eventlog(self)-> List[CCEventLogEntry]:
        """
        Get eventlog function to fetch event logs

        Returns:
          array: list of CCEventLogEntry
        """
        self.get_raw_eventlogs()
        self.cleanup_channel()

        if self._raw_eventlogs == "":
            LOG.info("No eventlog found.")
            return None

        eventlog_dict = json.loads(self._raw_eventlogs)
        cc_event_logs = self.parse_eventlogs(eventlog_dict)

        return cc_event_logs

    @classmethod
    def get_platform_eventlog(cls, eventlog_type=eventlog_server_pb2.CATEGORY.TDX_EVENTLOG,
            start_position=None, count=None) -> List[CCEventLogEntry]:
        """
        Get eventlogs from platform perspective.
        Currently, support event log fetching on Intel TDX and TPM.

        Args:
            eventlog_type(EventlogType): type of event log to fetch
            start_position(int): start position of event log to fetch
            count(int): number of event logs to fetch

        Returns:
            array: list of CCEventlogEntry
        """
        if not EventlogType.is_valid_type(eventlog_type):
            raise ValueError("Invalid eventlog type specified")

        if start_position is not None:
            if not isinstance(start_position, int) or start_position < 0:
                raise ValueError("Invalid value specified for start_position")

        if count is not None:
            if not isinstance(count, int) or count <= 0:
                raise ValueError("Invalid value specified for count")

        eventlog_class = cls()
        eventlog_class.setup_eventlog_request(eventlog_server_pb2.LEVEL.PAAS, eventlog_type,
                start_position, count)
        cc_event_logs = eventlog_class.get_eventlog()

        return cc_event_logs

    @classmethod
    def get_container_eventlog(cls):
        """
        Get eventlogs from container perspective

        """
        raise NotImplementedError("Not implemented")


class EventlogType:

    # Get TDX event logs
    TYPE_TDX = eventlog_server_pb2.CATEGORY.TDX_EVENTLOG
    # Get TPM event logs
    TYPE_TPM = eventlog_server_pb2.CATEGORY.TPM_EVENTLOG

    _type_dict = None

    @classmethod
    def event_log_type_dict(cls):
        """
        Class method to construct the event log typedict
        """
        if cls._type_dict is not None:
            return cls._type_dict

        # first time initialization
        cls._type_dict = {}
        for key, value in cls.__dict__.items():
            if key.startswith('TYPE_'):
                # pylint: disable=E1137
                cls._type_dict[value] = key
        return cls._type_dict

    @classmethod
    def is_valid_type(cls, value):
        """
        Class method to check if value is a valid eventlog type
        """
        cls.event_log_type_dict()
        if cls._type_dict is None:
            return False
        for key, _ in cls._type_dict.items():
            if key == value:
                return True
        return False
