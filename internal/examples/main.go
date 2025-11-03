package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/samber/lo"

	"github.com/pubgo/opendoc/opendoc"
	"github.com/pubgo/opendoc/security"
	"github.com/pubgo/opendoc/templates"
)

type TestQueryReqAAA struct {
	ID       int     `path:"id" validate:"required" json:"id" description:"id of model" default:"1"`
	Name     string  `required:"true" json:"name" validate:"required" doc:"name of model" default:"test"`
	Name1    *string `required:"true" json:"name1" validate:"required" doc:"name1 of model" default:"test"`
	Token    string  `header:"token" json:"token" default:"test"`
	Optional string  `query:"optional" json:"optional"`
}

func main() {
	doc := opendoc.New(func(swag *opendoc.Swagger) {
		swag.Title = "this service web title "
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

	data := lo.Must1(doc.MarshalYAML())
	lo.Must0(os.WriteFile("internal/examples/openapi.yaml", data, 0644))

	app := http.NewServeMux()
	templates.InitRouter(app, doc, templates.DefaultCfg())

	fmt.Println("http://localhost:8082/debug/apidocs")
	lo.Must0(http.ListenAndServe("localhost:8082", app))
}
