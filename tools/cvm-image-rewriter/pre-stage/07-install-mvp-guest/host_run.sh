#!/bin/bash

# Go to this dir
pushd "$(dirname "$(readlink -f "$0")")" || exit 0

# shellcheck disable=SC1091
source ../../scripts/common.sh
ARTIFACTS_HOST=./artifacts
ARTIFACTS_GUEST=/srv

# CVM_MVP_GUEST_REPO examples: 2023ww45/mvp-tdx-stack-guest-ubuntu-22.04, guest_repo
if [[ -z "$CVM_MVP_GUEST_REPO" ]]; then
    warn "CVM_MVP_GUEST_REPO is not set, skip"
    warn "Please put MVP guest repo into $ARTIFACTS_HOST"
    warn "Then set CVM_MVP_GUEST_REPO to the repo relative path"
    exit 0
fi
info "CVM_MVP_GUEST_REPO: $CVM_MVP_GUEST_REPO"

# Check if the local repo exists
MVP_GUEST_REPO_LOCAL="$ARTIFACTS_HOST/$CVM_MVP_GUEST_REPO"
if [[ ! -d "$MVP_GUEST_REPO_LOCAL" ]]; then
    warn "MVP guest local repo $MVP_GUEST_REPO_LOCAL does not exist, skip"
    exit 0
fi
# Check if it is a valid MVP repo
if ! compgen -G "$MVP_GUEST_REPO_LOCAL/jammy/amd64/linux-image-*mvp*.deb"; then
    warn "$MVP_GUEST_REPO_LOCAL is invalid, skip"
    exit 0
fi
info "MVP guest local repo $MVP_GUEST_REPO_LOCAL check passed"

# Copy MVP local repo from host to guest
virt-copy-in -a "${GUEST_IMG}" "$MVP_GUEST_REPO_LOCAL" "$ARTIFACTS_GUEST"
ok "MVP guest local repo $MVP_GUEST_REPO_LOCAL copied to guest $ARTIFACTS_GUEST"

# Generate cloud-config
mkdir -p cloud-init/x-shellscript/
cat > cloud-init/x-shellscript/07-install-mvp-guest.sh << EOL
#!/bin/bash

PACKAGE_DIR="$ARTIFACTS_GUEST/$(basename "$CVM_MVP_GUEST_REPO")/jammy/"
pushd \$PACKAGE_DIR || exit 0
apt install ./amd64/linux-image-unsigned-*.deb ./amd64/linux-modules-*.deb \
        ./amd64/linux-headers-*.deb ./all/linux-headers-*.deb --allow-downgrades -y
popd || exit 0
EOL
ok "Cloud config cloud-init/x-shellscript/07-install-mvp-guest.sh generated"

popd || exit 0
