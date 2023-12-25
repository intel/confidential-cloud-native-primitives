# IMA enable simple

IMA is short Integrity Measurement Architecture. This plugin is a simple demo to enable IMA feature in the guest image, which can be regarded as a reference implementation. Please create you own IMA plugin in the production env.

In this case, we have following assumption
- users have their customized IMA policy, instead of built-in policy. 
- only option `ima_appraise` is included.
- feature i_version mount is enabled.

More info and options related to IMA can be found in Ref[1]. 

Let us go through this plugin and help you write a customized one.

## Customized IMA policy

The policy MUST be placed in `files/etc/ima/` and named by `ima-policy`, otherwise it can not be recognized and loaded into the guest kernel. The policy file will be copied to `/etc/ima/ima-policy` in the guest image.

This plugin uses the policy same as the built-in policy `tcb_appraise`. Other built-in policies can be found in Ref[2].

It is recommend to use a customized policy other than those default ones in production env.

## Feature IMA enabling

To enable IMA, this plugin conducts following 2 steps.

1. Enabling IMA update

To enable IMA update, the script `host_run.sh` in the host invokes script `guest_enable_ima_fix.sh`, which will be executed in the guest image, by `virt-customize`.

The script `guest_enable_ima_fix.sh` appends option `ima_appraise=fix` in config file `/etc/default/grub.d/50-cloudimg-settings.cfg` and run `update-grub` to update kernel cmdline. If you guest image has other config files, please replace the config file path. 

Here, the IMA update is enabled, guest image can update IMA data at next launch, which typically refers to the next launch by the cloud-init in the following step.

2. Registering files hashes and enabling i_version mount

These two purposes will be completed by scripts `01-ima-enable-i_version-mount.sh` and `02-ima-register-file-hash.sh` in `cloud-init/x-shellscript`. The two-digital prefix of the script is necessary and indicates the execution order, that is to say, the script with smaller prefix will be executed earlier.

Script `01-ima-enable-i_version-mount.sh` registers files' hashes and update option `ima_appraise=fix` to `ima_appraise=enforce` in the kernel cmdline, to enable IMA at next launch.

Script `02-ima-register-file-hash.sh` updates `/etc/fstab`, where we only enable i_version on partitions with ext4 format in the case, and append option `rootflags=i_version` in kernel cmdline to enable i_version mount on root partition.

These two scripts will not be invoked by us. They will be integrated into `user-data` with other config or scripts for `cloud-init`. Then the `cloud-init` will copy them into the guest image and execute they in order at most once.

At last, the final `update-grub` is placed in `tools/cvm-image-rewriter/cloud-init/user-data.basic` to avoid redundant execution. Therefore, we comment the `update-grub` in the `01-ima-enable-i_version-mount.sh`.

Ref:
- [1] https://sourceforge.net/p/linux-IMA/wiki/Home/#enabling-IMA-measurement
- [2] https://wiki.gentoo.org/wiki/Integrity_Measurement_Architecture/Recipes
- [3] https://wiki.gentoo.org/wiki/Integrity_Measurement_Architecture
