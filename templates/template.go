package templates

import (
	_ "embed"
	"html/template"
	"net/http"

	"github.com/pubgo/funk/assert"
)

//go:embed redoc.html
var reDocFile string

//go:embed swagger.html
var swaggerFile string

var reDocTemplate = assert.Exit1(template.New("").Parse(reDocFile))
var swaggerTemplate = assert.Exit1(template.New("").Parse(swaggerFile))

func ReDocHandler(title, url string) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/html")
		assert.Must(reDocTemplate.Execute(writer, map[string]string{
			"title":           title,
			"openapi_url":     url,
			"openapi_options": `{}`,
		}))
	}
}

func SwaggerHandler(title, url string) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/html")
		assert.Must(swaggerTemplate.Execute(writer, map[string]string{
			"title":           title,
			"openapi_url":     url,
			"openapi_options": `{}`,
		}))
	}
}
