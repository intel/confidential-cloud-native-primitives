#!/bin/bash
# shellcheck disable=SC2086
set -e

# Common Definitions
TOP_DIR=$(dirname "$(readlink -f "$0")")
SCRIPTS_DIR="${TOP_DIR}/scripts"
TARGET_FILES_DIR="$(mktemp -d /tmp/cvm_target_files.XXXXXX)"
INPUT_IMG=""
OUTPUT_IMG="output.qcow2"
TIMEOUT=3
CONNECTION_SOCK=""
CONSOLE_OPT=""

# Scan all subdirectories from pre-stage and post-stage
pre_stage_dirs=("$TOP_DIR/pre-stage"/*/)
post_stage_dirs=("$TOP_DIR/post-stage"/*/)
# shellcheck disable=SC2034,SC2207
IFS=$'\n' sorted=($(sort <<<"${pre_stage_dirs[*]}")); unset IFS
# shellcheck disable=SC2034,SC2207
IFS=$'\n' sorted=($(sort <<<"${post_stage_dirs[*]}")); unset IFS

# Include common definitions and utilities
# shellcheck disable=SC1091
source ${SCRIPTS_DIR}/common.sh

#
# Display Usage information
#
usage() {

    cat <<EOM
Usage: $(basename "$0") [OPTION]...
Required
  -i <guest image>          Specify initial guest image file
Optional
  -t <number of minutes>    Specify the timeout of rewriting, 3 minutes default,
                            If enabling ima, recommend timeout >6 minutes
  -s <connection socket>    Default is connection URI is qemu:///system,
                            if install libvirt, you can specify to "/var/run/libvirt/libvirt-sock"
                            then the corresponding URI is "qemu+unix:///system?socket=/var/run/libvirt/libvirt-sock"
  -n                        Silence running for virt-install, no output

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
        #
        # If want to skip some steps, please create a file named "NOT_RUN"
        # under the plugin directory. For example:
        #
        # ``touch pre-stage/01-resize-image/NOT_RUN``
        #
        if [[ -f $path_item/NOT_RUN ]]; then
            info "Skip to run $path_item ... "
            continue
        fi

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
            cp "$path_item"/guest_run.sh "$TARGET_FILES_DIR"/opt/guest-scripts/"$(basename "$path_item")"_guest_run.sh
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
        --run-command 'cp -r /root/rootfs/* / && rm -rf /root/rootfs*' || error "Fail to customize the image..."

    rm /tmp/rootfs_overide.tar.gz
}

#
# Run the host_run.sh script from each subdirectories at pre-stage
#
run_pre_stage() {
    for path_item in "${pre_stage_dirs[@]}"
    do
        #
        # If want to skip some steps, please create a file named "NOT_RUN"
        # under the plugin directory. For example:
        #
        # ``touch pre-stage/01-resize-image/NOT_RUN``
        #
        if [[ -f $path_item/NOT_RUN ]]; then
            info "Skip to run $path_item ... "
            continue
        fi

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
        #
        # If want to skip some steps, please create a file named "NOT_RUN"
        # under the plugin directory. For example:
        #
        # ``touch pre-stage/01-resize-image/NOT_RUN``
        #
        if [[ -f $path_item/NOT_RUN ]]; then
            info "Skip to run $path_item ... "
            continue
        fi

        if [[ -f $path_item/host_run.sh ]]; then
            info "Execute the host_run.sh at $path_item"
            chmod +x $path_item/host_run.sh
            $path_item/host_run.sh
        fi
    done
    # TODO: image status check at post-stage
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

_generate_cloud_init_meta_data() {
    info "Generate cloud init meta-data"
    pushd ${TOP_DIR}/cloud-init
    awk -v inst_name="tdx-inst-$(date '+%Y-%m-%d-%H-%M-%S')"  \
        '{if ($1 ~ /instance-id/) $2=inst_name; print}' meta-data.template \
        > meta-data
    popd
    ok "Complete cloud init meta-data"
}

_cloud_init_user_data_fromat_check() {
    input_file=$1
    order_check=$2
    if [[ $order_check == "true" ]] && ! [[ $(basename $input_file) =~ ^[0-9]{2} ]]; then
        error "User data format file name should start with two digitals: $input_file\n\
        e.g. 03-example or 55-another"s
    fi

    file_type=$(head -1 $input_file)

    case "$file_type" in
    "#cloud-config")
        merge_opt=$(grep "merge_how" "$input_file" || true)
        example_path="tools/cvm-image-rewriter/cloud-init/user-data.basic"
        if  [[ "$merge_opt" == "" ]] || [[ "$merge_opt" == "#"* ]]; then
            error "Cloud config need entry 'merge_how',\n\
            add one or uncomment it in $input_file,\n\
            refer default in $example_path,\n\
            or https://cloudinit.readthedocs.io/en/latest/reference/merging.html"
        fi
        info "user data: cloud-config ---> " $input_file
    ;;
    \#!*)
        info "user data: x-shellscript ---> " $input_file
    ;;
    *)
    info "Supported user data: cloud-config & x-shellscript"
    error "Unsupported user data: $input_file, starting with $file_type"
    ;;
    esac
}

_generate_cloud_init_user_data() {
    CLD_INIT_CONFIG_SUFFIX=":cloud-config"
    CLD_INIT_SCRIPT_SUFFIX=":x-shellscript"
    USER_DATA_BASIC="./cloud-init/user-data.basic"
    USER_DATA_TARGET="./cloud-init/user-data"

    info "Generate cloud init user-data"

    # basic user data
    info "Start user data format check"

    _cloud_init_user_data_fromat_check $USER_DATA_BASIC false
    ARGS=" -a "$(realpath $USER_DATA_BASIC)$CLD_INIT_CONFIG_SUFFIX

    # find all cloud init user data in pre-stage
    for path_item in "${pre_stage_dirs[@]}"
    do
        #
        # If want to skip some steps, please create a file named "NOT_RUN"
        # under the plugin directory. For example:
        #
        # ``touch pre-stage/01-resize-image/NOT_RUN``
        #
        if [[ -f $path_item/NOT_RUN ]]; then
            info "Skip to run $path_item ... "
            continue
        fi

        sub_dirs=("$path_item"/*/)
        for dir in "${sub_dirs[@]}"
        do
            if [[ ${dir} == *"cloud-init"* ]]; then
                CLD_CONFIG_DIR=${dir}cloud-config
                CLD_SCRIPT_DIR=${dir}x-shellscript

                # find cloud init cloud-config
                for f in "$CLD_CONFIG_DIR"/*
                do
                    if [[ -f "$f" ]]; then
                        _cloud_init_user_data_fromat_check $f true
                        ARGS+=" -a $f$CLD_INIT_CONFIG_SUFFIX"
                    fi

                done

                # find cloud init x-shellscript
                for f in "$CLD_SCRIPT_DIR"/*
                do
                    if [[ -f "$f" ]]; then
                        _cloud_init_user_data_fromat_check $f true
                        ARGS+=" -a $f$CLD_INIT_SCRIPT_SUFFIX"
                    fi
                done
                break
            fi
        done
    done

    ok "Complete user data format check"

    cloud-init devel make-mime $ARGS > $USER_DATA_TARGET
    ok "Complete cloud init user-data"
}

do_cloud_init() {
    _generate_cloud_init_meta_data
    _generate_cloud_init_user_data

    info "Run cloud-init..."

    pushd ${TOP_DIR}/cloud-init
    info "Prepare cloud-init ISO image..."
    if [[ -f /tmp/ciiso.iso ]]; then
        rm /tmp/ciiso.iso -fr || sudo rm /tmp/ciiso.iso -fr
    fi
    #cloud-init devel make-mime -a ./cloud-config.yaml:cloud-config > ./user-data
    genisoimage -output /tmp/ciiso.iso -volid cidata -joliet -rock user-data meta-data
    ok "Generate the cloud-init ISO image..."
    popd

    CONNECT_URI="qemu:///system"
    if [[ -n ${CONNECTION_SOCK} ]]; then
        CONNECT_URI="qemu+unix:///system?socket=${CONNECTION_SOCK}"
    fi

    # Clean up before running virt-install
    virsh destroy tdx-config-cloud-init &> /dev/null || true
    virsh undefine tdx-config-cloud-init &> /dev/null || true

    virt-install --memory 4096 --vcpus 4 --name tdx-config-cloud-init \
        --disk ${OUTPUT_IMG} \
        --connect ${CONNECT_URI} \
        --disk /tmp/ciiso.iso,device=cdrom \
        --os-type Linux \
	    --os-variant ubuntu21.10 \
        --virt-type kvm \
        --graphics none \
        --import \
        --wait=$TIMEOUT \
        ${CONSOLE_OPT}
    # TODO: check return status
    ok "Complete cloud-init..."
    sleep 1

    virsh destroy tdx-config-cloud-init &> /dev/null || true
    virsh undefine tdx-config-cloud-init &> /dev/null || true
}

cleanup() {
    exit_code=$?
    info "cleanup"
    if [[ -d $TARGET_FILES_DIR ]]; then
        info "Delete temporary directory $TARGET_FILES_DIR..."
        rm -fr $TARGET_FILES_DIR
    fi
    virsh destroy tdx-config-cloud-init &> /dev/null || true
    virsh undefine tdx-config-cloud-init &> /dev/null || true
    exit $exit_code
}

#
# when run virt-customize, it need read access on the /boot/vmlinuz-${uname -r}
# of the host, so please make sure read permission or use root user to run
# this script.
#
check_kernel_image_access() {
    curr_kernel_image="/boot/vmlinuz-$(uname -r)"
    test -r "$curr_kernel_image" || {
        error "Need read permission on $curr_kernel_image, please run following"\
            "command: \n\t sudo chmod +r $curr_kernel_image \n" \
            "Or run this script under root user"
    }
}

process_args() {
    while getopts ":i:t:s:nh" option; do
        case "$option" in
        i) INPUT_IMG=$OPTARG ;;
        t) TIMEOUT=$OPTARG ;;
        s) CONNECTION_SOCK=$OPTARG;;
        n) CONSOLE_OPT="--noautoconsole";;
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

    if [[ -z "$INPUT_IMG" ]]; then
        error "Please specify the input guest image file via -i"
    else
        INPUT_IMG=$(readlink -f "$INPUT_IMG")
        if [[ ! -f "${INPUT_IMG}" ]]; then
            error "File not exist ${INPUT_IMG}"
        fi
    fi

    ok "================================="
    ok "Input image: ${INPUT_IMG}"
    ok "Output image: ${OUTPUT_IMG}"
    ok "Connection Sock: ${CONNECTION_SOCK}"
    ok "Console Opt: ${CONSOLE_OPT}"
    ok "================================="

    # Create output image
    cp "${INPUT_IMG}" "${OUTPUT_IMG}"

    if [[ -n "${CONNECTION_SOCK}" ]]; then
        [[ ! -f "${CONNECTION_SOCK}" ]] || {
            error  "File ${CONNECTION_SOCK} does not exist." ;
        }
        test -r "${CONNECTION_SOCK}" || {
            error "File ${CONNECTION_SOCK} could not be accessed, please make"\
                  " sure current user is belong to libvirtd group."
        }
    fi
}

check_tools qemu-img virt-customize virt-install genisoimage cloud-init git awk
check_kernel_image_access

trap cleanup EXIT

process_args "$@"

export GUEST_IMG=${OUTPUT_IMG}

do_pre_stage
do_cloud_init
do_post_stage

ok "Complete."
