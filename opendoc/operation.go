package opendoc

import (
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/samber/lo"

	"github.com/pubgo/opendoc/security"
)

type Operation struct {
	prefix              string
	path                string
	method              string
	summary             string
	description         string
	deprecated          bool
	requestContentType  []string
	responseContentType []string
	tags                []string
	operationID         string
	exclude             bool
	securities          []security.Security
	request             any
	response            any
	responses           map[string]*openapi3.ResponseRef
}

func (op *Operation) AddSecurity(security ...security.Security) *Operation {
	op.securities = append(op.securities, security...)
	return op
}

func (op *Operation) AddTags(tags ...string) *Operation {
	op.tags = append(op.tags, tags...)
	return op
}

func (op *Operation) SetExclude(exclude bool) *Operation {
	op.exclude = exclude
	return op
}

func (op *Operation) AddResponse(name string, resp any) *Operation {
	if op.responses == nil {
		op.responses = make(map[string]*openapi3.ResponseRef)
	}

	_, schema := genSchema(resp)
	op.responses[name] = &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: &name,
			Content:     openapi3.NewContentWithSchema(schema, []string{"application/json"}),
		},
	}
	return op
}

func (op *Operation) SetDescription(description string) *Operation {
	if description == "" {
		return op
	}

	op.summary = description
	return op
}

func (op *Operation) SetSummary(summary string) *Operation {
	if summary == "" {
		return op
	}

	op.summary = summary
	return op
}

func (op *Operation) SetPath(path string) *Operation {
	if path == "" {
		log.Panic("path should not be null")
	}

	path = strings.TrimSpace(path)
	path = strings.Trim(path, "/")
	op.path = filepath.Join(op.prefix, path)
	return op
}

func (op *Operation) SetOperation(operationID string) *Operation {
	if operationID == "" {
		log.Panic("operationID should not be nil")
	}

	op.operationID = operationID
	return op
}

func (op *Operation) SetModel(req, rsp any) *Operation {
	checkModelType(req)
	op.request = req

	checkModelType(rsp)
	op.response = rsp

	return op
}

func (op *Operation) Openapi(item *openapi3.PathItem) {
	if op.exclude {
		return
	}

	responses := genResponses(op.response, op.responseContentType...)
	if op.responses != nil {
		for k, v := range op.responses {
			responses.Set(k, v)
		}
	}

	operation := &openapi3.Operation{
		Tags:        lo.Uniq(op.tags),
		OperationID: op.operationID,
		Summary:     op.summary,
		Description: op.description,
		Deprecated:  op.deprecated,
		Responses:   responses,
		Parameters:  genParameters(op.request),
		Security:    getSecurityRequirements(op.securities),
	}

	switch op.method {
	case http.MethodGet:
		item.Get = operation
	case http.MethodPost:
		item.Post = operation
	case http.MethodDelete:
		item.Delete = operation
	case http.MethodPut:
		item.Put = operation
	case http.MethodPatch:
		item.Patch = operation
	case http.MethodHead:
		item.Head = operation
	case http.MethodOptions:
		item.Options = operation
	case http.MethodConnect:
		item.Connect = operation
	case http.MethodTrace:
		item.Trace = operation
	}

	switch op.method {
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		operation.RequestBody = genRequestBody(op.request, op.requestContentType...)
	}
}
