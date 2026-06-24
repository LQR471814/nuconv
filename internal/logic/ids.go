// logic holds the business logic of transforming types into nu-plugin
// code
package logic

import (
	"fmt"

	"github.com/dave/jennifer/jen"
)

// PKG_BUILTIN is the package that contains built-in types.
const PKG_BUILTIN = "builtin"

func BuiltinTypeID(builtin string) TypeID {
	return TypeID{
		Pkg:  PKG_BUILTIN,
		Name: builtin,
	}
}

/*
type NuType string

const (
	TYPE_ANY      NuType = "Any"
	TYPE_BINARY          = "Binary"
	TYPE_BLOCK           = "Block"
	TYPE_BOOL            = "Bool"
	TYPE_CELLPATH        = "Cellpath"
	TYPE_CLOSURE         = "Closure"
	TYPE_CUSTOM          = "Custom"
	TYPE_DATE            = "Date"
	TYPE_DURATION        = "Duration"
	TYPE_ERROR           = "Error"
	TYPE_FILESIZE        = "Filesize"
	TYPE_FLOAT           = "Float"
	TYPE_GLOB            = "Glob"
	TYPE_INT             = "Int"
	TYPE_LIST            = "List"
	TYPE_NOTHING         = "Nothing"
	TYPE_NUMBER          = "Number"
	TYPE_ONEOF           = "Oneof"
	TYPE_RANGE           = "Range"
	TYPE_RECORD          = "Record"
	TYPE_STRING          = "String"
	TYPE_TABLE           = "Table"
)
*/

// TypeID represents a fully-qualified type identifier
type TypeID struct {
	Pkg  PkgPath
	Name string
}

type Type interface {
	// ID returns the ID of the type
	TypeID() TypeID
	// Definition returns a statement which declares the type definition
	Definition() *jen.Statement
	// Parser returns a statement which declares the parsing function
	Parser() *jen.Statement
	// Serializer returns a statement which declares the serializer function
	Serializer() *jen.Statement
}

// TypeDefID returns the ID for the type's type definition
func TypeDefID(t TypeID) string {
	return fmt.Sprintf("NuDef%s", t.Name)
}

// ParserID returns the ID for the type's parser
func ParserID(t TypeID) string {
	return fmt.Sprintf("NuParse%s", t.Name)
}

// SerializerID returns the ID for the type's serializer
func SerializerID(t TypeID) string {
	return fmt.Sprintf("NuSerialize%s", t.Name)
}
