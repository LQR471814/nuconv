# nuconv

`nuconv` is a Go code generator for
[`github.com/ainvaltin/nu-plugin`](https://github.com/ainvaltin/nu-plugin).
It scans Go packages for annotated types and generates `nuconv.gen.go` files
with Nushell type definitions, parsers, and serializers.

## Usage

Mark a type for generation:

```go
// @nuconv:auto
type Config struct {
	Name    string `nu:"name"`
	Enabled bool
}
```

Run the generator with one or more Go package patterns:

```sh
go run . ./...
```

For each package with generated types, `nuconv` writes a `nuconv.gen.go` file.

Generated symbols follow this pattern:

- `NuDef<Type>`: Nushell type definition
- `NuParse<Type>`: convert a `nu.Value` into the Go type
- `NuSerialize<Type>`: convert the Go type into a `nu.Value`

## Annotations

- `@nuconv:auto` generates conversion code for the type.
- `@nuconv:manual` reserves the type so you can provide conversion code by hand.
- `nu:"name"` overrides the generated Nushell field name.
- `nu:"-"` skips a field.

Fields without a `nu` tag use snake_case names derived from the Go field name.

## Supported Types

`nuconv` supports common primitive types, `time.Time`, `time.Duration`, structs,
pointers, arrays, slices, maps, named types, and type aliases.

Channels, functions, tuples, interfaces, generics, and unions are not currently
supported.

## Development

Run tests with:

```sh
go test ./...
```
