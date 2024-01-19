#!/bin/bash

DIR=$(dirname "$(readlink -f "$0")")
# shellcheck disable=SC1091
. "$DIR"/scripts/comm.sh

# find valid tdx device
check_dev_tdx() {
    cur_dev=$(find /dev/ -name 'tdx*')
    case "$cur_dev" in
        "/dev/tdx-attest")
            error "The version of TDX module is too old. Please Upgrade greater TDX 1.0"
            ;;
        "/dev/tdx-guest" | "/dev/tdx_guest")
            ;;
        *)
            error "No Valid TDX Device. Make sure the TDX enabled"
            ;;
    esac
    echo "$cur_dev"
}

# make tdx device readable & writable
grant_dev_tdx() {
    cur_dev=$(check_dev_tdx)
    chmod 0666 "$cur_dev"
}
