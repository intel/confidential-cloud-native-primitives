#!/bin/bash

INITRAMFS_TOOLS_HOOKS_DIR=/etc/initramfs-tools/hooks

CURR_DIR=$(pwd)
CACHE_DIR=$CURR_DIR/cache
TEMP_DIR=$CURR_DIR/templates

copy_initramfs_tools_deps_into_image() {
    IMG=$1
    local -n _src_pkgs=$2
    local -n _dst_pkgs=$3
    len_of_arr=${#_src_pkgs[@]}

    for ((i=0;i<len_of_arr;i++))
    do
        if [[ ${_src_pkgs[i]} != "None" ]]; then
            # echo "|||" ${_src_pkgs[i]}:${_dst_pkgs[i]}
            virt-customize -a $IMG --copy-in ${_src_pkgs[i]}:${_dst_pkgs[i]}   
        fi
    done
}


create_initramfs_tools_hooks() {
    local -n src_pkgs=$1
    local -n dst_pkgs=$2
    len_of_arr=${#src_pkgs[@]}

    CACHED_HOOKS_DIR=$CACHE_DIR$INITRAMFS_TOOLS_HOOKS_DIR
    mkdir -p $CACHED_HOOKS_DIR
    CACHED_FILE=$CACHED_HOOKS_DIR/initrd-custom-update.sh

    cp $TEMP_DIR/initrd-custom-update.sh.template $CACHED_FILE

    for ((i=0;i<len_of_arr;i++))
    do
        if [ ${src_pkgs[i]} != "None" ]; then
            pkg=${dst_pkgs[i]}/$(basename ${src_pkgs[i]})
            echo "copy_exec $pkg" >> $CACHED_FILE  
        fi
    done
}



