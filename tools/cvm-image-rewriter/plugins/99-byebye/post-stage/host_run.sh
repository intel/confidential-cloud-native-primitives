#!/bin/bash

CURR_DIR=$(dirname "$(readlink -f "$0")")
TOP_DIR="${CURR_DIR}/../../../"
SCRIPTS_DIR="${TOP_DIR}/scripts"

# shellcheck disable=SC1091
source "${SCRIPTS_DIR}/common.sh"

ok "Success to create guest image ${GUEST_IMG}..."
