package v1

import (
	"dss/core/management"
	"dss/core/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

var Grpc *_Grpc

type _Grpc struct {
}

// Get
// @Summary Get Grpc Client
// @Tags Grpc
// @Accept  json
// @Produce  json
// @Param search query string false "Fuzzy Query"
// @Param page query string false "Current Page Default:1"
// @Param size query string false "Current Size Default:10"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /api/v1/grpc/client [get]
func (*_Grpc) Get(ctx *gin.Context) {
	var (
		g     = models.Gin{Ctx: ctx}
		param = models.ClientQueryFunc()
	)
	if err := ctx.ShouldBindQuery(&param); err != nil {
		g.Fail(http.StatusBadRequest, err)
		return
	}
	resp, err := management.GrpcManager.Get(*param)
	if err != nil {
		g.Fail(http.StatusBadRequest, err)
		return
	}
	g.Success(resp)
}
