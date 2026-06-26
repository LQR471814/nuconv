package code

import (
	"errors"
	"fmt"
	"go/ast"
	"go/types"
	"iter"
	"nuconv/internal/gen"
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

	UnderlyingPkgs map[gen.PkgPath]*packages.Package
	LogicalPkgs    map[gen.PkgPath]gen.Package
}

func NewParser() Parser {
	return Parser{
		pkgQueue:       map[string]struct{}{},
		UnderlyingPkgs: map[gen.PkgPath]*packages.Package{},
		LogicalPkgs:    map[gen.PkgPath]gen.Package{},
	}
}

func (p Parser) Packages() {
	// TODO: implement
}

type pkgParser struct {
	// reservedTypes are a list of named types either marked with
	// @nutype:manual or are built-in. these types should not have
	// serialize/deserialize code automatically generated for them
	reservedTypes map[gen.TypeID]struct{}
	// markedTypes are a list of types marked with @nutype:auto. these types
	// should have serialize/deserialize code automatically generated for them
	markedTypes map[string]struct{}
	pkg         *packages.Package
	out         *gen.Package
}

type pkgParseContext struct {
	// NearestAncestor stores the name of the nearest named-type ancestor
	NearestAncestor string
	// Path stores the path to the current type from the nearest named-type
	// ancestor
	Path string
}

// WithAccess returns a new copy of pkgParseContext with the path extended to
// cover the new access
func (p pkgParseContext) WithAccess(access string) pkgParseContext {
	return pkgParseContext{
		NearestAncestor: p.NearestAncestor,
		Path:            fmt.Sprintf("%s_%s", p.Path, access),
	}
}

func newPkgParser(pkg *packages.Package) (out pkgParser, err error) {
	errs := make([]error, len(pkg.Errors))
	for i, err := range pkg.Errors {
		errs[i] = errors.New(err.Error())
	}
	if len(errs) > 0 {
		err = fmt.Errorf("package errors: %w", errors.Join(errs...))
		return
	}
	out.pkg = pkg
	out.markedTypes = map[string]struct{}{}
	out.reservedTypes = map[gen.TypeID]struct{}{
		// these are named types we don't
		{
			Pkg:  "time",
			Name: "Time",
		}: {},
		{
			Pkg:  "time",
			Name: "Duration",
		}: {},
	}
	for _, file := range pkg.Syntax {
		for markedType := range iterMarkedTypes(file, marker_auto_gen) {
			out.markedTypes[markedType] = struct{}{}
		}
		for reservedType := range iterMarkedTypes(file, marker_manual_gen) {
			out.reservedTypes[gen.TypeID{
				Pkg:  gen.PkgPath(pkg.PkgPath),
				Name: reservedType,
			}] = struct{}{}
		}
	}
	return
}

func (p pkgParser) Parse() (out gen.Package, err error) {
	p.out = &gen.Package{
		Path:  gen.PkgPath(p.pkg.PkgPath),
		Types: make(map[string]gen.TypeEntry),
	}

	scope := p.pkg.Types.Scope()

	for name := range p.markedTypes {
		obj := scope.Lookup(name)
		typeName, ok := obj.(*types.TypeName)
		if !ok {
			continue
		}
		switch t := typeName.Type().(type) {
		case *types.Named:
			p.out.Types[name], err = p.convertNamedType(t, false)
		case *types.Alias:
			p.out.Types[name], err = p.convertTypeAlias(t, false)
		}
		if err != nil {
			err = fmt.Errorf("convert type top-level (%s): %w", name, err)
			return
		}
	}

	out = *p.out
	return
}

// convertTypeAlias converts a type alias into TypeEntry
func (p pkgParser) convertTypeAlias(alias *types.Alias, private bool) (out gen.TypeEntry, err error) {
	out = gen.TypeEntry{
		Anonymous: false,
		Private:   private,
	}
	out.Type, err = p.convertType(
		pkgParseContext{
			NearestAncestor: alias.Obj().Name(),
		},
		typeNameToID(alias.Obj()),
		alias.Underlying(),
	)
	return
}

// convertNamedType converts a type declaration into TypeEntry
func (p pkgParser) convertNamedType(named *types.Named, private bool) (out gen.TypeEntry, err error) {
	out = gen.TypeEntry{
		Anonymous: false,
		Private:   private,
	}
	out.Type, err = p.convertType(
		pkgParseContext{
			NearestAncestor: named.Obj().Name(),
		},
		typeNameToID(named.Obj()),
		named.Underlying(),
	)
	return
}

// convertType converts types.Type -> logic.Type in the context of ctx with a
// given TypeID
func (p pkgParser) convertType(ctx pkgParseContext, typeID gen.TypeID, t types.Type) (out gen.Type, err error) {
	switch t := t.(type) {
	case *types.Alias, *types.Named:
		err = fmt.Errorf("named or alias types should not be passed to convertType")
		return
	case *types.Chan:
		err = fmt.Errorf("chan is not a supported by nushell plugins")
		return
	case *types.Signature:
		err = fmt.Errorf("function is not supported by nushell plugins")
		return
	case *types.Tuple:
		err = fmt.Errorf("tuple is not supported by nushell plugins")
		return
	case *types.TypeParam, *types.Interface, *types.Union:
		// TODO: implement
		err = fmt.Errorf("interface/generic/union not yet implemented")
		return
	case *types.Basic:
		// we return zero value because a basic type will be handled separately
		// by the generator
		return
	case *types.Pointer:
		var elemTypeID gen.TypeID
		elemTypeID, err = p.resolveTypeID(ctx, t.Elem())
		if err != nil {
			err = fmt.Errorf("resolve type id (%v): %w", typeID, err)
			return
		}
		out = gen.OneofType{
			ID: typeID,
			Alts: []gen.TypeID{
				gen.BuiltinTypeID("nil"),
				elemTypeID,
			},
		}
		return
	case *types.Array:
		var elemTypeID gen.TypeID
		elemTypeID, err = p.resolveTypeID(ctx, t.Elem())
		if err != nil {
			err = fmt.Errorf("resolve type id (%v): %w", typeID, err)
			return
		}
		out = gen.ArrayType{
			ID:          typeID,
			ElementType: elemTypeID,
			Length:      t.Len(),
		}
		return
	case *types.Slice:
		var elemTypeID gen.TypeID
		elemTypeID, err = p.resolveTypeID(ctx, t.Elem())
		if err != nil {
			err = fmt.Errorf("resolve type id (%v): %w", typeID, err)
			return
		}
		out = gen.ListType{
			ID:          typeID,
			ElementType: elemTypeID,
		}
		return
	case *types.Map:
		_, err = p.resolveTypeID(ctx, t.Elem())
		if err != nil {
			err = fmt.Errorf("resolve type id (%v): %w", typeID, err)
			return
		}
		out = gen.RecordType{
			ID:     typeID,
			Fields: nil,
		}
		return
	case *types.Struct:
		record := gen.RecordType{
			ID:     typeID,
			Fields: nil,
		}

		for i := 0; i < t.NumFields(); i++ {
			field := t.Field(i)
			tag := reflect.StructTag(t.Tag(i))

			nuName := tag.Get("nu")
			if nuName == "-" {
				continue
			}
			if nuName == "" {
				nuName = pascalToSnakeCase(field.Name())
			}

			var fieldTypeID gen.TypeID
			fieldTypeID, err = p.resolveTypeID(ctx.WithAccess(field.Name()), field.Type())
			if err != nil {
				err = fmt.Errorf("resolve type id (%v): %w", fieldTypeID, err)
				return
			}

			record.Fields = append(record.Fields, gen.Field{
				GoName: field.Name(),
				NuName: nuName,
				Type:   fieldTypeID,
			})
		}

		out = record
		return
	}

	return
}

// resolveTypeID takes an arbitrary type and resolves it to a logic.TypeID,
// automatically creating an anonymous type if necessary
func (p pkgParser) resolveTypeID(ctx pkgParseContext, t types.Type) (out gen.TypeID, err error) {
	switch t := t.(type) {
	case *types.Named:
		out, err = p.createNamedOrAliasType(t)
		if err != nil {
			err = fmt.Errorf("create named (%s): %w", t.Obj().Name(), err)
			return
		}
		return
	case *types.Alias:
		out, err = p.createNamedOrAliasType(t)
		if err != nil {
			err = fmt.Errorf("create alias (%s): %w", t.Obj().Name(), err)
			return
		}
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
			out = gen.BuiltinTypeID(t.Name())
			return
		}
		err = fmt.Errorf(
			"primitive type '%s' does not have a nushell plugin implementation",
			t.Name(),
		)
		return
	case *types.Pointer:
		out, err = p.createAnonymousType(ctx.WithAccess("Pointer"), t)
		if err != nil {
			err = fmt.Errorf("create anonymous (ptr): %w", err)
			return
		}
		return
	case *types.Array:
		out, err = p.createAnonymousType(ctx.WithAccess("Array"), t)
		if err != nil {
			err = fmt.Errorf("create anonymous (arr): %w", err)
			return
		}
		return
	case *types.Slice:
		out, err = p.createAnonymousType(ctx.WithAccess("Slice"), t)
		if err != nil {
			err = fmt.Errorf("create anonymous (slice): %w", err)
			return
		}
		return
	case *types.Map:
		out, err = p.createAnonymousType(ctx.WithAccess("Map"), t)
		if err != nil {
			err = fmt.Errorf("create anonymous (map): %w", err)
			return
		}
		return
	case *types.Struct:
		out, err = p.createAnonymousType(ctx.WithAccess("Struct"), t)
		if err != nil {
			err = fmt.Errorf("create anonymous (struct): %w", err)
			return
		}
		return
	case *types.Chan:
		err = fmt.Errorf("chan does not have nushell plugin implementation")
		return
	case *types.Signature:
		err = fmt.Errorf("function does not have nushell plugin implementation")
		return
	case *types.Tuple:
		err = fmt.Errorf("tuple does not have nushell plugin implementation")
		return
	case *types.Interface, *types.TypeParam, *types.Union:
		// TODO: implement
		err = fmt.Errorf("interface/generic/union not yet implemented")
		return
	}
	err = fmt.Errorf("unsupported type (%T)", t)
	return
}

// createNamedOrAliasType imports the type as private if it doesn't already exist
func (p pkgParser) createNamedOrAliasType(t types.Type) (id gen.TypeID, err error) {
	switch t := t.(type) {
	case *types.Named:
		id = typeNameToID(t.Obj())
		if _, reserved := p.reservedTypes[id]; reserved {
			return
		}

		name := t.Obj().Name()

		_, exists := p.out.Types[name]
		if exists {
			return
		}

		var typ gen.Type
		typ, err = p.convertType(pkgParseContext{
			NearestAncestor: name,
		}, id, t.Underlying())
		if err != nil {
			err = fmt.Errorf("convert type (%v): %w", id, err)
			return
		}

		p.out.Types[t.Obj().Name()] = gen.TypeEntry{
			Type:      typ,
			Anonymous: false,
			Private:   true,
		}
		return
	case *types.Alias:
		id = typeNameToID(t.Obj())
		if _, reserved := p.reservedTypes[id]; reserved {
			return
		}

		name := t.Obj().Name()

		_, exists := p.out.Types[name]
		if exists {
			return
		}

		var typ gen.Type
		typ, err = p.convertType(pkgParseContext{
			NearestAncestor: name,
		}, id, t.Underlying())
		if err != nil {
			err = fmt.Errorf("convert type (%v): %w", id, err)
			return
		}

		p.out.Types[t.Obj().Name()] = gen.TypeEntry{
			Type:      typ,
			Anonymous: false,
			Private:   true,
		}
		return
	}

	err = fmt.Errorf("unsupported type (%T)", t)
	return
}

func (p pkgParser) createAnonymousType(ctx pkgParseContext, t types.Type) (id gen.TypeID, err error) {
	name := fmt.Sprintf("%s%s", ctx.NearestAncestor, ctx.Path)
	id = gen.TypeID{
		Pkg:  p.out.Path,
		Name: name,
	}

	conv, err := p.convertType(ctx, id, t)
	if err != nil {
		err = fmt.Errorf("convert type (%v): %w", id, err)
		return
	}

	p.out.Types[name] = gen.TypeEntry{
		Type:      conv,
		Anonymous: true,
		Private:   true,
	}
	return
}

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
