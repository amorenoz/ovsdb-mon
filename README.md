# ovsdb-mon

ovsdb-mon is an OVSDB monitoring and visulization tool based on [libovsdb](https://github.com/ovn-org/libovsdb)

### Building ovsdb-mon
A common usage for this tool is to monitor an OVN database. However, it is generic enough to be used with
any OVSDB schema. By default, ovsdb-mon uses the schema defined in *schemas/ovn-nb.ovsschema*. If you want to
use your own, simply download it from your ovsdb server and replace the existing one

     ovsdb-client get-schema ${SERVER} ${DATABASE} > schemas/my.ovsschema

Then, just build the program specifying a schema file

     make SCHEMA=schemas/my.ovsschema

This will use modelgen to generate a native model of the DB and use it to build ovsdb-mon

### Using ovsdb-mon
Usage of ovsdb-mon:

	./bin/ovsdb-mon [FLAGS] [COMMAND]
	FLAGS:
	  -auto
	        Autostart: If set to true, it will start monitoring from the begining
	  -db string
	        Database connection. Default: unix:/${OVS_RUNDIR}/ovnnb_db.sock
	  -monitor string
	        Only monitor these comma-separated tables
	  -no-monitor string
	        Do not monitor these comma-separated tables
	COMMAND:
	        If provided, it will run the command and exit. If not, it will enter interactive mode
	        For a full description of available commands use the command 'help'


By default it will open an interactive terminal where you can monitor the activity of the DB and inspect it:

	./bin/ovsdb-mon  --db tcp:172.19.0.4:6641
	OVN Monitoring Shell
	>>> help
	Commands:
	  clear      clear the screen
	  exit       exit the program
	  help       display help
	  list       List the content of a specific table
	  save       Save events
	  show       Print available tables
	  start      Start monitoring activity of the OVN DB
	  stop       Stop monitoring activity of the OVN DB


The result looks like this:

![Demo](doc/images/demo.gif)

### Kubernetes

Use the yaml and scripts provided in the dist folder in order to deploy pods that
provide a ready to use binary for the K8 cluster.

	[ -n ${KUBECONFIG} ] || echo Make sure kubectl command can reach the cluster
    cd dist
    source ./ovsdb-mon-ovn.source
    source ./ovsdb-mon-ovs.source

