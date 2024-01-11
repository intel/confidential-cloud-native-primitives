#!/bin/bash

DIR=$(dirname "$(readlink -f "$0")")
# shellcheck disable=SC1091
. "$DIR"/scripts/device.sh

CONFIG_DIR="$DIR/configs"

create_composes() {
    EVENTLOG_IMAGE=$1
    MEASUREMENT_IMAGE=$2
    QUOTE_IMAGE=$3

    DEV_TDX=$(check_dev_tdx)

    sed "s@\#EVENTLOG_IMAGE@$EVENTLOG_IMAGE@g" "$CONFIG_DIR"/eventlog-compose.yaml.template \
        > "$COMPOSE_CACHE_DIR"/eventlog-compose.yaml
    sed "s@\#MEASUREMENT_IMAGE@$MEASUREMENT_IMAGE@g" "$CONFIG_DIR"/measurement-compose.yaml.template \
        > "$COMPOSE_CACHE_DIR"/measurement-compose.yaml
    sed "s@\#QUOTE_IMAGE@$QUOTE_IMAGE@g" "$CONFIG_DIR"/quote-compose.yaml.template \
        > "$COMPOSE_CACHE_DIR"/quote-compose.yaml
    
    sed -i "s@\#DEV_TDX@$DEV_TDX@g" "$COMPOSE_CACHE_DIR"/measurement-compose.yaml
    sed -i "s@\#DEV_TDX@$DEV_TDX@g" "$COMPOSE_CACHE_DIR"/quote-compose.yaml

}

docker_compose_up() {
    if ! [ -d "$COMPOSE_CACHE_DIR" ]; then
        error "Compose Cache Dir not Exist: $COMPOSE_CACHE_DIR"
    fi

    CONFIGS=("$COMPOSE_CACHE_DIR"/*)
    for config in "${CONFIGS[@]}"
    do
        info "Compose $config Being Deployed"
    	docker compose -f "$config" up -d 
        ok "Compose $config Deployed"
    done
}

docker_compose_down() {
    if ! [ -d "$COMPOSE_CACHE_DIR" ]; then
        error "Compose Cache Dir not Exist: $COMPOSE_CACHE_DIR"
    fi
    
    CONFIGS=("$COMPOSE_CACHE_DIR"/*)
    for config in "${CONFIGS[@]}"
    do
        name_line=$(head -1 "$config")
        name="${name_line#name:}"
        # shellcheck disable=SC2086
        compose=$(docker compose ls | grep $name || true)
        if [[ "$compose" == "" ]]; then
            continue
        fi
        info "Compose $config Being Down"
            docker compose -f "$config"  down
        ok "Compose $config Down"
    done
}
