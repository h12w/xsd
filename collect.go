package xsd

import "bitbucket.org/pkg/inflect"

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
