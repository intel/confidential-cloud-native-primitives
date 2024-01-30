# Confidential Cloud-Native Primitives (CCNP)

![CI Check License](https://github.com/intel/confidential-cloud-native-primitives/actions/workflows/pr-license-python.yaml/badge.svg)
![CI Check Spelling](https://github.com/intel/confidential-cloud-native-primitives/actions/workflows/pr-doclint.yaml/badge.svg)
![CI Check Python](https://github.com/intel/confidential-cloud-native-primitives/actions/workflows/pr-pylint.yaml/badge.svg)
![CI Check Shell](https://github.com/intel/confidential-cloud-native-primitives/actions/workflows/pr-shell-check.yaml/badge.svg)
![CI Check Rust](https://github.com/intel/confidential-cloud-native-primitives/actions/workflows/pr-check-rust.yaml/badge.svg)
![CI Check Golang](https://github.com/intel/confidential-cloud-native-primitives/actions/workflows/pr-golang-check.yaml/badge.svg)
![CI Check Container](https://github.com/intel/confidential-cloud-native-primitives/actions/workflows/pr-container-check.yaml/badge.svg)
![CC Foundation Image Customize](https://github.com/intel/confidential-cloud-native-primitives/actions/workflows/image-rewriter.yaml/badge.svg)
[![OpenSSF Best Practices](https://www.bestpractices.dev/projects/8325/badge)](https://www.bestpractices.dev/projects/8325)

## 1. Introduction

Confidential Computing technology like Intel® TDX provides an isolated encryption runtime
environment to protect data-in-use based on hardware Trusted Execution Environment (TEE).
It requires a full chain integrity measurement on the launch-time or runtime environment
to guarantee "consistent behavior in an expected way" of confidential
computing environment for tenant's zero-trust use case.

This project is designed to provide cloud native measurement for the full measurement
chain from TEE TCB -> Firmware TCB -> Guest OS TCB -> Cloud Native TCB as follows:

![](/docs/cc-full-meaurement-chain.png)

_NOTE: Different from traditional trusted computing on non-confidential environment,
the measurement chain is not only started with Guest's `SRTM` (Static Root Of Measurement)
but it also needs to include the TEE TCB, because the CC VM environment is created by TEE
via `DRTM` (Dynamic Root of Measurement) like Intel TXT on the host._

From the perspective of a tenant's workload, `CCNP` will expose the [CC Trusted API](https://github.com/cc-api/cc-trusted-api)
as the unified interfaces across diverse trusted foundations like `RTMR+MRTD+CCEL`
and `PCR+TPM2`. Learn more details of CCNP design at [CCNP documentation](https://intel.github.io/confidential-cloud-native-primitives/).

![](/docs/ccnp-architecture-high-level.png)

Finally, the full trusted chain will be measured into a CC report as follows using
TDX as an example:

![](/docs/cc-full-measurement-tdreport.png)

_NOTE:_

- The measurement of TEE, Guest's boot, OS is per CC VM, but cluster/container measurement
might be per cluster/namespace/container for cloud native architecture.
- Please refer to structure [`TDREPORT`](https://github.com/tianocore/edk2/blob/master/MdePkg/Include/IndustryStandard/Tdx.h)
- The CCNP project collects container level primitives by implementing unified APIs defined in [CC Trusted API](https://github.com/cc-api/cc-trusted-api). The project will be moved to [CC Trusted API](https://github.com/cc-api/cc-trusted-api) in the near future. 


## 2. Installation

### 2.1 Configuration for Host and Guest

CCNP collects primitives of confidential cloud native environments running in confidential VMs(CVM), such as Intel® TDX guest. The primitives are not only from the TEE + CVM boot process + CVM OS but also from the environments running workloads, e.g. Kubernetes cluster or Docker containers. Thus, you need to check below configuration for both hosts and guests.

You can setup an Intel® TDX enlightened host and then boot a TD guest on it. The feasible configurations are as below.

|  CPU | Host OS  | Host packages  | Guest OS  | Guest packages  | Atttestation packages |
|---|---|---|---|---|---|
|  Intel® Emerald Rapids | Ubuntu 22.04| Build packages referring to [here](https://github.com/intel/tdx-tools/tree/tdx-1.5/build/ubuntu-22.04) | Ubuntu 22.04 | Build packages referring to [here](https://github.com/intel/tdx-tools/tree/tdx-1.5/build/ubuntu-22.04) | [here](https://download.01.org/intel-sgx/sgx-dcap/1.19/linux/distro/ubuntu22.04-server/)
| Intel® Emerald Rapids | Ubuntu 23.10 | Setup TDX host referring to [here](https://github.com/canonical/tdx) | Ubuntu 22.04 | Build packages referring to [here](https://github.com/intel/tdx-tools/tree/tdx-1.5/build/ubuntu-22.04)| Setup containerized [PCCS](https://github.com/intel/confidential-cloud-native-primitives/tree/main/container/pccs) and [QGS](https://github.com/intel/confidential-cloud-native-primitives/tree/main/container/qgs) on the host | 

_NOTE: The Platform certificate caching service (PCCS) is used to retrieve and cache PCK certificates locally to your cluster from Intel's Platform Certificate Service. This is necessary to attest the authenticity of a TD guest before a workload is started in it. The Quote Generate Service (QGS) runs on the host in a specialized enclave to generate and use TD quotes. For convenient setup these can run inside a Docker container. Learn more at https://download.01.org/intel-sgx/sgx-dcap/1.17/linux/docs/Intel_TDX_DCAP_Quoting_Library_API.pdf. The PCCS and QGS are used to get Quote for a TD guest. They need to be installed on TDX hosts._

### 2.2 Deploy CCNP Services in Confidential VM

_NOTE: the following installation will be performed in a confidential VM. Make sure you have confidential VM booted before moving forward._

It supports to deploy CCNP services as DaemonSets in Kubernetes cluster or docker containers on a single confidential VM. Please refer to below guides for different deployment environments.

- [CCNP deployment guide - K8S](deployment/README.md): on confidential VM node of Kubernetes cluster.

- [CCNP deployment guide - Docker](deployment/README.md): on confidential VM using docker compose.

### 2.3 Install SDK

CCNP SDK can be used by a workload for cloud native primitives collecting. It needs to be installed within the workload container image and called whenever the primitives are required. For example, in your workload written in Python, you can install the SDK from PyPI using the command:

```
pip install ccnp
```

Alternatively, the CCNP can be installed from source code with the following command. Make sure to clone the repository into your confidential VM and then run the following command:

```
cd sdk/python3
pip install -e .
```

### 2.4 Install CCNP Device Plugin
Follow the CCNP device plugin [Installation Guide](device-plugin/ccnp-device-plugin/README.md)

## 3. Contributing

This project welcomes contributions and suggestions. Most contributions require
you to agree to a Contributor License Agreement (CLA) declaring that you have the
right to, and actually do, grant us the rights to use your contribution. For details,
contact the maintainers of the project.

When you submit a pull request, a CLA-bot will automatically determine whether you
need to provide a CLA and decorate the PR appropriately (e.g., label, comment).
Simply follow the instructions provided by the bot. You will only need to do this
once across all repos using our CLA.

See [CONTRIBUTING.md](CONTRIBUTING.md) for details on building, testing, and contributing
to these libraries.

## 4. Provide Feedback

If you encounter any bugs or have suggestions, please file an issue in the Issues
section of the project.


_Note: This is pre-production software and, as such, it may be substantially modified as updated versions are made available._

## 5. Reference

[Trusted Computing](https://en.wikipedia.org/wiki/Trusted_Computing)

[TCG PC Client Platform TPM Profile Specification](https://trustedcomputinggroup.org/resource/pc-client-platform-tpm-profile-ptp-specification/)

[TCG PC Client Platform Firmware Profile Specification](https://trustedcomputinggroup.org/resource/pc-client-specific-platform-firmware-profile-specification/)

## 6. Contributors

<!-- spell-checker: disable -->

<!-- readme: contributors -start -->
<table>
<tr>
    <td align="center">
        <a href="https://github.com/Ruoyu-y">
            <img src="https://avatars.githubusercontent.com/u/70305231?v=4" width="100;" alt="Ruoyu-y"/>
            <br />
            <sub><b>Ruoyu Ying</b></sub>
        </a>
    </td>
    <td align="center">
        <a href="https://github.com/hairongchen">
            <img src="https://avatars.githubusercontent.com/u/105473940?v=4" width="100;" alt="hairongchen"/>
            <br />
            <sub><b>Hairongchen</b></sub>
        </a>
    </td>
    <td align="center">
        <a href="https://github.com/kenplusplus">
            <img src="https://avatars.githubusercontent.com/u/31843217?v=4" width="100;" alt="kenplusplus"/>
            <br />
            <sub><b>Lu Ken</b></sub>
        </a>
    </td>
    <td align="center">
        <a href="https://github.com/hjh189">
            <img src="https://avatars.githubusercontent.com/u/88485603?v=4" width="100;" alt="hjh189"/>
            <br />
            <sub><b>Jiahao  Huang</b></sub>
        </a>
    </td>
    <td align="center">
        <a href="https://github.com/ruomengh">
            <img src="https://avatars.githubusercontent.com/u/90233733?v=4" width="100;" alt="ruomengh"/>
            <br />
            <sub><b>Ruomeng Hao</b></sub>
        </a>
    </td>
    <td align="center">
        <a href="https://github.com/HaokunX-intel">
            <img src="https://avatars.githubusercontent.com/u/108452001?v=4" width="100;" alt="HaokunX-intel"/>
            <br />
            <sub><b>Haokun Xing</b></sub>
        </a>
    </td></tr>
<tr>
    <td align="center">
        <a href="https://github.com/hwang37">
            <img src="https://avatars.githubusercontent.com/u/36193324?v=4" width="100;" alt="hwang37"/>
            <br />
            <sub><b>Wang, Hongbo</b></sub>
        </a>
    </td>
    <td align="center">
        <a href="https://github.com/dongx1x">
            <img src="https://avatars.githubusercontent.com/u/34326010?v=4" width="100;" alt="dongx1x"/>
            <br />
            <sub><b>Xiaocheng Dong</b></sub>
        </a>
    </td>
    <td align="center">
        <a href="https://github.com/LeiZhou-97">
            <img src="https://avatars.githubusercontent.com/u/102779531?v=4" width="100;" alt="LeiZhou-97"/>
            <br />
            <sub><b>LeiZhou</b></sub>
        </a>
    </td>
    <td align="center">
        <a href="https://github.com/Yanbo0101">
            <img src="https://avatars.githubusercontent.com/u/110962880?v=4" width="100;" alt="Yanbo0101"/>
            <br />
            <sub><b>Yanbo Xu</b></sub>
        </a>
    </td>
    <td align="center">
        <a href="https://github.com/jialeif">
            <img src="https://avatars.githubusercontent.com/u/88661406?v=4" width="100;" alt="jialeif"/>
            <br />
            <sub><b>Jialei Feng</b></sub>
        </a>
    </td>
    <td align="center">
        <a href="https://github.com/jiere">
            <img src="https://avatars.githubusercontent.com/u/6448681?v=4" width="100;" alt="jiere"/>
            <br />
            <sub><b>Jie Ren</b></sub>
        </a>
    </td></tr>
<tr>
    <td align="center">
        <a href="https://github.com/rdower">
            <img src="https://avatars.githubusercontent.com/u/15023397?v=4" width="100;" alt="rdower"/>
            <br />
            <sub><b>Robert Dower</b></sub>
        </a>
    </td>
    <td align="center">
        <a href="https://github.com/zhlsunshine">
            <img src="https://avatars.githubusercontent.com/u/4101246?v=4" width="100;" alt="zhlsunshine"/>
            <br />
            <sub><b>Steve Zhang</b></sub>
        </a>
    </td>
    <td align="center">
        <a href="https://github.com/wenhuizhang">
            <img src="https://avatars.githubusercontent.com/u/2313277?v=4" width="100;" alt="wenhuizhang"/>
            <br />
            <sub><b>Wenhui Zhang</b></sub>
        </a>
    </td></tr>
</table>
<!-- readme: contributors -end -->

<!-- spell-checker: enable -->
