package security

import (
	"github.com/getkin/kin-openapi/openapi3"
)

type ApiKey struct {
	// Name of the header to be used.
	Name string
}

func (k ApiKey) Provider() AuthProviderType { return AuthProviderTypeApiKey }

func (k ApiKey) Scheme() *openapi3.SecurityScheme {
	return &openapi3.SecurityScheme{
		Type: "http",
		In:   "header",
		Name: k.Name,
	}
}
