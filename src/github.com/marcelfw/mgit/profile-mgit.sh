#!/bin/sh

# mcd changes directory to a git repository
# "mcd ."        changes to the local .mgit root
# "mcd <name>"   changes t a global directory
function mcd()
{
    if [ "$1" == "." ]; then
        # use local .mgit root
        cd `mgit echo "{{ .Path }}" | head -1`
    else
        # use global shortcut
        cd `mgit -s global -name "$1" echo "{{ .Path }}" | tail -1`
    fi
}
