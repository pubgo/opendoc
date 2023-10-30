package opendoc

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"

	"github.com/pubgo/opendoc/security"
)

func newService(name string, tags ...string) *Service {
	return &Service{name: name, tags: tags}
}

type Service struct {
	name        string
	tags        []string
	operations  []*Operation
	prefix      string
	securities  []security.Security
	contentType []string
}

func (s *Service) SetName(name string) *Service {
	s.name = name
	return s
}

func (s *Service) AddContentType(contentType ...string) *Service {
	s.contentType = append(s.contentType, contentType...)
	return s
}

func (s *Service) AddSecurity(security ...security.Security) *Service {
	s.securities = append(s.securities, security...)
	return s
}

func (s *Service) SetPrefix(prefix string) *Service {
	prefix = strings.TrimSpace(prefix)
	prefix = strings.Trim(prefix, "/")
	prefix = "/" + prefix

	if s.prefix == "" {
		s.prefix = prefix
	} else {
		s.prefix = filepath.Join(s.prefix, prefix)
	}

	return s
}

func (s *Service) newOpt(cb func(op *Operation)) *Operation {
	op := new(Operation)
	op.prefix = s.prefix
	op.tags = append(op.tags, s.name)
	op.tags = append(op.tags, s.tags...)
	op.securities = append(op.securities, s.securities...)
	op.requestContentType = append(op.requestContentType, s.contentType...)
	op.responseContentType = append(op.responseContentType, s.contentType...)
	s.operations = append(s.operations, op)

	cb(op)
	return op
}

func (s *Service) PutOf(cb func(op *Operation)) *Service {
	op := s.newOpt(cb)
	op.method = http.MethodPut
	return s
}

func (s *Service) GetOf(cb func(op *Operation)) *Service {
	op := s.newOpt(cb)
	op.method = http.MethodGet
	return s
}

func (s *Service) PatchOf(cb func(op *Operation)) *Service {
	op := s.newOpt(cb)
	op.method = http.MethodPatch
	return s
}

func (s *Service) DeleteOf(cb func(op *Operation)) *Service {
	op := s.newOpt(cb)
	op.method = http.MethodDelete
	return s
}

func (s *Service) PostOf(cb func(op *Operation)) *Service {
	op := s.newOpt(cb)
	op.method = http.MethodPost
	return s
}

func (s *Service) Openapi() map[string]*openapi3.PathItem {
	var routes = make(map[string]*openapi3.PathItem)
	for i := range s.operations {
		op := s.operations[i]
		routes[op.method+op.path] = op.Openapi()
	}
	return routes
}
