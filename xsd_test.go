// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xsd

import (
	"encoding/xml"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"testing"
)

func Test(t *testing.T) {
	var s Schema
	f, err := os.Open("data/gccxml.xsd")
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

func TestCheck(t *testing.T) {
	// src is the input for which we want to print the AST.
	src := `
package main

// T is a type T
type T struct {
	S string
}

// E is a type E
type E int
const (
	E1 E = 1
	E2 E = 2
)
`

	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		panic(err)
	}

	// Print the AST.
	ast.Print(fset, f)
}

func TestGen(t *testing.T) {
	f := &ast.File{
		Name: &ast.Ident{Name: "vast"},
		Decls: []ast.Decl{
			&ast.GenDecl{
				Tok:    token.CONST,
				Lparen: 1,
				Specs: []ast.Spec{
					&ast.ValueSpec{
						Names:  []*ast.Ident{{Name: "A"}},
						Type:   &ast.Ident{Name: "int"},
						Values: []ast.Expr{&ast.BasicLit{Kind: token.INT, Value: "1"}},
					},
					&ast.ValueSpec{
						Names:  []*ast.Ident{{Name: "B"}},
						Type:   &ast.Ident{Name: "int"},
						Values: []ast.Expr{&ast.BasicLit{Kind: token.INT, Value: "2"}},
					},
				},
			},
			&ast.GenDecl{
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: &ast.Ident{Name: "VAST"},
						Type: &ast.StructType{
							Fields: &ast.FieldList{
								List: []*ast.Field{
									{Names: []*ast.Ident{&ast.Ident{Name: "S"}}, Type: &ast.Ident{Name: "string"}},
								},
							},
						},
					},
				},
			},
		},
	}
	fset := token.NewFileSet()
	ast.Print(fset, f)
	err := printer.Fprint(os.Stdout, fset, f)
	if err != nil {
		t.Fatal(err)
	}
}
