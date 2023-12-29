# Confidential VM Customization Tool

This tool is used to customize the confidential VM guest including guest image,
config, OVMF firmware etc.

## 1. Overview

The confidential VM guest can be customized including follows:

![](/docs/cvm-customizations.png)

| Name | Type/Scope | Description |
| ---- | ---------- | ----------- |
| Launch Identity | Config | MROwner, MRConfig, MROwnerConfig |
| VM Configuration | Config | vCPU, memory, network config |
| Secure Boot Key | OVMF | the PK/DB/KEK for secure boot or Linux MoK |
| Config Variable | OVMF | the configurations in variable |
| Grub | Boot Loader | Grub kernel command, Grub modules |
| initrd | Boot Loader | Customize build-in binaries |
| IMA Policy | OS | Policy via loading systemd |
| Root File System | OS | RootFS customization |

## 2. Design

It is based on the [cloud-init](https://cloudinit.readthedocs.io/en/latest/)
framework, and the whole flow was divided into three stages:

- **Pre Stage**: prepare to run cloud-init. It will collect the files for target
  image, meta-data/x-shellscript/user-data for cloud-init's input.
- **Cloud-Init Stage**: it will run cloud init in sequences of
  - Generate meta files via `cloud-init make-mime`
  - Generate `ciiso.iso` via `genisoimage`
  - Run cloud-init via `virt-install`
- **Post Stage**: clean up and run post check

![](/docs/cvm-image-rewriter-flow.png)

## 2. Run

### 2.1 Prerequisite

1. This tool has been tested on `Ubuntu 22.04` and `Debian 10`. It is recommend to use
`Ubuntu 22.04`.

2. This tool can run on bare metal or virtual machine (with nest VM like `Intel VT-x`)

3. Please install following packages on Ubuntu/Debian:

    ```
    sudo apt install qemu-utils guestfs-tools virtinst genisoimage libvirt-daemon-system libvirt-daemon
    ```

4. Ensure current login user is in the group of libvirt

    ```
    sudo usermod -aG libvirt $USER
    ```

5. Ensure read permission on `/boot/vmlinuz-$(uname-r)`.

    ```
    sudo chmod o+r /boot/vmlinuz-*
    ```

6. The version of cloud-init is required > 23.0, so if the host distro could not
provide such cloud-init tool, you have to install by manual. For example, on a
debian 10 system, the version of default cloud-init is 20.0. Please do following
steps:
    ```
    wget http://ftp.cn.debian.org/debian/pool/main/c/cloud-init/cloud-init_23.3.1-1_all.deb
    sudo dpkg -i cloud-init_23.3.1-1_all.deb
    ```

7. If it is running with `libvirt/virt-daemon` hypervisor, then:

  - In file `/etc/libvirt/qemu.conf`, make sure `user` and `group` is `root` or
    current user.
  - If need customize the connection URL, you can specify via `-s` like `-s /var/run/libvirt/libvirt-sock`,
    please make sure current user belong to libvirt group via following commands:
    ```
    sudo usermod -aG libvirt $USER
    sudo systemctl daemon-reload
    sudo systemctl restart libvirtd
    ```

8. Please start the net `default` for libvirt via:

    ```
    virsh net-start default
    ```


### 2.1 Customize

```
$ ./run.sh -h
Usage: run.sh [OPTION]...
Required
  -i <guest image>          Specify initial guest image file
Optional
  -t <number of minutes>    Specify the timeout of rewriting, 3 minutes default,
                            If enabling ima, recommend timeout >6 minutes
  -s <connection socket>    Default is connection URI is qemu:///system,
                            if install libvirt, you can specify to "/var/run/libvirt/libvirt-sock"
                            then the corresponding URI is "qemu+unix:///system?socket=/var/run/libvirt/libvirt-sock"
  -n                        Silence running for virt-install, no output
```

**_NOTE_**:

- If want to skip to run specific plugins at `pre-stage` directory, please create
a file named as `NOT_RUN` at the plugin directory. For example:
    ```
    touch pre-stage/01-resize-image/NOT_RUN
    ```


### 2.2 Run Test

```
$ ./qemu-test.sh -h
Usage: qemu-test.sh [OPTION]...
Required
  -i <guest image>          Specify initial guest image file
```

## 3. Plugin

### 3.1 Existing Plugins

There are following customization plugins in Pre-Stage:

| Name | Descriptions |
| ---- | ------------ |
| 01-resize-image | Resize the input qcow2 image |
| 02-motd-welcome | Customize the login welcome message |
| 03-netplan | Customize the netplan.yaml |
| 04-user-authkey | Add auth key for user login instead of password |
| 05-readonly-data | Fix some file permission to ready-only |
| 07-install-mvp-guest | Fix some file permission to ready-only |
| 08-device-permission | Fix the permission for device node |
| 09-ccnp-uds-directory-permission | Fix the permission for CCNP UDS directory |
| 60-initrd-update | Update the initrd image |
| 98-ima-enable-simple | Enable IMA (Integrity Measurement Architecture) feature |

### 3.1 Design a new plugin

A plugin is put into the directory of [`pre-stage`](/tools/cvm-image-rewriter/pre-stage/),
with the number as directory name's prefix. So the execution of plugin will be
dispatched according to number sequence for example `99-test` is the final one.

A plugin includes several customization approaches:

1. File override: all files under `<plugin directory>/files` will be copied the
corresponding directory in target guest image.
2. Pre-stage execution on the host: the `<plugin directory>/host_run.sh` will be
executed before cloud-init stage
3. cloud-init customization: please put the config yaml into `<plugin directory>/cloud-init/cloud-config`,
and put the scripts to `<plugin directory>/cloud-init/x-shellscript`

Please refer [the sample plugin](/tools/cvm-image-rewriter/pre-stage/99-test/).
