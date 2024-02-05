#!/bin/bash

DIR=$(dirname "$(readlink -f "$0")")
CLD_SH_REGISTER_FILE_HASH="01-ima-register-file-hash.sh"
CLD_SH="$DIR/../cloud-init/x-shellscript/$CLD_SH_REGISTER_FILE_HASH"

if [[ -f "$CLD_SH" ]]; then
    rm "$CLD_SH"
fi