# Confidential Cloud-Native Primitives (CCNP)

![CI Check License](https://github.com/intel/confidential-cloud-native-primitives/actions/workflows/pr-license-python.yaml/badge.svg)
![CI Check Spelling](https://github.com/intel/confidential-cloud-native-primitives/actions/workflows/pr-doclint.yaml/badge.svg)
![CI Check Python](https://github.com/intel/confidential-cloud-native-primitives/actions/workflows/pr-pylint.yaml/badge.svg)
![CI Check Shell](https://github.com/intel/confidential-cloud-native-primitives/actions/workflows/pr-shell-check.yaml/badge.svg)
![CI Check Rust](https://github.com/intel/confidential-cloud-native-primitives/actions/workflows/pr-check-rust.yaml/badge.svg)
![CI Check Golang](https://github.com/intel/confidential-cloud-native-primitives/actions/workflows/pr-golang-check.yaml/badge.svg)
![CI Check Container](https://github.com/intel/confidential-cloud-native-primitives/actions/workflows/pr-container-check.yaml/badge.svg)

## 1. Introduction

Confidential Computing technology like Intel TDX provides isolated encryption runtime environment to protect data-in-use
based on hardware Trusted Execution Environment (TEE).
It requires a full chain integrity measurement on the launch-time or runtime environment to guarantee "consistently
behavior in expected way" (defined by [Trusted Computing](https://en.wikipedia.org/wiki/Trusted_Computing))
of confidential computing environment for tenant's zero-trust use case.

This project is designed to provide cloud native measurement for the full measurement chain from TEE TCB -> Firmware
TCB -> Guest OS TCB -> Cloud Native TCB as follows:

![](/docs/cc-full-meaurement-chain.png)

From the perspective of tenant's workload, it will expose the [CC Trusted API](https://github.com/cc-api/cc-trusted-api)
as the unified interfaces across diverse trusted foundations like `RTMR+TDMR+CCEL` and `PCR+TPM2`. The definitions and
structures follows standard specifications like [TCG PC Client Platform TPM Profile Specification](https://trustedcomputinggroup.org/resource/pc-client-platform-tpm-profile-ptp-specification/),
[TCG PC Client Platform Firmware Profile Specification](https://trustedcomputinggroup.org/resource/pc-client-specific-platform-firmware-profile-specification/)

![](/docs/ccnp-architecture-high-level.png)

This project should also be able deployed on diverse cloud native PaaS frameworks like confidential cluster, container, kubevirt etc.
An example landing architecture on confidential cluster is as follows:

![](/docs/ccnp-landing-confidential-cluster.png)


Finally, the full trusted chain will be measured into CC report as follows using TDX as example:

![](/docs/cc-full-measurement-tdreport.png)

Refer: https://github.com/tianocore/edk2/blob/master/MdePkg/Include/IndustryStandard/Tdx.h

## 2. Design

CCNP includes several micro-services as BaaS(Backend as a Service) to provides cloud native measurement, and exposed `CC trusted API`
via cloud native SDK:

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


## 3. Installation

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

## 4. Contributing

This project welcomes contributions and suggestions. Most contributions require you to agree to a Contributor License Agreement (CLA) declaring that you have the right to, and actually do, grant us the rights to use your contribution. For details, contact the maintainers of the project.

When you submit a pull request, a CLA-bot will automatically determine whether you need to provide a CLA and decorate the PR appropriately (e.g., label, comment). Simply follow the instructions provided by the bot. You will only need to do this once across all repos using our CLA.

See [CONTRIBUTING.md](CONTRIBUTING.md) for details on building, testing, and contributing to these libraries.

## 5. Provide Feedback

If you encounter any bugs or have suggestions, please file an issue in the Issues section of the project.


_Note: This is pre-release/prototype software and, as such, it may be substantially modified as updated versions are made available._

### Contributors
CONTRIBUTOR