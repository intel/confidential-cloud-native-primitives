# Service: CCNP Eventlog Server

To further verify the integrity and authenticity of the measurements in the confidential cloud native environment and its underlying platform, event logs are absolutely needed.
By reviewing the event logs, user can identify any errors or issues that may be preventing the confidential environment from functioning correctly.

This service (which is actually a GRPC server) will provide support to fetch the event logs for confidential cloud native environments, including TDX RTMR event logs and TPM event logs.



## Introduction

This service provides functionality to fetch event logs for confidential cloud native environments, both platform level(`PAAS` option) and container level (`SAAS` option). Using this service, user can fetch event logs for both TDX RTMR (`TDX_EVENTLOG` option) and TPM (`TPM_EVENTLOG` option).
Here shows the proto buf for the service:

```
enum CATEGORY {
    TDX_EVENTLOG = 0;
    TPM_EVENTLOG = 1;
}

enum LEVEL {
    PAAS = 0;
    SAAS = 1;
}

message GetEventlogRequest {
    LEVEL eventlog_level = 1;
    CATEGORY eventlog_category = 2;
    int32 start_position = 3;
    int32 count = 4;
}

message GetEventlogReply {
    string eventlog_data_loc = 1;
}

service Eventlog {
    rpc GetEventlog (GetEventlogRequest) returns (GetEventlogReply) {}
}
```

To select different level and category of event log, user shall send out request with different settings.
The service also supports fetching number of event logs start from a certain one. The `start_position` option and `count` option are provided for the usage.
User can find sample request in the Testing section.

Request to the service will return the location of the collected event logs. They are stored under the folder `/run/ccnp-eventlog` by default.



## Installation

The Eventlog server service can be deployed as either DaemonSet or sidecar in different user scenarios within a kubernetes cluster.


### Prerequisite

User need to have a kubernetes cluster ready to deploy the service. To simplify the deployment process, we provide Helm as one of the options for deployment. Please install Helm by following the [Helm official guide](https://helm.sh/docs/intro/install/). However, user can still install the service using the yaml file located in the manifests folder.

### Build docker image

The dockerfile for the service can be found under `container/eventlog-server` directory. Use the following command to build the image:

```
cd ../..
docker build -t ccnp-eventlog-server:0.1 -f container/eventlog-server/Dockerfile .
```
> Note: if you are using containerd as the default runtime for kubernetes. Please remember to use the following commands to import the image into containerd first:
```
docker save -o ccnp-eventlog-server.tar ccnp-eventlog-server:0.1
ctr -n=k8s.io image import ccnp-eventlog-server.tar
```

### Deploy as DaemonSet

In the scenario of confidential kubernetes cluster, it is nice to deploy the Eventlog server service as a DaemonSet to serve all the applications living inside that cluster node.
Run the following command to deploy the service:

```
cd ../../deployment
helm install charts/eventlog-server --generate-name
```
> Note: `ccnp` namespace may get duplicated in the case user installed multiple ccnp services. Please define the `namespaceCreate` option as `false` in the values.yaml before install the helm chart if certain case happens.

User can also choose the manifests file to deploy the service:
```
cd ../../deployment/manifests
kubectl create -f namespace.yaml
kubectl create -f eventlog-server-deployment.yaml
```

### Deploy as Sidecar

In the scenario of confidential containers, it is nice to deploy the Eventlog server service as sidecar working along with the confidential containers.
The deployment helm chart and manifest file are still working in progress.

### User application deployment

Make sure that the user application mount the same event log directory into the container, which defined by `eventlogDir` variable within the helm chart values.yaml.
So that it can get the event log fetched and use according to their usage.



## Testing

User can play with service on host from the source code by following the steps below:

1. Start the eventlog service

```
cd service/eventlog-server
make all

./eventlog-server
```

2. Play with the service

Use the `grpcurl` as the tool to play with the service. Please follow the [official documentation](https://github.com/fullstorydev/grpcurl) to install grpcurl

Get all TDX RTMR event logs from the platform level:
```
grpcurl -plaintext -d '{"eventlog_level": 0, "eventlog_category": 0}' -unix /run/ccnp/uds/eventlog.sock Eventlog/GetEventlog
```

Get 5 TDX RTMR event logs starting from the second one from platform level:
```
grpcurl -plaintext -d '{"eventlog_level": 0, "eventlog_category": 0, "start_position": 2, "count": 5}' -unix /run/ccnp/uds/eventlog.sock Eventlog/GetEventlog
```

User can find the fetched event logs under the mounted directory.

