__CCNP Installation Guide__
===

Confidential Cloud-Native Primitives (CCNP) uses cloud native design pattern to implement confidential computing primitives like event log, measurement, quote and attestation. The project itself contains several parts: the services, the SDK and related dependencies.This is a step-by-step guide for installing CCNP based on ubuntu, including creating TDVMs, creating a K8S cluster, building and installing CCNP. For more detailed information, please refer to [main page](https://github.com/intel/confidential-cloud-native-primitives).

![Architecture diagram](https://github.com/intel/confidential-cloud-native-primitives/blob/main/docs/ccnp_arch.png)

# Create TDVMs

 CCNP needs to run on TDVMs. So, users need to start at least one TDVM on the host server. Please set up the host environment and start TDVMs according to the [TDX white paper](https://www.intel.com/content/www/us/en/content-details/787041/whitepaper-linux-stacks-for-intel-trust-domain-extension-1-0.html).  
It is recommended to start TDVMs with Libvirt according to the  [TDX white paper](https://www.intel.com/content/www/us/en/content-details/787041/whitepaper-linux-stacks-for-intel-trust-domain-extension-1-0.html) and users should add the following configuration to the default xml template file.
```
<launchSecurity type='tdx'>
......
<Quote-Generation-Service>vsock:2:4050</Quote-Generation-Service>
</launchSecurity>
```

NOTE: The disk space of the TDVM should be greater than 40GB in order to build the images of CCNP.

# Create a K8S cluster
After TDVMs are started, users need to setup a K8S cluster in the TDVMs. Please refer to the [k8s official documentation](https://kubernetes.io/docs/home/) for detailed steps. Next, please install helm, a common K8S package management tool. Helm will be used to install CCNP in the next step. For detailed installation steps, please refer to the [helm official documentation](https://helm.sh/).  
NOTE: If the cluster has only one node (master node), the taint on the node needs to be removed.
## Set up proxy (optional)
Whether to set up the proxy depends on the user's network. If users want to set up a network proxy for TDVM, users need to set up no_proxy at the same time. In addition, users need to add the proxy and no_proxy to the configuration files of docker and containerd so that the K8S cluster can work fine. The following commands are for reference only.  
Set up proxy and no_proxy.
```
# 10.244.0.0/16 is the default pod IP for Flannel
# 10.96.0.0/12 is the K8S default service ip
# <TDVM ip segment> refers to the range of IP addresses used by TDVM, such as 10.10.0.0/16
export http_proxy=<Actual network proxy>
export https_proxy=<Actual network proxy>
export HTTP_PROXY=<Actual network proxy>
export HTTPS_PROXY=<Actual network proxy>
export no_proxy=$no_proxy,10.244.0.0/16,10.96.0.0/12,<TDVM ip segment> 
export NO_PROXY=$NO_PROXY,10.244.0.0/16,10.96.0.0/12,<TDVM ip segment>
```
Add the proxy and no_proxy to the configuration file of docker.
```
sudo mkdir -p /etc/systemd/system/docker.service.d
cat <<EOF | sudo tee /etc/systemd/system/docker.service.d/proxy.conf
[Service]
Environment="HTTP_PROXY=$HTTP_PROXY"
Environment="HTTPS_PROXY=$HTTPS_PROXY"
Environment="NO_PROXY=$NO_PROXY"
EOF

# Make the configuration take effect
sudo systemctl enable docker
sudo systemctl daemon-reload
sudo systemctl restart docker
```

 Edit the file "/lib/systemd/system/containerd.service" on each TDVM and add the following command under the field "[Service]" in the file.


```
sudo vi /lib/systemd/system/containerd.service

# Add the following command under the field "[Service]" in the file
Environment="HTTP_PROXY=<actual proxy>"
Environment="HTTPS_PROXY=<actual proxy>"
Environment="NO_PROXY=<actual noproxy>"

# Make the configuration take effect
sudo systemctl daemon-reload
sudo systemctl restart containerd.service
```


# Install CCNP

CCNP provides two scripts, __build-CCNP.sh__ and __deploy-CCNP.sh__ to quickly build images and deploy CCNP services. These scripts are [here](https://github.com/intel/confidential-cloud-native-primitives/tree/main/deployment). Please follow the steps below to install CCNP.

## Build images

Users can run the script, build-CCNP.sh, to quickly build images. The script obtains the environment variables HTTP_PROXY and HTTPS_PROXY from the command line. Therefore, please set up the proxy in advance or skip this if the proxy is not needed. Please run the script on the master node which will download the latest CCNP code and build the images.
```
./build-CCNP.sh
```
The script may take about 20 minutes, please be patient. And it will generate four docker images: __ccnp_eventlog_server__, __ccnp_measurement_server__, __ccnp-quote-server__ __and__ __ccnp-device-plugin__. You can use the command "docker images" to check them.

## Upload images

In order for each node to pull these images, users need:

1. Create a private image repository. You can create an internal private image repository such as registry, harbor.  And you can also use an external repository such as Docker Hub.

2. On the master node, add the private repository to the docker configuration file. And on all nodes, add the private repository to the containerd configuration file. Because users generate these images in docker on the master node and use these images in containerd on all nodes.

3. Tag these images (ccnp_eventlog_server, ccnp_measurement_server, ccnp-quote-server and ccnp-device-plugin) as private repository and upload them.

Take the example of creating an internal and password free registry repository.
```
# Create a registry repo on the master node
docker run --name=reqistry --restart=always -d -p 5000:5000 -v /opt/data/registry/:/var/lib/registry registry

# Add the private repository to the docker configuration file on the master node
# Add the following fields to the file /etc/docker/daemon.json
"insecure-registries": ["<actual ip>:<port>"],

# Make the configuration take effect
sudo systemctl daemon-reload
sudo systemctl restart docker

# Add the private repository to the containerd configuration file on all nodes
sudo mkdir -vp /etc/containerd/
sudo touch /etc/containerd/config.toml
sudo bash -c 'containerd config default > /etc/containerd/config.toml'
# Add the following fields to the file /etc/containerd/config.toml
[plugins."io.containerd.grpc.v1.cri".registry.mirrors]
  [plugins."io.containerd.grpc.v1.cri".registry.mirrors."<actual ip>:<port>"]
    endpoint = ["http://<actual ip>:<port>"]
# Make the configuration take effect
sudo systemctl daemon-reload
sudo systemctl restart containerd.service

# Tag these images as private repository and upload them
docker tag ccnp_eventlog_server:0.1 <actual ip>:<port>/ccnp_eventlog_server:0.1
docker push <actual ip>:<port>/ccnp_eventlog_server:0.1
```
## Deploy CCNP services and dependencies

Please run the following command on each node to set up installation environment.
```
sudo mkdir -p /etc/udev/rules.d
sudo touch /etc/udev/rules.d/90-tdx.rules
sudo bash -c 'echo "SUBSYSTEM==\"misc\",KERNEL==\"tdx-guest\",MODE=\"0666\"">/etc/udev/rules.d/90-tdx.rules'
sudo udevadm trigger
sudo mkdir -p  /run/ccnp/uds
sudo chmod o+w /run/ccnp/uds
```

Users can run the script, deploy-CCNP.sh, to quickly deploy CCNP services and dependencies. This script needs to be run in the same directory as the build-CCNP.sh script. The script obtains the environment variable PRIVATE_REPO from the command line and updates the images’ name in the YAML files, so that the containerd can pull the images from the private repository.
Please set up the environment variable PRIVATE_REPO and run the script on the master node. If the environment variable is not set, the script will use the default repo "docker.io/library".


```
# Please confirm that helm has been installed
export PRIVATE_REPO=<actual private repo>
./deploy-CCNP.sh
```

## Install CCNP SDK

There are two options to install CCNP SDK.  
Option 1: Users can run the following command on each work node to install CCNP client library for Python with PyPI.
```
pip install ccnp
```
Option 2: Users can also run the following commands to install CCNP SDK from source code.
```
cd confidential-cloud-native-primitives/sdk/python3/
pip install -r requirements.txt
pip install -e .
```
At this step, CCNP has been installed successfully. For more detailed information, including SDK usage, please refer to [here](https://intel.github.io/confidential-cloud-native-primitives/).

