#!/bin/bash
set -o errexit
: '
This is a script for setting up the CCNP ci environment. If you want to run the script please:
1 According to the CCNP documentation, create a TDVM
2 Install helm and kind on the TDVM
3 Follow the CCNP documentation to pre-configure the TDVM,including modifying the file "/etc/udev/rules.d/90-tdx.rules" and creating the folder "/run/ccnp/uds"
4 Git clone CCNP and run this script
'

CLUSTER_NAME=my-cluster
KIND_CONFIG=kind-config.yaml
DEVICE_PLUGIN=localhost:5001/ccnp-device-plugin:0.1
QUOTE=localhost:5001/ccnp-quote-server:latest
MEASUREMENT=localhost:5001/ccnp-measurement-server:latest
EVENTLOG=localhost:5001/ccnp-eventlog-server:latest
REG_NAME=kind-registry
REG_PORT=5001
REPO_CONFIGMAP=repo-configmap.yaml
LOCAL_REPO=localhost:5001
DOCKER_REPO=docker.io/library
NFD_URL=https://kubernetes-sigs.github.io/node-feature-discovery/charts
NFD_NS="node-feature-discovery"
WORK_DIR=$(cd "$(dirname "$0")"; pwd)

change_value(){
        sed -i  "s#${1}#${2}#g" "$3"
}

create_dir(){
        if [ ! -d "$1" ];then
                mkdir -p "$1"
        fi
}

create_file(){
        if [ ! -e "$1" ];then
                 touch  "$1"
        fi

}



#Set up the NO_PROXY
export NO_PROXY=$NO_PROXY,$REG_NAME

#Create a private repository
if [ "$(docker inspect -f '{{.State.Running}}' "${REG_NAME}" 2>/dev/null || true)" != 'true' ]; then
  docker run \
    -d --restart=always -p "127.0.0.1:${REG_PORT}:5000" --name "${REG_NAME}" \
    registry:2
fi

#Initialize a kind-K8S cluster and add the private repo to the configuration file
kind create cluster --name $CLUSTER_NAME  --config "${WORK_DIR}/${KIND_CONFIG}"

REGISTRY_DIR="/etc/containerd/certs.d/localhost:${REG_PORT}"
for node in $(kind get nodes --name $CLUSTER_NAME); do
  docker exec "${node}" mkdir -p "${REGISTRY_DIR}"
  cat <<EOF | docker exec -i "${node}" cp /dev/stdin "${REGISTRY_DIR}/hosts.toml"
[host."http://${REG_NAME}:5000"]
EOF
done
if [ "$(docker inspect -f='{{json .NetworkSettings.Networks.kind}}' "${REG_NAME}")" = 'null' ]; then
  docker network connect "kind" "${REG_NAME}"
fi
kubectl apply -f "${WORK_DIR}/$REPO_CONFIGMAP"



#Build and push images
pushd "${WORK_DIR}/../.."
docker build --build-arg http_proxy="$HTTP_PROXY" --build-arg https_proxy="$HTTPS_PROXY" \
        --build-arg no_proxy="$NO_PROXY" -t $QUOTE -f container/ccnp-quote-server/Dockerfile .
docker build --build-arg http_proxy="$HTTP_PROXY" --build-arg https_proxy="$HTTPS_PROXY" \
        --build-arg no_proxy="$NO_PROXY" -t $MEASUREMENT -f container/ccnp-measurement-server/Dockerfile .
docker build --build-arg http_proxy="$HTTP_PROXY" --build-arg https_proxy="$HTTPS_PROXY" \
        --build-arg no_proxy="$NO_PROXY" -t $EVENTLOG -f container/ccnp-eventlog-server/Dockerfile .
docker build  --build-arg http_proxy="$HTTP_PROXY" --build-arg https_proxy="$HTTPS_PROXY" \
        --build-arg no_proxy="$NO_PROXY" -t $DEVICE_PLUGIN -f container/ccnp-device-plugin/Dockerfile .


docker push ${DEVICE_PLUGIN}
docker push ${QUOTE}
docker push ${MEASUREMENT}
docker push ${EVENTLOG}

sed -i  "s#${DOCKER_REPO}#${LOCAL_REPO}#g" deployment/manifests/*
sed -i  "s#${DOCKER_REPO}#${LOCAL_REPO}#g" device-plugin/ccnp-device-plugin/deploy/helm/ccnp-device-plugin/values.yaml

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

#Wait for all pods and services to be ready
sleep 2m

#Install SDK
pushd sdk/python3/
pip install -r requirements.txt
pip install .
popd
popd

