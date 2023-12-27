#!/bin/bash
# Script to deploy and execute CCNP example container
# Attach the RTMR index after the script during execution to verify selected register

set -e

echo "Step 1:  Deploy CCNP example container for node measurement in Kubernetes"
kubectl apply -f ../manifests/ccnp-node-measurement-example-deployment.yaml
echo "Step 2:  Execute node measurement fetching and verification"
POD_NAME=$(kubectl get po | grep ccnp-node-measurement-example | awk '{ print $1 }')
if [ $# -eq 0 ]
then
	kubectl exec -it "$POD_NAME" -- python3 fetch_node_measurement.py
else
	kubectl exec -it "$POD_NAME" -- python3 fetch_node_measurement.py --verify-register-index "$@"
fi
