package v1

import (
	"github.com/gin-gonic/gin"
	"goportscan/core/models"
)

var Rule *_Rule

type _Rule struct {
}

func (*_Rule) Post(ctx *gin.Context) {
	var (
		g = models.Gin{Ctx: ctx}
	)
	g.Success(nil)
}
