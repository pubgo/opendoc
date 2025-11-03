package security

import (
	"github.com/getkin/kin-openapi/openapi3"
)

type Security interface {
	Provider() AuthProviderType
	Scheme() *openapi3.SecurityScheme
}

type AuthProviderType string

const (
	AuthProviderTypeBasic  AuthProviderType = "Basic"
	AuthProviderTypeApiKey AuthProviderType = "ApiKey"
	AuthProviderTypeBearer AuthProviderType = "Bearer"
	AuthProviderTypeOAuth2 AuthProviderType = "OAuth2"
	AuthProviderTypeOIDC   AuthProviderType = "OIDC"
)
