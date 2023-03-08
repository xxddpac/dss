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
// @Summary 获取gRPC客户端信息(IPV4、系统类型、系统版本、主机名、是否在线)
// @Tags Grpc
// @Accept  json
// @Produce  json
// @Param search query string false "模糊查询"
// @Param page query string false "当前页数,默认值:1"
// @Param size query string false "当前条数,默认值:10"
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
