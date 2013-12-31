// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xsd

import (
	"bitbucket.org/pkg/inflect"
	"io"
	"sort"
)

func (s ComplexTypes) Len() int {
	return len(s)
}

func (s ComplexTypes) Less(i, j int) bool {
	return s[i].GoName() < s[j].GoName()
}

func (s ComplexTypes) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Attributes) Len() int {
	return len(s)
}

func (s Attributes) Less(i, j int) bool {
	return s[i].GoName() < s[j].GoName()
}

func (s Attributes) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// take attributes out and embed common attributes.
func (s Schema) Gen2(w io.Writer, parent string) {
	sort.Sort(s.ComplexTypes)
	am := make(map[string]Attribute)
	for _, ct := range s.ComplexTypes {
		for _, a := range ct.Attributes {
			am[a.Name] = a
		}
	}

	for _, complexType := range s.ComplexTypes {
		complexType.Gen2(w, parent)
	}

	as := make(Attributes, 0, len(am))
	for _, a := range am {
		as = append(as, a)
	}
	sort.Sort(as)
	for _, a := range as {
		a.GenStruct(w)
	}
}

func (t ComplexType) Gen2(w io.Writer, parent string) {
	p(w, "type ", t.GoName(), " struct {")
	for _, attr := range t.Attributes {
		attr.Gen2(w)
	}
	for _, seq := range t.Sequences {
		seq.Gen(w, false)
	}
	p(w, parent)
	p(w, "}")
	pluralName := inflect.Pluralize(t.GoName())
	p(w, "type ", pluralName, " []*", t.GoName())
	p(w, "")
}

func (a Attribute) Gen2(w io.Writer) {
	p(w, a.GoName2())
}

func (a Attribute) GoName2() string {
	return a.GoName() + "__"
}

func (a Attribute) GenStruct(w io.Writer) {
	p(w, "type ", a.GoName2(), " struct {")
	p(w, a.GoName(), "_ ", a.GoType(), " `xml:\"", a.Name, ",attr\"`")
	p(w, "}")
	p(w, "func (s ", a.GoName2(), ")", a.GoName(), "()", a.GoType(), "{")
	p(w, "return s.", a.GoName(), "_")
	p(w, "}")
}
