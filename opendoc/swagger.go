package opendoc

import (
	"log"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/invopop/yaml"
)

type Swagger struct {
	rootPath       string
	Title          string
	Description    string
	Version        string
	TermsOfService string
	Routers        []*Service
	Servers        openapi3.Servers
	Contact        *openapi3.Contact
	License        *openapi3.License
}

func (s *Swagger) SetRootPath(path string) {
	if path == "" {
		log.Panic("path should not be null")
	}

	s.rootPath = "/" + strings.Trim(strings.TrimSpace(path), "/")
}

func (s *Swagger) ServiceOf(name string, cb func(srv *Service)) {
	srv := newService(name)
	srv.prefix = s.rootPath
	s.Routers = append(s.Routers, srv)
	cb(srv)
}

func (s *Swagger) WithService() *Service {
	srv := new(Service)
	srv.prefix = s.rootPath
	s.Routers = append(s.Routers, srv)
	return srv
}

func (s *Swagger) buildSwagger() *openapi3.T {
	t := &openapi3.T{
		OpenAPI:    "3.0.0",
		Servers:    s.Servers,
		Components: &components,
		Info: &openapi3.Info{
			Title:          s.Title,
			Description:    s.Description,
			TermsOfService: s.TermsOfService,
			Contact:        s.Contact,
			License:        s.License,
			Version:        s.Version,
		},
	}

	var opts []openapi3.NewPathsOption
	for i := range s.Routers {
		for k, v := range s.Routers[i].Openapi() {
			if v == nil {
				continue
			}

			opts = append(opts, openapi3.WithPath(k, v))
		}
	}
	t.Paths = openapi3.NewPaths(opts...)

	return t
}

func (s *Swagger) BuildSwagger() *openapi3.T {
	return s.buildSwagger()
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
	swagger := &Swagger{}
	for i := range handles {
		handles[i](swagger)
	}
	return swagger
}
