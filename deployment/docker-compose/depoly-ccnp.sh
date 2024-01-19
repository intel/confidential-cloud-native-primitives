#!/bin/bash

set -e

TAG="latest"
EVENTLOG_IMAGE="ccnp-eventlog-server"
MEASUREMENT_IMAGE="ccnp-measurement-server"
QUOTE_IMAGE="ccnp-quote-server"
REGISTRY=""

DIR=$(dirname "$(readlink -f "$0")")

#
# Display Usage information
#
usage() {

    cat <<EOM
Usage: $(basename "$0") [OPTION]...
    -r <registry prefix>    the prefix string for registry
    -g <tag>                container image tag
    -h                      show help info
EOM
}

process_args() {
    while getopts ":r:g:h" option; do
        case "$option" in
        r) REGISTRY=$OPTARG ;;
        g) TAG=$OPTARG ;;
        h)
            usage
            exit 0
            ;;
        *)
            echo "Invalid option '-$OPTARG'"
            usage
            exit 1
            ;;
        esac
    done

    EVENTLOG_IMAGE="$EVENTLOG_IMAGE:$TAG"
    MEASUREMENT_IMAGE="$MEASUREMENT_IMAGE:$TAG"
    QUOTE_IMAGE="$QUOTE_IMAGE:$TAG"

    if [[ ${REGISTRY: -1} == "/" ]]; then
        REGISTRY="${REGISTRY%/}"
    fi
    if [[ $REGISTRY != "" ]]; then
        EVENTLOG_IMAGE="$REGISTRY/$EVENTLOG_IMAGE"
        MEASUREMENT_IMAGE="$REGISTRY/$MEASUREMENT_IMAGE"
        QUOTE_IMAGE="$REGISTRY/$QUOTE_IMAGE"
    fi
}

process_args "$@"

# shellcheck disable=SC1091
. "$DIR"/scripts/cache.sh
create_cache_dir

# shellcheck disable=SC1091
. "$DIR"/scripts/docker_compose.sh
create_composes "$EVENTLOG_IMAGE" "$MEASUREMENT_IMAGE" "$QUOTE_IMAGE"

docker_compose_up
