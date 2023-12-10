#!/bin/bash

set -e

# Common Definitions
TOP_DIR="$(dirname $(readlink -f "$0"))"
SCRIPTS_DIR="${TOP_DIR}/scripts"
TARGET_FILES_DIR="$(mktemp -d /tmp/cvm_target_files.XXXXXX)"
INPUT_IMG=""
OUTPUT_IMG="output.qcow2"

# Scan all subdirectories from pre-stage and post-stage
pre_stage_dirs=("$TOP_DIR/pre-stage"/*/)
post_stage_dirs=("$TOP_DIR/post-stage"/*/)
IFS=$'\n' sorted=($(sort <<<"${pre_stage_dirs[*]}")); unset IFS
IFS=$'\n' sorted=($(sort <<<"${post_stage_dirs[*]}")); unset IFS

# Include common definitions and utilities
source ${SCRIPTS_DIR}/common.sh

#
# Display Usage information
#
usage() {

    cat <<EOM
Usage: $(basename "$0") [OPTION]...
Required
  -i <guest image>          Specify initial guest image file
EOM
}

#
# Prepare the files copying to target guest image.
#
# 1. Scan following content from all subdirectories under pre-stage
#    a) files ==> the root of target guest image
#    b) guest_run.sh ==> /opt/guest_scripts at target guest image
# 2. Copy all files to staging directory at $TAGET_FILES_DIR
# 3. Create rootfs_override.tar.gz from $TAGET_FILES_DIR
# 4. Copy rootfs_overide.tar.gz to target system and extract
#
prepare_target_files() {
    echo "Prepare target files ..."

    # Scan all files directory and copy the content to temporary directory
    for path_item in "${pre_stage_dirs[@]}"
    do
        # Copy the content from files directory to target guest images
        if [[ -d $path_item/files ]]; then
            info "Copy $path_item/files/ => $TARGET_FILES_DIR"
            cp $path_item/files/* $TARGET_FILES_DIR/ -fr
        fi

        # Copy all guest_run.sh from pre_stage dirs to target guest images
        if [[ -f $path_item/guest_run.sh ]]; then
            info "Copy $path_item/guest_run.sh ==> $TARGET_FILES_DIR/opt/guest-scripts/$(basename $path_item)_guest_run.sh"
            chmod +x $path_item/guest_run.sh
            mkdir -p $TARGET_FILES_DIR/opt/guest-scripts/
            cp $path_item/guest_run.sh $TARGET_FILES_DIR/opt/guest-scripts/$(basename $path_item)_guest_run.sh
        fi
    done

    info "List all files to be copied to target image at $TARGET_FILES_DIR/..."
    ls $TARGET_FILES_DIR

    info "Copy all files to target guest image ..."

    pushd $TARGET_FILES_DIR/
    tar cpzf /tmp/rootfs_overide.tar.gz .
    popd

    virt-customize -a ${OUTPUT_IMG} \
        --copy-in /tmp/rootfs_overide.tar.gz:/root/ \
        --run-command 'mkdir /root/rootfs/ && tar zxvf /root/rootfs_overide.tar.gz -C /root/rootfs' \
        --run-command 'cp -r /root/rootfs/* / && rm -rf /root/rootfs*'

    rm /tmp/rootfs_overide.tar.gz
}

#
# Run the host_run.sh script from each subdirectories at pre-stage
#
run_pre_stage() {
    for path_item in "${pre_stage_dirs[@]}"
    do
        if [[ -f $path_item/host_run.sh ]]; then
            info "Execute the host_run.sh at $path_item"
            chmod +x $path_item/host_run.sh
            $path_item/host_run.sh
        fi
    done
}


#
# Run the host_run.sh script from each subdirectories at post-stage
#
run_post_stage() {
    for path_item in "${post_stage_dirs[@]}"
    do
        if [[ -f $path_item/host_run.sh ]]; then
            info "Execute the host_run.sh at $path_item"
            chmod +x $path_item/host_run.sh
            $path_item/host_run.sh
        fi
    done
}

do_pre_stage() {
    info "Run pre-stage..."
    prepare_target_files
    run_pre_stage
}

do_post_stage() {
    info "Run post-stage..."
    run_post_stage
}

do_cloud_init() {
    info "Run cloud-init..."

    pushd ${TOP_DIR}/cloud-init
    info "Prepare cloud-init ISO image..."
    [ -e /tmp/ciiso.iso ] && rm /tmp/ciiso.iso
    #cloud-init devel make-mime -a ./cloud-config.yaml:cloud-config > ./user-data
    genisoimage -output /tmp/ciiso.iso -volid cidata -joliet -rock user-data meta-data
    ok "Generate the cloud-init ISO image..."
    popd

    virt-install --memory 4096 --vcpus 4 --name tdx-config-cloud-init \
        --disk ${OUTPUT_IMG} \
        --disk /tmp/ciiso.iso,device=cdrom \
        --os-type Linux \
	    --os-variant ubuntu21.10 \
        --virt-type kvm \
        --graphics none \
        --import \
        --wait=3
    ok "Complete cloud-init..."
    sleep 1

    virsh destroy tdx-config-cloud-init || true
    virsh undefine tdx-config-cloud-init || true
}

cleanup() {
    exit_code=$?
    info "cleanup"
    if [[ -d $TARGET_FILES_DIR ]]; then
        info "Delete temporary directory $TARGET_FILES_DIR..."
        rm -fr $TARGET_FILES_DIR
    fi
    exit $exit_code
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
        INPUT_IMG=$(readlink -f $INPUT_IMG)
        if [[ ! -f ${INPUT_IMG} ]]; then
            error "File not exist ${INPUT_IMG}"
        fi
    fi

    ok "================================="
    ok "Input image: ${INPUT_IMG}"
    ok "Output image: ${OUTPUT_IMG}"
    ok "================================="

    # Create output image
    cp ${INPUT_IMG} ${OUTPUT_IMG}
}

trap cleanup EXIT
process_args "$@"

export GUEST_IMG=${OUTPUT_IMG}

do_pre_stage
do_cloud_init
do_post_stage

ok "Complete."
