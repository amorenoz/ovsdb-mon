package main

import (
	"sync"

	goovn "github.com/ebay/go-ovn"
	"github.com/k0kubun/pp"
	ishell "gopkg.in/abiosoft/ishell.v2"
)

type OvnShell struct {
	mutex   *sync.RWMutex
	monitor bool
	orm     goovn.ORMClient
}

func (s *OvnShell) Monitor(monitor bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.monitor = monitor
}

func (s *OvnShell) OnCreated(m goovn.Model) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if s.monitor {
		pp.Printf("A %s has been added\n", m.Table())
		pp.Println(m)
		pp.Println("")
	}
}

func (s *OvnShell) OnDeleted(m goovn.Model) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if s.monitor {
		pp.Printf("A %s has been added\n", m.Table())
		pp.Println(m)
		pp.Println("")
	}
}

func (s *OvnShell) Run(orm goovn.ORMClient) {
	s.orm = orm
	shell := ishell.New()
	shell.Set("ovnShell", s)

	shell.Println("OVN Monitoring Shell")
	shell.AddCmd(&ishell.Cmd{
		Name: "start",
		Help: "Start monitoring activity of the OVN DB",
		Func: func(c *ishell.Context) {
			ovnShell := c.Get("ovnShell")
			if ovnShell == nil {
				c.Println("Error: No context")
			}
			ovnShell.(*OvnShell).Monitor(true)
		},
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "stop",
		Help: "Stop monitoring activity of the OVN DB",
		Func: func(c *ishell.Context) {
			ovnShell := c.Get("ovnShell")
			if ovnShell == nil {
				c.Println("Error: No context")
			}
			ovnShell.(*OvnShell).Monitor(false)
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "show",
		Help: "Print available tables",
		Func: func(c *ishell.Context) {
			ovnShell := c.Get("ovnShell")
			if ovnShell == nil {
				c.Println("Error: No context")
			}
			c.Println("Available Tables")
			c.Println("----------------")

			for name := range ovnShell.(*OvnShell).orm.GetSchema().Tables {
				c.Println(name)
			}
		},
	})
	shell.Run()
}

func newOvnShell(auto bool) *OvnShell {
	return &OvnShell{
		mutex:   new(sync.RWMutex),
		monitor: auto,
	}
}
