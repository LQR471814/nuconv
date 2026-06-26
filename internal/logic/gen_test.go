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

func TestHelpers(t *testing.T) {
	const expected = `package testpkg

import (
	"fmt"
	nuplugin "github.com/ainvaltin/nu-plugin"
)

func nuconvWrapErr[T any](v T) (T, error) {
	return v, nil
}
func tryCast[T any](v nuplugin.Value) (out T, err error) {
	out, ok := v.Value.(T)
	if !ok {
		err = fmt.Errorf("expected %T got %T", out, v)
		return
	}
	return
}
`

	f, err := os.Create(filepath.Join(test_data_dir, "helpers.go"))
	if err != nil {
		return
	}
	defer f.Close()

	g := GenContext{
		Out: jen.NewFilePath("testpkg"),
	}
	wrapErrFn(g)
	tryCastFn(g)

	buf := bytes.NewBuffer(nil)
	err = g.Out.Render(buf)
	if err != nil {
		t.Fatal(err)
	}
	_, err = io.Copy(f, buf)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, "", buf.String())
}
