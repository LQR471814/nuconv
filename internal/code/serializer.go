package code

import (
	"nuconv/internal/gen"

	"golang.org/x/tools/go/packages"
)

type Serializer struct {
	UnderlyingPkgs map[gen.PkgPath]*packages.Package
	LogicalPkgs    map[gen.PkgPath]gen.Package
}
