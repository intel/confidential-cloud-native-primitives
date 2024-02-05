#!/bin/bash

DIR=$(dirname "$(readlink -f "$0")")
FILE_LIST="$DIR/file_list"
CLD_SH_REGISTER_FILE_HASH="01-ima-register-file-hash.sh"
CLD_SH="$DIR/../cloud-init/x-shellscript/$CLD_SH_REGISTER_FILE_HASH"
CLD_SH_TEMPLATE=""
injects=""

read -r -d '' CLD_SH_TEMPLATE << EOM
#!/bin/bash

echo "=========== register file hash started ==========="

GRUB_FILE=/etc/default/grub.d/50-cloudimg-settings.cfg

# replaced by required files
PLACEHOLDER

sed -i 's/ima_appraise=fix/ima_appraise=log/' \$GRUB_FILE
sed -i 's/console=tty1 console=ttyS0/console=hvc0/' \$GRUB_FILE

echo "=========== register file hash finished ==========="

# the command is executed in basic script
# update-grub 
EOM

# filter specified files & directories
while IFS= read -r line || [ -n "$line" ]; do
    if [[ $line == "#"* ]] || [[ $line == "" ]]; then
        continue
    fi
    if [[ $line == *"/" ]]; then
        injects+="find $line -type f -uid 0 -exec dd if='{}' of=/dev/null count=0 status=none \\\;""\n"
    else
        name=$(basename "$line")
        path=${line%/*}/
        injects+="find $path -type f -name $name -uid 0 -exec dd if='{}' of=/dev/null count=0 status=none \\\;""\n"
    fi
done <"$FILE_LIST"

# generate script for cloud-init
# shellcheck disable=SC2001
echo "$CLD_SH_TEMPLATE" | sed -e "s@PLACEHOLDER@$injects@g" > "$CLD_SH"

# enable ima update
virt-customize -a "${GUEST_IMG}" \
    --run-command 'mkdir -p /etc/ima' \
    --run "$DIR"/guest_enable_ima_fix.sh

