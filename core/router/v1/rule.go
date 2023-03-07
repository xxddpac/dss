package v1

import (
	"dss/common/utils"
	"dss/core/dao"
	"dss/core/global"
	"dss/core/management"
	"dss/core/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"net/http"
)

var Rule *_Rule

type _Rule struct {
}

// Post
// @Summary Add New Rule
// @Tags Rule
// @Accept  json
// @Produce  json
// @Param param body models.Rule true "Request Body"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /api/v1/rule [post]
func (*_Rule) Post(ctx *gin.Context) {
	var (
		g       = models.Gin{Ctx: ctx}
		body    models.Rule
		ipCount int
	)
	if err := ctx.ShouldBindJSON(&body); err != nil {
		g.Fail(http.StatusBadRequest, err)
		return
	}
	if dao.Repo(global.ScanRule).Count(bson.M{"target_host": body.TargetHost}) != 0 {
		g.Fail(http.StatusBadRequest, fmt.Errorf("target host %s already exists ", body.TargetHost))
		return
	}
	switch body.Type {
	case global.Single:
		//192.168.1.1
		ipCount = 1
		if !utils.ParseIP(body.TargetHost) {
			g.Fail(http.StatusBadRequest, fmt.Errorf("err target host with single type"))
			return
		}
	case global.Range:
		//192.168.1.10-20
		startIp, startIpEndSuffix, ipRangeEndSuffix, b := utils.ParseIpRange(body.TargetHost)
		if !b {
			g.Fail(http.StatusBadRequest, fmt.Errorf("err parse target host %v", body.TargetHost))
			return
		}
		if !utils.ParseIP(startIp) || startIpEndSuffix >= ipRangeEndSuffix {
			g.Fail(http.StatusBadRequest, fmt.Errorf("err target host with range type"))
			return
		}
		ipCount = ipRangeEndSuffix - startIpEndSuffix + 1
	case global.Cidr:
		//192.168.1.0/20
		if !utils.ParseCidr(body.TargetHost) {
			g.Fail(http.StatusBadRequest, fmt.Errorf("err target host with cidr type"))
			return
		}
		ipCount = len(utils.GetIpListByCidr(body.TargetHost))
	}
	//1-4000
	start, end, b := utils.ParsePortRange(body.TargetPort)
	if !b {
		g.Fail(http.StatusBadRequest, fmt.Errorf("parse port range err"))
		return
	}
	if start < 1 || end > 65535 || start >= end {
		g.Fail(http.StatusBadRequest, fmt.Errorf("invalid port range value"))
		return
	}
	portCount := end - start + 1
	body.Count = portCount * ipCount
	if err := management.RuleManager.Post(body); err != nil {
		g.Fail(http.StatusBadRequest, err)
		return
	}
	g.Success(nil)
}

// Delete
// @Summary Delete Rule By ID
// @Tags Rule
// @Accept  json
// @Produce  json
// @Param id query string true "Rule ID"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /api/v1/rule [delete]
func (*_Rule) Delete(ctx *gin.Context) {
	var (
		g     = models.Gin{Ctx: ctx}
		param models.QueryID
	)
	if err := ctx.ShouldBindQuery(&param); err != nil {
		g.Fail(http.StatusBadRequest, err)
		return
	}
	if err := management.RuleManager.Delete(param); err != nil {
		g.Fail(http.StatusBadRequest, err)
		return
	}
	g.Success(nil)
}

// Get
// @Summary Get Rule List
// @Tags Rule
// @Accept  json
// @Produce  json
// @Param search query string false "Fuzzy Query"
// @Param type query int false "Rule Type 1:Single 2:Range 3:CIDR"
// @Param status query string false "Rule Status true/false"
// @Param page query string false "Current Page Default:1"
// @Param size query string false "Current Size Default:10"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /api/v1/rule [get]
func (*_Rule) Get(ctx *gin.Context) {
	var (
		g     = models.Gin{Ctx: ctx}
		param = models.RuleQueryFunc()
	)
	if err := ctx.ShouldBindQuery(&param); err != nil {
		g.Fail(http.StatusBadRequest, err)
		return
	}
	resp, err := management.RuleManager.Get(*param)
	if err != nil {
		g.Fail(http.StatusBadRequest, err)
		return
	}
	g.Success(resp)
}

// Put
// @Summary Modify Rule By ID
// @Tags Rule
// @Accept  json
// @Produce  json
// @Param id query string true "Rule ID"
// @Param param body models.Rule true "Request Body"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /api/v1/rule [put]
func (*_Rule) Put(ctx *gin.Context) {
	var (
		g     = models.Gin{Ctx: ctx}
		param models.QueryID
		body  models.Rule
	)
	if err := ctx.ShouldBindQuery(&param); err != nil {
		g.Fail(http.StatusBadRequest, err)
		return
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		g.Fail(http.StatusBadRequest, err)
		return
	}
	if err := management.RuleManager.Put(param, body); err != nil {
		g.Fail(http.StatusBadRequest, err)
		return
	}
	g.Success(nil)
}

// Enum
// @Summary Display rule type enum
// @Tags Rule
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /api/v1/rule/type/enum [get]
func (*_Rule) Enum(ctx *gin.Context) {
	var (
		g = models.Gin{Ctx: ctx}
	)
	g.Success(management.RuleManager.Enum())
}
