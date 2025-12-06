package opendoc

import (
	"github.com/getkin/kin-openapi/openapi3"
)

type (
	License = openapi3.License
	Contact = openapi3.Contact
	Servers = openapi3.Servers
	Server  = openapi3.Server
)

// NamedEnum returns the enumerated acceptable values with according string names.
type NamedEnum interface {
	NamedEnum() ([]any, []string)
}

// Enum returns the enumerated acceptable values.
type Enum interface {
	Enum() []any
}

// OneOfExposer exposes "oneOf" items as list of samples.
type OneOfExposer interface {
	JSONSchemaOneOf() []any
}

// AnyOfExposer exposes "anyOf" items as list of samples.
type AnyOfExposer interface {
	JSONSchemaAnyOf() []any
}

// AllOfExposer exposes "allOf" items as list of samples.
type AllOfExposer interface {
	JSONSchemaAllOf() []any
}

// NotExposer exposes "not" schema as a sample.
type NotExposer interface {
	JSONSchemaNot() any
}

// IfExposer exposes "if" schema as a sample.
type IfExposer interface {
	JSONSchemaIf() any
}

// ThenExposer exposes "then" schema as a sample.
type ThenExposer interface {
	JSONSchemaThen() any
}

// ElseExposer exposes "else" schema as a sample.
type ElseExposer interface {
	JSONSchemaElse() any
}
