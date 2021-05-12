package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	model "github.com/amorenoz/ovnmodel/model"
	"github.com/ovn-org/libovsdb/client"
	//"github.com/fatih/color"
)

const (
	ovnnbSocket = "ovnnb_db.sock"
)

var (
	ovs  *client.OvsdbClient
	db   = flag.String("db", "", "Database connection. Default: unix:/${OVS_RUNDIR}/ovnnb_db.sock")
	auto = flag.Bool("auto", false, "Autostart: If set to true, it will start monitoring from the begining")
)

type ormSignal struct{}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s [FLAGS] [COMMAND] \n", os.Args[0])
		fmt.Fprintf(os.Stderr, "FLAGS:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "COMMAND:\n")
		fmt.Fprintf(os.Stderr, "\tIf provided, it will run the command and exit. If not, it will enter interactive mode\n")
		fmt.Fprintf(os.Stderr, "\tFor a full description of available commands use the command 'help'")
	}
	flag.Parse()

	var addr string
	if *db != "" {
		addr = *db
	} else {
		var ovs_rundir = os.Getenv("OVS_RUNDIR")
		if ovs_rundir == "" {
			ovs_rundir = "/var/run/openvswitch"
		}
		addr = "unix:" + ovs_rundir + "/" + ovnnbSocket
	}

	dbModel, err := model.FullDatabaseModel()
	if err != nil {
		log.Fatal(err)
	}

	shell := newOvnShell(*auto, dbModel)
	ovs, err := client.Connect(addr, dbModel, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ovs.Disconnect()
	shell.Run(ovs, flag.Args()...)
}
