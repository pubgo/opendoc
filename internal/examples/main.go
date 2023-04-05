package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/recovery"
	"github.com/pubgo/opendoc/opendoc"
	"github.com/pubgo/opendoc/security"
)

func main() {
	defer recovery.Exit()

	var doc = opendoc.New(func(swag *opendoc.Swagger) {
		swag.Config.Title = "title-service"
		swag.Description = "this is description"
		swag.License = &opendoc.License{
			Name: "Apache License 2.0",
			URL:  "https://github.com/pubgo/opendoc/blob/dev/LICENSE",
		}

		swag.Contact = &opendoc.Contact{
			Name:  "long2ice",
			URL:   "https://github.com/pubgo/opendoc",
			Email: "long2ice@gmail.com",
		}

		swag.TermsOfService = "https://github.com/long2ice"
	})

	doc.ServiceOf("test", func(srv *opendoc.Service) {
		srv.SetPrefix("/api/v1")
		srv.AddSecurity(security.Basic{}, security.Bearer{})
		srv.PostOf(func(op *opendoc.Operation) {
			op.SetPath("no_model_opt", "/no_model")
			op.SetModel(new(TestNoModelReq), new(TestNoModelReq))
			op.SetSummary("Test no model")
			op.SetDescription("Test no model")
		})

		srv.GetOf(func(op *opendoc.Operation) {
			op.SetPath("article_list", "/v1/articles")
			op.SetModel(new(TestQueryReq), new(TestQueryRsp))
			op.SetSummary("get article list")
			op.SetDescription("Test query list model")
		})

		srv.PutOf(func(op *opendoc.Operation) {
			op.SetPath("article_update", "/v1/articles/{id}")
			op.SetModel(new(TestQueryReq1), new(TestQueryRsp))
			op.SetSummary("delete article")
			op.SetDescription("Test query list model")
		})
	})

	var app = fiber.New()
	doc.InitRouter(app)
	assert.Exit(app.Listen("localhost:8080"))
}
