package main

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/pubgo/opendoc/swagger"
)

func NewSwagger() *swagger.Swagger {
	return swagger.New("SwaGin", "Swagger + Gin = SwaGin", "0.1.0",
		swagger.License(&openapi3.License{
			Name: "Apache License 2.0",
			URL:  "https://github.com/pubgo/opendoc/blob/dev/LICENSE",
		}),
		swagger.Contact(&openapi3.Contact{
			Name:  "long2ice",
			URL:   "https://github.com/pubgo/opendoc",
			Email: "long2ice@gmail.com",
		}),
		swagger.TermsOfService("https://github.com/long2ice"),
	)
}
