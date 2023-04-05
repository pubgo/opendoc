package security

import (
	"github.com/getkin/kin-openapi/openapi3"
)

type OpenID struct {
	ConnectUrl string
}

func (i OpenID) Provider() string {
	return "OpenIdConnect"
}

func (i OpenID) Scheme() *openapi3.SecurityScheme {
	return &openapi3.SecurityScheme{
		Type:             "openIdConnect",
		OpenIdConnectUrl: i.ConnectUrl,
	}
}
