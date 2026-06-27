package pkgs

import (
	"encoding/json"
	"errors"
	"fmt"
	"go/token"
	"nuconv/internal/gen"
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

	expected := gen.Package{
		Path: example_pkg,
		Types: map[string]gen.TypeEntry{
			"Structure": {
				Type: gen.RecordType{
					ID: gen.TypeID{
						Pkg:  example_pkg,
						Name: "Structure",
					},
					Fields: []gen.Field{
						{
							GoName: "Name",
							Type:   gen.BuiltinTypeID("string"),
							NuName: "nu_name",
						},
						{
							GoName: "Foo",
							Type:   gen.TypeID{Pkg: example_pkg, Name: "Foo"},
							NuName: "foo",
						},
						{
							GoName: "Float",
							Type:   gen.BuiltinTypeID("float64"),
							NuName: "float",
						},
						{
							GoName: "Bool",
							Type:   gen.BuiltinTypeID("bool"),
							NuName: "bool",
						},
						{
							GoName: "Time",
							Type:   gen.TypeID{Pkg: "time", Name: "Time"},
							NuName: "time",
						},
						{
							GoName: "Duration",
							Type:   gen.TypeID{Pkg: "time", Name: "Duration"},
							NuName: "duration",
						},
						{
							GoName: "Map",
							Type: gen.TypeID{
								Pkg:  example_pkg,
								Name: "Structure_Map_Map",
							},
							NuName: "map",
						},
						{
							GoName: "Anonymous",
							Type: gen.TypeID{
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
				Type: gen.RecordType{
					ID: gen.TypeID{
						Pkg:  example_pkg,
						Name: "Structure_Map_Map",
					},
				},
				Anonymous: true,
				Private:   true,
			},
			"Structure_Anonymous_Struct": {
				Type: gen.RecordType{
					ID: gen.TypeID{
						Pkg:  example_pkg,
						Name: "Structure_Anonymous_Struct",
					},
					Fields: []gen.Field{
						{
							GoName: "Field",
							Type:   gen.BuiltinTypeID("string"),
							NuName: "field",
						},
					},
				},
				Anonymous: true,
				Private:   true,
			},
			"Foo": {
				Type: gen.RecordType{
					ID: gen.TypeID{
						Pkg:  example_pkg,
						Name: "Foo",
					},
					Fields: []gen.Field{
						{
							GoName: "Bar",
							Type: gen.TypeID{
								Pkg:  example_pkg,
								Name: "Foo_Bar_Pointer",
							},
							NuName: "bar",
						},
						{
							GoName: "Baz",
							Type: gen.TypeID{
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
				Type: gen.OptionalType{
					ID: gen.TypeID{
						Pkg:  example_pkg,
						Name: "Foo_Bar_Pointer",
					},
					Elem: gen.BuiltinTypeID("string"),
				},
				Anonymous: true,
				Private:   true,
			},
			"Foo_Baz_Slice": {
				Type: gen.ListType{
					ID: gen.TypeID{
						Pkg:  example_pkg,
						Name: "Foo_Baz_Slice",
					},
					ElementType: gen.BuiltinTypeID("int"),
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
