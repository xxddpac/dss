package v1

import (
	"dss/core/management"
	"dss/core/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

var Scan *_Scan

type _Scan struct {
}

// Get
// @Summary Get Scan Result
// @Tags Scan
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
func (*_Scan) Get(ctx *gin.Context) {
	var (
		g     = models.Gin{Ctx: ctx}
		param = models.ScanQueryFunc()
	)
	if err := ctx.ShouldBindQuery(&param); err != nil {
		g.Fail(http.StatusBadRequest, err)
		return
	}
	resp, err := management.ScanManager.Get(*param)
	if err != nil {
		g.Fail(http.StatusBadRequest, err)
		return
	}
	g.Success(resp)
}

// Remind
// @Summary Compare yesterday with today,if new port open in today will notify by workChat
// @Tags Scan
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /api/v1/port/remind [get]
func (*_Scan) Remind(ctx *gin.Context) {
	var (
		g = models.Gin{Ctx: ctx}
	)
	go management.ScanManager.Remind()
	g.Success(nil)
}

// Clear
// @Summary Clear data more than 7 days
// @Tags Scan
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /api/v1/port/clear [delete]
func (*_Scan) Clear(ctx *gin.Context) {
	var (
		g = models.Gin{Ctx: ctx}
	)
	go management.ScanManager.Clear()
	g.Success(nil)
}

// Trend
// @Summary last 7 days scan trend
// @Tags Scan
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /api/v1/port/trend [get]
func (*_Scan) Trend(ctx *gin.Context) {
	var (
		g = models.Gin{Ctx: ctx}
	)
	resp, err := management.ScanManager.Trend()
	if err != nil {
		g.Fail(http.StatusBadRequest, err)
		return
	}
	g.Success(resp)
}

// Location
// @Summary GroupBy Location
// @Tags Scan
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /api/v1/port/location [get]
func (*_Scan) Location(ctx *gin.Context) {
	var (
		g = models.Gin{Ctx: ctx}
	)
	resp, err := management.ScanManager.Location()
	if err != nil {
		g.Fail(http.StatusBadRequest, err)
		return
	}
	g.Success(resp)
}
