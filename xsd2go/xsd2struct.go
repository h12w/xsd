package main

import (
	"encoding/xml"
	"fmt"
	"os"

	"h12.me/xsd"
)

func main() {
	files := os.Args[1:]
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
