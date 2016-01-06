// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xsd

import (
	"fmt"
	"io"
	"reflect"
	"regexp"
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
	return goName(a.Name)
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
	return goName(e.Name)
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

func goName(s string) string {
	s = snakeToCamel(s)
	if strings.HasSuffix(s, "Id") {
		s = strings.TrimSuffix(s, "Id") + "ID"
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

var upperLower = regexp.MustCompile(`[\p{Lu}][\p{Ll}]`)
var lowerUpper = regexp.MustCompile(`([\p{Ll}])([\p{Lu}])`)

// convert name from camel case to snake case
func camelToSnake(s string) string {
	s = strings.TrimPrefix(upperLower.ReplaceAllString(s, `_${0}`), "_")
	s = lowerUpper.ReplaceAllString(s, `${1}_${2}`)
	s = strings.ToLower(s)
	return s
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

type XMLNodeType int

const (
	XMLElement XMLNodeType = iota
	XMLAttr
	XMLCharData
	XMLInnerXML
	XMLComment
	XMLOmitted
)

type XMLTag struct {
	Type      XMLNodeType
	Name      string
	Omitempty bool
}

func ParseXMLTag(s string) XMLTag {
	values := strings.Split(reflect.StructTag(strings.Trim(s, "`")).Get("xml"), ",")
	t := XMLTag{Name: values[0]}
	for _, value := range values[1:] {
		switch value {
		case "-":
			t.Type = XMLOmitted
		case "attr":
			t.Type = XMLAttr
		case "chardata":
			t.Type = XMLCharData
		case "innerxml":
			t.Type = XMLInnerXML
		case "comment":
			t.Type = XMLComment
		case "omitempty":
			t.Omitempty = true
		}
	}
	return t
}

func (t XMLTag) String() string {
	var values []string
	values = append(values, t.Name)
	switch t.Type {
	case XMLOmitted:
		return `xml:"-"`
	case XMLAttr:
		values = append(values, "attr")
	case XMLCharData:
		values = append(values, "chardata")
	case XMLInnerXML:
		values = append(values, "innerxml")
	case XMLComment:
		values = append(values, "comment")
	}
	if t.Omitempty {
		values = append(values, "omitempty")
	}
	return fmt.Sprintf(`xml:"%s"`, strings.Join(values, ","))
}

type BSONNodeType int

const (
	BSONNormal BSONNodeType = iota
	BSONInline
	BSONOmitted
)

type BSONTag struct {
	Type      BSONNodeType
	Name      string
	Omitempty bool
}

func (t BSONTag) String() string {
	if t.Type == BSONOmitted {
		return `bson:"-"`
	}
	values := []string{t.Name}
	if t.Omitempty {
		values = append(values, "omitempty")
	}
	return fmt.Sprintf(`bson:"%s"`, strings.Join(values, ","))
}

type JSONNodeType int

const (
	JSONNormal JSONNodeType = iota
	JSONOmitted
)

type JSONTag struct {
	Type      JSONNodeType
	Name      string
	Omitempty bool
}

func (t JSONTag) String() string {
	if t.Type == JSONOmitted {
		return `json:"-"`
	}
	values := []string{t.Name}
	if t.Omitempty {
		values = append(values, "omitempty")
	}
	return fmt.Sprintf(`json:"%s"`, strings.Join(values, ","))
}
