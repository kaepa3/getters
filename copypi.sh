#!/bin/sh
if [ $# -lt 1 ]; then
    echo arg error: $*
    exit 1
fi
chmod 777 getters-linux-arm-5
scp -p getters-linux-arm-5 pi@$1:/home/pi/priser/
