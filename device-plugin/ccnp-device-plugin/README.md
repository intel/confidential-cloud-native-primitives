# CCNP device plugin

The CCNP device plugin is based on Kubernetes plugin framework to expose host's TEE devices and other required resources to PODs.
This will enable the services in the PODs to be able to communicate to the device for quote, measurements etc.
And will also facilitate other CCNP requirements like mount certain directories to CCNP service PODs and workload PODs.

## Introduction

Currently, the CCNP device plugin has following capabilities:
- expose tdx guest device node in TDVM to PODs
- mount Unix Domain Socket dir /run/ccnp/uds into CCNP daemonset PODs and workload PODs to enable intra-node communication

The CCNP device plugin will respond to following resource request defined in POD definition yaml:
```
    resources:
      limits:
        tdx.intel.com/tdx-guest: 1    
```

## Installation

### Prerequisite
The CCNP device plugin need to deploy on VM nodes with guest TEE devices(currently only TDX guest device is supported). So the deployment
of the plugin daemonset is based on the node label set by [Node Feature Discovery](https://github.com/kubernetes-sigs/node-feature-discovery/).
So we need to install the NFD and corresponding label rules.

1. setup following udev rule to enable other user in the node to read or write to tdx guest device node
```
cat /etc/udev/rules.d/90-tdx.rules
SUBSYSTEM=="misc",KERNEL=="tdx-guest",MODE="0666"

```
After adding the rule, you can restart the node or run following command to trigger the update:
```
udevadm trigger
```

2. prepare the shared Unix Domain Socket directory to be mounted to both ccnp service pods and workload pods
```
mkdir -p /run/ccnp/uds
chmod o+w /run/ccnp/uds
```

3. deploy NFD

before NFD v0.14 release is ready, please build own image and deploy:
> Note: please change the imagePullPolicy to 'IfNotPresent' in kustomization file(worker-daemonset and master) at deployment/base if you build your image locally.
```
git clone https://github.com/kubernetes-sigs/node-feature-discovery.git
cd node-feature-discovery/
make image
kubectl apply -k .
```

> Note: when node-feature-discovery new [release v0.14](https://github.com/kubernetes-sigs/node-feature-discovery/issues/1250) is ready, below command can be used to deploy NFD with TDVM support:

```
kubectl apply -k https://github.com/kubernetes-sigs/node-feature-discovery/deployment/overlays/default?ref=v0.14
```

4. deploy NFD label rules
```
kubectl apply -f device-plugin/ccnp-device-plugin/deploy/node-feature-rules.yaml
```

After deployment, following label can be found in the VM node:
```
kubectl get node -o json | jq .items[].metadata.labels | grep tdx-guest
  "intel.feature.node.kubernetes.io/tdx-guest": "enabled",
```
Above label can be used as node selector by CCNP device plugin daemonset and CCNP services daemonset.


### Build docker image
The Dockerfile for the service can be found under device-plugin/ccnp-device-plugin/container directory. 
Use the following command to build the image:
```
cd device-plugin/ccnp-device-plugin/
docker build -t ccnp-device-plugin:0.1 -f container/Dockerfile .
```

> Note: if you are using containerd as the default runtime for kubernetes, don't forget to use the following commands to import the image into containerd first:
```
docker save -o ccnp-device-plugin.tar ccnp-device-plugin:0.1
ctr -n=k8s.io image import ccnp-device-plugin.tar
```

### Deploy as DaemonSet
Use below helm command to deploy:
> Note: you may need to edit settings in helm [value.yaml](deploy/helm/ccnp-device-plugin/value.yaml) according to you cluster status.
```
cd device-plugin/ccnp-device-plugin/
helm install ccnp-device-plugin deploy/helm/ccnp-device-plugin

```

After the deployment, for TDVM node, you can see below resource info:
```
kubectl describe node 
...
Capacity:
  cpu:                      8
...
  memory:                   7687708Ki
  pods:                     110
  tdx.intel.com/tdx-guest:  110
Allocatable:
  cpu:                      8
...
  memory:                   7585308Ki
  pods:                     110
  tdx.intel.com/tdx-guest:  110
...
Allocated resources:
  (Total limits may be over 100 percent, i.e., over committed.)
  Resource                 Requests     Limits
  --------                 --------     ------
  cpu                      1250m (15%)  600m (7%)
  memory                   510Mi (6%)   690Mi (9%)
...
  tdx.intel.com/tdx-guest  0            0
...
```

### Testing
User can deploy a CCNP quote service with tdx-guest resource request in the DaemonSet definition yaml:
```
...
        resources:
          limits:
            tdx.intel.com/tdx-guest: 1
...
```

And after the quote server POD is started, following resource and directory can be found in the container of the POD:
```
ls -l /dev/tdx-guest
crw-rw-rw- 1 root root 10, 126 Jul 12 04:58 /dev/tdx-guest

ls -l /run/ccnp/uds
total 0
srwxr-xr-x 1 ccnp ccnp 0 Jul 12 04:58 quote-server.sock
```