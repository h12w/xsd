// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xsd

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"bitbucket.org/pkg/inflect"
)

func (s *Schema) Gen(w io.Writer) {
	c := newCollector()
	s.collect(c)
	sort.Sort(c.types)
	for _, typ := range c.types {
		typ.Gen(w)
	}
}

func (t *enumType) Gen(w io.Writer) {
	fmt.Fprintf(w, "type %s %s\n", t.Name, t.Type)
	fmt.Fprintln(w, "const (")
	for _, kv := range t.KV {
		fmt.Fprintf(w, "%s %s = %s\n", kv.Key, t.Name, kv.Value)
	}
	fmt.Fprintln(w, ")")
}

func (t pluralType) Gen(w io.Writer) {
	p(w, "type ", t.Name, " []", t.Type)
	p(w)
}

func cleanDoc(s string) string {
	ss := strings.Split(s, "\n")
	for i := range ss {
		ss[i] = strings.TrimSpace(ss[i])
	}
	return strings.Join(ss, " ")
}

func (t ComplexType) Gen(w io.Writer) {
	if doc := t.Annotation.Documentation; doc != "" {
		doc = cleanDoc(doc)
		if !strings.HasPrefix(doc, t.GoName()) {
			doc = t.GoName() + " is " + doc
		}
		p(w, "// "+doc)
	}
	p(w, "type ", t.GoName(), " struct {")
	if t.SimpleContent != nil {
		t.SimpleContent.Gen(w, t.GoName())
	}
	for _, attr := range t.Attributes {
		attr.Gen(w, t.GoName())
	}
	for _, seq := range t.Sequences {
		seq.Gen(w, false)
	}
	for _, choice := range t.Choices {
		choice.Gen(w, false)
	}
	p(w, "}")
	p(w, "")
}

func (s *SimpleContent) Gen(w io.Writer, namespace string) {
	typ := goType(s.Extension.Base)
	p(w, "Value ", typ, " `xml:\",chardata\"`")
	for _, attr := range s.Extension.Attributes {
		attr.Gen(w, namespace)
	}
}

func (t ComplexType) GoName() string {
	return goType(t.Name)
}

func (a Attribute) Gen(w io.Writer, namespace string) {
	omitempty := ""
	typ := a.GoType(namespace)
	if a.Use == "optional" {
		omitempty = ",omitempty"
		typ = omitType(typ)
	}

	if a.Annotation.Documentation != "" {
		p(w, "")
		doc := cleanDoc(a.Annotation.Documentation)
		p(w, "// "+doc)
	}
	p(w, a.GoName(), " ", typ, " `xml:\"", a.Name, ",attr"+omitempty+"\"`")
}

func (a Attribute) GoName() string {
	return snakeToCamel(a.Name)
}

func (a Attribute) GoType(namespace string) string {
	if a.Type != "" {
		return goType(a.Type)
	}
	if goType(a.SimpleType.Restriction.Base) == "NMTOKEN" {
		return namespace + a.GoName()
	}
	return goType(a.SimpleType.Restriction.Base)
}

func (s Sequence) Gen(w io.Writer, plural bool) {
	if s.MaxOccurs == "unbounded" {
		plural = true
	}
	for _, seq := range s.Sequences {
		seq.Gen(w, plural)
	}
	for _, choice := range s.Choices {
		choice.Gen(w, plural)
	}
	for _, elem := range s.Elements {
		elem.Gen(w, plural)
	}
}

func (c Choice) Gen(w io.Writer, plural bool) {
	for _, e := range c.Elements {
		e.Gen(w, plural)
	}
}

func (e Element) Gen(w io.Writer, plural bool) {
	omitempty := ""
	if e.MinOccurs == "0" {
		omitempty = ",omitempty"
	}
	if e.MaxOccurs == "unbounded" {
		plural = true
	}
	if e.GoType() == "" {
		e.Type = e.Name
		defer func() { e.Type = "" }()
	}
	if e.Annotation.Documentation != "" {
		p(w, "")
		p(w, "// "+e.Annotation.Documentation)
	}
	if plural {
		pluralName := inflect.Pluralize(e.GoName())
		pluralType := "[]" + e.GoType()
		p(w, pluralName, " ", pluralType, " `xml:\"", e.Name, omitempty+"\"`")
	} else {
		typ := e.GoType()
		if e.MinOccurs == "0" {
			typ = omitType(typ)
		}
		p(w, e.GoName(), " ", typ, " `xml:\"", e.Name, omitempty+"\"`")
	}
}

func (e Element) GoName() string {
	return snakeToCamel(e.Name)
}

func (e Element) GoType() string {
	return goType(e.Type)
}

func trimNamespace(s string) string {
	m := strings.Split(s, ":")
	if len(m) == 2 {
		return m[1]
	}
	return s
}

func snakeToCamel(s string) string {
	ss := strings.Split(s, "_")
	for i := range ss {
		ss[i] = strings.Title(ss[i])
	}
	return strings.Join(ss, "")
}

func omitType(s string) string {
	switch s {
	case "int", "string", "bool":
		return s
	}
	return "*" + s
}

func goType(s string) string {
	s = trimNamespace(s)
	switch s {
	case "integer":
		return "int"
	case "boolean":
		return "bool"
	case "string":
		return "string"
	case "decimal":
		return "float32"
	}
	s = strings.TrimSuffix(s, "Type")
	s = strings.TrimSuffix(s, "type")
	return snakeToCamel(s)
}

func p(w io.Writer, v ...interface{}) {
	fmt.Fprint(w, v...)
	fmt.Fprintln(w)
}
