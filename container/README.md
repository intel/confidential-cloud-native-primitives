# Docker Image Builder

There are several docker image files in the sub directories of current directory. Each sub directory contains Dockerfile of an image. Please find description of each image as follows.

|  Sub directory | Image name  | Description  | 
|---|---|---|
|  ccnp-device-plugin | ccnp-device-plugin  | CCNP Device plugin |
|  ccnp-eventlog-server | ccnp-eventlog-server | Eventlog server |
|  ccnp-measurement-server | ccnp-measurement-server  | Measurement server |
|  ccnp-quote-server | ccnp-quote-server  | Quote server |
|  ccnp-node-measurement-example | ccnp-node-measurement-example  | Example image of getting eventlog and measurement using CCNP SDK |
|  pccs | pccs  | PCCS docker image for Intel® TDX remote attestation. Not required for CCNP usage.|
|  qgs | qgs  | QGS docker image for Intel® TDX remote attestation. Not required for CCNP usage. |


### Build Docker images

`build.sh` is a tool to build all above images and push them to a user-specified image registry. It supports below parameters for different use scenarios.

|  Parameter | Description  | Options  | Default Value  | Required\|Optional  |
|---|---|---|---|---|
|  -a | Action the script will execute: `all` means build and publish images; `build` means build images only; `publish` means publish images; `save` means docker save images to local file.  | build\|publish\|save\|all  | all  | Optional  |
|  -r |  Image registry. Images will be pushed to this image registry. |   |   | Required  |
|  -c | Image name to build. Image names will be the same as sub directory names under directory `container`. By default it will be `all` meaning build all images under `container` | Sub directory name under directory `container`.  |  all |  Optional |
|  -g | Image tag  |   |  latest |  Optional |
|  -f | Flag to build images with parameter "--no-cache"  |   |  |  Optional |
|  -p | Flag to build pccs docker image. Building pccs docker image requires interactive configuration. Please refer to [README.md](../container/pccs/README.md).  |   |   |  Optional |
|  -q | Flag to build qgs docker image. Please refer to [README.md](../container/qgs/README.md).  |  |   |  Optional |
|  -h | Show this usage guide.  |  |   |  Optional |

_NOTE: The script need to run on a server with docker installed._

_NOTE: please set `HTTP_PROXY`, `HTTPS_PROXY`, `NO_PROXY` in docker daemon if they are needed. Please refer to [Configure the Docker daemon to use a proxy server](https://docs.docker.com/config/daemon/systemd/#httphttps-proxy)._

Below are usage examples for different scenarios. Please replace the parameters with your input.

```
# Build all CCNP images with tag 0.3 and push them to remote registry test-registry.intel.com
$ sudo ./build.sh -r test-registry.intel.com/test -g 0.3

# Build images only with tag 0.3
$ sudo ./build.sh -a build -g 0.3

# Build ccnp-measurement-server image with tag 0.3 and push them to remote registry test-registry.intel.com
$ sudo ./build.sh -c ccnp-measurement-server -r test-registry.intel.com/test -g 0.3

# Build pccs image with tag 0.3 and push it to remote registry test-registry.intel.com
$ sudo ./build.sh -c pccs -r test-registry.intel.com/test -g 0.3 -p

# Build qgs image with tag 0.3 and push it to remote registry test-registry.intel.com
$ sudo ./build.sh -c qgs -r test-registry.intel.com/test -g 0.3 -q
```

After the script is running successfully, it's supposed to see corresponding CCNP docker images.

```
$ sudo docker images
ccnp-node-measurement-example   <your image tag>
ccnp-eventlog-server            <your image tag>
ccnp-measurement-server         <your image tag>
ccnp-quote-server               <your image tag>
ccnp-device-plugin              <your image tag>
```
