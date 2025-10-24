package templates

import (
	_ "embed"
	"html/template"
	"net/http"

	"github.com/pubgo/opendoc/opendoc"
	"github.com/samber/lo"
)

//go:embed redoc.html
var reDocFile string

//go:embed swagger.html
var swaggerFile string

//go:embed rapidoc.html
var rApiDocFile string

var (
	reDocTemplate       = lo.Must(template.New("").Parse(reDocFile))
	swaggerTemplate     = lo.Must(template.New("").Parse(swaggerFile))
	rApiDocFileTemplate = lo.Must(template.New("").Parse(rApiDocFile))
)

func RApiDocHandler(url string) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/html")
		lo.Must0(rApiDocFileTemplate.Execute(writer, map[string]string{
			"openapi_url":     url,
			"openapi_options": `{}`,
		}))
	}
}

func ReDocHandler(title, url string) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/html")
		lo.Must0(reDocTemplate.Execute(writer, map[string]string{
			"title":           title,
			"openapi_url":     url,
			"openapi_options": `{}`,
		}))
	}
}

func SwaggerHandler(title, url string) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/html")
		lo.Must0(swaggerTemplate.Execute(writer, map[string]string{
			"title":           title,
			"openapi_url":     url,
			"openapi_options": `{}`,
		}))
	}
}

func InitRouter(r *http.ServeMux, s *opendoc.Swagger, cfg Config) {
	title := s.Title
	r.Handle(cfg.OpenapiRouter, SwaggerHandler(title, cfg.OpenapiUrl))
	r.Handle(cfg.OpenapiRedocRouter, ReDocHandler(title, cfg.OpenapiUrl))
	r.Handle(cfg.OpenapiRApiDocRouter, RApiDocHandler(cfg.OpenapiUrl))
	r.Handle(cfg.OpenapiUrl, openapiDataHandler(s))
}

func openapiDataHandler(s *opendoc.Swagger) http.HandlerFunc {
	bytes := lo.Must1(s.MarshalYAML())
	return func(writer http.ResponseWriter, request *http.Request) {
		lo.Must1(writer.Write(bytes))
	}
}
