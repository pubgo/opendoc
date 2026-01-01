# OpenDoc - Go OpenAPI 文档生成工具

## 简介

`OpenDoc` 是一个用于 Go 项目的 OpenAPI 文档生成工具，它能够自动生成符合 OpenAPI 3.0 规范的 API 文档，并提供请求模型验证功能。该项目与具体 HTTP 框架无关，可以与任何 HTTP 框架（Gin、Echo、Fiber 等）一起使用，提供了一套完整的 API 文档解决方案。

## 特性

- 自动生成 OpenAPI 3.0 规范的 API 文档
- 支持多种安全认证方式（Basic、Bearer、ApiKey、OpenID、OAuth2）
- 请求参数自动验证和模型映射
- 支持服务分组和路由管理
- 提供 Swagger UI、Redoc 和 RapiDoc 文档展示
- 支持子应用挂载和独立文档
- 可在生产环境中禁用文档功能
- 与 HTTP 框架无关，可集成到任何 Go HTTP 项目

## 安装

```shell
go get -u github.com/pubgo/opendoc
```

## 快速开始

以下是一个简单的使用示例：

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/samber/lo"

	"github.com/pubgo/opendoc/opendoc"
	"github.com/pubgo/opendoc/security"
	"github.com/pubgo/opendoc/templates"
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
	Data any    `json:"data"`
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
			op.SetModel(new(TestQueryReq), new(TestQueryRsp))
			op.SetSummary("create article")
		})

		srv.GetOf(func(op *opendoc.Operation) {
			op.SetPath("/articles")
			op.SetOperation("article_list")
			op.SetModel(new(TestQueryReq), new(TestQueryRsp))
			op.SetSummary("get article list")
			op.AddResponse("Test", new(TestQueryReq))
		})

		srv.PutOf(func(op *opendoc.Operation) {
			op.SetPath("/articles/{id}")
			op.SetOperation("article_update")
			op.SetModel(new(TestQueryReq), new(TestQueryRsp))
			op.SetSummary("update article")
		})

		srv.DeleteOf(func(op *opendoc.Operation) {
			op.SetPath("/articles/{id}")
			op.SetOperation("article_delete")
			op.SetModel(new(TestQueryReq), new(TestQueryRsp))
			op.SetSummary("delete article")
		})
	})

	http.HandleFunc("/docs/", templates.SwaggerHandler("API Documentation", "/openapi.json"))
	http.HandleFunc("/redoc/", templates.ReDocHandler("API Documentation", "/openapi.json"))
	http.HandleFunc("/rapidoc/", templates.RApiDocHandler("/openapi.json"))
	http.HandleFunc("/openapi.json", func(w http.ResponseWriter, r *http.Request) {
		swagger := doc.BuildSwagger()
		w.Header().Set("Content-Type", "application/json")
		_ = swagger.MarshalJSON()
	})

	fmt.Println("Server starting at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
```

## 核心概念

### Swagger 配置

使用 `opendoc.New()` 创建 Swagger 实例，可以配置项目的基本信息：

```go
doc := opendoc.New(func(swag *opendoc.Swagger) {
    swag.Title = "API Title"
    swag.Description = "API Description"
    swag.Version = "1.0.0"
    swag.License = &opendoc.License{
        Name: "Apache License 2.0",
        URL:  "https://license-url"
    }
    swag.Contact = &opendoc.Contact{
        Name:  "Contact Name",
        URL:   "https://contact-url",
        Email: "contact@example.com"
    }
})
```

### 服务定义

使用 `ServiceOf` 方法定义服务组：

```go
doc.ServiceOf("service-name", func(srv *opendoc.Service) {
    srv.SetPrefix("/api/v1")  // 设置路由前缀
    srv.AddTags("User", "Account")  // 添加标签
    srv.AddSecurity(security.Basic{})  // 添加安全认证
})
```

### 操作定义

在服务中定义 API 操作：

```go
srv.GetOf(func(op *opendoc.Operation) {
    op.SetPath("/users")
    op.SetOperation("get_users")
    op.SetModel(requestStruct, responseStruct)
    op.SetSummary("获取用户列表")
    op.SetDescription("获取所有用户信息")
})
```

### 安全认证

OpenDoc 支持多种安全认证方式：

- `security.Basic{}` - HTTP Basic 认证
- `security.Bearer{}` - Bearer Token 认证
- `security.ApiKey{}` - API Key 认证
- `security.OpenID{}` - OpenID Connect 认证
- `security.OAuth2{}` - OAuth2 认证

```go
srv.AddSecurity(security.Basic{}, security.Bearer{})
```

## 参数映射

OpenDoc 支持多种参数映射方式：

- `path` - 路径参数
- `query` - 查询参数
- `header` - 请求头参数
- `cookie` - Cookie 参数
- `json` - JSON 请求体参数

```go
type RequestStruct struct {
    ID     int    `path:"id" validate:"required" json:"id"`
    Name   string `query:"name" validate:"required" json:"name"`
    Token  string `header:"authorization" json:"token"`
    Locale string `cookie:"locale" json:"locale"`
}
```

## 与其他 HTTP 框架集成

OpenDoc 与 HTTP 框架无关，可以轻松集成到任何 Go HTTP 框架中，如 Gin、Echo、Fiber 等。只需将 OpenDoc 生成的 OpenAPI 规范提供给相应的 HTTP 框架即可。

## 文档模板

OpenDoc 提供了多种文档展示模板：

- Swagger UI - 标准的 Swagger 文档界面
- Redoc - 现代化的 API 文档界面
- RapiDoc - 快速、轻量级的文档界面

## 构建和运行

使用以下命令运行示例：

```bash
make run-example
```

或者直接运行：

```bash
go run internal/examples/*.go
```

## 许可证

该项目根据 [Apache-2.0](https://github.com/pubgo/opendoc/blob/master/LICENSE) 许可证授权。