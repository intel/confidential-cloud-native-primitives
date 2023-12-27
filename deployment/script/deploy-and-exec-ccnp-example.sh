#!/bin/bash
# Script to deploy and execute CCNP example container
# Attach the RTMR index after the script during execution to verify selected register

set -e

DEFAULT_DOCKER_REPO=docker.io/library
DEFAULT_TAG=latest
TEMP_MANIFEST_FILE=/tmp/ccnp-node-measurement-example-deployment.yaml

usage() { echo "Usage: $0 [-r <registry-prefix>] [-g <image-tag>] [-i <register-index-to-verify>]"; exit 1; }
while getopts ":r:g:i:h" option; do
        case "${option}" in
            r) registry=${OPTARG};;
            g) tag=${OPTARG};;
            i) index=${OPTARG};;
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
kubectl apply -f $TEMP_MANIFEST_FILE
sleep 3

echo "Step 2:  Execute node measurement fetching and verification"
POD_NAME=$(kubectl get po | grep ccnp-node-measurement-example | awk '{ print $1 }')
if [[ -z "$POD_NAME" ]]; then
	echo "ccnp-node-measurement-example pod not found!"
fi
IFS=' ' read -r -a arr <<< "${index}"
if [[ ${#arr[@]} -eq 0 ]]
then
	kubectl exec -it "$POD_NAME" -- python3 fetch_node_measurement.py
else
	kubectl exec -it "$POD_NAME" -- python3 fetch_node_measurement.py --verify-register-index "${arr[@]}"
fi
