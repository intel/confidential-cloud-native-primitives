#!/usr/bin/python3 -u
# SPDX-license-identifier: Apache-2.0
# Copyright (c) 2023 Intel Corporation. All rights reserved.

"""
This file contains the test cases against the get_eventlog API
"""

import logging
import pytest
import grpc

from ccnp.eventlog.eventlog_sdk import (
    EventlogType,
    EventlogUtility,
)

__author__ = "Ruoyu Ying"
__copyright__ = "Ruoyu Ying"
__license__ = "Apache-2.0"

class TestEventlog():

    logging.basicConfig(level=logging.DEBUG)

    @pytest.mark.eventlog
    def test_eventlog_platform(self):
        """
        API test for platform level eventlog
        """
        # test get platform level eventlog using default params
        event_logs = EventlogUtility.get_platform_eventlog()
        assert len(event_logs) > 0

        # test get platform level eventlog with type 'TYPE_TDX'
        event_logs = EventlogUtility.get_platform_eventlog(EventlogType.TYPE_TDX)
        assert len(event_logs) > 0

        # test get eventlogs and check values inside
        event_logs = EventlogUtility.get_platform_eventlog(EventlogType.TYPE_TDX)
        assert len(event_logs) > 0
        assert event_logs[2].evt_type is not None
        assert event_logs[2].evt_type_str is not None
        assert event_logs[2].evt_size is not None
        assert event_logs[2].reg_idx is not None
        assert event_logs[2].alg_id is not None
        assert event_logs[2].event is not None
        assert event_logs[2].digest is not None

        # test get platform level eventlog with type 'TYPE_TDX', start_position and count
        event_logs = EventlogUtility.get_platform_eventlog(EventlogType.TYPE_TDX, 2, 1)
        assert len(event_logs) == 1
        assert event_logs[0].evt_type is not None
        assert event_logs[0].evt_type_str is not None
        assert event_logs[0].evt_size is not None
        assert event_logs[0].reg_idx is not None
        assert event_logs[0].alg_id is not None
        assert event_logs[0].event is not None
        assert event_logs[0].digest is not None

    @pytest.mark.eventlog
    def test_eventlog_tpm(self):
        """
        API test for eventlog with category
        TPM_PCR - success cases

        * Not implemented, just reserve for testing
        """
        # test get platform level eventlog with type 'TYPE_TPM'
        with pytest.raises(grpc.RpcError) as e:
            EventlogUtility.get_platform_eventlog(EventlogType.TYPE_TPM)
        exec_msg = e.value.details()
        assert exec_msg == \
                "stat /sys/kernel/security/tpm0/binary_bios_measurements: no such file or directory"

    @pytest.mark.eventlog
    def test_container_eventlog(self):
        """
        API test for container eventlog

        * Not implemented, just reserve for testing
        """
        with pytest.raises(NotImplementedError) as e:
            EventlogUtility.get_container_eventlog()
        exec_msg = e.value.args[0]
        assert exec_msg == "Not implemented"

    @pytest.mark.eventlog
    def test_eventlog_raise_error(self):
        """
        API test for eventlog - catch error cases
        """
        # test get eventlog with invalid type
        with pytest.raises(ValueError) as e:
            EventlogUtility.get_platform_eventlog("InvalidType")
        exec_msg = e.value.args[0]
        assert exec_msg == "Invalid eventlog type specified"

        # test get eventlog with invalid start_position
        with pytest.raises(ValueError) as e:
            EventlogUtility.get_platform_eventlog(EventlogType.TYPE_TDX, -1)
        exec_msg = e.value.args[0]
        assert exec_msg == "Invalid value specified for start_position"

        # test get eventlog with invalid count
        with pytest.raises(ValueError) as e:
            EventlogUtility.get_platform_eventlog(EventlogType.TYPE_TDX, 1, -1)
        exec_msg = e.value.args[0]
        assert exec_msg == "Invalid value specified for count"

        # test get eventlog with invalid count - exceed length
        with pytest.raises(grpc.RpcError) as e:
            EventlogUtility.get_platform_eventlog(EventlogType.TYPE_TDX, 1, 10000)
        exec_msg = e.value.details()
        assert exec_msg == "Invalid count exceeds event log length"
