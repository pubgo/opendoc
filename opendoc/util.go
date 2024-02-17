package opendoc

import (
	"fmt"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/fatih/structtag"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/goccy/go-json"
	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/log"
	"github.com/pubgo/opendoc/security"
	"k8s.io/kube-openapi/pkg/util"
)

func getTag(tags *structtag.Tags, key string, fn func(tag *structtag.Tag)) {
	var tag, err = tags.Get(key)
	if err == nil && tag.Key != "" {
		fn(tag)
	}
}

func checkModelType(model interface{}) {
	var t reflect.Type
	if _t, ok := model.(reflect.Type); ok {
		t = _t
	} else {
		t = reflect.TypeOf(model)
	}

	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	assert.If(t.Kind() != reflect.Struct, "The native type of model should be struct")
}

func getSchemaName(val interface{}) string {
	return util.ToRESTFriendlyName(getCanonicalTypeName(val))
}

func getComponentName(name string) string {
	return fmt.Sprintf("#/components/schemas/%s", name)
}

func getCanonicalTypeName(val interface{}) string {
	var model reflect.Type
	if typ, ok := val.(reflect.Type); ok {
		model = typ
	} else {
		model = reflect.TypeOf(val)
	}

	for model.Kind() == reflect.Pointer {
		model = model.Elem()
	}

	if model.PkgPath() == "" {
		return model.Name()
	}

	path := model.PkgPath()
	if strings.Contains(path, "/vendor/") {
		path = path[strings.Index(path, "/vendor/")+len("/vendor/"):]
	}

	path = strings.Trim(strings.TrimPrefix(path, "vendor"), "/")
	return path + "." + model.Name()
}

func getSecurityRequirements(securities []security.Security) *openapi3.SecurityRequirements {
	securityRequirements := openapi3.NewSecurityRequirements()
	for _, s := range securities {
		securityRequirements.With(openapi3.NewSecurityRequirement().Authenticate(s.Provider()))
		components.SecuritySchemes[s.Provider()] = &openapi3.SecuritySchemeRef{Value: s.Scheme()}
	}
	return securityRequirements
}

func genSchema(val interface{}) (ref string, schema *openapi3.Schema) {
	var model reflect.Type
	if _t, ok := val.(reflect.Type); ok {
		model = _t
	} else {
		model = reflect.TypeOf(val)
	}

	for model.Kind() == reflect.Pointer {
		model = model.Elem()
	}

	assert.If(model.Kind() == reflect.Interface, "type:%s kind should not be interface", model)

	switch model {
	case reflect.TypeOf([]byte{}):
		return "", openapi3.NewBytesSchema()
	case reflect.TypeOf(multipart.FileHeader{}):
		return "", &openapi3.Schema{Type: openapi3.TypeString, Format: "binary"}
	case reflect.TypeOf([]*multipart.FileHeader{}):
		schema = openapi3.NewArraySchema()
		schema.Items = openapi3.NewSchemaRef("", &openapi3.Schema{Type: openapi3.TypeString, Format: "binary"})
		return "", schema
	case reflect.TypeOf(time.Time{}):
		return "", openapi3.NewDateTimeSchema()
	case reflect.TypeOf(time.Duration(0)):
		return "", &openapi3.Schema{Type: openapi3.TypeString, Format: "duration"}
	case reflect.TypeOf(net.IP{}):
		return "", &openapi3.Schema{Type: openapi3.TypeString, Format: "ipv4"}
	case reflect.TypeOf(url.URL{}):
		return "", &openapi3.Schema{Type: openapi3.TypeString, Format: "uri"}
	}

	switch v := val.(type) {
	case OneOfExposer:
		var refs []*openapi3.SchemaRef
		for _, s := range v.JSONSchemaOneOf() {
			ref, schema := genSchema(s)
			if ref != "" {
				refs = append(refs, openapi3.NewSchemaRef(ref, nil))
			} else {
				refs = append(refs, &openapi3.SchemaRef{Value: schema})
			}
		}

		return "", &openapi3.Schema{OneOf: refs}
	case AnyOfExposer:
		var refs []*openapi3.SchemaRef
		for _, s := range v.JSONSchemaAnyOf() {
			ref, schema := genSchema(s)
			if ref != "" {
				refs = append(refs, openapi3.NewSchemaRef(ref, nil))
			} else {
				refs = append(refs, &openapi3.SchemaRef{Value: schema})
			}
		}

		return "", &openapi3.Schema{AnyOf: refs}

	case Enum:
		return "", &openapi3.Schema{Enum: v.Enum()}
	}

	switch model.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Uint, reflect.Uint8, reflect.Uint16:
		schema = openapi3.NewIntegerSchema()
	case reflect.Int32, reflect.Uint32:
		schema = openapi3.NewInt32Schema()
	case reflect.Int64, reflect.Uint64:
		schema = openapi3.NewInt64Schema()
	case reflect.String:
		schema = openapi3.NewStringSchema()
	case reflect.Float32, reflect.Float64:
		schema = openapi3.NewFloat64Schema()
	case reflect.Bool:
		schema = openapi3.NewBoolSchema()
	case reflect.Array, reflect.Slice:
		schema = openapi3.NewArraySchema()
		schema.Items = openapi3.NewSchemaRef(genSchema(model.Elem()))
	case reflect.Map:
		schema = openapi3.NewObjectSchema()
		schema.Items = openapi3.NewSchemaRef(genSchema(model))
	case reflect.Struct:
		schemaName := getSchemaName(val)
		if ss := components.Schemas[schemaName]; ss != nil {
			return getComponentName(schemaName), ss.Value
		}

		schema = openapi3.NewObjectSchema()
		components.Schemas[schemaName] = openapi3.NewSchemaRef("", schema)

		for i := 0; i < model.NumField(); i++ {
			field := model.Field(i)
			tags := assert.Must1(structtag.Parse(string(field.Tag)))
			if isParameter(tags) {
				continue
			}

			tag, err := tags.Get(jsonTag)
			if err != nil || tag.Name == "-" {
				continue
			}

			if !tag.HasOption(omitempty) {
				schema.Required = append(schema.Required, tag.Name)
			}

			ref, fieldSchema := genSchema(field.Type)
			if ref != "" {
				schema.Properties[tag.Name] = openapi3.NewSchemaRef(ref, nil)
				continue
			}

			fieldSchema = fieldSchema.WithNullable()
			fieldSchema.AllowEmptyValue = true

			getTag(tags, nullable, func(_ *structtag.Tag) { fieldSchema.Nullable = true })
			getTag(tags, readOnly, func(_ *structtag.Tag) { fieldSchema.ReadOnly = true })
			getTag(tags, writeOnly, func(_ *structtag.Tag) { fieldSchema.WriteOnly = true })
			getTag(tags, required, func(_ *structtag.Tag) { fieldSchema.AllowEmptyValue = false })
			getTag(tags, doc, func(tag *structtag.Tag) { fieldSchema.Description = tag.Name })
			getTag(tags, description, func(tag *structtag.Tag) { fieldSchema.Description = tag.Name })
			getTag(tags, format, func(tag *structtag.Tag) { fieldSchema.Format = tag.Name })
			getTag(tags, deprecated, func(tag *structtag.Tag) { fieldSchema.Deprecated = true })
			getTag(tags, defaultName, func(tag *structtag.Tag) { fieldSchema.Default = tag.Name })
			getTag(tags, example, func(tag *structtag.Tag) {
				if err := json.Unmarshal([]byte(tag.Value()), &fieldSchema.Example); err != nil {
					log.Err(err).Str("tag-value", tag.Value()).Msg("failed to unmarshal example")
				}
			})
			getTag(tags, validate, func(tag *structtag.Tag) {
				desc := fieldSchema.Description
				if desc == "" {
					desc = tag.Name
				} else {
					desc = desc + " validate:" + tag.Name
				}
				fieldSchema.Description = desc
			})
			schema.Properties[tag.Name] = openapi3.NewSchemaRef("", fieldSchema)
		}
		return getComponentName(schemaName), schema
	}

	return "", schema
}

func genRequestBody(model interface{}, contentType ...string) *openapi3.RequestBodyRef {
	if len(contentType) == 0 {
		contentType = []string{"application/json"}
	}

	_, schema := genSchema(model)
	body := &openapi3.RequestBodyRef{Value: openapi3.NewRequestBody()}
	body.Value.Required = true
	body.Value.Content = openapi3.NewContentWithSchema(schema, contentType)
	return body
}

func genResponses(response interface{}, contentType ...string) *openapi3.Responses {
	if len(contentType) == 0 {
		contentType = []string{"application/json"}
	}

	_, schema := genSchema(response)
	content := openapi3.NewContentWithSchema(schema, contentType)
	var docText = http.StatusText(http.StatusOK)
	var rsp = &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: &docText,
			Content:     content,
		},
	}

	ret := openapi3.NewResponses()
	ret.Set("200", rsp)
	ret.Set("default", rsp)
	return ret
}

func isParameter(val *structtag.Tags) bool {
	var params = []string{queryTag, uriTag, pathTag, headerTag, cookieTag}
	for i := range params {
		if _, err := val.Get(params[i]); err == nil {
			return true
		}
	}
	return false
}

func genParameters(val interface{}) openapi3.Parameters {
	assert.If(val == nil, "val is nil")

	var model reflect.Type
	if _t, ok := val.(reflect.Type); ok {
		model = _t
	} else {
		model = reflect.TypeOf(val)
	}

	for model.Kind() == reflect.Pointer {
		model = model.Elem()
	}

	assert.If(model.Kind() != reflect.Struct, "type:%s kind should be struct", model)

	parameters := openapi3.NewParameters()
	for i := 0; i < model.NumField(); i++ {
		field := model.Field(i)
		tags := assert.Must1(structtag.Parse(string(field.Tag)))
		if !isParameter(tags) {
			continue
		}

		ref, schema := genSchema(field.Type)
		if ref != "" {
			parameters = append(parameters, &openapi3.ParameterRef{
				Value: &openapi3.Parameter{Schema: openapi3.NewSchemaRef(ref, schema)},
			})
			continue
		}

		getTag(tags, defaultName, func(tag *structtag.Tag) { schema.Default = tag.Name })

		parameter := new(openapi3.Parameter)
		getTag(tags, queryTag, func(tag *structtag.Tag) {
			parameter = openapi3.NewQueryParameter(tag.Name)
			if !tag.HasOption(omitempty) {
				parameter.Required = true
			}
		})

		getTag(tags, headerTag, func(tag *structtag.Tag) {
			parameter = openapi3.NewHeaderParameter(tag.Name)
			if !tag.HasOption(omitempty) {
				parameter.Required = true
			}
		})

		getTag(tags, cookieTag, func(tag *structtag.Tag) {
			parameter = openapi3.NewCookieParameter(tag.Name)
			if !tag.HasOption(omitempty) {
				parameter.Required = true
			}
		})

		getTag(tags, required, func(tag *structtag.Tag) { parameter.Required = true })
		getTag(tags, uriTag, func(tag *structtag.Tag) { parameter = openapi3.NewPathParameter(tag.Name) })
		getTag(tags, pathTag, func(tag *structtag.Tag) { parameter = openapi3.NewPathParameter(tag.Name) })

		if parameter.In == "" {
			continue
		}

		parameter.Schema = openapi3.NewSchemaRef(ref, schema)
		getTag(tags, doc, func(tag *structtag.Tag) { parameter.Description = tag.Name })
		getTag(tags, description, func(tag *structtag.Tag) { parameter.Description = tag.Name })
		getTag(tags, validate, func(tag *structtag.Tag) {
			desc := parameter.Description
			if desc == "" {
				desc = tag.Name
			} else {
				desc = desc + " validate:" + tag.Name
			}
			parameter.Description = desc
		})
		getTag(tags, deprecated, func(tag *structtag.Tag) { parameter.Deprecated = true })
		parameters = append(parameters, &openapi3.ParameterRef{Value: parameter})
	}
	return parameters
}

// unescape unescapes an extended JSON pointer
func unescape(s string) string {
	s = strings.ReplaceAll(s, "~2", "*")
	s = strings.ReplaceAll(s, "~1", "/")
	return strings.ReplaceAll(s, "~0", "~")
}

// espaceSJSONPath escapes a sjson path
func espaceSJSONPath(s string) string {
	// https://github.com/tidwall/sjson/blob/master/sjson.go#L47
	s = strings.ReplaceAll(s, "|", "\\|")
	s = strings.ReplaceAll(s, "#", "\\#")
	s = strings.ReplaceAll(s, "@", "\\@")
	s = strings.ReplaceAll(s, "*", "\\*")

	return strings.ReplaceAll(s, "?", "\\?")
}

// Specific JSON pointer encoding here
// ~0 => ~
// ~1 => /
// ... and vice versa

const (
	encRefTok0 = `~0`
	encRefTok1 = `~1`
	decRefTok0 = `~`
	decRefTok1 = `/`
)

// Unescape unescapes a json pointer reference token string to the original representation
func Unescape(token string) string {
	step1 := strings.ReplaceAll(token, encRefTok1, decRefTok1)
	step2 := strings.ReplaceAll(step1, encRefTok0, decRefTok0)
	return step2
}

// Escape escapes a pointer reference token string
func Escape(token string) string {
	step1 := strings.ReplaceAll(token, decRefTok0, encRefTok0)
	step2 := strings.ReplaceAll(step1, decRefTok1, encRefTok1)
	return step2
}
