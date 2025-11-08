package handlers

import (
	"sccsmsserver/db/pg"
	"sccsmsserver/i18n"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Get Dashboard data handler
func GetDashboardDataHandler(c *gin.Context) {
	d := new(pg.DashBoardData)
	err := c.ShouldBind(d)
	if err != nil {
		zap.L().Error("GetDashboardDataHandler invaid parms", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Operator ID
	operatorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, d)
		return
	}
	// Get Dashboard Data
	resStatus, _ = d.Get(operatorID)
	// Response
	ResponseWithMsg(c, resStatus, d)
}

// Get Risk Trends Data handler
func GetRiskTrendDataHandler(c *gin.Context) {
	rtd := new(pg.RiskTrendData)
	err := c.ShouldBind(rtd)
	if err != nil {
		zap.L().Error("GetRiskTrendDataHandler invaid parms", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Data
	resStatus, _ := rtd.Get()
	// Response
	ResponseWithMsg(c, resStatus, rtd)
}
