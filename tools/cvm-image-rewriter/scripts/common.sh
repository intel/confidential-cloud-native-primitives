#!/bin/bash

SCRIPTS_DIR="$(dirname $(readlink -f "$0"))"
TOP_DIR="$(dirname ${SCRIPTS_DIR})"

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
