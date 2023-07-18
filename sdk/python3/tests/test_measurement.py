#!/usr/bin/python3 -u
# SPDX-license-identifier: Apache-2.0
# Copyright (c) 2023 Intel Corporation. All rights reserved.

"""
This file contains the test cases against the get_measurement API
"""

import logging
import pytest

from ccnp.measurement.measurement_sdk import (
    MeasurementType,
    MeasurementUtility,
)

_invalid_report_data_1 = "MTIzNDU2NzgxMjM0NTY3ODEyMzQ1Njc4MTIzNDU2N\
    zgxMjM0NTY3ODEyMzQ1Njc4MTIzNDU2NzgxMjM0NTY3ODEyMzQ1Njc4MTIzNDU2\
    NzgxMjM0NTY3ODEyMzQ1Njc4MTIzNDU2NzgxMjM0NTY3ODEyNTY3ODEyMzQ1Njc\
    4MTIzNDU2NzgxMjM0NTY3ODEyMzQ1Njc4MTIzNDU2NzgxMjM0NTY3ODEyMzQ1\
    Njc4MTIzNDU2NzgxMjM0NTY3ODEyMzQ1Njc4MTIzNDU2NzgxMjM0NTY3ODEyMzQ\
    1Njc4MTIzNDU2NzgxMjM0NTY3ODEyMzQ1Njc4MTIzNDU2NzgxMjM0NTY3ODEyMzQ1Njc4Cg=="
_invalid_report_data_2 = 1

__author__ = "Ruoyu Ying"
__copyright__ = "Ruoyu Ying"
__license__ = "Apache-2.0"

class TestMeasurement():

    logging.basicConfig(level=logging.DEBUG)

    @pytest.mark.measurement
    def test_measurement_tee_report(self):
        """
        API test for measurement with type
        TEE report
        """
        # test get tee report with default param
        report = MeasurementUtility.get_platform_measurement()
        assert report is not None

        # test get tee report with param set as 'TYPE_TEE_REPORT'
        report_1 = MeasurementUtility.get_platform_measurement(MeasurementType.TYPE_TEE_REPORT)
        assert report_1 is not None

        # test get tee report with param set as 'TYPE_TEE_REPORT'
        # and report data set as "test"
        report_2 = MeasurementUtility.get_platform_measurement(MeasurementType.TYPE_TEE_REPORT,\
                                                               "test")
        assert report_2 is not None

        # test get measurement with invalid report data
        with pytest.raises(ValueError) as e:
            MeasurementUtility.get_platform_measurement(MeasurementType.TYPE_TEE_REPORT,\
                                                        _invalid_report_data_1)
        exec_msg = e.value.args[0]
        assert exec_msg == "Invalid report data specified"

        with pytest.raises(ValueError) as e:
            MeasurementUtility.get_platform_measurement(MeasurementType.TYPE_TEE_REPORT,\
                                                        _invalid_report_data_2)
        exec_msg = e.value.args[0]
        assert exec_msg == "Invalid report data specified"

    @pytest.mark.measurement
    def test_measurement_tdx_rtmr(self):
        """
        API test for measurement with type
        TDX RTMR - success cases
        """
        # test get tee report with param set as 'TYPE_TDX_RTMR' and use default index
        measurement_1 = MeasurementUtility.get_platform_measurement(MeasurementType.TYPE_TDX_RTMR)
        assert measurement_1 is not None

        # test get tee report with param set as 'TYPE_TDX_RTMR'
        measurement_2 = MeasurementUtility.get_platform_measurement(MeasurementType.TYPE_TDX_RTMR, "", 1)
        assert measurement_2 is not None

    @pytest.mark.measurement
    def test_measurement_tpm(self):
        """
        API test for measurement with category
        TPM_PCR - success cases

        * Not implemented, just reserve for testing
        """
        # test get measurement with param set as 'TYPE_TPM_PCR' and use default index
        #measurement_1 = MeasurementUtility.get_platform_measurement(MeasurementType.TYPE_TPM_PCR)
        #assert measurement_1 is not None

        # test get measurement with param set as 'TYPE_TPM_PCR'
        #measurement_2 = MeasurementUtility.get_platform_measurement(MeasurementType.TYPE_TPM_PCR, "", 1)
        #assert measurement_2 is not None

    @pytest.mark.measurement
    def test_container_measurement(self):
        """
        API test for container measurement

        * Not implemented, just reserve for testing
        """
        with pytest.raises(NotImplementedError) as e:
            MeasurementUtility.get_container_measurement()
        exec_msg = e.value.args[0]
        assert exec_msg == "Not implemented"

    @pytest.mark.measurement
    def test_measurement_raise_error(self):
        """
        API test for measurement - catch error cases
        """
        # test get measurement with invalid type
        with pytest.raises(ValueError) as e:
            MeasurementUtility.get_platform_measurement("InvalidType")
        exec_msg = e.value.args[0]
        assert exec_msg == "Invalid measurement type specified"

        # test get measurement with invalid register index
        with pytest.raises(ValueError) as e:
            MeasurementUtility.get_platform_measurement(MeasurementType.TYPE_TPM_PCR, "", -1)
        exec_msg = e.value.args[0]
        assert exec_msg == "Invalid value specified for register index"

        with pytest.raises(ValueError) as e:
            MeasurementUtility.get_platform_measurement(MeasurementType.TYPE_TPM_PCR, "", 17)
        exec_msg = e.value.args[0]
        assert exec_msg == "Invalid value specified for register index"
