#!/bin/bash
set -ex

if [ "${1:0:1}" = '-' ] || [ "${1:0:1}" = '' ]; then
	exec vote-server
fi

exec $@
