package logic

import (
	"fmt"
	"path"

	"github.com/dave/jennifer/jen"
	"github.com/zeebo/xxh3"
)

// PkgPath is a fully-qualified pkg path
type PkgPath string

// Name returns the package name
func (p PkgPath) Name() string {
	return path.Base(string(p))
}

type TypeEntry struct {
	// Type is the type of the TypeEntry
	Type Type
	// Anonymous is true if the type represents an anonymous type. If Anonymous
	// is true, then Private must also be true.
	Anonymous bool
	// Private is true if the type identifiers do not need to be exposed
	Private bool
}

type Package struct {
	Path PkgPath
	// Types maps names to Type
	Types map[string]TypeEntry
}

func NewPackage(path PkgPath) Package {
	return Package{
		Path:  path,
		Types: make(map[string]TypeEntry),
	}
}

type GenContext struct {
	Out *jen.File
}

type Imports struct {
	RelativeTo PkgPath
	Packages   map[PkgPath]Package
	Imported   map[PkgPath]string
}

func NewImports(relativeTo PkgPath, pkgs map[PkgPath]Package) Imports {
	return Imports{
		RelativeTo: relativeTo,
		Packages:   pkgs,
		Imported:   make(map[PkgPath]string),
	}
}

// ResolvePkgName maps a TypeID to the pkg name it is imported as, if it is a
// local type, importName will be ""
func (i Imports) ResolvePkgName(id TypeID) (importName string, err error) {
	if id.Pkg == i.RelativeTo {
		return
	}
	target, ok := i.Packages[id.Pkg]
	if !ok {
		err = fmt.Errorf("package %s was not found", id.Pkg)
		return
	}
	importName = fmt.Sprintf("pkg%d", xxh3.Hash([]byte(id.Pkg)))
	i.Imported[target.Path] = importName
	return
}

// ResolveTypeQual resolves a TypeID to an identifier or a Qual with the
// appropriate package import
func (i Imports) ResolveTypeQual(id TypeID, stmt *jen.Statement) (qual *jen.Statement, err error) {
	importName, err := i.ResolvePkgName(id)
	if err != nil {
		return
	}
	if stmt != nil {
		if importName == "" {
			qual = stmt.Id(id.Name)
		} else {
			qual = stmt.Qual(importName, id.Name)
		}
	} else {
		if importName == "" {
			qual = jen.Id(id.Name)
		} else {
			qual = jen.Qual(importName, id.Name)
		}
	}
	return
}

type Generator struct {
	// Packages maps path to package
	Packages map[PkgPath]Package
}

// ResolveType resolves a TypeID to its type
func (g Generator) ResolveType(id TypeID) (t TypeEntry, err error) {
	pkg, ok := g.Packages[id.Pkg]
	if !ok {
		err = fmt.Errorf("package '%s' not found", id.Pkg)
		return
	}
	t, ok = pkg.Types[id.Name]
	if !ok {
		err = fmt.Errorf("type '%s' not found in package '%s'", id.Name, id.Pkg)
		return
	}
	return
}
