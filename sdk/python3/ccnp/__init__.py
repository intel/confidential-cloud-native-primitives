"""CCNP framework to enable TEE related operations in cloud native environments"""

__version__ = "0.0.1"

from .eventlog.eventlog_sdk import EventlogUtility as EventlogClient
# pylint: disable=E0611
from .eventlog.eventlog_server_pb2 import LEVEL as EventlogType
from .eventlog.eventlog_server_pb2 import CATEGORY as EventlogCategory

from .measurement.measurement_sdk import MeasurementUtility as MeasurementClient
from .measurement.measurement_server_pb2 import TYPE as MeasurementType
from .measurement.measurement_server_pb2 import CATEGORY as MeasurementCategory
