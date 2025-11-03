package security

import (
	"github.com/getkin/kin-openapi/openapi3"
)

type Basic struct{}

type User struct {
	// Username basic auth username
	Username string

	// Password basic auth password
	Password string
}

func (b Basic) Provider() AuthProviderType { return AuthProviderTypeBasic }

func (b Basic) Scheme() *openapi3.SecurityScheme {
	return &openapi3.SecurityScheme{
		Type:   "http",
		Scheme: "basic",
	}
}
