// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xsd

import (
	"fmt"
	"io"
	"strings"
)

func cleanDoc(s string) string {
	ss := strings.Split(s, "\n")
	for i := range ss {
		ss[i] = strings.TrimSpace(ss[i])
	}
	return strings.Join(ss, " ")
}

func (t ComplexType) GoName() string {
	return goType(t.Name)
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
	case "int", "string", "bool", "AnyURI":
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
