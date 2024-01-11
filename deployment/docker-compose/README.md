# Docker Compose Deployment

The CCNP can be deployed in the TDVM by docker compose.

## 1. Prerequisite

For deploying the CCNP in the TDX environment, the measurement dependencies should be installed on the host and the guest. This section only emphasizes some key points briefly. If you never use the TDX measurement before, please refer to the section `Measurement & Attesation` in the [White Paper](https://www.intel.com/content/www/us/en/content-details/790888/whitepaper-linux-stacks-for-intel-trust-domain-extensions-1-5.html) to setup it.


### 1.1 Host Side 

The host should enable TDX, and setup service PCCS & service QGS which help generate TDX quote in the guest.

### 1.2 Guest Side

The TDX device in the guest has different names in certain versions
- `/dev/tdx-attest`: the very early name, the version is not supported. 
- `/dev/tdx-guest` or `/dev/tdx_guest`: the name for newer version, it is supported.

To enable the containers read and write the TDX device, change the access privilege of the TDX device.

```
chmod 0666 $(find /dev/ -name "tdx*")
```

This tool will generate files to create some docker composes which are place in `/tmp/docker_ccnp`. Please make sure the directory clear.

For convenience, run the script `prepare.sh` directly as root.

```
sudo ./prepare.sh
```

## 2. Use CCNP

### 2.1 Deploy CCNP

Use the script `deploy-ccnp.sh` to deploy the CCNP services. 

```
./deploy-ccnp.sh
```

In default, the script will launch three containerized services
- eventlog server: from image `ccnp-eventlog-server:latest`
- measurement server: from image `ccnp-measurement-server:latest`
- quote server: from image `quote-server:latest`

Please make sure these container images exist on the guest. If you want to build the images in local, read the [README.md](../../container/README.md) here.

This script has some options.

```
Usage: $(basename "$0") [OPTION]...
    -r <registry prefix>    the prefix string for registry
    -g <tag>                container image tag
    -h                      show help info
```

Specify the registry and tag if necessary.

Typically, some successful message will log as below.

```
INFO: Cache Dir Being Created: /tmp/docker_ccnp
SUCCESS: Cache Dir Created: /tmp/docker_ccnp
INFO: Compose /tmp/docker_ccnp/composes/eventlog-compose.yaml Being Deployed
[+] Running 3/3
 ✔ Network eventlog-server-ctr_default                   Created           0.1s
 ✔ Container eventlog-server-ctr-init-eventlog-server-1  Exited            0.2s
 ✔ Container eventlog-server-ctr-eventlog-server-1       Started           0.1s
SUCCESS: Compose /tmp/docker_ccnp/composes/eventlog-compose.yaml Deployed
INFO: Compose /tmp/docker_ccnp/composes/measurement-compose.yaml Being Deployed
[+] Running 2/2
 ✔ Network measuerment-server-ctr_default                 Created          0.1s
 ✔ Container measuerment-server-ctr-measurement-server-1  Started          0.1s
SUCCESS: Compose /tmp/docker_ccnp/composes/measurement-compose.yaml Deployed
INFO: Compose /tmp/docker_ccnp/composes/quote-compose.yaml Being Deployed
[+] Running 2/2
 ✔ Network quote-server-ctr_default           Created                      0.1s
 ✔ Container quote-server-ctr-quote-server-1  Started                      0.1s
SUCCESS: Compose /tmp/docker_ccnp/composes/quote-compose.yaml Deployed

```

### 2.2 Run Example 

The script `exec-ccnp-example.sh` will help launch container from image `ccnp-node-measurement-example:latest` and request info from above three services. 

```
./exec-ccnp-example.sh -d
```

Typically, some successful message will log as below.

```
INFO: Execute example Container ccnp-node-measurement-example
INFO: Example Container No Avaliable. Attempt Deploy It
[+] Running 3/3
 ✔ Network node-measurement-example-ctr_default                            Created0.1s
 ✔ Container node-measurement-example-ctr-init-node-measurement-example-1  Exited0.2s
 ✔ Container node-measurement-example-ctr-node-measurement-example-1       Started0.1s
SUCCESS: Example Container Avaliable. Compose file: /tmp/docker_ccnp/composes/ccnp-node-measurement-example.yaml
SUCCESS: Measurement Log Saved in File /tmp/docker_ccnp/measurement.log
SUCCESS: Example Container ccnp-node-measurement-example Executed
INFO: Example Container Being Deleted
[+] Running 3/3
 ✔ Container node-measurement-example-ctr-node-measurement-example-1       Removed10.4s
 ✔ Container node-measurement-example-ctr-init-node-measurement-example-1  Removed0.0s
 ✔ Network node-measurement-example-ctr_default                            Removed0.1s
SUCCESS: Example Container Deleted

```

The example container will launch and stop automatically, and save the measurement log in `/tmp/docker_ccnp/measurement.log`. The log file is similar to the [sample file](../../docs/sample-output-for-node-measurement-tool-full.txt).

The script provides some options. 

```
Usage: $(basename "$0") [OPTION]...
    -r <registry prefix>    the prefix string for registry
    -g <tag>                container image tag
    -d			            delete example container
    -o			            request from host
    -h                      show help info
```

Specify the registry and tag if necessary.

The option `-o` will send requests to three service for the guest directly.

The option `-d` indicates that stopping the example after get measurement log.

### 2.3 Clean Up

The script `cleanup.sh` will help stop three containerized services and remove cache.

```
./cleanup.sh
```

Typically, some successful message will log as below.

```
INFO: Compose /tmp/docker_ccnp/composes/eventlog-compose.yaml Being Down
[+] Running 3/3
 ✔ Container eventlog-server-ctr-eventlog-server-1       Removed           0.4s
 ✔ Container eventlog-server-ctr-init-eventlog-server-1  Removed           0.0s
 ✔ Network eventlog-server-ctr_default                   Removed           0.2s
SUCCESS: Compose /tmp/docker_ccnp/composes/eventlog-compose.yaml Down
INFO: Compose /tmp/docker_ccnp/composes/measurement-compose.yaml Being Down
[+] Running 2/2
 ✔ Container measuerment-server-ctr-measurement-server-1  Removed          0.4s
 ✔ Network measuerment-server-ctr_default                 Removed          0.2s
SUCCESS: Compose /tmp/docker_ccnp/composes/measurement-compose.yaml Down
INFO: Compose /tmp/docker_ccnp/composes/quote-compose.yaml Being Down
[+] Running 2/2
 ✔ Container quote-server-ctr-quote-server-1  Removed                     10.3s
 ✔ Network quote-server-ctr_default           Removed                      0.2s
SUCCESS: Compose /tmp/docker_ccnp/composes/quote-compose.yaml Down
INFO: Cache Dir Being Removed
SUCCESS: Cache Dir Removed

```

