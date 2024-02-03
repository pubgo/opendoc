package templates

import (
	_ "embed"
	"fmt"
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

func RapiDocHandler(title, url string) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/html")
		writer.Write([]byte(fmt.Sprintf(`<!doctype html>
<html>
<head>
	<title>%s</title>
  <meta charset="utf-8">
  <script type="module" src="https://unpkg.com/rapidoc@9.1.4/dist/rapidoc-min.js"></script>
</head>
<body>
  <rapi-doc
		spec-url="%s"
		render-style="read"
    show-header="false"
    primary-color="#f74799"
    nav-accent-color="#47afe8"
  > </rapi-doc>
</body>
</html>`, title, url)))
	}
}

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
