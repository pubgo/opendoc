package security

import (
	"github.com/getkin/kin-openapi/openapi3"
)

type Bearer struct{}

func (b Bearer) Provider() AuthProviderType { return AuthProviderTypeBearer }

func (b Bearer) Scheme() *openapi3.SecurityScheme {
	return &openapi3.SecurityScheme{
		Type:         "http",
		Scheme:       "bearer",
		BearerFormat: "JWT",
	}
}
