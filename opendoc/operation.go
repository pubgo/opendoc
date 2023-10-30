package opendoc

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/pubgo/funk/assert"

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
	request             interface{}
	response            interface{}
	responses           map[string]*openapi3.ResponseRef
}

func (op *Operation) AddSecurity(security ...security.Security) *Operation {
	op.securities = append(op.securities, security...)
	return op
}

func (op *Operation) SetExclude(exclude bool) *Operation {
	op.exclude = exclude
	return op
}

func (op *Operation) AddResponse(name string, resp interface{}) *Operation {
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

func (op *Operation) SetPath(operationID string, path string) *Operation {
	assert.If(operationID == "", "operationID should not be nil")
	assert.If(path == "", "path should not be null")

	op.operationID = operationID

	path = strings.TrimSpace(path)
	path = strings.Trim(path, "/")
	op.path = filepath.Join(op.prefix, path)
	return op
}

func (op *Operation) SetOperation(operationID string) *Operation {
	assert.If(operationID == "", "operationID should not be nil")
	op.operationID = operationID
	return op
}

func (op *Operation) SetModel(req interface{}, rsp interface{}) *Operation {
	checkModelType(req)
	op.request = req

	checkModelType(rsp)
	op.response = rsp

	return op
}

func (op *Operation) Openapi() *openapi3.PathItem {
	if op.exclude {
		return nil
	}

	responses := genResponses(op.response, op.responseContentType...)
	if op.responses != nil {
		for k, v := range op.responses {
			responses[k] = v
		}
	}

	operation := &openapi3.Operation{
		Tags:        op.tags,
		OperationID: op.operationID,
		Summary:     op.summary,
		Description: op.description,
		Deprecated:  op.deprecated,
		Responses:   responses,
		Parameters:  genParameters(op.request),
		Security:    getSecurityRequirements(op.securities),
	}

	item := new(openapi3.PathItem)
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

	requestBody := genRequestBody(op.request, op.requestContentType...)
	switch op.method {
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		operation.RequestBody = requestBody
	}

	return item
}
