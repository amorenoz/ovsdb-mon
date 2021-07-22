package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/kylelemons/godebug/diff"
	"github.com/kylelemons/godebug/pretty"
	"github.com/ovn-org/libovsdb/client"
	"github.com/ovn-org/libovsdb/model"
	ishell "gopkg.in/abiosoft/ishell.v2"
)

type eventType string

const (
	updateEvent eventType = "UPDATE"
	addEvent    eventType = "ADD"
	deleteEvent eventType = "DELETE"
)

type OvnEvent struct {
	Timestamp time.Time
	Event     eventType
	Table     string
	Old       model.Model
	New       model.Model
}

type OvnShell struct {
	mutex   *sync.RWMutex
	monitor bool
	ovs     *client.OvsdbClient
	dbModel *model.DBModel
	events  []OvnEvent
}

func (s *OvnShell) Monitor(monitor bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.monitor = monitor
}

func (s *OvnShell) printEvent(event OvnEvent) {
	fmt.Printf("New \033[1m%s\033[0m event on table: \033[1m%s\033[0m\n", event.Event, event.Table)
	switch event.Event {
	case updateEvent:
		fmt.Println(colordiff(event.Old, event.New))
	case addEvent:
		fmt.Printf("\x1b[32m%s\x1b[0m\n", pretty.CompareConfig.Sprint(event.New))
	case deleteEvent:
		fmt.Printf("\x1b[31m%s\x1b[0m\n", pretty.CompareConfig.Sprint(event.Old))
	}
	fmt.Print("\n")
}

func (s *OvnShell) OnAdd(table string, m model.Model) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if s.monitor {
		event := OvnEvent{
			Timestamp: time.Now(),
			Event:     addEvent,
			Table:     table,
			New:       m,
		}
		s.printEvent(event)
		s.events = append(s.events, event)
	}
}

func (s *OvnShell) OnUpdate(table string, old, new model.Model) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if s.monitor {
		event := OvnEvent{
			Timestamp: time.Now(),
			Event:     updateEvent,
			Table:     table,
			New:       new,
			Old:       old,
		}
		s.printEvent(event)
		s.events = append(s.events, event)
	}
}

func (s *OvnShell) OnDelete(table string, m model.Model) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if s.monitor {
		event := OvnEvent{
			Timestamp: time.Now(),
			Event:     deleteEvent,
			Table:     table,
			Old:       m,
		}
		s.printEvent(event)
		s.events = append(s.events, event)
	}
}

func (s *OvnShell) Save(filePath string) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	content, err := json.MarshalIndent(s.events, "", "    ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filePath, content, 0644)
}

func (s *OvnShell) Run(ovs *client.OvsdbClient, args ...string) {
	s.ovs = ovs
	ovs.Cache.AddEventHandler(s)
	if err := ovs.MonitorAll(""); err != nil {
		panic(err)
	}
	shell := ishell.New()
	if shell == nil {
		panic("Failed to create shell")
	}
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
		Name: "save",
		Help: "Save events",
		Func: func(c *ishell.Context) {
			ovnShell := c.Get("ovnShell")
			if len(c.Args) != 1 {
				c.Println("Usage: save <path>")
				return
			}
			filePath := c.Args[0]
			if err := ovnShell.(*OvnShell).Save(filePath); err != nil {
				c.Println(err)
			} else {
				c.Println("File saved")
			}

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
				err = ovs.List(valueList.Interface())
				if err != nil && err != client.ErrNotFound {
					c.Println(err)
					return
				}

				// Render the result table
				err = printer.Append(reflect.Indirect(valueList).Interface())
				if err != nil {
					c.Println(err)
				}
				printer.Render()
				// Print the result table through shell so it can be paged
				if err := c.ShowPaged(buffer.String()); err != nil {
					panic(err)
				}
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
		if err := shell.Process(args...); err != nil {
			panic(err)
		}
	} else {
		shell.Run()
	}
}

func newOvnShell(auto bool, dbmodel *model.DBModel) *OvnShell {
	return &OvnShell{
		mutex:   new(sync.RWMutex),
		monitor: auto,
		dbModel: dbmodel,
	}
}

// colordiff is similar to what pretty.compare does but with colors
func colordiff(a, b interface{}) string {
	config := pretty.CompareConfig
	alines := strings.Split(config.Sprint(a), "\n")
	blines := strings.Split(config.Sprint(b), "\n")

	buf := new(strings.Builder)
	for _, c := range diff.DiffChunks(alines, blines) {
		for _, line := range c.Added {
			fmt.Fprintf(buf, "\x1b[32m+%s\x1b[0m\n", line)
		}
		for _, line := range c.Deleted {
			fmt.Fprintf(buf, "\x1b[31m-%s\x1b[0m\n", line)
		}
		for _, line := range c.Equal {
			fmt.Fprintf(buf, " %s\n", line)
		}
	}
	return strings.TrimRight(buf.String(), "\n")
}
