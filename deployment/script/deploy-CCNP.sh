#!/bin/bash

DOCKER_REPO=docker.io/library
NFD_NS=node-feature-discovery
NFD_URL=https://kubernetes-sigs.github.io/node-feature-discovery/charts
WORK_DIR=$(cd "$(dirname "$0")" || exit; pwd)

pushd "${WORK_DIR}/../.." || exit

#If private repo is used, modify the images' names in the yaml files
if [ -n "$PRIVATE_REPO" ]; then
sed -i  "s#${DOCKER_REPO}#${PRIVATE_REPO}#g" deployment/manifests/*
sed -i  "s#${DOCKER_REPO}#${PRIVATE_REPO}#g" device-plugin/ccnp-device-plugin/deploy/helm/ccnp-device-plugin/values.yaml
fi

#Check if "helm" is installed
if dpkg -s "helm" &> /dev/null; then
    echo "helm is installed"
else
    echo "Please install helm"
    exit 2
fi

#Deploy CCNP Dependencies
helm repo add nfd $NFD_URL
helm repo update
helm install nfd/node-feature-discovery --namespace $NFD_NS --create-namespace --generate-name
kubectl apply -f  device-plugin/ccnp-device-plugin/deploy/node-feature-rules.yaml
helm install ccnp-device-plugin  device-plugin/ccnp-device-plugin/deploy/helm/ccnp-device-plugin

#Deploy CCNP services
kubectl create -f deployment/manifests/namespace.yaml
kubectl create -f deployment/manifests/eventlog-server-deployment.yaml
kubectl create -f deployment/manifests/measurement-server-deployment.yaml
kubectl create -f deployment/manifests/quote-server-deployment.yaml
popd || exit
