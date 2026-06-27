package foo

import (
	"fmt"
	nuplugin "github.com/ainvaltin/nu-plugin"
	types "github.com/ainvaltin/nu-plugin/types"
	bar "nuconv/example/bar"
	"time"
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

var NuDefFoo = types.Record(types.RecordDef{
	"bar": NuDefFoo_Bar_Pointer,
	"baz": NuDefFoo_Baz_Slice,
})

func NuParseFoo(v nuplugin.Value) (out bar.Foo, err error) {
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
func NuSerializeFoo(v bar.Foo) (out nuplugin.Value, err error) {
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

type Structure_Map_Map = map[string]string

var NuDefStructure_Map_Map = types.Record(types.RecordDef{})

func NuParseStructure_Map_Map(v nuplugin.Value) (out Structure_Map_Map, err error) {
	typed, err := tryCast[map[string]nuplugin.Value](v.Value)
	if err != nil {
		err = fmt.Errorf("cast: %w", err)
		return
	}
	out = make(map[string]string)
	for key, e := range typed {
		out[key], err = tryCast[string](e.Value)
		if err != nil {
			err = fmt.Errorf("parse map entry (%s:%v): %w", key, e, err)
			return
		}
	}
	return
}
func NuSerializeStructure_Map_Map(v Structure_Map_Map) (out nuplugin.Value, err error) {
	tmp := make(map[string]nuplugin.Value)
	for key, e := range v {
		tmp[key], err = nuconvWrapErr(nuplugin.ToValue(e))
		if err != nil {
			err = fmt.Errorf("serialize map entry (%s:%v): %w", key, e, err)
			return
		}
	}
	out = nuplugin.ToValue(tmp)
	return
}

type Structure_Anonymous_Struct = struct {
	Field string
}

var NuDefStructure_Anonymous_Struct = types.Record(types.RecordDef{"field": types.String()})

func NuParseStructure_Anonymous_Struct(v nuplugin.Value) (out Structure_Anonymous_Struct, err error) {
	typed, err := tryCast[map[string]nuplugin.Value](v.Value)
	if err != nil {
		err = fmt.Errorf("cast: %w", err)
		return
	}
	out.Field, err = tryCast[string](typed["field"].Value)
	if err != nil {
		err = fmt.Errorf("parse Field (field): %w", err)
		return
	}
	return
}
func NuSerializeStructure_Anonymous_Struct(v Structure_Anonymous_Struct) (out nuplugin.Value, err error) {
	tmp := make(map[string]nuplugin.Value)
	tmp["field"], err = nuconvWrapErr(nuplugin.ToValue(v.Field))
	if err != nil {
		err = fmt.Errorf("serialize Field (field): %w", err)
		return
	}
	out = nuplugin.ToValue(tmp)
	return
}

var NuDefStructure = types.Record(types.RecordDef{
	"anonymous": NuDefStructure_Anonymous_Struct,
	"bool":      types.Bool(),
	"duration":  types.Duration(),
	"float":     types.Float(),
	"foo":       bar.NuDefFoo,
	"map":       NuDefStructure_Map_Map,
	"nu_name":   types.String(),
	"time":      types.Date(),
})

func NuParseStructure(v nuplugin.Value) (out Structure, err error) {
	typed, err := tryCast[map[string]nuplugin.Value](v.Value)
	if err != nil {
		err = fmt.Errorf("cast: %w", err)
		return
	}
	out.Name, err = tryCast[string](typed["nu_name"].Value)
	if err != nil {
		err = fmt.Errorf("parse Name (nu_name): %w", err)
		return
	}
	out.Foo, err = bar.NuParseFoo(typed["foo"])
	if err != nil {
		err = fmt.Errorf("parse Foo (foo): %w", err)
		return
	}
	out.Float, err = tryCast[float64](typed["float"].Value)
	if err != nil {
		err = fmt.Errorf("parse Float (float): %w", err)
		return
	}
	out.Bool, err = tryCast[bool](typed["bool"].Value)
	if err != nil {
		err = fmt.Errorf("parse Bool (bool): %w", err)
		return
	}
	out.Time, err = tryCast[time.Time](typed["time"].Value)
	if err != nil {
		err = fmt.Errorf("parse Time (time): %w", err)
		return
	}
	out.Duration, err = tryCast[time.Duration](typed["duration"].Value)
	if err != nil {
		err = fmt.Errorf("parse Duration (duration): %w", err)
		return
	}
	out.Map, err = NuParseStructure_Map_Map(typed["map"])
	if err != nil {
		err = fmt.Errorf("parse Map (map): %w", err)
		return
	}
	out.Anonymous, err = NuParseStructure_Anonymous_Struct(typed["anonymous"])
	if err != nil {
		err = fmt.Errorf("parse Anonymous (anonymous): %w", err)
		return
	}
	return
}
func NuSerializeStructure(v Structure) (out nuplugin.Value, err error) {
	tmp := make(map[string]nuplugin.Value)
	tmp["nu_name"], err = nuconvWrapErr(nuplugin.ToValue(v.Name))
	if err != nil {
		err = fmt.Errorf("serialize Name (nu_name): %w", err)
		return
	}
	tmp["foo"], err = bar.NuSerializeFoo(v.Foo)
	if err != nil {
		err = fmt.Errorf("serialize Foo (foo): %w", err)
		return
	}
	tmp["float"], err = nuconvWrapErr(nuplugin.ToValue(v.Float))
	if err != nil {
		err = fmt.Errorf("serialize Float (float): %w", err)
		return
	}
	tmp["bool"], err = nuconvWrapErr(nuplugin.ToValue(v.Bool))
	if err != nil {
		err = fmt.Errorf("serialize Bool (bool): %w", err)
		return
	}
	tmp["time"], err = nuconvWrapErr(nuplugin.ToValue(v.Time))
	if err != nil {
		err = fmt.Errorf("serialize Time (time): %w", err)
		return
	}
	tmp["duration"], err = nuconvWrapErr(nuplugin.ToValue(v.Duration))
	if err != nil {
		err = fmt.Errorf("serialize Duration (duration): %w", err)
		return
	}
	tmp["map"], err = NuSerializeStructure_Map_Map(v.Map)
	if err != nil {
		err = fmt.Errorf("serialize Map (map): %w", err)
		return
	}
	tmp["anonymous"], err = NuSerializeStructure_Anonymous_Struct(v.Anonymous)
	if err != nil {
		err = fmt.Errorf("serialize Anonymous (anonymous): %w", err)
		return
	}
	out = nuplugin.ToValue(tmp)
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
