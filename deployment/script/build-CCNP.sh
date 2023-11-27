#!/bin/bash
set -o errexit

DEVICE_PLUGIN=ccnp-device-plugin:0.1
QUOTE=ccnp-quote-server:0.1
MEASUREMENT=ccnp-measurement-server:0.1
EVENTLOG=ccnp-eventlog-server:0.1
WORK_DIR=$(cd "$(dirname "$0")"; pwd)

#Build images
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
popd

