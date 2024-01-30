#!/bin/bash

DIR=$(dirname "$(readlink -f "$0")")
ETC_DIR="$DIR/../files/etc"

if [[ -d "$ETC_DIR" ]]; then
    rm -rf "$ETC_DIR"
fi