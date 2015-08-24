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

type Type interface {
	TypeName() string
	Gen(w io.Writer)
}
type Types []Type

func (a Types) Len() int           { return len(a) }
func (a Types) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Types) Less(i, j int) bool { return a[i].TypeName() < a[j].TypeName() }

type collector struct {
	set   map[string]bool
	types Types
}

func newCollector() *collector {
	return &collector{
		make(map[string]bool),
		nil,
	}
}

func (c *collector) add(t Type) {
	if !c.set[t.TypeName()] {
		c.types = append(c.types, t)
		c.set[t.TypeName()] = true
	}
}

func (c *collector) needPlural(name string) {
	c.add(pluralType{
		Name: inflect.Pluralize(name),
		Type: name,
	})
}

func (s *Schema) Gen(w io.Writer) {
	c := newCollector()
	s.collect(c)
	sort.Sort(c.types)
	for _, t := range c.types {
		t.Gen(w)
	}
}

func (s *Schema) collect(c *collector) {
	for _, element := range s.Elements {
		element.collect(c)
	}
	for _, complexType := range s.ComplexTypes {
		complexType.collect(c)
	}
}

func (t ComplexType) collect(c *collector) {
	for _, sequence := range t.Sequences {
		sequence.collect(c)
	}
	for _, choice := range t.Choices {
		choice.collect(c)
	}
	c.add(t)
}

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

func (t pluralType) Gen(w io.Writer) {
	p(w, "type ", t.Name, " []", t.Type)
	p(w)
}

func (s Sequence) collect(c *collector) {
	for _, element := range s.Elements {
		element.collect(c)
	}
	for _, choice := range s.Choices {
		choice.collect(c)
	}
}

func (s Choice) collect(c *collector) {
	for _, element := range s.Elements {
		element.collect(c)
	}
}

func (e Element) collect(c *collector) {
	if e.ComplexType != nil {
		if e.ComplexType.Name == "" {
			e.ComplexType.Name = e.Name
		}
		e.ComplexType.collect(c)
	}
	if e.MaxOccurs == "unbounded" {
		if e.GoType() != "" {
			c.needPlural(e.GoType())
		} else {
			c.needPlural(e.GoName())
		}
	}
}

func (t ComplexType) Gen(w io.Writer) {
	p(w, "type ", t.GoName(), " struct {")
	for _, attr := range t.Attributes {
		attr.Gen(w)
	}
	for _, seq := range t.Sequences {
		seq.Gen(w, false)
	}
	for _, choice := range t.Choices {
		choice.Gen(w, false)
	}
	if t.SimpleContent != nil {
		t.SimpleContent.Gen(w)
	}
	p(w, "}")
	p(w, "")
}

func (s *SimpleContent) Gen(w io.Writer) {
	for _, attr := range s.Extension.Attributes {
		attr.Gen(w)
	}
}

func (t ComplexType) GoName() string {
	return goType(t.Name)
}

func p(w io.Writer, v ...interface{}) {
	fmt.Fprint(w, v...)
	fmt.Fprintln(w)
}

func (a Attribute) Gen(w io.Writer) {
	p(w, a.GoName(), " ", a.GoType(), " `xml:\"", a.Name, ",attr\"`")
}

func (a Attribute) GoName() string {
	return snakeToCamel(a.Name)
}

func (a Attribute) GoType() string {
	if a.Type != "" {
		return goType(a.Type)
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
	if e.MaxOccurs == "unbounded" {
		plural = true
	}
	if e.GoType() == "" {
		e.Type = e.Name
		defer func() { e.Type = "" }()
	}
	if plural {
		pluralName := inflect.Pluralize(e.GoType())
		p(w, pluralName, " ", pluralName, " `xml:\"", e.Name, "\"`")
	} else {
		p(w, e.GoName(), " ", e.GoType(), " `xml:\"", e.Name, "\"`")
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

func goType(s string) string {
	s = trimNamespace(s)
	switch s {
	case "integer":
		return "int"
	case "boolean":
		return "bool"
	case "string":
		return s
	}
	s = strings.TrimSuffix(s, "Type")
	s = strings.TrimSuffix(s, "type")
	return snakeToCamel(s)
}
