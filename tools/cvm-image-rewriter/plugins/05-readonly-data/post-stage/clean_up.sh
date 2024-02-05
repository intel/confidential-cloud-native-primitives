#!/bin/bash

DIR=$(dirname "$(readlink -f "$0")")
CLD_DIR="$DIR/../cloud-init"

if [[ -d "$CLD_DIR" ]]; then
    rm -rf "$CLD_DIR"
fi