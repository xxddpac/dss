package v1

import (
	"dss/core/management"
	"dss/core/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

var Task *_Task

type _Task struct {
}

// Post
// @Summary Task Execute
// @Tags Task
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /api/v1/task [post]
func (*_Task) Post(ctx *gin.Context) {
	var (
		g     = models.Gin{Ctx: ctx}
		query models.QueryID
	)
	if err := ctx.ShouldBindQuery(&query); err != nil {
		g.Fail(http.StatusBadRequest, err)
		return
	}
	if err := management.TaskManager.Post(query); err != nil {
		g.Fail(http.StatusBadRequest, err)
		return
	}
	g.Success(nil)
}

func (*_Task) Get(ctx *gin.Context) {
	var (
		g     = models.Gin{Ctx: ctx}
		param = models.TaskQueryFunc()
	)
	if err := ctx.ShouldBindQuery(&param); err != nil {
		g.Fail(http.StatusBadRequest, err)
		return
	}
	resp, err := management.TaskManager.Get(*param)
	if err != nil {
		g.Fail(http.StatusBadRequest, err)
		return
	}
	g.Success(resp)
}
