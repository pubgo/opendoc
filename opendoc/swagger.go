package opendoc

import (
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
	"github.com/invopop/yaml"
	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/version"

	"github.com/pubgo/opendoc/templates"
)

type Swagger struct {
	rootPath string

	Config         *Config
	Description    string
	Version        string
	TermsOfService string
	Routers        []*Service
	Servers        openapi3.Servers
	Contact        *openapi3.Contact
	License        *openapi3.License
}

func (s *Swagger) SetRootPath(path string) {
	assert.If(path == "", "path should not be null")
	s.rootPath = "/" + strings.Trim(strings.TrimSpace(path), "/")
}

func (s *Swagger) ServiceOf(name string, cb func(srv *Service)) {
	var srv = newService(name)
	srv.prefix = s.rootPath
	s.Routers = append(s.Routers, srv)
	cb(srv)
}

func (s *Swagger) Service() *Service {
	var srv = new(Service)
	srv.prefix = s.rootPath
	s.Routers = append(s.Routers, srv)
	return srv
}

func (s *Swagger) buildSwagger() *openapi3.T {
	if s.Config == nil {
		s.Config = DefaultCfg()
	}

	var t = &openapi3.T{
		OpenAPI:    "3.0.0",
		Servers:    s.Servers,
		Components: &components,
		Info: &openapi3.Info{
			Title:          s.Config.Title,
			Description:    s.Description,
			TermsOfService: s.TermsOfService,
			Contact:        s.Contact,
			License:        s.License,
			Version:        s.Version,
		},
	}

	var routes = make(map[string]*openapi3.PathItem)
	for i := range s.Routers {
		for k, v := range s.Routers[i].Openapi() {
			if v == nil {
				continue
			}
			routes[k] = v
		}
	}
	t.Paths = routes

	return t
}

func (s *Swagger) InitRouter(r fiber.Router) {
	r.Get(s.Config.OpenapiRouter, templates.SwaggerHandler(s.Config.Title, s.Config.OpenapiUrl))
	r.Get(s.Config.OpenapiRedocRouter, templates.ReDocHandler(s.Config.Title, s.Config.OpenapiUrl))
	r.Get(s.Config.OpenapiUrl, s.OpenapiDataHandler())
}

func (s *Swagger) OpenapiDataHandler() fiber.Handler {
	var bytes = assert.Must1(s.MarshalYAML())
	return func(ctx *fiber.Ctx) error {
		var _, err = ctx.Write(bytes)
		return err
	}
}

func (s *Swagger) MarshalJSON() ([]byte, error) {
	return s.buildSwagger().MarshalJSON()
}

func (s *Swagger) MarshalYAML() ([]byte, error) {
	b, err := s.MarshalJSON()
	if err != nil {
		return nil, err
	}

	return yaml.JSONToYAML(b)
}

func New(handles ...func(swag *Swagger)) *Swagger {
	swagger := &Swagger{
		Config:      DefaultCfg(),
		Description: fmt.Sprintf("project:%s version:%s commit:%s", version.Project(), version.Version(), version.CommitID()),
		Version:     version.Version(),
	}

	for i := range handles {
		handles[i](swagger)
	}

	return swagger
}
