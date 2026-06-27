package gen

import (
	"path"

	"github.com/dave/jennifer/jen"
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
