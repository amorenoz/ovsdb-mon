package main

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/k0kubun/pp"
	"github.com/ovn-org/libovsdb"
	ishell "gopkg.in/abiosoft/ishell.v2"
)

type OvnShell struct {
	mutex   *sync.RWMutex
	monitor bool
	ovs     *libovsdb.OvsdbClient
	dbModel *libovsdb.DBModel
}

func (s *OvnShell) Monitor(monitor bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.monitor = monitor
}

func (s *OvnShell) OnAdd(table string, m libovsdb.Model) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if s.monitor {
		pp.Printf("An element of table %s has been added\n", table)
		pp.Println(m)
		pp.Println("")
	}
}

func (s *OvnShell) OnUpdate(table string, old, new libovsdb.Model) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if s.monitor {
		pp.Printf("An element of table %s has been updated. New value: \n", table)
		pp.Println(new)
		pp.Println("")
	}
}
func (s *OvnShell) OnDelete(table string, m libovsdb.Model) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if s.monitor {
		pp.Printf("An element of table %s has been deleted\n", table)
		pp.Println(m)
		pp.Println("")
	}
}

func (s *OvnShell) Run(ovs *libovsdb.OvsdbClient, args ...string) {
	s.ovs = ovs
	ovs.Cache.AddEventHandler(s)
	if err := ovs.MonitorAll(""); err != nil {
		panic(err)
	}
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

			for name := range ovnShell.(*OvnShell).ovs.Schema.Tables {
				c.Println(name)
			}
		},
	})

	// List Command
	// Add a subcommand per table
	listCmd := ishell.Cmd{
		Name: "list",
		Help: "List the content of a specific table",
	}

	// To generate the completer for each subcommand we store all the possible fields per table
	tableFields := make(map[string][]string)
	for tname, mtype := range s.dbModel.Types() {
		fields := []string{}
		for i := 0; i < mtype.Elem().NumField(); i++ {
			fields = append(fields, mtype.Elem().Field(i).Name)
		}
		tableFields[tname] = fields
	}

	for name := range s.dbModel.Types() {
		// Trick to be able to use name inside closures
		tableName := name
		subTableCmd := ishell.Cmd{
			Name: name,
			Help: fmt.Sprintf("%s [Field1 Field2 ...]", name),
			LongHelp: fmt.Sprintf(
				"List the content of Table %s", name) +
				fmt.Sprintf("\n\n%s [Field1 Field2 ...]", name) +
				"\n\t[Field1 Field2 ...]: List of fields to show (default: all fields will be shown)" +
				fmt.Sprintf("\n\t\tPossible Fields: %s", strings.Join(tableFields[name], ", ")),
			Func: func(c *ishell.Context) {
				ovnShell := c.Get("ovnShell")
				if ovnShell == nil {
					c.Println("Error: No context")
				}
				// Use a buffer to store the generated output table
				buffer := bytes.Buffer{}
				mtype := ovnShell.(*OvnShell).dbModel.Types()[c.Cmd.Name]
				printer, err := NewStructPrinter(&buffer, mtype.Elem(), c.Args...)
				if err != nil {
					c.Println(err)
					return
				}

				valueList := reflect.New(reflect.SliceOf(mtype.Elem()))
				ovs := ovnShell.(*OvnShell).ovs
				err = ovs.API.List(valueList.Interface())
				if err != nil && err != libovsdb.ErrNotFound {
					c.Println(err)
					return
				}

				// Render the result table
				printer.Append(reflect.Indirect(valueList).Interface())
				printer.Render()
				// Print the result table through shell so it can be paged
				c.ShowPaged(buffer.String())
			},
			Completer: func(args []string) []string {
				return tableFields[tableName]
			},
		}
		listCmd.AddCmd(&subTableCmd)
	}
	shell.AddCmd(&listCmd)

	// If we have arguments, just run them and exit
	if len(args) > 0 {
		shell.Process(args...)
		return
	}
	shell.Run()
}

func newOvnShell(auto bool, dbmodel *libovsdb.DBModel) *OvnShell {
	return &OvnShell{
		mutex:   new(sync.RWMutex),
		monitor: auto,
		dbModel: dbmodel,
	}
}
