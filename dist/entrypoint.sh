#!/bin/bash
set -ex

export OVS_RUNDIR=${OVS_RUNDIR:-/run/ovn}

git clone https://github.com/amorenoz/ovsdb-mon

pushd ovsdb-mon

declare -A db_schemas=( ["OVN_Northbound"]="ovnnb_db.sock" ["OVN_Southbound"]="ovnsb_db.sock" ["Open_vSwitch"]="db.sock")

for k in "${!db_schemas[@]}"; do
    #The following trick makes it work on upstream ovn-kubernetes as it uses different paths, it does nothing on openshift
    ln -s /run/openvswitch/${db_schemas[${k}]} ${OVS_RUNDIR}/${db_schemas[${k}]} 2> /dev/null || true
    if [ -e "${OVS_RUNDIR}/${db_schemas[${k}]}" ]; then
        ovsdb-client get-schema "unix:${OVS_RUNDIR}/${db_schemas[${k}]}" ${k} > ${k}.schema
        SCHEMA=${k}.schema make build
        mv -v ./bin/ovsdb-mon /usr/local/bin/ovsdb-mon.${k}
    fi
done
touch /tmp/build_finished

popd

trap : TERM INT; sleep infinity & wait
