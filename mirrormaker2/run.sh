# !/usr/bin/env bash

function die() {
    for i in "${@}"; do
        echo -e "${i}"
    done
    exit 1
}

[[ -z "${MM2_TARGET}" ]] && die "No path has been set for mm2.properties"

if [[ -n "${MM2_SOURCE}" ]] && [[ -f "${MM2_SOURCE}" ]]; then
    mm2-init -source  "${MM2_SOURCE}" -target "${MM2_TARGET}"
else
    mm2-init -target "${MM2_TARGET}"
fi

# Start kafka mirror maker.
/opt/kafka/bin/connect-mirror-maker.sh "${MM2_TARGET}"

