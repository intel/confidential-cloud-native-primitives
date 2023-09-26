# Confidential Cloud-Native Primitives (CCNP)

![CI Check License](https://github.com/intel/confidential-cloud-native-primitives/actions/workflows/pr-license-python.yaml/badge.svg)
![CI Check Spelling](https://github.com/intel/confidential-cloud-native-primitives/actions/workflows/pr-doclint.yaml/badge.svg)
![CI Check Python](https://github.com/intel/confidential-cloud-native-primitives/actions/workflows/pr-pylint.yaml/badge.svg)
![CI Check Shell](https://github.com/intel/confidential-cloud-native-primitives/actions/workflows/pr-shell-check.yaml/badge.svg)
![CI Check Rust](https://github.com/intel/confidential-cloud-native-primitives/actions/workflows/pr-check-rust.yaml/badge.svg)
![CI Check Golang](https://github.com/intel/confidential-cloud-native-primitives/actions/workflows/pr-golang-check.yaml/badge.svg)
![CI Check Container](https://github.com/intel/confidential-cloud-native-primitives/actions/workflows/pr-container-check.yaml/badge.svg)

## Introduction

VM(Virtual Machine) based confidential computing like Intel TDX provides isolated encryption runtime environment based on
hardware Trusted Execution Environment (TEE) technologies. To land cloud native computing into confidential environment,
there are lots of different PaaS frameworks such as confidential cluster, confidential container, which brings challenges
for enabling and TEE measurement.
This project uses cloud native design pattern to implement confidential computing primitives like event log, measurement,
quote and attestation. It also provides new features design to address new challenges like how to auto scale trustworthy,
how to reduce TCB size, etc.

The project itself contains several parts: the services, the SDK and related dependencies

- Services are designed to hide the complexity of different TEE platforms and provides common interfaces and scalability
for cloud-native environment to address the fetching the fetching of quote, measurement and event log.
- SDK is provided to simplify the use of the service interface for development, it covers communication to the service
and parses the results from the services. With such SDK, users can perform related actions with one simple API call.
- A ccnp device plugin is provided as the dependency for services such as Quote Server and Measurement Server. It will help with
device mount and folder injection within the service.

SDK PyPI package can be found [here](https://pypi.org/project/ccnp/). Please check our [documentation](https://intel.github.io/confidential-cloud-native-primitives/) for more details.

![](docs/ccnp_arch.png)

*Note: For Intel TDX, it bases on Linux TDX Software Stack at [tdx-tools](https://github.com/intel/tdx-tools), the corresponding white
paper is at [Whitepaper: Linux* Stacks for Intel® Trust Domain Extension 1.0](https://www.intel.com/content/www/us/en/content-details/779108/whitepaper-linux-stacks-for-intel-trust-domain-extension-1-0.html).*


## Installation

Here provides the description on the installation steps for the services and the SDK.

For Services, we provided installation using Helm or yaml. And there are also several deployment modes available for
the services. User can find the details in the 'Installation' section within the README file under each service folder.

- Quote Server: [Installation guide](service/quote-server/README.md)
- Measurement Server: [Installation guide](service/measurement-server/README.md)
- Event Log Server: [Installation guide](service/eventlog-server/README.md)

For SDK, user can simply install from PyPI using command:

```
pip install ccnp
```

Or to install from source code with the following command:

```
cd sdk/python3
pip install -e .
```

For the ccnp device plugin, user can find the installation guide under the 'Installation' section [here](device-plugin/ccnp-device-plugin/README.md)

## Contributing

This project welcomes contributions and suggestions. Most contributions require you to agree to a Contributor License Agreement (CLA) declaring that you have the right to, and actually do, grant us the rights to use your contribution. For details, contact the maintainers of the project.

When you submit a pull request, a CLA-bot will automatically determine whether you need to provide a CLA and decorate the PR appropriately (e.g., label, comment). Simply follow the instructions provided by the bot. You will only need to do this once across all repos using our CLA.

See [CONTRIBUTING.md](CONTRIBUTING.md) for details on building, testing, and contributing to these libraries.

## Provide Feedback

If you encounter any bugs or have suggestions, please file an issue in the Issues section of the project.


_Note: This is pre-release/prototype software and, as such, it may be substantially modified as updated versions are made available._
