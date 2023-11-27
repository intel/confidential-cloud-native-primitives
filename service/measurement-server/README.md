# Service: CCNP Measurement Server

Measurements are key assets in the confidential computing world for user to verify and determine the trustworthiness of the environment and its underlying platform.
Users can utilize different kinds of measurements to monitor or validate the operations performed within the platform.

This service (which is actually a GRPC server) will provide such interface to fetch the measurements for confidential cloud native environments, including TDX RTMR measurements, TPM PCR measurements and TEE reports.



## Introduction

This service provides functionality to fetch measurements for confidential cloud native environments, both platform level(`PAAS` option) and container level (`SAAS` option). Using this service, user can fetch different
categories of measurements including: TDX RTMR measurements (`TDX_RTMR` option), TPM PCR measurements (`TPM` option) and TEE reports from different vendors (`TEE_REPORT` option).
Here shows the proto buf for the service:

```
enum TYPE {
    PAAS = 0;
    SAAS = 1;
}

enum CATEGORY {
    TEE_REPORT = 0;
    TPM = 1;
    TDX_RTMR = 2;
}

message GetMeasurementRequest {
    TYPE measurement_type = 1;
    CATEGORY measurement_category = 2;
    string report_data = 3;
    int32 register_index = 4;

}

message GetMeasurementReply {
    string measurement = 1;
}

service Measurement {
    rpc GetMeasurement (GetMeasurementRequest) returns (GetMeasurementReply) {}
}
```

To select different level and category of measurements, user shall send out request with different settings.
User can find sample request in the Testing section.

The collected measurements are returned as json string to the client.



## Installation

The Measurement service can be deployed as either DaemonSet or sidecar in different user scenarios within a kubernetes cluster.


### Prerequisite

User need to have a kubernetes cluster ready to deploy the services. To simplify the deployment process, we provide Helm as one of the options to deploy the service. Please install Helm by following the [Helm official guide](https://helm.sh/docs/intro/install/). However, user can also use the yaml files located in the manifests folder for deployment.
Also, the ccnp device plugin need to installed before the installation of measurement server. Please refer to its [deployment guide](../../device-plugin/ccnp-device-plugin/README.md) for installation.

### Build docker image

The dockerfile for the service can be found under `container/measurement-server` directory. Use the following command to build the image:

```
cd ../..
docker build -t ccnp-measurement-server:0.1 -f container/measurement-server/Dockerfile .
```
> Note: if you are using containerd as the default runtime for kubernetes. Please remember to use the following commands to import the image into containerd first:
```
docker save -o ccnp-measurement-server.tar ccnp-measurement-server:0.1
ctr -n=k8s.io image import ccnp-measurement-server.tar
```

### Deploy as DaemonSet

In the scenario of confidential kubernetes cluster, it is nice to deploy the measurement service as a DaemonSet to serve all the applications living inside that cluster node.
Run the following command to deploy the service using helm chart:

```
cd ../../deployment
helm install charts/measurement-server --generate-name
```
> Note: `ccnp` namespace may get duplicated in the case user installed multiple ccnp services. Please define the `namespaceCreate` option as `false` in the values.yaml before install the helm chart if certain case happens.

User can also choose the manifests file to deploy the service:
```
cd ../../deployment/manifests
kubectl create -f namespace.yaml
kubectl create -f measurement-server-deployment.yaml
```

### Deploy as Sidecar

In the scenario of confidential containers, it is nice to deploy the measurement service as sidecar working along with the confidential containers.
The deployment helm chart and manifests are still in progress.

### User application deployment

Make sure that the user application requests such resource in the yaml file to take use of sdk to contact the measurement server: `tdx.intel.com/tdx-guest: 1`.



## Testing

User can play with the service on host from the source code by following the steps below:

1. Start the measurement service

```
cd service/measurement-server
make all

./measurement-server
```

2. Play with the service

Use the `grpcurl` as the tool to play with the service. Please follow the [official documentation](https://github.com/fullstorydev/grpcurl) to install grpcurl.

Get TDX RTMR measurement of register index equal to 0:
```
grpcurl -plaintext -d '{"measurement_type": 0, "measurement_category": 2, "register_index": 0}' -unix /run/ccnp/uds/measurement.sock measurement.Measurement/GetMeasurement
```

Get TEE report according to the platform capability:
```
grpcurl -plaintext -d '{"measurement_type": 0, "measurement_category": 0}' -unix /run/ccnp/uds/measurement.sock measurement.Measurement/GetMeasurement
```

User can find the fetched measurements as base64 encoded string returned as response.


