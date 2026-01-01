package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/invopop/yaml"
	"github.com/pubgo/opendoc/opendoc"
	"github.com/pubgo/opendoc/security"
	"github.com/pubgo/opendoc/templates"
	"github.com/samber/lo"
)

type TestQueryReq struct {
	ID       int     `path:"id" validate:"required" json:"id" description:"id of model" default:"1"`
	Name     string  `required:"true" json:"name" validate:"required" doc:"name of model" default:"test"`
	Name1    *string `required:"true" json:"name1" validate:"required" doc:"name1 of model" default:"test"`
	Token    string  `header:"token" json:"token" default:"test"`
	Optional string  `query:"optional" json:"optional"`
}

type TestQueryRsp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data map[string]interface{} `json:"data"`
}

type TestQueryReq1 struct {
	Name     string  `required:"true" json:"name" validate:"required" doc:"name of model" default:"test"`
	Name1    *string `required:"true" json:"name1" validate:"required" doc:"name1 of model" default:"test"`
	Token    string  `header:"token" json:"token" default:"test"`
	Optional string  `query:"optional" json:"optional"`
}

func main() {
	doc := opendoc.New(func(swag *opendoc.Swagger) {
		swag.Title = "this service web title "
		swag.Description = "this is description"
		swag.Version = "1.0.0"
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
			op.SetDescription("Creates a new article with the provided data")
		})

		srv.GetOf(func(op *opendoc.Operation) {
			op.SetPath("/articles")
			op.SetOperation("article_list")
			op.SetModel(new(TestQueryReq1), new(TestQueryRsp))
			op.SetSummary("get article list")
			op.SetDescription("Retrieves a list of articles")
		})

		srv.PutOf(func(op *opendoc.Operation) {
			op.SetPath("/articles/{id}")
			op.SetOperation("article_update")
			op.SetModel(new(TestQueryReq), new(TestQueryRsp))
			op.SetSummary("update article")
			op.SetDescription("Updates an existing article by ID")
		})

		srv.DeleteOf(func(op *opendoc.Operation) {
			op.SetPath("/articles/{id}")
			op.SetOperation("article_delete")
			op.SetModel(new(TestQueryReq), new(TestQueryRsp))
			op.SetSummary("delete article")
			op.SetDescription("Deletes an article by ID")
		})
	})

	lo.Must0(os.WriteFile("./internal/examples/openapi.yaml",
		lo.Must(yaml.Marshal(lo.Must(doc.BuildSwagger().MarshalYAML()))), 0644))

	http.HandleFunc("/docs/", templates.SwaggerHandler("API Documentation", "/openapi.json"))
	http.HandleFunc("/redoc/", templates.ReDocHandler("API Documentation", "/openapi.json"))
	http.HandleFunc("/rapidoc/", templates.RApiDocHandler("/openapi.json"))
	http.HandleFunc("/openapi.json", func(w http.ResponseWriter, r *http.Request) {
		swagger := doc.BuildSwagger()
		w.Header().Set("Content-Type", "application/json")
		data, _ := swagger.MarshalJSON()
		w.Write(data)
	})

	fmt.Println("Server starting at :8080")
	fmt.Println("Visit http://localhost:8080/docs/ for API documentation")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
