package xsd

import (
	"go/ast"
	"go/token"
	"strings"
)

func addBSONTags(decls []ast.Decl) []ast.Decl {
	for _, decl := range decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		if genDecl.Tok != token.TYPE {
			continue
		}
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}
			for _, field := range structType.Fields.List {
				xmlTag := ParseXMLTag(field.Tag.Value)
				name := camelToSnake(field.Names[0].Name)
				jsonTag := JSONTag{
					Name:      name,
					Omitempty: xmlTag.Omitempty,
				}
				if xmlTag.Type == XMLOmitted {
					jsonTag.Type = JSONOmitted
				}
				field.Tag = tag(xmlTag.String() + " " + jsonTag.String())
			}
		}
	}
	return decls
}

func elevateSubArrays(decls []ast.Decl) []ast.Decl {
	arrayTypes := make(map[string]*ast.Field)
	for _, decl := range decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		if genDecl.Tok != token.TYPE {
			continue
		}
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}
			if len(structType.Fields.List) != 1 {
				continue
			}
			field := structType.Fields.List[0]
			fieldType, ok := field.Type.(*ast.Ident)
			if !ok {
				continue
			}
			if strings.HasPrefix(fieldType.Name, "[]") {
				arrayTypes[typeSpec.Name.Name] = field
			}
		}
	}
	for _, decl := range decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		if genDecl.Tok != token.TYPE {
			continue
		}
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}
			for _, field := range structType.Fields.List {
				if typ, ok := field.Type.(*ast.Ident); ok {
					typeName := typ.Name
					typeName = strings.TrimPrefix(typeName, "*")
					if f, ok := arrayTypes[typeName]; ok {
						parentTag := ParseXMLTag(field.Tag.Value)
						childTag := ParseXMLTag(f.Tag.Value)
						*field = *f
						field.Tag = tag(XMLTag{
							Name:      parentTag.Name + ">" + childTag.Name,
							Type:      childTag.Type,
							Omitempty: parentTag.Omitempty || childTag.Omitempty,
						}.String())
					}
				}
			}
		}
	}

	var newDecls []ast.Decl
nextDecl:
	for _, decl := range decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok {
			if genDecl.Tok == token.TYPE {
				for _, spec := range genDecl.Specs {
					if typeSpec, ok := spec.(*ast.TypeSpec); ok {
						if _, ok := arrayTypes[typeSpec.Name.Name]; ok {
							continue nextDecl
						}
					}
				}
			}
		}
		newDecls = append(newDecls, decl)
	}
	return newDecls
}
