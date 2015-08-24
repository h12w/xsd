package xsd

import "io"

type Type interface {
	TypeName() string
	Gen(w io.Writer)
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
