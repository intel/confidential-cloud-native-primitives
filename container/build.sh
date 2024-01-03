#!/bin/bash

set -e

curr_dir=$(readlink -f "$(dirname "${BASH_SOURCE[0]}")")
top_dir=$(dirname "${curr_dir}")
action="all"
registry=""
container="all"
tag="latest"
docker_build_clean_param=""
pccs=false
qgs=false
all_containers=()

#
# Scan the subdirectories, each one is for a container need to be built
#
# @param $1 The root directory
#
function scan_all_containers {
    cd "${1}"
    mapfile -t dirs < <(ls -d -- */)
    for item in "${dirs[@]}"
    do
        all_containers[${#all_containers[@]}]="${item::-1}"
    done
}

#
# Print usage
#
function usage {
    cat << EOM
usage: $(basename "$0") [OPTION]...
    -a <build|publish|save|all>  all is default, which not include save. Please execute save explicity if need.
    -r <registry prefix> the prefix string for registry
    -c <container name> same as directory name
    -g <tag> container image tag
    -f Clean build
    -p Flag to build PCCS
    -q Flag to build QGS
EOM
    exit 1
}

#
# Process arguments
#
function process_args {
    while getopts ":a:r:c:g:hfpq" option; do
        case "${option}" in
            a) action=${OPTARG};;
            r) registry=${OPTARG};;
            c) container=${OPTARG};;
            g) tag=${OPTARG};;
            h) usage;;
            f) docker_build_clean_param="--no-cache";;
	    p) pccs=true;;
	    q) qgs=true;;
            *)
        esac
    done

    if [[ ! "$action" =~ ^(build|publish|save|all) ]]; then
        echo "invalid type: $action"
        usage
    fi

    if [ ! -d "${container}" ]; then
        echo "${container} does not exist."
    fi

    if [[ "$registry" == "" ]]; then
        if [[ -z "$CCNP_REGISTRY" ]]; then
            echo "Error: Please specify your docker registry via -r <registry prefix> or set environment variable CCNP_REGISTRY."
            exit 1
        else
            registry=$CCNP_REGISTRY
        fi
    fi
}

#
# Build a container image
#
# @param $1: the (directory) name of container
#
function build_a_image {
    echo "Build container image => ${registry}/${1}:${tag}"

    cd "${curr_dir}"/"${1}"
    if [ -f "pre-build.sh" ]; then
        echo "Execute pre build script at ${curr_dir}/${1}/pre-build.sh"
        ./pre-build.sh || { echo 'Fail to execute pre-build.sh'; exit 1; }
    fi

    cd "${top_dir}"
    if [[ -n "${docker_build_clean_param}" ]]; then
        docker build \
             --build-arg http_proxy \
             --build-arg https_proxy \
             --build-arg no_proxy \
             --build-arg pip_mirror \
             -f container/"${1}"/Dockerfile \
             . \
             -t "${registry}/${1}:${tag}" \
             "${docker_build_clean_param}" || \
             { echo "Fail to build docker ${registry}/${1}:${tag}"; exit 1; }
    else
        docker build \
             --build-arg http_proxy \
             --build-arg https_proxy \
             --build-arg no_proxy \
             --build-arg pip_mirror \
             -f container/"${1}"/Dockerfile \
             . \
             -t "${registry}/${1}:${tag}" || \
             { echo "Fail to build docker ${registry}/${1}:${tag}"; exit 1; }
    fi

    echo "Complete build image => ${registry}/${1}:${tag}"

    cd "${curr_dir}"/"${1}"
    if [ -f "post-build.sh" ]; then
        echo "Execute post build script at ${curr_dir}/${1}/post-build.sh"
        ./post-build.sh || { echo "Fail to execute post-build.sh"; exit 1; }
    fi

    echo -e "\n\n"
}

#
# Build all containers
#
function build_images {
    if [[ "$container" == "all" ]]; then
        for item in "${all_containers[@]}"
        do
	    if [[ ${item} != "pccs" && ${item} != "qgs" ]] || [[ ${item} == "pccs" && ${pccs} == true ]] || [[ ${item} == "qgs" && ${qgs} == true ]] ; then
            build_a_image "$item"
	    fi
        done
    else
        build_a_image "$container"
    fi
}

#
# Publish a container image to given registry via "-r" or environment variable CCNP_REGISTRY
#
# @param $1 the name of container
#
function publish_a_image {
    echo "Publish container image: ${registry}/${1}:${tag} ..."
    docker push "${registry}/${1}:${tag}" || \
        { echo "Fail to push docker ${registry}/${1}:${tag}"; exit 1; }
    echo -e "Complete publish container image ${registry}/${1}:${tag} ...\n"
}

#
# Publish all container images
#
function publish_images {
    if [[ "$container" == "all" ]]; then
        for item in "${all_containers[@]}"
        do
            publish_a_image "$item"
        done
    else
        publish_a_image "$container"
    fi
}

#
# Save container image binary to a file, which can be restored later
#
function save_a_image {
    echo "Save container image ${registry}/${1}:${tag} => ${top_dir}/images/ ... "
    mkdir -p "${top_dir}"/images/
    docker save -o "${top_dir}/images/${1}-${tag}".tar "${registry}/${1}:${tag}"
    docker save "${registry}/${1}:${tag}" | gzip > "${top_dir}/images/${1}-${tag}.tgz"

    #
    # Please use following command to restore
    #
    # gunzip -c mycontainer.tgz | docker load
    # or
    # gunzip -c mycontainer.tgz
    # sudo ctr images import mycontainer.tar
    #
}

function save_images {
    if [[ "$container" == "all" ]]; then
        for item in "${all_containers[@]}"
        do
            save_a_image "$item"
        done
    else
        save_a_image "$container"
    fi
}

scan_all_containers "${curr_dir}"
process_args "$@"
echo ""
echo "-------------------------"
echo "action: ${action}"
echo "container: ${container}"
echo "tag: ${tag}"
echo "registry: ${registry}"
echo "-------------------------"
echo ""

if [[ "$action" =~ ^(build|all) ]]; then
    build_images
fi

if [[ "$action" =~ ^(publish|all) ]]; then
    publish_images
fi

if [[ "$action" =~ ^(save) ]]; then
    save_images
fi
