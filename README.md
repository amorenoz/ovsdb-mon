# ovnmodel

Experiment with native OVN models. It is based on yet-to-merge branches of [go-ovn](https://github.com/ebay/go-ovn) and [libovsdb](https://github.com/eBay/libovsdb)

This repository contains two programs:
- **modelgen**: Generates go-ovn-compatible bindings given a ovsschema file
- **ovnmon**:  Uses the models generated by *modelgen* to monitor a ovn database

The whole idea of this project is to test / demonstrate the goodies of the improvements that are being developed in go-ovn and libovsdb.

## How to use it
Generate the native go-ovn model:

    make model

This will use the *ovn-nb.ovsschema* provided in the repo. If you want to use your own schema run:

     make modelgen
     ovsdb-client get-schema ${SERVER} ${DATABASE} > myDB.ovsschema
     ./bin/modelgen -o model -p model myDB.ovsschema

Build the simple monitoring app:

	 make

Run the app against the database:

	./bin/ovnmon --db tcp:172.18.0.4:6641

The result looks like this:

            A "Logical_Switch" has been added!
           &model.LogicalSwitch{
             UUID: "421d4947-0ea0-4f00-b412-b11f46664e4b",
             Acls: []string{
               "3bc7d4ff-955d-4abd-843a-a46684c053bd",
             },
             ExternalIds:  map[string]string{},
             LoadBalancer: []string{
               "243fc1d3-2eef-4da2-8853-d57389e63881",
               "5d5f5d9c-0c28-4565-b288-747d69058f9f",
               "7ef402bf-a4b3-4f5e-8faa-9b4ee8e30578",
               "803a08dd-9fa5-4d7d-8f79-bdb36e7fe007",
               "9c70db1c-4254-4825-ad1a-969d32fd860e",
               "c3c7d07e-ee86-4b73-abc6-665a83f791c2",
             },
             DnsRecords:       []string{},
             ForwardingGroups: []string{},
             Name:             "ovn-control-plane",
             OtherConfig:      map[string]string{
               "subnet": "10.244.1.0/24",
             },
             Ports: []string{
               "680148d5-f99f-4330-9d30-937caae8318c",
               "f2d58d9c-9832-41b5-a30f-19710dcc7c1b",
             },
             QosRules: []string{},
           }
           A "Address_Set" has been added!
           &model.AddressSet{
             UUID:        "6e8efd17-8377-4ce7-907f-5d40fad3b845",
             Addresses:   []string{},
             ExternalIds: map[string]string{
               "name": "kube-public_v4",
             },
             Name: "a18363165982804349389",
           }

