import pytest
import base64
from ccnp import Measurement
from ccnp import MeasurementType
from pytdxattest.tdreport  import TdReport


class TestTdreport:
    def test_without_data(self):
        ccnp_report = Measurement.get_platform_measurement()
        py_report = TdReport.get_td_report()
        assert base64.b64decode(ccnp_report) == py_report.data

