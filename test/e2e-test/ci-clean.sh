#!/bin/bash
set -o errexit

CLUSTER_NAME=my-cluster
REG_NAME=kind-registry
DEVICE_PLUGIN=localhost:5001/ccnp-device-plugin:0.1
QUOTE=localhost:5001/ccnp-quote-server:0.1
MEASUREMENT=localhost:5001/ccnp-measurement-server:0.1
EVENTLOG=localhost:5001/ccnp-eventlog-server:0.1

kind delete cluster --name $CLUSTER_NAME
rm /run/ccnp/uds/*
docker stop $REG_NAME
docker rm $REG_NAME
sleep 1m
docker rmi ${DEVICE_PLUGIN} ${QUOTE} ${MEASUREMENT} ${EVENTLOG}
for i in $(docker images | grep "none" | awk '{print $3}');do
  docker rmi "$i"
done
pip uninstall -y ccnp
