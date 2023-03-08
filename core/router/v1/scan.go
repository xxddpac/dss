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
// @Summary 获取扫描结果列表
// @Tags Scan
// @Accept  json
// @Produce  json
// @Param search query string false "模糊查询"
// @Param location query string false "根据规则定义的location筛选扫描结果"
// @Param date query string false "根据日期(year-month-day)筛选扫描结果"
// @Param page query string false "当前页数,默认值:1"
// @Param size query string false "当前条数,默认值:10"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /api/v1/scan [get]
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
// @Summary 作为计划任务每日执行,对比今日与昨日是否有新增端口开放,若有则通过企业微信发送提醒
// @Tags Scan
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /api/v1/scan/remind [get]
func (*_Scan) Remind(ctx *gin.Context) {
	var (
		g = models.Gin{Ctx: ctx}
	)
	go management.ScanManager.Remind()
	g.Success(nil)
}

// Clear
// @Summary 作为计划任务每日执行,只保留最近7天扫描结果
// @Tags Scan
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /api/v1/scan/clear [delete]
func (*_Scan) Clear(ctx *gin.Context) {
	var (
		g = models.Gin{Ctx: ctx}
	)
	go management.ScanManager.Clear()
	g.Success(nil)
}

// Trend
// @Summary 最近7天开放端口数量趋势
// @Tags Scan
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /api/v1/scan/trend [get]
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
// @Summary 获取规则提供location作为查询参数
// @Tags Scan
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /api/v1/scan/location [get]
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
