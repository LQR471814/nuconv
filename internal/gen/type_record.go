package gen

import (
	"fmt"

	"github.com/dave/jennifer/jen"
)

// RecordType maps go struct <-> nu record
type RecordType struct {
	ID     TypeID
	Fields []Field
}

// ID returns the ID of the type
func (recordtype RecordType) TypeID() TypeID {
	return recordtype.ID
}

// Definition returns a statement which declares the type definition
func (recordtype RecordType) Definition(g GenContext) (out *jen.Statement, err error) {
	typeFields := jen.Dict{}
	for _, field := range recordtype.Fields {
		typeFields[jen.Lit(field.NuName)] = field.Type.DefQual(nil)
	}
	out = defTemplate(g, recordtype).
		Qual(nu_types_pkg, "Record").
		Call(jen.Qual(nu_types_pkg, "RecordDef").Values(typeFields))
	return
}

// Parser returns a statement which declares the parsing function
func (recordtype RecordType) Parser(g GenContext) (out *jen.Statement, err error) {
	cast, errHandle := castValueBlock(jen.Map(jen.String()).Qual(nu_pkg, "Value"))

	block := []jen.Code{cast, errHandle}
	for _, field := range recordtype.Fields {
		block = append(block,
			field.Type.ParserQual(
				jen.List(jen.Id(id_out).Dot(field.GoName), jen.Id(id_err_val)).
					Op("="),
				jen.Id(id_typed).Index(jen.Lit(field.NuName)),
			),
			jen.If(jen.Id(id_err_val).Op("!=").Nil()).Block(
				jen.Id(id_err_val).Op("=").Qual("fmt", "Errorf").Call(
					jen.Lit(fmt.Sprintf(
						"parse %s (%s): %%w",
						field.GoName,
						field.NuName,
					)),
					jen.Id(id_err_val),
				),
				jen.Return(),
			),
		)
	}
	block = append(block, jen.Return())

	out = parserTypeTemplate(
		g,
		recordtype,
		block...,
	)
	return
}

// Serializer returns a statement which declares the serializer function
func (recordtype RecordType) Serializer(g GenContext) (out *jen.Statement, err error) {
	block := []jen.Code{
		jen.Id(id_tmp).Op(":=").Make(jen.Map(jen.String()).Qual(nu_pkg, "Value")),
	}
	for _, field := range recordtype.Fields {
		block = append(block,
			field.Type.SerializerQual(
				jen.List(
					jen.Id(id_tmp).Index(jen.Lit(field.NuName)),
					jen.Id(id_err_val),
				).Op("="),
				jen.Id(id_value).Dot(field.GoName),
			),
			jen.If(jen.Id(id_err_val).Op("!=").Nil()).Block(
				jen.Id(id_err_val).Op("=").Qual("fmt", "Errorf").Call(
					jen.Lit(fmt.Sprintf(
						"serialize %s (%s): %%w",
						field.GoName,
						field.NuName,
					)),
					jen.Id(id_err_val),
				),
				jen.Return(),
			),
		)
	}
	block = append(block,
		jen.Id(id_out).Op("=").Qual(nu_pkg, "ToValue").Call(jen.Id(id_tmp)),
		jen.Return(),
	)
	out = serializerTypeTemplate(g, recordtype, block...)
	return
}

// Field maps a go struct field <-> nu field
type Field struct {
	// NuName is the name of the nushell-facing field name.
	NuName string
	// GoName is the golang-facing field name.
	GoName string
	// Type is the TypeID of the field's type.
	Type TypeID
}
