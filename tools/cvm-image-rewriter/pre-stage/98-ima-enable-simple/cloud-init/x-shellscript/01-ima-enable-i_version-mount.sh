#!/bin/bash

echo "=========== enable i_version mount started ==========="

FSTAB_FILE=/etc/fstab
FSTAB_FILE_NEW=$FSTAB_FILE".new"
awk '{if($3 == "ext4" && $4 !~ /iversion/) $4 = $4",iversion"; print}' \
    $FSTAB_FILE \
    > $FSTAB_FILE_NEW
mv $FSTAB_FILE_NEW $FSTAB_FILE

GRUB_FILE=/etc/default/grub.d/50-cloudimg-settings.cfg
GRUB_FILE_NEW=$GRUB_FILE".new"
awk '{if($0 !~ /i_version/ && $1 ~ /GRUB_CMDLINE_LINUX_DEFAULT/) sub("=\"","=\"rootflags=i_version "); print}' \
    $GRUB_FILE \
    > $GRUB_FILE_NEW
mv $GRUB_FILE_NEW $GRUB_FILE

echo "=========== enable i_version mount finished ==========="
# update-grub