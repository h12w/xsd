package xsd

import (
	"go/ast"
)

type Type interface {
	TypeName() string
	Decls() []ast.Decl
}
type Types []Type

func (a Types) Len() int           { return len(a) }
func (a Types) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Types) Less(i, j int) bool { return a[i].TypeName() < a[j].TypeName() }

type pluralType struct {
	Name string
	Type string
}

func (t pluralType) TypeName() string {
	return t.Name
}

func (t ComplexType) TypeName() string {
	return t.GoName()
}

type enumType struct {
	Name string
	Type string
	KV   []KV
	Doc  string
}
type KV struct {
	Key   string
	Value string
}

func (t *enumType) TypeName() string {
	return t.Name
}
