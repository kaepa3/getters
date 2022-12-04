#!/bin/sh
if [ $# -lt 1 ]; then
    echo arg error: $*
    exit 1
fi

scp -p getters-linux-arm-5 pi@$1:/home/pi/priser/

 getters-linux-arm-7
 getters-linux-arm64
