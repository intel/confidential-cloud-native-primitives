"""CCNP framework to enable TEE related operations in cloud native environments"""

__version__ = "0.0.1"

from .eventlog.eventlog_sdk import EventlogUtility as Eventlog
from .eventlog.eventlog_sdk import EventlogType

from .measurement.measurement_sdk import MeasurementUtility as Measurement
from .measurement.measurement_sdk import MeasurementType

from .quote.quote_sdk import Quote
