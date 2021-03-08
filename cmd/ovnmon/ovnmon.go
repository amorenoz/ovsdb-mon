package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	model "github.com/amorenoz/ovnmodel/model"
	goovn "github.com/ebay/go-ovn"
	//"github.com/fatih/color"
	"github.com/k0kubun/pp"
)

const (
	ovnnbSocket = "ovnnb_db.sock"
)

var (
	orm goovn.ORMClient
	db  = flag.String("db", "", "Database connection. Default: unix:/${OVS_RUNDIR}/ovnnb_db.sock")
)

type ormSignal struct{}

func (s ormSignal) OnCreated(m goovn.Model) {
	pp.Printf("A %s has been added!\n", m.Table())
	pp.Println(m)
	fmt.Printf("")
}

func (s ormSignal) OnDeleted(m goovn.Model) {
	pp.Printf("A %s has been removed!\n", m.Table())
	pp.Println(m)
	fmt.Printf("")
}

func main() {
	flag.Parse()
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		fmt.Printf("Got signal %s", sig)
		done <- true
	}()

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

	config := goovn.Config{
		Db:          goovn.DBNB,
		Addr:        addr,
		ORMSignalCB: ormSignal{},
		DBModel:     dbModel,
	}
	orm, err := goovn.NewORMClient(&config)
	if err != nil {
		panic(err)
	}
	defer orm.Close()
	fmt.Println("Waiting for signal or new Logical Routers")
	<-done
	fmt.Println("Exiting")
}
