#!/bin/bash

set -e

DIR=$(dirname "$(readlink -f "$0")")

# shellcheck disable=SC1091
. "$DIR"/scripts/cache.sh
check_cache_dir
ok "Cache Dir Clear"

# shellcheck disable=SC1091
. "$DIR"/scripts/device.sh
grant_dev_tdx
ok "Dev TDX Valid"

info "Make Sure Service QGS&PCCS is Avaliable to Get Quote"
