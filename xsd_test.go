// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xsd

import (
	"encoding/xml"
	"os"
	"testing"
)

func Test(t *testing.T) {
	var s Schema
	f, err := os.Open("xsd/gccxml.xsd")
	checkError(err)
	defer f.Close()
	checkError(xml.NewDecoder(f).Decode(&s))
	of, err := os.Create("out_go")
	checkError(err)
	p(of, "package gccxml")
	p(of, "")
	s.Gen2(of, "xmlDoc__")
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
