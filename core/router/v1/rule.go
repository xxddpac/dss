package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goportscan/common/utils"
	"goportscan/core/global"
	"goportscan/core/management"
	"goportscan/core/models"
	"net/http"
	"strings"
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
	switch body.Type {
	case global.Single:
		//192.168.1.1
		if !utils.ParseIP(body.TargetHost) {
			g.Fail(http.StatusBadRequest, fmt.Errorf("err target host with single type"))
			return
		}
	case global.Range:
		//192.168.1.10-20
		ipRange := strings.Split(body.TargetHost, "-")
		if len(ipRange) != 2 {
			g.Fail(http.StatusBadRequest, fmt.Errorf("err parse target host %v", body.TargetHost))
			return
		}
		start := ipRange[0]
		end := ipRange[1]
		if !utils.ParseIP(start) || start >= end {
			g.Fail(http.StatusBadRequest, fmt.Errorf("err target host with range type"))
			return
		}
	case global.Cidr:
		//192.168.1.0/20
		if !utils.ParseCidr(body.TargetHost) {
			g.Fail(http.StatusBadRequest, fmt.Errorf("err target host with cidr type"))
			return
		}
	}
	//1-4000
	portRange := strings.Split(body.TargetPort, "-")
	if len(portRange) != 2 {
		g.Fail(http.StatusBadRequest, fmt.Errorf("parse port range err"))
		return
	}
	start := utils.StrToInt(portRange[0])
	end := utils.StrToInt(portRange[1])
	if start < 1 || end > 65535 || start >= end {
		g.Fail(http.StatusBadRequest, fmt.Errorf("invalid port range value"))
		return
	}
	if err := management.RuleManager.Post(body); err != nil {
		g.Fail(http.StatusBadRequest, err)
		return
	}
	g.Success(nil)
}
