package logic

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/dave/jennifer/jen"
	"github.com/stretchr/testify/require"
)

const test_data_dir = "test_data"

func setupTestData(t *testing.T) {
	err := os.MkdirAll(test_data_dir, 0777)
	if err != nil {
		t.Fatal(err)
	}
}

func testType(t *testing.T, outputFile string, typ Type, expected string) {
	setupTestData(t)

	f, err := os.Create(filepath.Join(test_data_dir, outputFile))
	if err != nil {
		return
	}
	defer f.Close()

	c := GenContext{
		Out: jen.NewFilePath("testpkg"),
	}

	typ.Definition(c)
	typ.Parser(c)
	typ.Serializer(c)

	buf := bytes.NewBuffer(nil)
	err = c.Out.Render(buf)
	if err != nil {
		t.Fatal(err)
	}
	str := buf.String()

	_, err = io.Copy(f, bytes.NewBufferString(str))
	if err != nil {
		t.Fatal(err)
	}
	require.Equal(t, expected, str)
}

func TestArrayType(t *testing.T) {
	const expect = `package testpkg

import (
	"fmt"
	nuplugin "github.com/ainvaltin/nu-plugin"
	types "github.com/ainvaltin/nu-plugin/types"
)

var NuDefArray = types.List(types.String())

func NuParseArray(v nuplugin.Value) (out Array, err error) {
	typed, err := tryCast[[]nuplugin.Value](v)
	if err != nil {
		err = fmt.Errorf("cast: %w", err)
		return
	}
	if len(typed) != 3 {
		err = fmt.Errorf("array has wrong number of elements, got (%d), expected (%d)", len(typed), 3)
		return
	}
	out = make(Array)
	for i := range 3 {
		out[i], err = nuconvWrapErr(string(typed[i]))
		if err != nil {
			err = fmt.Errorf("parse element (%d): %w", i, err)
			return
		}
	}
	return
}
func NuSerializeArray(v Array) (out nuplugin.Value, err error) {
	tmp := make([]nuplugin.Value, 3)
	for i := range 3 {
		tmp[i], err = nuconvWrapErr(nuplugin.ToValue(v[i]))
		if err != nil {
			err = fmt.Errorf("serialize element (%d): %w", i, err)
			return
		}
	}
	out = nuplugin.ToValue(tmp)
	return
}
`

	foo := ArrayType{
		ID: TypeID{
			Pkg:  "testpkg",
			Name: "Array",
		},
		Length:      3,
		ElementType: BuiltinTypeID("string"),
	}

	testType(t, "array.go", foo, expect)
}
