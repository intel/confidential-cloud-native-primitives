#!/bin/bash
#
# Create TD VM from libvirt template
#
set -e

CURR_DIR=$(readlink -f "$(dirname "$0")")

GUEST_IMG="tdx-guest-ubuntu22.04.qcow2"
GUEST_NAME="tdx-guest"
GUEST_ROOTDIR=/tmp/libvirt-vms
TEMPLATE="${CURR_DIR}/tdx-libvirt-ubuntu-host.xml.template"
FORCE=false
VCPU_NUM=1
MEM_SIZE=4

usage() {
    cat << EOM
Usage: $(basename "$0") [OPTION]...
  -i <guest image file>     Default is tdx-guest-ubuntu22.04.qcow2 under current directory
  -n <guest name>           Name of TD guest
  -t <template file>        Default is ./tdx-libvirt-ubuntu-host.xml.template
  -f                        Force recreate
  -v <vcpu number>          VM vCPU number
  -m <memory size in GB>    VM memory size in GB
  -h                        Show this help
EOM
}

pre-check() {
    # Check whether current user belong to libvirt
    if [[ ! $(id -nG "$USER") == *"libvirt"* ]]; then
    echo WARNING! Please add user "$USER" into group "libvirt" via \"sudo usermod -aG libvirt "$USER"\"
    return 1
    fi
}

process_args() {
    while getopts ":i:n:t:v:m:fh" option; do
        case "$option" in
            i) GUEST_IMG=$OPTARG;;
            n) GUEST_NAME=$OPTARG;;
            t) TEMPLATE=$OPTARG;;
            v) VCPU_NUM=$OPTARG;;
            m) MEM_SIZE=$OPTARG;;
            f) FORCE=true;;
            h) usage
               exit 0
               ;;
            *)
               echo "Invalid option '-$OPTARG'"
               usage
               exit 1
               ;;
        esac
    done

    echo "====================================================================="
    echo " Use template   : ${TEMPLATE}"
    echo " Guest name     : ${GUEST_NAME}"
    echo " Guest XML      : ${GUEST_ROOTDIR}/${GUEST_NAME}.xml"
    echo " Guest Image    : ${GUEST_ROOTDIR}/${GUEST_NAME}.qcow2"
    echo " Force Recreate : ${FORCE}"
    echo "====================================================================="

    if [[ ! -f ${GUEST_IMG} ]]; then
        echo "Error: Guest image ${GUEST_IMG} does not exist"
        exit 1
    fi

    if [[ ${FORCE} == "true" ]]; then
        echo "> Clean up the old guest... "
        virsh destroy "${GUEST_NAME}" || true
        sleep 2
        virsh undefine "${GUEST_NAME}" || true
        sleep 2
        rm "${GUEST_ROOTDIR}/${GUEST_NAME}.xml" -fr || true
    fi

    if [[ -f ${GUEST_ROOTDIR}/${GUEST_NAME}.xml ]]; then
        echo "Error: Guest XML ${GUEST_ROOTDIR}/${GUEST_NAME}.xml already exist."
        echo "Error: you can delete the old one via 'rm ${GUEST_ROOTDIR}/${GUEST_NAME}.xml'"
        exit 1
    fi

    if [[ ! -f ${TEMPLATE} ]]; then
        echo "Template ${TEMPLATE} does not exist".
        echo "Please specify via -t"
        exit 1
    fi

    # Validate the number of vCPUs
    if ! [[ ${VCPU_NUM} =~ ^[0-9]+$ && ${VCPU_NUM} -gt 0 ]]; then
        echo "Error: Invalid number of vCPUs: ${VCPU_NUM}"
        usage
        exit 1
    fi

    # Validate the size of memory
    if ! [[ ${MEM_SIZE} =~ ^[0-9]+$ && ${MEM_SIZE} -gt 0 ]]; then
        echo "Error: Invalid memory size: ${MEM_SIZE}"
        usage
        exit 1
    fi
}

create-vm() {
    mkdir -p ${GUEST_ROOTDIR}/
    echo "> Create ${GUEST_ROOTDIR}/${GUEST_NAME}.qcow2..."
    cp "${GUEST_IMG}" "${GUEST_ROOTDIR}/${GUEST_NAME}.qcow2"
    echo "> Create ${GUEST_ROOTDIR}/${GUEST_NAME}.xml..."
    cp "${TEMPLATE}" "${GUEST_ROOTDIR}/${GUEST_NAME}.xml"

    echo "> Modify configurations..."
    sed -i "s/.*<name>.*/<name>${GUEST_NAME}<\/name>/" "${GUEST_ROOTDIR}/${GUEST_NAME}.xml"
    sed -i "s#/path/to/image#${GUEST_ROOTDIR}/${GUEST_NAME}.qcow2#" "${GUEST_ROOTDIR}/${GUEST_NAME}.xml"
    sed -i "s/REPLACE_VCPU_NUM/${VCPU_NUM}/g" "${GUEST_ROOTDIR}/${GUEST_NAME}.xml"
    sed -i "s/REPLACE_MEM_SIZE/${MEM_SIZE}/g" "${GUEST_ROOTDIR}/${GUEST_NAME}.xml"
}

start-vm() {
    echo "> Create VM domain..."
    virsh define "${GUEST_ROOTDIR}/${GUEST_NAME}.xml"
    sleep 2
    echo "> Start VM..."
    virsh start "${GUEST_NAME}"
    sleep 2
    echo "> Connect console..."
    virsh console "${GUEST_NAME}"
}

pre-check
process_args "$@"
create-vm
start-vm
