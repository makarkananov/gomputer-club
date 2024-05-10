#!/bin/sh
set -e

if [ -z "$1" ]; then
    echo "Usage: entrypoint.sh <input_file>"
    exit 1
fi

/gomputerClub "$1"
