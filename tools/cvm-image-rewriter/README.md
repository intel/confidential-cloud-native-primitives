# Confidential VM Customization Tool

This tool is used to customize the confidential VM guest including guest image,
config, OVMF firmware etc.

## 1. Overview
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

### 2.1 Customize

```
$ ./run.sh -h
Usage: run.sh [OPTION]...
Required
  -i <guest image>          Specify initial guest image file
Optional
  -t <number of minutes>    Specify the timeout of rewriting, 3 minutes default,
                            If enabling ima, recommend timeout >6 minutes
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
| 60-initrd-update | Update the initrd image |
| 98-ima-enable-simple | Enable IMA (Integrity Measurement Architecture) feature |

### 3.1 Design a new plugin

TBD