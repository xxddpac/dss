package v1

import (
	"dss/core/management"
	"dss/core/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

var Task *_Task

type _Task struct {
}

// Post
// @Summary 执行扫描任务
// @Tags Task
// @Accept  json
// @Produce  json
// @Param id query string true "规则ID"
// @Param run_type query int true "任务类型枚举,1:计划任务调度执行,2:手动执行"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /api/v1/task [post]
func (*_Task) Post(ctx *gin.Context) {
	var (
		g     = models.Gin{Ctx: ctx}
		query models.TaskParam
	)
	if err := ctx.ShouldBindQuery(&query); err != nil {
		g.Fail(http.StatusBadRequest, err)
		return
	}
	if err := management.TaskManager.Post(query); err != nil {
		g.Fail(http.StatusBadRequest, err)
		return
	}
	g.Success(nil)
}

// Get
// @Summary Get 获取任务信息(执行时间、执行进度、执行状态、任务类型)
// @Tags Task
// @Accept  json
// @Produce  json
// @Param search query string false "模糊查询"
// @Param status query int false "任务状态枚举,1:等待中,2:检测中,3:已完成,4:出错"
// @Param run_type query int false "任务类型枚举,1:计划任务调度执行,2:手动执行"
// @Param page query string false "当前页数,默认值:1"
// @Param size query string false "当前条数,默认值:10"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /api/v1/task [get]
func (*_Task) Get(ctx *gin.Context) {
	var (
		g     = models.Gin{Ctx: ctx}
		param = models.TaskQueryFunc()
	)
	if err := ctx.ShouldBindQuery(&param); err != nil {
		g.Fail(http.StatusBadRequest, err)
		return
	}
	resp, err := management.TaskManager.Get(*param)
	if err != nil {
		g.Fail(http.StatusBadRequest, err)
		return
	}
	g.Success(resp)
}

// Enum
// @Summary 获取任务所有枚举信息
// @Tags Task
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /api/v1/task/enum [get]
func (*_Task) Enum(ctx *gin.Context) {
	var (
		g = models.Gin{Ctx: ctx}
	)
	g.Success(management.TaskManager.Enum())
}
