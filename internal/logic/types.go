// logic holds the business logic of transforming types into nu-plugin
// code
package logic

import (
	"fmt"

	"github.com/dave/jennifer/jen"
)

func newStatement() *jen.Statement {
	return &jen.Statement{}
}

func BuiltinTypeID(builtin string) TypeID {
	return TypeID{Name: builtin}
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

// Qual creates a jen.Statement which accesses the type qualified
// by its package name
//
// NOTE: this works even if the Pkg == the current package the
// file is being generated in
func (id TypeID) Qual(prev *jen.Statement) *jen.Statement {
	if prev == nil {
		prev = newStatement()
	}
	// empty pkg indicates a builtin type
	if id.Pkg == "" {
		return prev.Id(id.Name)
	}
	return prev.Qual(string(id.Pkg), id.Name)
}

// DefQual creates a jen.Statement which accesses a type's nu type definition
func (id TypeID) DefQual(prev *jen.Statement) *jen.Statement {
	if prev == nil {
		prev = newStatement()
	}
	// empty pkg indicates a builtin type
	switch id.Pkg {
	case "":
		switch id.Name {
		case "string":
			return prev.Qual(nu_types_pkg, "String").Call()
		case "bool":
			return prev.Qual(nu_types_pkg, "Bool").Call()
		case "int", "int8", "int16", "int32", "int64",
			"uint", "uint8", "uint16", "uint32", "uint64":
			return prev.Qual(nu_types_pkg, "Int").Call()
		case "float32", "float64":
			return prev.Qual(nu_types_pkg, "Float").Call()
		case "error":
			return prev.Qual(nu_types_pkg, "Error").Call()
		}
	case "time":
		switch id.Name {
		case "Time":
			return prev.Qual(nu_types_pkg, "Date").Call()
		case "Duration":
			return prev.Qual(nu_types_pkg, "Duration").Call()
		}
	}
	return prev.Qual(string(id.Pkg), NewTypeDefID(id))
}

// ParserQual creates a jen.Statement which accesses a type's nu parser
// function
func (id TypeID) ParserQual(prev, value *jen.Statement) *jen.Statement {
	if prev == nil {
		prev = newStatement()
	}
	// empty pkg indicates a builtin type
	if id.Pkg == "" {
		return prev.Id(id_try_cast_fn).Types(jen.Id(id.Name)).Call(value)
	}
	return prev.Qual(string(id.Pkg), NewParserID(id)).Call(value)
}

// SerializerQual creates a jen.Statement which accesses a type's nu serializer
// function
func (id TypeID) SerializerQual(prev, value *jen.Statement) *jen.Statement {
	if prev == nil {
		prev = newStatement()
	}
	// empty pkg indicates a builtin type
	if id.Pkg == "" {
		return prev.Id(id_wrap_err_fn).Call(jen.Qual(nu_pkg, "ToValue").Call(value))
	}
	return prev.Qual(string(id.Pkg), NewSerializerID(id)).Call(value)
}

func (id TypeID) String() string {
	if id.Pkg == "" {
		return id.Name
	}
	return fmt.Sprintf("%s.%s", id.Pkg, id.Name)
}

type Type interface {
	// ID returns the ID of the type
	TypeID() TypeID
	// Definition returns a statement which declares the type definition
	//
	// NOTE: the reason why it returns an entire statement instead of just the
	// expression is so that if any additional variables necessary for the
	// definition must be declared, they can be declared
	Definition(g GenContext) (*jen.Statement, error)
	// Parser returns a statement which declares the parsing function
	Parser(g GenContext) (*jen.Statement, error)
	// Serializer returns a statement which declares the serializer function
	Serializer(g GenContext) (*jen.Statement, error)
}

// NewTypeDefID returns the ID for the type's type definition
func NewTypeDefID(t TypeID) string {
	return fmt.Sprintf("NuDef%s", t.Name)
}

// NewParserID returns the ID for the type's parser
func NewParserID(t TypeID) string {
	return fmt.Sprintf("NuParse%s", t.Name)
}

// NewSerializerID returns the ID for the type's serializer
func NewSerializerID(t TypeID) string {
	return fmt.Sprintf("NuSerialize%s", t.Name)
}
