# Confidential VM Customization Tool

This tool is used to customize the confidential VM guest including guest image,
config, OVMF firmware etc.

## 1. Overview

The confidential VM guest can be customized as follows:

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

### 2.1 Existing Plugins

There are following customization plugins in Plugins providing customization to base image.

| Name | Descriptions | Required for CCNP deployment |
| ---- | ------------ | ------------ |
| 01-resize-image | Resize the input qcow2 image | N |
| 02-motd-welcome | Customize the login welcome message | N |
| 03-netplan | Customize the netplan.yaml | N |
| 04-user-authkey | Add auth key for user login instead of password | N |
| 05-readonly-data | Fix some file permission to ready-only | N |
| 06-install-tdx-guest-kernel | Install MVP TDX guest kernel | Y |
| 07-device-permission | Fix the permission for device node | Y |
| 08-ccnp-uds-directory-permission | Fix the permission for CCNP UDS directory | Y |
| 60-initrd-update | Update the initrd image | N |
| 97-sample | plugin customization example | N |
| 98-ima-enable-simple | Enable IMA (Integrity Measurement Architecture) feature | N |

### 2.2 Design a new plugin

A plugin is put into the directory of [`plugins`](/tools/cvm-image-rewriter/plugins/),
with the number as directory name's prefix. So the execution of plugin will be
dispatched according to number sequence for example `99-byebye` is the final one.

A plugin includes several customization approaches:

1. File override: all files under `<plugin directory>/files` will be copied into the
corresponding directory in the target guest image.
2. Pre-stage execution on the host: the `<plugin directory>/pre-stage/host_run.sh` will be
executed before cloud-init stage.
3. cloud-init customization: please put the config yaml in `<plugin directory>/cloud-init/cloud-config`,
and put the scripts in `<plugin directory>/cloud-init/x-shellscript`.

Please refer to [the sample plugin](/tools/cvm-image-rewriter/plugins/97-sample/).

## 3. How to Run the tool

### 3.1 Prerequisite

1. This tool has been tested on `Ubuntu 22.04` and `Debian 10`. It is recommend to use
`Ubuntu 22.04`.

2. This tool can run on bare metal or within a virtual machine using nesting as detailed in [Section 3.4](#3.4-Run-in-Nested-VM-(Optional)).

3. Please install the following packages on Ubuntu/Debian.

    ```
    sudo apt install qemu-utils guestfs-tools virtinst genisoimage libvirt-daemon-system libvirt-daemon
    ```
    If `guestfs-tools` is not available in your distribution, you may need to install some additional packages on Debian.

    ```
    sudo apt-get install guestfsd libguestfs-tools
    ```

4. Ensure current login user is in the group of libvirt.

    ```
    sudo usermod -aG libvirt $USER
    ```

5. Ensure read permission on `/boot/vmlinuz-$(uname-r)`.

    ```
    sudo chmod o+r /boot/vmlinuz-*
    ```

6. The version of cloud-init is required > 23.0, so if the host distro could not
provide such cloud-init tool, you have to install it manually. For example, on a
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
    please make sure the current user belongs to the libvirt group via the following commands:
    ```
    sudo usermod -aG libvirt $USER
    sudo systemctl daemon-reload
    sudo systemctl restart libvirtd
    ```

8. Please start the net `default` for libvirt via:

    ```
    virsh net-start default
    ```

### 3.2 Run the tool

The tool provides several plugins to customize the initial image. It will generate an `output.qcow2` under current directory.

Before running the tool, please choose the plugins that are needed.You can skip any plugin by creating a file "NOT_RUN" under the current directory.
For example:

    ```
    touch plugins/01-resize-image/NOT_RUN
    ```

If the guest image is used for CCNP deployment, it's recommended to run the below plugin combination depending on which guest image type is used.
Others are not required by CCNP and can be skipped.
|  Base image | 01  | 02  | 03  | 04  | 05 | 06 | 07 | 08 | 60 | 98 |
|---|---|---|---|---|---|---|---|---|---|---|
|  Ubuntu base image | | | | | | Y| Y| Y| | |
| TD enlightened image | | | | | | | Y| Y| | |

**NOTE:**
  - All plugins need to be executed in numerical order.
  - TD enlightened image means the image already has a TDX kernel. If not, plugin 06 is required to install a TDX kernel.
  - Plugin 7 and Plugin 8 need to be executed before deploying CCNP to provide device permissions for CCNP.
  - Plugin 60 requires copying or generating all files to the root directory first. When users customize plugins, please ensure that the plugin number with this requirement is placed before 60.
  - Plugin 98 needs to be executed after all other plugins have completed. The number of the user-customized plugin must be before 98.
  - Other plugins are optional for CCNP deployment. 

The tool supports parameters as below.
```
$ ./run.sh -h
Usage: run.sh [OPTION]...
Required
  -i <guest image>          Specify initial guest image file
Optional
  -t <number of minutes>    Specify the timeout of rewriting, 3 minutes default,
                            If enabling IMA, recommend timeout >6 minutes
  -s <connection socket>    Default connection URI is qemu:///system,
                            if install libvirt, you can specify to "/var/run/libvirt/libvirt-sock"
                            then the corresponding URI is "qemu+unix:///system?socket=/var/run/libvirt/libvirt-sock"
  -n                        Silent running for virt-install with no output
  -h                        Show usage
```

For example:
```
# Run the tool with an initial guest image and set timeout as 10 minutes.
$ ./run.sh -i <initial guest image> -t 10
```

### 3.3 Boot a VM

After above tool is running successfully, you can boot a VM using the generated `output.qcow2` using `qemu-test.sh` or `start-virt.sh`.

- Boot TD or normal VM using `qemu-test.sh`.
  ```
  $ sudo ./qemu-test.sh -h
  Usage: qemu-test.sh [OPTION]...
  Required
    -i <guest image>          Specify initial guest image file
  ```

  For example:
  ```
  # Boot a TD
  $ sudo ./qemu-test.sh -i output.qcow2 -t td -p <qemu monitor port> -f <ssh_forward port>

  # Boot a normal VM
  $ sudo ./qemu-test.sh -i output.qcow2 -p <qemu monitor port> -f <ssh_forward port>
  ```
- Boot TD using `start-virt.sh`.

  ```
  $ sudo ./start-virt.sh -h
  Usage: start-virt.sh [OPTION]...
    -i <guest image file>     Default is tdx-guest-ubuntu22.04.qcow2 under current directory
    -n <guest name>           Name of TD guest
    -t <template file>        Default is ./tdx-libvirt-ubuntu-host.xml.template
    -f                        Force recreate
    -v <vcpu number>          VM vCPU number
    -m <memory size in GB>    VM memory size in GB
    -h                        Show this help
  ```

  For example:
  ```
  # Boot a TD with specified name and CPU/memory
  $ sudo ./qemu-test.sh -i output.qcow2 -n <libvirt domain name> -v <vCPU number> -m <memory size in GiB>
  ```

### 3.4 Run in Nested VM (Optional)

This tool can also be run in a guest VM on the host, in case that users need to prepare a clean host environment.  

1. Enable Nested Virtualization

Given that some plugins will consume more time in a low-performance guest VM, it is recommended to enable nested virtualization feature on the host.

First, check if the nested virtualization is enabled. If the file `/sys/module/kvm_intel/parameters/nested` show `Y` or `1`, it indicates that the feature is enabled. 

```
cat /sys/module/kvm_intel/parameters/nested
```

If the feature is not enabled, create the file ` /etc/modprobe.d/kvm.conf`, appending `options kvm_intel nested=1` to it and reboot the host.

```
echo "options kvm_intel nested=1" > /etc/modprobe.d/kvm.conf
```

2. Launch the guest VM

When we launch the guest VM, it is recommended to allocate more than `8G` memory for the guest VM, because this tool will occupy at least `4G` memory. And more CPU cores will improve the guest VM performance, typically the number of CPU cores is at least `4`.

3. Install dependencies

At last, install dependencies in the guest VM before running this tools.

It is an example for a basic Ubuntu 22.04 guest VM.

```
sudo apt install qemu-utils libguestfs-tools virtinst genisoimage cloud-init qemu-kvm libvirt-daemon-system
```
