package v1

import (
	"dss/core/management"
	"dss/core/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

var ScanItem *_ScanItem

type _ScanItem struct {
}

func (*_ScanItem) Get(ctx *gin.Context) {
	var (
		g     = models.Gin{Ctx: ctx}
		param = models.ScanItemQueryFunc()
	)
	if err := ctx.ShouldBindQuery(&param); err != nil {
		g.Fail(http.StatusBadRequest, err)
		return
	}
	resp, err := management.ScanItemManager.Get(*param)
	if err != nil {
		g.Fail(http.StatusBadRequest, err)
		return
	}
	g.Success(resp)
}

func (*_ScanItem) Post(ctx *gin.Context) {
	var (
		g     = models.Gin{Ctx: ctx}
		param models.ScanItem
	)
	if err := ctx.ShouldBindJSON(&param); err != nil {
		g.Fail(http.StatusBadRequest, err)
		return
	}
	if err := management.ScanItemManager.Post(param); err != nil {
		g.Fail(http.StatusBadRequest, err)
		return
	}
	g.Success(nil)
}

func (*_ScanItem) Put(ctx *gin.Context) {
	var (
		g     = models.Gin{Ctx: ctx}
		param models.ScanItem
		query models.QueryID
	)
	if err := ctx.ShouldBindJSON(&param); err != nil {
		g.Fail(http.StatusBadRequest, err)
		return
	}
	if err := ctx.ShouldBindQuery(&query); err != nil {
		g.Fail(http.StatusBadRequest, err)
		return
	}
	if err := management.ScanItemManager.Put(param, query); err != nil {
		g.Fail(http.StatusBadRequest, err)
		return
	}
	g.Success(nil)
}

func (*_ScanItem) Delete(ctx *gin.Context) {
	var (
		g     = models.Gin{Ctx: ctx}
		query models.QueryID
	)
	if err := ctx.ShouldBindQuery(&query); err != nil {
		g.Fail(http.StatusBadRequest, err)
		return
	}
	if err := management.ScanItemManager.Delete(query); err != nil {
		g.Fail(http.StatusBadRequest, err)
		return
	}
	g.Success(nil)
}

func (*_ScanItem) Query(ctx *gin.Context) {
	var (
		g     = models.Gin{Ctx: ctx}
		query models.QueryID
	)
	if err := ctx.ShouldBindQuery(&query); err != nil {
		g.Fail(http.StatusBadRequest, err)
		return
	}
	result, err := management.ScanItemManager.Query(query)
	if err != nil {
		g.Fail(http.StatusBadRequest, err)
		return
	}
	g.Success(result)
}

func (*_ScanItem) Enum(ctx *gin.Context) {
	var (
		g = models.Gin{Ctx: ctx}
	)
	g.Success(management.ScanItemManager.Enum())
}
