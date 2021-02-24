package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of modelgen:\n")
	fmt.Fprintf(os.Stderr, "\tovsidl [flags] OVS_SCHEMA\n")
	fmt.Fprintf(os.Stderr, "For more information, see:\n")
	fmt.Fprintf(os.Stderr, "\t   TODO INSERT LINK \n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

const (
	version = "0.1.1"
)

var (
	outDirP  = flag.String("o", ".", "Directory where the generated files shall be stored")
	pkgNameP = flag.String("p", "modelgen", "Package name")
	dryRun   = flag.Bool("d", false, "Dry run")
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("modelgen: ")
	flag.Usage = Usage
	flag.Parse()
	outDir := *outDirP
	pkgName := *pkgNameP

	/*Option handling*/
	outDir, err := filepath.Abs(outDir)
	if err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(outDir, 0700); err != nil {
		log.Fatal(err)
	}

	if len(flag.Args()) != 1 {
		flag.Usage()
		os.Exit(2)
	}

	schemaFile, err := os.Open(flag.Args()[0])
	if err != nil {
		log.Fatal(err)
	}
	defer schemaFile.Close()

	schemaBytes, err := ioutil.ReadAll(schemaFile)
	if err != nil {
		log.Fatal(err)
	}

	generator := DbModelGenerator{
		Generator: Generator{pkgName: pkgName},
	}

	if err := generator.Parse(schemaBytes); err != nil {
		log.Fatal(err)
	}
	if err := generator.Generate(); err != nil {
		log.Fatal(err)
	}

	outFile := filepath.Join(outDir, generator.FileName())
	src := generator.Format()
	if err := write_file(outFile, src); err != nil {
		log.Panic(err)
	}

	for _, tableGen := range generator.Tables() {
		if err := tableGen.Generate(); err != nil {
			log.Panic(err)
		}
		src = tableGen.Format()
		outFile = filepath.Join(outDir, tableGen.FileName())
		if err := write_file(outFile, src); err != nil {
			log.Panic(err)
		}
	}
}

func write_file(filename string, src []byte) error {
	if *dryRun {
		fmt.Printf("----Content of file %s\n", filename)
		fmt.Printf(string(src))
		fmt.Printf("\n")
		return nil
	} else {
		return ioutil.WriteFile(filename, src, 0644)
	}
}
