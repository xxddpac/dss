package v1

import (
	"dss/core/management"
	"dss/core/models"
	"github.com/gin-gonic/gin"
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
// @Router /api/v1/task [post]
func (*_Task) Post(ctx *gin.Context) {
	var (
		g = models.Gin{Ctx: ctx}
	)
	go management.TaskManager.Post()
	g.Success(nil)
}
