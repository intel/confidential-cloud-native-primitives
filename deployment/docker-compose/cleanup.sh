#!/bin/bash

set -e

DIR=$(dirname "$(readlink -f "$0")")

# shellcheck disable=SC1091
. "$DIR"/scripts/docker_compose.sh
docker_compose_down

# shellcheck disable=SC1091
. "$DIR"/scripts/cache.sh
remove_cache_dir
