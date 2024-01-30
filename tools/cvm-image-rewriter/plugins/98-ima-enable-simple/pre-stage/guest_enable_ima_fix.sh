#!/bin/sh

FSTAB_FILE=/etc/fstab
FSTAB_FILE_NEW=$FSTAB_FILE".new"
awk '{if($3 == "ext4" && $4 !~ /iversion/) $4 = $4",iversion"; print}' \
    $FSTAB_FILE \
    > $FSTAB_FILE_NEW
mv $FSTAB_FILE_NEW $FSTAB_FILE

GRUB_FILE=/etc/default/grub.d/50-cloudimg-settings.cfg
GRUB_FILE_NEW=$GRUB_FILE".new"

# Remove ima_appraise and rootflags if exists
sed -i 's/ima_appraise=\(fix\|enforce\|log\|off\)//' $GRUB_FILE
sed -i 's/ima_hash=sha384//' $GRUB_FILE
sed -i 's/rootflags=i_version//' $GRUB_FILE
sed -i 's/console=tty1 console=ttyS0//' $GRUB_FILE
sed -i 's/console=hvc0//' $GRUB_FILE 

awk '{if($1 ~ /GRUB_CMDLINE_LINUX_DEFAULT/) sub("=\"","=\"console=tty1 console=ttyS0 ima_appraise=fix ima_hash=sha384 rootflags=i_version "); print}' \
    $GRUB_FILE \
    > $GRUB_FILE_NEW
mv $GRUB_FILE_NEW $GRUB_FILE

update-grub
