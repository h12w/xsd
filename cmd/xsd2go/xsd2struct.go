package main

import (
	"encoding/xml"
	"fmt"
	"go/printer"
	"go/token"
	"log"
	"os"

	"h12.io/xsd"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("xsd2struct [package name] [file]...")
	}
	pkg := os.Args[1]
	files := os.Args[2:]
	for _, file := range files {
		if err := gen(pkg, file); err != nil {
			log.Fatal(err)
		}
	}
}

func gen(pkg, file string) error {
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	var s xsd.Schema
	if err := xml.NewDecoder(f).Decode(&s); err != nil {
		panic(err)
	}
	astFile := s.Ast(pkg)
	return printer.Fprint(os.Stdout, token.NewFileSet(), astFile)
}
