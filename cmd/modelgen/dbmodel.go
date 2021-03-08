package main

import (
	"log"
	//"os"
	"encoding/json"

	"github.com/ebay/libovsdb"
)

type DbModelGenerator struct {
	Generator
	schema libovsdb.DatabaseSchema
	tables []TableGenerator
}

func (g *DbModelGenerator) FileName() string {
	return "model.go"
}

// TODO: Check validity
func (g *DbModelGenerator) Parse(bytes []byte) error {
	if err := json.Unmarshal(bytes, &g.schema); err != nil {
		return err
	}
	for name, table := range g.schema.Tables {
		g.tables = append(g.tables, TableGenerator{
			Generator: Generator{
				pkgName: g.pkgName,
			},
			name:  name,
			table: table,
		})
	}
	return nil
}

func (g *DbModelGenerator) Tables() []TableGenerator {
	return g.tables
}

func (g *DbModelGenerator) Dump() {
	log.Printf("%p", &g)
	log.Printf("%s", g.schema.Name)

	log.Printf("DB NAME %s \n", g.schema.Name)
	log.Printf("DB Version%s \n", g.schema.Version)
	for tableName, table := range g.schema.Tables {
		log.Printf("\tTable: %s\n", tableName)
		log.Printf("\t\tIndexes: %v\n", table.Indexes)

		log.Printf("\t\tColumns:\n")
		for columnName, column := range table.Columns {
			log.Printf("\t\t\tName: %s:\n", columnName)
			log.Printf("\t\t\tColumn:: %s:\n", column)
		}
	}
}

func (g *DbModelGenerator) Generate() error {
	g.PrintHeader()
	// Add imports
	imports := [][]string{
		{"goovn", "github.com/ebay/go-ovn"},
	}
	g.PrintImports(imports)
	g.Printf("\n")
	g.Printf("// DB Model \n")
	g.PrintDBModel()
	return nil
}

/*PrintDBModel generates a Model constructor such as this:

func DBModel() goovn.DBModel {
    return goovn.NewDBModel ([]Model {
	&LogicalSwitch{},
	&LogicalRouter{},
...
    })
}
*/
func (g *DbModelGenerator) PrintDBModel() {
	g.Printf("func DBModel() (*goovn.DBModel, error){\n")
	g.Printf("\treturn goovn.NewDBModel ([]goovn.Model {\n")
	for _, table := range g.tables {
		g.Printf("\t\t&%s{},\n", table.TableTypeName())
	}
	g.Printf("\t})\n")
	g.Printf("}\n")
}

// TODO: Create a API constructor
//func (g *DbModelGenerator) PrintApi() {
//	apiName := camelCase(g.schema.Name)
//	g.Printf("type %s struct {\n", apiName)
//	for _, table := range g.tables {
//		g.Printf("\t%s %s\n", table.TableTypeName(), table.TableApiStruct())
//	}
//	g.Printf("}\n")
//	g.Printf("\n")
//	g.Printf("func New%sApi(c goovn.Client) %s {\n", apiName, apiName)
//	g.Printf("\t return %s {\n", apiName)
//	for _, table := range g.tables {
//		g.Printf("\t\t%s:%s,\n", table.TableTypeName(), table.ApiInitString(3))
//	}
//	g.Printf("\t}\n")
//	g.Printf("}\n")
//}

func (g *DbModelGenerator) apiName() string {
	return camelCase(g.schema.Name)
}
