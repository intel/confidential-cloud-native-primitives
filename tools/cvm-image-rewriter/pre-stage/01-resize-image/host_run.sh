#!/bin/bash

CURR_DIR=$(dirname "$(readlink -f "$0")")

# shellcheck disable=SC1091
source "${CURR_DIR}/../../scripts/common.sh"

info "Guest Image is at ${GUEST_IMG}..."

if [[ -z ${GUEST_SIZE} ]]; then
    warn "SKIP: Guest size is not defined via environment variable 'GUEST_SIZE'"
else
    qemu-img resize "${GUEST_IMG}" "${GUEST_SIZE}"
    virt-customize -a "${GUEST_IMG}" \
            --run-command 'growpart /dev/sda 1' \
            --run-command 'resize2fs /dev/sda1'
fi
