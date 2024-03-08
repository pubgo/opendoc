package main

import (
	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/recovery"
	"github.com/pubgo/opendoc/opendoc"
	"github.com/pubgo/opendoc/security"
	"net/http"
)

type TestQueryReqAAA struct {
	ID       int     `path:"id" validate:"required" json:"id" description:"id of model" default:"1"`
	Name     string  `required:"true" json:"name" validate:"required" doc:"name of model" default:"test"`
	Name1    *string `required:"true" json:"name1" validate:"required" doc:"name1 of model" default:"test"`
	Token    string  `header:"token" json:"token" default:"test"`
	Optional string  `query:"optional" json:"optional"`
}

func main() {
	defer recovery.Exit()

	var doc = opendoc.New(func(swag *opendoc.Swagger) {
		swag.Config.Title = "this service web title "
		swag.Description = "this is description"
		swag.License = &opendoc.License{
			Name: "Apache License 2.0",
			URL:  "https://github.com/pubgo/opendoc/blob/master/LICENSE",
		}

		swag.Contact = &opendoc.Contact{
			Name:  "barry",
			URL:   "https://github.com/pubgo/opendoc",
			Email: "kooksee@163.com",
		}

		swag.TermsOfService = "https://github.com/pubgo"
	})

	doc.ServiceOf("test article service", func(srv *opendoc.Service) {
		srv.SetPrefix("/api/v1")
		srv.AddSecurity(security.Basic{}, security.Bearer{})
		srv.PostOf(func(op *opendoc.Operation) {
			op.SetPath("/articles")
			op.SetOperation("article_create")
			op.SetModel(new(TestQueryReq1), new(TestQueryRsp))
			op.SetSummary("create article")
		})

		srv.GetOf(func(op *opendoc.Operation) {
			op.SetPath("/articles")
			op.SetOperation("article_list")
			op.SetModel(new(TestQueryReq), new(TestQueryRsp))
			op.SetSummary("get article list")
			op.AddResponse("Test", new(TestQueryReqAAA))
		})

		srv.PutOf(func(op *opendoc.Operation) {
			op.SetPath("/articles/{id}")
			op.SetOperation("article_update")
			op.SetModel(new(TestQueryReq1), new(TestQueryRsp))
			op.SetSummary("update article")
			op.AddResponse("error", &TestFileReq{})
		})
	})

	//data := assert.Must1(doc.MarshalYAML())
	//assert.Exit(os.WriteFile("openapi.yaml", data, 0644))

	var app = http.NewServeMux()
	doc.InitRouter(app)
	assert.Exit(http.ListenAndServe("localhost:8080", app))
}
