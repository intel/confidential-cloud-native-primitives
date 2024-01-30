#!/bin/bash

# Go to this dir
pushd "$(dirname "$(readlink -f "$0")")" || exit 0

# shellcheck disable=SC1091
source ../../../scripts/common.sh

# Check CVM_USER, CVM_AUTH_KEY
CVM_USER="${CVM_USER:-cvm}"
info "Config user: $CVM_USER"

if [[ -z "$CVM_AUTH_KEY" ]]; then
    warn "SKIP: CVM_AUTH_KEY is not defined via environment variable 'CVM_AUTH_KEY'"
    exit 0
fi
info "ssh pubkey: $CVM_AUTH_KEY"

# Generate cloud-config
mkdir -p ../cloud-init/cloud-config/
cat > ../cloud-init/cloud-config/04-user-authkey.yaml << EOL
#cloud-config
merge_how:
  - name: list
    settings: [append]
  - name: dict
    settings: [no_replace, recurse_list]
users:
  - default
  - name: $CVM_USER
    groups: sudo
    sudo: ALL=(ALL) NOPASSWD:ALL
    shell: /bin/bash
    ssh_authorized_keys:
      - $CVM_AUTH_KEY
EOL

popd || exit 0
