package main

import (
	"mime/multipart"
)

type TestQueryReq struct {
	Name     string `query:"name" validate:"required" json:"name" description:"name of model" default:"test"`
	Token    string `header:"token" validate:"required" json:"token" default:"test"`
	Optional string `query:"optional" json:"optional"`
	Name1    string `required:"true" json:"name1" validate:"required" doc:"name of model" default:"test"`
}

type TestQueryListReq struct {
	Name  string `query:"name" validate:"required" json:"name" description:"name of model" default:"test"`
	Token string `header:"token" validate:"required" json:"token" default:"test"`
}

type TestQueryPathReq struct {
	Name  string `query:"name" validate:"required" json:"name" description:"name of model" default:"test"`
	ID    int    `uri:"id" validate:"required" json:"id" description:"id of model" default:"1"`
	Token string `header:"token" validate:"required" json:"token" default:"test"`
}

type TestFormReq struct {
	ID   int    `query:"id" validate:"required" json:"id" description:"id of model" default:"1"`
	Name string `form:"name" validate:"required" json:"name" description:"name of model" default:"test"`
	List []int  `form:"list" validate:"required" json:"list" description:"list of model"`
}

type TestNoModelReq struct {
	Authorization string `header:"authorization" validate:"required" json:"authorization" default:"authorization"`
	Token         string `header:"token" binding:"required" json:"token" default:"token"`
}

type TestFileReq struct {
	File *multipart.FileHeader `form:"file" validate:"required" description:"file upload"`
}

type TestQueryRsp struct {
	Name     string        `required:"true" json:"name" doc:"name of model" default:"test"`
	Token    string        `required:"true" json:"token" default:"test"`
	Optional *string       `json:"optional"`
	Types    []string      `json:"types,omitempty" doc:"类型" required:"true" example:"[\"a\",\"b\"]" readOnly:"true"`
	Req      *TestQueryReq `json:"req" required:"true"`
}

type TestQueryReq1 struct {
	ID       int           `path:"id" validate:"required" json:"id" description:"id of model" default:"1"`
	Name     string        `required:"true" json:"name" validate:"required" doc:"name of model" default:"test"`
	Name1    *string       `required:"true" json:"name1" validate:"required" doc:"name1 of model" default:"test"`
	Token    string        `header:"token" json:"token" default:"test"`
	Optional string        `query:"optional" json:"optional"`
	Rsp      *TestQueryRsp `json:"rsp" required:"true"`
}
