#!/bin/bash

DIR=$(dirname "$(readlink -f "$0")")

virt-customize -a "${GUEST_IMG}" \
    --run-command 'mkdir -p /etc/ima' \
    --run "$DIR"/guest_enable_ima_fix.sh
    