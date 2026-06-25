package code

import (
	"encoding/json"
	"errors"
	"fmt"
	"go/token"
	"nuconv/internal/logic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/packages"
)

func loadExample() (pkg *packages.Package, err error) {
	cfg := &packages.Config{
		Mode: packages.NeedName |
			packages.NeedFiles |
			packages.NeedSyntax |
			packages.NeedTypes |
			packages.NeedTypesInfo,
		Fset: token.NewFileSet(),
		Dir:  ".",
	}

	pkgs, err := packages.Load(cfg, "../../example")
	if err != nil {
		return
	}
	if len(pkgs) != 1 {
		err = fmt.Errorf("incorrect number of packages loaded (%d)", len(pkgs))
		return
	}
	pkg = pkgs[0]

	errs := make([]error, len(pkg.Errors))
	for i, err := range pkg.Errors {
		errs[i] = errors.New(err.Error())
	}
	if len(errs) > 0 {
		err = fmt.Errorf("package errors: %w", errors.Join(errs...))
		return
	}

	return
}

func TestPkgParser(t *testing.T) {
	pkg, err := loadExample()
	if err != nil {
		err = fmt.Errorf("load example: %w", err)
		t.Fatal(err)
	}

	parser, err := newPkgParser(pkg)
	if err != nil {
		err = fmt.Errorf("new pkg parser: %w", err)
		t.Fatal(err)
	}

	out, err := parser.Parse()
	if err != nil {
		err = fmt.Errorf("parse: %w", err)
		t.Fatal(err)
	}

	const example_pkg = "nuconv/example"

	expected := logic.Package{
		Path: example_pkg,
		Types: map[string]logic.TypeEntry{
			"Structure": {
				Type: logic.RecordType{
					ID: logic.TypeID{
						Pkg:  example_pkg,
						Name: "Structure",
					},
					Fields: []logic.Field{
						{
							GoName: "Name",
							Type:   logic.BuiltinTypeID("string"),
							NuName: "nu_name",
						},
						{
							GoName: "Foo",
							Type:   logic.TypeID{Pkg: example_pkg, Name: "Foo"},
							NuName: "foo",
						},
						{
							GoName: "Float",
							Type:   logic.BuiltinTypeID("float64"),
							NuName: "float",
						},
						{
							GoName: "Bool",
							Type:   logic.BuiltinTypeID("bool"),
							NuName: "bool",
						},
						{
							GoName: "Time",
							Type:   logic.TypeID{Pkg: "time", Name: "Time"},
							NuName: "time",
						},
						{
							GoName: "Duration",
							Type:   logic.TypeID{Pkg: "time", Name: "Duration"},
							NuName: "duration",
						},
						{
							GoName: "Map",
							Type: logic.TypeID{
								Pkg:  example_pkg,
								Name: "Structure_Map_Map",
							},
							NuName: "map",
						},
						{
							GoName: "Anonymous",
							Type: logic.TypeID{
								Pkg:  example_pkg,
								Name: "Structure_Anonymous_Struct",
							},
							NuName: "anonymous",
						},
					},
				},
				Anonymous: false,
				Private:   false,
			},
			"Structure_Map_Map": {
				Type: logic.RecordType{
					ID: logic.TypeID{
						Pkg:  example_pkg,
						Name: "Structure_Map_Map",
					},
				},
				Anonymous: true,
				Private:   true,
			},
			"Structure_Anonymous_Struct": {
				Type: logic.RecordType{
					ID: logic.TypeID{
						Pkg:  example_pkg,
						Name: "Structure_Anonymous_Struct",
					},
					Fields: []logic.Field{
						{
							GoName: "Field",
							Type:   logic.BuiltinTypeID("string"),
							NuName: "field",
						},
					},
				},
				Anonymous: true,
				Private:   true,
			},
			"Foo": {
				Type: logic.RecordType{
					ID: logic.TypeID{
						Pkg:  example_pkg,
						Name: "Foo",
					},
					Fields: []logic.Field{
						{
							GoName: "Bar",
							Type: logic.TypeID{
								Pkg:  example_pkg,
								Name: "Foo_Bar_Pointer",
							},
							NuName: "bar",
						},
						{
							GoName: "Baz",
							Type: logic.TypeID{
								Pkg:  example_pkg,
								Name: "Foo_Baz_Slice",
							},
							NuName: "baz",
						},
					},
				},
				Anonymous: false,
				Private:   true,
			},
			"Foo_Bar_Pointer": {
				Type: logic.OneofType{
					ID: logic.TypeID{
						Pkg:  example_pkg,
						Name: "Foo_Bar_Pointer",
					},
					Alternatives: []logic.TypeID{
						{Pkg: "builtin", Name: "nil"},
						{Pkg: "builtin", Name: "string"},
					},
				},
				Anonymous: true,
				Private:   true,
			},
			"Foo_Baz_Slice": {
				Type: logic.ListType{
					ID: logic.TypeID{
						Pkg:  example_pkg,
						Name: "Foo_Baz_Slice",
					},
					ElementType: logic.TypeID{
						Pkg:  "builtin",
						Name: "int",
					},
				},
				Anonymous: true,
				Private:   true,
			},
		},
	}

	expect, err := json.MarshalIndent(expected, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	got, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	if !assert.Equal(t, expect, got) {
		t.Log("expected", string(expect))
		t.Log("got", string(got))
	}
	require.Equal(t, expected, out)
}
