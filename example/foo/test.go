package foo

import (
	"nuconv/example/bar"
	"time"

	"github.com/ainvaltin/nu-plugin"
)

// this is a doc comment
//
// @nuconv:auto
type Structure struct {
	Name      string `nu:"nu_name"`
	Foo       bar.Foo
	Float     float64
	Bool      bool
	Time      time.Time
	Duration  time.Duration
	Map       map[string]string
	Anonymous struct {
		Field string
	}
	Ignored chan bool `nu:"-" json:"-"`
}

// @nuconv:manual
type Custom struct {
}

func NuParseCustom(c Custom) (nu.Value, error) {
	return nu.ToValue(0), nil
}

func NuSerializeCustom(nu.Value) (Custom, error) {
	return Custom{}, nil
}
