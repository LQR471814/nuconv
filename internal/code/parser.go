package code

import (
	"errors"
	"fmt"
	"go/ast"
	"go/types"
	"iter"
	"nu-plugin-gen/internal/logic"
	"reflect"
	"strings"

	"golang.org/x/tools/go/packages"
)

const (
	marker_auto_gen   = "@nutype:auto"
	marker_manual_gen = "@nutype:manual"
)

type Parser struct {
	// pkgQueue stores queued packages to-be-parsed
	pkgQueue map[string]struct{}

	UnderlyingPkgs map[logic.PkgPath]*packages.Package
	LogicalPkgs    map[logic.PkgPath]logic.Package
}

func NewParser() Parser {
	return Parser{
		pkgQueue:       map[string]struct{}{},
		UnderlyingPkgs: map[logic.PkgPath]*packages.Package{},
		LogicalPkgs:    map[logic.PkgPath]logic.Package{},
	}
}

func (p Parser) Packages() {
}

type pkgParser struct {
	markedTypes map[string]struct{}
	pkg         *packages.Package
	out         *logic.Package
}

func newPkgParser(pkg *packages.Package) (out pkgParser, err error) {
	errs := make([]error, len(pkg.Errors))
	for i, err := range pkg.Errors {
		errs[i] = fmt.Errorf(err.Error())
	}
	if len(errs) > 0 {
		err = fmt.Errorf("package errors: %w", errors.Join(errs...))
		return
	}
	out.markedTypes = map[string]struct{}{}
	for _, file := range pkg.Syntax {
		for markedType := range iterMarkedTypes(file, marker_auto_gen) {
			out.markedTypes[markedType] = struct{}{}
		}
	}
	return
}

func (p pkgParser) Parse() (out logic.Package, err error) {
	p.out = &logic.Package{
		Path:  logic.PkgPath(p.pkg.PkgPath),
		Types: make(map[string]logic.TypeEntry),
	}

	scope := p.pkg.Types.Scope()

	for name := range p.markedTypes {
		obj := scope.Lookup(name)
		typeName, ok := obj.(*types.TypeName)
		if !ok {
			continue
		}
		named, ok := typeName.Type().(*types.Named)
		if !ok {
			continue
		}
		out.Types[name] = p.convertType(named)
	}

	return
}

func (p pkgParser) convertType(named *types.Named) (out logic.TypeEntry) {
	typeID := typeNameToID(named.Obj())

	switch underlying := named.Underlying().(type) {
	case *types.Basic:
		// we return zero value because a basic type will be handled separately
		// by the generator
		return
	case *types.Interface:
		// TODO: implement
		return
	case *types.Pointer:
	case *types.Map:
	case *types.Array:
	case *types.Slice:
	case *types.Struct:
		record := logic.RecordType{
			ID:     typeID,
			Fields: make([]logic.Field, underlying.NumFields()),
		}

		for i := 0; i < underlying.NumFields(); i++ {
			field := underlying.Field(i)
			tag := reflect.StructTag(underlying.Tag(i))

			nuName := tag.Get("nu")
			if nuName == "-" {
				continue
			}
			if nuName == "" {
				nuName = pascalToSnakeCase(field.Name())
			}

			record.Fields[i] = logic.Field{
				GoName: field.Name(),
				NuName: nuName,
			}

			fmt.Printf(
				"  field: %s type=%s exported=%v tag=%q\n",
				field.Name(),
				field.Type().String(),
				field.Exported(),
				underlying.Tag(i),
			)
		}

		// out = record
	}

	return
}

// resolveTypeID takes an arbitrary type and resolves it to a logic.TypeID,
// automatically creating an anonymous type if necessary
func (p pkgParser) resolveTypeID(t types.Type, parentName string) (out logic.TypeID, err error) {
	switch t := t.(type) {
	case *types.Named:
		out = typeNameToID(t.Obj())
		return
	case *types.Alias:
		out = typeNameToID(t.Obj())
		return
	case *types.Basic:
		switch t.Name() {
		case "any",
			"bool",
			"string",
			"int",
			"int8",
			"int16",
			"int32",
			"int64",
			"uint",
			"uint8",
			"uint16",
			"uint32",
			"uint64",
			"float32",
			"float64":
			out = logic.BuiltinTypeID(t.Name())
			return
		}
		err = fmt.Errorf(
			"primitive type '%s' does not have a nushell plugin implementation",
			t.Name(),
		)
		return
	case *types.Array:
	case *types.Struct:
	case *types.Pointer:
	case *types.Slice:
	case *types.Map:
	case *types.Chan:
		err = fmt.Errorf("chan does not have nushell plugin implementation")
		return
	case *types.Signature:
		err = fmt.Errorf("function does not have nushell plugin implementation")
		return
	case *types.Tuple:
		err = fmt.Errorf("tuple does not have nushell plugin implementation")
		return
	case *types.Interface, *types.TypeParam:
		// TODO: implement
		err = fmt.Errorf("interface/generic not yet implemented")
		return
	}
	err = fmt.Errorf("unsupported type (%T)", t)
	return
}

func (p pkgParser) createAnonymousType(t types.Type, field, parentName string) (id logic.TypeID, err error) {
	name := fmt.Sprintf("%s_%s", parentName, field)
	id = logic.TypeID{
		Pkg:  "",
		Name: name,
	}

	switch t := t.(type) {
	case *types.Array:
		var elemTypeID logic.TypeID
		elemTypeID, err = p.resolveTypeID(t, name)
		if err != nil {
			return
		}
		p.out.Types[name] = logic.TypeEntry{
			Type: logic.ArrayType{
				ID:          id,
				ElementType: elemTypeID,
				Length:      t.Len(),
			},
			Anonymous: true,
			Private:   true,
		}
		return
	case *types.Struct:
		p.out.Types[name] = logic.TypeEntry{
			Type: logic.RecordType{
				ID: id,
			},
			Anonymous: true,
			Private:   true,
		}
		return
	case *types.Pointer:
		var elemTypeID logic.TypeID
		elemTypeID, err = p.resolveTypeID(t, name)
		if err != nil {
			return
		}
		p.out.Types[name] = logic.TypeEntry{
			Type: logic.OneofType{
				ID: id,
				Alternatives: []logic.TypeID{
					elemTypeID,
				},
			},
			Anonymous: true,
			Private:   true,
		}
		return
	case *types.Slice:
		p.out.Types[name] = logic.TypeEntry{
			Anonymous: true,
			Private:   true,
		}
		return
	case *types.Map:
		p.out.Types[name] = logic.TypeEntry{
			Anonymous: true,
			Private:   true,
		}
		return
	}
	err = fmt.Errorf("cannot create anonymous type of '%T'", t)
	return
}

func typeNameToID(typeName *types.TypeName) logic.TypeID {
	return logic.TypeID{
		Pkg:  logic.PkgPath(typeName.Pkg().Path()),
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
