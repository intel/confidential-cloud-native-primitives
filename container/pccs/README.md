# Deploy PCCS service on Docker

PCCS (Provisioning Certificate Caching Service) service implementation comes from
[DCAP](https://github.com/intel/SGXDataCenterAttestationPrimitives/blob/master/QuoteGeneration/pccs/README.md).
Currently, the package of PCCS only support several distros. Using docker to deploy the PCCS service can be an alternative for some unsupported distros, like RHEL9.


## 1. PCCS Service Usage Guide

### 1.1 Setup PCCS Configuration

When building the PCCS docker image will be required to setup the configuration of PCCS.

1. Please prepare an Intel PCS API key in advance. If you do not have an API key for the Intel SGX PCS service, obtain a key: Go to https://api.portal.trustedservices.intel.com/provisioning-certification login (or create an account), and click 'Subscribe'. Two API keys are generated (for key rollover). Copy one of the API keys for the configuration steps.

2. Install PCCS with following steps. During installation, answer “Y” when asked if the PCCS should be installed now, “Y” when asked if PCCS should be configured now, and enter API key generated in step 1 when asked for “Intel PCS API key”. Answer the remaining questions according to your needs.

### 1.2 Start PCCS Service

Note: Configure the restart policy to always,which makes PCCS service to keep running after server reboot.

```bash
docker run -d --privileged -v /sys/firmware/efi/:/sys/firmware/efi/ --name pccs --restart always --net host <your registry>
```

- Check if PCCS service works

```console
$ docker ps
CONTAINER ID   IMAGE      COMMAND                 CREATED         STATUS         PORTS      NAMES
90a3777d813e   pccs       "node pccs_server.js"   9 minutes ago   Up 9 minutes   8081/tcp   pccs
```

### 1.3 Register SGX Platform

PCKIDRetrieval tool has already integrated into the PCCS docker image. Therefore, after pccs is activated, registration can be triggered directly.

```bash
docker exec -it pccs /opt/intel/sgx-pck-id-retrieval-tool/PCKIDRetrievalTool
```

### 1.4 Check PCCS Service Log

Debug the pccs service, when registration failed.

```bash
docker logs pccs
```

**NOTE: If you see message about "Platform Manifest not available" or the PCCS service complains that "Error: No cache data for this platform" in the log, you may need to perform SGX Factory Reset in BIOS and run PCKIDRetrievalTool again.**

## 2. Optional advanced operations

If the pccs docker need to be removed (Normally, we do not recommend remove the pccs docker.), please copy out the **pckcache.db** firstly.

```bash
docker cp pccs:/opt/intel/sgx-dcap-pccs/pckcache.db .
```

The **pckcache.db** can be reused when restarts pccs service in the next time.

```bash
docker run -d --privileged -v /sys/firmware/efi/:/sys/firmware/efi/ -v /path/to/pckcache.db:/opt/intel/sgx-dcap-pccs/pckcache.db --name pccs --restart always --net host <your registry>
```
