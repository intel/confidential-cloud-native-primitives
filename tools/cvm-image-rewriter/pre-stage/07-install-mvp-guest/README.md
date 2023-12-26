# Install MVP TDX guest

This plugin is used to install MVP TDX guest kernel from a given local repo.

## Prerequisite
- Please put MVP guest repo into ./artifacts
- Set CVM_MVP_GUEST_REPO to the repo relative path, for example 2023ww45/mvp-tdx-stack-guest-ubuntu-22.04, guest_repo
- If the original image is smaller then 1.5G, please also set the environment variable GUEST_SIZE larger.

## Example output

```
CVM_MVP_GUEST_REPO: 2023ww45/mvp-tdx-stack-guest-ubuntu-22.04
MVP guest local repo ./artifacts/2023ww45/mvp-tdx-stack-guest-ubuntu-22.04 check passed
SUCCESS: MVP guest local repo ./artifacts/2023ww45/mvp-tdx-stack-guest-ubuntu-22.04 copied to guest /srv
SUCCESS: Cloud config cloud-init/x-shellscript/07-install-mvp-guest.sh generated
```
