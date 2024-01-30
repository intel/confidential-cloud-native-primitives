#!/bin/bash

DIR=$(dirname "$(readlink -f "$0")")
TEMPLATE="#!/bin/sh\n\n. /usr/share/initramfs-tools/hook-functions\n"

HOOKS_DIR="$DIR/../files/etc/initramfs-tools/hooks"
SCRIPT_NAME=initrd-custom-update.sh
SCRIPT_PATH=$HOOKS_DIR/$SCRIPT_NAME

mkdir -p "$HOOKS_DIR"
echo -e "$TEMPLATE" > "$SCRIPT_PATH"
chmod a+x "$SCRIPT_PATH"

mapfile -t files < <(find "$DIR/../files" -type f)
for f in "${files[@]}"
do
    if [[ $f == *$SCRIPT_NAME ]]; then
        continue
    fi
    path=${f#*/files}
    echo "copy_exec $path" >> "$SCRIPT_PATH"
done

