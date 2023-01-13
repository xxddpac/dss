package v1

import (
	"github.com/gin-gonic/gin"
	"goportscan/core/management"
	"goportscan/core/models"
	"net/http"
)

var Rule *_Rule

type _Rule struct {
}

func (*_Rule) Post(ctx *gin.Context) {
	var (
		g    = models.Gin{Ctx: ctx}
		body models.Rule
	)
	if err := ctx.ShouldBindJSON(&body); err != nil {
		g.Fail(http.StatusBadRequest, err)
		return
	}
	if err := management.RuleManager.Post(body); err != nil {
		g.Fail(http.StatusBadRequest, err)
		return
	}
	g.Success(nil)
}
