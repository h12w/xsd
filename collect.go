package xsd

import (
	"strconv"

	"bitbucket.org/pkg/inflect"
)

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

func (s *Schema) collect(c *collector) {
	for _, element := range s.Elements {
		element.collect(c, "")
	}
	for _, complexType := range s.ComplexTypes {
		complexType.collect(c, "")
	}
}

func (t ComplexType) collect(c *collector, namespace string) {
	for _, sequence := range t.Sequences {
		sequence.collect(c, namespace)
	}
	for _, choice := range t.Choices {
		choice.collect(c, namespace)
	}
	if t.SimpleContent != nil {
		t.SimpleContent.collect(c, namespace)
	}
	for _, attr := range t.Attributes {
		attr.collect(c, namespace)
	}
	c.add(t)
}

func (s *SimpleContent) collect(c *collector, namespace string) {
	s.Extension.collect(c, namespace)
}

func (e *Extension) collect(c *collector, namespace string) {
	for _, attr := range e.Attributes {
		attr.collect(c, namespace)
	}
}

func (a *Attribute) collect(c *collector, namespace string) {
	if a.SimpleType != nil {
		a.SimpleType.collect(c, namespace+a.GoName())
	}
}

func (s *SimpleType) collect(c *collector, namespace string) {
	if s.Restriction != nil {
		s.Restriction.collect(c, namespace)
	}
}

func (r *Restriction) collect(c *collector, namespace string) {
	if len(r.Enumerations) > 0 {
		t := &enumType{
			Name: namespace,
		}
		switch goType(r.Base) {
		case "NMTOKEN":
			t.Type = "string"
			for _, enum := range r.Enumerations {
				t.KV = append(t.KV, KV{Key: snakeToCamel(enum.Value), Value: strconv.Quote(enum.Value)})
			}
		}
		c.add(t)
	}
}

func (s Sequence) collect(c *collector, namespace string) {
	for _, element := range s.Elements {
		element.collect(c, namespace)
	}
	for _, choice := range s.Choices {
		choice.collect(c, namespace)
	}
}

func (s Choice) collect(c *collector, namespace string) {
	for _, element := range s.Elements {
		element.collect(c, namespace)
	}
}

func (e Element) collect(c *collector, namespace string) {
	if e.ComplexType != nil {
		if e.ComplexType.Name == "" {
			e.ComplexType.Name = e.Name
		}
		if e.ComplexType.Annotation.Documentation == "" {
			e.ComplexType.Annotation = e.Annotation
		}
		e.ComplexType.collect(c, e.GoName())
	}
	/*
		if e.MaxOccurs == "unbounded" {
			typ := e.GoType()
			if e.GoType() == "" {
				typ = e.GoName()
			}
			if inflect.Pluralize(typ) != namespace {
				c.needPlural(typ)
			}
		}
	*/
}
