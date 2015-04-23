// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xsd

import (
	"bitbucket.org/pkg/inflect"
	"fmt"
	"io"
	"strings"
)

func (s Schema) Gen(w io.Writer) {
	for _, complexType := range s.ComplexTypes {
		complexType.Gen(w)
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
	p(w, "}")
	pluralName := inflect.Pluralize(t.GoName())
	p(w, "type ", pluralName, " []", t.GoName())
	p(w, "")
}

func (t ComplexType) GoName() string {
	return snakeToCamel(t.Name)
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
	return a.Type
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
	if e.Type == "" {
		e.Type = e.Name
		defer func() { e.Type = "" }()
	}
	if plural {
		pluralName := inflect.Pluralize(e.GoName())
		p(w, pluralName, " ", pluralName, " `xml:\"", e.Name, "\"`")
	} else {
		p(w, e.GoName(), " ", e.GoType(), " `xml:\"", e.Name, "\"`")
	}
}

func (e Element) GoName() string {
	return snakeToCamel(e.Name)
}

func (e Element) GoType() string {
	return snakeToCamel(e.Type)
}

func snakeToCamel(s string) string {
	ss := strings.Split(s, "_")
	for i := range ss {
		ss[i] = strings.Title(ss[i])
	}
	return strings.Join(ss, "")
}
