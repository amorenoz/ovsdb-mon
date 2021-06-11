package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	model "github.com/amorenoz/ovnmon/model"
	"github.com/ovn-org/libovsdb/client"
	//"github.com/fatih/color"
)

const (
	ovnnbSocket = "ovnnb_db.sock"
)

var (
	db   = flag.String("db", "", "Database connection. Default: unix:/${OVS_RUNDIR}/ovnnb_db.sock")
	auto = flag.Bool("auto", false, "Autostart: If set to true, it will start monitoring from the beginning")
)

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
		var ovsRundir = os.Getenv("OVS_RUNDIR")
		if ovsRundir == "" {
			ovsRundir = "/var/run/openvswitch"
		}
		addr = "unix:" + ovsRundir + "/" + ovnnbSocket
	}

	dbModel, err := model.FullDatabaseModel()
	if err != nil {
		log.Fatal(err)
	}

	shell := newOvnShell(*auto, dbModel)
	ovs, err := client.Connect(context.Background(), dbModel, client.WithEndpoint(addr))
	if err != nil {
		log.Fatal(err)
	}
	defer ovs.Disconnect()
	shell.Run(ovs, flag.Args()...)
}
