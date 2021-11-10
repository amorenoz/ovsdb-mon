package main

import (
	"fmt"
	"io"
	"reflect"
	"sort"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/ovn-org/libovsdb/client"
	"github.com/ovn-org/libovsdb/model"
	"k8s.io/apimachinery/pkg/util/sets"
)

// StructPrinter is a wrapper around tablewriter that
// uses the struct's field names as columns
type StructPrinter struct {
	cols  []string
	table *tablewriter.Table
}

// Append adds a list of objects to be printed
func (sp *StructPrinter) Append(objList interface{}) error {
	objListVal := reflect.ValueOf(objList)
	if objListVal.Type().Kind() != reflect.Slice {
		return fmt.Errorf("Append expects a slice of objects. Instead go %s", objListVal.Type())
	}

	data := make([][]string, 0, objListVal.Len())
	for i := 0; i < objListVal.Len(); i++ {
		objVal := objListVal.Index(i)
		objData := make([]string, 0, len(sp.cols))
		for _, col := range sp.cols {
			field := objVal.FieldByName(col)
			if !field.IsValid() {
				continue
			}
			if field.Type().Kind() == reflect.Ptr && !field.IsNil() {
				field = field.Elem()
			}
			objData = append(objData, fmt.Sprintf("%v", field.Interface()))
		}
		data = append(data, objData)
	}
	sp.table.AppendBulk(data)
	return nil
}

// Render prints the table content
func (sp *StructPrinter) Render() {
	sp.table.Render()
}

// NewStructPrinter generates a StructPrinter based on a Writer for a given reflect.Type
// Optionally, a list of field selectors can be given. If so, only those columns
// will be printed.
func NewStructPrinter(writer io.Writer, stype reflect.Type, fieldSel ...string) (*StructPrinter, error) {
	var cols []string
	table := tablewriter.NewWriter(writer)

	if len(fieldSel) > 0 {
		cols = fieldSel
	} else {
		for i := 0; i < stype.NumField(); i++ {
			field := stype.Field(i).Name
			cols = append(cols, field)
		}
	}
	table.SetHeader(cols)
	return &StructPrinter{
		cols:  cols,
		table: table,
	}, nil
}

func getTablesToMonitor(dbModel *model.DBModel, monitorTables string, noMonitorTables string) ([]client.TableMonitor, error) {
	prettyTableNames := func() string {
		tableNames := make([]string, 0, len(dbModel.Types()))
		for tableName := range dbModel.Types() {
			tableNames = append(tableNames, tableName)
		}
		sort.Strings(tableNames)
		return strings.Join(tableNames, ", ")
	}

	// TODO (flaviof): as a follow-up feature, we could allow the caller to specify what fields of the
	// provided table(s) we would like to monitor on. That would give us even more control over what
	// notification we get. I'm thinking something along the lines: table1:field1:field2,table2,table3:name

	// Prepare a map with lower caps of tables for easy lookup
	tables := map[string]string{}
	for table := range dbModel.Types() {
		tables[strings.ToLower(table)] = table
	}

	// Assemble tables to be monitored
	tablesWanted := sets.String{}
	if monitorTables == "" {
		for _, table := range tables {
			tablesWanted.Insert(table)
		}
	} else {
		for _, currMonitorTable := range strings.Split(monitorTables, ",") {
			if currMonitorTable == "" {
				continue
			}
			table, exists := tables[strings.ToLower(currMonitorTable)]
			if !exists {
				return nil, fmt.Errorf("table '%s' is unknown. Available tables are: %s", currMonitorTable, prettyTableNames())
			}
			tablesWanted.Insert(table)
		}
	}

	// Remove anything mentioned in noMonitorTables
	for _, currNoMonitorTable := range strings.Split(noMonitorTables, ",") {
		if currNoMonitorTable == "" {
			continue
		}
		table, exists := tables[strings.ToLower(currNoMonitorTable)]
		if !exists {
			return nil, fmt.Errorf("table '%s' is unknown. Available tables are: %s", currNoMonitorTable, prettyTableNames())
		}
		tablesWanted.Delete(table)
	}

	if len(tablesWanted) == 0 {
		return nil, fmt.Errorf("no tables to monitor, that is kinda sad")
	}

	tablesToMonitor := make([]client.TableMonitor, 0, len(tablesWanted))
	for table := range tablesWanted {
		tableMonitor := client.TableMonitor{Table: table}
		tablesToMonitor = append(tablesToMonitor, tableMonitor)
	}

	return tablesToMonitor, nil
}
