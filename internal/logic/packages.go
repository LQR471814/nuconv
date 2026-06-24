package logic

import (
	"fmt"
	"path"
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
