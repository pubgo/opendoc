package security

import (
	"github.com/getkin/kin-openapi/openapi3"
)

type Security interface {
	Provider() string
	Scheme() *openapi3.SecurityScheme
}
