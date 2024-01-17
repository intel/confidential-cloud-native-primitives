# Deploy QGS service on Docker

QGS (Quote Generation Service) implementation comes from
[DCAP](https://github.com/intel/SGXDataCenterAttestationPrimitives/tree/master/QuoteGeneration/quote_wrapper/qgs).
Currently, the package of QGS only support several distros. Using docker to deploy the QGS service can be an alternative for some unsupported distros, like RHEL9.

## 1. QGS Service Usage Guide

### 1.1 Start QGS Service

```bash
docker run -d --privileged --name qgs --restart always --net host <your registry>
```
- Check if QGS service works

```console
$ docker ps
CONTAINER ID   IMAGE      COMMAND                 CREATED         STATUS         PORTS      NAMES
90a3777d813e   qgs        "/opt/intel/tdx-qgs/…"  9 minutes ago   Up 9 minutes              qgs
```
