// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Doc: http://www.w3schools.com/schema/schema_elements_ref.asp

It is convecient to be used together with a XML Schema learner like: github.com/kore/XML-Schema-learner
*/
package xsd

import (
	"encoding/xml"
)

type Schema struct {
	XMLName         xml.Name         `xml:"schema"`
	SimpleTypes     []SimpleType     `xml:"simpleType"`
	ComplexTypes    ComplexTypes     `xml:"complexType"`
	Groups          []Group          `xml:"group"`
	Attributes      Attributes       `xml:"attribute"`
	AttributeGroups []AttributeGroup `xml:"attributeGroup"`
	Elements        []Element        `xml:"element"`
	Notations       []Notation       `xml:"notation"`
	Annotations     []Annotation     `xml:"annotation"`
}

type SimpleType struct {
	Id          string       `xml:"id,attr"`
	Name        string       `xml:"name,attr"`
	Restriction *Restriction `xml:"restriction"`
	List        *List        `xml:"list"`
	Union       *Union       `xml:"union"`
}

type Restriction struct {
	Base         string       `xml:"base,attr"`
	Enumerations Enumerations `xml:"enumeration"`
}

type Enumerations []Enumeration
type Enumeration struct {
	Value string `xml:"value,attr"`
}

type List struct {
}

type Union struct {
}

type ComplexType struct {
	Id              string           `xml:"id,attr"`
	Name            string           `xml:"name,attr"`
	Groups          []Group          `xml:"group"`
	All             []All            `xml:"all"`
	Choices         []Choice         `xml:"choice"`
	Sequences       []Sequence       `xml:"sequence"`
	Attributes      Attributes       `xml:"attribute"`
	AttributeGroups []AttributeGroup `xml:"attributeGroup"`
	SimpleContent   *SimpleContent   `xml:"simpleContent"`
	Annotation      Annotation       `xml:"annotation"`
}
type ComplexTypes []ComplexType

type SimpleContent struct {
	Extension Extension `xml:"extension"`
}

type Extension struct {
	Base       string     `xml:"base,attr"`
	Attributes Attributes `xml:"attribute"`
}

type Attribute struct {
	Id         string      `xml:"id,attr"`
	Name       string      `xml:"name,attr"`
	Type       string      `xml:"type,attr"`
	Use        string      `xml:"use,attr"`
	SimpleType *SimpleType `xml:"simpleType"`
	Annotation Annotation  `xml:"annotation"`
}
type Attributes []Attribute

type Element struct {
	Id          string       `xml:"id,attr"`
	Name        string       `xml:"name,attr"`
	Type        string       `xml:"type,attr"`
	ComplexType *ComplexType `xml:"complexType"`
	MinOccurs   string       `xml:"minOccurs,attr"`
	MaxOccurs   string       `xml:"maxOccurs,attr"`
	Annotation  Annotation   `xml:"annotation"`
}

type Choice struct {
	Id        string     `xml:"id,attr"`
	Elements  []Element  `xml:"element"`
	Groups    []Group    `xml:"group"`
	Choices   []Choice   `xml:"choice"`
	Sequences []Sequence `xml:"sequence"`
}

type Sequence struct {
	Id        string     `xml:"id,attr"`
	Elements  []Element  `xml:"element"`
	Groups    []Group    `xml:"group"`
	Choices   []Choice   `xml:"choice"`
	Sequences []Sequence `xml:"sequence"`
	MinOccurs string     `xml:"minOccurs,attr"`
	MaxOccurs string     `xml:"maxOccurs,attr"`
}

type All struct {
	Elements []Element `xml:"element"`
}

type Group struct {
	Choices   []Choice   `xml:"choice"`
	Sequences []Sequence `xml:"sequence"`
}

type AttributeGroup struct {
	Attributes      Attributes       `xml:"attribute"`
	AttributeGroups []AttributeGroup `xml:"attributeGroup"`
}

type Notation struct {
}

type Annotation struct {
	Documentation string `xml:"documentation"`
}
