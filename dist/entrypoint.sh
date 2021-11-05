#!/bin/bash
set -ex

export OVS_RUNDIR=${OVS_RUNDIR:-/run/ovn}

pushd ovsdb-mon

declare -A db_schemas=( ["OVN_Northbound"]="ovnnb_db.sock" ["OVN_Southbound"]="ovnsb_db.sock" ["Open_vSwitch"]="db.sock")

declare -a ovs_runpaths=("/run/openvswitch" "/run/ovn" "/var/run/openvswitch" "/var/lib/ovn/" "/var/lib/openvswitch/ovn")

for k in "${!db_schemas[@]}"; do
    for path in "${ovs_runpaths[@]}"; do
        if [ -e "${path}/${db_schemas[${k}]}" ]; then
            ovsdb-client get-schema "unix:${path}/${db_schemas[${k}]}" ${k} > ${k}.schema
            SCHEMA=${k}.schema make build
            mv -v ./bin/ovsdb-mon /usr/local/bin/bin-ovsdb-mon.${k}
            cat >/usr/local/bin/ovsdb-mon.${k} <<EOF 
#!/bin/sh
/usr/local/bin/bin-ovsdb-mon.${k} -db "unix:${path}/${db_schemas[${k}]}" \$@
EOF
            chmod +x /usr/local/bin/ovsdb-mon.${k}
            break
        fi
    done
done
touch /tmp/build_finished

popd

trap : TERM INT; sleep infinity & wait
