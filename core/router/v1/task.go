package v1

import (
	"github.com/gin-gonic/gin"
	"goportscan/core/management"
	"goportscan/core/models"
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
