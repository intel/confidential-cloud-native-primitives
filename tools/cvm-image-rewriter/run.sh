#!/bin/bash

set -e

TOP_DIR="$(dirname $(readlink -f "$0"))"
SCRIPTS_DIR="${TOP_DIR}/scripts"
TARGET_FILES_DIR="$(mktemp -d /tmp/cvm_target_files.XXXXXX)"
INPUT_IMG=""
OUTPUT_IMG="output.qcow2"

# Scan directories in pre-stage and post-stage
pre_stage_dirs=("$CURR_DIR/pre-stage"/*/)
post_stage_dirs=("$CURR_DIR/post-stage"/*/)
IFS=$'\n' sorted=($(sort <<<"${pre_stage_dirs[*]}")); unset IFS
IFS=$'\n' sorted=($(sort <<<"${post_stage_dirs[*]}")); unset IFS

source ${SCRIPTS_DIR}/common.sh

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

    pushd $TARGET_FILES_DIR/
    tar cpzf /tmp/rootfs_overide.tar.gz .
    popd

    virt-customize -a ${OUTPUT_IMG} \
        --copy-in /tmp/rootfs_overide.tar.gz:/root/ \
        --run-command 'tar zxvf /root/rootfs_overide.tar.gz -C /'
}

run_pre_stage() {
    export GUEST_IMG=${OUTPUT_IMG}
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

do_pre_stage
do_cloud_init
do_post_stage
cleanup

ok "Complete."
