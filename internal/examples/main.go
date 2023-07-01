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
			op.SetPath("article_create", "/articles")
			op.SetModel(new(TestQueryReq1), new(TestQueryRsp))
			op.SetDescription("create article")
		})

		srv.GetOf(func(op *opendoc.Operation) {
			op.SetPath("article_list", "/articles")
			op.SetModel(new(TestQueryReq), new(TestQueryRsp))
			op.SetDescription("get article list")
		})

		srv.PutOf(func(op *opendoc.Operation) {
			op.SetPath("article_update", "/articles/{id}")
			op.SetModel(new(TestQueryReq1), new(TestQueryRsp))
			op.SetDescription("update article")
		})
	})

	var app = fiber.New()
	doc.InitRouter(app)
	assert.Exit(app.Listen("localhost:8080"))
}
