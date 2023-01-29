package v1

import (
	"dss/core/management"
	"dss/core/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

var Port *_Port

type _Port struct {
}

// Get
// @Summary Get Port Scan Result
// @Tags PortScan
// @Accept  json
// @Produce  json
// @Param search query string false "Fuzzy Query"
// @Param location query string false "Location Select"
// @Param date query string false "Date Select"
// @Param page query string false "Current Page Default:1"
// @Param size query string false "Current Size Default:10"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /api/v1/port [get]
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

func (*_Port) Remind(ctx *gin.Context) {
	//todo
}

func (*_Port) Stats(ctx *gin.Context) {
	//todo
}

func (*_Port) Clear(ctx *gin.Context) {
	var (
		g = models.Gin{Ctx: ctx}
	)
	go management.PortManager.Clear()
	g.Success(nil)
}

func (*_Port) Trend(ctx *gin.Context) {
	//todo
}

// Location
// @Summary GroupBy Location
// @Tags PortScan
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /api/v1/port/location [get]
func (*_Port) Location(ctx *gin.Context) {
	var (
		g = models.Gin{Ctx: ctx}
	)
	resp, err := management.PortManager.Location()
	if err != nil {
		g.Fail(http.StatusBadRequest, err)
		return
	}
	g.Success(resp)
}
