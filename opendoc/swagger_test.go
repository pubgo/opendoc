package opendoc

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
	"k8s.io/kube-openapi/pkg/util"
)

func TestRefName(t *testing.T) {
	assert.Equal(t,
		"com.github.getkin.kin-openapi.openapi3.License",
		util.ToRESTFriendlyName(util.GetCanonicalTypeName(new(openapi3.License))))

	assert.Equal(t,
		"github.com/getkin/kin-openapi/openapi3.License",
		util.GetCanonicalTypeName(new(openapi3.License)))
}

type testQueryRsp struct {
	Name     string        `required:"true" json:"name" doc:"name of model" default:"test"`
	Token    string        `required:"true" json:"token" default:"test"`
	Optional *string       `json:"optional"`
	Req      *testQueryReq `json:"req" required:"true"`
}

type testQueryReq struct {
	Name     string        `required:"true" json:"name" doc:"name of model" default:"test"`
	Token    string        `header:"token" json:"token" default:"test"`
	Optional string        `query:"optional" json:"optional"`
	Rsp      *testQueryRsp `json:"rsp" required:"true"`
}

func TestGenSchema(t *testing.T) {
	ref, s := genSchema(testQueryReq{})
	assert.NotNil(t, s)
	assert.Equal(t, "#/components/schemas/com.github.pubgo.opendoc.testQueryReq", ref)

	data, err := json.Marshal(s)
	assert.NoError(t, err)
	assert.Equal(t,
		`{"properties":{"name":{"default":"test","description":"name of model","nullable":true,"type":"string"},"rsp":{"$ref":"#/components/schemas/com.github.pubgo.opendoc.testQueryRsp"}},"required":["name","rsp"],"type":"object"}`,
		string(data),
	)

	p := genParameters(testQueryReq{})
	data, err = json.Marshal(p)
	assert.NoError(t, err)
	assert.Equal(t,
		`[{"in":"header","name":"token","required":true,"schema":{"default":"test","type":"string"}},{"in":"query","name":"optional","required":true,"schema":{"type":"string"}}]`,
		string(data),
	)
}
