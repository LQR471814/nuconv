package code

import (
	"nuconv/internal/logic"

	"golang.org/x/tools/go/packages"
)

type Serializer struct {
	UnderlyingPkgs map[logic.PkgPath]*packages.Package
	LogicalPkgs    map[logic.PkgPath]logic.Package
}


