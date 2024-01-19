#!/bin/bash

set -e

DIR=$(dirname "$(readlink -f "$0")")
# shellcheck disable=SC1091
. "$DIR"/scripts/device.sh
CONFIG_DIR="$DIR"/configs 	  

EXAMPLE_IMAGE=ccnp-node-measurement-example
TAG=latest
REGISTRY=""
DEV_TDX="/dev/tdx_guest"
DELETE_CTR=false
FROM_HOST=false

#
# Display Usage information
#
usage() {
    cat <<EOM
Usage: $(basename "$0") [OPTION]...
    -r <registry prefix>    the prefix string for registry
    -g <tag>                container image tag
    -d                      delete example container
    -o                      request from host
    -h                      show help info
EOM
}

process_args() {
    while getopts ":r:g:dho" option; do
        case "$option" in
        r) REGISTRY=$OPTARG ;;
        g) TAG=$OPTARG ;;
	    d) DELETE_CTR=true ;;
	    o) FROM_HOST=true ;;
        h)
            usage
            exit 0
            ;;
        *)
            echo "Invalid option '-$OPTARG'"
            usage
            exit 1
            ;;
        esac
    done

    EXAMPLE_IMAGE="$EXAMPLE_IMAGE:$TAG"

    if [[ "${REGISTRY: -1}" == "/" ]]; then
        REGISTRY="${REGISTRY%/}"
    fi
    if [[ "$REGISTRY" != "" ]]; then
        EXAMPLE_IMAGE="$REGISTRY/$EXAMPLE_IMAGE"
    fi

    DEV_TDX=$(check_dev_tdx)
}

validate_on_guest() {
    if [[ "$FROM_HOST" == "false" ]]; then
        return
    fi
    GRPCURL=grpcurl
    if [[ "$GRPCURL" == ""  ]]; then
       error "grpcurl missing. Please install it."
    fi
    info "Validate CCNP on Guest"
    # eventlog
    info "Require TDX RTMR Event Log"
    ret=$("$GRPCURL" -plaintext \
        -d '{"eventlog_level": 0, "eventlog_category": 0}' \
        -unix /tmp/docker_ccnp/run/ccnp/uds/eventlog.sock \
        Eventlog/GetEventlog)
    echo -e "$ret"
    ok "Eventlog server validated"
    
    # measurement
    info "Require TDX RTMR Measurement"
    ret=$("$GRPCURL" -plaintext \
            -d '{"measurement_type": 0, "measurement_category": 0}' \
            -unix /tmp/docker_ccnp/run/ccnp/uds/measurement.sock \
    	measurement.Measurement/GetMeasurement)
    echo -e "$ret"
    ok "Measurement Server Validated"
    
    # quote
    info "Require TDX Quote"
    ret=$("$GRPCURL" \
            -authority "dummy" \
            -d '{"user_data": "MTIzNDU2NzgxMjM0NTY3ODEyMzQ1Njc4MTIzNDU2NzgxMjM0NTY3ODEyMzQ1Njc4", "nonce":"IXUKoBO1UM3c1wopN4sY"}'  \
            -plaintext \
            -unix /tmp/docker_ccnp/run/ccnp/uds/quote-server.sock  \
    	quoteserver.GetQuote.GetQuote)
    echo -e "$ret"
    ok "Quote Server Validated"
    
    ok "CCNP Validated on Guest"
}


delete_example_ctr() {
    if [[ "$DELETE_CTR" == "false" ]]; then
        return
    fi

    info "Example Container Being Deleted"
    docker compose -f "$COMPOSE_CACHE_DIR"/ccnp-node-measurement-example.yaml down
    ok "Example Container Deleted"
}

validate_on_container() {
    info "Execute example Container ccnp-node-measurement-example"
    ctr_id=$(docker ps | grep node-measurement-example-ctr | awk '{print $1}')    
    if [[ "$ctr_id" == "" ]]; then
   	info "Example Container No Avaliable. Attempt Deploy It"
        sed "s@\#EXAMPLE_IMAGE@$EXAMPLE_IMAGE@g" "$CONFIG_DIR"/ccnp-node-measurement-example.yaml.template \
                    > "$COMPOSE_CACHE_DIR"/ccnp-node-measurement-example.yaml
        sed -i "s@\#DEV_TDX@$DEV_TDX@g" "$COMPOSE_CACHE_DIR"/ccnp-node-measurement-example.yaml
        docker compose -f "$COMPOSE_CACHE_DIR"/ccnp-node-measurement-example.yaml up -d
    fi

    ctr_id=$(docker ps | grep node-measurement-example-ctr | awk '{print $1}')
    if [[ "$ctr_id" == "" ]]; then
       error "Example Container Deploy Failed"
    fi

    ok "Example Container Avaliable. Compose file: $COMPOSE_CACHE_DIR/ccnp-node-measurement-example.yaml"
    ok "Example Container Avaliable. Compose file: $COMPOSE_CACHE_DIR/ccnp-node-measurement-example.yaml"
    docker exec -it "$ctr_id" python3 fetch_node_measurement.py > "$CCNP_CACHE_DIR"/measurement.log
    ok "Measurement Log Saved in File $CCNP_CACHE_DIR/measurement.log"
    ok "Example Container ccnp-node-measurement-example Executed"
}

process_args "$@"

validate_on_guest
validate_on_container
delete_example_ctr

