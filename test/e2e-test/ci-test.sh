#!/bin/bash
set -o errexit
: '
This is an E2E test script, if you want to run the script please:
1 According to the CCNP documentation, create a TDVM
2 Install helm and kind on the TDVM
3 Set up the environment according to the CCNP documentation,including configuring the file "/etc/udev/rules.d/90-tdx.rules" and creating the folder "/run/ccnp/uds" on the TDVM
4 Git clone CCNP and run this script
5 Clean up the environment
'

CLUSTER_NAME=my-cluster
KIND_CONFIG=kind-config.yaml
DEVICE_PLUGIN=localhost:5001/ccnp-device-plugin:0.1
QUOTE=localhost:5001/ccnp-quote-server:0.1
MEASUREMENT=localhost:5001/ccnp_measurement_server:0.1
EVENTLOG=localhost:5001/ccnp_eventlog_server:0.1
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



#Set up the environment
export NO_PROXY=$NO_PROXY,$REG_NAME

#Create a private repository
if [ "$(docker inspect -f '{{.State.Running}}' "${REG_NAME}" 2>/dev/null || true)" != 'true' ]; then
  docker run \
    -d --restart=always -p "127.0.0.1:${REG_PORT}:5000" --name "${REG_NAME}" \
    registry:2
fi

#Initialize a kind-K8S cluster and add the private repo to the configuration file
kind create cluster --name $CLUSTER_NAME  --config $KIND_CONFIG

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
kubectl apply -f $REPO_CONFIGMAP



#Build and push images

pushd "${WORK_DIR}/../.."
docker build --build-arg http_proxy="$HTTP_PROXY" --build-arg https_proxy="$HTTPS_PROXY" \
        --build-arg no_proxy="$NO_PROXY" -t $QUOTE -f container/quote-server/Dockerfile .
docker build --build-arg http_proxy="$HTTP_PROXY" --build-arg https_proxy="$HTTPS_PROXY" \
        --build-arg no_proxy="$NO_PROXY" -t $MEASUREMENT -f container/measurement-server/Dockerfile .
docker build --build-arg http_proxy="$HTTP_PROXY" --build-arg https_proxy="$HTTPS_PROXY" \
        --build-arg no_proxy="$NO_PROXY" -t $EVENTLOG -f container/eventlog-server/Dockerfile .
pushd  device-plugin/ccnp-device-plugin/
docker build  --build-arg http_proxy="$HTTP_PROXY" --build-arg https_proxy="$HTTPS_PROXY" \
        --build-arg no_proxy="$NO_PROXY" -t $DEVICE_PLUGIN -f container/Dockerfile .
popd


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

#Install SDK
pushd sdk/python3/
pip install -r requirements.txt
pip install -e .
pip install pytest pytdxattest
popd
popd
sleep 2m

#Run test cases
pytest  test_eventlog.py test_tdquote.py  test_tdreport.py

