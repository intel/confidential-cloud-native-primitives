#!/bin/bash

DIR=$(dirname "$(readlink -f "$0")")
FILE_LIST="$DIR/file_list"
CLD_SH_READONLY_FILE="01-file-readonly.sh"
CLD_SH="$DIR/../cloud-init/x-shellscript/$CLD_SH_READONLY_FILE"
CLD_SH_TEMPLATE=""
injects=""

read -r -d '' CLD_SH_TEMPLATE << EOM
#!/bin/bash
# replaced by required files
PLACEHOLDER
EOM

# filter specified files & directories
while IFS= read -r line || [ -n "$line" ]; do
    if [[ $line == "#"* ]] || [[ $line == "" ]]; then
        continue
    fi
    if [[ $line == *"/" ]]; then
        injects+="chmod -R a=r $line""\n"
    else
        injects+="chmod a=r $line""\n"
    fi
done <"$FILE_LIST"

mkdir -p "$DIR/../cloud-init/x-shellscript"
# shellcheck disable=SC2001
echo "$CLD_SH_TEMPLATE" | sed -e "s@PLACEHOLDER@$injects@g" > "$CLD_SH"