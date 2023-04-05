package templates

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"

	"github.com/gofiber/fiber/v2"
	"github.com/pubgo/funk/assert"
	"github.com/pubgo/opendoc/config"
)

//go:embed redoc.html
var reDocFile string

//go:embed swagger.html
var swaggerFile string

var reDocTemplate = assert.Exit1(template.New("").Parse(reDocFile))
var swaggerTemplate = assert.Exit1(template.New("").Parse(swaggerFile))

func RapiDocHandler(cfg *config.Config) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		ctx.Response().Header.Set("Content-Type", "text/html")
		ctx.Write([]byte(fmt.Sprintf(`<!doctype html>
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
</html>`, cfg.Title, cfg.OpenapiUrl)))
		return nil
	}
}

func ReDocHandler(cfg *config.Config) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		ctx.Response().Header.Set("Content-Type", "text/html")
		return reDocTemplate.Execute(ctx, map[string]string{
			"title":           cfg.Title,
			"openapi_url":     cfg.OpenapiUrl,
			"openapi_options": string(assert.Must1(json.Marshal(cfg.OpenapiOpt))),
		})
	}
}

func SwaggerHandler(cfg *config.Config) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		ctx.Response().Header.Set("Content-Type", "text/html")
		return swaggerTemplate.Execute(ctx, map[string]string{
			"title":           cfg.Title,
			"openapi_url":     cfg.OpenapiUrl,
			"openapi_options": string(assert.Must1(json.Marshal(cfg.OpenapiOpt))),
		})
	}
}
