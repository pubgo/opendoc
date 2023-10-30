package opendoc

import (
	"github.com/getkin/kin-openapi/openapi3"
)

type License = openapi3.License
type Contact = openapi3.Contact
type Servers = openapi3.Servers
type Server = openapi3.Server

// NamedEnum returns the enumerated acceptable values with according string names.
type NamedEnum interface {
	NamedEnum() ([]interface{}, []string)
}

// Enum returns the enumerated acceptable values.
type Enum interface {
	Enum() []interface{}
}

// OneOfExposer exposes "oneOf" items as list of samples.
type OneOfExposer interface {
	JSONSchemaOneOf() []interface{}
}

// AnyOfExposer exposes "anyOf" items as list of samples.
type AnyOfExposer interface {
	JSONSchemaAnyOf() []interface{}
}

// AllOfExposer exposes "allOf" items as list of samples.
type AllOfExposer interface {
	JSONSchemaAllOf() []interface{}
}

// NotExposer exposes "not" schema as a sample.
type NotExposer interface {
	JSONSchemaNot() interface{}
}

// IfExposer exposes "if" schema as a sample.
type IfExposer interface {
	JSONSchemaIf() interface{}
}

// ThenExposer exposes "then" schema as a sample.
type ThenExposer interface {
	JSONSchemaThen() interface{}
}

// ElseExposer exposes "else" schema as a sample.
type ElseExposer interface {
	JSONSchemaElse() interface{}
}
