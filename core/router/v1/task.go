package v1

import (
	"github.com/gin-gonic/gin"
	"goportscan/core/management"
	"goportscan/core/models"
)

var Task *_Task

type _Task struct {
}

func (*_Task) Post(ctx *gin.Context) {
	var (
		g = models.Gin{Ctx: ctx}
	)
	go management.TaskManager.Post()
	g.Success(nil)
}
