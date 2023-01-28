package v1

import (
	"github.com/gin-gonic/gin"
	"goportscan/core/management"
	"goportscan/core/models"
	"net/http"
)

var Port *_Port

type _Port struct {
}

func (*_Port) Get(ctx *gin.Context) {
	var (
		g     = models.Gin{Ctx: ctx}
		param = models.ScanQueryFunc()
	)
	if err := ctx.ShouldBindQuery(&param); err != nil {
		g.Fail(http.StatusBadRequest, err)
		return
	}
	resp, err := management.PortManager.Get(*param)
	if err != nil {
		g.Fail(http.StatusBadRequest, err)
		return
	}
	g.Success(resp)
}
