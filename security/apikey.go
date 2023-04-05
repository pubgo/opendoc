package security

import (
	"github.com/getkin/kin-openapi/openapi3"
)

type ApiKey struct {
	Name string
}

func (k ApiKey) Provider() string {
	return "ApiKey"
}

func (k ApiKey) Scheme() *openapi3.SecurityScheme {
	return &openapi3.SecurityScheme{
		Type: "http",
		In:   "header",
		Name: k.Name,
	}
}
