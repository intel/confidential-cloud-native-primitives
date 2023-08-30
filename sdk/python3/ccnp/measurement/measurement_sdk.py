# pylint: disable=duplicate-code
# Copyright (c) 2023, Intel Corporation. All rights reserved.<BR>
# SPDX-license-identifier: Apache-2.0
"""This module provides the functions to talk to measurement-server and fetch measurements"""

import logging
import os

import grpc
# pylint: disable=E1101
from ccnp.measurement import measurement_server_pb2
from ccnp.measurement import measurement_server_pb2_grpc

LOG = logging.getLogger(__name__)
TIMEOUT = 5


class MeasurementUtility:
    """
    Common utility for measurement related actions
    """

    def __init__(self, target="unix:/run/ccnp/uds/measurement.sock"):
        if target[:5] != "unix:":
            raise ValueError("Invalid server path, only unix domain socket supported.")

        if not os.path.exists(target.replace('unix:', '')):
            raise FileNotFoundError('Measurement socket does not exist.')
        self._channel = grpc.insecure_channel(target)
        try:
            grpc.channel_ready_future(self._channel).result(timeout=TIMEOUT)
        except grpc.FutureTimeoutError as err:
            raise ConnectionRefusedError('Connection to measurement server failed') from err
        self._stub = measurement_server_pb2_grpc.MeasurementStub(self._channel)
        self._request = measurement_server_pb2.GetMeasurementRequest()

    def setup_measurement_request(self, measurement_type=0, measurement_category=0,
                                  report_data=None, register_index=0):
        """ Function to generate a get_measurement request
        Setup a measurement request for get_measurement API

        Args:
            measurement_type(TYPE): type of measurement to fetch - PaaS or SaaS
            measurement_category(CATEGORY): category of measurement to fetch
            report_data(str): user data to get wrapped as part of tee report
            register_index(int): register index to fetch measurements
        """
        self._request = measurement_server_pb2.GetMeasurementRequest(
            measurement_type=measurement_type,
            measurement_category=measurement_category,
            report_data=report_data,
            register_index=register_index
        )

    def cleanup_channel(self):
        """ Clean up channel used for grpc """
        self._channel.close()

    def get_measurement(self):
        """
        Get measurement

        Args:
          request (GetMeasurementRequest): request data
          stub (MeasurementStub): the stub to call server

        Returns:
          string: base64 encoded string of measurement
        """

        reply_data = self._stub.GetMeasurement(self._request)
        if reply_data.measurement == "":
            LOG.info("Failed to get measurement from server.")
            return ""

        LOG.info("Fetch measurement successfully.")
        return reply_data.measurement

    @classmethod
    def get_platform_measurement(cls, measurement_type=measurement_server_pb2.CATEGORY.TEE_REPORT,
            report_data=None, register_index=None) -> str:
        """
        Get measurements from platform perspective.
        Currently, support measurement fetching on TEE reports, Intel TDX RTMR and TPM.

        Args:
            measurement_type(MeasurementType): type of measurement to fetch
            report_data(str): user data to be wrapped as part of TEE report
            register_index(int): register index used to fetch TDX RTMR or TPM PCR measurement

        Returns:
            string: base64 encoded measurement string
        """
        if not MeasurementType.is_valid_type(measurement_type):
            raise ValueError("Invalid measurement type specified")

        if report_data is not None:
            if not isinstance(report_data, str) or len(report_data) > 64:
                raise ValueError("Invalid report data specified")

        if register_index is not None:
            if not isinstance(register_index, int) or register_index < 0 or register_index > 16:
                raise ValueError("Invalid value specified for register index")

        measurement_class = cls()
        measurement_class.setup_measurement_request(measurement_server_pb2.TYPE.PAAS,
                measurement_type, report_data, register_index)
        return measurement_class.get_measurement()

    @classmethod
    def get_container_measurement(cls) -> str:
        """
        Get measurements from container perspective

        """
        raise NotImplementedError("Not implemented")

class MeasurementType:

    # Get TEE report
    TYPE_TEE_REPORT = measurement_server_pb2.CATEGORY.TEE_REPORT
    # Get TDX RTMR measurement (of a specific register)
    TYPE_TDX_RTMR = measurement_server_pb2.CATEGORY.TDX_RTMR
    # Get TPM PCR measurement (of a specific register)
    TYPE_TPM_PCR =measurement_server_pb2.CATEGORY.TPM

    _type_dict = None

    @classmethod
    def measurement_type_dict(cls):
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
        cls.measurement_type_dict()
        if cls._type_dict is None:
            return False
        for key, _ in cls._type_dict.items():
            if key == value:
                return True
        return False
