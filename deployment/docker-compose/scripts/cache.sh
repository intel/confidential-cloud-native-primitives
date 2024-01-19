#!/bin/bash

DIR=$(dirname "$(readlink -f "$0")")
# shellcheck disable=SC1091
. "$DIR"/scripts/comm.sh

check_cache_dir() {
    if [[ -d "$CCNP_CACHE_DIR" ]]; then
    	error "Cache Dir $CCNP_CACHE_DIR Exists. Please Back & Delete It"
    fi
}

create_cache_dir() {
    info "Cache Dir Being Created: $CCNP_CACHE_DIR"
    mkdir -p "$CCNP_CACHE_DIR"
    mkdir -p "$CCNP_CACHE_DIR/run/ccnp-eventlog"
    mkdir -p "$CCNP_CACHE_DIR/run/ccnp/uds"
    mkdir -p "$CCNP_CACHE_DIR/eventlog-entry-dir"
    mkdir -p "$CCNP_CACHE_DIR/eventlog-data-dir"
    mkdir -p "$COMPOSE_CACHE_DIR"

    chmod 777 -R "$CCNP_CACHE_DIR"
    ok "Cache Dir Created: $CCNP_CACHE_DIR"
}

remove_cache_dir() {
    info "Cache Dir Being Removed"
    if [[ -d "$CCNP_CACHE_DIR" ]]; then
    	rm -rf "$CCNP_CACHE_DIR"
    fi
    ok "Cache Dir Removed"
}
