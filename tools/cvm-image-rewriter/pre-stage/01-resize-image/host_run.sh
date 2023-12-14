#!/bin/bash

echo "Guest Image is at ${GUEST_IMG}..."

qemu-img resize "${GUEST_IMG}" 50G
virt-customize -a "${GUEST_IMG}" \
        --run-command 'growpart /dev/sda 1' \
        --run-command 'resize2fs /dev/sda1'
