package pkgs

import (
	"go/ast"
	"go/types"
	"iter"
	"nuconv/internal/gen"
	"strings"
)

func typeNameToID(typeName *types.TypeName) gen.TypeID {
	return gen.TypeID{
		Pkg:  gen.PkgPath(typeName.Pkg().Path()),
		Name: typeName.Name(),
	}
}

func iterMarkedTypes(file *ast.File, tag string) iter.Seq[string] {
	return func(yield func(string) bool) {
		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)

			if !ok || genDecl.Tok.String() != "type" {
				continue
			}

			for _, spec := range genDecl.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}

				doc := docText(typeSpec.Doc)
				if doc == "" {
					doc = docText(genDecl.Doc)
				}

				if !strings.Contains(doc, tag) {
					continue
				}

				if !yield(typeSpec.Name.Name) {
					return
				}
			}
		}
	}
}

func docText(group *ast.CommentGroup) string {
	if group == nil {
		return ""
	}
	return group.Text()
}

func pascalToSnakeCase(pascalcase string) string {
	var result strings.Builder
	for i, c := range pascalcase {
		if c >= 'A' && c <= 'Z' {
			if i > 0 {
				result.WriteString("_")
			}
			result.WriteRune(c + ('a' - 'A'))
			continue
		}
		result.WriteRune(c)
	}
	return result.String()
}
