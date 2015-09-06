package xsd

import (
	"go/ast"
	"go/token"
	"sort"
	"strings"

	"bitbucket.org/pkg/inflect"
)

func (s *Schema) Ast(name string) *ast.File {
	f := &ast.File{
		Name: &ast.Ident{Name: name},
	}
	c := newCollector()
	s.collect(c)
	sort.Sort(c.types)
	var decls []ast.Decl
	for _, typ := range c.types {
		decls = append(decls, typ.Decls()...)
	}
	f.Decls = addBSONTags(elevateSubArrays(decls))
	return f
}

func (t *enumType) Decls() []ast.Decl {
	typeDecl := &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: &ast.Ident{Name: t.Name},
				Type: &ast.Ident{Name: t.Type},
			},
		},
		Doc: comment(t.Doc),
	}

	constDecl := &ast.GenDecl{
		Tok:    token.CONST,
		Lparen: 1,
	}
	kind := token.Lookup(t.Type)
	for _, kv := range t.KV {
		constDecl.Specs = append(constDecl.Specs,
			&ast.ValueSpec{
				Names:  []*ast.Ident{{Name: kv.Key}},
				Type:   &ast.Ident{Name: t.Name},
				Values: []ast.Expr{&ast.BasicLit{Kind: kind, Value: kv.Value}},
			})
	}
	return []ast.Decl{typeDecl, constDecl}
}

func (t pluralType) Decls() []ast.Decl {
	return []ast.Decl{&ast.GenDecl{Tok: token.TYPE, Specs: []ast.Spec{&ast.TypeSpec{
		Name: &ast.Ident{Name: t.Name},
		Type: &ast.Ident{Name: t.Type},
	}}}}
}

func (t ComplexType) Decls() []ast.Decl {
	doc := t.Annotation.Documentation
	if doc != "" {
		doc = cleanDoc(doc)
		if !strings.HasPrefix(doc, t.GoName()) {
			doc = t.GoName() + " is " + doc
		}
	}
	if doc == "" {
		doc = " "
	}
	var fields []*ast.Field
	if t.SimpleContent != nil {
		fields = append(fields, t.SimpleContent.Fields(t.GoName())...)
	}
	for _, attr := range t.Attributes {
		fields = append(fields, attr.Field(t.GoName()))
	}
	for _, seq := range t.Sequences {
		fields = append(fields, seq.Fields(false)...)
	}
	for _, choice := range t.Choices {
		fields = append(fields, choice.Fields(false)...)
	}
	return []ast.Decl{&ast.GenDecl{Tok: token.TYPE, Specs: []ast.Spec{&ast.TypeSpec{
		Name: &ast.Ident{Name: t.GoName()},
		Type: &ast.StructType{Fields: &ast.FieldList{List: fields}},
	}},
		Doc: comment(doc),
	}}
}

func (s *SimpleContent) Fields(namespace string) []*ast.Field {
	typ := goType(s.Extension.Base)
	fields := []*ast.Field{
		&ast.Field{
			Names: []*ast.Ident{{Name: "Value"}},
			Type:  &ast.Ident{Name: typ},
			Tag:   tag(XMLTag{Type: XMLCharData}.String()),
		},
	}
	for _, attr := range s.Extension.Attributes {
		fields = append(fields, attr.Field(namespace))
	}
	return fields
}

func (a Attribute) Field(namespace string) *ast.Field {
	typ := a.GoType(namespace)
	omitempty := false
	if a.Use == "optional" {
		omitempty = true
		typ = omitType(typ)
	}
	doc := ""
	if a.Annotation.Documentation != "" {
		doc = cleanDoc(a.Annotation.Documentation)
	}
	return &ast.Field{
		Names: []*ast.Ident{{Name: a.GoName()}},
		Type:  &ast.Ident{Name: typ},
		Tag:   tag(XMLTag{Name: a.Name, Type: XMLAttr, Omitempty: omitempty}.String()),
		Doc:   comment(doc),
	}
}

func (s Sequence) Fields(plural bool) []*ast.Field {
	var fields []*ast.Field
	if s.MaxOccurs == "unbounded" {
		plural = true
	}
	for _, seq := range s.Sequences {
		fields = append(fields, seq.Fields(plural)...)
	}
	for _, choice := range s.Choices {
		fields = append(fields, choice.Fields(plural)...)
	}
	for _, elem := range s.Elements {
		fields = append(fields, elem.Field(plural))
	}
	return fields
}

func (c Choice) Fields(plural bool) []*ast.Field {
	var fields []*ast.Field
	for _, elem := range c.Elements {
		fields = append(fields, elem.Field(plural))
	}
	return fields
}

func (e Element) Field(plural bool) *ast.Field {
	omitempty := false
	if e.MinOccurs == "0" {
		omitempty = true
	}
	if e.MaxOccurs == "unbounded" {
		plural = true
	}
	if e.GoType() == "" {
		e.Type = e.Name
		defer func() { e.Type = "" }()
	}
	doc := ""
	if e.Annotation.Documentation != "" {
		doc = e.Annotation.Documentation
	}
	if plural {
		pluralName := inflect.Pluralize(e.GoName())
		pluralType := "[]" + e.GoType()
		return &ast.Field{
			Names: []*ast.Ident{{Name: pluralName}},
			Type:  &ast.Ident{Name: pluralType},
			Tag:   tag(XMLTag{Name: e.Name, Omitempty: omitempty}.String()),
			Doc:   comment(doc),
		}
	}
	typ := e.GoType()
	if e.MinOccurs == "0" {
		typ = omitType(typ)
	}
	return &ast.Field{
		Names: []*ast.Ident{{Name: e.GoName()}},
		Type:  &ast.Ident{Name: typ},
		Tag:   tag(XMLTag{Name: e.Name, Omitempty: omitempty}.String()),
		Doc:   comment(doc),
	}
}

func comment(doc string) *ast.CommentGroup {
	if doc == "" {
		return nil
	}
	return &ast.CommentGroup{List: []*ast.Comment{{Text: "\n// " + doc, Slash: 1}}}
}

func tag(value string) *ast.BasicLit {
	if value == "" {
		return nil
	}
	return &ast.BasicLit{Kind: token.STRING, Value: "`" + value + "`\n"}
}
