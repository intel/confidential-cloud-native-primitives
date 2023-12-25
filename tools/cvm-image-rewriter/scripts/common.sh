#!/bin/bash
#
# Common Scripts
#

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

# Check whether the tools are installed from the provided list
#
# args:
#   array of tool name  -    the name list to be checked
check_tools() {
    arr=("$@")
    is_missing=false
    for i in "${arr[@]}";
    do
        [[ "$(command -v "$i")" ]] || { info "MISSING: $i is not installed" 1>&2 ; is_missing=true ;}
    done

    [[ $is_missing != true ]] || { error "Please install missing tools"; }
}
