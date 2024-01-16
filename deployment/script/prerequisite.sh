#!/bin/bash

# This script implements the prerequisites for deploying CCNP, including installing docker, helm, python3-pip, 
# and setting the access permissions of the TD device node and the ccnp working directory on the TD node.

set -e

INSTALL_DOCKER=true
INSTALL_HELM=true
INSTALL_PIP=true
CCNP_UDEV=true
CCNP_UDS=true

UDEV_FILE=/etc/udev/rules.d
TDX_RULES_FILE=${UDEV_FILE}/90-tdx.rules
CCNP_CONF=/usr/lib/tmpfiles.d/ccnp.conf

function check_env {
    if command -v docker &> /dev/null; then
        INSTALL_DOCKER=false
        echo "Skip: Docker has been installed."       
    fi

    if  command -v helm &> /dev/null; then
        INSTALL_HELM=false
        echo "Skip: Helm has been installed."
    fi

    if  command -v pip &> /dev/null; then
        INSTALL_PIP=false
        echo "Skip: Python3-pip has been installed."
    fi

    if [ -e "$TDX_RULES_FILE" ]; then
        CCNP_UDEV=flase
        echo "Skip: CCNP udev rules has been set."
    fi

    if [ -e "$CCNP_CONF" ]; then
        CCNP_UDS=flase
        echo "Skip: CCNP uds dir has been prepared."
    fi
}

function install_docker {
    # install GPG key
    install -m 0755 -d /etc/apt/keyrings
    rm -f /etc/apt/keyrings/docker.gpg
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg
    chmod a+r /etc/apt/keyrings/docker.gpg

    # install repo
    # shellcheck disable=SC1091
    echo \
    "deb [arch=\"$(dpkg --print-architecture)\" signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
    $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
    tee /etc/apt/sources.list.d/docker.list > /dev/null
    apt-get update > /dev/null

    # install docker
    apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

    # start docker
    systemctl enable docker
    systemctl daemon-reload
    systemctl start docker
}

function install_helm {
    # install repo
    curl -fsSL https://baltocdn.com/helm/signing.asc | gpg --dearmor | tee /usr/share/keyrings/helm.gpg > /dev/null
    echo \
    "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/helm.gpg] https://baltocdn.com/helm/stable/debian/ all main" | \
    tee /etc/apt/sources.list.d/helm-stable-debian.list > /dev/null
    apt-get update > /dev/null

    # install helm
    apt-get install -y helm
}

function install_pip {
    # install python3-pip
    apt install -y python3-pip 
}

function ccnp_udev_rules { 
    mkdir -p ${UDEV_FILE}
    touch ${TDX_RULES_FILE}

    # check TDX device node
    OUTPUT=$(ls /dev/tdx* 2>&1 || true)

    if [[ "$OUTPUT" == *"/dev/tdx-guest"* ]]; then
        echo "SUBSYSTEM==\"misc\",KERNEL==\"tdx-guest\",MODE=\"0666\"" > $TDX_RULES_FILE
    elif [[ "$OUTPUT" == *"/dev/tdx_guest"* ]]; then
        echo "SUBSYSTEM==\"misc\",KERNEL==\"tdx_guest\",MODE=\"0666\"" > $TDX_RULES_FILE
    else
        echo "Error: No TDX deivce node with msg: $OUTPUT"
        exit 1
    fi

    # make the udev setup effective
    udevadm trigger
}

function ccnp_uds_dir {
    touch ${CCNP_CONF}
    echo "D /run/ccnp/uds 0757 - - -">${CCNP_CONF}

    # make the directory setup effective
    systemd-tmpfiles --create
}

function install_prereqs {
    echo "-----------Check Environment..."
    check_env

    if [[ "$INSTALL_DOCKER" = true ]]; then
        echo "-----------Install Docker..."
        install_docker
    fi

    if [[ "$INSTALL_HELM" = true ]]; then
        echo "-----------Install Helm..."
        install_helm
    fi

    if [[ "$INSTALL_PIP" = true ]]; then
        echo "-----------Install Python3-pip..."
        install_pip
    fi

    if [[ "$CCNP_UDEV" = true ]]; then
        echo "-----------Setup udev rules for CCNP device plugin..."
        ccnp_udev_rules
    fi

    if [[ "$CCNP_UDS" = true ]]; then
        echo "-----------Prepare the shared Unix Domain Socket directory for CCNP..."
        ccnp_uds_dir
    fi
}

install_prereqs
