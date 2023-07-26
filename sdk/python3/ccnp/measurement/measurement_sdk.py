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
        """ Function to generate a get_measurement request """
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
