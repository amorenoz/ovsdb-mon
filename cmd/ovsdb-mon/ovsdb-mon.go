package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	model "github.com/amorenoz/ovsdb-mon/model"
	"github.com/ovn-org/libovsdb/client"
	//"github.com/fatih/color"
)

const (
	ovnnbSocket = "ovnnb_db.sock"
)

var (
	db              = flag.String("db", "", "Database connection. Default: unix:/${OVS_RUNDIR}/ovnnb_db.sock")
	auto            = flag.Bool("auto", false, "Autostart: If set to true, it will start monitoring from the beginning")
	monitorTables   = flag.String("monitor", "", "Only monitor these comma-separated tables")
	noMonitorTables = flag.String("no-monitor", "", "Do not monitor these comma-separated tables")
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

	tablesToMonitor, err := getTablesToMonitor(dbModel, *monitorTables, *noMonitorTables)
	if err != nil {
		log.Fatal(err)
	}

	shell := newOvsdbShell(*auto, dbModel, tablesToMonitor)
	c, err := client.NewOVSDBClient(dbModel, client.WithEndpoint(addr))
	if err != nil {
		log.Fatal(err)
	}
	err = c.Connect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	defer c.Disconnect()
	shell.Run(c, flag.Args()...)
}
