package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type TestQuery struct {
	Model *QueryModel
}

func (t *TestQuery) NewModel() interface{} {
	return &QueryModel{}
}

func (t *TestQuery) Handler(c *gin.Context) {
	c.JSON(http.StatusOK, t.Model)
}

type TestQueryPath struct {
	Model *QueryPathModel
}

func (t *TestQueryPath) NewModel() interface{} {
	return &QueryPathModel{}
}

func (t *TestQueryPath) Handler(c *gin.Context) {
	c.JSON(http.StatusOK, t.Model)
}

type TestForm struct {
	Model *FormModel
}

func (t *TestForm) NewModel() interface{} {
	return &FormModel{}
}

func (t *TestForm) Handler(c *gin.Context) {
	c.JSON(http.StatusOK, t.Model)
}

type TestBody struct {
	Model *BodyModel
}

func (t *TestBody) NewModel() interface{} {
	return &BodyModel{}
}

func (t *TestBody) Handler(c *gin.Context) {
	c.JSON(http.StatusOK, t.Model)
}