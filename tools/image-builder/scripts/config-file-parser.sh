#!/bin/bash

# TODO: parse_initrd_src_paths() {...}

parse_initrd_dst_paths () {
    config_file_path=$1
    
    # TODO: error check
    dst_pkgs=$(yq '.initrd.[].dst' $config_file_path)
    echo $dst_pkgs
}

parse_initrd_src_paths () {
    config_file_path=$1
    
    # TODO: error check
    src_pkgs=$(yq '.initrd.[].src' $config_file_path)
    echo $src_pkgs
}