package main

import (
	"encoding/xml"
	"fmt"
	"os"

	"h12.me/xsd"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("xsd2struct [package name] [file]...")
	}
	pkg := os.Args[1]
	files := os.Args[2:]
	fmt.Println("package", pkg)
	for _, file := range files {
		gen(file)
	}
}

func gen(file string) {
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	var s xsd.Schema
	if err := xml.NewDecoder(f).Decode(&s); err != nil {
		panic(err)
	}
	s.Gen(os.Stdout)
	fmt.Println()
}
