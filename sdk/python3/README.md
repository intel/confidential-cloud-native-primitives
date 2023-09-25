# Confidential Cloud-Native Primitives SDK for Python

The Confidential Cloud-Native Primitives (CCNP) project is the solution targeted on simplifying the use of Trusted Execution Environment (TEE) in cloud-native environment. Currently, there are 2 parts included in CCNP, the services and the SDK.

- Service is designed to hide the complexity of different TEE platforms and provides common interfaces and scalability for cloud-native environment.
- SDK is to simplify the use of the service interface for development, it covers communication to the service and parses the results from the services.

The service supports attestation, measurement fetching and event log collecting of various platforms including Intel Trusted Domain Extensions (TDX), Trusted Platform Modules (TPM) and AMD SEV-SNP. More platforms will be supported later.

Attestation is a common process within TEE platform and TPM to verify if the software binaries were properly instantiated on a trusted platform. Third parties can leverage the attestation process to identify the trustworthiness of the platform (by checking the measurements or event logs) as well as the software running on it, in order to decide whether they shall put their confidential information/workload onto the platform.

CCNP, as the overall framework for attestation, measurement and event log fetching, provides user with both customer-facing SDK and overall framework. By leveraging this SDK, user can easily retrieve different kinds of measurements or evidence such as event logs. Working along with different verification services (such as Amber) and configurable policies, user can validate the trustworthiness of the  platform and make further decision.

[Source code][source_code]
| [Package (PyPI)][ccnp_pypi]
| [API reference documentation][api_doc]

## Getting started

### Prerequisites
In order to work properly, user need to have the backend services ready on the TEE or TPM enabled platform first. Please refer to each deployment guide reside in the [service](../../service/) folder to install the backend services.

### Install the package
User can install the CCNP client library for Python with PyPI:

```
pip install ccnp
```

To install from source code, user can use the following command:

```
pip install -e .
```

## Key concepts and usage
There are three major functionalities provided in this SDK:

* [Quote fetching](#quote)
* [Measurement fetching](#measurement)
* [Event log fetching](#event-log)

### Quote

Using this SDK, user could fetch the quote from different platforms, the service detect the platform automatically and return the type and the quote.

#### Quote type for platform

* TYPE_TDX - This provides the quote fetching based on Intel TDX.
* TYPE_TPM - This provides the quote fetching based on TPM.

#### Example usage of quote SDK

The interface input of quote is `nonce` and `user_data`, both of them are optional and will be measured in quote.
Here are the example usages of quote SDK:

* Fetch quote without any inputs
```python
from ccnp import Quote

quote = Quote.get_quote()

print(quote.quote_type)
print(quote.quote)

```

* Fetch quote with a `nonce`
```python
import secrets
from ccnp import Quote

nonce = secrets.token_urlsafe()
quote = Quote.get_quote(nonce=nonce)

print(quote.quote_type)
print(quote.quote)

```

* Fetch quote with a `nonce` and `user_data`
```python
import base64
import secrets
from ccnp import Quote

nonce = secrets.token_urlsafe()
user_data = base64.b64encode(b'This data should be measured.')
quote = Quote.get_quote(nonce=nonce, user_data=user_data)

print(quote.quote_type)
print(quote.quote)

# For TD quote, it includes RTMRs, TD report, etc.
if quote.quote_type == Quote.TYPE_TDX:
  print(quote.rtmrs)
  print(quote.tdreport)
```

### Measurement

Using this SDK, user could fetch various measurements from different perspective and categories.
Basic support on measurement focus on the platform measurements, including TEE report, values within TDX RTMR registers or values reside in TPM PCR registers.
There's also advanced support to provide measurement for a certain workload or container. The feature is still developing in progress.

#### MeasurementType for platform

The measurement SDK supports fetching different types of evidence depending on the environment.
Currently, CCNP supports the following categories of measurements:

* TYPE_TEE_REPORT - This provides the report fetching on various Trusted Execution Environment from all kinds of vendors, including Intel TDX, AMD SEV (Working in Progress), etc.
* TYPE_TDX_RTMR - This provides the measurement fetching on TDX RTMR. Users could fetch the measurement from one single RTMR register with its index.
* TYPE_TPM_PCR - This provides th measurement fetching on TPM PCR. Users could fetch measurement from one single PCR register with its index.

#### Example usage of measurement SDK

Here are the example usages for measurement SDK:

* Fetch TEE report base on platform
```python
from ccnp import Measurement
from ccnp import MeasurementType

# Fetch TEE report without user data
report = Measurement.get_platform_measurement()

# Fetch TEE report with user data
data = "testing"
report = Measurement.get_platform_measurement(MeasurementType.TYPE_TEE_REPORT, data)

```

* Fetch single RTMR measurement for platform
```python
from ccnp import Eventlog
from ccnp import Measurement
from ccnp import MeasurementType

# Fetch the value reside in register 1 of RTMR
rtmr_measurement = Measurement.get_platform_measurement(MeasurementType.TYPE_TDX_RTMR, None, 1)
```

* Fetch container measurement (Working in Progress)
```python
from ccnp import Measurement
from ccnp import MeasurementType

container_measurement = Measurement.get_container_measurement()
```

### Event log

Using this SDK, user can fetch the event logs to assist the attestation/verification process. It also enables two different categories of event logs - for the platform or for a single workload/container.
From platform perspective, it can support different Trusted Execution Environment and TPM. This sdk can also do fetching on certain number of event logs.

#### EventlogType for platform

* TYPE_TDX - This provides the event log fetching based on Intel TDX.
* TYPE_TPM - This provides the event log fetching based on TPM.

#### Example usage of Eventlog SDK

Here are the example usages of eventlog SDK:

* Fetch event log of Intel TDX platform for platform and check the information inside
```python
from ccnp import Eventlog
from ccnp import EventlogType

# default type for get_platform_eventlog() is 'TYPE_TDX'
logs = Eventlog.get_platform_eventlog()
# same as setting type as 'TYPE_TDX'
logs = Eventlog.get_platform_eventlog(EventlogType.TYPE_TDX)

# show total length
print(len(logs))

# fetch event log attributes
print(logs[2].evt_type)
print(logs[2].evt_type_str)
print(logs[2].evt_size)
print(logs[2].reg_idx)
print(logs[2].alg_id)
print(logs[2].event)
print(logs[2].digest)

# fetch 5 event logs from the second one
logs = Eventlog.get_platform_eventlog(EventlogType.TYPE_TDX, 2, 5)

# show log length, which shall equal to 5
print(len(logs))
```

* Fetch event log of TPM platform (Working in Progress)
```python
from ccnp import Eventlog
from ccnp import EventlogType

# set type for get_platform_eventlog() as 'TYPE_TPM'
logs = Eventlog.get_platform_eventlog(EventlogType.TYPE_TPM)
```

* Fetch event log for certain container (Working in Progress)
```python
from ccnp import Eventlog
from ccnp import EventlogType

logs = Eventlog.get_container_eventlog()
```

## End-to-end examples

TBA.

## Troubleshooting

Troubleshooting information for the CCNP SDK can be found here.

## Next steps
For more information about the Confidential Cloud-Native Primitives, please see our documentation page.

## Contributing
This project welcomes contributions and suggestions. Most contributions require you to agree to a Contributor License Agreement (CLA) declaring that you have the right to, and actually do, grant us the rights to use your contribution. For details, visit the Contributor License Agreement site.

When you submit a pull request, a CLA-bot will automatically determine whether you need to provide a CLA and decorate the PR appropriately (e.g., label, comment). Simply follow the instructions provided by the bot. You will only need to do this once across all repos using our CLA.

See [CONTRIBUTING.md](../../CONTRIBUTING.md) for details on building, testing, and contributing to these libraries.

## Provide Feedback
If you encounter any bugs or have suggestions, please file an issue in the Issues section of the project.

<!-- LINKS -->
[source_code]: https://github.com/intel/confidential-cloud-native-primitives/tree/main/sdk/python3
[ccnp_pypi]: https://pypi.org/project/ccnp/
[api_doc]: https://intel.github.io/confidential-cloud-native-primitives/_rst/sdk.readme.html
