package main

import (
	"fmt"
	"io"
	"reflect"

	"github.com/olekukonko/tablewriter"
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
		for _, sel := range fieldSel {
			found := false
			for i := 0; i < stype.NumField(); i++ {
				if stype.Field(i).Name == sel {
					found = true
					break
				}
			}
			if !found {
				return nil, fmt.Errorf("Field %s not found in Type %s", sel, stype.Name())
			}
		}
		cols = fieldSel
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
