package security

import (
	"github.com/getkin/kin-openapi/openapi3"
)

type Basic struct {
}

type User struct {
	Username string
	Password string
}

func (b Basic) Provider() string {
	return "Basic"
}

func (b Basic) Scheme() *openapi3.SecurityScheme {
	return &openapi3.SecurityScheme{
		Type:   "http",
		Scheme: "basic",
	}
}
