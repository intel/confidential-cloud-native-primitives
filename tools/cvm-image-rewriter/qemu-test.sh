#!/bin/bash

TOP_DIR="$(dirname $(readlink -f "$0"))"
SCRIPTS_DIR="${TOP_DIR}/scripts"

INPUT_IMG="output.qcow2"

source ${SCRIPTS_DIR}/common.sh

usage() {
    cat <<EOM
Usage: $(basename "$0") [OPTION]...
Required
  -i <guest image>          Specify initial guest image file
EOM
}

process_args() {
    while getopts ":i:h" option; do
        case "$option" in
        i) INPUT_IMG=$OPTARG ;;
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

    if [[ -z $INPUT_IMG ]]; then
        error "Please specify the input guest image file via -i"
    else
        INPUT_IMG=$(readlink -f ${INPUT_IMG})
        if [[ ! -f ${INPUT_IMG} ]]; then
            error "File not exist ${INPUT_IMG}. Please specify the input"\
                  "guest image file via -i"
        fi
    fi
}

process_args "$@"

# remap CTRL-C to CTRL ]
echo "Remapping CTRL-C to CTRL-]"
stty intr ^]

virt-customize -a ${INPUT_IMG} --root-password password:123456
qemu-system-x86_64  \
  -machine accel=kvm,type=q35 \
  -cpu host \
  -m 2G \
  -nographic \
  -monitor pty\
  -device virtio-net-pci,netdev=net0 \
  -netdev user,id=net0,hostfwd=tcp::2222-:22 \
  -drive if=virtio,format=qcow2,file=${INPUT_IMG}

# restore CTRL-C mapping
stty intr ^c
