package main

import (
	"flag"
	"log"
	"os"

	model "github.com/amorenoz/ovnmodel/model"
	goovn "github.com/ebay/go-ovn"
	//"github.com/fatih/color"
)

const (
	ovnnbSocket = "ovnnb_db.sock"
)

var (
	orm  goovn.ORMClient
	db   = flag.String("db", "", "Database connection. Default: unix:/${OVS_RUNDIR}/ovnnb_db.sock")
	auto = flag.Bool("auto", false, "Autostart: If set to true, it will start monitoring from the begining")
)

type ormSignal struct{}

func main() {
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

	dbModel, err := model.DBModel()
	if err != nil {
		log.Fatal(err)
	}

	shell := newOvnShell(*auto)
	config := goovn.Config{
		Db:          goovn.DBNB,
		Addr:        addr,
		ORMSignalCB: shell,
		DBModel:     dbModel,
	}
	orm, err := goovn.NewORMClient(&config)
	if err != nil {
		panic(err)
	}
	defer orm.Close()
	shell.Run(orm)
}
