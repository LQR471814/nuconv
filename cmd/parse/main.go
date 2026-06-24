package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"log"
	"strings"

	"golang.org/x/tools/go/packages"
)

func main() {
	cfg := &packages.Config{
		Mode: packages.NeedName |
			packages.NeedFiles |
			packages.NeedSyntax |
			packages.NeedTypes |
			packages.NeedTypesInfo,
		Fset: token.NewFileSet(),
		Dir:  ".",
	}

	pkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		panic(err)
	}

	for _, pkg := range pkgs {
		for _, err := range pkg.Errors {
			fmt.Println("package error:", err)
		}

		log.Println(pkg.PkgPath)

		annotated := make(map[string]struct{})
		for _, file := range pkg.Syntax {
			parseFile(file, "@nutype", annotated)
		}
		parsePkg(pkg, annotated)
	}
}

func parseFile(file *ast.File, tag string, annotated map[string]struct{}) {
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

			annotated[typeSpec.Name.Name] = struct{}{}
		}
	}
}

func parsePkg(pkg *packages.Package, annotated map[string]struct{}) {
	if len(pkg.Errors) > 0 {
		for _, e := range pkg.Errors {
			log.Println(e)
		}
		return
	}

	scope := pkg.Types.Scope()

	for _, name := range scope.Names() {
		_, isAnnot := annotated[name]
		if !isAnnot {
			continue
		}

		obj := scope.Lookup(name)

		typeName, ok := obj.(*types.TypeName)
		if !ok {
			continue
		}

		named, ok := typeName.Type().(*types.Named)
		if !ok {
			continue
		}

		st, ok := named.Underlying().(*types.Struct)
		if !ok {
			continue
		}

		fmt.Println("struct:", name, "name", named.Obj().Name(), "path", named.Obj().Pkg().Path())

		for i := 0; i < st.NumFields(); i++ {
			field := st.Field(i)

			fmt.Printf(
				"  field: %s type=%s exported=%v tag=%q\n",
				field.Name(),
				field.Type().String(),
				field.Exported(),
				st.Tag(i),
			)
		}
	}
}

func docText(group *ast.CommentGroup) string {
	if group == nil {
		return ""
	}
	return group.Text()
}
