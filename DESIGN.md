# Multi Packages

- []package:
	- []types:
		- gen: serialize/parse
		- deps: []types

Type resolution:

- all type methods are generated to pkg they belong to
- types that depend on other pkgs should autoimport

# Single Package

*Expected result:*

- type definitions
- parsers
- serializers

*Runtime constructs to use:*

- nu-plugin/types.*(...) -> type definitions
- runtime values (nu.Record, string, int, etc...)

*Keep in mind:*

- error handling (show parse error path)
- type resolution

# Non-exposed types

Some types rely on other types that are not exposed with a
`@nutype` marker, these other types can have type def / parser /
serializer but must be named according to non-exposed syntax.

This is some format like: `__private_<typical name>` to prevent it
from showing up in LSP and exports.

# Primitive types

Primitive types do not require extra parsers/serializers as their
handling is built-in with `nu.ToValue` or `val, ok :=
nu.Value().(<type>)` and built-in `types.*()` functions.

As such, the type generator should not bother to "resolve" them,
and simply short-circuit the logic.

# Anonymous types

Anonymous types can be given a name so that their usage is less
verbose in generated code.

This name should be calculated by their access path from the
parent named type (which can be a named anonymous type).

The parser should store these under a "synthetic anonymous types"
registry. They should be treated like a [[#non-exposed types]].

