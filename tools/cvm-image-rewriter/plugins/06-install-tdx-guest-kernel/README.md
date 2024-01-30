# Install TDX guest kernel

This plugin is used to install a TDX guest kernel from a given local repository.

# Prerequisite

Prepare the local repository and confirm that there are Debian packages related to the TDX kernel in the `/jammy/amd64/` directory of this repository. It is recommended to place this local repository in the `pre-stage/artifacts/` directory.
```
mkdir -p ./pre-stage/artifacts
mv <your guest repo> ./pre-stage/artifacts/
```

Set `${CVM_TDX_GUEST_REPO}` to the repository absolute path, or this plugin will be skipped. 
```
export CVM_TDX_GUEST_REPO=$(pwd)/pre-stage/artifacts/<your guest repo>

# Or
export CVM_TDX_GUEST_REPO=<your local guest repo>
```


_NOTE: IF the original image is smaller than 1.5G, please set the environment variable GUEST\_SIZE to a larger value, as this will result in the execution of plugin 01._
