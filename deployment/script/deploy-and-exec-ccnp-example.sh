#!/bin/bash
# Script to deploy and execute CCNP example container
# Attach the RTMR index after the script during execution to verify selected register

set -e

DEFAULT_DOCKER_REPO=docker.io/library
DEFAULT_TAG=latest
TEMP_MANIFEST_FILE=/tmp/ccnp-node-measurement-example-deployment.yaml
DELETE_DEPLOYMENT=false

usage() { echo "Usage: $0 [-r <registry-prefix>] [-g <image-tag>] [-d delete existing deployment] [-i <register-index-to-verify>]"; exit 1; }
while getopts ":r:g:i:dh" option; do
        case "${option}" in
            r) registry=${OPTARG};;
            g) tag=${OPTARG};;
            i) index=${OPTARG};;
            d) DELETE_DEPLOYMENT=true;;
            h) usage;;
            *) echo "Invalid option: -${OPTARG}" >&2
               usage
               ;;
        esac
    done

echo "Step 1:  Deploy CCNP example container for node measurement in Kubernetes"
# replace registry and image tag according to user input
cp ../manifests/ccnp-node-measurement-example-deployment.yaml $TEMP_MANIFEST_FILE
if [[ -n "$registry" ]]; then
	sed -i  "s#${DEFAULT_DOCKER_REPO}#${registry}#g" $TEMP_MANIFEST_FILE
fi
if [[ -n "$tag" ]];then
	sed -i "s#${DEFAULT_TAG}#${tag}#g" $TEMP_MANIFEST_FILE
fi

if [ $DELETE_DEPLOYMENT == true ]
then
    echo "==> Cleaning up ccnp-node-measurement-example deployment"
    kubectl delete -f $TEMP_MANIFEST_FILE
fi

echo "==> Creating ccnp-node-measurement-example deployment"
kubectl apply -f $TEMP_MANIFEST_FILE
for i in {1..5}
do
    POD_NAME=$(kubectl get po | grep ccnp-node-measurement-example | grep Running | awk '{ print $1 }')
    if [[ -z "$POD_NAME" ]]
    then
        sleep 3
        echo "Retrying $i time ..."
    else
        break
    fi
done

if [[ -z "$POD_NAME" ]]; then
    echo "No ccnp-node-measurement-example pod with status running! Please check your deployment."
    exit 1
fi
echo ""

echo "Step 2:  Execute node measurement fetching and verification"
IFS=' ' read -r -a arr <<< "${index}"
if [[ ${#arr[@]} -eq 0 ]]
then
	kubectl exec -it "$POD_NAME" -- python3 fetch_node_measurement.py
else
	kubectl exec -it "$POD_NAME" -- python3 fetch_node_measurement.py --verify-register-index "${arr[@]}"
fi
