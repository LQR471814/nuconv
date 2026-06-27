package bar

import (
	"fmt"
	nuplugin "github.com/ainvaltin/nu-plugin"
	types "github.com/ainvaltin/nu-plugin/types"
)

func nuconvWrapErr[T any](v T) (T, error) {
	return v, nil
}
func tryCast[T any](v any) (out T, err error) {
	out, ok := v.(T)
	if !ok {
		err = fmt.Errorf("expected %T got %T", out, v)
		return
	}
	return
}

type Foo_Bar_Pointer = *string

var NuDefFoo_Bar_Pointer = types.OneOf(types.Nothing(), types.String())

func NuParseFoo_Bar_Pointer(v nuplugin.Value) (out Foo_Bar_Pointer, err error) {
	if v.Value == nil {
		return
	}
	var tmp string
	tmp, err = tryCast[string](v.Value)
	if err != nil {
		err = fmt.Errorf("parse underlying: %w", err)
		return
	}
	out = &tmp
	return
}
func NuSerializeFoo_Bar_Pointer(v Foo_Bar_Pointer) (out nuplugin.Value, err error) {
	if v == nil {
		out = nuplugin.ToValue(nil)
		return
	}
	out, err = nuconvWrapErr(nuplugin.ToValue(*v))
	if err != nil {
		err = fmt.Errorf("serialize underlying: %w", err)
		return
	}
	return
}

type Foo_Baz_Slice = []int

var NuDefFoo_Baz_Slice = types.List(types.Int())

func NuParseFoo_Baz_Slice(v nuplugin.Value) (out Foo_Baz_Slice, err error) {
	typed, err := tryCast[[]nuplugin.Value](v.Value)
	if err != nil {
		err = fmt.Errorf("cast: %w", err)
		return
	}
	out = make(Foo_Baz_Slice, len(typed))
	for i, e := range typed {
		out[i], err = tryCast[int](e.Value)
		if err != nil {
			err = fmt.Errorf("parse element (%d): %w", i, err)
			return
		}
	}
	return
}
func NuSerializeFoo_Baz_Slice(v Foo_Baz_Slice) (out nuplugin.Value, err error) {
	tmp := make([]nuplugin.Value, len(v))
	for i, e := range v {
		tmp[i], err = nuconvWrapErr(nuplugin.ToValue(e))
		if err != nil {
			err = fmt.Errorf("serialize element (%d): %w", i, err)
			return
		}
	}
	out = nuplugin.ToValue(tmp)
	return
}

var NuDefFoo = types.Record(types.RecordDef{
	"bar": NuDefFoo_Bar_Pointer,
	"baz": NuDefFoo_Baz_Slice,
})

func NuParseFoo(v nuplugin.Value) (out Foo, err error) {
	typed, err := tryCast[map[string]nuplugin.Value](v.Value)
	if err != nil {
		err = fmt.Errorf("cast: %w", err)
		return
	}
	out.Bar, err = NuParseFoo_Bar_Pointer(typed["bar"])
	if err != nil {
		err = fmt.Errorf("parse Bar (bar): %w", err)
		return
	}
	out.Baz, err = NuParseFoo_Baz_Slice(typed["baz"])
	if err != nil {
		err = fmt.Errorf("parse Baz (baz): %w", err)
		return
	}
	return
}
func NuSerializeFoo(v Foo) (out nuplugin.Value, err error) {
	tmp := make(map[string]nuplugin.Value)
	tmp["bar"], err = NuSerializeFoo_Bar_Pointer(v.Bar)
	if err != nil {
		err = fmt.Errorf("serialize Bar (bar): %w", err)
		return
	}
	tmp["baz"], err = NuSerializeFoo_Baz_Slice(v.Baz)
	if err != nil {
		err = fmt.Errorf("serialize Baz (baz): %w", err)
		return
	}
	out = nuplugin.ToValue(tmp)
	return
}
