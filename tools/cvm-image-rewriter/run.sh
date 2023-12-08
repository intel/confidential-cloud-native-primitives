#!/bin/bash

set -e

CURR_DIR="$(dirname $(readlink -f "$0"))"
TARGET_FILES_DIR="$(mktemp -d /tmp/cvm_target_files.XXXXXX)"
GUEST_IMG=""
OUTPUT_IMG="output.qcow2"

# Scan directories in pre-stage and post-stage
pre_stage_dirs=("$CURR_DIR/pre-stage"/*/)
post_stage_dirs=("$CURR_DIR/post-stage"/*/)
IFS=$'\n' sorted=($(sort <<<"${pre_stage_dirs[*]}")); unset IFS
IFS=$'\n' sorted=($(sort <<<"${post_stage_dirs[*]}")); unset IFS

info() {
    echo -e "\e[1;33m$*\e[0;0m"
}

ok() {
    echo -e "\e[1;32mSUCCESS: $*\e[0;0m"
}

error() {
    echo -e "\e[1;31mERROR: $*\e[0;0m"
    exit 1
}

warn() {
    echo -e "\e[1;33mWARN: $*\e[0;0m"
}

usage() {
    cat <<EOM
Usage: $(basename "$0") [OPTION]...
Required
  -i <guest image>          Specify initial guest image file
EOM
}

prepare_target_files() {
    echo "Prepare target files ..."

    # Scan all files directory and copy the content to temporary directory
    for path_item in "${pre_stage_dirs[@]}"
    do
        if [[ -d $path_item/files ]]; then
            echo "Copy $path_item/files/ => $TARGET_FILES_DIR"
            cp $path_item/files/* $TARGET_FILES_DIR/ -fr
        fi
    done

    info "List all files to be copied to target image at $TARGET_FILES_DIR/..."
    ls $TARGET_FILES_DIR

    info "Copy all files to target guest image ..."
    # TODO: Copy all content from temporary directory to target guest image

    pushd $TARGET_FILES_DIR/
    tar cpjf /tmp/rootfs_overide.tar.bz2 .
    popd

    virt-customize -a ${OUTPUT_IMG} \
        --copy-in /tmp/rootfs_overide.tar.bz2:/tmp
}

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

do_pre_stage() {
    info "Run pre-stage..."
    prepare_target_files
    run_pre_stage
}

do_post_stage() {
    info "Run post-stage..."
}

do_cloud_init() {
    info "Run cloud-init..."
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
        i) GUEST_IMG=$OPTARG ;;
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

    if [[ -z $GUEST_IMG ]]; then
        error "Please specify the input guest image file via -i"
    else
        GUEST_IMG=$(readlink -f $GUEST_IMG)
        if [[ ! -f ${GUEST_IMG} ]]; then
            error "File not exist ${GUEST_IMG}"
        fi
    fi

    ok "================================="
    ok "Input image: ${GUEST_IMG}"
    ok "Output image: ${OUTPUT_IMG}"
    ok "================================="

    # Create output image
    cp ${GUEST_IMG} ${OUTPUT_IMG}
}

trap cleanup EXIT
process_args "$@"

do_pre_stage
do_cloud_init
do_post_stage
cleanup

ok "Complete."