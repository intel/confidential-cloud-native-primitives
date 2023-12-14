#!/bin/bash

echo "=========== register file hash started ==========="

GRUB_FILE=/etc/default/grub.d/50-cloudimg-settings.cfg

time find / -path /proc -prune -o -fstype ext4 -type f -uid 0 -exec dd if='{}' of=/dev/null count=0 status=none \;

sed -i 's/ima_appraise=fix/ima_appraise=enforce/' $GRUB_FILE

echo "=========== register file hash finished ==========="

# the command is executed in basic script
# update-grub 