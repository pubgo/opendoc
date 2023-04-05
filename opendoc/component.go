package opendoc

import (
	"github.com/getkin/kin-openapi/openapi3"
)

var components = openapi3.NewComponents()

func init() {
	components.SecuritySchemes = make(openapi3.SecuritySchemes)
	components.Schemas = make(map[string]*openapi3.SchemaRef)
	components.Examples = make(map[string]*openapi3.ExampleRef)
	components.Responses = make(map[string]*openapi3.ResponseRef)
	components.Links = make(map[string]*openapi3.LinkRef)
	components.RequestBodies = make(map[string]*openapi3.RequestBodyRef)
	components.Headers = make(map[string]*openapi3.HeaderRef)
	components.Parameters = make(map[string]*openapi3.ParameterRef)
}
