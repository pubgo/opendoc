package templates

import (
	_ "embed"
	"fmt"
	"html/template"

	"github.com/gofiber/fiber/v2"
	"github.com/pubgo/funk/assert"
)

//go:embed redoc.html
var reDocFile string

//go:embed swagger.html
var swaggerFile string

var reDocTemplate = assert.Exit1(template.New("").Parse(reDocFile))
var swaggerTemplate = assert.Exit1(template.New("").Parse(swaggerFile))

func RapiDocHandler(title, url string) fiber.Handler {
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
</html>`, title, url)))
		return nil
	}
}

func ReDocHandler(title, url string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		ctx.Response().Header.Set("Content-Type", "text/html")
		return reDocTemplate.Execute(ctx, map[string]string{
			"title":           title,
			"openapi_url":     url,
			"openapi_options": `{}`,
		})
	}
}

func SwaggerHandler(title, url string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		ctx.Response().Header.Set("Content-Type", "text/html")
		return swaggerTemplate.Execute(ctx, map[string]string{
			"title":           title,
			"openapi_url":     url,
			"openapi_options": `{}`,
		})
	}
}
